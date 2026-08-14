package doctor

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

type managedContract struct {
	ContractVersion string            `json:"contractVersion"`
	Files           map[string]string `json:"files"`
}

func auditContract(repository Repository) (ContractReport, error) {
	content, err := repository.ReadText(".aruo/managed.json")
	if errors.Is(err, fs.ErrNotExist) {
		return ContractReport{}, nil
	}
	if err != nil {
		return ContractReport{}, err
	}
	report := ContractReport{Present: true, Valid: true}
	var managed managedContract
	if err := json.Unmarshal([]byte(content), &managed); err != nil {
		addContractFinding(&report, ".aruo/managed.json is not valid JSON", "Run a future safe Aruo repair/update workflow or restore the manifest from version control.")
		report.Valid = false
		return report, nil //nolint:nilerr // Invalid user JSON is a report finding, not an operational failure.
	}
	report.Version = managed.ContractVersion
	if managed.ContractVersion == "" || len(managed.Files) == 0 {
		addContractFinding(&report, "managed contract metadata is incomplete", "Reinitialize the contract in a clean repository or restore valid managed metadata.")
		report.Valid = false
	}
	paths := make([]string, 0, len(managed.Files))
	for name := range managed.Files {
		paths = append(paths, name)
	}
	sort.Strings(paths)
	for _, name := range paths {
		expected := managed.Files[name]
		file := ContractFile{Path: name, Status: "VERIFIED"}
		switch {
		case name == "aruo.yaml":
			file.Status = "INVALID_OWNERSHIP"
			addContractFinding(&report, "aruo.yaml must remain application-owned", "Remove aruo.yaml from .aruo/managed.json; Aruo updates must not overwrite application intent.")
		case !safeEvidencePath(name):
			file.Status = "INVALID_PATH"
			addContractFinding(&report, fmt.Sprintf("managed path %q is unsafe", name), "Keep managed paths repository-relative without traversal.")
		case !strings.HasPrefix(expected, "sha256:"):
			file.Status = "INVALID_DIGEST"
			addContractFinding(&report, fmt.Sprintf("managed digest for %q is invalid", name), "Restore a sha256 digest from a trusted Aruo installation.")
		default:
			managedContent, readErr := repository.ReadText(name)
			switch {
			case errors.Is(readErr, fs.ErrNotExist):
				file.Status = "MISSING"
				addContractFinding(&report, fmt.Sprintf("managed contract file %q is missing", name), "Restore the managed file before relying on the AI contract.")
			case readErr != nil:
				return report, readErr
			default:
				digest := sha256.Sum256([]byte(managedContent))
				actual := "sha256:" + hex.EncodeToString(digest[:])
				if actual != expected {
					file.Status = "MODIFIED"
					addContractFinding(&report, fmt.Sprintf("managed contract file %q was modified", name), "Restore it from version control or use a future Aruo update workflow; put application-specific policy outside managed files.")
				}
			}
		}
		report.Files = append(report.Files, file)
	}
	if report.BlockingFindings > 0 {
		report.Valid = false
	}
	return report, nil
}

func addContractFinding(report *ContractReport, message, action string) {
	report.Findings = append(report.Findings, ContractFinding{Severity: "error", Message: message, Action: action})
	report.BlockingFindings++
}
