package view

import (
	"bufio"
	"fmt"
	"lazyollama/db"
	"lazyollama/model"
	"lazyollama/ollama"
	"os"
	"time"
)

func AddNewChat() {
	firstMessage := true
	var chatId int64

	startChatLoop(firstMessage, chatId)
}

func OpenChat(chatId int64) {
	sqliteClient, err := db.NewSQLiteClient()
	if err != nil {
		panic(err)
	}

	messages, _ := sqliteClient.FetchMessages(chatId)

	for _, m := range messages {
		fmt.Println(m.Content)
	}

	startChatLoop(false, chatId)
}

func startChatLoop(firstMessage bool, chatId int64) {
	scanner := bufio.NewScanner(os.Stdin)
	sqliteClient, err := db.NewSQLiteClient()
	if err != nil {
		panic(err)
	}

	for {
		fmt.Print("> ")
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			panic(err)
		}

		ollama := ollama.Ollama{
			Model: os.Getenv("DEFAULT_MODEL"),
		}

		gen, err := ollama.Generate(scanner.Text())
		if err != nil {
			panic(err)
		}

		fmt.Print("\n")

		if firstMessage {
			chatId, err = sqliteClient.CreateChat(scanner.Text())
			if err != nil {
				panic(err)
			}

			firstMessage = false
		}

		userMessage := model.Message{
			Sender:    "user",
			Content:   scanner.Text(),
			CreatedAt: time.Now().String(),
			ChatId:    chatId,
		}
		go sqliteClient.CreateMessage(userMessage)

		modelMessage := model.Message{
			Sender:    "model",
			Content:   gen.Response,
			CreatedAt: time.Now().String(),
			ChatId:    chatId,
		}
		go sqliteClient.CreateMessage(modelMessage)
	}
}
