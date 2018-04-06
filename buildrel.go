package main

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/xoebus/rules_bosh/internal/buildtar"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("buildrel: ")
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	jobs := multiFlag{}
	packages := multiFlag{}
	flags := flag.NewFlagSet("buildrel", flag.ExitOnError)
	name := flags.String("name", "", "name of the release")
	stemcellDistro := flags.String("stemcellDistro", "", "distro of the stemcell")
	stemcellVersion := flags.String("stemcellVersion", "", "version of the stemcell")
	flags.Var(&jobs, "job", "repeated jobs for the release")
	flags.Var(&packages, "package", "repeated packages for the release")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *name == "" {
		return errors.New("-name must be specified")
	}
	if *stemcellDistro == "" {
		return errors.New("-stemcellDistro must be specified")
	}
	if *stemcellVersion == "" {
		return errors.New("-stemcellVersion must be specified")
	}

	gw, err := gzip.NewWriterLevel(os.Stdout, gzip.BestSpeed)
	if err != nil {
		return err
	}
	tb := buildtar.NewBuilder(gw)

	manifest := Manifest{
		Name:               *name,
		Version:            "0.0.0+dev.1",
		CommitHash:         "0000000",
		UncommittedChanges: true,
	}

	sort.Strings(jobs)
	for _, job := range jobs {
		if err := tb.AddFile(job, buildtar.Hermetic(), buildtar.Prefix("./jobs/"), buildtar.Mode(os.FileMode(0644))); err != nil {
			return err
		}

		jobName := strings.TrimSuffix(filepath.Base(job), filepath.Ext(job))
		sha, err := shaFile(job)
		if err != nil {
			return err
		}
		manifest.Jobs = append(manifest.Jobs, Job{
			Name:        jobName,
			Fingerprint: sha,
			Sha1:        fmt.Sprintf("sha256:%s", sha),
		})
	}
	sort.Strings(packages)
	for _, pkg := range packages {
		if err := tb.AddFile(pkg, buildtar.Hermetic(), buildtar.Prefix("./compiled_packages/"), buildtar.Mode(os.FileMode(0644))); err != nil {
			return err
		}
		pkgName := strings.TrimSuffix(filepath.Base(pkg), filepath.Ext(pkg))
		sha, err := shaFile(pkg)
		if err != nil {
			return err
		}
		manifest.Packages = append(manifest.Packages, Package{
			Name:        pkgName,
			Fingerprint: sha,
			Sha1:        fmt.Sprintf("sha256:%s", sha),
			Stemcell:    fmt.Sprintf("%s/%s", *stemcellDistro, *stemcellVersion),
		})
	}

	f, err := ioutil.TempFile("", "releasemanifest")
	if err != nil {
		return err
	}
	defer f.Close()
	defer os.Remove(f.Name())
	if err := json.NewEncoder(f).Encode(manifest); err != nil {
		return err
	}
	if err := tb.AddFile(f.Name(), buildtar.Hermetic(), buildtar.Prefix("./"), buildtar.Rename("release.MF"), buildtar.Mode(os.FileMode(0644))); err != nil {
		return err
	}

	if err := tb.Close(); err != nil {
		return err
	}
	if err := gw.Close(); err != nil {
		return err
	}

	return nil
}

func shaFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

type Manifest struct {
	Name               string    `json:"name"`
	Version            string    `json:"version"`
	UncommittedChanges bool      `json:"uncommitted_changes"`
	CommitHash         string    `json:"commit_hash"`
	Jobs               []Job     `json:"jobs"`
	Packages           []Package `json:"compiled_packages"`
}

type Job struct {
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
	Sha1        string `json:"sha1"`
}

type Package struct {
	Name         string     `json:"name"`
	Fingerprint  string     `json:"fingerprint"`
	Sha1         string     `json:"sha1"`
	Stemcell     string     `json:"stemcell"`
	Dependencies []struct{} `json:"dependencies"`
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
