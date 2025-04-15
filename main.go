package main

import (
	"flag"
	"fmt"
	"lazyollama/db"
	"lazyollama/view"
	"os"
	"strconv"
)

var usage = `Usage: lazyollama <command>

Avaliable Commands:
  new		create a new chat
  list 	    	list the existing chats
  select <id> 	select a chat
  delete <id> 	delete a chat
  model         show the current active model
  model <name> 	change the active model
  help 		show more information
`

var sqliteClient *db.SQLiteClient

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		showUsage()
	}

	var err error
	sqliteClient, err = db.NewSQLiteClient()
	if err != nil {
		panic(err)
	}

	_, err = os.Stat("database.sqlite")
	if err != nil {
		sqliteClient.CreateTables()
	}

	model, err := sqliteClient.FetchModel()
	if err != nil {
		panic(err)
	}

	err = os.Setenv("DEFAULT_MODEL", *model)
	if err != nil {
		panic(err)
	}

	option := args[0]
	switch option {
	case "", "help":
		showUsage()

	case "new":
		view.AddNewChat()

	case "list":
		view.AddList()

	case "delete":
		if len(args) < 2 {
			showUsage()
		}

		if args[1] == "all" {
			sqliteClient.DeleteAllChats()
			os.Exit(0)
		}

		id := getId(args)

		sqliteClient.DeleteChat(id)

	case "select":
		id := getId(args)

		chat, err := sqliteClient.FetchChatById(id)
		if err != nil {
			fmt.Println("Not found")
			os.Exit(0)
		}

		view.OpenChat(int64(chat.Id))

	case "model":
		model := getModel(args)

		sqliteClient.UpdateModel(model)
		fmt.Println(model)
	}
}

func showUsage() {
	fmt.Print(usage)
	os.Exit(0)
}

func getId(args []string) int {
	id, err := strconv.Atoi(args[1])
	if err != nil {
		panic(err)
	}

	return id
}

func getModel(args []string) string {
	if len(args) < 2 {
		fmt.Println(os.Getenv("DEFAULT_MODEL"))
		os.Exit(0)
	}

	return args[1]
}
