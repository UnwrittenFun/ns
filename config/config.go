package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	Path string
	raw  map[string]interface{}
}

func New(dir, filename string) *Config {
	var cfg = &Config{
		Path: filepath.Join(dir, filename),
		raw:  make(map[string]interface{}),
	}

	cfg.Read()

	return cfg
}

func (c *Config) Read() {
	if _, err := os.Stat(c.Path); os.IsNotExist(err) {
		return
	}

	raw, err := ioutil.ReadFile(c.Path)
	if err != nil {
		fmt.Println("Warning: Failed to read config file, ignoring")
		return
	}

	if err := json.Unmarshal(raw, &c.raw); err != nil {
		fmt.Println("Warning: Failed to parse json config, ignoring")
		return
	}
}

func (c *Config) Save() error {
	jsonCfg, err := json.Marshal(c.raw)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(c.Path, jsonCfg, os.ModePerm)
}

func (c *Config) Set(key string, value interface{}) {
	c.raw[key] = value
}

func (c *Config) Get(key string, fallback interface{}) interface{} {
	if val, ok := c.raw[key]; ok {
		return val
	}

	return fallback
}

func (c *Config) GetString(key string, fallback string) string {
	return c.Get(key, fallback).(string)
}
