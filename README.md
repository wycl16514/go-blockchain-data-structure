前面章节我们了解了区块链的基本原理，涉足到其底层加密算法，现在我们可以开始将区块链技术应用起来，我们看看如何建立一个基于区块链技术的，去中心化的交易平台。这里我们使用go来实现区块链的基本算法，使用gin完成后台restful服务功能，然后使用reactjs完成前端开发，通过具体的应用，我们看看区块链技术如何保证信息的不可篡改性，同时我们也逐步了解多节点共识，网络同步等重要概念。

首先我们先定义系统后台的基本数据结构，在本地目录创建GO-BLOCKCHAIN，然后增加一个名为blockchain.go的文件，首先要做的是定义数据结构：
```
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
```
在前面章节我们描述过区块链数据结构的定义，同时也接触了智能合约，因此上面结构的定义，各个字段的内容想必都比较清楚。Block对应区块，Transation用来记录交易信息，它包含了交易的数额，接收者和发送者，PROOF_OF_WORK_ZEROS对应挖矿时要在哈希前面添加的0的个数，由于”挖矿“会消耗很多算力，因此我们这里只使用4个0作为例子。

当我们开发的后台成功挖到一个区块后，它会形成一个Block结构并添加到chan队列里，当前所有交易会放置在pendingTransations队列中，一旦区块挖出来后，我们会把当前交易转移到区块中，这样交易信息就可以实现不可更改性，我们看看如何创建一个区块，增加实现函数如下：
```
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
```
上面函数的实现用于创建一个区块，其中nonce变量对应我们以前说过的用于计算开头给定个数0的数字，prevoiusBlockHash对应上一个区块的哈希值，hash对应本区块的哈希值，注意看当区块生成时，它会把当前存在的交易信息存储到区块对应的交易列表中，完成了上面代码后，我们创建blockchain_test.go文件，对已经完成的接口进行单元测试，代码如下：
```
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
```
在TestCreateNewBlock函数中，它调用CreateNewBlock接口创建区块，然后判断创建的区块对应各个字段是否跟预期的一致，我们这里使用的go-cmp来比较两个结构体。接下来我们要实现的是获取当前区块链中最后一个区块的接口，代码如下：
```
func (bc *BlockChain) GetLastBlock() *Block {
	if len(bc.chain) == 0 {
		return nil 
	}

	return bc.chain[len(bc.chain) - 1]
}
```
它的逻辑很简单，就是直接返回chain列表中最后一个元素，同理我们也添加一个单元测试检验它实现的逻辑，在blockchain_test.go里面添加代码如下：
```
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
```
在上面的测试函数中，我们创建两个区块并依次加入链表，然后调用GetLastBlock接口获取最后一个区块，最后判断获取区块的字段与我们最后加入的区块是否一致。下面我们看看如何创建交易数据，由于交易信息只有三个字段，因此它的生成逻辑
```
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
```
由于交易信息只有三个字段，因此创建时只要传入三个参数，这里需要注意的是，有新交易数据生成时，我们需要将它添加到pendingTransation队列，然后等到新区块生成后才将他们加入到区块里。该接口对应的单元测试如下：
```
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
```
测试时我们调用CreateNewTransation接口，传入对应参数，接着检验对应字段是否一致，这里我们要特别判断生成的新交易是否放置在pendingTransation列表。同时我们还要确保一旦生成一个新区块后，代码需要将当前交易信息添加到区块中。下一个我们要完成比较复杂的逻辑，那就是对区块计算哈希，代码如下：
```
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
```
接口传入的参数是区块在链中的下标，计算哈希需要三部分信息，第一是上一个区块的哈希，第二是交易信息，我们需要将交易信息转换为JSON格式，然后再转换为字符串才能用于计算哈希，第三就是用于创建哈希头部给定0的数值，注意现在这个函数还半成品，因为nonce的值还没有确定，我们需要确保给定nonce的值能使得哈希后，所得结果前头有4个0，因此我们还需要添加”挖矿“接口：
```
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
```
挖矿接口的逻辑在前面章节也实现过，我们从0开始不断的增加nonce的值，直到计算出来的哈希值在开头有4个0位置。接下来我们看看对应的单元测试：
```
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
```
看过前面章节的同学应该比较容易的了解这里实现的代码逻辑。我们构造这个系统的目的还在与掌握分布式系统中一些非常重要的算法或概念，例如网络同步，共识达成等，下一节我们使用gin来完成后台服务，向外提供调用接口.
