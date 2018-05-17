package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// StreamFileSystem asdf
type StreamFileSystem struct {
	root   string
	prefix string
	fs     http.FileSystem
}

// NewStreamFileSystem returns a new filesysteam
func NewStreamFileSystem(root, prefix string) *StreamFileSystem {
	return &StreamFileSystem{
		root:   root,
		prefix: prefix,
	}
}

func (s *StreamFileSystem) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.DirList(w, r)
}

func (s *StreamFileSystem) Stream(w http.ResponseWriter, r *http.Request) {
	p := path.Clean(r.URL.Path)
	if strings.HasPrefix(p, "/stream") {
		p = p[len("/stream"):]
	}

	p = filepath.Join(s.root, p, "/")

	f, err := os.Open(p)
	if err != nil {
		http.Error(w, "file not found or so "+err.Error(), http.StatusNotFound)
		return
	}

	d, err := f.Stat()
	if err != nil {
		http.Error(w, "f.Stat() error", http.StatusInternalServerError)
		return
	}

	if d.IsDir() {
		dirList(w, f)
		return
	}

	http.ServeFile(w, r, d.Name())
}

// DirList returns all files in the path
func (s *StreamFileSystem) DirList(w http.ResponseWriter, r *http.Request) {
	// clean path an replace prefix
	p := path.Clean(r.URL.Path)
	if strings.HasPrefix(p, s.prefix) {
		p = p[len(s.prefix):]
	}

	p = filepath.Join(s.root, p, "/")

	f, err := os.Open(p)
	if err != nil {
		http.Error(w, "file not found or so "+err.Error(), http.StatusNotFound)
		return
	}

	d, err := f.Stat()
	if err != nil {
		http.Error(w, "f.Stat() error", http.StatusInternalServerError)
		return
	}

	if d.IsDir() {
		dirList(w, f)
		return
	}

	http.ServeFile(w, r, d.Name())
}

func dirList(w http.ResponseWriter, f http.File) {
	dirs, err := f.Readdir(-1)
	if err != nil {
		// TODO: log err.Error() to the Server.ErrorLog, once it's possible
		// for a handler to get at its Server via the ResponseWriter. See
		// Issue 12438.
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, htmlHeader)
	fmt.Fprintf(w, "<ul>\n")
	fmt.Fprintf(w, "<li><a href=\"%s\">%s</a></li>", ".", ".")
	fmt.Fprintf(w, "<li><a href=\"%s\">%s</a></li>", "..", "..")

	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() {
			name += "/"
		}
		// name may contain '?' or '#', which must be escaped to remain
		// part of the URL path, and not indicate the start of a query
		// string or fragment.
		url := url.URL{Path: name}
		fmt.Fprintf(w, `<li><a href="/stream/%s">VLC Link</a> - <a href="%s">%s</a></li>`, url.String(), url.String(), name)
	}
	fmt.Fprintf(w, "</ul>\n")
	fmt.Fprintf(w, htmlFooter)
}
