package mirrormeta

import (
	"testing"
)

// TODO: tests first
func Test_Parse(t *testing.T) {
	type test struct {
		input       string
		expected    ParsedMeta
		expectedErr bool
	}

	tests := []test{
		{
			"type:{ name?: string, email?: string }",
			ParsedMeta{
				Type:     ref("{ name?: string, email?: string }"),
				Skip:     ref(false),
				Name:     ref(""),
				Optional: ref(false),
			},
			false,
		},
		{
			"type:string, name:email, optional:true",
			ParsedMeta{
				Type:     ref("string"),
				Skip:     ref(false),
				Name:     ref("email"),
				Optional: ref(true),
			},
			false,
		},
		{
			"type:string, name:email, optional:0",
			ParsedMeta{
				Type:     ref("string"),
				Skip:     ref(false),
				Name:     ref("email"),
				Optional: ref(false),
			},
			false,
		},
	}
	_ = tests
}
