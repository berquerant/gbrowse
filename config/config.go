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
	Phases []Phase `json:"phases,omitempty"`
}

type Phase string

func (p Phase) String() string { return string(p) }

const (
	Punknown       = "unknown"
	Pbranch        = "branch"
	PdefaultBranch = "default_branch"
	Ptag           = "tag"
	Pcommit        = "commit"
)

func NewPhase(value string) Phase {
	switch value {
	case Pbranch, PdefaultBranch, Ptag, Pcommit:
		return Phase(value)
	default:
		return Punknown
	}
}
