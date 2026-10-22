package ui

import (
	"strings"
	"testing"
)

func TestTesterAccessibility(t *testing.T) {
	on := true
	name := "Ada"
	tt := NewTester(func(c *Context) {
		Column(c).Children(func() {
			Text(c, "Settings").Role(RoleHeading).Level(1)
			Button(c, "Save")
			Checkbox(c.Key("box"), &on, "Notify me")
			TextInput(c.Key("name"), &name).Label("Name")
			Box(c).Size(10, 10).Draw(func(*Painter, Rect) {}) // unnamed: unseen
		})
	}, 320, 240)
	els := tt.Accessibility()
	find := func(role, name string) AccessElement {
		t.Helper()
		for _, e := range els {
			if e.Role == role && e.Name == name {
				return e
			}
		}
		t.Fatalf("no %s %q in\n%s", role, name, tt.AccessibilityText())
		return AccessElement{}
	}
	if h := find("heading", "Settings"); h.Level != 1 {
		t.Errorf("heading %+v", h)
	}
	if b := find("button", "Save"); !b.Pressable || b.Bounds.W <= 0 {
		t.Errorf("button %+v", b)
	}
	if c := find("check box", "Notify me"); !c.Checked {
		t.Errorf("check box %+v", c)
	}
	if f := find("text field", "Name"); f.Value != "Ada" {
		t.Errorf("field %+v", f)
	}
	text := tt.AccessibilityText()
	if !strings.Contains(text, "button “Save”\n") || !strings.Contains(text, "check box “Notify me” (checked)\n") {
		t.Errorf("text:\n%s", text)
	}
}
