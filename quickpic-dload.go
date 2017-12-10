package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

func main() {
	var (
		ulist     = flag.String("list", "cloud_gallery_urls.txt", "load URLs from the list in the file")
		outdir    = flag.String("outdir", "quickpic-backup", "save the image files to this directory (it will be created if not exists)")
		overwrite = flag.Bool("rewrite", false, "overwrite old files with same names on each new run of the utility")
		verbose   = flag.Bool("v", false, "verbose output")
	)
	flag.Parse()
	fd, err := os.Open(*ulist)
	defer fd.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open file %s with URLs list: %s\n", *ulist, err)
		os.Exit(1)
	}
	os.MkdirAll(*outdir, 0755)
	var rawURL string
	for {
		_, err := fmt.Fscanln(fd, &rawURL)
		if err != nil {
			break
		}
		fname, found := check(rawURL, *outdir, *overwrite, *verbose)
		if !*overwrite && found {
			continue
		}
		resp := load(rawURL)
		save(fname, resp, *verbose)
		resp.Close()
	}
}

func check(rawURL, outdir string, overwrite, verbose bool) (string, bool) {
	var name string
	if parsedUrl, err := url.Parse(rawURL); err == nil {
		name = parsedUrl.Query()["s"][0]
	}
	if name == "" {
		name = "xxx" // TODO fix it
	}
	fname := path.Join(outdir, name+".jpeg")
	if overwrite {
		return fname, false
	}
	if fd, err := os.Open(fname); err == nil {
		fd.Close()
		if verbose {
			fmt.Fprintf(os.Stdout, "%s already found — skip download\n", fname)
		}
		return fname, true
	}
	return fname, false
}

func load(rawURL string) io.ReadCloser {
	resp, err := http.Get(rawURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open %s: %s\n", rawURL, err)
		return nil
	}
	return resp.Body
}

func save(dname string, src io.Reader, verbose bool) {
	fd, err := os.Create(dname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't save the file into %s: %s\n", dname, err)
		return
	}
	_, err = io.Copy(fd, src)
	defer fd.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't copy to %s: %s\n", dname, err)
	}
	if verbose {
		fmt.Fprintf(os.Stdout, "%s saved\n", dname)
	}
}
