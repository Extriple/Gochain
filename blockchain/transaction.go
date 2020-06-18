package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/PatrykSer/Gochain/digitalWallet"
)

//Podstawowa struktura transakcji

type Transaction struct {
	ID     []byte
	Input  []InputStr
	Output []OutputStr
}

//Metoda serializacji transakcji
func (trans Transaction) Serialization() []byte {
	var encoded bytes.Buffer

	encod := gob.NewEncoder(&encoded)
	err := encod.Encode(trans)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

//Metoda tworząca hash dla transakcji
func (trans *Transaction) Hash() []byte {

	var hash [32]byte

	transactionCopy := *trans
	transactionCopy.ID = []byte{}

	hash = sha256.Sum256(transactionCopy.Serialization())

	return hash[:]
}

//Metoda SetID tworzy skrót na podstawie bajtu --> dajemy przez funckję Encode --> Sha256 --> powstaje hash
func (tranc *Transaction) SetID() {
	var encode bytes.Buffer
	var hash [32]byte

	encoded := gob.NewEncoder(&encode)
	err := encoded.Encode(tranc)
	Handle(err)

	hash = sha256.Sum256(encode.Bytes())
	tranc.ID = hash[:]
}

//Transackja na bazie moment = 1 input i 1 output

func Cryptobase(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coin to %s", to)

	}

	input := InputStr{[]byte{}, -1, nil, []byte(data)}
	output := NewOutputTransaction(100, to)

	tranc := Transaction{nil, []InputStr{input}, []OutputStr{*output}}
	//Id = hash
	tranc.SetID()

	return &tranc
}

//Metoda tworząca nowe transakcje
func NewTransaction(from, to string, amount int, chain *Blockchain) *Transaction {
	var inputs []InputStr
	var outputs []OutputStr

	wallet, err := digitalWallet.CreateNewWallet()
	Handle(err)
	wal := wallet.GetNewWallet(from)
	publicKeyHash := digitalWallet.HashingPublicKey(wal.PublicKey)
	acc, validOut := chain.FindSpendingOutput(publicKeyHash, amount)

	if acc < amount {
		log.Panic("Error: not enough coins")
	}

	for transactionID, outs := range validOut {
		txID, err := hex.DecodeString(transactionID)
		Handle(err)

		for _, out := range outs {
			input := InputStr{ txID, out, nil, wal.PublicKey }
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, *NewOutputTransaction(amount, to))

	if acc > amount {

		outputs = append(outputs, *NewOutputTransaction(acc-amount, from))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	chain.SignedTransaction(&tx, wal.PrivateKey)

	return &tx
}

//Funckja sprawdzająca, czy nasza transakcja jest na bazie kryptowalut
func (transc *Transaction) isCryptobase() bool {
	return len(transc.Input) == 1 && len(transc.Input[0].ID) == 0 && transc.Input[0].Out == -1
}

//Metoda podpisująca i werfikująca transackje
func (trans *Transaction) Sign(privateKey ecdsa.PrivateKey, lastTransactions map[string]Transaction) {
	if trans.isCryptobase() {
		return
	}

	for _, input := range trans.Input {
		if lastTransactions[hex.EncodeToString(input.ID)].ID == nil {
			log.Panic("ERROR: The last transaction is not correct")
		}
	}
	transCopy := trans.Trimmed()

	//Każdy input ma podpis ustawiony na 0
	//Metoda serializacji transakcji szyfruje dane
	//Przekazanie hashu do identyfikatora transakcji

	for inputID, input := range transCopy.Input {

		prevTrans := lastTransactions[hex.EncodeToString(input.ID)]
		transCopy.Input[inputID].Signature = nil
		transCopy.Input[inputID].PubKey = prevTrans.Output[input.Out].PublicKeyHash
		transCopy.ID = transCopy.Hash()
		transCopy.Input[inputID].PubKey = nil

		//Losowa liczba wygenerowane i przekazana do klucza prywatnego
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, transCopy.ID)
		Handle(err)
		//Umieszczamy gotowy podpis w odniesieniu do naszych danych wejsciowych
		signature := append(r.Bytes(), s.Bytes()...)
		trans.Input[inputID].Signature = signature

	}
}

//Metoda weryfikacji, która uwzględnia mapę poprzednich transkacji
func (trans *Transaction) Verification(prevTrans map[string]Transaction) bool {
	// sprawdzamy czy dana transakcja posiada baze moneterną/kryptograficzną
	if trans.isCryptobase() {
		return true
	}

	for _, input := range trans.Input {
		if prevTrans[hex.EncodeToString(input.ID)].ID == nil {
			log.Panic("Last transaction not correct")
		}
	}

	transCopy := trans.Trimmed()
	curveElliptic := elliptic.P256()

	//Ta sama pętla co przy metodzie napisanej wyżwj, ponieważ  potrzebujemy tej samej transakcji, która jest już  podpisana
	for inputID, input := range trans.Input {
		prevTrans := prevTrans[hex.EncodeToString(input.ID)]
		transCopy.Input[inputID].Signature = nil
		transCopy.Input[inputID].PubKey = prevTrans.Output[input.Out].PublicKeyHash
		transCopy.ID = transCopy.Hash()
		transCopy.Input[inputID].PubKey = nil

		//Rozpakowujemy wszystkie dane, które zostały umieszczone w podpisie i pubie
		//klucz publiczny = para współrzędnych

		p := big.Int{}
		x := big.Int{}

		//Długośc pdpisu
		signLen := len(input.Signature)
		p.SetBytes(input.Signature[:(signLen / 2)])
		x.SetBytes(input.Signature[(signLen / 2):])

		r := big.Int{}
		w := big.Int{}

		keylen := len(input.PubKey)

		r.SetBytes(input.PubKey[:(keylen / 2)])
		w.SetBytes(input.PubKey[(keylen / 2):])
		//Weryfikacja poprzez ECDSA
		PubKeyInRAW := ecdsa.PublicKey{curveElliptic, &r, &w}

		if ecdsa.Verify(&PubKeyInRAW, transCopy.ID, &p, &x) == false {
			return false
		}
	}

	return true

}

//Metoda służąca do kopiowania  transakcji i generowania jej
func (trans *Transaction) Trimmed() Transaction {
	//Dwie tablice z danymi wejsciowmi i wyjsciowymi  wraz z pętlami , które odnosza sie do struktury danych
	var input []InputStr
	var output []OutputStr

	for _, in := range trans.Input {
		input = append(input, InputStr{in.ID, in.Out, nil, nil})
	}

	for _, out := range trans.Output {
		output = append(output, OutputStr{out.Value, out.PublicKeyHash})
	}

	transCopy := Transaction{trans.ID, input, output}

	return transCopy
}

//Metoda konwersji transkacji na stringa
func (trans Transaction) String() string {
	var line []string

	line = append(line, fmt.Sprintf("----Transaction %x:", trans.ID))

	for i, input := range trans.Input {

		line = append(line, fmt.Sprintf("   Input %d", i))
		line = append(line, fmt.Sprintf("   TransactionID %x", input.ID))
		line = append(line, fmt.Sprintf("       Out:    %d", input.Out))
		line = append(line, fmt.Sprintf("       Signature:  %x", input.Signature))
		line = append(line, fmt.Sprintf("           PubKey:     %x", input.PubKey))
	}

	for i, output := range trans.Output {
		line = append(line, fmt.Sprintf("   Output &d:", i))
		line = append(line, fmt.Sprintf("       Value:  %d", output.Value))
		line = append(line, fmt.Sprintf("       Script: %x", output.PublicKeyHash))

	}

	return strings.Join(line, "\n")
}
