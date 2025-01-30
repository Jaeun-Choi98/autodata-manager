package bcnet

import (
	"crypto/sha256"
	"encoding/json"
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
	Blocks    []Block
	BlocksStr []string
	mu        sync.RWMutex
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
	genesisBlockStr, err := json.Marshal(genesisBlock)
	if err != nil {
		log.Println(err)
	}
	return &BlockChain{Blocks: []Block{genesisBlock}, BlocksStr: []string{string(genesisBlockStr)}}
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

func (bc *BlockChain) AddBlock(transactions []string) error {
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
	newBlockStr, err := json.Marshal(newBlock)
	if err != nil {
		log.Println(err)
	}
	bc.Blocks = append(bc.Blocks, newBlock)
	bc.BlocksStr = append(bc.BlocksStr, string(newBlockStr))
	return nil
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

type PeerInfo struct {
	info map[string][]string // email -> 컨소시엄 이름
}

func NewPeerInfo() *PeerInfo {
	return &PeerInfo{
		info: make(map[string][]string),
	}
}

func (p *PeerInfo) Add(email, consortium string) {
	if _, exists := p.info[email]; !exists {
		p.info[email] = make([]string, 0)
	}
	p.info[email] = append(p.info[email], consortium)
}

func (p *PeerInfo) Remove(email, consortium string) {
	newSlice := make([]string, 0)
	for _, v := range p.info[email] {
		if v != consortium {
			newSlice = append(newSlice, v)
		}
	}
	p.info[email] = newSlice
}

type BlockChainNetwork struct {
	Blockchains map[string]*BlockChain // 컨소시엄 이름 -> 블록체인
	Consortiums map[string]*Consortium // 컨소시엄 이름 -> 컨소시엄 메타데이터터
	Peerinfo    *PeerInfo
	JWT_KEY     string
}

func NewBlockChainNetwork() *BlockChainNetwork {
	return &BlockChainNetwork{
		Blockchains: make(map[string]*BlockChain),
		Consortiums: make(map[string]*Consortium),
		Peerinfo:    NewPeerInfo(),
		JWT_KEY:     os.Getenv("JWT_KEY"),
	}
}

type Request struct {
	Token      string
	Cmd        string
	Consortium string
}

func (bcnet *BlockChainNetwork) InitPeer(req *Request) ([]string, error) {
	tokenString := strings.TrimPrefix(req.Token, "Bearer ")
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
		if ret, exists := bcnet.Peerinfo.info[email]; exists {
			return ret, nil
		}
	}
	return []string{}, nil
}

func (bcnet *BlockChainNetwork) AddBlockChain(req *Request, transactions []string) error {
	tokenString := strings.TrimPrefix(req.Token, "Bearer ")
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
		if c, exists := bcnet.Consortiums[req.Consortium]; exists {
			if c.ApprovedPeers[email] {
				bcnet.Blockchains[req.Consortium].AddBlock(transactions)
				return nil
			}
		}
		return fmt.Errorf("'%s' doesn't exists or '%s' doesn't belong to the '%s'", req.Consortium, email, req.Consortium)
	}
	return fmt.Errorf("invalid token")
}

func (bcnet *BlockChainNetwork) GetBlockChain(req *Request) (*BlockChain, error) {
	tokenString := strings.TrimPrefix(req.Token, "Bearer ")
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
		if c, exists := bcnet.Consortiums[req.Consortium]; exists {
			if c.ApprovedPeers[email] {
				return bcnet.Blockchains[req.Consortium], nil
			}
		}
		return nil, fmt.Errorf("'%s' doesn't exists or '%s' doesn't belong to the '%s'", req.Consortium, email, req.Consortium)
	}
	return nil, fmt.Errorf("invalid token")
}

func (bcnet *BlockChainNetwork) ValidateRequest(req *Request) error {
	tokenString := strings.TrimPrefix(req.Token, "Bearer ")

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
		switch req.Cmd {
		case "participate":
			if !bcnet.Consortiums[req.Consortium].ApprovedPeers[email] {
				bcnet.Consortiums[req.Consortium].ApprovedPeers[email] = true
				bcnet.Peerinfo.Add(email, req.Consortium)
				return nil
			}
			return fmt.Errorf("already participated")
		case "make":
			if role != "Admin" {
				return fmt.Errorf("access denied for employee role")
			}
			if _, exists := bcnet.Consortiums[req.Consortium]; !exists {
				bcnet.Consortiums[req.Consortium] = NewConsortium(req.Consortium, email)
				bcnet.Blockchains[req.Consortium] = NewBlockChain()
				bcnet.Consortiums[req.Consortium].ApprovedPeers[email] = true
				bcnet.Peerinfo.Add(email, req.Consortium)
				return nil
			}
			return fmt.Errorf("already existed '%s'", req.Consortium)
		case "exit":
			delete(bcnet.Consortiums[req.Consortium].ApprovedPeers, email)
			bcnet.Peerinfo.Remove(email, req.Consortium)
			return nil
		default:
			return fmt.Errorf("invalid request command")
		}
	}
	log.Println("invalid token")
	return fmt.Errorf("invalid token")
}
