package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

//Schemat proof of work:
//1. Bierzemy Date z naszego bloku
//2. Usatwaimy tzn licznil, który na początku wynosi 0
//3. Tworzymy skrót czyt. hash, który jest spójny z licznikiem
//4. Sprawdzamy poprawność wymagań jakie musi posiadać nasz hash

// Wymagania proof of work:
//1. Pierwsze klika bajtów musi miec wartość 0
//2. Czas wydobycia bloku pozostaje taki sam
//3, Stawka czyt. nagroda również się nie zmienia

//Dowolona wartość

const Difficulty = 10

//Struktura proof of work

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

//Metoda tworząca działanie nowego proof of work

func StatNewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	proofOfWork := &ProofOfWork{b, target}

	return proofOfWork
}

//Funckja tworząca nowy hash

func (proofOfWork *ProofOfWork) NewInitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			proofOfWork.Block.LastHash,
			proofOfWork.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return data
}

//Funkcja oblicznieniowa służąca do sprawdzenia poprawności proof of work

func (proofOfWork *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := proofOfWork.NewInitData(nonce)
		hash := sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		//Konwertujemy skrót na wieszką liczbę
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(proofOfWork.Target) == -1 {
			break
		} else {
			nonce++

		}
	}
	fmt.Println()

	return nonce, hash[:]
}

//Funckja sprawdzenia poprawności czyli czy wygamagania jakie musi spełnić algorytm proof of work i hash i czy są one popranie wykonane

func (proofOfWork *ProofOfWork) ValidateRequirement() bool {

	var hashInt big.Int

	data := proofOfWork.NewInitData(proofOfWork.Block.Nonce)

	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(proofOfWork.Target) == -1
}

//Funckcja operacyjna służąca do pomocy szyfrowania, która przyjmuje liczbę całkowitą

func ToHex(number int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, number)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}
