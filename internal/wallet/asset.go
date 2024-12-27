package wallet

import (
    "math/big"
)

// Asset represents a generic asset within a wallet
type Asset struct {
    Type    string
    Symbol  string
    Balance *big.Int
}

// NewAsset initializes a new asset with the given type and symbol
func NewAsset(assetType, symbol string) *Asset {
    return &Asset{
        Type:    assetType,
        Symbol:  symbol,
        Balance: big.NewInt(0),
    }
}
