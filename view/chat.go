package view

import (
	"bufio"
	"fmt"
	"lazyollama/commands"
	"lazyollama/db"
	"lazyollama/model"
	"lazyollama/ollama"
	"os"
	"os/exec"
	"time"
)

func startOllama() {
	fmt.Println("Starting Ollama...")
	cmd := exec.Command("ollama", "serve")
	go cmd.CombinedOutput()
}

func AddNewChat() {
	startOllama()

	firstMessage := true
	var chatId int64

	startChatLoop(firstMessage, chatId, nil)
}

func OpenChat(chatId int64) {
	startOllama()

	sqliteClient, err := db.NewSQLiteClient()
	if err != nil {
		panic(err)
	}

	messages, _ := sqliteClient.FetchMessages(chatId)

	for _, m := range messages {
		fmt.Println(m.Content)
	}

	startChatLoop(false, chatId, messages)
}

func startChatLoop(firstMessage bool, chatId int64, lastMessages []model.Message) {
	scanner := bufio.NewScanner(os.Stdin)
	sqliteClient, err := db.NewSQLiteClient()
	if err != nil {
		panic(err)
	}

	lastChatMessages := []ollama.MessageChat{}
	for _, m := range lastMessages {
		lastChatMessages = append(lastChatMessages, ollama.MessageChat{
			Role:    m.Sender,
			Content: m.Content,
		})
	}

	messages := []ollama.MessageChat{}
	messages = append(messages, lastChatMessages...)

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

		switch scanner.Text() {
		case "/leetcodehack":
			gen, userReq, err := commands.GenerateLeetHack(&ollama)
			if err != nil {
				panic(err)
			}

			messages = append(messages, *userReq)
			messages = append(messages, gen.Message)
		case "/copycode":
			if len(messages) > 0 {
				message := messages[len(messages)-1]

				err := commands.CopyCode(message)
				if err != nil {
					fmt.Println("Nothing to copy")
					continue
				}

				fmt.Println("Copied to clipboard")
			}
		default:
			gen, userReq, err := ollama.Generate(scanner.Text(), messages)
			if err != nil {
				panic(err)
			}

			messages = append(messages, *userReq)
			messages = append(messages, gen.Message)
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
			Content:   messages[len(messages)-1].Content,
			CreatedAt: time.Now().String(),
			ChatId:    chatId,
		}
		go sqliteClient.CreateMessage(modelMessage)
	}
}
