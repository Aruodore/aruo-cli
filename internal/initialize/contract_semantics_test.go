package initialize

import (
	"strings"
	"testing"

	"github.com/aruodore/aruo-cli/internal/contractmeta"
	"gopkg.in/yaml.v3"
)

func TestEmbeddedContractSchemaAndSafetyInvariants(t *testing.T) {
	t.Parallel()
	content, err := contractFiles.ReadFile("contract/contract.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var contract struct {
		APIVersion string            `yaml:"apiVersion"`
		Kind       string            `yaml:"kind"`
		Version    int               `yaml:"version"`
		Terms      map[string]string `yaml:"normativeTerms"`
		Ownership  struct {
			Managed          []string `yaml:"managed"`
			ApplicationOwned []string `yaml:"applicationOwned"`
		} `yaml:"ownership"`
	}
	if err := yaml.Unmarshal(content, &contract); err != nil {
		t.Fatalf("contract.yaml is invalid: %v", err)
	}
	if contract.APIVersion != "aruo.dev/v1alpha1" || contract.Kind != "EngineeringContract" || contract.Version != 2 || contractmeta.CurrentVersion != "2" {
		t.Fatalf("contract identity = %#v, code version = %q", contract, contractmeta.CurrentVersion)
	}
	if contract.Terms["MUST"] == "" || contract.Terms["SHOULD"] == "" || contract.Terms["CONDITIONAL"] == "" {
		t.Fatalf("normative terms = %#v, want MUST/SHOULD/CONDITIONAL", contract.Terms)
	}
	if !contains(contract.Ownership.Managed, "AGENTS.md") || !contains(contract.Ownership.ApplicationOwned, "AGENTS.local.md") || !contains(contract.Ownership.ApplicationOwned, "aruo.yaml") {
		t.Fatalf("ownership = %#v", contract.Ownership)
	}

	requireContractText(t, "contract/AGENTS.md", []string{
		"Do not add infrastructure solely to satisfy an inapplicable rule",
		"If authorization cannot be requested, do not perform the action",
		"preserve application-specific and unknown fields",
	})
	requireContractText(t, "contract/rules/security.md", []string{
		"Validation and safe use are separate controls",
		"constrained outbound destinations",
		"Never invent cryptographic protocols",
	})
	requireContractText(t, "contract/rules/testing.md", []string{
		"repository-provided checks relevant to the changed area",
		"Do not add tooling solely to satisfy this list",
	})
	requireContractText(t, "contract/rules/architecture.md", []string{
		"otherwise-unconstrained files",
		"Do not rename unrelated files",
		"MUST NOT introduce a distributed or deployment boundary",
	})

	architecture, err := contractFiles.ReadFile("contract/rules/architecture.md")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(architecture), "Keep the application a modular monolith") {
		t.Fatal("architecture contract restored the unconditional modular-monolith rule")
	}
}

func requireContractText(t *testing.T, path string, required []string) {
	t.Helper()
	content, err := contractFiles.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, phrase := range required {
		if !strings.Contains(string(content), phrase) {
			t.Errorf("%s is missing safety invariant %q", path, phrase)
		}
	}
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
