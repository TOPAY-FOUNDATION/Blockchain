package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// SignMessage signs a message using a private key
func SignMessage(privateKey *ecdsa.PrivateKey, message string) (string, error) {
	hash := sha256.Sum256([]byte(message))
	signature, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %v", err)
	}
	return hex.EncodeToString(signature), nil
}

// VerifySignature verifies a signature using a public key
func VerifySignature(publicKey *ecdsa.PublicKey, message, signatureHex string) bool {
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false
	}

	hash := sha256.Sum256([]byte(message))
	return ecdsa.VerifyASN1(publicKey, hash[:], signature)
}
