// Regenerates docs/*.md from the cobra command tree.
//
// Usage: cd tools/gendocs && go run .
//
// Lives in its own go module so the root module's go.mod isn't polluted with
// cobra/doc's transitive deps (go-md2man, blackfriday, yaml).
package main

import (
	"log"
	"os"

	"github.com/radiusmethod/awsd/src/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	out := "../../docs"
	if len(os.Args) > 1 {
		out = os.Args[1]
	}
	if err := doc.GenMarkdownTree(cmd.RootCmd(), out); err != nil {
		log.Fatal(err)
	}
}
