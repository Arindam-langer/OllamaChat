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
- qwen2.5:3b
- nomic-embed-text:latest

## Prerequisites & Setup

1. **Ollama**: Install and run Ollama with an embedding model.
   ```bash
   ollama pull nomic-embed-text
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

## Running the App

```bash
# To ingest the PDFs into the vector database
go run ./main ingest

# To start the chat
go run ./main chat
```

## Architecture:
### Training from PDFs
PDF → extract text → split → embed → store vectors → similarity search → send context to LLM
- successfully extracted the contents of a pdf from a file.
- now splitting it.
- chunking completed.
- successfully embedded chunks using Ollama's embedding API (`nomic-embed-text`) via `langchaingo`.
-  ingestion completed using pgvector with docker i know and it has the vector similarity search thingy installed.

chatting is working now!! you can do `go run ./main chat` and it actually talks back using the context from ingested PDFs. feels good man.

current to do: need to make cool UI for it using bubble tea TUI.
