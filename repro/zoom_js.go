package main

import (
	"strings"
	"syscall/js"
)

// zoomSetting is the zoom in the page's URL, as ?zoom=2.
func zoomSetting() string {
	return strings.TrimPrefix(js.Global().Get("location").Get("search").String(), "?zoom=")
}
