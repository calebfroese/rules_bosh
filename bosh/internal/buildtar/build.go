package buildtar

import (
	"archive/tar"
	"io"
	"os"
	"strings"
	"time"
)

type Builder struct {
	tw *tar.Writer
}

type AddOption func(*tar.Header)

func NewBuilder(w io.Writer) *Builder {
	tw := tar.NewWriter(w)

	return &Builder{
		tw: tw,
	}
}

func (b *Builder) Close() error {
	return b.tw.Close()
}

func (b *Builder) AddFile(path string, opts ...AddOption) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	hdr, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return err
	}
	for _, opt := range opts {
		opt(hdr)
	}
	if err := b.tw.WriteHeader(hdr); err != nil {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	if _, err := io.Copy(b.tw, f); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func Prefix(prefix string) AddOption {
	return func(hdr *tar.Header) {
		hdr.Name = prefix + hdr.Name
	}
}

func Rename(to string) AddOption {
	return func(hdr *tar.Header) {
		name := hdr.Name
		parts := strings.Split(name, "/")
		parts[len(parts)-1] = to
		hdr.Name = strings.Join(parts, "/")
	}
}

func Mode(mode os.FileMode) AddOption {
	return func(hdr *tar.Header) {
		hdr.Mode = int64(mode)
	}
}

func Hermetic() AddOption {
	return func(hdr *tar.Header) {
		hdr.Mode = int64(os.FileMode(0400))
		hdr.Uid = 0
		hdr.Gid = 0
		hdr.Uname = "root"
		hdr.Gname = "root"
		hdr.ModTime = time.Time{}
	}
}
