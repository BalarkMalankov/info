package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"time"
)

type HistoriaClinica struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Paciente     string `json:"paciente"`
	CreationDate string `json:"creation_date"`
	Historia     string `json:"historia"`
	IsGenesis    bool   `json:"is_genesis"`
}

type Block struct {
	Pos       int
	Data      HistoriaClinica
	Timestamp string
	Hash      string
	PrevHash  string
}

func (b *Block) generateHash() {
	// get string val of the Data
	bytes, _ := json.Marshal(b.Data)
	// concatenate the dataset
	data := string(b.Pos) + b.Timestamp + string(bytes) + b.PrevHash
	hash := sha256.New()
	hash.Write([]byte(data))
	b.Hash = hex.EncodeToString(hash.Sum(nil))
}

type Blockchain struct {
	blocks []*Block
}

var BlockChain *Blockchain

func CreateBlock(prevBlock *Block, Item HistoriaClinica) *Block {
	block := &Block{}
	block.Pos = prevBlock.Pos + 1
	block.Timestamp = time.Now().String()
	block.Data = Item
	block.PrevHash = prevBlock.Hash
	block.generateHash()

	return block
}

func (bc *Blockchain) AddBlock(data HistoriaClinica) {
	// get previous block
	prevBlock := bc.blocks[len(bc.blocks)-1]
	// create new block
	block := CreateBlock(prevBlock, data)
	bc.blocks = append(bc.blocks, block)
}

func GenesisBlock() *Block {
	return CreateBlock(&Block{}, HistoriaClinica{IsGenesis: true})

}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{GenesisBlock()}}
}

func validBlock(block, prevBlock *Block) bool {

	if prevBlock.Hash != block.PrevHash {
		return false
	}

	if !block.validateHash(block.Hash) {
		return false
	}

	if prevBlock.Pos+1 != block.Pos {
		return false
	}
	return true
}

func (b *Block) validateHash(hash string) bool {
	b.generateHash()
	if b.Hash != hash {
		return false
	}
	return true
}

func getBlockchain(w http.ResponseWriter, r *http.Request) {
	jbytes, err := json.MarshalIndent(BlockChain.blocks, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	// write JSON string
	io.WriteString(w, string(jbytes))
}
func writeBlock(w http.ResponseWriter, r *http.Request) {
	var checkoutItem HistoriaClinica
	if err := json.NewDecoder(r.Body).Decode(&checkoutItem); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not write Block: %v", err)
		w.Write([]byte("could not write block"))
		return
	}
	// create block
	BlockChain.AddBlock(checkoutItem)
	resp, err := json.MarshalIndent(checkoutItem, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not marshal payload: %v", err)
		w.Write([]byte("could not write block"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
func newHistoriaClinica(w http.ResponseWriter, r *http.Request) {
	var historiaClinica HistoriaClinica
	if err := json.NewDecoder(r.Body).Decode(&historiaClinica); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not create: %v", err)
		w.Write([]byte("could not create new HistoriaClinica"))
		return
	}
	// We'll create an ID, concatenating the ISDBand publish date
	// This isn't an efficient way but it serves for this tutorial
	h := md5.New()
	io.WriteString(h, historiaClinica.Historia+historiaClinica.CreationDate)
	historiaClinica.ID = fmt.Sprintf("%x", h.Sum(nil))

	// send back payload
	resp, err := json.MarshalIndent(historiaClinica, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not marshal payload: %v", err)
		w.Write([]byte("could not save historiaClinica data"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

var ipsRed []string

func registerNode(w http.ResponseWriter, r *http.Request) {
	var ip String

	if err := json.NewDecoder(r.Body).Decode(&ip); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not create: %v", err)
		w.Write([]byte("could not create new ip"))
		return
	}

	ipsRed = append(ipsRed, ip)

}

func getAllNodes(w http.ResponseWriter, r *http.Request) {

	jbytes, err := json.MarshalIndent(ipsRed, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	// write JSON string
	io.WriteString(w, string(jbytes))
}

func main() {

	ipGenesis := "0.0.0.0" //aqui va la ip principal
	isGenesis := true
	BlockChain = NewBlockchain()

	r := mux.NewRouter()

	if GetOutboundIP() != ipGenesis {
		isGenesis = false

		url := ipGenesis + "/registerNode"

		var jsonStr = json.newEncoder().Encode(GetOutboundIP())

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		cliente := http.client{
			Timeout: time.Second * 2,
		}

		req, err := http.NewRequest(http.MethodGet, ipGenesis+"/getAllNodes", nil)

		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("User-Agent", "spacecount-tutorial")

		res, getErr := spaceClient.Do(req)
		if getErr != nil {
			log.Fatal(getErr)
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		jsonErr := json.Unmarshal(body, &ipsRed)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

	} else {
		r.HandleFunc("/registerNode", registerNode).Methods("POST")
		r.HandleFunc("/getAllNodes", getAllNodes).Methods("GET")

	}

	// register router

	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newHistoriaClinica).Methods("POST")

	// dump the state of the Blockchain to the console
	go func() {
		reader := bufio.NewReader(os.Stdin)
		auxText := ""
		for {
			fmt.Println("Bienvenido, seleccione (1) para agregar una historia clinica, (2) para listar todas")

			text, _ := reader.ReadString('\n')

			switch text {
			case "1":
				var nuevaHistoria HistoriaClinica
				fmt.Println("Ingrese el nombre del paciente: ")
				auxText = reader.ReadString('\n')
				nuevaHistoria.Paciente = auxText
				fmt.Println("Ingrese el titulo: ")
				auxText = reader.ReadString('\n')
				nuevaHistoria.Title = auxText
				fmt.Println("Ingrese la historia clinica: ")
				auxText = reader.ReadString('\n')
				nuevaHistoria.Historia = auxText
				nuevaHistoria.IsGenesis = false
				nuevaHistoria.CreationDate = time.Now().Format(time.RFC850)

				url := ipGenesis + "/new"

				var jsonStr = json.newEncoder().Encode(nuevaHistoria)

				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
				req.Header.Set("X-Custom-Header", "myvalue")
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()
				break
			case "2":
				for _, block := range BlockChain.blocks {
					fmt.Printf("Prev. hash: %x\n", block.PrevHash)
					bytes, _ := json.MarshalIndent(block.Data, "", " ")
					fmt.Printf("Data: %v\n", string(bytes))
					fmt.Printf("Hash: %x\n", block.Hash)
					fmt.Println()
				}

			}

		}

	}()

}
