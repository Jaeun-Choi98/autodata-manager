package bcnet

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestBlockChainNetwork(t *testing.T) {
	err := godotenv.Load("../../.env")
	assert.NoError(t, err)
	secretKey := os.Getenv("JWT_KEY")

	generateToken := func(email, role string) string {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":  email,
			"role": role,
			"exp":  time.Now().Add(time.Hour * 12).Unix(),
		})
		tokenString, _ := token.SignedString([]byte(secretKey))
		return "Bearer " + tokenString
	}

	// Initialize BlockChainNetwork
	bcnet := NewBlockChainNetwork()

	// Test data
	adminToken := generateToken("admin@example.com", "Admin")
	peerToken := generateToken("peer@example.com", "Employee")
	consortiumName := "TestConsortium"

	// Test creating a consortium
	t.Run("Create Consortium", func(t *testing.T) {
		req := &Request{
			token:      adminToken,
			cmd:        "make",
			consortium: consortiumName,
		}
		err := bcnet.ValidateRequest(req)
		assert.NoError(t, err)
		assert.Contains(t, bcnet.Consortiums, consortiumName)
		assert.Equal(t, "admin@example.com", bcnet.Consortiums[consortiumName].Creator)
	})

	// Test approving a peer
	t.Run("Approve Peer", func(t *testing.T) {
		req := &Request{
			token:      peerToken,
			cmd:        "participate",
			consortium: consortiumName,
		}
		err := bcnet.ValidateRequest(req)
		assert.NoError(t, err)
		assert.True(t, bcnet.Consortiums[consortiumName].ApprovedPeers["peer@example.com"])
	})

	// Test getting a blockchain
	t.Run("Get Blockchain", func(t *testing.T) {
		bcnet.Blockchains[consortiumName] = NewBlockChain()
		req := &Request{
			token:      peerToken,
			consortium: consortiumName,
		}
		blockchain, err := bcnet.GetBlockChain(req)
		assert.NoError(t, err)
		assert.NotNil(t, blockchain)
		assert.Equal(t, 1, len(blockchain.Blocks)) // Genesis block exists
	})

	// Test invalid consortium access
	t.Run("Invalid Consortium Access", func(t *testing.T) {
		req := &Request{
			token:      peerToken,
			consortium: "InvalidConsortium",
		}
		_, err := bcnet.GetBlockChain(req)
		assert.Error(t, err)
	})
}
