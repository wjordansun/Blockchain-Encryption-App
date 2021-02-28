package blockchain

import (
	"fmt"
	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks" //path that stores key-value database.
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash	[]byte
	Database 	*badger.DB
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	// writes the blockchain into the tmp folder
	opt := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opt)
	Handle(err)

	erro := db.Update(func(txn *badger.Txn) error{
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound { //if no lasthash
			fmt.Println("No existing blockchain found")
			genesis := Genesis()								//creates genesis block
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
							// key 			//value
			Handle(err)

			err = txn.Set([]byte("lh"), genesis.Hash)
							// key 			//value
			lastHash = genesis.Hash

			return err
		} else {					// key
			item, err := txn.Get([]byte("lh")) //uses key to get last item in DB
			Handle(err)
			err = item.Value(func(val []byte) error {
				lastHash = append([]byte{}, val...)
				return nil
			}) // gets the last hash to put in lastHash
			return err
		}
	})

	Handle(erro)

	blockchain := BlockChain{lastHash, db}
	return &blockchain


}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error { // reads lastHash
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return err
		})

		return err
	})
	Handle(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	Handle(err)

}

func (chain *BlockChain) Iterator() *BlockChainIterator { //blockchain is turned into a 
	iter := &BlockChainIterator{chain.LastHash, chain.Database} // bloclchain iterator

	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block
	var encodedBlock []byte

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		erro := item.Value(func(val []byte) error {
			encodedBlock = append([]byte{}, val...)
			block = Deserialize(encodedBlock)
			return nil
		})

		return erro
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}




