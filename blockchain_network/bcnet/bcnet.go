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
	}
	genesisBlock.Hash = ""
	genesisBlockStr, err := json.Marshal(genesisBlock)
	if err != nil {
		log.Println("Error marshalling genesis block:", err)
	}
	return &BlockChain{
		Blocks:    []Block{genesisBlock},
		BlocksStr: []string{string(genesisBlockStr)},
	}
}

func calculateHash(block Block) string {
	record := fmt.Sprintf("%d%s%v%s", block.Index, block.Timestamp, block.Transactions, block.PrevHash)
	hash := sha256.Sum256([]byte(record))
	return fmt.Sprintf("%x", hash[:])
}

func (bc *BlockChain) AddBlock(transactions []string) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := Block{
		Index:        len(bc.Blocks),
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		PrevHash:     calculateHash(prevBlock),
	}
	newBlock.Hash = calculateHash(newBlock)

	newBlockStr, err := json.Marshal(newBlock)
	if err != nil {
		log.Println("Error marshalling new block:", err)
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
	info map[string][]string // email -> consortium names
}

func NewPeerInfo() *PeerInfo {
	return &PeerInfo{
		info: make(map[string][]string),
	}
}

func (p *PeerInfo) Add(email, consortium string) {
	p.info[email] = append(p.info[email], consortium)
}

func (p *PeerInfo) Remove(email, consortium string) {
	newSlice := make([]string, 0, len(p.info[email]))
	for _, v := range p.info[email] {
		if v != consortium {
			newSlice = append(newSlice, v)
		}
	}
	p.info[email] = newSlice
}

type BlockChainNetwork struct {
	Blockchains map[string]*BlockChain // consortium name -> blockchain
	Consortiums map[string]*Consortium // consortium name -> consortium metadata
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

func parseJWT(jwtKey, tokenStr string) (jwt.MapClaims, error) {
	// Remove "Bearer " prefix if present.
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (bcnet *BlockChainNetwork) InitPeer(req *Request) ([]string, error) {
	claims, err := parseJWT(bcnet.JWT_KEY, req.Token)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	email, _ := claims.GetSubject()
	if consortiums, exists := bcnet.Peerinfo.info[email]; exists {
		return consortiums, nil
	}
	return []string{}, nil
}

func (bcnet *BlockChainNetwork) AddBlockChain(req *Request, transactions []string) error {
	claims, err := parseJWT(bcnet.JWT_KEY, req.Token)
	if err != nil {
		log.Println(err)
		return err
	}

	email, _ := claims.GetSubject()
	consortium, exists := bcnet.Consortiums[req.Consortium]
	if !exists || !consortium.ApprovedPeers[email] {
		return fmt.Errorf("'%s' doesn't exist or '%s' doesn't belong to the consortium '%s'", req.Consortium, email, req.Consortium)
	}

	bcnet.Blockchains[req.Consortium].AddBlock(transactions)
	return nil
}

func (bcnet *BlockChainNetwork) GetBlockChain(req *Request) (*BlockChain, error) {
	claims, err := parseJWT(bcnet.JWT_KEY, req.Token)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	email, _ := claims.GetSubject()
	consortium, exists := bcnet.Consortiums[req.Consortium]
	if !exists || !consortium.ApprovedPeers[email] {
		return nil, fmt.Errorf("'%s' doesn't exist or '%s' doesn't belong to the consortium '%s'", req.Consortium, email, req.Consortium)
	}

	return bcnet.Blockchains[req.Consortium], nil
}

func (bcnet *BlockChainNetwork) ValidateRequest(req *Request) error {
	claims, err := parseJWT(bcnet.JWT_KEY, req.Token)
	if err != nil {
		log.Println(err)
		return err
	}

	email, _ := claims.GetSubject()
	role, _ := claims["role"].(string)

	switch req.Cmd {
	case "participate":
		consortium, exists := bcnet.Consortiums[req.Consortium]
		if !exists {
			return fmt.Errorf("consortium '%s' does not exist", req.Consortium)
		}
		if consortium.ApprovedPeers[email] {
			return fmt.Errorf("already participated")
		}
		consortium.ApprovedPeers[email] = true
		bcnet.Peerinfo.Add(email, req.Consortium)
		return nil

	case "make":
		if role != "Admin" {
			return fmt.Errorf("access denied for employee role")
		}
		if _, exists := bcnet.Consortiums[req.Consortium]; exists {
			return fmt.Errorf("consortium '%s' already exists", req.Consortium)
		}
		bcnet.Consortiums[req.Consortium] = NewConsortium(req.Consortium, email)
		bcnet.Blockchains[req.Consortium] = NewBlockChain()
		// consortium 생성자는 자동으로 참여자로 등록.
		bcnet.Consortiums[req.Consortium].ApprovedPeers[email] = true
		bcnet.Peerinfo.Add(email, req.Consortium)
		return nil

	case "exit":
		consortium, exists := bcnet.Consortiums[req.Consortium]
		if !exists {
			return fmt.Errorf("consortium '%s' does not exist", req.Consortium)
		}
		delete(consortium.ApprovedPeers, email)
		bcnet.Peerinfo.Remove(email, req.Consortium)
		return nil

	default:
		return fmt.Errorf("invalid request command")
	}
}
