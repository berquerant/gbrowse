package parse

import (
	"fmt"
	"strconv"
	"strings"
)

//go:generate go tool dataclass -type "InternalTarget" -field "Path string|Linum int" -output parameter_dataclass_generated.go

type Target struct {
	value InternalTarget
}

func (t *Target) Path() string {
	return t.value.Path()
}

func (t *Target) Linum() (int, bool) {
	v := t.value.Linum()
	return v, v > 0
}

func (t *Target) String() string {
	if t.value.Linum() < 1 {
		return t.value.Path()
	}
	return fmt.Sprintf("%s:%d", t.value.Path(), t.value.Linum())
}

func NewTarget(path string, linum int) *Target {
	return &Target{
		value: NewInternalTarget(path, linum),
	}
}

func NewPathTarget(path string) *Target {
	return NewTarget(path, -1)
}

func ReadTarget(value string) (*Target, error) {
	if xs := strings.SplitN(value, ":", 2); len(xs) == 2 {
		linum, err := strconv.Atoi(xs[1])
		if err != nil {
			return nil, fmt.Errorf("invalid target %s, %w", value, err)
		}
		return NewTarget(xs[0], linum), nil
	}

	return NewPathTarget(value), nil
}
