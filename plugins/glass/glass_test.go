package glass

import (
	"image"
	"testing"

	"github.com/egoist/mygo/ui"
)

// pixelAt returns the color of img at (x, y).
func pixelAt(img *image.RGBA, x, y int) [3]int {
	p := img.RGBAAt(x, y)
	return [3]int{int(p.R), int(p.G), int(p.B)}
}

func TestGlassPaintsOverWhatIsBehind(t *testing.T) {
	g := Glass{}
	view := func(c *ui.Context) {
		ui.Box(c).Fill().Background(ui.RGB(20, 120, 220)).Children(func() {
			ui.Row(c).Size(200, 44).Radius(22).Material(g).Children(func() {
				ui.Text(c, "On glass")
			})
		})
	}
	tt := ui.NewTester(view, 300, 100)
	tt.SetScale(2)
	if !tt.HasText("On glass") {
		t.Error("the glass's children did not paint")
	}
	// The glass lightens the blue behind it, which shows through.
	light := pixelAt(tt.Image(), 300, 70)
	if light[0] < 100 || light[2] < light[0] {
		t.Errorf("light glass over blue is %v", light)
	}
	tt.SetDark(true)
	if dark := pixelAt(tt.Image(), 300, 70); dark[2] >= light[2] || dark[0] >= light[0] {
		t.Errorf("dark glass over blue is %v, not darker than light glass's %v", dark, light)
	}
	tt.SetDark(false)
	g.Style = Clear
	tt.Frame()
	if clear := pixelAt(tt.Image(), 300, 70); clear[0] >= light[0] {
		t.Errorf("clear glass over blue is %v, not clearer than regular glass's %v", clear, light)
	}
}

func TestInteractiveGlassGrowsWhilePressed(t *testing.T) {
	view := func(c *ui.Context) {
		ui.Box(c).Fill().Background(ui.RGB(0, 0, 0)).Children(func() {
			ui.Box(c).Absolute().Left(50).Top(50).Size(200, 100).Radius(20).Material(Glass{Interactive: true})
		})
	}
	tt := ui.NewTester(view, 300, 200)
	tt.SetScale(2)
	// It grows at once.
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	// The pixel left of the glass's left edge, in the middle.
	outside := func() int { return pixelAt(tt.Image(), 99, 200)[0] }
	before := outside()
	tt.Press(150, 100)
	tt.Frame()
	if grown := outside(); grown <= before+50 {
		t.Errorf("pressed, the pixel left of the glass is %d, as before (%d)", grown, before)
	}
	tt.Release(150, 100)
	tt.Frame()
	if after := outside(); after != before {
		t.Errorf("after the release, the pixel left of the glass is %d, not %d", after, before)
	}
}

func TestPaintInDrawing(t *testing.T) {
	view := func(c *ui.Context) {
		ui.Box(c).Fill().Background(ui.RGB(0, 0, 0)).Draw(func(p *ui.Painter, r ui.Rect) {
			Paint(p, ui.Rect{X: 20, Y: 20, W: 100, H: 60}, 20, Glass{})
		})
	}
	tt := ui.NewTester(view, 200, 120)
	// Light glass over black is mid gray.
	if g := pixelAt(tt.Image(), 70, 50); g[0] < 100 || g[0] > 180 {
		t.Errorf("light glass over black is %v", g)
	}
}

func TestString(t *testing.T) {
	if s := (Glass{Style: Clear, Tint: ui.RGB(1, 2, 3), Interactive: true}).String(); s != "glass clear tint(#010203ff) interactive" {
		t.Errorf("String: %q", s)
	}
}

// Panes up to 100 DIPs, as controls, keep the measured material; those
// from 200, as sidebars, are clearer; between, it eases from one to the
// other.
func TestPaneMaterial(t *testing.T) {
	near := func(a, b float32) bool { return a-b < 1e-5 && b-a < 1e-5 }
	for _, m := range []float32{20, 44, 72, 100} {
		for _, s := range []float32{1, 2} {
			bezel := min(36, m/2)
			blur := min(max(m*0.035, 1), 10) * s
			light := paneMaterial(m, s, false, false)
			if light.blur != blur || light.bezel != bezel*s || light.refraction != 1.6*bezel*s || light.rim != 0.35 ||
				light.low != 0.541 || light.high != 1 {
				t.Errorf("light %v DIPs at %vx: %+v", m, s, light)
			}
			dark := paneMaterial(m, s, false, true)
			if dark.blur != blur || dark.rim != 0.35 || dark.low != 0.15 || dark.high != 0.51 {
				t.Errorf("dark %v DIPs at %vx: %+v", m, s, dark)
			}
		}
	}
	for _, m := range []float32{200, 300, 900} {
		light, dark := paneMaterial(m, 2, false, false), paneMaterial(m, 2, false, true)
		if !near(light.blur, 10) || !near(light.bezel, 96) || !near(light.rim, 0.5) || !near(light.low, 0.35) || light.high != 1 {
			t.Errorf("light %v DIPs: %+v", m, light)
		}
		if !near(dark.blur, 10) || !near(dark.high, 0.6) || dark.low != 0.15 || !near(dark.rim, 0.5) {
			t.Errorf("dark %v DIPs: %+v", m, dark)
		}
	}
	// Between, every parameter lies between the two and moves one way.
	prev := paneMaterial(100, 1, false, false)
	for m := float32(110); m <= 200; m += 10 {
		mat := paneMaterial(m, 1, false, false)
		if mat.low > prev.low || mat.rim < prev.rim || mat.bezel < prev.bezel || mat.low < 0.35 || mat.rim > 0.5 {
			t.Errorf("%v DIPs: %+v after %+v", m, mat, prev)
		}
		prev = mat
	}
	if mid := paneMaterial(150, 1, false, false); !near(mid.low, (0.541+0.35)/2) || !near(mid.rim, 0.425) {
		t.Errorf("150 DIPs, half way: %+v", mid)
	}
	// Clear glass's tone is the same at any size.
	if small, large := paneMaterial(60, 1, true, false), paneMaterial(400, 1, true, false); small.low != large.low || small.high != large.high || small.blur != large.blur {
		t.Errorf("clear: %+v and %+v", small, large)
	}
}
