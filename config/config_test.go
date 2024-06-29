package config

import (
	"testing"

	"go.trulyao.dev/mirror/helper"
)

func Test_MergeConfig(t *testing.T) {
	type TestCase struct {
		original Config
		target   Config
		expected Config
	}

	testCases := []TestCase{
		{
			original: Config{
				Enabled: helper.Bool(true),
				Targets: []Target{
					{"test.swift", "bar/baz", LanguageSwift},
					{"test.ts", "foo/bar", LanguageTypescript},
				},
				ExpandObjectTypes:     helper.Bool(true),
				UseTypeForObjects:     helper.Bool(true),
				PreferUnknown:         helper.Bool(true),
				AllowUnexportedFields: helper.Bool(true),
				IndentationType:       Tab,
				SpaceCount:            helper.Int(4),
			},

			target: Config{
				Enabled: helper.Bool(false),
				Targets: []Target{
					{"test.ts", "foo2/bar", LanguageTypescript},
				},
				ExpandObjectTypes:     helper.Bool(false),
				UseTypeForObjects:     helper.Bool(true),
				PreferUnknown:         helper.Bool(true),
				AllowUnexportedFields: helper.Bool(false),
				IndentationType:       Space,
				SpaceCount:            helper.Int(4),
			},

			expected: Config{
				Enabled: helper.Bool(false),
				Targets: []Target{
					{"test.swift", "bar/baz", LanguageSwift},
					{"test.ts", "foo2/bar", LanguageTypescript},
				},
				ExpandObjectTypes:     helper.Bool(false),
				UseTypeForObjects:     helper.Bool(true),
				PreferUnknown:         helper.Bool(true),
				AllowUnexportedFields: helper.Bool(false),
				IndentationType:       Space,
				SpaceCount:            helper.Int(4),
			},
		},
	}

	for _, tc := range testCases {
		result := tc.original.Merge(tc.target)

		// Since most of the things in the config are pointers, we can't use reflect.DeepEqual and need to compare manually

		if *result.Enabled != *tc.expected.Enabled {
			t.Errorf("Expected Enabled to be %v, got %v", *tc.expected.Enabled, *result.Enabled)
		}

		if len(result.Targets) != len(tc.expected.Targets) {
			t.Errorf("Expected %d targets, got %d", len(tc.expected.Targets), len(result.Targets))
		}

		// Targes are order dependent to be tested here due to the way they are rebuilt during a merge, so changing the order of the targets will cause this test to fail
		for i, target := range result.Targets {
			if target != tc.expected.Targets[i] {
				t.Errorf("Expected target %v, got %v", tc.expected.Targets[i], target)
			}
		}

		if *result.ExpandObjectTypes != *tc.expected.ExpandObjectTypes {
			t.Errorf("Expected ExpandObjectTypes to be %v, got %v", *tc.expected.ExpandObjectTypes, *result.ExpandObjectTypes)
		}

		if *result.UseTypeForObjects != *tc.expected.UseTypeForObjects {
			t.Errorf("Expected UseTypeForObjects to be %v, got %v", *tc.expected.UseTypeForObjects, *result.UseTypeForObjects)
		}

		if *result.PreferUnknown != *tc.expected.PreferUnknown {
			t.Errorf("Expected PreferUnknown to be %v, got %v", *tc.expected.PreferUnknown, *result.PreferUnknown)
		}

		if *result.AllowUnexportedFields != *tc.expected.AllowUnexportedFields {
			t.Errorf("Expected AllowUnexportedFields to be %v, got %v", *tc.expected.AllowUnexportedFields, *result.AllowUnexportedFields)
		}

		if result.IndentationType != tc.expected.IndentationType {
			t.Errorf("Expected IndentationType to be %v, got %v", tc.expected.IndentationType, result.IndentationType)
		}

		if *result.SpaceCount != *tc.expected.SpaceCount {
			t.Errorf("Expected SpaceCount to be %v, got %v", *tc.expected.SpaceCount, *result.SpaceCount)
		}
	}
}
