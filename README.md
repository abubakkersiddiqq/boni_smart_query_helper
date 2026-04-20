# Bino Smart Query Helper

A simple web tool that helps users create better, more detailed queries for Bino (WhatsApp-based local search service).

Type a need (e.g. "restaurant", "hospital", "haircut Koramangala") → get 5 optimized suggestions → click to send directly via WhatsApp to +91 98000 81110.

## Live

[Bino Smart Query Helper](https://binosmartqueryhelper-production.up.railway.app/)

## Features

- AI-generated natural-language queries using OpenRouter (free model)
- Direct wa.me links to send queries to Bino via WhatsApp
- Minimal Go backend + HTML/CSS frontend

## Tech

- Go (net/http, no frameworks)
- OpenRouter API (Llama 3.3 70B free)
- Docker, Railway
- Static CSS served from `/static/`

## Run locally

### With Docker (recommended)

```bash
git clone https://github.com/abubakkersiddiqq/bino
cd bino
docker build -t bino .
docker run -p 8080:8080 -e OPENROUTER_API_KEY=sk-or-v1-your-key bino
```

### Without Docker

```bash
export OPENROUTER_API_KEY=sk-or-v1-your-key
go run main.go
```
