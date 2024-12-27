package api

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"tpy-blockchain/internal/blockchain"
	"tpy-blockchain/internal/wallet"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes setups all the routes for the API server
func RegisterRoutes(router *gin.Engine, chain *blockchain.Blockchain) {
	router.GET("/blocks", getBlocksHandler(chain))
	router.POST("/wallets/new", createWalletHandler())
	router.POST("/wallets/import", importWalletHandler())
	router.GET("/wallets/balance", getWalletBalanceHandler(chain))
    router.POST("/wallets/transaction", createWalletTransactionHandler(chain))

}

func getWalletBalanceHandler(chain *blockchain.Blockchain) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the wallet address from query parameters
		address := c.Query("address")
		if address == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Address is required"})
			return
		}

		// Get the balance for the specified address
		balance, err := chain.GetBalance(address)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		// Return the balance as a JSON response
		c.JSON(http.StatusOK, gin.H{
			"address": address,
			"balance": balance.String(), // Convert *big.Int to string for JSON
		})
	}
}

func getBlocksHandler(chain *blockchain.Blockchain) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"blocks": chain.Blocks})
	}
}

func createWalletHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		w, err := wallet.NewWallet()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"address":   w.Address,
			"mnemonic":  w.Mnemonic,
			"publicKey": hex.EncodeToString(crypto.FromECDSAPub(w.PublicKey)),
		})
	}
}

func importWalletHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Mnemonic string `json:"mnemonic"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}
		w, err := wallet.RecoverWallet(req.Mnemonic)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to recover wallet: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"address":   w.Address,
			"mnemonic":  req.Mnemonic,
			"publicKey": hex.EncodeToString(crypto.FromECDSAPub(w.PublicKey)),
		})
	}
}

// Handler for creating wallet-to-wallet transactions
func createWalletTransactionHandler(chain *blockchain.Blockchain) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Sender      string  `json:"sender"`
			Receiver    string  `json:"receiver"`
			Amount      float64 `json:"amount"`
			TokenSymbol string  `json:"token_symbol"`
		}

		// Parse the incoming JSON request
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request data: %v", err)})
			return
		}

		// Fetch the sender's wallet
		senderWallet, err := chain.GetWallet(req.Sender)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Sender wallet not found: %v", err)})
			return
		}

		// Convert amount to *big.Int
		amountInt := big.NewInt(int64(req.Amount))

		// Create a new transaction
		transaction, err := blockchain.NewTransaction(senderWallet, req.Receiver, amountInt, req.TokenSymbol)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Transaction creation failed: %v", err)})
			return
		}

		// Add the transaction to the blockchain
		if err := chain.AddTransaction(transaction); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add transaction: %v", err)})
			return
		}

		// Respond with success
		c.JSON(http.StatusOK, gin.H{
			"message":     "Transaction successfully processed",
			"transaction": transaction,
		})
	}
}