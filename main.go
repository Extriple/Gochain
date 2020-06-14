package main

import (
	"fmt"
	"strconv"

	"github.com/PatrykSer/Gochain/blockchain"
)

func main() {

	chain := blockchain.InitBlockchain()

	chain.AddBlocks("This is 1 block after Genesis")
	chain.AddBlocks("This is 2 block after Genesis")
	chain.AddBlocks("This is 3 block after Genesis")

	for _, block := range chain.Blocks {

		fmt.Printf("LastHash: %X\n", block.LastHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %X\n:", block.Hash)

		proofOfWork := blockchain.StatNewProof(block)
		fmt.Printf("ProofOfWork: %s\n", strconv.FormatBool(proofOfWork.ValidateRequirement()))
		fmt.Println()
	}

}
