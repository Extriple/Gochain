package blockchain

//Struktura blockchaina
type Blockchain struct {
	Blocks []*Block
}

// Struktura bloku
type Block struct {
	Hash     []byte
	Data     []byte
	LastHash []byte
	Nonce    int
}

//Metoda, która tworzy nowy blok
func CreateNewBlock(data string, lastHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), lastHash, 0}
	proofOfWork := StatNewProof(block)
	nonce, hash := proofOfWork.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

//Metoda, która dodaje nowo powstały blok do blockchaina
func (chain *Blockchain) AddBlocks(data string) {
	lastBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateNewBlock(data, lastBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
}

//Genesis
func Genesis() *Block {
	return CreateNewBlock("Genesis", []byte{})
}

//Metoda, która zespala wszystkie wyżej wymienione bloki i tworzy naszego blockchaina
func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}}
}
