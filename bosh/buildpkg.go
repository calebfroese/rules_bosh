package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/xoebus/rules_bosh/bosh/internal/buildtar"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("buildpkg: ")
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	files := multiFlag{}
	flags := flag.NewFlagSet("buildpkg", flag.ExitOnError)
	output := flags.String("output", "", "path to place output")
	flags.Var(&files, "file", "repeated files to add to the package")
	if err := flags.Parse(args); err != nil {
		return err
	}

	out, err := os.Create(*output)
	if err != nil {
		return err
	}
	defer out.Close()

	gw := gzip.NewWriter(out)
	tb := buildtar.NewBuilder(gw)

	sort.Strings(files)
	for _, file := range files {
		if err := tb.AddFile(file, buildtar.Hermetic(), buildtar.Prefix("./"), buildtar.Mode(os.FileMode(0755))); err != nil {
			return err
		}
	}

	if err := tb.Close(); err != nil {
		return err
	}
	if err := gw.Close(); err != nil {
		return err
	}

	return nil
}

type multiFlag []string

func (f *multiFlag) Set(val string) error {
	*f = append(*f, val)
	return nil
}

func (m *multiFlag) String() string {
	if len(*m) == 0 {
		return ""
	}
	return fmt.Sprint(*m)
}
