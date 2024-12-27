package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
	"tpy-blockchain/internal/common"
	"tpy-blockchain/internal/wallet"
)

type Block struct {
	Index         int                          `json:"index"`
	Timestamp     string                       `json:"timestamp"`
	Transactions  []*Transaction               `json:"transactions"`
	Wallets       map[string]*wallet.Wallet    `json:"wallets"` // Include Wallets
	Tokens        map[string]*common.UtilityToken `json:"tokens"` // Include Tokens
	Nonce         int                          `json:"nonce"`
	PreviousHash  string                       `json:"previousHash"`
	Hash          string                       `json:"hash"`
}

// NewBlock initializes a new block with the given parameters
func NewBlock(index int, previousHash string, transactions []*Transaction) *Block {
	block := &Block{
		Index:        index,
		Timestamp:    time.Now().Format(time.RFC3339),
		Transactions: transactions,
		Wallets:      make(map[string]*wallet.Wallet),
		Tokens:       make(map[string]*common.UtilityToken),
		PreviousHash: previousHash,
		Nonce:        0,
	}
	block.Hash = CalculateHash(block)
	return block
}

// AddTransaction adds a transaction to the block
func (block *Block) AddTransaction(tx *Transaction) error {
	if err := ValidateTransaction(tx); err != nil {
		return fmt.Errorf("transaction validation failed: %v", err)
	}
	block.Transactions = append(block.Transactions, tx)
	return nil
}

// AddWallet adds a wallet to the block
func (block *Block) AddWallet(w *wallet.Wallet) {
	block.Wallets[w.Address] = w
}

// AddToken adds a token to the block
func (block *Block) AddToken(token *common.UtilityToken) {
	block.Tokens[token.Symbol] = token
}

// CalculateHash computes a hash of the block
func CalculateHash(block *Block) string {
	record := fmt.Sprintf("%08d", block.Index) + block.Timestamp + fmt.Sprintf("%08d", block.Nonce) + block.PreviousHash
	for _, tx := range block.Transactions {
		record += tx.Hash
	}
	h := sha256.New()
	h.Write([]byte(record))
	return hex.EncodeToString(h.Sum(nil))
}

// MineBlock performs proof-of-work mining for the block
func (block *Block) MineBlock(difficulty int) {
	target := ""
	for i := 0; i < difficulty; i++ {
		target += "0"
	}

	for block.Hash[:difficulty] != target {
		block.Nonce++
		block.Hash = CalculateHash(block)
	}
	fmt.Printf("Block mined with nonce %d: %s\n", block.Nonce, block.Hash)
}

// ValidateTransactions checks the validity of all transactions in the block
func (block *Block) ValidateTransactions() error {
	for _, tx := range block.Transactions {
		if !VerifyTransaction(tx) {
			return fmt.Errorf("invalid transaction: %v", tx)
		}
	}
	return nil
}

// VerifyTransaction validates a transaction's signature, hash, and data
func VerifyTransaction(tx *Transaction) bool {
	publicKey, err := wallet.DecodePublicKey(tx.Sender)
	if err != nil {
		fmt.Printf("Failed to decode public key: %v\n", err)
		return false
	}

	// Verify the signature
	isValid := wallet.VerifySignature(publicKey, tx.Signature, []byte(tx.Hash))
	if !isValid {
		fmt.Println("Signature verification failed for transaction:", tx)
		return false
	}

	// Ensure the hash matches the calculated hash
	calculatedHash := tx.calculateHash()
	if tx.Hash != calculatedHash {
		fmt.Println("Transaction hash mismatch:", tx)
		return false
	}

	return true
}

// ValidateTransaction checks a single transaction
func ValidateTransaction(tx *Transaction) error {
	if !VerifyTransaction(tx) {
		return fmt.Errorf("transaction is invalid")
	}
	if tx.Amount.Sign() <= 0 {
		return fmt.Errorf("transaction amount must be positive")
	}
	return nil
}
