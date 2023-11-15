package zoau

import (
	"fmt"
)

// Execute an MVS Program.
// Returns the stdout or the stderr and the return code
func Execute(pgm string, pgmArgs *string, dds []DDStatement, args *Args) (string, int) {
	options := make([]string, 0)

	options = append(options, parseUniversalArgs(*args)...)

	if pgmArgs != nil {
		options = append(options, "--args="+*pgmArgs)
	}

	options = append(options, "--pgm="+pgm)

	for _, dd := range dds {
		options = append(options, fmt.Sprintf("--%s=%s", dd.Name, dd.Definition.buildArgsString()))
	}
	out, rc, _ := execZaouCmd("mvscmd", options)
	return out, rc
}

// Execute an authorized MVS Program.
// Returns the stdout or the stderr and the return code
func ExecuteAuthorized(pgm string, pgmArgs *string, dds []DDStatement, args *Args) (string, int) {
	options := make([]string, 0)

	options = append(options, parseUniversalArgs(*args)...)
	if pgmArgs != nil {
		options = append(options, "--args="+*pgmArgs)
	}

	options = append(options, "--pgm="+pgm)

	for _, dd := range dds {
		options = append(options, fmt.Sprintf("--%s=%s", dd.Name, dd.Definition.buildArgsString()))
	}

	out, rc, _ := execZaouCmd("mvscmdauth", options)
	return out, rc
}
