package notification

import (
	"reflect"
	"testing"
)

func TestNormalizeTargets(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want []string
	}{
		{"empty", nil, []string{}},
		{"trims whitespace", []string{" abc ", "def"}, []string{"abc", "def"}},
		{"drops blanks", []string{"abc", "  ", ""}, []string{"abc"}},
		{"dedupes", []string{"abc", "abc", " abc "}, []string{"abc"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeTargets(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("normalizeTargets(%v) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}
