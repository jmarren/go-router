package gorouter

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type RW struct {
	http.ResponseWriter
	*http.Request
}

func (rw *RW) Refresh() {
	rw.ResponseWriter.Header().Set("HX-Refresh", "true")
}

func (rw *RW) Retarget(target string) {
	rw.ResponseWriter.Header().Set("HX-Retarget", target)
}

func (rw *RW) Reselect(target string) {
	rw.ResponseWriter.Header().Set("HX-Reselect", target)
}

func (rw *RW) ReplaceUrl(url string) {
	rw.ResponseWriter.Header().Set("HX-Replace-Url", url)
}

func (rw *RW) Reswap(swapMethod string) {
	rw.ResponseWriter.Header().Set("HX-Reswap", swapMethod)
}

func (rw *RW) Redirect(location string) {
	rw.ResponseWriter.Header().Set("HX-Redirect", location)
}

func (rw *RW) IsHxRequest() bool {
	return rw.Request.Header.Get("HX-Request") == "true"
}

func (rw *RW) ContainsSubPath(subPath string) bool {
	url, _ := url.Parse(rw.CurrentUrl())
	return strings.Contains(url.Path, subPath)
}

func (rw *RW) CurrentUrl() string {
	return rw.Request.Header.Get("HX-Current-Url")
}

func (rw *RW) ExecutedScripts() []string {

	executedStr := rw.Request.Header.Get("HX-Executed")

	var executed []string

	json.Unmarshal([]byte(executedStr), &executed)

	return executed
}
