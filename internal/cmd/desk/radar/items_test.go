package radar

import (
	"strings"
	"testing"
)

func TestRenderers(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"thesis with title", renderThesisItem("abc123", "AI boom play"), "thesis:abc123 AI boom play"},
		{"thesis without title", renderThesisItem("abc123", ""), "thesis:abc123"},
		{"thesis trims label", renderThesisItem("abc123", "   spaced   "), "thesis:abc123 spaced"},
		{"driver", renderDriverItem("d1", "Fed Policy"), "driver:d1 Fed Policy"},
		{"url", renderURLItem("https://example.com", "Example"), "url:https://example.com Example"},
		{"asset uppercases", renderAssetItem("aapl"), "asset:AAPL"},
		{"asset trims", renderAssetItem("  btc  "), "asset:BTC"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s: got %q want %q", c.name, c.got, c.want)
		}
	}
}

func TestClipTo160(t *testing.T) {
	short := "hello"
	if got, trunc := clipTo160(short); got != short || trunc {
		t.Errorf("short: got %q trunc=%v", got, trunc)
	}

	exactly := strings.Repeat("a", maxItemLength)
	if got, trunc := clipTo160(exactly); got != exactly || trunc {
		t.Errorf("exact: got len=%d trunc=%v", len([]rune(got)), trunc)
	}

	over := strings.Repeat("b", maxItemLength+50)
	got, trunc := clipTo160(over)
	if !trunc {
		t.Errorf("over: expected truncation")
	}
	if runes := []rune(got); len(runes) != maxItemLength {
		t.Errorf("over: got rune length %d want %d", len(runes), maxItemLength)
	}
	if !strings.HasSuffix(got, "…") {
		t.Errorf("over: expected ellipsis suffix, got %q", got[len(got)-4:])
	}

	// multibyte: 200 CJK chars should clip to 160 runes (not 160 bytes).
	cjk := strings.Repeat("字", 200)
	clipped, trunc := clipTo160(cjk)
	if !trunc {
		t.Errorf("cjk: expected truncation")
	}
	if runes := []rune(clipped); len(runes) != maxItemLength {
		t.Errorf("cjk: got rune length %d want %d", len(runes), maxItemLength)
	}
}

func TestFindIndex(t *testing.T) {
	items := []string{"asset:AAPL", "thesis:abc 1", "driver:xyz"}
	if got := findIndex(items, "thesis:abc 1"); got != 1 {
		t.Errorf("existing: got %d want 1", got)
	}
	if got := findIndex(items, "not there"); got != -1 {
		t.Errorf("missing: got %d want -1", got)
	}
	// case-sensitive: "asset:aapl" != "asset:AAPL"
	if got := findIndex(items, "asset:aapl"); got != -1 {
		t.Errorf("case: got %d want -1 (findIndex is case-sensitive)", got)
	}
}

func TestPickOne(t *testing.T) {
	// zero
	if _, err := pickOne(map[string]string{"--a": "", "--b": ""}); err == nil {
		t.Error("zero: expected error")
	}
	// exactly one
	got, err := pickOne(map[string]string{"--a": "x", "--b": ""})
	if err != nil || got != "--a" {
		t.Errorf("one: got=%q err=%v", got, err)
	}
	// multiple
	if _, err := pickOne(map[string]string{"--a": "x", "--b": "y"}); err == nil {
		t.Error("multi: expected error")
	}
	// whitespace-only counts as empty
	got, err = pickOne(map[string]string{"--a": "   ", "--b": "y"})
	if err != nil || got != "--b" {
		t.Errorf("ws: got=%q err=%v", got, err)
	}
}
