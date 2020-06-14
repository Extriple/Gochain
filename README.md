# Gochain
Blockchain in Golang
Blockchain w golangu,

=======================
Początek:

-dodajemy go.mod: go mod init [link do githuba]
- Tworzymy Block, Blockchain,ExteractBlock, Genesis, AddBlock, CreateBlock, InitBlockchain w main go
- przenosimy wyżej wymienione funckję/metody do Block.go
-Tworzymy proofOfWork

LSH --> locality sensitive hashing --> biblioteka zawierające różne alogorytmy związane z funckją wrazliwości.
uint --> conajmniej 32 bity, typ całkowity

Proof of Work: 
-StartNewProof --> tworzymy nasz alogorytm proof of work 
- NewIntData --> tworzymy nowego hasha
- ToHex --> funkcja pomocniczna służąca do szyfrowania 

na tym etapie mamy usatwiony nasz licznik 

- tworzymy funckję obliczenią start
- usuwamy ExteractBlock,a następnie modyfikujemy funckję CreateBlock poprzez dodanie metody proof of work 
