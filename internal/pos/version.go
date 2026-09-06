package pos

import "runtime/debug"

// BuildRevision reads the git commit Go embeds in the binary at build time
// (buildvcs, on by default since Go 1.18), so the landing page can show
// which commit is actually running with no custom build step. The Pi does
// not pull from git or build (scripts/deploy-pi.sh): the binary that lands
// there is the binary that runs, so this is the only way to tell, from the
// tablet itself, which commit that was. Revision is empty when the binary
// was not built inside a git checkout; dirty is true when the working tree
// that built it had uncommitted changes.
func BuildRevision() (revision string, dirty bool) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", false
	}
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
			if len(revision) > 7 {
				revision = revision[:7]
			}
		case "vcs.modified":
			dirty = setting.Value == "true"
		}
	}
	return revision, dirty
}
