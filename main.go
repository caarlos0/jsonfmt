package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"strconv"
	"strings"

	"github.com/editorconfig/editorconfig-core-go/v2"
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
			return doRun(args, cmd.PersistentFlags().Changed("indent"))
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

func doRun(globs []string, indentChanged bool) error {
	var rerr error

	files, err := findFiles(globs)
	if err != nil {
		return err
	}

	for _, file := range files {
		bts, err := file.Read()
		if err != nil {
			return err
		}

		var out bytes.Buffer
		if err := json.Indent(&out, bytes.TrimSpace(bts), "", computeIndent(file.Name(), indentChanged)); err != nil {
			return fmt.Errorf("failed to format json file: %s: %v", file.Name(), err)
		}
		if _, err := out.Write([]byte{'\n'}); err != nil {
			return fmt.Errorf("failed to write: %v", err)
		}

		if bytes.Equal(bts, out.Bytes()) {
			continue
		}

		if write {
			if err := file.Write(out.Bytes()); err != nil {
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

func computeIndent(filename string, indentChanged bool) string {
	if indentChanged || filename == "-" {
		return indent
	}

	ec, err := editorconfig.GetDefinitionForFilename(filename)

	if err != nil {
		return indent
	}

	indentStyle := ec.IndentStyle
	indentSize := ec.IndentSize
	tabWidth := ec.TabWidth

	if indentStyle == "tab" || indentSize == "tab" {
		return "\t"
	}

	if indentSize != "" {
		if n, err := strconv.Atoi(indentSize); err == nil && n > 0 {
			return strings.Repeat(" ", n)
		}
	}

	if tabWidth > 0 {
		return strings.Repeat(" ", tabWidth)
	}

	return indent
}

func findFiles(globs []string) ([]file, error) {
	if len(globs) == 1 && (globs)[0] == "-" {
		return []file{stdInOut{}}, nil
	}

	var files []file
	var rerr error

	for _, glob := range globs {
		matches, err := fileglob.Glob(glob, fileglob.MaybeRootFS, fileglob.MatchDirectoryIncludesContents)
		if err != nil {
			return files, err
		}
		if len(matches) == 0 {
			rerr = multierror.Append(rerr, fmt.Errorf("no matches found: %s", glob))
		}

		for _, match := range matches {
			files = append(files, realFile{match})
		}
	}

	return files, rerr
}

type file interface {
	Name() string
	Read() ([]byte, error)
	Write(b []byte) error
}

type stdInOut struct{}

func (f stdInOut) Name() string          { return "-" }
func (f stdInOut) Read() ([]byte, error) { return io.ReadAll(os.Stdin) }
func (f stdInOut) Write(b []byte) error {
	_, err := os.Stdout.Write(b)
	return err
}

type realFile struct {
	path string
}

func (f realFile) Name() string          { return f.path }
func (f realFile) Read() ([]byte, error) { return os.ReadFile(f.path) }
func (f realFile) Write(b []byte) error {
	stat, err := os.Stat(f.Name())
	if err != nil {
		return err
	}
	return os.WriteFile(f.Name(), b, stat.Mode())
}
