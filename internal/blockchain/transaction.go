package blockchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"tpy-blockchain/internal/wallet"

	"github.com/ethereum/go-ethereum/crypto"
)

type Transaction struct {
    Sender      string
    Receiver    string
    Amount      *big.Int
    Signature   string
    Hash        string
    TokenSymbol string
}

// NewTransaction creates a new transaction, validates balances, and ensures token compatibility.
func NewTransaction(senderWallet *wallet.Wallet, receiver string, amount *big.Int, tokenSymbol string) (*Transaction, error) {
    // Validate the sender's token balance
    if balance, ok := senderWallet.Balances[tokenSymbol]; !ok || balance.Cmp(amount) < 0 {
        return nil, fmt.Errorf("insufficient balance for token %s", tokenSymbol)
    }

    // Create and populate the transaction
    tx := &Transaction{
        Sender:      senderWallet.Address,
        Receiver:    receiver,
        Amount:      amount,
        TokenSymbol: tokenSymbol,
    }

    // Generate the transaction hash
    tx.Hash = tx.calculateHash()

    // Sign the transaction using the sender's wallet
    signature, err := senderWallet.Sign([]byte(tx.Hash))
    if err != nil {
        return nil, fmt.Errorf("failed to sign transaction: %v", err)
    }
    tx.Signature = signature

    return tx, nil
}

func (tx *Transaction) calculateHash() string {
    record := fmt.Sprintf("%s:%s:%s:%s", tx.Sender, tx.Receiver, tx.Amount.String(), tx.TokenSymbol)
    h := sha256.New()
    h.Write([]byte(record))
    return hex.EncodeToString(h.Sum(nil))
}

func VerifySignature(publicKey *ecdsa.PublicKey, signatureHex string, data []byte) bool {
    sigBytes, err := hex.DecodeString(signatureHex)
    if err != nil {
        return false
    }

    hash := crypto.Keccak256Hash(data)

    r := new(big.Int).SetBytes(sigBytes[:32])
    s := new(big.Int).SetBytes(sigBytes[32:])
    return ecdsa.Verify(publicKey, hash.Bytes(), r, s)
}
