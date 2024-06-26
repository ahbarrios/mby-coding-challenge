package prompt

import (
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/ollama.json
var ollamaGoldenResponse string

func TestOLLama_Acknowledge(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, ollamaGoldenResponse)
	}))
	t.Cleanup(ts.Close)

	a := OLLamaAssistant(AssistantOptions{
		URL: ts.URL,
	})
	got, err := a.Acknowledge("Hello")
	assert.NoError(t, err, "Expect no error from Acknowledge")
	assert.Equal(t, "Hello! How are you today?", got)
}
