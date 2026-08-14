package initialize

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/aruodore/aruo-cli/internal/contractmeta"
)

// A managed contract text change must intentionally update this snapshot and
// CurrentVersion. Keeping snapshots keyed by version makes a versionless text
// change fail review visibly instead of silently changing every new install.
var managedContractSnapshots = map[string]string{
	"2": "3f7c48b11cbd719c8dbde52009f34a8435c34e9717b22f6dbef9fdf1d6581b0b",
}

func TestManagedContractSnapshotMatchesVersion(t *testing.T) {
	t.Parallel()
	required, supported := contractmeta.RequiredFiles(contractmeta.CurrentVersion)
	if !supported {
		t.Fatalf("current contract version %q has no inventory", contractmeta.CurrentVersion)
	}
	hash := sha256.New()
	for _, destination := range required {
		if destination == ".aruo/stack.yaml" {
			continue
		}
		embeddedPath := strings.TrimPrefix(destination, ".aruo/")
		content, err := contractFiles.ReadFile("contract/" + embeddedPath)
		if err != nil {
			t.Fatalf("read %s: %v", destination, err)
		}
		hash.Write([]byte(destination))
		hash.Write([]byte{0})
		hash.Write(content)
		hash.Write([]byte{0})
	}
	got := hex.EncodeToString(hash.Sum(nil))
	want, exists := managedContractSnapshots[contractmeta.CurrentVersion]
	if !exists {
		t.Fatalf("contract version %q has no reviewed text snapshot; got %s", contractmeta.CurrentVersion, got)
	}
	if got != want {
		t.Fatalf("managed contract text changed without a reviewed version snapshot: got %s, want %s", got, want)
	}
}
