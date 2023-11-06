package buildtime

import (
	_ "embed"
	"strings"
	"time"
)

//
// -- VERSION
//

//go:generate sh get_version.sh

//go:embed version.txt
var version string

func Version() string {
	return strings.TrimSpace(version)
}

//
// -- COMPILE TIME
//

//go:generate sh -c "TZ=UTC date > compile_time.txt"

//go:embed compile_time.txt
var compileTime string

func CompiledAt() time.Time {
	val := strings.TrimSpace(compileTime)

	date, err := time.Parse(time.UnixDate, val)
	if err != nil {
		return time.Time{}
	}

	return date
}
