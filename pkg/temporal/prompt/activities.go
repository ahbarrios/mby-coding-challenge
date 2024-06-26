package prompt

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type message struct {
	Role    string   `json:"role"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
}

type payload struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type response struct {
	Model     string  `json:"model"`
	CreatedAt string  `json:"created_at"`
	Message   message `json:"message"`
	Done      bool    `json:"done"`
}

// Assistant Activity that helps to acknowledge chat messages connecting to LLMs model server
// that use OpenAI architecture for *text generation*.
//
// The only supported API for challenge purpose is `api/chat` [1].
//
// E.g. Servers
// - Ollama
// - OpenAI
// - Transformers
//
// [1]: https://platform.openai.com/docs/api-reference/chat
type Assistant interface {
	// Acknowledge activity will be used as OpenAI compatible client to call prompt servers
	// and obtain a AI response from anyone of it.
	Acknowledge(string) (string, error)
}

// AssistantOptions will be use to configure any [Assistant] implementation
type AssistantOptions struct {
	// URL the URL of the prompt server that support OpenAI Restful API
	URL string
	// APIKey the APIKey of the authenticated account that will be used for connecting
	APIKey string
	// Model the LLMs pretrained models that support text-generation
	Model string
}

// body build a [http] request payload for `api/chat` OpenAI alike server with text-only payload
//
// TODO support history API call and persisted chats
func (cfg *AssistantOptions) body(msg string) (io.Reader, error) {
	bs, err := json.Marshal(&payload{
		Model: cfg.Model,
		Messages: []message{
			{
				Role:    "user",
				Content: msg,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(bs), nil
}

// parse will helps to parse any [response] from `api/chat` OpenAI alike server [http.Response]
func (cfg *AssistantOptions) parse(body io.ReadCloser) (*response, error) {
	bs, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var res response
	if err := json.Unmarshal(bs, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

type Transformers struct {
	AssistantOptions
}

// Acknowledge this is a custom transformers implementation using HuggingFace
// Docker space made by me for this challenge [1].
//
// - [1]: https://huggingface.co/spaces/ahbarrios/faq-canada-immigration
func (m *Transformers) Acknowledge(q string) (string, error) {
	path, err := url.Parse(m.URL)
	if err != nil {
		return "", err
	}

	values := path.Query()
	values.Add("text", q)
	path.RawQuery = values.Encode()

	res, err := http.Get(path.String())
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	bs, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var out struct {
		Response string `json:"output"`
	}
	if err := json.Unmarshal(bs, &out); err != nil {
		return "", err
	}

	return out.Response, nil
}

// FAQCanadaImmigration it creates and [Transformers] compatible [Assistant]
//
// - cfg.URL is required
func FAQCanadaImmigration(cfg AssistantOptions) Assistant {
	if cfg.URL == "" {
		panic("URL is not present")
	}
	return &Transformers{cfg}
}

// OLLama implements [Assistant] for OLLama open source LLMs server
type OLLama struct {
	AssistantOptions
}

func (m *OLLama) Acknowledge(q string) (string, error) {
	p, err := m.body(q)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, m.URL, p)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	rs, err := m.parse(res.Body)
	if err != nil {
		return "", err
	}

	if rs == nil || rs.Message.Content == "" {
		return "", errors.New("no message return from server")
	}
	return rs.Message.Content, nil
}

// OLLamaAssistant it creates and [OLLama] compatible [Assistant]
//
// - cfg.URL is required
func OLLamaAssistant(cfg AssistantOptions) Assistant {
	if cfg.URL == "" {
		panic("URL is not present")
	}
	if cfg.Model == "" {
		cfg.Model = "llama3"
	}
	return &OLLama{cfg}
}
