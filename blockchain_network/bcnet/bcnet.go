package bcnet

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Block struct {
	Index        int
	Timestamp    string
	Transactions []string
	PrevHash     string
	Hash         string
}

type BlockChain struct {
	Blocks []Block
	mu     sync.RWMutex
}

func NewBlockChain() *BlockChain {
	genesisBlock := Block{
		Index:        0,
		Timestamp:    time.Now().String(),
		Transactions: []string{"Genesis Block"},
		PrevHash:     "",
		Hash:         "",
	}
	genesisBlock.Hash = ""
	return &BlockChain{Blocks: []Block{genesisBlock}}
}

func calculateHash(block Block) string {
	record := fmt.Sprintf("%d%s%v%s", block.Index, block.Timestamp, block.Transactions, block.PrevHash)
	hash := sha256.Sum256([]byte(record))
	return fmt.Sprintf("%x", hash[:])
}

// func ValidBlock(prev, new Block) bool {
// 	if prev.Hash != new.PrevHash {
// 		return false
// 	}
// 	if new.Hash != calculateHash(new) {
// 		return false
// 	}
// 	return true
// }

func (bc *BlockChain) AddBlock(transactions []string) Block {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := Block{
		Index:        len(bc.Blocks),
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		PrevHash:     calculateHash(prevBlock),
		Hash:         "",
	}
	newBlock.Hash = calculateHash(newBlock)
	return newBlock
}

type Consortium struct {
	Name          string
	Creator       string
	ApprovedPeers map[string]bool
}

func NewConsortium(name, email string) *Consortium {
	return &Consortium{
		Name:          name,
		Creator:       email,
		ApprovedPeers: make(map[string]bool),
	}
}

type BlockChainNetwork struct {
	Blockchains map[string]*BlockChain // 컨소시엄 이름 -> 블록체인
	Consortiums map[string]*Consortium // 컨소시엄 이름 -> 컨소시엄 메타데이터터
	JWT_KEY     string
}

func NewBlockChainNetwork() *BlockChainNetwork {
	return &BlockChainNetwork{
		Blockchains: make(map[string]*BlockChain),
		Consortiums: make(map[string]*Consortium),
		JWT_KEY:     os.Getenv("JWT_KEY"),
	}
}

type Request struct {
	token      string
	cmd        string
	consortium string
}

func (bcnet *BlockChainNetwork) GetBlockChain(req *Request) (*BlockChain, error) {
	tokenString := strings.TrimPrefix(req.token, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(bcnet.JWT_KEY), nil
	})
	if err != nil {
		log.Printf("failed to parse token: %v", err)
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email, _ := claims.GetSubject()
		if c, exsit := bcnet.Consortiums[req.consortium]; exsit {
			if c.ApprovedPeers[email] {
				return bcnet.Blockchains[req.consortium], nil
			}
		}
		return nil, fmt.Errorf("'%s' doesn't exsit or '%s' doesn't bleong to the '%s'", req.consortium, email, req.consortium)
	}
	return nil, fmt.Errorf("invalid token")
}

func (bcnet *BlockChainNetwork) ValidateRequest(req *Request) error {
	tokenString := strings.TrimPrefix(req.token, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(bcnet.JWT_KEY), nil
	})
	if err != nil {
		log.Printf("failed to parse token: %v", err)
		return err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email, _ := claims.GetSubject()
		role := claims["role"].(string)
		switch req.cmd {
		case "participate":
			bcnet.Consortiums[req.consortium].ApprovedPeers[email] = true
			return nil
		case "make":
			if role != "Admin" {
				return fmt.Errorf("access denied for employee role")
			}
			bcnet.Consortiums[req.consortium] = NewConsortium(req.consortium, email)
			return nil
		default:
			return fmt.Errorf("invalid request command")
		}
	}
	log.Println("invalid token")
	return fmt.Errorf("invalid token")
}
