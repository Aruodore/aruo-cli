package doctor

import (
	"testing"
	"testing/fstest"
)

func TestAuditIntent(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		manifest    string
		files       map[string]string
		present     bool
		valid       bool
		blocking    int
		evidence    EvidenceStatus
		capability  string
		findingType string
	}{
		{name: "missing manifest is visible but not blocking", valid: false},
		{
			name:       "repository path is verified",
			manifest:   manifestWith("quality-gate: { status: SOLVED, evidence: scripts/check.sh }"),
			files:      map[string]string{"scripts/check.sh": "#!/bin/sh\n"},
			present:    true,
			valid:      true,
			evidence:   EvidenceVerified,
			capability: "quality-gate",
		},
		{
			name:       "semantic evidence remains declared",
			manifest:   manifestWith("structured-logging: { status: SOLVED, evidence: json-request-logs }"),
			present:    true,
			valid:      true,
			evidence:   EvidenceDeclared,
			capability: "structured-logging",
		},
		{
			name:        "missing path disproves solved claim",
			manifest:    manifestWith("quality-gate: { status: SOLVED, evidence: scripts/check.sh }"),
			present:     true,
			valid:       true,
			blocking:    1,
			evidence:    EvidenceMissing,
			capability:  "quality-gate",
			findingType: "error",
		},
		{
			name:        "required responsibility blocks",
			manifest:    manifestWith("authentication: { status: REQUIRED, reason: identity-not-selected }"),
			present:     true,
			valid:       true,
			blocking:    1,
			evidence:    EvidenceNotApplicable,
			capability:  "authentication",
			findingType: "required",
		},
		{
			name:       "optional responsibility with reason does not block",
			manifest:   manifestWith("email: { status: OPTIONAL, reason: product-dependent }"),
			present:    true,
			valid:      true,
			evidence:   EvidenceNotApplicable,
			capability: "email",
		},
		{
			name:        "unknown status invalidates manifest",
			manifest:    manifestWith("auth: { status: COMPLETE, evidence: auth.ts }"),
			present:     true,
			blocking:    1,
			evidence:    EvidenceNotApplicable,
			capability:  "auth",
			findingType: "error",
		},
		{
			name:        "unsafe traversal is rejected",
			manifest:    manifestWith("auth: { status: SOLVED, evidence: ../auth.ts }"),
			present:     true,
			valid:       true,
			blocking:    1,
			evidence:    EvidenceMissing,
			capability:  "auth",
			findingType: "error",
		},
		{
			name:        "malformed yaml blocks",
			manifest:    "apiVersion: [\n",
			present:     true,
			blocking:    1,
			findingType: "error",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			files := fstest.MapFS{}
			if test.manifest != "" {
				files["aruo.yaml"] = &fstest.MapFile{Data: []byte(test.manifest)}
			}
			for name, content := range test.files {
				files[name] = &fstest.MapFile{Data: []byte(content)}
			}
			repository, err := NewRepository(files)
			if err != nil {
				t.Fatal(err)
			}
			report, err := auditIntent(repository)
			if err != nil {
				t.Fatal(err)
			}
			if report.Present != test.present || report.Valid != test.valid || report.BlockingFindings != test.blocking {
				t.Fatalf("report state = present %t, valid %t, blocking %d; want %t, %t, %d", report.Present, report.Valid, report.BlockingFindings, test.present, test.valid, test.blocking)
			}
			if test.capability != "" {
				if len(report.Capabilities) != 1 || report.Capabilities[0].Name != test.capability || report.Capabilities[0].EvidenceStatus != test.evidence {
					t.Fatalf("capabilities = %#v; want %s with %s evidence", report.Capabilities, test.capability, test.evidence)
				}
			}
			if test.findingType != "" {
				if len(report.Findings) == 0 || report.Findings[0].Severity != test.findingType {
					t.Fatalf("findings = %#v; want first severity %s", report.Findings, test.findingType)
				}
			}
		})
	}
}

func TestAuditIntentSortsCapabilities(t *testing.T) {
	t.Parallel()
	repository, err := NewRepository(fstest.MapFS{
		"aruo.yaml": &fstest.MapFile{Data: []byte(manifestWith("zeta: { status: OPTIONAL, reason: later }\n    alpha: { status: DEFERRED, reason: now }"))},
	})
	if err != nil {
		t.Fatal(err)
	}
	report, err := auditIntent(repository)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Capabilities) != 2 || report.Capabilities[0].Name != "alpha" || report.Capabilities[1].Name != "zeta" {
		t.Fatalf("capabilities = %#v; want deterministic alphabetical order", report.Capabilities)
	}
}

func manifestWith(capabilities string) string {
	return "apiVersion: aruo.dev/v1alpha1\n" +
		"template: { id: test, profile: lean }\n" +
		"intent:\n  capabilities:\n    " + capabilities + "\n"
}
