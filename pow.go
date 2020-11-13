package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"strconv"
	"time"
)

/*
模拟比特币POW算法
1、生成创世区块 加入区块链中
2、循环判断生成的区块是否合法
		合法：将区块加入区块链中
		不合法：继续生成新的区块

*/
const difficulty = 5

type Block struct {
	//不参与SHA-256运算
	Index int    //区块链中数据记录的位置，区块高度
	Hash  string //代表这个数据记录的SHA256标识符

	//参与SHA-256运算
	version   int    //版本号
	PreHash   string //上一条记录的SHA256标识符
	Merkle    string //Merkle 根
	TimeStamp string //时间戳
	Target    string //目标值
	Nonce     string //PoW 挖矿中符合条件的随机数
}

//存放区块数据的集合
var Blockchain []Block

//生成区块
func generateBlock(oldBlock Block) Block {
	var newBlock Block

	//获取当前时间并设置新区块信息
	t := time.Now()
	newBlock.version = oldBlock.version
	newBlock.Merkle = oldBlock.Merkle
	newBlock.Index = oldBlock.Index + 1
	newBlock.TimeStamp = t.String()
	newBlock.PreHash = oldBlock.Hash
	newBlock.Target = oldBlock.Target

	//循环nonce递增，尝试hash
	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		//判断计算的hash是否头部满足要求
		if !isHashVaild(calculateHash(newBlock), newBlock.Target) {
			//	fmt.Println(calculateHash(newBlock), "do more work!")
			//time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(calculateHash(newBlock), "work done!")
			newBlock.Hash = calculateHash(newBlock)
			break
		}
	}
	return newBlock
}

//验证hash是否合法，hash小于目标难度
func isHashVaild(hash, target string) bool {
	return hash < target
}

//计算hash
//两次SHA-256算法返回hash
func calculateHash(block Block) string {
	record := strconv.Itoa(block.version) + block.TimeStamp + block.Merkle + block.PreHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	firsthash := h.Sum(nil)
	secondhash := hex.EncodeToString(firsthash)
	g := sha256.New()
	g.Write([]byte(secondhash))
	hashed := g.Sum(nil)
	return hex.EncodeToString(hashed)
}

//验证区块是否合法
func isBlockvaild(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PreHash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

//主函数
func main() {
	//获取当前时间戳
	t := time.Now()
	genesisBlock := Block{}
	//创世区块
	genesisBlock = Block{version: 1, PreHash: "",
		Merkle: "d32f29bae1691d327b5eb0c604bbfeab90bede4ec5d98a2e1d3b2495fc0eb46e", TimeStamp: t.String(),
		Target: "0000021g89544adsfgdfgf9ee29b2e148bc84b5099fd600af3569b6075a835c5", Nonce: "", Index: 0,
		Hash: calculateHash(genesisBlock)}
	Blockchain = append(Blockchain, genesisBlock)
	//格式化输出区块内容
	spew.Dump(genesisBlock)
	for {
		if isBlockvaild(genesisBlock, Blockchain[len(Blockchain)-1]) {
			if len(Blockchain) != 0 {
				Blockchain = append(Blockchain, genesisBlock)
			}
			//格式化输出区块内容
			spew.Dump(genesisBlock)
		} else {
			genesisBlock = generateBlock(Blockchain[len(Blockchain)-1])
		}
	}
}
