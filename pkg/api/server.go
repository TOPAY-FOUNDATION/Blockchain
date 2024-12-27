package api

import (
    "tpy-blockchain/internal/blockchain"
    "github.com/gin-gonic/gin"
)

// StartServer initializes and starts the Gin server on the specified port.
func StartServer(port string, chain *blockchain.Blockchain) {
    router := gin.Default()
    RegisterRoutes(router, chain)
    router.Run(":" + port) // Starts the HTTP server
}
