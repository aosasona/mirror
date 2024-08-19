package mirror

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"reflect"
	"strings"

	"go.trulyao.dev/mirror/v2/config"
	"go.trulyao.dev/mirror/v2/generator/typescript"
	"go.trulyao.dev/mirror/v2/parser"
	"go.trulyao.dev/mirror/v2/types"
)

type Mirror struct {
	parser types.ParserInterface
	config config.Config
}

var (
	ErrNoSources        = errors.New("no sources provided")
	ErrNoTargetsDefined = errors.New("no targets provided, at least one target must be defined")
)

// New() returns a new instance of the Mirror struct
func New(mirrorConfig config.Config, optionalParser ...types.ParserInterface) *Mirror {
	var p types.ParserInterface

	p = parser.New()
	if len(optionalParser) > 0 {
		p = optionalParser[0]
	}

	m := &Mirror{config: mirrorConfig}
	m.SetParser(p)

	return m
}

// Parser() returns the parser used by the current mirror instance
func (m *Mirror) Parser() types.ParserInterface {
	return m.parser
}

// Config() returns the config used by the current mirror instance
func (m *Mirror) Config() config.Config {
	return m.config
}

// Count() returns the number of sources to generate code for
func (m *Mirror) Count() int {
	return m.parser.Count()
}

// SetEnabled() sets the enabled status of the mirror instance
func (m *Mirror) SetEnabled(enabled bool) *Mirror {
	m.config.Enabled = enabled
	return m
}

// AddSource() adds a source to the list of sources to generate code for
func (m *Mirror) AddSource(s any) *Mirror {
	m.parser.AddSource(reflect.TypeOf(s))
	return m
}

// AddSources() adds multiple sources to the list of sources to generate code for
func (m *Mirror) AddSources(s ...any) *Mirror {
	for _, source := range s {
		m.AddSource(source)
	}

	return m
}

// ResetTargets() resets the targets to an empty list
func (m *Mirror) ResetTargets() *Mirror {
	m.config.Targets = []types.TargetInterface{}
	return m
}

// ResetSources() resets the sources to an empty list
func (m *Mirror) ResetSources() *Mirror {
	m.parser.Reset()
	return m
}

// AddTarget() adds a target to the list of targets to generate code for
func (m *Mirror) AddTarget(t types.TargetInterface) *Mirror {
	m.config.AddTarget(t)
	return m
}

// SetParser() overrides the default parser with a custom parser or any other parser that implements the ParserInterface
func (m *Mirror) SetParser(p types.ParserInterface) *Mirror {
	p.SetConfig(parser.Config{
		FlattenEmbeddedStructs: m.config.FlattenEmbeddedStructs,
		EnableCaching:          m.config.EnableParserCache,
	})

	m.parser = p

	return m
}

// GenerateAndSaveAll() generates code for all sources and saves them to the target files
func (m *Mirror) GenerateAndSaveAll() error {
	if !m.config.Enabled {
		return nil
	}

	if m.Count() == 0 {
		return ErrNoSources
	}

	if len(m.config.Targets) == 0 {
		return ErrNoTargetsDefined
	}

	for _, target := range m.config.Targets {
		var (
			code string
			err  error
		)

		if code, err = m.GenerateforTarget(target); err != nil {
			slog.Error(
				fmt.Sprintf("failed to generate `%s` code", target.Language()),
				slog.String("error", err.Error()),
				slog.String("target", target.Name()),
				slog.String("path", target.Path()),
			)
		}

		if code == "" {
			continue
		}

		if err = m.SaveToFile(target, code); err != nil {
			slog.Error(
				fmt.Sprintf("failed to save `%s` code", target.Language()),
				slog.String("error", err.Error()),
				slog.String("target", target.Name()),
				slog.String("path", target.Path()),
			)
		}
	}

	return nil
}

// GenerateforTarget generates code for a single target returning the fully generated code and an error if any
func (m *Mirror) GenerateforTarget(target types.TargetInterface) (string, error) {
	if !m.config.Enabled {
		return "", nil
	}

	if err := target.Validate(); err != nil {
		return "", err
	}

	dirStat, err := os.Stat(target.Path())
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.New("output path does not exist")
		}

		return "", err
	}

	if !dirStat.IsDir() {
		return "", errors.New("output path is not a directory")
	}

	gen := target.Generator()
	if err = gen.SetParser(m.parser); err != nil {
		return "", err
	}

	generatedTypes, err := gen.GenerateAll()
	if err != nil {
		return "", err
	}

	return target.Header() + "\n" + strings.Join(generatedTypes, "\n\n"), nil
}

// GenerateN generates code for the nth element in the parsed items list
func (m *Mirror) GenerateN(target types.TargetInterface, n int) (string, error) {
	if !m.config.Enabled {
		return "", nil
	}

	if err := target.Validate(); err != nil {
		return "", nil
	}

	gen := target.Generator()
	if err := gen.SetParser(m.parser); err != nil {
		return "", nil
	}

	return gen.GenerateN(n)
}

// SaveToFile saves the generated code to the target file
func (m *Mirror) SaveToFile(target types.TargetInterface, code string) error {
	file, err := os.Create(path.Join(target.Path(), target.Name()))
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(code); err != nil {
		return err
	}

	return nil
}

// Check that all built-in implementations match the interface types
var _ types.ParserInterface = &parser.Parser{}

var (
	_ types.TargetInterface    = &typescript.Config{}
	_ types.GeneratorInterface = &typescript.Generator{}
)
