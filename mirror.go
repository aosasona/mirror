package mirror

import (
	"errors"

	"go.trulyao.dev/mirror/config"
	"go.trulyao.dev/mirror/generator/typescript"
	"go.trulyao.dev/mirror/parser"
	"go.trulyao.dev/mirror/types"
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

func (m *Mirror) AddTarget(t types.TargetInterface) {
	m.config.AddTarget(t)
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

// func (m *Mirror) hasDuplicates() bool {
// }

// Check that all built-in types match the interface types
var _ types.ParserInterface = &parser.Parser{}

var _ types.TargetInterface = &typescript.Config{}

var _ types.GeneratorInterface = &typescript.Generator{}
