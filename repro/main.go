// Zoomdemo shows a diff view in a window zoomed as an app zooms by
// scaling its root: laid out 1/z the window's size, then scaled by z. The
// zoom comes from $ZOOM on the desktop, or the page's URL (?zoom=2) on the web.
package main

import (
	"fmt"
	"log"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

// mygoVersion is the version of mygo built in, its replacement's if replaced.
func mygoVersion() string {
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, m := range bi.Deps {
			if m.Path == "github.com/egoist/mygo" {
				if m.Replace != nil {
					return m.Replace.Path + " " + m.Replace.Version
				}
				return m.Path + " " + m.Version
			}
		}
	}
	return "mygo"
}

type demo struct {
	z    float32
	diff *editor.DiffState
}

func (s *demo) view(c *ui.Context) {
	width, height := c.Size()
	z := s.z
	ui.Column(c).Absolute().Top(0).Left(0).Width(width/z).Height(height/z).TransformOrigin(0, 0).Scale(z, z).
		Background(c.Theme().Background).Children(func() {
		ui.Row(c).Height(120).Padding(0, 16, 0, 16).AlignItems(ui.Center).Children(func() {
			ui.Text(c, fmt.Sprintf("%s · zoom %v×", mygoVersion(), z)).Bold()
		})
		editor.DiffView(c, s.diff).Grow(1).FillWidth()
	})
}

func main() {
	z := float32(1)
	if f, err := strconv.ParseFloat(zoomSetting(), 32); err == nil && f > 0 {
		z = float32(f)
	}
	var b strings.Builder
	for i := range 400 {
		fmt.Fprintf(&b, "line %d of a long file\n", i)
	}
	s := &demo{z: z, diff: &editor.DiffState{Old: editor.NewBuffer(""), New: editor.NewBuffer(b.String()), Layout: editor.DiffUnified}}
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{Title: "Zoom demo", Content: ui.View(s.view)})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
