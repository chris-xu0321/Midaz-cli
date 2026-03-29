package topic

import (
	"encoding/json"
	"testing"
)

func TestNormalizeTopic(t *testing.T) {
	apiResponse := `{
		"id": "top1", "name": "AI Infrastructure", "bias": "bullish",
		"standing_thesis": "test thesis",
		"view_url": "http://localhost:3000/topics/top1",
		"threads": [{"id": "t1"}, {"id": "t2"}, {"id": "t3"}],
		"snapshot": {}
	}`

	data, meta, err := normalizeTopic([]byte(apiResponse))
	if err != nil {
		t.Fatal(err)
	}

	// Check meta
	if meta["view_url"] != "http://localhost:3000/topics/top1" {
		t.Errorf("expected view_url in meta, got %v", meta["view_url"])
	}
	if meta["thread_count"] != 3 {
		t.Errorf("expected thread_count=3, got %v", meta["thread_count"])
	}

	// Check data: view_url removed
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		t.Fatal("data should be a map")
	}
	if _, exists := dataMap["view_url"]; exists {
		t.Error("view_url should be removed from data")
	}

	// threads should still be in data
	if _, exists := dataMap["threads"]; !exists {
		t.Error("threads should remain in data")
	}

	// Verify data is valid JSON
	if _, err := json.Marshal(data); err != nil {
		t.Fatalf("data should be marshallable: %v", err)
	}
}
