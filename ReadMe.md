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
frontend : dont know


## Architecture:
PDF → extract text → split → embed → store vectors → similarity search → send context to LLM