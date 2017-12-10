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
		ulist  = flag.String("list", "cloud_gallery_urls.txt", "load URLs from the list in the file")
		outdir = flag.String("outdir", "quickpic-backup", "save the image files to this directory (it will be created if not exists)")
	)
	flag.Parse()
	fd, err := os.Open(*ulist)
	defer fd.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open file %s with URLs list: %s\n", *ulist, err)
		os.Exit(1)
	}
	os.MkdirAll(*outdir, 0755)
	var rawUrl string
	for {
		_, err := fmt.Fscanln(fd, &rawUrl)
		if err != nil {
			break
		}
		name, resp := load(rawUrl)
		save(path.Join(*outdir, name+".jpeg"), resp)
		resp.Close()
	}
}

func load(rawUrl string) (string, io.ReadCloser) {
	resp, err := http.Get(rawUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open %s: %s\n", rawUrl, err)
		return "", nil
	}
	var name string
	if parsedUrl, err := url.Parse(rawUrl); err == nil {
		name = parsedUrl.Query()["s"][0]
	}
	if name == "" {
		name = "xxx"
	}
	return name, resp.Body
}

func save(dname string, src io.Reader) {
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
}
