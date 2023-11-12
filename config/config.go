package config

import "encoding/json"

func Default() *Config {
	var c Config
	return &c
}

func Parse(b []byte) (*Config, error) {
	c := Default()
	if err := json.Unmarshal(b, c); err != nil {
		return nil, err
	}
	return c, nil
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
