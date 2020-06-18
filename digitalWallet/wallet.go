package digitalWallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

//Stałe odnoszące sie do sumy kontrolnej
const (
	checksum = 4
	version  = byte(0x00)
)

//Struktura portfela
type Wallet struct {
	//algorytm podpisywania cyfrowego
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

//Metoda generująca adres dla każdego z użytkowników
func (w Wallet) Address() []byte {
	publicHash := HashingPublicKey(w.PublicKey)

	versionHash := append([]byte{version}, publicHash...)
	checksum := Checksum(versionHash)

	completeHash := append(versionHash, checksum...)

	address := Encode58Base(completeHash)

	//fmt.Printf("public key: %x\n", w.PublicKey)
	//fmt.Printf("public hash: %x\n", publicHash)
	//fmt.Printf("address: %x\n", address)

	return address
}

// Funckja sprawdzająca poprawność naszego adresu w celu dodawania signatury
func ValidateAddress(address string) bool {
	publicKeyHash := DecodeBase58([]byte(address))
	actuallyCheckSum := publicKeyHash[len(publicKeyHash)-checksum:]
	version := publicKeyHash[0]
	publicKeyHash = publicKeyHash[1 : len(publicKeyHash)-checksum]
	targetChecksum := Checksum(append([]byte{version}, publicKeyHash...))

	return bytes.Compare(actuallyCheckSum, targetChecksum) == 0
}

//Funckja tworząca klucz prywatny
func NewPrivateKey() (ecdsa.PrivateKey, []byte) {
	//Typ krzywej eliptycznej
	ellipitcCurve := elliptic.P256()

	private, err := ecdsa.GenerateKey(ellipitcCurve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	publicKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, publicKey
}

//Funckja tworząca  portfel
func BuildWallet() *Wallet {
	privateKey, publicKey := NewPrivateKey()
	digitalWallet := Wallet{privateKey, publicKey}

	return &digitalWallet
}

//Funckja przekształcająca klucz publiczny w bajtowy hash
func HashingPublicKey(publicKey []byte) []byte {
	publicHash := sha256.Sum256(publicKey)
	//Zapisujemy public hash
	hashing := ripemd160.New()
	_, err := hashing.Write(publicHash[:])
	if err != nil {
		log.Panic(err)
	}
	publicRipemd := hashing.Sum(nil)

	return publicRipemd
}

//Funckja sumy kontrolnej
func Checksum(loadPaynament []byte) []byte {
	hash1 := sha256.Sum256(loadPaynament)
	hash2 := sha256.Sum256(hash1[:])

	return hash2[:checksum]
}
