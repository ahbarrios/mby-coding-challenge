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

//go:embed testdata/transformers.json
var transformersGoldenResponse string

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

func TestTransformers_Acknowledge(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, transformersGoldenResponse)
	}))
	t.Cleanup(ts.Close)

	a := FAQCanadaImmigration(AssistantOptions{
		URL: ts.URL,
	})
	got, err := a.Acknowledge("Visa")
	assert.NoError(t, err, "Expect no error from Acknowledge")
	assert.Equal(t, "Resubmit documents by mail, print for paper-based programs, or fill out a new IMM1344 for family sponsorship.", got)
}
