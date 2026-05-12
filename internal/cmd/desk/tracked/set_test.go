package tracked

import (
	"reflect"
	"testing"
)

func TestParseAssetList(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{"empty", "", nil},
		{"single", "NVDA", []string{"NVDA"}},
		{"comma_sep_with_spaces", "NVDA, GLD,  US10Y", []string{"NVDA", "GLD", "US10Y"}},
		{"lowercase_uppercased", "nvda,gld", []string{"NVDA", "GLD"}},
		{"dedupe_preserves_first", "NVDA,GLD,NVDA,gld", []string{"NVDA", "GLD"}},
		{"drops_empty_tokens", "NVDA,,GLD,   ", []string{"NVDA", "GLD"}},
		{"newline_separator_unsupported_here", "NVDA\nGLD", []string{"NVDA\nGLD"}}, // intentional: set.go normalizes newlines before calling
	}
	for _, c := range cases {
		got := parseAssetList(c.in)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: got %v want %v", c.name, got, c.want)
		}
	}
}

func TestMergeUnique(t *testing.T) {
	cases := []struct {
		name string
		a, b []string
		want []string
	}{
		{"both_empty", nil, nil, []string{}},
		{"a_only", []string{"NVDA"}, nil, []string{"NVDA"}},
		{"b_only", nil, []string{"GLD"}, []string{"GLD"}},
		{"disjoint", []string{"NVDA"}, []string{"GLD"}, []string{"NVDA", "GLD"}},
		{"overlap_keeps_a_order", []string{"NVDA", "GLD"}, []string{"GLD", "TLT"}, []string{"NVDA", "GLD", "TLT"}},
		{"normalizes_case", []string{"nvda"}, []string{"NVDA"}, []string{"NVDA"}},
		{"trims_whitespace", []string{" nvda "}, []string{"GLD"}, []string{"NVDA", "GLD"}},
	}
	for _, c := range cases {
		got := mergeUnique(c.a, c.b)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: got %v want %v", c.name, got, c.want)
		}
	}
}

func TestRemoveFrom(t *testing.T) {
	cases := []struct {
		name       string
		from, drop []string
		want       []string
	}{
		{"empty_from", nil, []string{"NVDA"}, []string{}},
		{"empty_drop", []string{"NVDA", "GLD"}, nil, []string{"NVDA", "GLD"}},
		{"drops_one", []string{"NVDA", "GLD", "TLT"}, []string{"GLD"}, []string{"NVDA", "TLT"}},
		{"case_insensitive", []string{"NVDA", "GLD"}, []string{"gld"}, []string{"NVDA"}},
		{"drop_all", []string{"NVDA", "GLD"}, []string{"NVDA", "GLD"}, []string{}},
	}
	for _, c := range cases {
		got := removeFrom(c.from, c.drop)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: got %v want %v", c.name, got, c.want)
		}
	}
}
