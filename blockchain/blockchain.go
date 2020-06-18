package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/dgraph-io/badger"
)

//scieżka do bazy danych
const (
	dbPath = "/tmp/blocks"
	//Stała z plikiem manifest
	dbConstFile = "/tmp/blocks/MANIFEST"
	genesisData = "First transaction from Genesis"
)

//Struktura blockchaina
//ps. bloki będą przechowywały w swojej pamieci ostatni hash i zostanie to zapisane w bazie danych bagder
type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

//Struktura iteracji
type BlockchainIterator struct {
	NowHash  []byte
	Database *badger.DB
}

//Funkcja, która pozwala ustalić czy istnieje nasza baza danych
func ExistDB() bool {
	if _, err := os.Stat(dbConstFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func ContinueBlockchain(address string) *Blockchain {
	//Baza danych nie istnieje
	if ExistDB() == false {
		fmt.Println("No existing blockchain found ")
		runtime.Goexit()
	}

	var lastHash []byte
	opts := badger.DefaultOptions(path.Dir("/tmp/blocks"))
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	//Tworzymy baze danych
	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lastHash"))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)

		return err
	})

	Handle(err)
	chain := Blockchain{lastHash, db}

	return &chain
}

//Metoda, która zespala wszystkie wyżej wymienione bloki i tworzy naszego blockchaina
func InitBlockchain(address string) *Blockchain {
	var lastHash []byte

	//Sprawdzanie czy mamy DB
	if ExistDB() {
		fmt.Println("Blockchain exists")
		runtime.Goexit()
	}

	//Określamy, gdzie chcemy przechowywać nasze plili bazy danych
	opts := badger.DefaultOptions(path.Dir("/tmp/blocks"))
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	//Tworzymy baze danych
	db, err := badger.Open(opts)
	Handle(err)

	//Zapis i odczyt
	err = db.Update(func(txn *badger.Txn) error {
		//Za wydobycie żetonów górnik zostanie nagrodzony monetami, które może dać do bazy monentarnej
		cryptobaseTransaction := Cryptobase(address, genesisData)
		genesis := Genesis(cryptobaseTransaction)
		//Tworzymy genesis block
		fmt.Println("Genesis  create succesful")
		//Serializacja genesisBlock
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		//Usatawiamy lastHash wzięty z genesisBlock jak początek
		err = txn.Set([]byte("lastHash"), genesis.Hash)

		lastHash = genesis.Hash

		return err
	})
	//Pierwszy błąd, który będzie poza DB
	Handle(err)

	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

//Metoda, która dodaje nowo powstały blok do blockchaina
func (chain *Blockchain) AddBlocks(transac []*Transaction) {
	var lastHash []byte

	//Transakcja tylko do odczytu
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lastHash"))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)

		return err
	})

	Handle(err)

	newBlock := CreateNewBlock(transac, lastHash)
	err = chain.Database.Update(func(txn *badger.Txn) error {
		//Transakcja zapis/odczyt
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lastHash"), newBlock.Hash)
		//Ustawiamy ostatni hash jako wartość nowego hashu
		chain.LastHash = newBlock.Hash

		return err
	})

	Handle(err)
}

// Konwertacja blockchaina w iterator struktury
func (chain *Blockchain) Iterator() *BlockchainIterator {
	iterator := &BlockchainIterator{chain.LastHash, chain.Database}

	return iterator
}

func (iterator *BlockchainIterator) Next() *Block {
	var block *Block

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.NowHash)
		Handle(err)
		encodeTheBlock, err := item.ValueCopy(nil)
		block = Deserialize(encodeTheBlock)

		return err
	})
	Handle(err)

	iterator.NowHash = block.LastHash

	return block
}

//Metoda, która znajduje wszystkie niewydane transakcje
func (chain *Blockchain) UpsentTranscation(pubKeyHash []byte)[]Transaction {
	var upsentTransaction []Transaction

	spent := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, transaction := range block.Transaction {
			transactionID := hex.EncodeToString(transaction.ID)
			//Dostęp do danych wyjsciowych
		OutputStr:
			for outputTransaction, output := range transaction.Output {
				// mapa nie byc równa 0
				if spent[transactionID] != nil {
					for _, spentOutput := range spent[transactionID] {
						//indext wyjsciowy jest równy transanckji wyjsciowej
						if spentOutput == outputTransaction {
							continue OutputStr
						}
					}
				}
				if output.LokedKey(pubKeyHash) {
					upsentTransaction = append(upsentTransaction, *transaction)
				}
			}
			//Sprawdzamy czy transakcja jest na bazie monetarnej
			if transaction.isCryptobase() == false {
				for _, in := range transaction.Input {
					if in.KeyUser(pubKeyHash) {
						inputTansactionID := hex.EncodeToString(in.ID)
						spent[inputTansactionID] = append(spent[inputTansactionID], in.Out)
					}
				}
			}
		}
		if len(block.LastHash) == 0 {
			break
		}
	}
	//transackje przypisane do konta
	return upsentTransaction
}

//Znajdujemy wszystkie nie wykorzytsane transakcje UTXO 
func (chain *Blockchain) Find(pubKeyHash []byte) []OutputStr {
	var outputTransactionOs []OutputStr
	unspentTransaction := chain.UpsentTranscation(pubKeyHash)

	for _, transaction := range unspentTransaction {
		for _, out := range transaction.Output {
			if out.LokedKey(pubKeyHash) {
				outputTransactionOs = append(outputTransactionOs, out)
			}
		}
	}

	return outputTransactionOs

}

//Metoda tworząca normalne transakcje, które nie są oparte na momentach
func (chain *Blockchain) FindSpendingOutput(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutput := make(map[string][]int)
	unspentTransaction := chain.UpsentTranscation(pubKeyHash)
	acc := 0

Work:
	//Niewydane transakcje
	for _, transaction := range unspentTransaction {
		transactionID := hex.EncodeToString(transaction.ID)
		//Iteracje przez wyjscie
		for outputTransaction, out := range transaction.Output {
			if out.LokedKey(pubKeyHash) && acc < amount {
				acc += out.Value
				unspentOutput[transactionID] = append(unspentOutput[transactionID], outputTransaction)

				if acc >= amount {
					break Work
				}
			}
		}
	}

	return acc, unspentOutput
}

//Metoda znajdowania nowych transakcji
func (blockchain *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	iterator := blockchain.Iterator()

	for {
		block := iterator.Next()

		//Porównanie identyfikatorów
		for _, trans := range block.Transaction {
			if bytes.Compare(trans.ID, ID) == 0 {
				return *trans, nil
			}
		}

		if len(block.LastHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction does not exist")
}

//Metoda podpisywania transkacji
func (blockchain *Blockchain) SignedTransaction(trans *Transaction, privateKey ecdsa.PrivateKey) {
	//mapa transakcji
	lastTransactions := make(map[string]Transaction)

	//Znajdowanie transakcji po ID
	for _, input := range trans.Input {
		lastTransaction, err := blockchain.FindTransaction(input.ID)
		Handle(err)

		lastTransactions[hex.EncodeToString(lastTransaction.ID)] = lastTransaction
	}

	trans.Sign(privateKey, lastTransactions)

}

//Metoda weryfikacji transackji
func (blockchain *Blockchain) VerifyTransaction(trans *Transaction) bool {
	lastTransactions := make(map[string]Transaction)

	for _, input := range trans.Input {
		lastTransaction, err := blockchain.FindTransaction(input.ID)
		Handle(err)
		lastTransactions[hex.EncodeToString(lastTransaction.ID)] = lastTransaction
	}

	return trans.Verification(lastTransactions)
}
