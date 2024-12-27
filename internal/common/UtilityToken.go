package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/google/uuid"
)

type UtilityToken struct {
    Name         string
    Symbol       string
    TotalSupply  *big.Int
    Decimals     uint
    Balances     map[string]*big.Int
    VotingPower  map[string]*big.Int
    Proposals    []*Proposal
    Address      string
}

type Proposal struct {
    ID          string
    Title       string
    Description string
    Votes       map[string]*big.Int
    YesVotes    *big.Int
    NoVotes     *big.Int
    Status      string
}

// Use common.UtilityToken instead of defining it locally
func NewUtilityToken() *UtilityToken {
    name := "TOPAY"
    symbol := "TPY"
    decimals := uint(18)
    supply := uint64(120000000)
    totalSupply := new(big.Int).Mul(big.NewInt(int64(supply)), new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))

    address := generateUniqueAddress()

    return &UtilityToken{
        Name:        name,
        Symbol:      symbol,
        TotalSupply: totalSupply,
        Decimals:    decimals,
        Balances:    make(map[string]*big.Int),
        VotingPower: make(map[string]*big.Int),
        Proposals:   []*Proposal{},
        Address:     address,
    }
}

func generateUniqueAddress() string {
    id := uuid.New().String() // Generate a UUID
    hash := sha256.Sum256([]byte(id))
    return hex.EncodeToString(hash[:])[:40] // Use first 40 hex characters
}


// Transfer transfers tokens from the sender to the receiver.
func (token *UtilityToken) Transfer(sender, receiver string, amount *big.Int) error {
    // Check if the sender has a balance
    if _, ok := token.Balances[sender]; !ok {
        return fmt.Errorf("sender address not found")
    }
    // Ensure the sender has enough balance
    if token.Balances[sender].Cmp(amount) < 0 {
        return fmt.Errorf("insufficient balance")
    }

    // Deduct from sender's balance
    token.Balances[sender].Sub(token.Balances[sender], amount)

    // Add to receiver's balance
    if _, ok := token.Balances[receiver]; !ok {
        token.Balances[receiver] = new(big.Int)
    }
    token.Balances[receiver].Add(token.Balances[receiver], amount)

    return nil
}

func (token *UtilityToken) AddProposal(title, description string) *Proposal {
    proposal := &Proposal{
        ID:          generateProposalID(),
        Title:       title,
        Description: description,
        Votes:       make(map[string]*big.Int),
        YesVotes:    big.NewInt(0),
        NoVotes:     big.NewInt(0),
        Status:      "Active",
    }
    token.Proposals = append(token.Proposals, proposal)
    return proposal
}

func generateProposalID() string {
    id := uuid.New().String()
    return id
}

func (token *UtilityToken) Vote(proposalID, voterAddress string, voteYes bool) error {
    var proposal *Proposal
    for _, p := range token.Proposals {
        if p.ID == proposalID {
            proposal = p
            break
        }
    }
    if proposal == nil {
        return fmt.Errorf("proposal not found")
    }
    if proposal.Status != "Active" {
        return fmt.Errorf("proposal is not active")
    }

    votingPower := token.VotingPower[voterAddress]
    if votingPower == nil || votingPower.Sign() == 0 {
        return fmt.Errorf("no voting power")
    }

    if voteYes {
        proposal.YesVotes.Add(proposal.YesVotes, votingPower)
    } else {
        proposal.NoVotes.Add(proposal.NoVotes, votingPower)
    }

    proposal.Votes[voterAddress] = new(big.Int).Set(votingPower)
    return nil
}

func (token *UtilityToken) CloseProposal(proposalID string) error {
    var proposal *Proposal
    for _, p := range token.Proposals {
        if p.ID == proposalID {
            proposal = p
            break
        }
    }
    if proposal == nil {
        return fmt.Errorf("proposal not found")
    }
    if proposal.Status != "Active" {
        return fmt.Errorf("proposal is already closed")
    }

    if proposal.YesVotes.Cmp(proposal.NoVotes) > 0 {
        proposal.Status = "Passed"
    } else {
        proposal.Status = "Rejected"
    }
    return nil
}
