package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io/ioutil"
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

var tarOpts = []buildtar.AddOption{buildtar.Hermetic(), buildtar.Prefix("./"), buildtar.Mode(os.FileMode(0755))}

func run(args []string) error {
	files := multiFlag{}
	flags := flag.NewFlagSet("buildpkg", flag.ExitOnError)
	output := flags.String("output", "", "path to place output")
	uncompiled := flags.Bool("uncompiled", false, "make an uncompiled package")
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

	if *uncompiled {
		f, err := ioutil.TempFile("", "packaging")
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := f.Write([]byte(packagingScript)); err != nil {
			return err
		}
		if err := tb.AddFile(f.Name(), append(tarOpts, buildtar.Rename("packaging"))...); err != nil {
			return err
		}
	}

	sort.Strings(files)
	for _, file := range files {
		if err := tb.AddFile(file, tarOpts...); err != nil {
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

var packagingScript = `#!/bin/bash

set -e
set -u

cp -r ${BOSH_COMPILE_TARGET}/* ${BOSH_INSTALL_TARGET}
rm ${BOSH_INSTALL_TARGET}/packaging
`
