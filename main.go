package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/goreleaser/fileglob"
	"github.com/hashicorp/go-multierror"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/cobra"
)

var (
	version  = "master"
	write    bool
	failfast bool
	indent   string
	rootCmd  = &cobra.Command{
		Use:          "jsonfmt",
		Short:        "Like gofmt, but for JSON.",
		Long:         `A fast and 0-options way to format JSON files`,
		Args:         cobra.ArbitraryArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				args = []string{"**/*.json"}
			}
			return doRun(args)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&write, "write", "w", false, "write changes to the files")
	rootCmd.PersistentFlags().BoolVarP(&failfast, "failfast", "f", false, "exit on first error")
	rootCmd.PersistentFlags().StringVarP(&indent, "indent", "i", "  ", "indentation string")

	rootCmd.Version = version
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func doRun(globs []string) error {
	var rerr error

	files, err := findFiles(globs)
	if err != nil {
		return err
	}

	for _, file := range files {
		bts, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			return err
		}

		var out bytes.Buffer
		if err := json.Indent(&out, bytes.TrimSpace(bts), "", indent); err != nil {
			return fmt.Errorf("failed to format json file: %s: %v", file.Name(), err)
		}
		if _, err := out.Write([]byte{'\n'}); err != nil {
			return fmt.Errorf("failed to write: %v", err)
		}

		if bytes.Equal(bts, out.Bytes()) {
			continue
		}

		if write {
			if err := os.WriteFile(file.Name(), out.Bytes(), 0); err != nil {
				return fmt.Errorf("failed to write file: %s: %v", file.Name(), err)
			}
			continue
		}

		diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
			A:        difflib.SplitLines(string(bts)),
			B:        difflib.SplitLines(out.String()),
			FromFile: "Original",
			ToFile:   "Formatted",
			Context:  3,
		})
		if err != nil {
			return fmt.Errorf("failed to generate diff: %s: %v", file.Name(), err)
		}

		rerr = multierror.Append(rerr, fmt.Errorf("file %s differs:\n\n%s", file.Name(), diff))

		if failfast {
			break
		}
	}
	return rerr
}

func findFiles(globs []string) ([]*os.File, error) {
	if len(globs) == 1 && (globs)[0] == "-" {
		return []*os.File{os.Stdin}, nil
	}

	var files []*os.File
	var rerr error

	for _, glob := range globs {
		matches, err := fileglob.Glob(glob, fileglob.MaybeRootFS, fileglob.MatchDirectoryIncludesContents)
		if err != nil {
			return files, err
		}
		if len(matches) == 0 {
			multierror.Append(rerr, fmt.Errorf("no matches found: %s", glob))
		}

		for _, match := range matches {
			f, err := os.Open(match)
			if err != nil {
				multierror.Append(rerr, err)
			}
			files = append(files, f)
		}
	}

	return files, rerr
}
