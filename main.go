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

// nolint: gochecknoglobals
var (
	version = "master"

	app      = kingpin.New("jsonfmt", "Like gofmt, but for JSON")
	globs    = app.Arg("files", "glob of the files you want to check").Default("**/*.json").Strings()
	write    = app.Flag("write", "write changes to the files").Short('w').Bool()
	indent   = app.Flag("indent", "characteres used to indend the json").Short('i').Default("  ").String()
	failfast = app.Flag("failfast", "fail on the first error").Bool()
)

func main() {
	app.Version(fmt.Sprintf("%s version %s", app.Name, version))
	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')

	kingpin.MustParse(app.Parse(os.Args[1:]))

	var failed bool
Globs:
	for _, glob := range *globs {
		matches, err := zglob.Glob(glob)
		app.FatalIfError(err, "failed to parse glob: %s", glob)
		if len(matches) == 0 {
			app.Errorf("no matches found: %s", glob)
			failed = true
			if *failfast {
				break Globs
			}
		}

		for _, match := range matches {
			bts, err := ioutil.ReadFile(match)
			app.FatalIfError(err, "failed to read file: %s", match)

			var out bytes.Buffer
			err = json.Indent(&out, bytes.TrimSpace(bts), "", *indent)
			app.FatalIfError(err, "failed to format json file: %s", match)
			out.Write([]byte{'\n'})

			if bytes.Equal(bts, out.Bytes()) {
				continue
			}

			if *write {
				err := ioutil.WriteFile(match, out.Bytes(), 0)
				app.FatalIfError(err, "failed to write json file: %s", match)
				continue
			}

			diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
				A:        difflib.SplitLines(string(bts)),
				B:        difflib.SplitLines(out.String()),
				FromFile: "Original",
				ToFile:   "Formatted",
				Context:  3,
			})
			app.FatalIfError(err, "failed to diff file: %s", match)
			app.Errorf("file %s differs:\n%s\n", match, diff)

			failed = true
			if *failfast {
				break Globs
			}
		}
	}
	if failed {
		app.Fatalf("some files may not be properly formated, check above")
	}
}
