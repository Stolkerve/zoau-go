package zoau

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ZOAU dmod function to be used by zos_blockinfile Ansible module
func BlockInFile(dataset string, args *BlockInFileArgs) error {
	options := []string{"-b"}
	state := true

	if args != nil {
		if args.Lock {
			options = append(options, "-l")
		}
		if args.Force {
			options = append(options, "-f")
		}
		if args.State != nil {
			state = *args.State
		}
		if args.Encoding != nil {
			options = append(options, "-c", *args.Encoding)
		}
		if args.Market != nil {
			options = append(options, "-m", *args.Market)
		}
	}

	if state {
		if args.Block == nil {
			return errors.New("Block is required when state=true.")
		} else if args.InsAft != nil {
			if *args.InsAft == "EOF" {
				options = append(
					options,
					fmt.Sprintf(`"$ a\%s"`, *args.Block),
					dataset,
				)
			} else {
				options = append(options,
					"-s",
					"-e",
					fmt.Sprintf(`"/%s/a\\%s/$"`, *args.InsAft, *args.Block),
					"-e",
					fmt.Sprintf(`"$ a\%s"`, *args.Block),
					dataset,
				)
			}
		} else if args.InsBef != nil {
			if *args.InsBef == "BOF" {
				options = append(
					options,
					fmt.Sprintf(`"1 i\\%s"`, *args.Block),
					dataset,
				)
			} else {
				options = append(
					options,
					"-s",
					"-e",
					fmt.Sprintf(`"/%s/i\\%s/$"`, *args.InsBef, *args.Block),
					"-e",
					fmt.Sprintf(`"$ a\\%s"`, *args.Block),
					dataset,
				)
			}
		} else {
			return errors.New("InsAft or InsBef is required when state=true")
		}
	} else {
		options = append(
			options,
			"//d",
			dataset,
		)
	}

	_, _, err := execZaouCmd("dmod", options)
	if err != nil {
		return err
	}
	return nil
}

// Compare two datasets, output the ISRSUPC output.
func Compare(source string, target string, args *CompareArgs) (*string, error) {
	options := make([]string, 0)
	if args != nil {
		if args.IgnoreCase {
			options = append(options, "-i")
		}
		if args.Columnns != nil {
			options = append(options, "-c", fmt.Sprintf("%d:%d", args.Columnns.Start, args.Columnns.End))
		}
		if args.Lines != nil {
			options = append(options, "-C", fmt.Sprintf("%d:%d", args.Lines.Start, args.Lines.End))
		}
	}
	options = append(options, source, target)

	stdout, rc, err := execZaouCmd("ddiff", options)

	if rc == 0 {
		return nil, nil
	} else if rc == 1 {
		return &stdout, nil
	}
	return nil, err
}

// Copy a z/OS source (dataset, HFS file) to a z/OS target.
func Copy(source string, target string, args *CopyArgs) error {
	options := make([]string, 0)
	if args != nil {
		if args.Force {
			options = append(options, "-f")
		}
		if args.Alias {
			options = append(options, "-I")
		}
		if args.Executable {
			options = append(options, "-X")
		}
		if args.Binary {
			options = append(options, "-B")
		}
		if args.TextMode {
			options = append(options, "-T")
		}
	}

	options = append(options, source, target)

	_, _, err := execZaouCmd("dcp", options)
	return err
}

// Create a z/OS dataset.
func Create(name string, args *CreateArgs) (*Dataset, error) {
	options := make([]string, 0)

	if args != nil {
		if args.Type != nil {
			options = append(options, "-t", *args.Type)
		}
		if args.PrimarySpace != nil {
			options = append(options, "-s", *args.PrimarySpace)
		}
		if args.SecondarySpace != nil {
			options = append(options, "-e", *args.SecondarySpace)
		}
		if args.DirectoryBlocks != nil {
			options = append(options, "-b", fmt.Sprintf("%d", *args.DirectoryBlocks))
		}
		if args.BlockSize != nil {
			options = append(options, "-B", fmt.Sprintf("%d", *args.BlockSize))
		}
		if args.RecordFormat != nil {
			options = append(options, "-r", *args.RecordFormat)
		}
		if args.RecordLength != nil {
			options = append(options, "-l", fmt.Sprintf("%d", *args.RecordLength))
		}
		if args.StorageClassName != nil {
			options = append(options, "-c", *args.StorageClassName)
		}
		if args.DataClassName != nil {
			options = append(options, "-D", *args.DataClassName)
		}
		if args.Keys != nil {
			options = append(options, "-k", fmt.Sprintf("%d:%d", args.Keys.KeyLength, args.Keys.KeyOffset))
		}
		if args.Volumes != nil {
			options = append(options, "-V", *args.Volumes)
		}
	}

	options = append(options, name)

	_, rc, err := execZaouCmd("dtouch", options)
	if rc >= 8 {
		return nil, err
	}
	if out, err := ListingDataset(name, nil); err != nil {
		return nil, err
	} else {
		return &out[0], nil
	}
}

// Delete a z/OS dataset.
// Return false if dataset specified by dataset pattern does not exist.
// Otherwise task completed without error.
func Delete(datasets ...string) (bool, error) {
	_, returnCode, err := execZaouCmd("drm", datasets)

	if returnCode == 1 {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// Delete members contained in a dataset.
// Return false if dataset specified by dataset pattern does not exist.
// Otherwise task completed without error.
func DeleteMember(pattern string) (bool, error) {
	options := make([]string, 0)
	options = append(options, pattern)

	_, returnCode, err := execZaouCmd("mrm", options)

	if returnCode == 1 {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// Check whether or not a dataset exists.
func Exist(dataset string) (bool, error) {
	if out, err := ListingDataset(dataset, nil); err != nil {
		return false, err
	} else {
		return len(out) > 0, nil
	}
}

// Find dataset that contains member within a concatenation. Returns the first dataset that contains member.
// Return the dataset containing the member.
func FindMember(member string, concatentation string) (string, error) {
	return execSimpleStringCmd("dwhence", []string{member, concatentation})
}

// Replace text within a dataset.
func FindReplace(dataset string, find string, replace string) error {
	options := []string{
		fmt.Sprintf(`"s/%s/%s/g"`, find, replace),
		dataset,
	}
	_, _, err := execZaouCmd("dsed", options)
	return err
}

// ZOAU dsed function to be used by zos_lineinfile Ansible module.
func LineInFile(dataset string, line string, args *LineInFileArgs) error {
	options := make([]string, 0)
	state := true
	matchCharacter := "$"

	if args.State != nil {
		state = *args.State
	}

	if args.Lock {
		options = append(options, "-l")
	}

	if args.Force {
		options = append(options, "-f")
	}

	if args.Encoding != nil {
		options = append(options, "-c", *args.Encoding)
	}

	if args.FirstMatch {
		matchCharacter = "1"
	}

	if state {
		if args.Regex != nil {
			if args.InsAft != nil {
				if *args.InsAft == "EOF" {
					options = append(options,
						"-s",
						"-e",
						fmt.Sprintf(`"/%s/c\\%s/%s"`, *args.Regex, line, matchCharacter),
						"-e",
						fmt.Sprintf(`"$ a\\%s"`, line),
						dataset,
					)
				} else {
					options = append(options,
						"-s",
						"-e",
						fmt.Sprintf(`"/%s/c\\%s/%s"`, *args.Regex, line, matchCharacter),
						"-e",
						fmt.Sprintf(`"/%s/a\\%s/%s"`, *args.InsAft, line, matchCharacter),
						"-e",
						fmt.Sprintf(`"$ a\\%s"`, line),
						dataset,
					)
				}
			} else if args.InsBef != nil {
				if *args.InsBef == "BOF" {
					options = append(options,
						"-s",
						"-e",
						fmt.Sprintf(`"/%s/c\\%s/%s"`, *args.Regex, line, matchCharacter),
						"-e",
						fmt.Sprintf(`"1 i\\%s"`, line),
						dataset,
					)
				} else {
					options = append(options,
						"-s",
						"-e",
						fmt.Sprintf(`"/%s/c\\%s/%s"`, *args.Regex, line, matchCharacter),
						"-e",
						fmt.Sprintf(`"/%s/i\\%s/%s"`, *args.InsBef, line, matchCharacter),
						"-e",
						fmt.Sprintf(`"$ a\\%s"`, line),
						dataset,
					)
				}
			} else {
				options = append(options,
					fmt.Sprintf(`"/%s/c\\%s/%s"`, *args.Regex, line, matchCharacter),
					dataset,
				)
			}
		} else {
			if args.InsAft != nil {
				if *args.InsAft == "EOF" {
					options = append(options,
						fmt.Sprintf(`"$ a\\%s"`, line),
						dataset,
					)
				} else {
					options = append(options,
						"-s",
						"-e",
						fmt.Sprintf(`"/%s/a\\%s/%s"`, *args.InsAft, line, matchCharacter),
						"-e",
						fmt.Sprintf(`"$ a\\%s"`, line),
						dataset,
					)
				}
			} else if args.InsBef != nil {
				if *args.InsBef == "BOF" {
					options = append(options,
						fmt.Sprintf(`"1 i\\%s"`, line),
						dataset,
					)
				} else {
					options = append(options,
						"-s",
						"-e",
						fmt.Sprintf(`"/%s/i\\%s/%s"`, *args.InsBef, line, matchCharacter),
						"-e",
						fmt.Sprintf(`"$ a\\%s"`, line),
						dataset,
					)
				}
			} else {
				return errors.New("Incorrect parameters")
			}
		}
	} else {
		if args.Regex != nil {
			if len(line) != 0 {
				options = append(options,
					"-s",
					"-e",
					fmt.Sprintf(`"/%s/d"`, *args.Regex),
					"-e",
					fmt.Sprintf(`"/%s/d"`, line),
					dataset,
				)
			} else {
				options = append(options,
					fmt.Sprintf(
						`"/%s/d"`,
						*args.Regex,
					),
					dataset,
				)
			}
		} else {
			options = append(options,
				fmt.Sprintf(`"/%s/d"`, line),
				dataset,
			)
		}
	}

	_, _, err := execZaouCmd("dsed", options)
	return err
}

// Get a list of members from a dataset.
func ListMembers(pattern string) ([]string, error) {
	options := []string{pattern}
	stdout, returnCode, err := execZaouCmd("mls", options)
	if err != nil {
		return nil, err
	}

	if returnCode != 0 {
		return nil, nil
	}

	return strings.Split(stdout, "\n"), nil
}

// Returns a listing of the datasets matching the supplied pattern.
func ListingDataset(pattern string, args *ListingArgs) ([]Dataset, error) {
	options := []string{"-l", "-u", "-s", "-b"}
	if args != nil {
		if args.NameOnly {
			options = []string{}
		}
		if args.Migrate && args.NameOnly {
			options = []string{"-m"}
		} else {
			return nil, errors.New("To display migrated datasets, requires NameOnly to be true.")
		}
	}

	options = append(options, pattern)

	stdout, returnCode, err := execZaouCmd("dls", options)
	if returnCode == 1 {
		return []Dataset{}, nil
	}

	if err != nil {
		return nil, err
	}

	output := make([]Dataset, 0)

	unparsedLines := strings.Split(stdout, "\n")
	for _, unparsedLine := range unparsedLines {
		parsedLine := ParseLine(unparsedLine)
		if len(parsedLine) == 1 {
			output = append(output, Dataset{
				Name: parsedLine[0],
			})
			continue
		}
		if len(parsedLine) == 9 {
			v, err := ParseListingValues(parsedLine)
			if err != nil {
				return nil, errors.New("Unexpected error occurred while parsing")
			}
			if args != nil {
				if args.Volume != nil {
					if *args.Volume == v.Volume {
						output = append(output, v)
					}
					continue
				}
			}
			output = append(output, v)
		}
	}
	return output, nil
}

// Move (rename) a dataset.
func Move(source string, target string) error {
	options := []string{source, target}

	_, _, err := execZaouCmd("dmv", options)
	return err
}

// Move (rename) a member.
func MoveMember(dataset string, source string, target string) error {
	options := []string{dataset, source, target}

	_, _, err := execZaouCmd("mmv", options)
	return err
}

// Get the string contents of a dataset.
func Read(dataset string, args *ReadArgs) (string, error) {
	options := make([]string, 0)
	if args != nil {
		if args.FromLine != nil {
			options = append(options, "-n", fmt.Sprintf("+%d", *args.FromLine))
		} else if args.Tail != nil {
			options = append(options, "-n", fmt.Sprintf("-%d", *args.Tail))
		} else {
			options = append(options, "-n", "+1")
		}
	}

	options = append(options, dataset)

	return execSimpleStringCmd("dtail", options)
}

// A function to display the head content of a non-VSAM dataset. Gets the head content of a dataset.
// Nlines: Read the first nlines lines from the dataset.
func ReadHead(dataset string, Nlines *uint) (string, error) {
	options := []string{"-n", "+1"}

	options = append(options, dataset)

	if Nlines != nil {
		options = []string{"|", "head", "-n", fmt.Sprintf("%d", *Nlines)}
	}

	return execSimpleStringCmd("dtail", options)
}

// Search a dataset using ISRSUPC.
func Search(dataset string, value string, args *SearchArgs) (*string, error) {
	options := make([]string, 0)
	if args != nil {
		if args.DisplayLines {
			options = append(options, "-n")
		}
		if args.IgnoreCase {
			options = append(options, "-i")
		}
		if args.PrintDatasets {
			options = append(options, "-v")
		}
		if args.Lines != nil {
			options = append(options, "-C", fmt.Sprintf("%d", *args.Lines))
		}
	}

	options = append(options, value, dataset)

	stdout, _, err := execZaouCmd("dgrep", options)
	if err != nil {
		return nil, err
	}
	if args.CountLines {
		out := strconv.Itoa(strings.Count(stdout, "\n"))
		return &out, nil
	}
	return &stdout, nil
}

// Return the high level qualifier (HLQ) of the active TSO environment
func Hlq() (string, error) {
	return execSimpleStringCmd("hlq", nil)
}

// Creates a temporary dataset name.
// hlq:     The HLQ of the temporary dataset name.
func TmpName(hlq *string) (string, error) {
	options := make([]string, 0)
	if hlq != nil {
		options = append(options, *hlq)
	}
	return execSimpleStringCmd("mvstmp", options)
}

// Unzips a .dzp file.
func UnZip(file string, hlq string, args *UnZipArgs) error {
	options := make([]string, 0)
	if args != nil {
		if args.Size != nil {
			options = append(options, "-s", *args.Size)
		}
		if args.Volume != nil {
			options = append(options, "-V")
		}
		if args.Dataset {
			options = append(options, "-D")
		}
		if args.Overwrite {
			options = append(options, "-o")
		}
		if args.SmsForTmp {
			options = append(options, "-u")
		}
		if args.Include != nil {
			options = append(options, "-i", *args.Include)
		}
		if args.Exclude != nil {
			options = append(options, "-e", *args.Exclude)
		}
		if args.StorageClassName != nil {
			options = append(options, " -S", *args.StorageClassName)
		}
		if args.ManagementClassName != nil {
			options = append(options, "-m", *args.ManagementClassName)
		}
		if args.DestVolume != nil {
			options = append(options, "-t", *args.DestVolume)
		}

	}
	options = append(options, "-H", hlq, file)
	if args != nil {
		if args.Volume != nil {
			options = append(options, *args.Volume)
		} else if args.SrcVolume != nil {
			options = append(options, *args.SrcVolume)
		}
	}

	_, _, err := execZaouCmd("dunzip", options)
	return err
}

// Write content to a z/OS data set.
func Write(dataset string, content string, _append bool) error {
	options := make([]string, 0)
	if _append {
		options = append(options, "-a")
	}
	options = append(options, content, dataset)
	_, _, err := execZaouCmd("decho", options)
	return err
}

// Zip datasets into an HFS file.
func Zip(file string, target string, args *ZipArgs) error {
	options := make([]string, 0)
	if args != nil {
		if args.Size != nil {
			options = append(options, "-s", *args.Size)
		}
		if args.Volume != nil {
			options = append(options, "-V")
		}
		if args.Force {
			options = append(options, "-f")
		}
		if args.Overwrite {
			options = append(options, "-o")
		}
		if args.Dataset {
			options = append(options, "-D")
		}
		if args.Exclude != nil {
			options = append(options, "-e", *args.Exclude)
		}
		if args.StorageClassName != nil {
			options = append(options, "-S", *args.StorageClassName)
		}
		if args.ManagementClassName != nil {
			options = append(options, "-m", *args.ManagementClassName)
		}
		if args.DestVolume != nil {
			options = append(options, "-t", *args.DestVolume)
		}
	}

	options = append(options, file, target)

	if args != nil {
		if args.SrcVolume != nil {
			options = append(options, *args.SrcVolume)
		}
	}

	_, _, err := execZaouCmd("dzip", options)
	return err
}
