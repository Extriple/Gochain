package main

import (
	"os"
	

	"github.com/PatrykSer/Gochain/digitalWallet"
)

func main() {
	defer os.Exit(0)
	//cmd :=commandLine.LineStruct {}
	//cmd.running()
	w := digitalWallet.BuildWallet()
	w.Address()

}