package thread

import (
	"encoding/json"
	"testing"
)

func TestNormalizeThread(t *testing.T) {
	apiResponse := `{
		"id": "t1", "title": "Test Thread", "thesis": "test", "bias": "bullish",
		"status": "active", "topic_id": "top1",
		"view_url": "http://localhost:3000/threads/t1",
		"topic_url": "http://localhost:3000/topics/top1",
		"claims": [{"id": "c1"}, {"id": "c2"}],
		"market_links": [{"market_id": "m1"}],
		"supporting_count": 2,
		"contradicting_count": 0,
		"has_market_link": true,
		"market_link_count": 1,
		"snapshot": {}
	}`

	data, meta, err := normalizeThread([]byte(apiResponse))
	if err != nil {
		t.Fatal(err)
	}

	// Check meta
	if meta["view_url"] != "http://localhost:3000/threads/t1" {
		t.Errorf("expected view_url in meta, got %v", meta["view_url"])
	}
	if meta["topic_url"] != "http://localhost:3000/topics/top1" {
		t.Errorf("expected topic_url in meta, got %v", meta["topic_url"])
	}
	if meta["claim_count"] != 2 {
		t.Errorf("expected claim_count=2, got %v", meta["claim_count"])
	}
	if meta["supporting_count"] != 2 {
		t.Errorf("expected supporting_count=2, got %v", meta["supporting_count"])
	}
	if meta["contradicting_count"] != 0 {
		t.Errorf("expected contradicting_count=0, got %v", meta["contradicting_count"])
	}
	if meta["market_link_count"] != 1 {
		t.Errorf("expected market_link_count=1, got %v", meta["market_link_count"])
	}

	// Check data: view_url, topic_url, has_market_link, market_link_count removed
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		t.Fatal("data should be a map")
	}
	if _, exists := dataMap["view_url"]; exists {
		t.Error("view_url should be removed from data")
	}
	if _, exists := dataMap["topic_url"]; exists {
		t.Error("topic_url should be removed from data")
	}
	if _, exists := dataMap["has_market_link"]; exists {
		t.Error("has_market_link should be removed from data")
	}
	if _, exists := dataMap["market_link_count"]; exists {
		t.Error("market_link_count should be removed from data")
	}

	// supporting_count and contradicting_count should STAY in data
	if _, exists := dataMap["supporting_count"]; !exists {
		t.Error("supporting_count should remain in data")
	}
	if _, exists := dataMap["contradicting_count"]; !exists {
		t.Error("contradicting_count should remain in data")
	}

	// Verify data is valid JSON
	if _, err := json.Marshal(data); err != nil {
		t.Fatalf("data should be marshallable: %v", err)
	}
}
