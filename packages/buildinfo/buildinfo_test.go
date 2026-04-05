package buildinfo

import "testing"

func TestSummaryUsesDefaults(t *testing.T) {
	oldVersion := Version
	oldCommit := Commit
	oldBuildDate := BuildDate
	Version = ""
	Commit = ""
	BuildDate = ""
	defer func() {
		Version = oldVersion
		Commit = oldCommit
		BuildDate = oldBuildDate
	}()

	if got := VersionLine(); got != "iot4b version dev" {
		t.Fatalf("VersionLine() = %q", got)
	}

	if got := Summary(); got != "iot4b version dev (commit unknown, built unknown)" {
		t.Fatalf("Summary() = %q", got)
	}
}
