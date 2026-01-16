package mem

import (
	"sync"

	"github.com/cloudwego/eino/schema"
)

var SimpleMemoryMap = make(map[string]*SimpleMemory)
var mu sync.Mutex

func GetSimpleMemory(id string) *SimpleMemory {
	mu.Lock()
	defer mu.Unlock()
	if mem, ok := SimpleMemoryMap[id]; ok {
		return mem
	} else {
		newMem := &SimpleMemory{
			ID:            id,
			Messages:      []*schema.Message{},
			MaxWindowSize: 6,
		}
		SimpleMemoryMap[id] = newMem
		return newMem
	}
}

type SimpleMemory struct {
	ID            string            `json:"id"`
	Messages      []*schema.Message `json:"messages"`
	MaxWindowSize int
	mu            sync.Mutex
}

func (c *SimpleMemory) SetMessages(msg *schema.Message) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Messages = append(c.Messages, msg)
	if len(c.Messages) > c.MaxWindowSize {
		excess := len(c.Messages) - c.MaxWindowSize
		if excess%2 != 0 {
			excess++
		}
		c.Messages = c.Messages[excess:]
	}
}
func (c *SimpleMemory) GetMessages() []*schema.Message {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.Messages
}
