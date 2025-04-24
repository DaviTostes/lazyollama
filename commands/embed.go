package commands

import (
	"errors"
	"lazyollama/db"
	"lazyollama/model"
	"lazyollama/ollama"
	"os"
	"path/filepath"
	"strings"
)

func Embed(chatId int64, filePath string) error {
	if !strings.EqualFold(filepath.Ext(filePath), ".txt") {
		return errors.New("Only .txt supported")
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	queryEmb, err := ollama.GetEmbedding(string(content))
	if err != nil {
		return err
	}

	emb := model.Embedding{
		Text:   string(content),
		Emb:    queryEmb,
		ChatId: chatId,
	}

	client, err := db.NewSQLiteClient()
	if err != nil {
		return err
	}

	err = client.CreateEmbedding(emb)
	if err != nil {
		return err
	}

	return nil
}
