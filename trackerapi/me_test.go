package trackerapi

import (
	"bytes"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMe_promptsForUsernameAndPassword(t *testing.T) {
	var byteBuffer bytes.Buffer
	Me(&byteBuffer)
	var prompts = byteBuffer.String()
	assert.Contains(t, prompts, "Username: ")
	assert.Contains(t, prompts, "Password: ")
}