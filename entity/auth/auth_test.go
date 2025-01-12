package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAutoMigrateAuthSchema(t *testing.T) {
	err := AutoMigrateAuthSchema()
	assert.NoError(t, err, "failed 'AutoMigrateAuthSchema'")
}
