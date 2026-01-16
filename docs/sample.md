# DocReader User Guide

## Introduction
DocReader is an intelligent document assistant powered by RAG (Retrieval-Augmented Generation) technology. It allows you to upload documents and ask questions about their content.

## Features
- Document Upload: Support for Markdown, PDF, TXT formats
- Intelligent Q&A: Ask questions in natural language
- Context Memory: Multi-turn conversation support
- Streaming Response: Real-time answer generation

## How to Use

### Uploading Documents
1. Use the `/api/upload` endpoint to upload your documents
2. Supported formats: .md, .pdf, .txt
3. Documents will be automatically indexed for searching

### Asking Questions
1. Send your question to `/api/chat` or `/api/chat/stream`
2. The system will search relevant documents and generate answers
3. For better results, ask specific questions

## API Reference

### POST /api/chat
Normal chat endpoint for single-turn Q&A.

Request:
```json
{
  "id": "user-123",
  "question": "What is DocReader?"
}
```

### POST /api/chat/stream
Streaming chat endpoint with Server-Sent Events.

### POST /api/upload
Upload documents to knowledge base.

## Troubleshooting

### Common Issues
- If answers are not accurate, try uploading more relevant documents
- For best results, use well-structured markdown documents
- Ensure Milvus vector database is running before use
