package ollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
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

func (o Ollama) Generate(msg string) (Response, error) {
	body, _ := json.Marshal(map[string]any{
		"model":  o.Model,
		"prompt": msg,
	})
	payload := bytes.NewBuffer(body)

	req, err := http.Post("http://localhost:11434/api/generate", "application/json", payload)
	if err != nil {
		return Response{}, err
	}
	defer req.Body.Close()

	data := Response{
		Model:      "",
		Created_at: time.Now().String(),
		Response:   "",
		Done:       true,
	}

	isCode := false

	scanner := bufio.NewScanner(req.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var current Response
		if err := json.Unmarshal([]byte(line), &current); err != nil {
			return Response{}, err
		}

		if !isCode && strings.Contains(current.Response, "```") {
			isCode = true
			current.Response = "\033[32m"
		} else if isCode && strings.Contains(current.Response, "```") {
			isCode = false
			current.Response = "\033[0m"
		}

		fmt.Print(current.Response)
		data.Response += current.Response
	}

	return data, nil
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
