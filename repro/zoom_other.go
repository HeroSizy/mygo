//go:build !js

package main

import "os"

// zoomSetting is $ZOOM.
func zoomSetting() string { return os.Getenv("ZOOM") }
