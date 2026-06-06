# ChatUI
----
## Description:
Making a GUI or a Clean TUI for when chatting with Ollama models that are running locally.

## MVP
1. we install an uncensored or no safety open weights model using ollama and do ```ollama serve ``` after that model is installed.
2. now we have ollama api present now we just need a backend that can send queries to ollama and give us a good response.
3. Pretty good UI, just something that works.
4. fine tuning it on PDFs.

## Tech stack
golang
frontend : going with bubble tea cli 

## Models in use:
- qwen2.5:3b (or whatever Ollama model you want for chat)
- nomic-embed-text:latest (for embedding text chunks)

## Prerequisites & Setup

1. **Ollama**: Install and run Ollama with the embedding and chat models.
   ```bash
   ollama pull nomic-embed-text
   ollama pull qwen2.5:3b
   ollama serve
   ```
2. **Postgres Vector Database**: Start the local vector store using Docker.
   ```bash
   docker-compose up -d
   ```
3. **Initialize pgvector**: You must create the `vector` extension in the database before running ingestion.
   ```bash
   docker exec -it chattui-pgvector psql -U postgres -d vectordb -c "CREATE EXTENSION IF NOT EXISTS vector;"
   ```

4. **Environment Configuration**: Create a `.env` file at the root. Here is what you need:
   ```env
   EMBED_MODEL=nomic-embed-text
   CHAT_MODEL=qwen2.5:3b
   DATABASE_URL=postgres://postgres:password@localhost:5432/vectordb?sslmode=disable
   DATASET=Test_dataset
   SYSTEM_PROMPT="You are a helpful assistant that answers questions based on the provided context documents..."
   ```

## Running the App

```bash
# To start the interactive TUI application (replaces the old CLI)
go run ./main
```

## Architecture:
### Training from PDFs
PDF → extract text → split → embed → store vectors → similarity search → send context to LLM
- successfully extracted the contents of a pdf from a file.
- now splitting it.
- chunking completed.
- successfully embedded chunks using Ollama's embedding API (`nomic-embed-text`) via `langchaingo`.
- ingestion completed using pgvector with docker i know and it has the vector similarity search thingy installed.

chatting is working now!! you can do `go run ./main` and it actually talks back using the context from ingested PDFs. feels good man.


## How to use the TUI:

The tool runs completely in your terminal. Use `j`/`k` or the arrow keys to navigate the menu and press `enter` to select:

- **Run Chat**: Starts a conversation using your ingested document contexts. Type your question into the multi-line input box and press `enter` to send it. It streams query embeddings, runs similarity searches, and returns the response from Ollama. Press `esc` to go back to the menu.
- **Ingest**: Scans the dataset directory (configured in `.env`) for PDFs, chunks them, generates vector embeddings, and stores them in Postgres. Shows a loading spinner while doing it.
- **Show Files**: Lists all the PDFs in your dataset directory. Scroll with `j`/`k` and hit `enter` to open the selected PDF in your system's default PDF viewer (uses `xdg-open`).
- **Flush DB**: Clears all documents and embeddings from Postgres. It asks you to confirm with `y`/`n` first, then shows a progress bar while doing the clean up.
- **Exit**: Quits the program cleanly.
