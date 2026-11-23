package pages

import "testing"

// the viewer's own rows read "You ...", so each kind needs a different sentence for them
func TestActivityKeysHaveBothVariants(t *testing.T) {
	for kind, k := range activityKeys {
		if k.them == "" || k.you == "" || k.you == k.them {
			t.Errorf("%s: them %q, you %q", kind, k.them, k.you)
		}
	}
}
