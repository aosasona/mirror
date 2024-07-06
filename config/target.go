package config

// A general language interface to make it harder to pass in a wrong language or extend the built-in languages and backends in the future
// There will clearly be neglibile performance impact but it should not matter much here
type LanguageInterface interface {
	Name() string
	Extension() string
}

type (
	typescript string
	swift      string
)

const (
	LangTypescript typescript = "typescript"
	LangSwift      swift      = "swift"
)

func (l typescript) Name() string      { return "typescript" }
func (l typescript) Extension() string { return "ts" }

func (l swift) Name() string      { return "swift" }
func (l swift) Extension() string { return "swift" }

// A Target is the language and file to generate types for
type Target struct {
	Filename        string
	PathToDirectory string
	Language        LanguageInterface
}

func NewTarget(filename, path string, language LanguageInterface) Target {
	return Target{
		Filename:        filename,
		PathToDirectory: path,
		Language:        language,
	}
}
