package commands

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var (
	catfileShort = "Provide content or type and size information for repository objects."
	catfileLong  = "The cat-file command allows you to inspect the contents of Delta objects stored in the repository. It can be used to retrieve the content of a blob, the structure of a tree, the details of a commit, or the information about a tag. The command is particularly useful for debugging or scripting purposes when you need to access raw object data."
	pretty       bool
)

var catfile = &cobra.Command{
	Use:   "cat-file",
	Short: catfileShort,
	Long:  catfileLong,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "usage: delta cat-file (-t [--allow-unknown-type] | -s [--allow-unknown-type] | -e | -p | <type> | --textconv | --filters) [--path=<path>] <object>\n")
			os.Exit(1)
		}

		var hash string = args[0]

		matched, err := regexp.Match("^[a-fA-F0-9]{40}$", []byte(hash))
		if err != nil {
			fmt.Fprintf(os.Stderr, "usage: delta cat-file (-t [--allow-unknown-type] | -s [--allow-unknown-type] | -e | -p | <type> | --textconv | --filters) [--path=<path>] <object>\n")
			return
		}
		if !matched {
			fmt.Fprintf(os.Stderr, "fatal: Not a valid object name: %s\n", hash)
			return
		}

		filepath := filepath.Join(".delta/objects", hash[:2], hash[2:])
		file, err := os.Open(filepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fatal: Unable to read object: %s\n", filepath)
			return
		}
		defer file.Close()

		if pretty {
			reader, err := zlib.NewReader(file)
			reader.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not read file %v", err)
			}
			containts, err := io.ReadAll(reader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not read from reader %v", err)
			}
			parts := strings.Split(string(containts), "\x00")
			fmt.Printf("%s\n", parts[1])
		}
	},
}

func init() {
	catfile.Flags().BoolVarP(&pretty, "pretty", "p", false, "pretty-print object's content")
	rootCmd.AddCommand(catfile)
}
