package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	lru "github.com/hashicorp/golang-lru"
)

type StaticResourceHandlerOption func(h *StaticHandler)

type StaticHandler struct {
	dir                     string
	pathPrefix              string
	extensionContentTypeMap map[string]string
	cache                   *lru.Cache
	maxFileSize             int
}

type fileCacheItem struct {
	fileName    string
	fileSize    int
	contentType string
	data        []byte
}

/** static route handler **/
func (h *StaticHandler) Handle(c *Context) {
	reqPath := strings.TrimPrefix(c.R.URL.Path, h.pathPrefix)
	reqFilePath := filepath.Join(h.dir, reqPath)

	/** try to get data from cache **/
	if item, ok := h.readFromCache(reqPath); ok {
		fmt.Printf("read data from cache...")
		h.writeItemAsResponse(item, c.W)
		return
	}

	/** read local file **/
	f, err := os.Open(reqFilePath)
	if err != nil {
		c.W.WriteHeader(http.StatusInternalServerError)
		return
	}

	/** check extension **/
	ext := strings.TrimPrefix(filepath.Ext(f.Name()), ".")
	t, ok := h.extensionContentTypeMap[ext]
	if !ok {
		c.W.WriteHeader(http.StatusBadRequest)
		return
	}

	/** load file **/
	data, err := ioutil.ReadAll(f)
	if err != nil {
		c.W.WriteHeader(http.StatusInternalServerError)
		return
	}

	item := &fileCacheItem{
		fileSize:    len(data),
		data:        data,
		contentType: t,
		fileName:    reqPath,
	}

	/** cache file **/
	h.cacheFile(item)
	h.writeItemAsResponse(item, c.W)
}

/** read data from cache **/
func (h *StaticHandler) readFromCache(fileName string) (*fileCacheItem, bool) {
	if h.cache != nil {
		if item, ok := h.cache.Get(fileName); ok {
			return item.(*fileCacheItem), true
		}
	}
	return nil, false
}

/** response with file **/
func (h *StaticHandler) writeItemAsResponse(item *fileCacheItem, writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", item.contentType)
	writer.Header().Set("Content-Length", fmt.Sprintf("%d", item.fileSize))
	writer.Write(item.data)
}

/** save file data in cache **/
func (h *StaticHandler) cacheFile(item *fileCacheItem) {
	if h.cache != nil && item.fileSize < h.maxFileSize {
		h.cache.Add(item.fileName, item)
	}
}

/** set cache **/
func WithFileCache(maxFileSizeThreshold int, maxCacheFileCnt int) StaticResourceHandlerOption {
	return func(h *StaticHandler) {
		c, err := lru.New(maxCacheFileCnt)
		if err != nil {
			fmt.Printf("could not create LRU, we won't cache static file")
		}
		h.maxFileSize = maxFileSizeThreshold
		h.cache = c
	}
}

/** set extension **/
func WithMoreExtension(extMap map[string]string) StaticResourceHandlerOption {
	return func(h *StaticHandler) {
		for ext, contentType := range extMap {
			h.extensionContentTypeMap[ext] = contentType
		}
	}
}

/** create a new staticHandler **/
func NewStaticHandler(dir string, pathPrefix string, options ...StaticResourceHandlerOption) *StaticHandler {
	h := &StaticHandler{
		dir:        dir,
		pathPrefix: pathPrefix,
		extensionContentTypeMap: map[string]string{
			"jpeg": "image/jpeg",
			"jpe":  "image/jpeg",
			"jpg":  "image/jpeg",
			"png":  "image/png",
			"pdf":  "image/pdf",
		},
	}
	/** set option **/
	for _, o := range options {
		o(h)
	}
	return h
}
