package doctor

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"
	"testing/fstest"

	"github.com/aruodore/aruo-cli/internal/contractmeta"
)

func TestAuditContractVerifiesCompleteManagedInventory(t *testing.T) {
	t.Parallel()
	repository, err := NewRepository(validContractRepository(t, "2"))
	if err != nil {
		t.Fatal(err)
	}
	report, err := auditContract(repository)
	if err != nil {
		t.Fatal(err)
	}
	requiredFiles, _ := contractmeta.RequiredFiles("2")
	if !report.Present || !report.Valid || report.BlockingFindings != 0 || len(report.Files) != len(requiredFiles) {
		t.Fatalf("report = %#v, want complete verified contract", report)
	}
	for _, file := range report.Files {
		if file.Status != "VERIFIED" {
			t.Errorf("file %s status = %s, want VERIFIED", file.Path, file.Status)
		}
	}
}

func TestAuditContractBlocksModifiedManagedFile(t *testing.T) {
	t.Parallel()
	files := validContractRepository(t, "2")
	files["AGENTS.md"] = &fstest.MapFile{Data: []byte("modified")}
	repository, err := NewRepository(files)
	if err != nil {
		t.Fatal(err)
	}
	report, err := auditContract(repository)
	if err != nil {
		t.Fatal(err)
	}
	if report.Valid || report.BlockingFindings != 1 || statusFor(report, "AGENTS.md") != "MODIFIED" {
		t.Fatalf("report = %#v, want one blocking modified file", report)
	}
}

func TestAuditContractRejectsApplicationOwnedIntent(t *testing.T) {
	t.Parallel()
	files := validContractRepository(t, "2")
	manifest := decodeManagedManifest(t, files)
	manifest.Files["aruo.yaml"] = "sha256:0000"
	files[".aruo/managed.json"] = &fstest.MapFile{Data: marshalJSON(t, manifest)}
	repository, err := NewRepository(files)
	if err != nil {
		t.Fatal(err)
	}
	report, err := auditContract(repository)
	if err != nil {
		t.Fatal(err)
	}
	if report.BlockingFindings != 1 || statusFor(report, "aruo.yaml") != "INVALID_OWNERSHIP" {
		t.Fatalf("report = %#v, want ownership violation", report)
	}
}

func TestAuditContractRejectsIncompleteInventory(t *testing.T) {
	t.Parallel()
	files := validContractRepository(t, "2")
	manifest := decodeManagedManifest(t, files)
	delete(manifest.Files, ".aruo/rules/security.md")
	files[".aruo/managed.json"] = &fstest.MapFile{Data: marshalJSON(t, manifest)}
	repository, err := NewRepository(files)
	if err != nil {
		t.Fatal(err)
	}
	report, err := auditContract(repository)
	if err != nil {
		t.Fatal(err)
	}
	if report.Valid || report.BlockingFindings != 1 {
		t.Fatalf("report = %#v, want missing-inventory finding", report)
	}
}

func TestAuditContractRejectsUnsupportedVersion(t *testing.T) {
	t.Parallel()
	files := validContractRepository(t, "999")
	repository, err := NewRepository(files)
	if err != nil {
		t.Fatal(err)
	}
	report, err := auditContract(repository)
	if err != nil {
		t.Fatal(err)
	}
	if report.Valid || report.BlockingFindings != 1 || report.Version != "999" {
		t.Fatalf("report = %#v, want unsupported-version finding", report)
	}
}

func validContractRepository(t *testing.T, version string) fstest.MapFS {
	t.Helper()
	files := fstest.MapFS{}
	manifest := managedContract{ContractVersion: version, Files: map[string]string{}}
	requiredFiles, supported := contractmeta.RequiredFiles(version)
	if !supported {
		requiredFiles, _ = contractmeta.RequiredFiles(contractmeta.CurrentVersion)
	}
	for _, name := range requiredFiles {
		content := []byte("managed content for " + name)
		digest := sha256.Sum256(content)
		manifest.Files[name] = "sha256:" + hex.EncodeToString(digest[:])
		files[name] = &fstest.MapFile{Data: content}
	}
	files[".aruo/managed.json"] = &fstest.MapFile{Data: marshalJSON(t, manifest)}
	return files
}

func decodeManagedManifest(t *testing.T, files fstest.MapFS) managedContract {
	t.Helper()
	var manifest managedContract
	if err := json.Unmarshal(files[".aruo/managed.json"].Data, &manifest); err != nil {
		t.Fatal(err)
	}
	return manifest
}

func marshalJSON(t *testing.T, value any) []byte {
	t.Helper()
	content, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return content
}

func statusFor(report ContractReport, path string) string {
	for _, file := range report.Files {
		if file.Path == path {
			return file.Status
		}
	}
	return ""
}
