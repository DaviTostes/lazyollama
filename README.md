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

### `/embed`

Embed a `.txt` file into the current chat. This file will be used as part of the context.

### `/usecontext`

Enable the model to use the embedded context when generating a response. This triggers the RAG (Retrieval-Augmented Generation) process using your attached `.txt` files.

## How Context Embedding Works (RAG)

lazyollama supports basic Retrieval-Augmented Generation (RAG) using local .txt files.

When you run /embed, the file is split into chunks and embedded using a language model. These embeddings are stored locally.

When you run /usecontext, the tool retrieves the most relevant chunks using cosine similarity — it compares the meaning of your current message with the embedded content to find related pieces.

It uses a top-K strategy (e.g., top 3 most similar chunks) to pick what goes into the final prompt. This helps the model answer based on your documents, like personal notes or guides.

## 🤖 Model Suggestions

- Use **gemma:3b** or similar for fast, lightweight general tasks and casual chat.
- Use **mistral** or **qwen2.5-coder** for coding-heavy tasks like `/leetcodehack`.

## 📥 Installation

Clone the repo and build:

```sh
git clone https://github.com/davitostes/lazyollama.git
cd lazyollama
go build -o lazyollama
```

## ❤️ Contributing

PRs and ideas are welcome. Open an issue or send improvements!

## 📝 License

MIT License
