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
	target string
}

// Request Headers
func (rw *RW) Target() string {
	return rw.Request.Header.Get("HX-Target")
}

func (rw *RW) TriggerName() string {
	return rw.Request.Header.Get("HX-Trigger-Name")
}

func (rw *RW) TriggerId() string {
	return rw.Request.Header.Get("HX-Trigger")
}

func (rw *RW) IsHistoryRestoreRequest() bool {
	return rw.Request.Header.Get("HX-History-Restore-Request") == "true"
}

func (rw *RW) Prompt() string {
	return rw.Request.Header.Get("HX-Prompt")
}

func (rw *RW) Boosted() bool {
	return rw.Request.Header.Get("HX-Boosted") == "true"
}

// Response Headers
func (rw *RW) Refresh() {
	rw.ResponseWriter.Header().Set("HX-Refresh", "true")
}

func (rw *RW) Retarget(target string) {
	rw.target = target
	// rw.ResponseWriter.Header().Set("HX-Retarget", target)
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

func (rw *RW) CurrentUrl() string {
	return rw.Request.Header.Get("HX-Current-Url")
}

// Custom Headers
func (rw *RW) ExecutedScripts() []string {

	executedStr := rw.Request.Header.Get("HX-Executed")

	var executed []string

	json.Unmarshal([]byte(executedStr), &executed)

	return executed
}

// Other
func (rw *RW) PathHasPrefix(subPath string) bool {
	url, _ := url.Parse(rw.CurrentUrl())
	path := url.Path
	return strings.HasPrefix(path, subPath)
}
