package zoau_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Stolkerve/zoau-go"
)

func TestCrud(t *testing.T) {
	firstLine := "Michurao"
	secondLine := "Mucuchies"
	thirdLine := "Guavidia"
	text := "Las cinco aguilas"

	id := os.Getenv("USER")

	DSN := id + ".ZOAU1"

	success, err := zoau.Delete(DSN)
	if success {
		t.Logf("Dataset %s deleted successfuly", DSN)
	} else {
		t.Logf("Dataset %s failed to delete, didn't exist.", DSN)
	}
	if err != nil {
		t.Fatalf("Fail to delete %s", DSN)
	}

	if _, err = zoau.Create(DSN, &zoau.CreateArgs{Type: zoau.String(zoau.DS_ORG_SEQ), PrimarySpace: zoau.String("10")}); err != nil {
		t.Fatalf("Fail to create %s. Err: %v", DSN, err)
	} else {
		t.Logf("Dataset %s created successfuly", DSN)
	}

	if err := zoau.Write(DSN, firstLine, false); err != nil {
		t.Fatalf("Fail to write in %s, the first line", DSN)
	} else {
		t.Logf("Firts line `%s` was append in dataset %s successfuly", firstLine, DSN)
	}

	if output, err := zoau.Read(DSN, nil); err != nil {
		t.Fatalf("Fail to read %s", DSN)
	} else {
		outputTrim := (*output)[:len(firstLine)]
		if outputTrim != firstLine {
			t.Fatalf("expected: %s, got %s", firstLine, *output)
		} else {
			t.Logf("First line match")
		}
	}

	if err := zoau.Write(DSN, secondLine, true); err != nil {
		t.Fatalf("Fail to append in %s, the second line", DSN)
	} else {
		t.Logf("Second line `%s` was append in dataset %s successfuly", secondLine, DSN)
	}
	if err := zoau.Write(DSN, thirdLine, true); err != nil {
		t.Fatalf("Fail to append in %s, the third line", DSN)
	} else {
		t.Logf("Third line `%s` was append in dataset %s successfuly", thirdLine, DSN)
	}

	if output, err := zoau.Read(DSN, nil); err != nil {
		t.Fatalf("Fail to read %s", DSN)
	} else {
		outputSplited := strings.Split(*output, "\n")
		outputTrim := fmt.Sprintf(
			`%s\n%s\n%s`,
			outputSplited[0][:len(firstLine)],
			outputSplited[1][:len(secondLine)],
			outputSplited[2][:len(thirdLine)],
		)
		expected := fmt.Sprintf(`%s\n%s\n%s`, firstLine, secondLine, thirdLine)
		if outputTrim != expected {
			t.Fatalf("expected: %s, got %s", expected, *output)
		} else {
			t.Logf("First, second and third line match")
		}
	}

	var tail uint = 1
	if output, err := zoau.Read(DSN, &zoau.ReadArgs{
		Tail: &tail,
	}); err != nil {
		t.Fatalf("Fail to read from the tail -1 %s", DSN)
	} else {
		outputTrim := (*output)[:len(thirdLine)]
		if outputTrim != thirdLine {
			t.Fatalf("expected: %s, got %s", firstLine, *output)
		} else {
			t.Logf("Tail -1 line match")
		}
	}

	if err := zoau.Write(DSN, text, false); err != nil {
		t.Fatalf("Fail to write in %s, the text: %s", DSN, text)
	} else {
		t.Logf("The text `%s` was write in dataset %s successfuly", text, DSN)
	}

	if output, err := zoau.Read(DSN, nil); err != nil {
		t.Fatalf("Fail to read %s", DSN)
	} else {
		outputTrim := (*output)[:len(text)]
		if outputTrim != text {
			t.Fatalf("expected: %s, got %s", text, *output)
		} else {
			t.Logf("The text match")
		}
	}
}

func TestHlqTempName(t *testing.T) {
	id := os.Getenv("USER")
	hlq, err := zoau.Hlq()
	if err != nil {
		t.Fatal("Fail to fetch hlq")
	} else {
		if *hlq != id {
			t.Fatalf("hlq returned %s, expected %s", *hlq, id)
		} else {
			t.Log("hlq test success")
		}
	}

	if res, err := zoau.TmpName(nil); err != nil {
		t.Fatal("Fail to create a temporary dataset")
	} else {
		numOfDots := strings.Count(*res, ".")
		if !strings.HasPrefix(*res, "MVSTMP.") && len(*res) != 33 || numOfDots < 3 {
			t.Fatalf("string returned by TmpName expected to start with `MVSTMP.`, got %s", *res)
		} else {
			t.Logf("Temporary %s dataset  created successfuly", *res)
		}
	}

	if res, err := zoau.TmpName(hlq); err != nil {
		t.Fatal("Fail create temporary with hlq dataset")
	} else {
		numOfDots := strings.Count(*res, ".")
		if !strings.HasPrefix(*res, *hlq) && len(*res) != 33 || numOfDots != 3 {
			t.Fatalf("string returned by TmpName expected to start with `%s`, got %s", *hlq, *res)
		} else {
			t.Logf("Temporary with hlq %s dataset created successfuly", *res)
		}
	}
}

func TestDelete(t *testing.T) {
	id := os.Getenv("USER")
	dsns := []string{
		id + ".ZOAU3A",
		id + ".ZOAU3B",
		id + ".ZOAU3C",
		id + ".ZOAU3D",
		id + ".ZOAU3E",
	}

	dsnPattern := id + ".ZOAU3*"

	t.Logf("Deleting %s", dsnPattern)
	if _, err := zoau.Delete(dsnPattern); err != nil {
		t.Fatalf("Fail to delete %s", dsnPattern)
	}

	for _, ds := range dsns {
		if _, err := zoau.Create(ds, &zoau.CreateArgs{
			Type: zoau.String(zoau.DS_ORG_SEQ), PrimarySpace: zoau.String("10"),
		}); err != nil {
			t.Fatalf("Fail to create %s", ds)
		} else {
			t.Logf("Dataset %s created successfuly", ds)
		}
	}

	t.Logf("Listing %s", dsnPattern)
	if res, err := zoau.ListingDataset(dsnPattern, nil); err != nil {
		t.Fatalf("Fail to listing %s", dsnPattern)
	} else {
		if len(res) != len(dsns) {
			t.Fatalf("The listing of %s return a set of %d, expected %d", dsnPattern, len(res), len(dsns))
		}
		for i := 0; i < len(res); i++ {
			if res[i].Name != dsns[i] {
				t.Fatalf("Unexpected dataset `%s` in listing of %s", res[i].Name, dsnPattern)
			}
		}
	}
	t.Log("All datasets were listing")

	t.Logf("Deleting %s", dsns[0])
	if _, err := zoau.Delete(dsns[0]); err != nil {
		t.Fatalf("Fail to delete %s", dsns[0])
	}

	t.Logf("Listing %s", dsnPattern)
	if res, err := zoau.ListingDataset(dsnPattern, nil); err != nil {
		t.Fatalf("Fail to listing %s", dsnPattern)
	} else {
		if len(res) != len(dsns[1:]) {
			t.Fatalf("The listing of %s return a set of %d, expected %d", dsnPattern, len(res), len(dsns))
		}
		for i := 0; i < len(res); i++ {
			if res[i].Name != dsns[1:][i] {
				t.Fatalf("Unexpected dataset `%s` in listing of %s", res[i].Name, dsnPattern)
			}
		}
	}
	t.Log("All datasets were listing")

	t.Logf("Deleting %s", dsnPattern)
	if _, err := zoau.Delete(dsnPattern); err != nil {
		t.Fatalf("Fail to delete %s", dsnPattern)
	}

	if res, err := zoau.ListingDataset(dsnPattern, nil); err != nil {
		t.Fatalf("Fail to listing %s", dsnPattern)
	} else {
		if len(res) != 0 {
			t.Fatalf("The listing of %s return a set of %d, expected %d", dsnPattern, len(res), 0)
		}
	}
}

func TestCopy(t *testing.T) {
	id := os.Getenv("USER")
	ds1 := id + ".ZOAU1a"
	ds2 := id + ".ZOAU1b"
	dsp := id + ".ZOAU1?"

	t.Logf("Deleting %s", dsp)
	if _, err := zoau.Delete(dsp); err != nil {
		t.Fatalf("Fail to delete %s", dsp)
	}

	t.Log("Copy of a non-existent USS source file")
	if err := zoau.Copy("/etc/profilxyx", ds1, nil); err == nil {
		t.Fatalf("This copy must fail")
	}

	t.Log("Copy a USS source file")
	if err := zoau.Copy("/etc/profile", ds1, nil); err != nil {
		t.Fatalf("copy failed: %s", err)
	}

	t.Log("Delete work datasets")
	if _, err := zoau.Delete(dsp); err != nil {
		t.Fatalf("Fail to delete %s", dsp)
	}

	t.Log("Create")
	if _, err := zoau.Create(ds1, &zoau.CreateArgs{Type: zoau.String(zoau.DS_ORG_SEQ)}); err != nil {
		t.Fatalf("Fail to create %s", ds1)
	}

	t.Log("Write another line")
	line := "This is the first line"
	if err := zoau.Write(ds1, line, false); err != nil {
		t.Fatalf("Fail to write %s in %s", line, ds1)
	}

	t.Log("Copy dataset as binary")
	if err := zoau.Copy(ds1, ds2, &zoau.CopyArgs{
		Binary: true,
	}); err != nil {
		t.Fatalf("Fail to copy %s in %s", ds1, ds2)
	}

	if res, err := zoau.Read(ds2, &zoau.ReadArgs{
		Tail: zoau.Uint(1),
	}); err != nil {
		t.Fatalf("Fail to listing %s", ds2)
	} else {
		outputTrim := (*res)[:len(line)]
		if outputTrim != line {
			t.Fatalf("expected: %s, got %s", line, *res)
		}
	}
}

func TestCompare(t *testing.T) {
	id := os.Getenv("USER")
	ds1 := id + ".ZOAU1a"
	ds2 := id + ".ZOAU1b"
	ds3 := id + ".ZOAU1d"
	ds4 := id + ".ZOAU1e"
	dsp := id + ".ZOAU1?"

	t.Log("Delete  datasets")
	if _, err := zoau.Delete(dsp); err != nil {
		t.Fatalf("Fail to delete %s", dsp)
	}

	if _, err := zoau.Create(ds1, &zoau.CreateArgs{
		Type: zoau.String(zoau.DS_ORG_SEQ),
	}); err != nil {
		t.Fatalf("Fail to create %s", ds1)
	}
	if _, err := zoau.Create(ds2, &zoau.CreateArgs{
		Type: zoau.String(zoau.DS_ORG_SEQ),
	}); err != nil {
		t.Fatalf("Fail to create %s", ds2)
	}

	t.Log("Write to ds2")
	line := "This is the first line"
	if err := zoau.Write(ds2, line, false); err != nil {
		t.Fatalf("Fail to write %s in %s", line, ds2)
	}

	t.Log("Copy dataset ds2 to ds1")
	if err := zoau.Copy(ds2, ds1, nil); err != nil {
		t.Fatalf("Fail to copy %s in %s", ds2, ds1)
	}

	t.Log("Comparing ds1 to ds2")
	if res, err := zoau.Compare(ds1, ds2, nil); err != nil {
		t.Fatalf("Fail to compere %s to %s", ds2, ds1)
	} else {
		if res != nil {
			t.Log(*res)
			t.Fatalf("%s and %s must be equals", ds2, ds1)
		}
	}

	t.Log("Write extra line to ds2")
	if err := zoau.Write(ds2, line, true); err != nil {
		t.Fatalf("Fail to write %s in %s", line, ds2)
	}

	if res, err := zoau.Compare(ds1, ds2, nil); err != nil {
		t.Fatalf("Fail to compere %s to %s", ds2, ds1)
	} else {
		if res == nil {
			t.Fatalf("%s and %s must not be equals", ds2, ds1)
		}
	}

	t.Log("Write extra line to ds1")
	if err := zoau.Write(ds1, strings.ToUpper(line), true); err != nil {
		t.Fatalf("Fail to write %s in %s", strings.ToUpper(line), ds1)
	}

	t.Log("Comparing ds1 to ds2")
	if res, err := zoau.Compare(ds1, ds2, &zoau.CompareArgs{IgnoreCase: true}); err != nil {
		t.Fatalf("Fail to compere %s to %s", ds2, ds1)
	} else {
		if res != nil {
			t.Log(*res)
			t.Fatalf("%s and %s must be equals", ds2, ds1)
		}
	}

	t.Log("Comparing ds1 to ds2")
	if res, err := zoau.Compare(ds1, ds2, &zoau.CompareArgs{IgnoreCase: false}); err != nil {
		t.Fatalf("Fail to compere %s to %s", ds2, ds1)
	} else {
		if res == nil {
			t.Fatalf("%s and %s must not be equals", ds2, ds1)
		}
	}

	t.Log("Comparing ds1 to ds2")
	if res, err := zoau.Compare(ds1, ds2, &zoau.CompareArgs{Columnns: &zoau.Point{
		Start: 1, End: 4,
	}, IgnoreCase: true}); err != nil {
		t.Fatalf("Fail to compere %s to %s", ds2, ds1)
	} else {
		if res != nil {
			t.Log(*res)
			t.Fatalf("%s and %s must be equals", ds2, ds1)
		}
	}

	t.Log("Comparing ds3 to ds4")
	if _, err := zoau.Compare(ds3, ds4, &zoau.CompareArgs{IgnoreCase: false}); err == nil {
		t.Fatalf("Fail to compere %s to %s", ds3, ds4)
	}
}

func TestFindReplace(t *testing.T) {
	id := os.Getenv("USER")
	ds := id + ".ZOAU5.SEQ"
	lines := "This is the first line.\nThis is the second line.\nThis is the third line."
	expectedLines := []string{"This was the first LINE", "This was the second LINE", "This was the third LINE"}

	t.Logf(`Delete %s if exists`, ds)
	if _, err := zoau.Delete(ds); err != nil {
		t.Fatalf("Fail to delete %s", ds)
	}

	if _, err := zoau.Create(ds, &zoau.CreateArgs{Type: zoau.String(zoau.DS_ORG_SEQ), PrimarySpace: zoau.String("10")}); err != nil {
		t.Fatalf("Fail to create %s. Err: %v", ds, err)
	} else {
		t.Logf("Dataset %s created successfuly", ds)
	}

	t.Log("Write lines")
	if err := zoau.Write(ds, lines, false); err != nil {
		t.Fatalf("Fail to write %s in %s. Err: %s", lines, ds, err)
	}

	t.Log("FindReplace")
	if err := zoau.FindReplace(ds, "line.", "LINE"); err != nil {
		t.Fatalf("Fail to find and replace. Err %s", err)
	}

	t.Log("FindReplace")
	if err := zoau.FindReplace(ds, "This is", "This was"); err != nil {
		t.Fatalf("Fail to find and replace. Err %s", err)
	}

	t.Log("Verify")
	if res, err := zoau.Read(ds, nil); err != nil {
		t.Fatalf("Fail to read %s", ds)
	} else {
		outputLines := strings.Split(*res, "\n")
		for i := 0; i < len(expectedLines); i++ {
			if outputLines[i][:len(expectedLines[i])] != expectedLines[i] {
				t.Fatalf("expected: %s, got %s", expectedLines[i], outputLines[i])
			}
		}
	}
}

func TestListingParse(t *testing.T) {
	detailsArgs := `
Z38816.LOAD                                  2023/11/06 po  U        0  4096 ZXPM01           ??       283320
BGYSC2006I Unable to obtain dataset information for dataset Z38816.P0397638.T0524969.C0000001 on volume ZXPM06.
`
	expectedParsedValues := []string{
		"Z38816.LOAD", "2023/11/06", "po", "U", "0", "4096", "ZXPM01", "??", "283320",
	}
	expectedOutput := zoau.Dataset{
		Name:           "Z38816.LOAD",
		LastReferenced: "2023/11/06",
		Dsorg:          "po",
		Recfm:          "U",
		Lrecl:          0,
		BlockSize:      4096,
		Volume:         "ZXPM01",
		UsedSpace:      nil,
		TotalSpace:     283320,
	}
	unparsedLines := strings.Split(detailsArgs, "\n")
	for _, line := range unparsedLines {
		parsedLine := zoau.ParseLine(line)
		if len(parsedLine) == 9 {
			for i, parsedValue := range parsedLine {
				if expectedParsedValues[i] != parsedValue {
					t.Fatalf("%s is not equal to %s", expectedParsedValues[i], parsedValue)
				}
			}
			output, err := zoau.ParseListingValues(parsedLine)
			if err != nil {
				t.Fatal(err)
			}

			if expectedOutput.Name != output.Name {
				t.Fatalf("Expected Name: %s. go %s", expectedOutput.Name, output.Name)
			}
			if expectedOutput.LastReferenced != output.LastReferenced {
				t.Fatalf("Expected LastReferenced: %s. go %s", expectedOutput.LastReferenced, output.LastReferenced)
			}
			if expectedOutput.Dsorg != output.Dsorg {
				t.Fatalf("Expected Dsorg: %s. go %s", expectedOutput.Dsorg, output.Dsorg)
			}
			if expectedOutput.Recfm != output.Recfm {
				t.Fatalf("Expected Recfm: %s. go %s", expectedOutput.Recfm, output.Recfm)
			}
			if expectedOutput.Lrecl != output.Lrecl {
				t.Fatalf("Expected Lrecl: %v. go %v", expectedOutput.Lrecl, output.Lrecl)
			}
			if expectedOutput.BlockSize != output.BlockSize {
				t.Fatalf("Expected BlockSize: %v. go %v", expectedOutput.BlockSize, output.BlockSize)
			}
			if expectedOutput.Volume != output.Volume {
				t.Fatalf("Expected Volume: %s. go %s", expectedOutput.Volume, output.Volume)
			}
			if expectedOutput.UsedSpace != output.UsedSpace {
				t.Fatalf("Expected UsedSpace: %v. go %v", expectedOutput.UsedSpace, output.UsedSpace)
			}
			if expectedOutput.TotalSpace != output.TotalSpace {
				t.Fatalf("Expected TotalSpace: %v. go %v", expectedOutput.TotalSpace, output.TotalSpace)
			}
		}
	}
}
