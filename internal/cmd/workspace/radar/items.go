package radar

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/SparkssL/Midaz-cli/internal/client"
	"github.com/SparkssL/Midaz-cli/internal/output"
)

const (
	maxRadarItems = 12  // matches Seer apps/api/src/services/radar.ts MAX_ITEMS
	maxItemLength = 160 // matches Seer MAX_ITEM_LEN
)

// fetchCurrentItems returns the current radar items from GET /api/ws/settings.
func fetchCurrentItems(ctx context.Context, c *client.Client) ([]string, error) {
	resp, err := c.Get(ctx, "/api/ws/settings", nil)
	if err != nil {
		return nil, err
	}
	var parsed struct {
		RadarItems []string `json:"radar_items"`
	}
	if err := json.Unmarshal(resp.Body, &parsed); err != nil {
		return nil, output.Errorf(output.ExitInternal, "internal",
			"failed to parse /api/ws/settings response: %s", err)
	}
	return parsed.RadarItems, nil
}

// pushItems writes the full items list back via PATCH /api/ws/radar.
// Returns the server's parsed response for echoing in the success envelope.
func pushItems(ctx context.Context, c *client.Client, items []string) (map[string]any, error) {
	resp, err := c.Patch(ctx, "/api/ws/radar", map[string]any{"radar_items": items})
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, &out); err != nil {
			return nil, output.Errorf(output.ExitInternal, "internal",
				"failed to parse /api/ws/radar response: %s", err)
		}
	}
	return out, nil
}

// resolveTitle fetches an entity by id and extracts a human-readable label.
// path is something like "/api/theses/" or "/api/topics/". titleKey is "title"
// for theses, "name" for topics.
func resolveTitle(ctx context.Context, c *client.Client, path, id, titleKey string) (string, error) {
	resp, err := c.Get(ctx, path+url.PathEscape(id), nil)
	if err != nil {
		return "", err
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(resp.Body, &m); err != nil {
		return "", output.Errorf(output.ExitInternal, "internal",
			"failed to parse %s response: %s", path, err)
	}
	raw, ok := m[titleKey]
	if !ok {
		return "", nil
	}
	var s string
	_ = json.Unmarshal(raw, &s)
	return strings.TrimSpace(s), nil
}

// renderThesisItem formats a thesis reference as a radar line.
func renderThesisItem(id, title string) string {
	return renderRef("thesis", id, title)
}

// renderTopicItem formats a topic reference as a radar line.
func renderTopicItem(id, name string) string {
	return renderRef("topic", id, name)
}

// renderURLItem formats an external URL + title as a radar line.
func renderURLItem(u, title string) string {
	return renderRef("url", u, title)
}

// renderAssetItem formats a ticker as a radar line.
func renderAssetItem(ticker string) string {
	return "asset:" + strings.ToUpper(strings.TrimSpace(ticker))
}

func renderRef(kind, ref, label string) string {
	ref = strings.TrimSpace(ref)
	label = strings.TrimSpace(label)
	if label == "" {
		return fmt.Sprintf("%s:%s", kind, ref)
	}
	return fmt.Sprintf("%s:%s %s", kind, ref, label)
}

// clipTo160 truncates a line to maxItemLength runes, appending an ellipsis if
// truncated. Returns (clipped, true) when truncation happened.
func clipTo160(s string) (string, bool) {
	runes := []rune(s)
	if len(runes) <= maxItemLength {
		return s, false
	}
	return string(runes[:maxItemLength-1]) + "…", true
}

// findIndex returns the 0-based index of the first item equal to line, or -1.
func findIndex(items []string, line string) int {
	for i, it := range items {
		if it == line {
			return i
		}
	}
	return -1
}
