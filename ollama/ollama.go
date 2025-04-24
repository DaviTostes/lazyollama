package ollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Ollama struct {
	Model string
}

type Response struct {
	Model      string `json:"model"`
	Created_at string `json:"created_at"`
	Response   string `json:"response"`
	Done       bool   `json:"done"`
}

type ResponseChat struct {
	Model      string      `json:"model"`
	Created_at string      `json:"created_at"`
	Message    MessageChat `json:"message"`
	Done       bool        `json:"done"`
}

type MessageChat struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

func (o Ollama) Generate(msg string, messages []MessageChat) (*ResponseChat, *MessageChat, error) {
	userMessage := MessageChat{
		Role:    "user",
		Content: msg,
	}

	body, _ := json.Marshal(map[string]any{
		"model":    o.Model,
		"messages": append(messages, userMessage),
	})

	payload := bytes.NewBuffer(body)

	req, err := http.Post("http://localhost:11434/api/chat", "application/json", payload)
	if err != nil {
		return nil, nil, err
	}

	defer req.Body.Close()

	if req.StatusCode != 200 {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, nil, err
		}

		var errorResponse ErrorResponse
		if err = json.Unmarshal(body, &errorResponse); err != nil {
			return nil, nil, err
		}

		return nil, nil, errors.New(errorResponse.Error)
	}

	var data ResponseChat
	isCode := false

	scanner := bufio.NewScanner(req.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var current ResponseChat
		if err := json.Unmarshal([]byte(line), &current); err != nil {
			return nil, nil, err
		}

		content := ""
		foundBacktick := strings.Contains(current.Message.Content, "``")

		if !isCode && foundBacktick {
			isCode = true
			content += "\033[32m"
		}

		if isCode && foundBacktick && !strings.Contains(content, "\033[32m") {
			isCode = false
			content += "\033[0m"
		}

		content += current.Message.Content

		fmt.Print(content)
		data.Message.Content += content
	}

	return &data, &userMessage, nil
}

func (o Ollama) GenerateChatName(context string) (string, error) {
	body, _ := json.Marshal(map[string]any{
		"model":  o.Model,
		"prompt": context,
		"stream": false,
	})
	payload := bytes.NewBuffer(body)

	req, err := http.Post("http://localhost:11434/api/generate", "application/json", payload)
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	body, err = io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}

	var resp Response
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}

	return resp.Response, err
}

func GetEmbedding(text string) ([]float64, error) {
	body, err := json.Marshal(map[string]string{
		"model":  "nomic-embed-text",
		"prompt": text,
	})
	if err != nil {
		return nil, err
	}

	payload := bytes.NewBuffer(body)

	req, err := http.Post("http://localhost:11434/api/embeddings", "application/json", payload)
	if err != nil {
		return nil, err
	}

	defer req.Body.Close()

	respBody, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var resp EmbeddingResponse
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Embedding, nil
}
