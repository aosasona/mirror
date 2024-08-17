package mirror

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"go.trulyao.dev/mirror/config"
	"go.trulyao.dev/mirror/generator/typescript"
	"go.trulyao.dev/mirror/parser"
	"go.trulyao.dev/mirror/types"
)

type Mirror struct {
	parser types.ParserInterface
	config config.Config
}

var (
	ErrNoSources        = errors.New("no sources provided")
	ErrNoTargetsDefined = errors.New("no targets provided, at least one target must be defined")
)

func New(mirrorConfig config.Config, optionalParser ...types.ParserInterface) *Mirror {
	var p types.ParserInterface

	p = parser.New()
	if len(optionalParser) > 0 {
		p = optionalParser[0]
	}

	return &Mirror{config: mirrorConfig, parser: p}
}

func (m *Mirror) Count() int {
	return m.parser.Count()
}

func (m *Mirror) AddSource(s any) {
	m.parser.AddSource(reflect.TypeOf(s))
}

func (m *Mirror) AddSources(s ...any) {
	for _, source := range s {
		m.AddSource(source)
	}
}

func (m *Mirror) AddTarget(t types.TargetInterface) {
	m.config.AddTarget(t)
}

func (m *Mirror) OverrideParser(p types.ParserInterface) {
	slog.Info("overriding parser")
	m.parser = p
}

func (m *Mirror) GenerateAll() error {
	if m.Count() == 0 {
		return ErrNoSources
	}

	if len(m.config.Targets) == 0 {
		return ErrNoTargetsDefined
	}

	for _, target := range m.config.Targets {
		if err := m.GenerateforTarget(target); err != nil {
			slog.Error(
				fmt.Sprintf("failed to generate `%s` code", target.Language()),
				slog.String("error", err.Error()),
				slog.String("target", target.Name()),
				slog.String("path", target.Path()),
			)
		}
	}

	return nil
}

func (m *Mirror) GenerateforTarget(target types.TargetInterface) error {
	if err := target.Validate(); err != nil {
		return err
	}

	dirStat, err := os.Stat(target.Path())
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("output path does not exist")
		}

		return err
	}

	if !dirStat.IsDir() {
		return errors.New("output path is not a directory")
	}

	// TODO: complete generation logic
	// fullPath := path.Join(target.Path(), target.Name())

	return nil
}

// Check that all built-in implementations match the interface types
var _ types.ParserInterface = &parser.Parser{}

var (
	_ types.TargetInterface    = &typescript.Config{}
	_ types.GeneratorInterface = &typescript.Generator{}
)
