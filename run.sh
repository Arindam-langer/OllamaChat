#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

echo "Starting OllamaChat..."

# 1. Create a default .env if it doesn't exist
if [ ! -f .env ]; then
  echo "No .env file found. Creating default .env..."
  cat <<EOT > .env
EMBED_MODEL=nomic-embed-text
DATABASE_URL=postgres://postgres:password@localhost:5432/vectordb?sslmode=disable
DATASET=Test_dataset
CHAT_MODEL=qwen2.5:3b
SYSTEM_PROMPT="You are a helpful assistant that answers questions based on the provided context documents. Use the context below to answer the user's question accurately and concisely. If the context doesn't contain relevant information, say so honestly and try your best to help."
EOT
fi

# 2. Start pgvector Docker container if not running
if ! docker ps --format '{{.Names}}' | grep -q "^chattui-pgvector$"; then
  echo "Starting pgvector database container..."
  docker-compose up -d
fi

# 3. Wait for Postgres to be ready
echo "Waiting for database to be ready..."
until docker exec chattui-pgvector pg_isready -U postgres -d vectordb >/dev/null 2>&1; do
  sleep 1
done

# 4. Initialize the pgvector extension
echo "Initializing pgvector extension..."
docker exec -i chattui-pgvector psql -U postgres -d vectordb -c "CREATE EXTENSION IF NOT EXISTS vector;" >/dev/null 2>&1

# 5. Build and launch the TUI application
if [ -n "$SUDO_USER" ]; then
  echo "Launching TUI as user $SUDO_USER..."
  # Pass user's home and environment variables so xdg-open/display work correctly
  sudo -E -u "$SUDO_USER" env HOME="/home/$SUDO_USER" go run ./main
else
  echo "Launching TUI..."
  go run ./main
fi
