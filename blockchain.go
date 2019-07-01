package main

import(
	"fmt"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net"
	"bufio"
)

type Block struct {
	PrevHash string
	Data string
	Hash string
}

func newBlock(data, prevHash string) *Block {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s%s", data, prevHash)))
	hash := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return &Block{prevHash, data, hash}
}

type BlockChain struct{
	Blocks []*Block
}

func newBlockChain() *BlockChain {
	return &BlockChain{[]*Block{newBlock("","0")}}
}

var bc *BlockChain = newBlockChain()

func handleBC(con net.Conn) {
	defer con.Close()
	r := bufio.NewReader(con)
	_, _ = r.ReadString('\n')
	buf, _ := json.Marshal(bc)
	fmt.Fprintln(con, string(buf))
}

func servBC() {
	ln, _ := net.Listen("tcp", "localhost:8000")
	defer ln.Close()
	for {
		con, _ := ln.Accept()
		go handleBC(con)
	}
}

func main() {
	//b := newBlock("","0")
	//fmt.Println(b)
}