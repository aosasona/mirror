package mirror

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"go.trulyao.dev/mirror/config"
	"go.trulyao.dev/mirror/generator"
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

// for convenience
type Config = config.Config

func New(c config.Config) (*Mirror, error) {
	if len(c.Targets) == 0 {
		return nil, ErrNoTargetsDefined
	}

	return &Mirror{config: c}, nil
}

func (m *Mirror) Count() int {
	return len(m.sources)
}

func (m *Mirror) AddSource(s any) {
	m.sources = append(m.sources, s)
}

// TODO: fix this
func (m *Mirror) Commit(output string) error {
	if !m.config.Enabled {
		return nil
	}

	if len(m.config.Targets) == 0 {
		return ErrNoTargetsDefined
	}

	// TODO: compare hashes instead
	if m.areSameBytesContent(output) {
		fmt.Println("No changes detected, skipping...")
		return nil
	}

	err := os.WriteFile(m.config.OutputFileOrDefault(), []byte(output), 0644)
	if err != nil {
		return err
	}

	m.sources = Sources{}
	return nil
}

func (m *Mirror) Generate() (string, error) {
	if !m.config.EnabledOrDefault() {
		return "", nil
	}

	var output string

	if len(m.sources) == 0 {
		return "", ErrNoSources
	}

	gn := generator.NewGenerator(generator.Opts{
		UseTypeForObjects:     m.config.UseTypeForObjectsOrDefault(),
		ExpandStructs:         m.config.ExpandObjectTypesOrDefault(),
		PreferUnknown:         m.config.PreferUnknownOrDefault(),
		AllowUnexportedFields: m.config.AllowUnexportedFieldsOrDefault(),
	})

	for _, src := range m.sources {
		result := gn.Generate(src)
		output += result + ";\n\n"
	}

	output = FileHeader + "\n" + strings.TrimSpace(output)

	return output, nil
}

// This is mainly for testing purposes but you can use it to generate a single type
func (m *Mirror) GenerateSingle(src any) (string, error) {
	if !m.config.EnabledOrDefault() {
		return "", nil
	}

	gn := generator.NewGenerator(generator.Opts{
		UseTypeForObjects:     m.config.UseTypeForObjectsOrDefault(),
		ExpandStructs:         m.config.ExpandObjectTypesOrDefault(),
		PreferUnknown:         m.config.PreferUnknownOrDefault(),
		AllowUnexportedFields: m.config.AllowUnexportedFieldsOrDefault(),
	})

	return gn.Generate(src) + ";", nil
}

func (m *Mirror) Execute(log ...bool) error {
	output, err := m.Generate()
	if err != nil {
		return err
	}

	if len(log) > 0 && log[0] {
		fmt.Println(output)
	}

	return m.Commit(output)
}

// Calling Register will register the sources passed to it (doesn't replace the existing sources)
func (m *Mirror) Register(sources ...any) error {
	if !m.config.EnabledOrDefault() {
		return nil
	}

	if len(sources) == 0 {
		return ErrNoSources
	}

	m.sources = append(m.sources, sources...)

	return nil
}
