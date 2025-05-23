package db

import (
	"database/sql"
	"encoding/json"
	"lazyollama/model"
	"lazyollama/ollama"
	"os"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

type SQLiteClient struct {
	db *sql.DB
}

func NewSQLiteClient() (*SQLiteClient, error) {
	db, err := sql.Open("sqlite", os.Getenv("LAZY_DIR")+"/database.sqlite")
	if err != nil {
		return nil, err
	}

	return &SQLiteClient{db: db}, nil
}

func (client SQLiteClient) CreateTables() error {
	sql := `
	CREATE TABLE chats(
		id INTEGER PRIMARY KEY,
		createdAt TEXT NOT NULL,
		desc TEXT NOT NULL
	);

	CREATE TABLE messages(
		id INTEGER PRIMARY KEY,
		createdAt TEXT NOT NULL,
		sender TEXT NOT NULL,
		content TEXT NOT NULL,
		chatId INTEGER NOT NULL,
		FOREIGN KEY(chatId) REFERENCES chats(id) ON DELETE CASCADE
	);

	CREATE TABLE model(
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
	);

	CREATE TABLE embeddings(
		id INTEGER PRIMARY KEY,
		text STRING NOT NULL,
		emb STRING NOT NULL,
		chatId INTEGER NOT NULL,
		FOREIGN KEY(chatId) REFERENCES chats(id) ON DELETE CASCADE
	);

	INSERT INTO model(id, name) VALUES(1, "llama3.2") 
	`

	_, err := client.db.Exec(sql)
	return err
}

func (client SQLiteClient) CreateChat(text string) (int64, error) {
	ollama := ollama.Ollama{
		Model: os.Getenv("DEFAULT_MODEL"),
	}

	gen, err := ollama.GenerateChatName(`
Based on the user's first message below, generate a short and relevant description for this conversation. The description should capture the topic or intent of the message.
Respond only with the description — no explanations or extra text.
Keep the description not longer than 5 words


User's message:
		` + text)
	if err != nil {
		panic(err)
	}

	sql := `
		INSERT INTO chats(createdAt, desc)
		VALUES (?, ?);
	`

	result, err := client.db.Exec(sql, time.Now(), gen)
	if err != nil {
		return 0, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return id, err
}

func (client SQLiteClient) FetchChats() ([]model.Chat, error) {
	query := "SELECT * FROM chats"

	rows, err := client.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []model.Chat

	for rows.Next() {
		c := &model.Chat{}

		err := rows.Scan(&c.Id, &c.CreatedAt, &c.Desc)
		if err != nil {
			return nil, err
		}

		chats = append(chats, *c)
	}

	return chats, err
}

func (client SQLiteClient) FetchChatById(id int) (*model.Chat, error) {
	query := "SELECT * FROM chats WHERE id = ?"

	row := client.db.QueryRow(query, id)

	c := &model.Chat{}
	err := row.Scan(&c.Id, &c.CreatedAt, &c.Desc)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (client SQLiteClient) DeleteChat(id int) error {
	query := `
		DELETE FROM chats WHERE id = ?;
		DELETE FROM messages WHERE chatId = ?;
		DELETE FROM embeddings WHERE chatId = ?;
	`

	_, err := client.db.Exec(query, id, id, id)
	return err
}
func (client SQLiteClient) DeleteAllChats() error {
	query := "DELETE FROM chats"

	_, err := client.db.Exec(query)
	return err
}

func (client SQLiteClient) CreateMessage(message model.Message) error {
	sql := `
		INSERT INTO messages(createdAt, sender, content, chatId)
		VALUES (?, ?, ?, ?);
	`

	_, err := client.db.Exec(
		sql,
		message.CreatedAt,
		message.Sender,
		message.Content,
		message.ChatId,
	)
	return err
}

func (client SQLiteClient) FetchMessages(chatId int64) ([]model.Message, error) {
	query := "SELECT * FROM messages WHERE chatId = ? ORDER BY createdAt"

	rows, err := client.db.Query(query, chatId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.Message

	for rows.Next() {
		m := &model.Message{}

		err := rows.Scan(&m.Id, &m.CreatedAt, &m.Sender, &m.Content, &m.ChatId)
		if err != nil {
			return nil, err
		}

		messages = append(messages, *m)
	}

	return messages, err
}

func (client SQLiteClient) FetchModel() (*string, error) {
	query := "SELECT * FROM model WHERE id = 1"

	row := client.db.QueryRow(query)

	m := &model.Model{}
	err := row.Scan(&m.Id, &m.Name)
	if err != nil {
		return nil, err
	}

	return &m.Name, nil
}

func (client SQLiteClient) UpdateModel(name string) error {
	query := `UPDATE model SET name = ?`

	_, err := client.db.Exec(query, name)
	if err != nil {
		return err
	}

	return nil
}

func (client SQLiteClient) CreateEmbedding(emb model.Embedding) error {
	sql := `
		INSERT INTO embeddings(text, emb, chatId)
		VALUES (?, ?, ?);
	`

	jsonEmb, err := json.Marshal(emb.Emb)
	if err != nil {
		return err
	}

	_, err = client.db.Exec(sql, emb.Text, string(jsonEmb), emb.ChatId)
	return err
}

func (client SQLiteClient) FetchEmbeddings(chatId int) ([]model.Embedding, error) {
	query := "SELECT * FROM embeddings WHERE chatId = ?"

	rows, err := client.db.Query(query, chatId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []model.Embedding

	for rows.Next() {
		var jsonEmb string
		e := &model.Embedding{}

		err := rows.Scan(&e.Id, &e.Text, &jsonEmb, &e.ChatId)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(jsonEmb), &e.Emb)
		if err != nil {
			return nil, err
		}

		chats = append(chats, *e)
	}

	return chats, err
}
