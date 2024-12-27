package consensus

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// ProofOfWork represents the proof-of-work system
type ProofOfWork struct {
	Difficulty int // Number of leading zeros required in the hash
}

// NewProofOfWork initializes a new PoW system
func NewProofOfWork(difficulty int) *ProofOfWork {
	return &ProofOfWork{
		Difficulty: difficulty,
	}
}

// Mine performs mining by solving the PoW challenge
func (pow *ProofOfWork) Mine(index int, timestamp string, transactions string, previousHash string) (string, int) {
	nonce := 0
	var hash string
	target := strings.Repeat("0", pow.Difficulty)

	for {
		data := fmt.Sprintf("%d%s%s%s%d", index, timestamp, transactions, previousHash, nonce)
		hash = calculateHash(data)

		// Check if the hash satisfies the difficulty
		if strings.HasPrefix(hash, target) {
			break
		}
		nonce++
	}

	fmt.Printf("Block mined! Hash: %s, Nonce: %d\n", hash, nonce)
	return hash, nonce
}

// ValidateProof checks if a given hash satisfies the PoW difficulty
func (pow *ProofOfWork) ValidateProof(hash string) bool {
	target := strings.Repeat("0", pow.Difficulty)
	return strings.HasPrefix(hash, target)
}

// calculateHash generates a SHA256 hash for the given data
func calculateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
