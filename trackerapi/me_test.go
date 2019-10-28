package trackerapi

import (
	"os"
	"bufio"
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

var defaultTestClient = NewTestClient(func(req *http.Request) *http.Response {
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),
		Header:     make(http.Header),
	}
})

func TestMe_promptsForUsernameAndPassword(t *testing.T) {
	var tempDirectory, _ = ioutil.TempDir(".", "test")
	var outputByteBuffer bytes.Buffer
	var inputByteBuffer bytes.Buffer

	Me(&outputByteBuffer, bufio.NewReader(&inputByteBuffer), defaultTestClient, tempDirectory)
	var prompts = outputByteBuffer.String()
	assert.Contains(t, prompts, "Username: ")
	assert.Contains(t, prompts, "Password: ")
	os.RemoveAll(tempDirectory)
}

func TestMe_submitsSuppliedUsernameAndPasswordToLogin(t *testing.T) {
	var tempDirectory, _ = ioutil.TempDir(".", "test")
	var outputByteBuffer bytes.Buffer
	var inputByteBuffer bytes.Buffer
	inputByteBuffer.WriteString("someUserName\nsomePassword\n")
	client := NewTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "https://www.pivotaltracker.com/services/v5/me", req.URL.String())
		assert.Equal(t, "Basic "+base64.StdEncoding.EncodeToString([]byte("someUserName:somePassword")), req.Header.Get("Authorization"))
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),
			Header:     make(http.Header),
		}
	})
	Me(&outputByteBuffer, bufio.NewReader(&inputByteBuffer), client, tempDirectory)
	os.RemoveAll(tempDirectory)
}

func TestMe_writesTheLoginTokenToAFile(t *testing.T) {
	var tempDirectory, _ = ioutil.TempDir(".", "test")
	var outputByteBuffer bytes.Buffer
	var inputByteBuffer bytes.Buffer
	const apiToken = "someApiToken"
	client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString("{\"api_token\": \"" + apiToken + "\"}")),
			Header:     make(http.Header),
		}
	})
	Me(&outputByteBuffer, bufio.NewReader(&inputByteBuffer), client, tempDirectory)
	var fileContents, _ = ioutil.ReadFile(tempDirectory + "/.tracker")
	assert.Equal(t, apiToken, string(fileContents))
	os.RemoveAll(tempDirectory)
}
