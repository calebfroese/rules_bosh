package main

import (
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/xoebus/rules_bosh/internal/buildtar"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("buildjob: ")
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	templates := multiFlag{}
	flags := flag.NewFlagSet("buildjob", flag.ExitOnError)
	manifest := flags.String("manifest", "", "path to the job spec file")
	monit := flags.String("monit", "", "path to the job monit file")
	flags.Var(&templates, "template", "repeated template files for the job")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *manifest == "" {
		return errors.New("-manifest must be specified")
	}
	if *monit == "" {
		return errors.New("-monit must be specified")
	}

	gw := gzip.NewWriter(os.Stdout)
	tb := buildtar.NewBuilder(gw)
	if err := tb.AddFile(*manifest, buildtar.Hermetic(), buildtar.Prefix("./"), buildtar.Rename("job.MF"), buildtar.Mode(os.FileMode(0644))); err != nil {
		return err
	}
	if err := tb.AddFile(*monit, buildtar.Hermetic(), buildtar.Prefix("./"), buildtar.Mode(os.FileMode(0644))); err != nil {
		return err
	}
	sort.Strings(templates)
	for _, template := range templates {
		if err := tb.AddFile(template, buildtar.Hermetic(), buildtar.Prefix("./templates/"), buildtar.Mode(os.FileMode(0644))); err != nil {
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
