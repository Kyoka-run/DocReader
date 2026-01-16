package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`

	ChatModel struct {
		APIKey  string `yaml:"api_key"`
		BaseURL string `yaml:"base_url"`
		Model   string `yaml:"model"`
	} `yaml:"chat_model"`

	EmbeddingModel struct {
		APIKey  string `yaml:"api_key"`
		BaseURL string `yaml:"base_url"`
		Model   string `yaml:"model"`
	} `yaml:"embedding_model"`

	Milvus struct {
		Address    string `yaml:"address"`
		Collection string `yaml:"collection"`
	} `yaml:"milvus"`

	FileDir string `yaml:"file_dir"`
}

var (
	cfg  *Config
	once sync.Once
)

func Load(path string) error {
	var err error
	once.Do(func() {
		data, e := os.ReadFile(path)
		if e != nil {
			err = e
			return
		}
		cfg = &Config{}
		err = yaml.Unmarshal(data, cfg)
	})
	return err
}

func Get() *Config {
	return cfg
}
