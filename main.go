package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alecthomas/kingpin"
	zglob "github.com/mattn/go-zglob"
	"github.com/pmezard/go-difflib/difflib"
)

var (
	version = "master"
	commit  = "none"
	date    = "unknown"

	app   = kingpin.New("jsonfmt", "Like gofmt, but for JSON")
	write = app.Flag("write", "write changes to the files").Short('w').Bool()
	globs = app.Arg("files", "glob of the files you want to check").Default("**/*.json").Strings()
)

func main() {
	app.Version(fmt.Sprintf("%v, commit %v, built at %v", version, commit, date))
	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')

	kingpin.MustParse(app.Parse(os.Args[1:]))

	var failed bool
	for _, glob := range *globs {
		matches, err := zglob.Glob(glob)
		app.FatalIfError(err, "failed to parse glob: %s", glob)

		for _, match := range matches {
			file, err := os.Open(match)
			app.FatalIfError(err, "failed to open file: %s", match)
			defer file.Close()
			bts, err := ioutil.ReadAll(file)
			app.FatalIfError(err, "failed to read file: %s", match)
			_ = file.Close()

			var tmp json.RawMessage
			app.FatalIfError(json.Unmarshal(bts, &tmp), "failed to parse json file: %s", match)
			out, err := json.MarshalIndent(tmp, "", "  ")                 // TODO: support to customize indent
			app.FatalIfError(err, "failed to parse json file: %s", match) // XXX: improve error msg
			out = append(out, '\n')
			if !bytes.Equal(bts, out) {
				if *write {
					app.FatalIfError(ioutil.WriteFile(match, out, 0), "failed to write json file: %s", match)
					continue
				}

				diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
					A:        difflib.SplitLines(string(bts)),
					B:        difflib.SplitLines(string(out)),
					FromFile: "Original",
					ToFile:   "Formatted",
					Context:  3,
				})
				app.FatalIfError(err, "failed to diff file: %s", match)
				app.Errorf("file %s differs:\n%s\n", match, diff)
				failed = true
			}
		}
	}
	if failed {
		app.Fatalf("some files are not properly formated, check above")
	}
}
