package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"tpy-blockchain/internal/common"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

// Wallet represents a cryptocurrency wallet with associated data.
type Wallet struct {
	Mnemonic   string
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    string
	Balances   map[string]*big.Int             // Balances for tokens (e.g., "TPY": 0)
	Tokens     map[string]*common.UtilityToken // Tokens held by the wallet
}

func (w *Wallet) TransferTokens(tokenSymbol, receiver string, amount *big.Int) error {
	token, ok := w.Tokens[tokenSymbol]
	if !ok {
		return fmt.Errorf("token %s not found in wallet", tokenSymbol)
	}

	if w.Balances[tokenSymbol].Cmp(amount) < 0 {
		return fmt.Errorf("insufficient balance for token %s", tokenSymbol)
	}

	// Perform the transfer
	err := token.Transfer(w.Address, receiver, amount)
	if err != nil {
		return fmt.Errorf("transfer failed: %v", err)
	}

	// Update wallet's local balance
	w.Balances[tokenSymbol].Sub(w.Balances[tokenSymbol], amount)
	return nil
}

func (w *Wallet) VoteOnProposal(tokenSymbol, proposalID string, voteYes bool) error {
	token, ok := w.Tokens[tokenSymbol]
	if !ok {
		return fmt.Errorf("token %s not found in wallet", tokenSymbol)
	}

	votingPower := w.Balances[tokenSymbol]
	if votingPower.Sign() == 0 {
		return fmt.Errorf("no voting power for token %s", tokenSymbol)
	}

	// Cast the vote
	err := token.Vote(proposalID, w.Address, voteYes)
	if err != nil {
		return fmt.Errorf("failed to cast vote: %v", err)
	}

	return nil
}

func (w *Wallet) AddToken(token *common.UtilityToken) {
	w.Tokens[token.Symbol] = token
	if _, ok := w.Balances[token.Symbol]; !ok {
		w.Balances[token.Symbol] = big.NewInt(0) // Initialize balance if not already present
	}
}

func (w *Wallet) Sign(data []byte) (string, error) {
	// Create a hash of the data to sign
	hash := crypto.Keccak256Hash(data)

	// Sign the hash with the private key
	signature, err := crypto.Sign(hash.Bytes(), w.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %v", err)
	}

	// Convert the signature to a hex-encoded string
	return hex.EncodeToString(signature), nil
}

// LoadWordlist loads a 2048-word BIP-39 wordlist from a file.
func LoadWordlist() ([]string, error) {
	// Define the hardcoded path to the wordlist file
	filePath := "Seed-2048/English.txt"

	// Read the wordlist file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read wordlist file: %v", err)
	}

	// Split the file content into individual words
	words := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(words) != 2048 {
		return nil, errors.New("invalid wordlist: must contain exactly 2048 words")
	}

	return words, nil
}

func NewWallet() (*Wallet, error) {
	// Generate a private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	// Derive public key and address
	publicKey := &privateKey.PublicKey
	address := crypto.PubkeyToAddress(*publicKey).Hex()

	// Generate a mnemonic
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return nil, fmt.Errorf("failed to generate entropy: %v", err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf("failed to generate mnemonic: %v", err)
	}

	// Save the wallet data
	err = saveWalletToFile(privateKey, mnemonic, address)
	if err != nil {
		return nil, fmt.Errorf("failed to save wallet: %v", err)
	}

	// Return the wallet
	return &Wallet{
		Mnemonic:   mnemonic,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
		Balances:   make(map[string]*big.Int),
	}, nil
}

// saveWalletToFile saves the private key and mnemonic to a single file in the .wallets/ directory.
func saveWalletToFile(privateKey *ecdsa.PrivateKey, mnemonic, address string) error {
	// Ensure the .wallets directory exists
	walletsDir := "wallets"
	err := os.MkdirAll(walletsDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create wallets directory: %v", err)
	}

	// Convert the private key to PEM format
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %v", err)
	}

	pemBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privateKeyPem := pem.EncodeToMemory(pemBlock)

	// Prepare the file content, including the address
	fileContent := fmt.Sprintf(
		"Address:\n%s\n\nMnemonic:\n%s\n\nPrivate Key (PEM):\n%s",
		address,
		mnemonic,
		string(privateKeyPem),
	)

	// Save the wallet to a file
	filename := fmt.Sprintf("%s/%s_wallet.txt", walletsDir, address)
	err = os.WriteFile(filename, []byte(fileContent), 0600)
	if err != nil {
		return fmt.Errorf("failed to write wallet file: %v", err)
	}

	return nil
}

// LoadPrivateKeyFromFile loads a private key from a local PEM file.
func LoadPrivateKeyFromFile(filename string) (*ecdsa.PrivateKey, error) {
	// Read the PEM file
	pemData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %v", err)
	}

	// Decode the PEM block
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, errors.New("invalid private key PEM file")
	}

	// Parse the private key
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return privateKey, nil
}

func DecodePublicKey(pubKeyHex string) (*ecdsa.PublicKey, error) {
	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %v", err)
	}

	pubKey, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	switch pub := pubKey.(type) {
	case *ecdsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("decoded key is not an ECDSA public key")
	}
}

// RecoverWallet recreates a wallet using a mnemonic.
func RecoverWallet(mnemonic string) (*Wallet, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic")
	}

	seed := bip39.NewSeed(mnemonic, "")
	privateKey, err := crypto.ToECDSA(seed[:32])
	if err != nil {
		return nil, err
	}

	publicKey := &privateKey.PublicKey
	address := crypto.PubkeyToAddress(*publicKey).Hex()

	return &Wallet{
		Mnemonic:   mnemonic,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
		Balances:   make(map[string]*big.Int),
	}, nil
}

// VerifySignature validates a signature against the provided data and public key.
func VerifySignature(publicKey *ecdsa.PublicKey, signatureHex string, data []byte) bool {
	// Decode the hexadecimal signature
	sigBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false
	}

	// Hash the data
	hash := sha256.Sum256(data)

	// Split the signature into r and s values
	r := new(big.Int).SetBytes(sigBytes[:32])
	s := new(big.Int).SetBytes(sigBytes[32:])

	// Verify the signature
	return ecdsa.Verify(publicKey, hash[:], r, s)
}
