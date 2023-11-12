package config

import (
	"encoding/json"
	"io"
	"os"
)

func Default() *Config {
	var c Config
	return &c
}

func Parse(b []byte, c *Config) (*Config, error) {
	if c == nil {
		c = Default()
	}
	if err := json.Unmarshal(b, c); err != nil {
		return nil, err
	}
	return c, nil
}

func parseFile(filePath string, c *Config) (*Config, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return Parse(b, c)
}

func ParseStringOrFile(b string, c *Config) *Config {
	if c == nil {
		c = Default()
	}
	if b == "" {
		return c
	}
	fc, err := parseFile(b, c)
	if err == nil {
		return fc
	}
	sc, err := Parse([]byte(b), c)
	if err != nil {
		return Default()
	}
	return sc
}

type Config struct {
	Phases      []Phase      `json:"phases,omitempty"`
	Definitions []Definition `json:"defs,omitempty"`
}

type Definition struct {
	ID      string   `json:"id"`
	Command []string `json:"cmd"`
}

type Phase string

func NewPhase(value string) Phase { return Phase(value) }
func (p Phase) String() string    { return string(p) }

const (
	Punknown       = "unknown"
	Pbranch        = "branch"
	PdefaultBranch = "default_branch"
	Ptag           = "tag"
	Pcommit        = "commit"
)
