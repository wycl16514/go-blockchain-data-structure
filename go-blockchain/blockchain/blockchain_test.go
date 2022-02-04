package blockchain 

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/google/go-cmp/cmp"
	"fmt"
)

func TestNewBlockChain(t *testing.T) {
	blockchain := NewBlockChain()
	require.Equal(t, len(blockchain.chain), 0)
	require.Equal(t, len(blockchain.pendingTransations), 0)
}

func blockComparer(x , y Block ) bool {
	check_fields := x.index == y.index && x.nonce == y.nonce && x.hash == y.hash &&
	x.previousBlockHash == y.previousBlockHash 
	if check_fields != true {
		return false  
	}
	
	if len(x.transations) != len(y.transations) {
		return false 
	}

	for i := 0; i < len(x.transations); i++ {
		if cmp.Equal(x.transations[i], y.transations[i]) != true {
			return false 
		}
	}

    return true
}

func TestCreateNewBlock(t *testing.T) {
	blockchain := NewBlockChain()
	block_return := blockchain.CreateNewBlock(2389, "OIUOEREDHKHKD", "78s97d4x6dsf")
	block_want := &Block {
		index: 1,
		nonce: 2389,
		previousBlockHash: "OIUOEREDHKHKD",
		hash: "78s97d4x6dsf",
	}
    comparer := cmp.Comparer(blockComparer)
	diff := cmp.Diff(block_return, block_want, comparer)
	require.Equal(t, diff, "")
	require.Equal(t, len(blockchain.chain), 1)
}

func TestGetLastBlock(t *testing.T) {
	blockchain := NewBlockChain()
	require.Nil(t, blockchain.GetLastBlock())

	blockchain.CreateNewBlock(2389, "OIUOEREDHKHKD", "78s97d4x6dsf")
	blockchain.CreateNewBlock(2899, "UINIUN90ANSDF", "99889HBAIUSBDF")
	block_return := blockchain.GetLastBlock()
	block_want := &Block {
		index: 2,
		nonce: 2899,
		previousBlockHash: "UINIUN90ANSDF",
		hash: "99889HBAIUSBDF",
	}
    comparer := cmp.Comparer(blockComparer)
	diff := cmp.Diff(block_return, block_want, comparer)
	require.Equal(t, diff, "")
}



func TestCreateNewTransation(t *testing.T) {
	blockchain := NewBlockChain()
	next_block_idx := blockchain.CreateNewTransation(100, "ALEXHT854", "JENN5BG")
	require.Equal(t, next_block_idx, uint64(1))
	require.Equal(t, len(blockchain.pendingTransations), 1)

	transation_added := blockchain.pendingTransations[0]
	transation_want := Transation {
		Amount: 100,
		Sender: "ALEXHT854",
		Recipient: "JENN5BG",
	}
	diff := cmp.Diff(transation_added, transation_want)
	require.Equal(t, diff, "")

	//挖到一个区块后，交易信息要放置到区块里面
	block := blockchain.CreateNewBlock(2389, "OIUOEREDHKHKD", "78s97d4x6dsf")
	require.Equal(t, len(block.transations), 1)
	require.Equal(t, len(blockchain.pendingTransations), 0)
	diff = cmp.Diff(block.transations[0], transation_want)
}

func TestHashBlock(t *testing.T) {
	blockchain := NewBlockChain()
	blockchain.CreateNewTransation(100, "ALEXHT854", "JENN5BG")
	block := blockchain.CreateNewBlock(2389, "OIUOEREDHKHKD", "78s97d4x6dsf")

	hash, err := blockchain.HashBlock(1)
	require.NotNil(t, err)
    hash , err = blockchain.HashBlock(0)
    require.Nil(t, err)

	fmt.Printf("hash : %s\n", hash)
	
	diff := cmp.Diff(hash, block.hash)
	require.Equal(t, diff, "")
}

func TestMiningBlock(t *testing.T) {
	blockchain := NewBlockChain()
	blockchain.CreateNewTransation(100, "ALEXHT854", "JENN5BG")
	block := blockchain.CreateNewBlock(0, "OIUOEREDHKHKD", "78s97d4x6dsf")

	nonce, err := blockchain.MiningBlock(1)
	require.NotNil(t, err)

	nonce, err = blockchain.MiningBlock(0)
	require.Nil(t, err)

	heading := ""
	for i := 0; i < PROOF_OF_WORK_ZEORS; i++ {
		heading += "0"
	}

	require.Equal(t, block.hash[0 : PROOF_OF_WORK_ZEORS], heading)

	fmt.Printf("nonce is %d, hash after mining the block: %s\n", nonce, block.hash)
}

