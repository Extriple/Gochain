package blockchain

import (
	"bytes"

	"github.com/PatrykSer/Gochain/digitalWallet"
)
type OutputStr struct {
	Value			int
	PublicKeyHash	[]byte
      }
      
type InputStr struct{
	ID			[]byte
	Out			int
	Signature		   []byte
	PubKey		    []byte
      }
      //Funckja tworząca klucz dla
      func (in* InputStr) KeyUser(pubKeyHash []byte) bool {
	      lockHash := digitalWallet.HashingPublicKey(in.PubKey)

	      return bytes.Compare(lockHash, pubKeyHash) == 0
      }
      //Metoda blokująca  dane wyjsciowe
func (output *OutputStr) Locked(address []byte) {
	publicKeyHash := digitalWallet.DecodeBase58(address)
	publicKeyHash = publicKeyHash[1 : len(publicKeyHash) -4]
	output.PublicKeyHash = publicKeyHash
}
//Funckja sprawdzająca czy blokoda jest wykonan w sposób prawidłowy 
func (out *OutputStr) LokedKey(pubKey []byte) bool {
	return  bytes.Compare(out.PublicKeyHash, pubKey) == 0
}
//MEtoda tworząca nowy wynik
func NewOutputTransaction(value int, address string)  *OutputStr {
	transactionOuput  := &OutputStr{value, nil}
	transactionOuput.Locked([]byte(address))

	return transactionOuput

}