package buildinfo

import "fmt"

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func versionValue() string {
	if Version == "" {
		return "dev"
	}
	return Version
}

func commitValue() string {
	if Commit == "" {
		return "unknown"
	}
	return Commit
}

func buildDateValue() string {
	if BuildDate == "" {
		return "unknown"
	}
	return BuildDate
}

func VersionLine() string {
	return fmt.Sprintf("iot4b version %s", versionValue())
}

func Summary() string {
	return fmt.Sprintf("%s (commit %s, built %s)", VersionLine(), commitValue(), buildDateValue())
}
