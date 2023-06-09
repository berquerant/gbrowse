// Code generated by "goconfig -field GitCommand string -option -output config_generated.go"; DO NOT EDIT.

package git

type ConfigItem[T any] struct {
	modified     bool
	value        T
	defaultValue T
}

func (s *ConfigItem[T]) Set(value T) {
	s.modified = true
	s.value = value
}
func (s *ConfigItem[T]) Get() T {
	if s.modified {
		return s.value
	}
	return s.defaultValue
}
func (s *ConfigItem[T]) Default() T {
	return s.defaultValue
}
func (s *ConfigItem[T]) IsModified() bool {
	return s.modified
}
func NewConfigItem[T any](defaultValue T) *ConfigItem[T] {
	return &ConfigItem[T]{
		defaultValue: defaultValue,
	}
}

type Config struct {
	GitCommand *ConfigItem[string]
}
type ConfigBuilder struct {
	gitCommand string
}

func (s *ConfigBuilder) GitCommand(v string) *ConfigBuilder {
	s.gitCommand = v
	return s
}
func (s *ConfigBuilder) Build() *Config {
	return &Config{
		GitCommand: NewConfigItem(s.gitCommand),
	}
}

func NewConfigBuilder() *ConfigBuilder { return &ConfigBuilder{} }
func (s *Config) Apply(opt ...ConfigOption) {
	for _, x := range opt {
		x(s)
	}
}

type ConfigOption func(*Config)

func WithGitCommand(v string) ConfigOption {
	return func(c *Config) {
		c.GitCommand.Set(v)
	}
}
