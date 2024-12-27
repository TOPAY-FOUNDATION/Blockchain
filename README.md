# TOPAY BLOCKCHAIN

This project implements a blockchain system with wallets, transactions, and governance mechanisms in Go. The system is entirely terminal-based, providing a Command-Line Interface (CLI) for interacting with the blockchain.

---

## **Features**

- **Blockchain Management**: View the blockchain and its blocks.
- **Wallet Creation**: Generate wallets with recovery mnemonics and private keys.
- **Token Management**: Built-in utility token (`TPY`) for governance and transfers.
- **Token Transfers**: Transfer tokens between wallets.
- **Transaction Handling**: (In Progress) Track and verify transactions in the blockchain.
- **Persistent Storage**: Blockchain data is saved to JSON files for persistence.

---

## **Prerequisites**

Ensure the following are installed on your system:

- [Go](https://golang.org/doc/install) (version 1.19 or above)
- [Git](https://git-scm.com/downloads) (for cloning the repository)

---

## **Run Instructions**

### **1. Clone the Repository**

Clone the repository and navigate into the project directory:

```bash
gh repo clone TOPAY-FOUNDATION/Blockchain
cd tpy-blockchain
```

### **2. Install Dependencies**

Run the following command to download and install all dependencies specified in `go.mod`:

```bash
go mod tidy
```

### **3. Run the Application**

Start the CLI application using the `main.go` entry point:

```bash
go run cmd/main.go
```

---

## **Using the CLI**

The CLI provides the following options:

```
--- TOPAY Blockchain CLI ---
1. Create Wallet
2. View Blockchain
3. Add Transaction
4. View Wallet Balance
5. Transfer Tokens
6. Exit
```

### **1. Create Wallet**
Generates a new wallet with:
- **Address**: A unique wallet address.
- **Mnemonic**: A recovery phrase for restoring the wallet.
- **Initial Balance**: Wallets are assigned an initial balance of 1000 `TPY`.

### **2. View Blockchain**
Displays all blocks in the blockchain, including:
- Block Index
- Block Hash
- Previous Block Hash
- Number of Transactions

### **3. Add Transaction**
(Add functionality in progress) Allows users to add custom transactions to the blockchain.

### **4. View Wallet Balance**
Displays the balance of a wallet by its address.

### **5. Transfer Tokens**
Transfers governance tokens (`TPY`) between wallets:
- **Sender Address**: The address sending the tokens.
- **Receiver Address**: The address receiving the tokens.
- **Amount**: The number of tokens to transfer.

---

## **Data Storage**

### `.Blocks/`
All blockchain data is stored in this directory. Files include:
- **`chain1.json`**: Contains the genesis block and subsequent blocks.
- **`chain2.json`**, etc.: Created when block limits are exceeded.

Each block includes:
- **Index**: Position in the chain.
- **Transactions**: List of transactions in the block.
- **Wallets**: Wallet data associated with the block.
- **Tokens**: Token balances and metadata.
- **Hash and Previous Hash**: Ensures integrity of the blockchain.

---

## **Troubleshooting**

1. **Dependencies Issue**:
   - Run `go mod tidy` to ensure all dependencies are installed.

2. **Wallet Not Saved**:
   - Verify that the `.Blocks` directory exists and is writable.

3. **Transaction Errors**:
   - Ensure wallet addresses and token balances are valid before performing transfers.

---

## **Planned Enhancements**

- **Transaction System**: Add verification and tracking of blockchain transactions.
- **Governance Features**: Enable token-based voting and proposals.
- **Improved Storage**: Implement database support for scalability.

---

Let us know if anything is missing! 🚀