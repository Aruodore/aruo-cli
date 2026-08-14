package contractmeta

import "testing"

func TestRequiredFilesReturnsDefensiveCopy(t *testing.T) {
	t.Parallel()
	files, supported := RequiredFiles(CurrentVersion)
	if !supported || len(files) == 0 {
		t.Fatal("current contract inventory is unavailable")
	}
	files[0] = "changed"
	again, _ := RequiredFiles(CurrentVersion)
	if again[0] == "changed" {
		t.Fatal("RequiredFiles exposed mutable package state")
	}
}

func TestRequiredFilesRejectsUnknownVersion(t *testing.T) {
	t.Parallel()
	if files, supported := RequiredFiles("999"); supported || files != nil {
		t.Fatalf("RequiredFiles returned %#v, %v for an unsupported version", files, supported)
	}
}
