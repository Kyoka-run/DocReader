# DocReader

An intelligent document reading assistant powered by RAG (Retrieval-Augmented Generation) technology. Upload your documents, build a private knowledge base, and get accurate answers through natural language queries.

## Features

- **Document Indexing** - Support for Markdown, PDF, and TXT formats with automatic chunking and vectorization
- **Semantic Search** - Vector-based retrieval using Milvus for accurate document matching
- **Intelligent Q&A** - ReAct Agent with LLM for context-aware responses
- **Multi-turn Conversation** - Context memory support for continuous dialogue
- **Streaming Response** - Real-time answer generation via Server-Sent Events (SSE)
- **Tool Calling** - Extensible tool system for enhanced capabilities

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.24 |
| Web Framework | Gin |
| AI Framework | CloudWeGo Eino |
| Vector Database | Milvus |
| LLM | DeepSeek (OpenAI-compatible API) |
| Embedding | Alibaba DashScope text-embedding-v4 |

## Architecture
```
┌─────────────────────────────────────────────────────────────┐
│                        User Request                          │
└─────────────────────────────┬───────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     Gin Web Server                           │
│                  /chat  /chat/stream  /upload                │
└─────────────────────────────┬───────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                 Chat Pipeline (Eino Graph)                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Vector    │  │   Prompt    │  │    ReAct Agent      │  │
│  │  Retrieval  │─▶│  Assembly   │─▶│   (LLM + Tools)     │  │
│  │  (Milvus)   │  │             │  │                     │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                        Response                              │
└─────────────────────────────────────────────────────────────┘
```

## Quick Start

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- API keys for LLM and Embedding services

### 1. Start Milvus
```bash
cd manifest/docker
docker-compose up -d
```

### 2. Configure

Edit `config/config.yaml`:
```yaml
chat_model:
  api_key: "your-llm-api-key"
  base_url: "https://ark.cn-beijing.volces.com/api/v3"
  model: "deepseek-v3-1-terminus"

embedding_model:
  api_key: "your-embedding-api-key"
  base_url: "https://dashscope.aliyuncs.com/compatible-mode/v1"
  model: "text-embedding-v4"
```

### 3. Index Documents

Place your documents in the `docs/` directory, then run:
```bash
go run internal/ai/cmd/knowledge_cmd/main.go
```

### 4. Test Chat (Optional)
```bash
go run internal/ai/cmd/chat_cmd/main.go
```

### 5. Start Server
```bash
go run main.go
```

Server runs on `http://localhost:6872`

## API Reference

### Health Check
```
GET /api/health
```

### Chat
```
POST /api/chat
Content-Type: application/json

{
  "id": "user-123",
  "question": "What is DocReader?"
}
```

Response:
```json
{
  "answer": "DocReader is an intelligent document assistant..."
}
```

### Chat Stream (SSE)
```
POST /api/chat/stream
Content-Type: application/json

{
  "id": "user-123",
  "question": "What is DocReader?"
}
```

Response: Server-Sent Events stream

### Upload Document
```
POST /api/upload
Content-Type: multipart/form-data

file: <your-document>
```

## Usage Example
```bash
# Upload a document
curl -X POST http://localhost:6872/api/upload \
  -F "file=@./my-document.md"

# Ask a question
curl -X POST http://localhost:6872/api/chat \
  -H "Content-Type: application/json" \
  -d '{"id": "user-1", "question": "Summarize the document"}'
```

## License

MIT