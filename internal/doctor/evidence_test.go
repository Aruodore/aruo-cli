package doctor

import (
	"testing"
	"testing/fstest"
)

func TestVerifySemanticEvidence(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		evidence   string
		files      fstest.MapFS
		recognized bool
		verified   bool
	}{
		{
			name:       "quality gate verifies required lifecycle checks",
			evidence:   "npm-run-check",
			files:      mapFiles(map[string]string{"package.json": `{"scripts":{"check":"npm run format:check && npm run lint && npm test && npm run build"}}`}),
			recognized: true,
			verified:   true,
		},
		{
			name:       "quality gate rejects incomplete script",
			evidence:   "npm-run-check",
			files:      mapFiles(map[string]string{"package.json": `{"scripts":{"check":"npm test"}}`}),
			recognized: true,
		},
		{
			name:       "vitest also requires a test file",
			evidence:   "vitest",
			files:      mapFiles(map[string]string{"package.json": `{"scripts":{"test":"vitest run"}}`}),
			recognized: true,
		},
		{
			name:       "vitest verifies runner and test",
			evidence:   "vitest",
			files:      mapFiles(map[string]string{"package.json": `{"scripts":{"test":"vitest run"}}`, "src/example.test.ts": "test('works', () => {})"}),
			recognized: true,
			verified:   true,
		},
		{
			name:       "strict typescript verifies compiler option",
			evidence:   "strict-typescript",
			files:      mapFiles(map[string]string{"tsconfig.json": "{\n  \"compilerOptions\": { \"strict\": true }\n}"}),
			recognized: true,
			verified:   true,
		},
		{
			name:       "dependency audit verifies workflow gate",
			evidence:   "ci-high-severity-gate",
			files:      mapFiles(map[string]string{".github/workflows/ci.yml": "steps:\n  - run: npm audit --audit-level=high\n"}),
			recognized: true,
			verified:   true,
		},
		{
			name:       "committed migration verifies sql",
			evidence:   "committed-drizzle-sql",
			files:      mapFiles(map[string]string{"server/db/migrations/0000.sql": "create table example;"}),
			recognized: true,
			verified:   true,
		},
		{
			name:       "health evidence requires both routes",
			evidence:   "live-and-ready",
			files:      mapFiles(map[string]string{"server/api/health/live.get.ts": "export default 1"}),
			recognized: true,
		},
		{
			name:       "health evidence verifies inline route pair",
			evidence:   "live-and-ready-routes",
			files:      mapFiles(map[string]string{"server/index.ts": `app.get("/api/health/live"); app.get("/api/health/ready")`}),
			recognized: true,
			verified:   true,
		},
		{
			name:     "runtime claim remains declared",
			evidence: "json-request-logs",
			files:    fstest.MapFS{},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			repository, err := NewRepository(test.files)
			if err != nil {
				t.Fatal(err)
			}
			verification, err := verifySemanticEvidence(repository, test.evidence)
			if err != nil {
				t.Fatal(err)
			}
			if verification.recognized != test.recognized || verification.verified != test.verified {
				t.Fatalf("verification = recognized %t, verified %t; want %t, %t", verification.recognized, verification.verified, test.recognized, test.verified)
			}
		})
	}
}

func TestAuditIntentBlocksFailedSemanticEvidence(t *testing.T) {
	t.Parallel()
	repository, err := NewRepository(mapFiles(map[string]string{
		"aruo.yaml":    manifestWith("quality-gate: { status: SOLVED, evidence: npm-run-check }"),
		"package.json": `{"scripts":{"check":"npm test"}}`,
	}))
	if err != nil {
		t.Fatal(err)
	}
	report, err := auditIntent(repository)
	if err != nil {
		t.Fatal(err)
	}
	if report.BlockingFindings != 1 || report.Capabilities[0].EvidenceStatus != EvidenceMissing {
		t.Fatalf("report = %#v, want one blocking missing-evidence finding", report)
	}
}

func mapFiles(files map[string]string) fstest.MapFS {
	result := make(fstest.MapFS, len(files))
	for name, content := range files {
		result[name] = &fstest.MapFile{Data: []byte(content)}
	}
	return result
}
