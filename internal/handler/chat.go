package handler

import (
	chat_pipeline2 "DocReader/internal/agent/chat_pipeline"
	"DocReader/pkg/common"
	"DocReader/pkg/log_call_back"
	"DocReader/pkg/mem"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	ID       string `json:"id"`
	Question string `json:"question"`
}

type ChatResponse struct {
	Answer string `json:"answer"`
}

// Health check
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Chat handles normal chat requests
func Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userMessage := &chat_pipeline2.UserMessage{
		ID:      req.ID,
		Query:   req.Question,
		History: mem.GetSimpleMemory(req.ID).GetMessages(),
	}

	runner, err := chat_pipeline2.BuildChatAgent(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	out, err := runner.Invoke(c.Request.Context(), userMessage, compose.WithCallbacks(log_call_back.LogCallback(nil)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save to memory
	mem.GetSimpleMemory(req.ID).SetMessages(schema.UserMessage(req.Question))
	mem.GetSimpleMemory(req.ID).SetMessages(schema.SystemMessage(out.Content))

	c.JSON(http.StatusOK, ChatResponse{Answer: out.Content})
}

// ChatStream handles streaming chat requests (SSE)
func ChatStream(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	userMessage := &chat_pipeline2.UserMessage{
		ID:      req.ID,
		Query:   req.Question,
		History: mem.GetSimpleMemory(req.ID).GetMessages(),
	}

	runner, err := chat_pipeline2.BuildChatAgent(c.Request.Context())
	if err != nil {
		c.SSEvent("error", err.Error())
		return
	}

	sr, err := runner.Stream(c.Request.Context(), userMessage, compose.WithCallbacks(log_call_back.LogCallback(nil)))
	if err != nil {
		c.SSEvent("error", err.Error())
		return
	}
	defer sr.Close()

	var fullResponse string

	c.Stream(func(w io.Writer) bool {
		chunk, err := sr.Recv()
		if errors.Is(err, io.EOF) {
			// Save to memory
			if fullResponse != "" {
				mem.GetSimpleMemory(req.ID).SetMessages(schema.UserMessage(req.Question))
				mem.GetSimpleMemory(req.ID).SetMessages(schema.SystemMessage(fullResponse))
			}
			c.SSEvent("done", "Stream completed")
			return false
		}
		if err != nil {
			c.SSEvent("error", err.Error())
			return false
		}

		fullResponse += chunk.Content
		c.SSEvent("message", chunk.Content)
		return true
	})
}

// FileUpload handles file upload to knowledge base
func FileUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Save file
	savePath := filepath.Join(common.FileDir, file.Filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"fileName": file.Filename,
		"filePath": savePath,
		"fileSize": file.Size,
	})
}

func init() {
	// Ensure docs directory exists
	if common.FileDir == "" {
		common.FileDir = "./docs"
	}
	os.MkdirAll(common.FileDir, 0755)
}
