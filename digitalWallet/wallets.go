package digitalWallet

import (
	"bytes"
	"encoding/gob"
	"crypto/elliptic"
	"fmt"
	"io/ioutil"
	"os"
	"log"
)

//Miejsce przechowywania portfeli
const fileWithWalletData = "./tmp/wallets.data"

//Strutura portfela jest kluczem prywatnym i publicznym  - więc łątwo mozna pobrac portfel  - dane
type Wallets struct {
	Wallets map[string]*Wallet
}

//Metoda tworzenia nowego portfela
func CreateNewWallet() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadingFile()

	return &wallets, err
}

//Metoda dodawania nowego portfela
func (wallets *Wallets) AddNewWallet() string {
	wallet := BuildWallet()
	address := fmt.Sprintf("%s", wallet.Address())

	//Nowo dodany portfel umieszczamy w mapie z adresami  [klucz]
	wallets.Wallets[address] = wallet

	return address
}

//Funckja dostarczjąca wszystkie adresy naszych portfeli 
func (wallets * Wallets) GetAllNewAddress() []string {
	var bigAddress[]string

	for address := range wallets.Wallets {
		bigAddress = append(bigAddress, address)
	}

	return bigAddress
}

//Funkcja zwracająca portfel ,a w zasadzie jego mapę po której szybko możemy odnaleźć adres
func (wallets Wallets) GetNewWallet(address string) Wallet {
	return *wallets.Wallets[address]
}

//Funkcja wczytująca pliki 
func (wal *Wallets) LoadingFile() error {
	//Sprawdzamy czy plik istnieje
	if _, err := os.Stat(fileWithWalletData); os.IsNotExist(err) {
		return err
	}

	var wallets Wallets
	//Sprawdzamy czy mozemy odczytać pliki 
	file, err := ioutil.ReadFile(fileWithWalletData)
	if err != nil {
		return err
	}
		//Czytnik bajtów 
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(file))
	err = decoder.Decode(&wallets)
	if err != nil { 
		return err
	}

	wal.Wallets = wallets.Wallets
	return nil 
}

//Funckja zapisywania danych 
func (wallets *Wallets) Saving() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())
		//Nowy koder 
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(wallets)
	if err != nil {
		log.Panic(err)

	}
		//Zapis 
	err = ioutil.WriteFile(fileWithWalletData, content.Bytes(), 0644)
	if err != nil { 
		log.Panic(err)
	}
}
