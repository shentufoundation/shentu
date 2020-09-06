package init

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// settings for CLI document generation
const (
	DocFlag      = "doc"
	DocFlagAbbr  = "d"
	DocFlagUsage = "generate cli document to <path>"
)

// GenDoc generates document for the current command-line interface to a designated folder.
func GenDoc(cmd *cobra.Command, outputDir string) {
	err := os.MkdirAll(outputDir, 0750)
	if err != nil {
		panic(err)
	}
	err = doc.GenMarkdownTree(cmd, (outputDir))
	if err != nil {
		panic(err)
	}
}
