package mirror

import (
	"errors"

	"go.trulyao.dev/mirror/config"
)

type Sources []any

type Mirror struct {
	config  config.Config
	sources Sources
}

var (
	ErrNoSources        = errors.New("no sources provided")
	ErrNoTargetsDefined = errors.New("no targets provided, at least one target must be defined")
)

func New(c config.Config) *Mirror {
	return &Mirror{config: c}
}

func (m *Mirror) Count() int {
	return len(m.sources)
}

func (m *Mirror) AddSource(s any) {
	m.sources = append(m.sources, s)
}

func (m *Mirror) AddSources(s ...any) {
	m.sources = append(m.sources, s...)
}

func (m *Mirror) GenerateAll() error {
	if len(m.sources) == 0 {
		return ErrNoSources
	}

	if len(m.config.Targets) == 0 {
		return ErrNoTargetsDefined
	}

	// for _, target := range m.config.Targets {
	//
	// }

	return nil
}
