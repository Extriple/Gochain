package commandLine

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/PatrykSer/Gochain/blockchain"
	"github.com/PatrykSer/Gochain/digitalWallet"
)

//Struktura wiersza poleceń
type LineStruct struct{}

//Metoda wyświetlająca nasze polecenia
func (cli *LineStruct) printCommandLine() {
	fmt.Println("Usage:")
	fmt.Println(" getbalance -address ADDRESS - get the balance for an address")
	fmt.Println(" createblockchain -address ADDRESS creates a blockchain and sends genesis reward to address")
	fmt.Println(" printchain - Prints the blocks in the chain")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT - Send amount of coins")
	fmt.Println("createWallet - Create new Wallet")
	fmt.Println("listaddress - List address in our wallet file")
}

//Metoda sprawdzająca poprawność argumentów
func (cli *LineStruct) validateArguments() {
	if len(os.Args) < 2 {
		cli.printCommandLine()
		runtime.Goexit()
	}
}

//Metoda generująca/pokazująca listę posiadanych adresów
func (cli *LineStruct) AddressList() {
	wallets, _ := digitalWallet.CreateNewWallet()
	allAddress := wallets.GetAllNewAddress()

	for _, address := range allAddress {
		fmt.Println(address)
	}
}

//Metoda tworząca nowt portfel
func (cli *LineStruct) createWallet() {
	wallets, _ := digitalWallet.CreateNewWallet()
	address := wallets.AddNewWallet()
	wallets.Saving()

	fmt.Printf("New address is : %s\n", address)
}

//Metoda pokazująca wartość blockchaina
func (cli *LineStruct) ShowChain() {
	chain := blockchain.ContinueBlockchain("")
	defer chain.Database.Close()
	iterator := chain.Iterator()

	for {

		block := iterator.Next()
		fmt.Printf("LastHash: %X\n", block.LastHash)
		fmt.Printf("Hash in Block: %s\n", block.Hash)

		proofOfWork := blockchain.StatNewProof(block)
		fmt.Printf("ProofOfWork: %s\n", strconv.FormatBool(proofOfWork.ValidateRequirement()))
		for _, trans := range block.Transaction {
			fmt.Println(trans)
		}
		fmt.Println()
		//Przerwanie pętli jezeli długość bloków jest równa 0
		if len(block.LastHash) == 0 {
			break
		}

	}

}

//Tworzymy nowy łańcuch bloków
func (cli *LineStruct) CreateBlockchain(address string) {
	if !digitalWallet.ValidateAddress(address) { 
		log.Panic("Address is not valid")
	}

	chain := blockchain.InitBlockchain(address)
	chain.Database.Close()
	fmt.Println("Finished")
}

//Metoda Balance
func (cli *LineStruct) Balance(address string) {
	if !digitalWallet.ValidateAddress(address) {
		log.Panic("Address is not Valid")
	}

	chain := blockchain.ContinueBlockchain(address)
	defer chain.Database.Close()

	balance := 0
	pubKeyHash := digitalWallet.DecodeBase58([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash) - 4]
	UTXOs := chain.Find(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

//Metoda wysyłająca nową transakcję
func (cli *LineStruct) Sending(from, to string, amount int) {
	if !digitalWallet.ValidateAddress(to) {
		log.Panic("Address is not Valid")
	}

	if !digitalWallet.ValidateAddress(from) {
		log.Panic("Address is not valid ")
	}
	chain := blockchain.ContinueBlockchain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlocks([]*blockchain.Transaction{tx})
	fmt.Println("Success!")
}

//Wywołuejmy wszystkie metody
func (cli *LineStruct) running() {
	cli.validateArguments()

	//Flagi
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createWallet", flag.ExitOnError)
	listAddressCmd := flag.NewFlagSet("listAddress", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	//Instrukcje swith
	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createWallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listAddress":
		err := listAddressCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printCommandLine()
		runtime.Goexit()
	}

	//Flagi
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.Balance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.CreateBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.ShowChain()
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}
	
	if listAddressCmd.Parsed() {
		cli.AddressList()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.Sending(*sendFrom, *sendTo, *sendAmount)
	}
}
