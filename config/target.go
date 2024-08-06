package config

// A general language interface to make it harder to pass in a wrong language or extend the built-in languages and backends in the future
// There will clearly be neglibile performance impact but it should not matter much here
type LanguageConfigInterface interface {
	Name() string
	Extension() string
}

type (
	typescript string
	swift      string
)

// TODO: replace this with the `Config` struct in each generator packages
// A Target is the language and file to generate types for
type Target struct {
	Filename        string
	PathToDirectory string
	Language        LanguageConfigInterface
}

func NewTarget(filename, path string, language LanguageConfigInterface) Target {
	return Target{
		Filename:        filename,
		PathToDirectory: path,
		Language:        language,
	}
}
