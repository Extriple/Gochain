package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

// Struktura bloku
type Block struct {
	Hash        []byte
	Transaction []*Transaction
	LastHash    []byte
	Nonce       int
}

//Metoda, która uwzględnia algorytm proof of work w stosunku do transackji 
func (b *Block) HashTransaction() []byte {
	var  transacHash [32]byte
	var transacHashes [][]byte

	for _, transac := range b.Transaction{
		transacHashes = append(transacHashes, transac.ID)
	}
	transacHash = sha256.Sum256(bytes.Join(transacHashes, []byte{}))

	return transacHash[:]
}

//Metoda, która tworzy nowy blok
func CreateNewBlock(tranc []*Transaction, lastHash []byte) *Block {
	block := &Block{[]byte{}, tranc, lastHash, 0}
	proofOfWork := StatNewProof(block)
	nonce, hash := proofOfWork.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

//Genesis
func Genesis(Cryptobase *Transaction) *Block {
	return CreateNewBlock([]*Transaction{Cryptobase}, []byte{})
}

// Funckje serializacji i deserializacji danych
func (b *Block) Serialize() []byte {
	var buffor bytes.Buffer
	encoder := gob.NewEncoder(&buffor)

	err := encoder.Encode(b)

	Handle(err)

	//zwracamy bajtową cześć naszego wyniku
	return buffor.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block

}

//Funckja pomocnicza wywołująca error/ funckja skrótowa.
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
