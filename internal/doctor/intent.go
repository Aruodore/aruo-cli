package doctor

import (
	"errors"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type intentManifest struct {
	APIVersion string `yaml:"apiVersion"`
	Template   struct {
		ID      string `yaml:"id"`
		Profile string `yaml:"profile"`
	} `yaml:"template"`
	Intent struct {
		Capabilities map[string]struct {
			Status   string `yaml:"status"`
			Evidence string `yaml:"evidence"`
			Reason   string `yaml:"reason"`
		} `yaml:"capabilities"`
	} `yaml:"intent"`
}

func auditIntent(repository Repository) (IntentReport, error) {
	content, err := repository.ReadText("aruo.yaml")
	if errors.Is(err, fs.ErrNotExist) {
		return IntentReport{
			Findings: []IntentFinding{{Severity: "warning", Message: "No aruo.yaml intent manifest found", Action: "Add a versioned intent manifest when this repository adopts the Aruo production contract."}},
		}, nil
	}
	if err != nil {
		return IntentReport{}, err
	}
	report := IntentReport{Present: true}
	var manifest intentManifest
	decoder := yaml.NewDecoder(strings.NewReader(content))
	if err := decoder.Decode(&manifest); err != nil {
		report.Findings = []IntentFinding{{Severity: "error", Message: "aruo.yaml is not valid YAML", Action: "Fix the manifest syntax before relying on its capability claims."}}
		report.BlockingFindings = 1
		return report, nil //nolint:nilerr // Invalid user YAML is a report finding, not an operational failure.
	}
	report.Valid = true
	report.APIVersion = manifest.APIVersion
	report.TemplateID = manifest.Template.ID
	report.Profile = manifest.Template.Profile
	if manifest.APIVersion != "aruo.dev/v1alpha1" {
		report.Valid = false
		addIntentFinding(&report, "", "error", fmt.Sprintf("unsupported apiVersion %q", manifest.APIVersion), "Use aruo.dev/v1alpha1 or migrate the manifest with a compatible Aruo release.", true)
	}

	names := make([]string, 0, len(manifest.Intent.Capabilities))
	for name := range manifest.Intent.Capabilities {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		declaration := manifest.Intent.Capabilities[name]
		capability := IntentCapability{Name: name, Status: CapabilityStatus(declaration.Status), Evidence: declaration.Evidence, Reason: declaration.Reason, EvidenceStatus: EvidenceNotApplicable}
		if !validCapabilityStatus(capability.Status) {
			report.Valid = false
			addIntentFinding(&report, name, "error", fmt.Sprintf("unknown capability status %q", declaration.Status), "Use SOLVED, REQUIRED, OPTIONAL, DEFERRED, or UNKNOWN.", true)
			report.Capabilities = append(report.Capabilities, capability)
			continue
		}
		switch capability.Status {
		case CapabilitySolved:
			switch {
			case capability.Evidence == "":
				capability.EvidenceStatus = EvidenceMissing
				addIntentFinding(&report, name, "error", "SOLVED capability has no evidence", "Name a repository-relative file or a precise declaration that supports this claim.", true)
			case looksLikeRepositoryPath(capability.Evidence):
				if !safeEvidencePath(capability.Evidence) {
					capability.EvidenceStatus = EvidenceMissing
					addIntentFinding(&report, name, "error", "capability evidence is not a safe repository-relative path", "Use a clean repository-relative path without traversal.", true)
				} else {
					exists, evidenceErr := repository.Exists(capability.Evidence)
					if evidenceErr != nil {
						return report, evidenceErr
					}
					if exists {
						capability.EvidenceStatus = EvidenceVerified
					} else {
						capability.EvidenceStatus = EvidenceMissing
						addIntentFinding(&report, name, "error", fmt.Sprintf("declared evidence %q does not exist", capability.Evidence), "Restore the evidence or change the capability status honestly.", true)
					}
				}
			default:
				capability.EvidenceStatus = EvidenceDeclared
			}
		case CapabilityRequired:
			if capability.Reason == "" {
				addIntentFinding(&report, name, "error", "REQUIRED capability has no reason", "Explain why this responsibility remains unresolved.", true)
			} else {
				addIntentFinding(&report, name, "required", capability.Reason, "Implement and verify this capability before the relevant production use, then mark it SOLVED with evidence.", true)
			}
		case CapabilityOptional, CapabilityDeferred, CapabilityUnknown:
			if capability.Reason == "" {
				addIntentFinding(&report, name, "warning", fmt.Sprintf("%s capability has no reason", capability.Status), "Document why this status is appropriate for the application.", false)
			}
		}
		report.Capabilities = append(report.Capabilities, capability)
	}
	return report, nil
}

func addIntentFinding(report *IntentReport, capability, severity, message, action string, blocking bool) {
	report.Findings = append(report.Findings, IntentFinding{Capability: capability, Severity: severity, Message: message, Action: action})
	if blocking {
		report.BlockingFindings++
	}
}

func validCapabilityStatus(status CapabilityStatus) bool {
	switch status {
	case CapabilitySolved, CapabilityRequired, CapabilityOptional, CapabilityDeferred, CapabilityUnknown:
		return true
	default:
		return false
	}
}

func looksLikeRepositoryPath(value string) bool {
	return strings.Contains(value, "/") || strings.HasPrefix(value, ".") || path.Ext(value) != ""
}

func safeEvidencePath(value string) bool {
	clean := path.Clean(value)
	return value != "" && clean == value && clean != "." && !path.IsAbs(clean) && clean != ".." && !strings.HasPrefix(clean, "../")
}
