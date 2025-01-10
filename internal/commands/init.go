package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initRepository = &cobra.Command{
	Use:   "init",
	Short: "Initialize a repository",
	Long:  "Initialize a delta repository on your local machine",
	Run: func(cmd *cobra.Command, args []string) {
		for _, dir := range []string{".delta", ".delta/objects", ".delta/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".delta/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")
	},
}

func init() {
	rootCmd.AddCommand(initRepository)
}