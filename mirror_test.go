package mirror

import (
	"testing"

	"go.trulyao.dev/mirror/config"
)

func Test_Fork(t *testing.T) {
	original := New(config.Config{})

	fork := original.Fork(config.Config{
		Enabled:           Bool(true),
		OutputFile:        String("test.ts"),
		UseTypeForObjects: Bool(true),
	}, false)

	if fork.config.EnabledOrDefault() != true {
		t.Errorf("Expected fork.config.EnabledOrDefault() to be true, got %v", fork.config.EnabledOrDefault())
	}

	if fork.config.OutputFileOrDefault() != "test.ts" {
		t.Errorf("Expected fork.config.OutputFileOrDefault() to be \"test.ts\", got %v", fork.config.OutputFileOrDefault())
	}

	if fork.config.UseTypeForObjectsOrDefault() != true {
		t.Errorf("Expected fork.config.UseTypeForObjectsOrDefault() to be true, got %v", fork.config.UseTypeForObjectsOrDefault())
	}

	if fork.config.ExpandObjectTypesOrDefault() != false {
		t.Errorf("Expected fork.config.ExpandObjectTypesOrDefault() to be false, got %v", fork.config.ExpandObjectTypesOrDefault())
	}
}

func Test_Generate(t *testing.T) {
	type Test struct {
		Name      string
		Data      any
		Expected  string
		ExpectErr bool
	}

	tests := []Test{
		{"test integer", 1, "export type int = number;", false},
	}

	m := New(config.Config{Enabled: Bool(true)})

	for _, test := range tests {
		got, err := m.GenerateSingle(test.Data)
		if err != nil && !test.ExpectErr {
			t.Errorf("Expected no error, got %v", err)
		}

		if got != test.Expected {
			t.Errorf("fail(%s): expected %s, got %s", test.Name, test.Expected, got)
		}
	}
}
