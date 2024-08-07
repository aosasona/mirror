package config

// A general language interface to make it harder to pass in a wrong language or extend the built-in languages and backends in the future
// There will clearly be neglibile performance impact but it should not matter much here
type LanguageConfigInterface interface {
	Name() string
	Extension() string
	SetHeader(string)
	AddCustomType(string, string)
}

type (
	typescript string
	swift      string
)

// TODO: replace this with the `Config` struct in each generator packages
// A Target is the language and file to generate types for
type Target struct {
	// Filename is the name of the file to generate
	Filename string

	// PathToDirectory is the path to the directory where the file should be generated (relative to the current working directory or preferably absolute)
	PathToDirectory string

	// Config is the language to generate the types for
	Config LanguageConfigInterface
}

func NewTarget(filename, path string, language LanguageConfigInterface) Target {
	return Target{
		Filename:        filename,
		PathToDirectory: path,
		Config:          language,
	}
}
