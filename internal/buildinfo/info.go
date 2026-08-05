// Package buildinfo exposes immutable build metadata supplied by the linker.
package buildinfo

// These values are replaced by release builds through -ldflags. Keep safe
// development defaults so local builds are honest and reproducible.
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

// Info identifies one Aruo build.
type Info struct {
	Version string
	Commit  string
	Date    string
}

// Current returns the metadata embedded in the running binary.
func Current() Info {
	return Info{Version: version, Commit: commit, Date: date}
}
