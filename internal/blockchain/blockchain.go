package blockchain

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"tpy-blockchain/internal/common"
	"tpy-blockchain/internal/wallet"
)

type Blockchain struct {
	Blocks        []*Block                 `json:"blocks"`
	Wallets       map[string]*wallet.Wallet `json:"wallets"`
	Tokens        map[string]*common.UtilityToken `json:"tokens"`
	Balances      map[string]*big.Int      `json:"balances"`      // Add Balances
	Transactions  []*Transaction           `json:"transactions"` // Add Transactions
	mutex         sync.Mutex
	blockDir      string
	blockLimit    int
}

func NewBlockchain() *Blockchain {
	blockDir := "Blocks"

	// Ensure the Blocks directory exists
	err := os.MkdirAll(blockDir, 0755)
	if err != nil {
		panic(fmt.Sprintf("Failed to create block directory: %v", err))
	}

	// Check for existing chain files
	files, err := os.ReadDir(blockDir)
	if err != nil {
		panic(fmt.Sprintf("Failed to read block directory: %v", err))
	}

	var highestChainIndex int
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "chain") && strings.HasSuffix(file.Name(), ".json") {
			chainNumber, err := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(file.Name(), "chain"), ".json"))
			if err == nil && chainNumber > highestChainIndex {
				highestChainIndex = chainNumber
			}
		}
	}

	// Load blockchain if chain files are found
	if highestChainIndex > 0 {
		fmt.Printf("Found existing blockchain files up to chain%d.json. Loading...\n", highestChainIndex)
		bc, err := LoadBlockchainFromFiles(blockDir, highestChainIndex)
		if err != nil {
			panic(fmt.Sprintf("Failed to load blockchain: %v", err))
		}
		return bc
	}

	// Create a new blockchain if no chain files are found
	genesisBlock := &Block{
		Index:        0,
		Timestamp:    time.Now().String(),
		Transactions: []*Transaction{},
		Wallets:      make(map[string]*wallet.Wallet), // Initialize Wallets
		Tokens:       make(map[string]*common.UtilityToken), // Initialize Tokens
		Nonce:        0,
		PreviousHash: "0",
	}	
	genesisBlock.Hash = CalculateHash(genesisBlock)

	bc := &Blockchain{
		Blocks:     []*Block{genesisBlock},
		Balances:   make(map[string]*big.Int),
		Tokens:     make(map[string]*common.UtilityToken),
		Wallets:    make(map[string]*wallet.Wallet), // Initialize Wallets
		blockDir:   blockDir,
		blockLimit: 1000,
	}

	// Save the genesis block
	if err := bc.SaveBlocksToFile(); err != nil {
		panic(fmt.Sprintf("Failed to save genesis block: %v", err))
	}

	fmt.Println("Genesis block created and saved as chain1.json.")
	return bc
}

func (bc *Blockchain) AddToken(name, symbol string, totalSupply *big.Int, decimals uint) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	// Check if the token already exists
	if _, exists := bc.Tokens[symbol]; exists {
		return fmt.Errorf("token with symbol %s already exists", symbol)
	}

	// Create and add the token
	token := &common.UtilityToken{
		Name:        name,
		Symbol:      symbol,
		TotalSupply: totalSupply,
		Decimals:    decimals,
		Balances:    make(map[string]*big.Int),
		VotingPower: make(map[string]*big.Int),
		Proposals:   []*common.Proposal{},
		Address:     generateUniqueTokenAddress(symbol),
	}

	bc.Tokens[symbol] = token
	return nil
}

func LoadBlockchainFromFiles(blockDir string, highestChainIndex int) (*Blockchain, error) {
	bc := &Blockchain{
		Blocks:       []*Block{},
		Wallets:      make(map[string]*wallet.Wallet),
		Tokens:       make(map[string]*common.UtilityToken),
		Balances:     make(map[string]*big.Int),
		Transactions: []*Transaction{},
		blockDir:     blockDir,
		blockLimit:   1000,
	}

	for i := 1; i <= highestChainIndex; i++ {
		filename := filepath.Join(blockDir, fmt.Sprintf("chain%d.json", i))

		// Read the file
		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %v", filename, err)
		}

		// Parse the JSON data
		var loadedData map[string]interface{}
		if err := json.Unmarshal(data, &loadedData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal data from file %s: %v", filename, err)
		}

		// Load Blocks
		if blocks, ok := loadedData["blocks"].([]interface{}); ok {
			for _, blockData := range blocks {
				blockBytes, _ := json.Marshal(blockData)
				var block Block
				if err := json.Unmarshal(blockBytes, &block); err != nil {
					return nil, fmt.Errorf("failed to unmarshal block: %v", err)
				}
				bc.Blocks = append(bc.Blocks, &block)
			}
		}

		// Load Wallets
		if wallets, ok := loadedData["wallets"].(map[string]interface{}); ok {
			for addr, walletData := range wallets {
				walletBytes, _ := json.Marshal(walletData)
				var w wallet.Wallet
				if err := json.Unmarshal(walletBytes, &w); err != nil {
					return nil, fmt.Errorf("failed to unmarshal wallet: %v", err)
				}
				bc.Wallets[addr] = &w
			}
		}

		// Load Tokens
		if tokens, ok := loadedData["tokens"].(map[string]interface{}); ok {
			for symbol, tokenData := range tokens {
				tokenBytes, _ := json.Marshal(tokenData)
				var token common.UtilityToken
				if err := json.Unmarshal(tokenBytes, &token); err != nil {
					return nil, fmt.Errorf("failed to unmarshal token: %v", err)
				}
				bc.Tokens[symbol] = &token
			}
		}

		// Load Balances
		if balances, ok := loadedData["balances"].(map[string]interface{}); ok {
			for addr, balance := range balances {
				if balanceStr, ok := balance.(string); ok {
					bal := new(big.Int)
					bal.SetString(balanceStr, 10)
					bc.Balances[addr] = bal
				}
			}
		}

		// Load Transactions
		if transactions, ok := loadedData["transactions"].([]interface{}); ok {
			for _, txData := range transactions {
				txBytes, _ := json.Marshal(txData)
				var tx Transaction
				if err := json.Unmarshal(txBytes, &tx); err != nil {
					return nil, fmt.Errorf("failed to unmarshal transaction: %v", err)
				}
				bc.Transactions = append(bc.Transactions, &tx)
			}
		}
	}

	fmt.Println("Blockchain successfully loaded from chain files.")
	return bc, nil
}

// Helper function to generate a unique token address
func generateUniqueTokenAddress(symbol string) string {
	return fmt.Sprintf("0x%s", symbol[:3])
}

func (bc *Blockchain) GetWallet(address string) (*wallet.Wallet, error) {
	balance, exists := bc.Balances[address]
	if !exists {
		return nil, fmt.Errorf("wallet with address %s not found", address)
	}
	return &wallet.Wallet{
		Address: address,
		Balances: map[string]*big.Int{
			"default": balance,
		},
	}, nil
}

func (bc *Blockchain) AddTransaction(transaction *Transaction) error {
	if transaction.Amount.Sign() <= 0 {
		return fmt.Errorf("transaction amount must be positive")
	}

	senderBalance, senderExists := bc.Balances[transaction.Sender]
	if !senderExists || senderBalance.Cmp(transaction.Amount) < 0 {
		return fmt.Errorf("insufficient balance for sender %s", transaction.Sender)
	}

	bc.Balances[transaction.Sender].Sub(senderBalance, transaction.Amount)
	receiverBalance, receiverExists := bc.Balances[transaction.Receiver]
	if !receiverExists {
		bc.Balances[transaction.Receiver] = big.NewInt(0)
		receiverBalance = bc.Balances[transaction.Receiver]
	}
	receiverBalance.Add(receiverBalance, transaction.Amount)

	bc.Transactions = append(bc.Transactions, transaction)
	return nil
}

func (bc *Blockchain) GetBlockByIndex(index int) (*Block, error) {
	if index < 0 || index >= len(bc.Blocks) {
		return nil, fmt.Errorf("block with index %d not found", index)
	}
	return bc.Blocks[index], nil
}

func (bc *Blockchain) GetBlocks() []*Block {
	return bc.Blocks
}

func (bc *Blockchain) GetBalance(address string) (*big.Int, error) {
	balance, exists := bc.Balances[address]
	if !exists {
		return nil, fmt.Errorf("address %s not found", address)
	}
	return balance, nil
}

func (bc *Blockchain) AddBlock(transactions []*Transaction) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	previousBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := &Block{
		Index:        len(bc.Blocks),
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		Nonce:        0,
		PreviousHash: previousBlock.Hash,
	}
	newBlock.Hash = CalculateHash(newBlock)
	bc.Blocks = append(bc.Blocks, newBlock)
	return nil
}

func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		previousBlock := bc.Blocks[i-1]
		if currentBlock.PreviousHash != previousBlock.Hash || currentBlock.Hash != CalculateHash(currentBlock) {
			return false
		}
	}
	return true
}

func (bc *Blockchain) SaveBlocksToFile() error {
	startIndex := len(bc.Blocks) - bc.blockLimit
	if startIndex < 0 {
		startIndex = 0
	}
	blocksToSave := bc.Blocks[startIndex:]

	dataToSave := map[string]interface{}{
		"blocks":  blocksToSave,
		"tokens":  bc.Tokens,   // Serialize tokens globally
		"wallets": bc.Wallets,  // Serialize wallets globally
	}

	fileIndex := (len(bc.Blocks)-1)/bc.blockLimit + 1
	filename := filepath.Join(bc.blockDir, fmt.Sprintf("chain%d.json", fileIndex))

	data, err := json.MarshalIndent(dataToSave, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal blockchain data: %v", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write to file %s: %v", filename, err)
	}
	return nil
}

func (bc *Blockchain) LoadBlocks() error {
	files, err := os.ReadDir(bc.blockDir)
	if err != nil {
		return fmt.Errorf("failed to read block directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filename := filepath.Join(bc.blockDir, file.Name())

			// Read the file
			data, err := os.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %v", filename, err)
			}

			// Unmarshal the blocks, tokens, and wallets
			var dataToLoad map[string]interface{}
			if err := json.Unmarshal(data, &dataToLoad); err != nil {
				return fmt.Errorf("failed to unmarshal data from file %s: %v", filename, err)
			}

			// Load blocks
			var loadedBlocks []*Block
			blocksData, _ := dataToLoad["blocks"].([]interface{})
			for _, block := range blocksData {
				blockBytes, _ := json.Marshal(block)
				var b Block
				json.Unmarshal(blockBytes, &b)
				loadedBlocks = append(loadedBlocks, &b)
			}
			bc.Blocks = append(bc.Blocks, loadedBlocks...)

			// Load tokens
			tokensData, _ := dataToLoad["tokens"].(map[string]interface{})
			for symbol, tokenData := range tokensData {
				tokenBytes, _ := json.Marshal(tokenData)
				var token common.UtilityToken
				json.Unmarshal(tokenBytes, &token)
				bc.Tokens[symbol] = &token
			}

			// Load wallets
			walletsData, _ := dataToLoad["wallets"].(map[string]interface{})
			for address, walletData := range walletsData {
				walletBytes, _ := json.Marshal(walletData)
				var w wallet.Wallet
				json.Unmarshal(walletBytes, &w)
				bc.Wallets[address] = &w
			}
		}
	}
	return nil
}
