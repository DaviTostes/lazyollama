package view

import (
	"bufio"
	"fmt"
	"lazyollama/commands"
	"lazyollama/db"
	"lazyollama/model"
	oll "lazyollama/ollama"
	"lazyollama/vecstore"
	"os"
	"os/exec"
	"strings"
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

	lastChatMessages := []oll.MessageChat{}
	for _, m := range lastMessages {
		lastChatMessages = append(lastChatMessages, oll.MessageChat{
			Role:    m.Sender,
			Content: m.Content,
		})
	}

	messages := []oll.MessageChat{}
	messages = append(messages, lastChatMessages...)

	for {
		fmt.Print("> ")
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			panic(err)
		}

		ollama := oll.Ollama{
			Model: os.Getenv("DEFAULT_MODEL"),
		}

		switch {
		case scanner.Text() == "/leetcodehack":
			gen, userReq, err := commands.GenerateLeetHack(&ollama)
			if err != nil {
				panic(err)
			}

			messages = append(messages, *userReq)
			messages = append(messages, gen.Message)

		case scanner.Text() == "/copycode":
			if len(messages) > 0 {
				message := messages[len(messages)-1]

				err := commands.CopyCode(message)
				if err != nil {
					fmt.Println("Nothing to copy")
					continue
				}

				fmt.Println("Copied to clipboard")
			}

		case strings.Contains(scanner.Text(), "/embed"):
			promptSplited := strings.Split(scanner.Text(), " ")
			if len(promptSplited) < 2 {
				fmt.Println("File path missing.")
				continue
			}

			if firstMessage {
				chatId, err = sqliteClient.CreateChat(scanner.Text())
				if err != nil {
					panic(err)
				}

				firstMessage = false
			}

			err := commands.Embed(chatId, promptSplited[1])
			if err != nil {
				panic(err)
			}

			fmt.Print("text embedded")

			messages = append(messages, oll.MessageChat{Role: "user", Content: promptSplited[1]})
			messages = append(messages, oll.MessageChat{Role: "model", Content: "text embedded"})

		default:
			prompt := scanner.Text()

			if strings.Contains(prompt, "/usecontext") {
				embeddings, err := sqliteClient.FetchEmbeddings(int(chatId))
				if err != nil {
					panic(err)
				}

				if len(embeddings) > 0 {
					queryEmb, err := oll.GetEmbedding(scanner.Text())
					if err != nil {
						panic(err)
					}

					topDocs := vecstore.FindTopKRelevant(embeddings, queryEmb, 2)

					var context strings.Builder
					for _, d := range topDocs {
						context.WriteString(d.Text + "\n")
					}

					prompt = fmt.Sprintf(
						"Use also this context to respond to the prompt:\n\n%s\n\nPrompt: %s",
						context.String(),
						prompt,
					)
				}

				prompt = strings.Replace(prompt, "/usecontenxt", "", 1)
			}

			gen, userReq, err := ollama.Generate(prompt, messages)
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
