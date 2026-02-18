# Bino Smart Query Helper

A simple web tool that helps users create better, more detailed queries for Bino (WhatsApp-based local search service).

Type a need (e.g. "hospital", "haircut Koramangala") → get 5 optimized suggestions → click to send directly via WhatsApp to +91 98000 81110.

## Features
- AI-generated natural-language queries using OpenRouter (free model)
- Direct wa.me links
- Minimal Go backend + HTML/CSS frontend

## Tech
- Go (net/http, no frameworks)
- OpenRouter API (Llama 3.3 70B free)
- Static CSS served from /static/

## Run locally
```bash
export OPENROUTER_API_KEY=sk-or-v1-your-key
go run main.go