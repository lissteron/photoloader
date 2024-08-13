package ui

import (
	"embed"
	"mime"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/cespare/xxhash/v2"
)

//go:embed dist/*
var frontend embed.FS

const (
	index = "dist/index.html"

	ifNoneMatchHeader = "If-None-Match"
	etagHeader        = "ETag"
)

func NewFS() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fileData, err := frontend.ReadFile("dist/" + strings.TrimLeft(r.URL.Path, "/"))
		if err != nil {
			if r.URL.Path != "/" {
				w.WriteHeader(http.StatusNotFound)

				return
			}

			fileData, err = frontend.ReadFile(index)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)

				return
			}
		}

		etag := makeETAG(fileData)
		if etag == r.Header.Get(ifNoneMatchHeader) {
			w.WriteHeader(http.StatusNotModified)

			return
		}

		w.Header().Add("Content-Type", getContentType(r.URL.Path))
		w.Header().Add("Cache-Control", "public")
		w.Header().Add(etagHeader, etag)

		_, _ = w.Write(fileData)
	})
}

func getContentType(filename string) string {
	if strings.HasSuffix(filename, ".ico") {
		return "image/x-icon"
	}

	return mime.TypeByExtension(path.Ext(filename))
}

func makeETAG(b []byte) string {
	const base = 32

	return strconv.FormatInt(int64(len(b)), base) + "-" + strconv.FormatUint(xxhash.Sum64(b), base)
}
