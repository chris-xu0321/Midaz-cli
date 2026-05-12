package position

import (
	"context"
	"testing"

	"github.com/SparkssL/Midaz-cli/internal/auth"
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
)

// TestResolveDeskSlug_FromCreds verifies the cached-slug fast path —
// when DeskSlug is already in auth.json, we don't touch the API. We
// stub f.Client to nil so any HTTP attempt would panic, asserting the
// fast path actually returns before falling through.
func TestResolveDeskSlug_FromCreds(t *testing.T) {
	cases := []struct {
		name  string
		creds *auth.Creds
		want  string
	}{
		{
			name:  "slug_wins",
			creds: &auth.Creds{Credentials: &auth.Credentials{DeskSlug: "my-slug", DeskID: "uuid-here"}},
			want:  "my-slug",
		},
		{
			name:  "id_fallback",
			creds: &auth.Creds{Credentials: &auth.Credentials{DeskID: "uuid-here"}},
			want:  "uuid-here",
		},
	}
	for _, c := range cases {
		got, err := resolveDeskSlug(context.Background(), &cmdutil.Factory{}, c.creds)
		if err != nil {
			t.Fatalf("%s: unexpected err: %v", c.name, err)
		}
		if got != c.want {
			t.Errorf("%s: got %q want %q", c.name, got, c.want)
		}
	}
}
