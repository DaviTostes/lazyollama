# lazyollama 🦙

**lazyollama** is a terminal-based interface written in Go to interact with [Ollama](https://ollama.com/) models — keeping your chats organized and accessible directly from the command line.

No heavy UI, no clutter — just a simple, efficient CLI.

## 🚀 Features

- Start and manage multiple chats with Ollama models
- Switch between chats and models on the fly
- Simple and scriptable interface
- Optional in-chat commands for power users

## 📦 Usage

```sh
lazyollama <command>
```

### Available Commands

| Command            | Description                        |
|--------------------|------------------------------------|
| `new`              | Create a new chat                  |
| `list`             | List existing chats                |
| `select <id>`      | Select a chat by ID                |
| `delete <id>`      | Delete a chat                      |
| `model`            | Show the current active model      |
| `model <name>`     | Change the active model            |
| `help`             | Show help information              |

## ⚙️ In-Chat Commands

You can use special commands inside the chat to trigger advanced features:

### `/leetcodehack`

Take a screenshot of a LeetCode problem and get the interpreted solution from the model.

**Dependencies:**

- [`hyprshot`](https://github.com/hyprwm/hyprshot) – for taking the screenshot
- [`tesseract`](https://github.com/tesseract-ocr/tesseract) – for OCR to transcribe the image

### `/copycode`

Automatically copies the **first code block** from the latest response to your clipboard.

**Dependencies:**
- `xclip` (X11) **or** `wl-clip` (Wayland) — for clipboard access

## 🤖 Model Suggestions

- Use **gemma:3b** or similar for fast, lightweight general tasks and casual chat.
- Use **mistral** or **qwen2.5-coder** for coding-heavy tasks like `/leetcodehack`.

## 📥 Installation

Clone the repo and build:

```sh
git clone https://github.com/yourusername/lazyollama.git
cd lazyollama
go build -o lazyollama
```

## ❤️ Contributing

PRs and ideas are welcome. Open an issue or send improvements!

## 📝 License

MIT License