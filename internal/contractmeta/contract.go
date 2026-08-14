// Package contractmeta defines the versioned ownership boundary shared by
// contract installation and auditing.
package contractmeta

const CurrentVersion = "2"

var requiredFilesByVersion = map[string][]string{
	"1": managedContractFiles(),
	"2": managedContractFiles(),
}

// RequiredFiles returns a copy of the exact managed file inventory for a
// supported contract version. The manifest itself is managed but omitted
// because it cannot contain its own digest.
func RequiredFiles(version string) ([]string, bool) {
	files, supported := requiredFilesByVersion[version]
	if !supported {
		return nil, false
	}
	return append([]string(nil), files...), true
}

func managedContractFiles() []string {
	return []string{
		".aruo/contract.yaml",
		".aruo/rules/api.md",
		".aruo/rules/architecture.md",
		".aruo/rules/data.md",
		".aruo/rules/delivery.md",
		".aruo/rules/observability.md",
		".aruo/rules/security.md",
		".aruo/rules/testing.md",
		".aruo/stack.yaml",
		"AGENTS.md",
	}
}
