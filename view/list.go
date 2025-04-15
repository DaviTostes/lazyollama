package view

import (
	"fmt"
	"lazyollama/db"
	"os"
)

func AddList() {
	sqliteClient, err := db.NewSQLiteClient()
	if err != nil {
		panic(err)
	}

	chats, err := sqliteClient.FetchChats()
	if err != nil {
		panic(err)
	}

	for _, c := range chats {
		fmt.Println(c.ToString())
	}
	os.Exit(0)
}
