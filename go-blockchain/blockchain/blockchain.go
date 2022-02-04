package blockchain

import (
	"time"
	"encoding/json"
	"crypto/sha256"
	"fmt"
	"encoding/hex"
	"strconv"
)

const PROOF_OF_WORK_ZEORS = 4

type Transation struct {
	Amount uint64  `json:"Amount"`
	Sender string   `json:"Sender"`
	Recipient string   `json:"Recipient"`
}

type Block struct {
	index uint64
	timestamp time.Time 
	transations []Transation
	nonce uint64 
	hash string 
	previousBlockHash string 
}

type BlockChain struct {
    chain []*Block 
	pendingTransations []Transation 
}

func NewBlockChain() *BlockChain {
	//return &BlockChain{}
	block := &BlockChain{}  
	return block
}

func (bc *BlockChain) CreateNewBlock(nonce uint64, previousBlockHash, hash string) *Block {
	new_block :=  &Block{
		index: uint64(len(bc.chain) + 1),
		timestamp: time.Now(),
		transations: bc.pendingTransations,
		nonce: nonce,
		hash: hash,
		previousBlockHash: previousBlockHash,
	}

	bc.pendingTransations = nil 
    bc.chain = append(bc.chain, new_block)
	return new_block
}

func (bc *BlockChain) GetLastBlock() *Block {
	if len(bc.chain) == 0 {
		return nil 
	}

	return bc.chain[len(bc.chain) - 1]
}

func (bc *BlockChain) CreateNewTransation(amount uint64, sender, recipient string) uint64 {
	transation := Transation {
		Amount: amount,
		Sender: sender, 
		Recipient: recipient, 
	}
    
	bc.pendingTransations = append(bc.pendingTransations, transation)
	//该交易信息应该挂载到下一个新添加的区块
	return uint64(len(bc.chain) + 1)
}

func (bc *BlockChain) HashBlock(block_idx uint64) (string, error) {
	if block_idx >= uint64(len(bc.chain)) {
		return "", fmt.Errorf("index out of range")
	}
	
	block := bc.chain[block_idx]
    hash_content := block.previousBlockHash 
	hash_content += strconv.FormatUint(block.nonce, 10)
	for transation := range block.transations {
		tran_json, err := json.Marshal(&transation)
		if err != nil {
			return "" , err 
		}

		hash_content += string(tran_json)
	}

	h := sha256.New()
	h.Write([]byte(hash_content))
    block.hash = hex.EncodeToString(h.Sum(nil))

	return block.hash, nil
}

func (bc *BlockChain) MiningBlock(block_index uint64) (uint64, error) {
	if block_index >= uint64(len(bc.chain)) {
		return 0, fmt.Errorf("index out of range")
	} 

	block := bc.chain[block_index]
	block.nonce = 0
	hash, err := bc.HashBlock(block_index)
	if err != nil {
		return 0, err
	}
	hash_head := ""
	for i := 0; i < PROOF_OF_WORK_ZEORS; i++ {
		hash_head += "0"
	}
	for hash[0 : PROOF_OF_WORK_ZEORS] != hash_head {
		block.nonce += 1
		hash, err = bc.HashBlock(block_index)
		if err != nil {
			return 0, err
		}
	}

	return block.nonce, nil
}