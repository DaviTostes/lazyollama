# 🦙 lazyollama

**lazyollama** is a terminal-based interface for managing conversations with Ollama models, built in Go.  
It allows you to create, view, and organize your chats in a simple and efficient way — directly from your terminal.

## 🚀 Features

- Start new conversations with Ollama models
- List, select, and delete existing chats
- Easily switch between models
- Minimal and fast CLI interface

## 📦 Installation

You can install `lazyollama` by cloning the repository and building it:

```bash
git clone https://github.com/yourusername/lazyollama.git
cd lazyollama
go build -o lazyollama
```

Then move the binary to your path if you'd like:

```bash
sudo mv lazyollama /usr/local/bin
```

## 🛠 Usage

```bash
lazyollama <command>
```

### Available Commands

| Command            | Description                        |
|--------------------|------------------------------------|
| `new`              | Create a new chat                  |
| `list`             | List the existing chats            |
| `select <id>`      | Select a chat by ID                |
| `delete <id>`      | Delete a chat by ID                |
| `model`            | Show the current active model      |
| `model <name>`     | Change the active model            |
| `help`             | Show more information              |

## 💡 Example

```bash
lazyollama new
lazyollama list
lazyollama select 2
lazyollama model llama3
```

## 🤝 Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---

Happy chatting! 🦙
