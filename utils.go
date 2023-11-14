package zoau

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

// String returns a pointer value for the string value passed in.
func String(v string) *string {
	return &v
}

func Uint(v uint) *uint {
	return &v
}

func execZaouCmd(proc string, params []string) (string, int, error) {
	//fmt.Println(proc, params)
	cmd := exec.Command(proc, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), cmd.ProcessState.ExitCode(), errors.New(string(output))
	}

	return string(output), cmd.ProcessState.ExitCode(), nil
}

func ParseLine(line string) []string {
	return strings.Fields(line)
}

func ParseListingValues(parsedLine []string) (Dataset, error) {
	lrecl, err := strconv.Atoi(parsedLine[4])
	if err != nil {
		return Dataset{}, err
	}

	blockSize, err := strconv.Atoi(parsedLine[5])
	if err != nil {
		return Dataset{}, err
	}

	var usedSpace *int
	usedSpaceTmp, err := strconv.Atoi(parsedLine[7])
	if err != nil {
		if parsedLine[7] == "??" {
			usedSpace = nil
		} else {
			return Dataset{}, err
		}
	} else {
		usedSpace = &usedSpaceTmp
	}

	totalSpace, err := strconv.Atoi(parsedLine[8])
	if err != nil {
		return Dataset{}, err
	}

	return Dataset{
		Name:           parsedLine[0],
		LastReferenced: parsedLine[1],
		Dsorg:          parsedLine[2],
		Recfm:          parsedLine[3],
		Lrecl:          lrecl,
		BlockSize:      blockSize,
		Volume:         parsedLine[6],
		UsedSpace:      usedSpace,
		TotalSpace:     totalSpace,
	}, nil
}
