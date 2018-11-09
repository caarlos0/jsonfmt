package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/pmezard/go-difflib/difflib"
)

var (
	version = "master"
	commit  = "none"
	date    = "unknown"

	app   = kingpin.New("jsonfmt", "Like gofmt, but for JSON")
	files = app.Arg("files", "glob of the files you want to check").Strings()

	write    = app.Flag("write", "write changes to the files").Short('w').Bool()
	indent   = app.Flag("indent", "characteres used to indend the json").Short('i').Default("  ").String()
	failfast = app.Flag("failfast", "fail on the first error").Bool()
)

func main() {
	app.Version(fmt.Sprintf("%v, commit %v, built at %v", version, commit, date))
	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')

	kingpin.MustParse(app.Parse(os.Args[1:]))

	if len(*files) == 0 {
		app.FatalUsage("list of files to check is empty")
	}

	var failed bool
	for _, file := range *files {
		bts, err := ioutil.ReadFile(file)
		app.FatalIfError(err, "failed to read file: %s", file)
		var out bytes.Buffer
		err = json.Indent(&out, bytes.TrimSpace(bts), "", *indent)
		app.FatalIfError(err, "failed to format json file: %s", file)
		out.Write([]byte{'\n'})
		if bytes.Equal(bts, out.Bytes()) {
			continue
		}
		if *write {
			err := ioutil.WriteFile(file, out.Bytes(), 0)
			app.FatalIfError(err, "failed to write json file: %s", file)
			continue
		}

		diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
			A:        difflib.SplitLines(string(bts)),
			B:        difflib.SplitLines(string(out.Bytes())),
			FromFile: "Original",
			ToFile:   "Formatted",
			Context:  3,
		})
		app.FatalIfError(err, "failed to diff file: %s", file)
		app.Errorf("file %s differs:\n%s\n", file, diff)
		failed = true
		if *failfast {
			break
		}
	}
	if failed {
		app.Fatalf("some files are not properly formated, check above")
	}
}
