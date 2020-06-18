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

Naszą bazą danych bedzie BadgerDB --> wybór był uwarunkowany tym, że jest to bardzo dobra baza danych napisana w języku w GO, która może pomieści tysiace terabajtóœ danych.
 Nie posiada table SQL.

 1. Tworzymy plik.go Blockchain --> tworzymy strukturę, funckję init oraz addblock
 2. Tworzymy funkcje serializacji i deserializacji naszych danych 
 gob -->pakiet służący do strumieniowania gobów czyli naszje zawartości --> wartości binarne są wymienniane miedzy encoderem i decoderem.
 Dalej przejdziemy do bloku Blockchain.go--> 
 1. Zmieniamy skrókturę blockchaina tak oby był powiązany z bazą danych Badger oraz zapamiętywał lasthash
 2. Dodajemy folder, w którym będzie zapisywać lasthash + scieżka do folderu 
 3. Tworzymy bazę danych z funckją update 

 Transackje
 1.Struktura
 2. Funckje Cryptobase oraz SetID
 3. Zamieniamy w strukturze bloku date na Transaction
 

CLI 

WALLET

SIGNATURE


hahs dla transakcji 
serialization 


UTXO 
