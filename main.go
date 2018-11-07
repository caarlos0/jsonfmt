package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alecthomas/kingpin"
	zglob "github.com/mattn/go-zglob"
	"github.com/sergi/go-diff/diffmatchpatch"
	json "github.com/virtuald/go-ordered-json"
)

var (
	version = "master"
	commit  = "none"
	date    = "unknown"

	app   = kingpin.New("jsonfmt", "Like gofmt, but for JSON")
	files = app.Flag("files", "glob of the files you want to check").Short('f').Default("**/*.json").String()
	write = app.Flag("write", "write changes to the files").Short('w').Bool()
)

func main() {
	app.Version(fmt.Sprintf("%v, commit %v, built at %v", version, commit, date))
	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')

	kingpin.MustParse(app.Parse(os.Args[1:]))

	matches, err := zglob.Glob(*files)
	app.FatalIfError(err, "failed to parse glob: %s", *files)

	var dmp = diffmatchpatch.New()

	for _, match := range matches {
		file, err := os.Open(match)
		app.FatalIfError(err, "failed to open file: %s", match)
		defer file.Close()
		bts, err := ioutil.ReadAll(file)
		app.FatalIfError(err, "failed to read file: %s", match)
		_ = file.Close()

		var tmp json.OrderedObject
		app.FatalIfError(json.Unmarshal(bts, &tmp), "failed to parse json file: %s", match)
		out, err := json.MarshalIndent(tmp, "", "  ")
		app.FatalIfError(err, "failed to parse json file: %s", match) // XXX: improve error msg
		out = append(out, '\n')
		if !bytes.Equal(bts, out) {
			if *write {
				app.FatalIfError(ioutil.WriteFile(match, out, 0), "failed to write json file: %s", match)
				continue
			}
			diffs := dmp.DiffMain(string(bts), string(out), false)
			app.Errorf("file %s differs:\n%s\n", match, dmp.DiffPrettyText(diffs))
		}
	}
}
