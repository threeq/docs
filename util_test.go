package docs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestColonPath(t *testing.T) {
	assert.Equal(t, "/api/:id", ColonPath("/api/{id}"))
	assert.Equal(t, "/api/:a/:b/:c", ColonPath("/api/{a}/{b}/{c}"))
}
