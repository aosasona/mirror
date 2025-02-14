package mirrormeta

import (
	"testing"
)

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
			"type:{ name?: string, email?: string], skip: true ",
			ParsedMeta{},
			true,
		},
		{
			"type:(string, skip: true,",
			ParsedMeta{},
			true,
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
		{
			"name:email, optional:0",
			ParsedMeta{
				Type:     ref(""),
				Skip:     ref(false),
				Name:     ref("email"),
				Optional: ref(false),
			},
			false,
		},
		{
			"nme:email, optional:0",
			ParsedMeta{},
			false,
		},
		{
			"",
			ParsedMeta{},
			false,
		},
		{
			"type:string, name:email, optional:0, skip:1",
			ParsedMeta{
				Type:     ref("string"),
				Skip:     ref(true),
				Name:     ref("email"),
				Optional: ref(false),
			},
			false,
		},
		{
			"name:email",
			ParsedMeta{Name: ref("email")},
			false,
		},
		{
			"optional:true",
			ParsedMeta{Optional: ref(true)},
			false,
		},
		{
			"skip:true",
			ParsedMeta{Skip: ref(true)},
			false,
		},
	}

	for _, tc := range tests {
		p := NewMetaParser(tc.input)
		meta, err := p.Parse()

		if tc.expectedErr && err == nil {
			t.Errorf("Expected error but got none")
		} else if !tc.expectedErr && err != nil {
			t.Errorf("Expected no error but got %v", err)
		} else if !tc.expectedErr && err == nil {
			if deref(meta.Name) != deref(tc.expected.Name) {
				t.Errorf("Expected name `%v` but got `%v`", deref(tc.expected.Name), deref(meta.Name))
			}

			if deref(meta.Type) != deref(tc.expected.Type) {
				t.Errorf("Expected type `%v` but got `%v`", deref(tc.expected.Type), deref(meta.Type))
			}

			if deref(meta.Skip) != deref(tc.expected.Skip) {
				t.Errorf("Expected skip `%v` but got `%v`", deref(tc.expected.Skip), deref(meta.Skip))
			}

			if deref(meta.Optional) != deref(tc.expected.Optional) {
				t.Errorf("Expected optional `%v` but got `%v`", deref(tc.expected.Optional), deref(meta.Optional))
			}
		}
	}
}
