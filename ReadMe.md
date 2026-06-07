# OllamaChat

A local RAG pipeline with a terminal UI. Point it at a folder of PDFs, ingest them, and then ask questions — it finds the relevant chunks and feeds them as context to a locally running Ollama model.

---

## Demo

<video src="assets/demo.mp4" autoplay loop muted playsinline width="100%"></video>

---

## What it does

- Ingests PDFs: extracts text, chunks it, generates vector embeddings via `nomic-embed-text`, and stores them in a pgvector Postgres database
- At query time: embeds your question, runs a similarity search, pulls the top matching chunks, and sends them as context to your chat model
- Runs entirely locally — no API keys, no internet required

---

## Tech stack

- Go + [langchaingo](https://github.com/tmc/langchaingo)
- [BubbleTea](https://github.com/charmbracelet/bubbletea) for the TUI
- pgvector (via Docker) for vector storage
- Ollama for embeddings and chat

## Models

- `nomic-embed-text` — for embedding chunks and queries
- `qwen2.5:3b` — default chat model (swap for any Ollama-compatible model)

---

## Quick Start

The easiest way to get going:

**1. Pull and serve the models**
```bash
ollama pull nomic-embed-text
ollama pull qwen2.5:3b
ollama serve
```

**2. Run the script**
```bash
chmod +x run.sh
./run.sh
```

This will create a default `.env`, spin up the Postgres vector DB, install the pgvector extension, and launch the TUI.

---

## Manual Setup

If you'd rather do it yourself:

**1. Pull the models**
```bash
ollama pull nomic-embed-text
ollama pull qwen2.5:3b
ollama serve
```

**2. Start the vector DB**
```bash
docker-compose up -d
```

**3. Initialize pgvector**
```bash
docker exec -it chattui-pgvector psql -U postgres -d vectordb -c "CREATE EXTENSION IF NOT EXISTS vector;"
```

**4. Create a `.env` file**
```env
EMBED_MODEL=nomic-embed-text
CHAT_MODEL=qwen2.5:3b
DATABASE_URL=postgres://postgres:password@localhost:5432/vectordb?sslmode=disable
DATASET=Test_dataset
SYSTEM_PROMPT="You are a helpful assistant..."
```

**5. Run it**
```bash
go run ./main
```

---

## Using the TUI

Navigate with `j`/`k` or arrow keys, select with `enter`:

- **Run Chat** — ask questions against your ingested PDFs. Type your question and hit `enter`. Streams back a response using retrieved context. `esc` to go back.
- **Ingest** — scans your dataset directory for PDFs, chunks them, embeds them, and stores everything in Postgres. Shows a spinner while it works.
- **Show Files** — lists PDFs in your dataset folder. Hit `enter` on any file to open it with your system's default PDF viewer.
- **Flush DB** — clears all stored documents and embeddings. Asks for `y`/`n` confirmation first.
- **Exit** — quits cleanly.

---

## Architecture

```
PDF → extract text → chunk → embed (nomic-embed-text) → store in pgvector
Query → embed → similarity search → top-k chunks → context + prompt → Ollama → response
```

This is RAG (Retrieval-Augmented Generation) — the model itself isn't modified, it just gets relevant document excerpts injected into each prompt.

---

## Prerequisites

- [Ollama](https://ollama.ai) installed and running
- Docker (for pgvector)
- Go 1.21+
