package initialize

import (
	"encoding/json"
	"errors"
	"io/fs"
	"sort"
)

type nodeManifest struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func detectStack(root fs.FS) (Stack, error) {
	stack := Stack{Ecosystem: "unknown"}
	if exists(root, "go.mod") {
		stack.Ecosystem = "go"
	}
	if exists(root, "pyproject.toml") || exists(root, "requirements.txt") {
		stack.Ecosystem = "python"
	}
	if exists(root, "package.json") {
		stack.Ecosystem = "node"
		manifest, err := readNodeManifest(root)
		if err != nil {
			return Stack{}, err
		}
		for _, framework := range []string{"next", "nuxt", "react", "vue", "vite"} {
			if _, ok := manifest.Dependencies[framework]; ok {
				stack.Frameworks = append(stack.Frameworks, framework)
				continue
			}
			if _, ok := manifest.DevDependencies[framework]; ok {
				stack.Frameworks = append(stack.Frameworks, framework)
			}
		}
		stack.PackageManager = detectPackageManager(root)
	}
	sort.Strings(stack.Frameworks)
	return stack, nil
}

func readNodeManifest(root fs.FS) (nodeManifest, error) {
	content, err := fs.ReadFile(root, "package.json")
	if err != nil {
		return nodeManifest{}, err
	}
	var manifest nodeManifest
	if err := json.Unmarshal(content, &manifest); err != nil {
		return nodeManifest{}, errors.New("package.json is not valid JSON")
	}
	return manifest, nil
}

func detectPackageManager(root fs.FS) string {
	for _, candidate := range []struct{ file, name string }{
		{"pnpm-lock.yaml", "pnpm"},
		{"yarn.lock", "yarn"},
		{"bun.lock", "bun"},
		{"bun.lockb", "bun"},
		{"package-lock.json", "npm"},
	} {
		if exists(root, candidate.file) {
			return candidate.name
		}
	}
	return "unknown"
}

func exists(root fs.FS, name string) bool {
	_, err := fs.Stat(root, name)
	return err == nil
}
