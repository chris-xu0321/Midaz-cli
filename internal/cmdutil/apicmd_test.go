package cmdutil

import (
	"encoding/json"
	"testing"
)

func TestInjectItemViewURL(t *testing.T) {
	pattern := func(id string) string { return "https://ex.test/market?driver=" + id }

	t.Run("injects into items missing view_url", func(t *testing.T) {
		body := []byte(`[{"id":"a","name":"first"},{"id":"b","name":"second"}]`)
		data, meta, err := InjectItemViewURL(pattern)(body)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if meta["count"] != 2 {
			t.Fatalf("expected count=2, got %v", meta["count"])
		}
		items, ok := data.([]interface{})
		if !ok || len(items) != 2 {
			t.Fatalf("expected 2 items, got %v", data)
		}
		first := items[0].(map[string]interface{})
		if first["view_url"] != "https://ex.test/market?driver=a" {
			t.Errorf("first.view_url = %v", first["view_url"])
		}
		second := items[1].(map[string]interface{})
		if second["view_url"] != "https://ex.test/market?driver=b" {
			t.Errorf("second.view_url = %v", second["view_url"])
		}
	})

	t.Run("preserves existing view_url", func(t *testing.T) {
		body := []byte(`[{"id":"a","view_url":"https://existing.test/x"}]`)
		data, _, err := InjectItemViewURL(pattern)(body)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		item := data.([]interface{})[0].(map[string]interface{})
		if item["view_url"] != "https://existing.test/x" {
			t.Errorf("expected existing view_url kept, got %v", item["view_url"])
		}
	})

	t.Run("skips items without id", func(t *testing.T) {
		body := []byte(`[{"name":"no-id"}]`)
		data, _, err := InjectItemViewURL(pattern)(body)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		item := data.([]interface{})[0].(map[string]interface{})
		if _, has := item["view_url"]; has {
			t.Error("view_url should not be injected when id is missing")
		}
	})

	t.Run("skips when pattern returns empty", func(t *testing.T) {
		body := []byte(`[{"id":"a"}]`)
		data, _, err := InjectItemViewURL(func(string) string { return "" })(body)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		item := data.([]interface{})[0].(map[string]interface{})
		if _, has := item["view_url"]; has {
			t.Error("view_url should not be injected when pattern returns empty")
		}
	})

	t.Run("empty array returns empty data with count 0", func(t *testing.T) {
		body := []byte(`[]`)
		data, meta, err := InjectItemViewURL(pattern)(body)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if meta["count"] != 0 {
			t.Fatalf("expected count=0, got %v", meta["count"])
		}
		items, ok := data.([]interface{})
		if !ok || len(items) != 0 {
			t.Fatalf("expected empty slice, got %v", data)
		}
	})

	t.Run("non-array body returns error", func(t *testing.T) {
		_, _, err := InjectItemViewURL(pattern)([]byte(`{"id":"a"}`))
		if err == nil {
			t.Fatal("expected error on non-array body")
		}
	})

	t.Run("result remains JSON-marshallable", func(t *testing.T) {
		body := []byte(`[{"id":"a","nested":{"k":1},"arr":[1,2]}]`)
		data, _, err := InjectItemViewURL(pattern)(body)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, err := json.Marshal(data); err != nil {
			t.Fatalf("result should marshal: %v", err)
		}
	})
}
