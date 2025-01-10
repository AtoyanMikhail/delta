package commands

import (
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	hashObjectShort = "Compute object ID and optionally creates a blob from a file."
	hashObjectLong  = `Computes the object ID value for an object with specified type with the contents of the named file (which can be outside of the
work tree), and optionally writes the resulting object into the object database. Reports its object ID to its standard output.
When <type> is not specified, it defaults to "blob".`
	hashObjectUsage = "delta cat-file (-t [--allow-unknown-type] | -s [--allow-unknown-type] | -e | -p | <type> | --textconv | --filters) [--path=<path>] <object>"
)

var hashObject = &cobra.Command{
	Use:   "hash-object",
	Short: hashObjectShort,
	Long:  hashObjectLong,
	Run: func(cmd *cobra.Command, args []string) {
		// args[0] is a file name that was passed as an argument to the command
		if len(args) < 1 {
			fmt.Fprint(os.Stderr, hashObjectUsage)
			return
		}

		file, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "fatal: Cannot open '%s' : %v", args[0], err)
			return
		}
		defer file.Close()
		stats, _ := os.Stat(args[0])

		content, _ := os.ReadFile(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v", err)
			return
		}
		contentAndHeader := fmt.Sprintf("blob %d\x00%s", stats.Size(), content)
		sha := sha1.Sum([]byte(contentAndHeader))
		hash := fmt.Sprintf("%x", sha)

		blobName := []rune(hash)
		blobPath := filepath.Join(".delta/objects/", string(blobName[:2]), string(blobName[2:]))

		os.Mkdir(filepath.Dir(blobPath), os.ModePerm)

		f, _ := os.Create(blobPath)
		defer f.Close()

		z := zlib.NewWriter(f)
		defer z.Close()
		if _, err := z.Write([]byte(contentAndHeader)); err != nil {
			fmt.Printf("Error writing to file: %v", err)
            return
		}

		fmt.Print(hash)
	},
}

func init() {
	rootCmd.AddCommand(hashObject)
}
