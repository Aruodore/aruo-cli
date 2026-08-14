package doctor

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"testing/fstest"
)

func TestAuditContractVerifiesManagedFiles(t *testing.T) {
	t.Parallel()
	content := "contract"
	digest := sha256.Sum256([]byte(content))
	manifest := fmt.Sprintf(`{"contractVersion":"1","files":{"AGENTS.md":"sha256:%s"}}`, hex.EncodeToString(digest[:]))
	repository, err := NewRepository(fstest.MapFS{
		".aruo/managed.json": &fstest.MapFile{Data: []byte(manifest)},
		"AGENTS.md":          &fstest.MapFile{Data: []byte(content)},
	})
	if err != nil {
		t.Fatal(err)
	}
	report, err := auditContract(repository)
	if err != nil {
		t.Fatal(err)
	}
	if !report.Present || !report.Valid || report.BlockingFindings != 0 || len(report.Files) != 1 || report.Files[0].Status != "VERIFIED" {
		t.Fatalf("report = %#v, want verified contract", report)
	}
}

func TestAuditContractBlocksModifiedManagedFile(t *testing.T) {
	t.Parallel()
	repository, err := NewRepository(fstest.MapFS{
		".aruo/managed.json": &fstest.MapFile{Data: []byte(`{"contractVersion":"1","files":{"AGENTS.md":"sha256:0000"}}`)},
		"AGENTS.md":          &fstest.MapFile{Data: []byte("modified")},
	})
	if err != nil {
		t.Fatal(err)
	}
	report, err := auditContract(repository)
	if err != nil {
		t.Fatal(err)
	}
	if report.Valid || report.BlockingFindings != 1 || report.Files[0].Status != "MODIFIED" {
		t.Fatalf("report = %#v, want blocking modified contract", report)
	}
}

func TestAuditContractRejectsApplicationOwnedIntent(t *testing.T) {
	t.Parallel()
	repository, err := NewRepository(fstest.MapFS{
		".aruo/managed.json": &fstest.MapFile{Data: []byte(`{"contractVersion":"1","files":{"aruo.yaml":"sha256:0000"}}`)},
	})
	if err != nil {
		t.Fatal(err)
	}
	report, err := auditContract(repository)
	if err != nil {
		t.Fatal(err)
	}
	if report.BlockingFindings != 1 || report.Files[0].Status != "INVALID_OWNERSHIP" {
		t.Fatalf("report = %#v, want ownership violation", report)
	}
}
