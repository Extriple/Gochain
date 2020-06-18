package digitalWallet

import (
	"log"

	"github.com/mr-tron/base58"
)



//Funckja kodująca i dekodująca wygenerowane wyżej skróty ---- algorytm base 58 
func Encode58Base(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}




func DecodeBase58(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	if err != nil {
		log.Panic(err)
	}

	return decode

	// 0 0 l L + /
}