package ui

import (
	"strings"

	"github.com/egoist/mygo/internal/platform"
)

// AccessElement is an element as assistive technology sees it, for tests
// of what a screen reader reads (Tester.Accessibility).
type AccessElement struct {
	// Role is what the element is, as "button", "text", "list item" or
	// "heading".
	Role string
	// Name is what a screen reader calls it, Value its value, as a text
	// field's text, and Description what it reads after them, as a
	// tooltip.
	Name, Value, Description string
	// Depth is how deep it is in the tree, from 0 at the top.
	Depth int
	// Bounds is its visible box, in the window's DIPs.
	Bounds Rect
	// Index and Count tell which of how many items of its list it is,
	// from 1; 0 when not in one.
	Index, Count int
	// Level is the rank of a heading, or how deep a tree's item is.
	Level                                                     int
	Focused, Disabled, Checked, Selected, Expanded, Offscreen bool
	// Pressable tells that assistive technology can press it.
	Pressable bool
}

// accessRoleNames name the roles of AccessElement.
var accessRoleNames = map[platform.AccessRole]string{
	platform.RoleGroup: "group", platform.RoleText: "text", platform.RoleButton: "button",
	platform.RoleLink: "link", platform.RoleCheckBox: "check box", platform.RoleRadio: "radio button",
	platform.RoleSwitch: "switch", platform.RoleSlider: "slider", platform.RoleProgress: "progress",
	platform.RoleTextField: "text field", platform.RoleImage: "image", platform.RoleList: "list",
	platform.RoleScroll: "scroll area", platform.RoleDialog: "dialog", platform.RolePopup: "popup",
	platform.RoleTooltip: "tooltip", platform.RolePopUpButton: "pop-up button", platform.RoleTabList: "tab list",
	platform.RoleTab: "tab", platform.RoleSplitter: "splitter", platform.RoleStatus: "status",
	platform.RoleTable: "table", platform.RoleRow: "row", platform.RoleCell: "cell",
	platform.RoleColumnHeader: "column header", platform.RoleTree: "tree", platform.RoleTreeItem: "tree item",
	platform.RoleListItem: "list item", platform.RoleMenuButton: "menu button", platform.RoleToolbar: "toolbar",
	platform.RoleRadioGroup: "radio group", platform.RoleToggleButton: "toggle button", platform.RoleComboBox: "combo box",
	platform.RoleDisclosure: "disclosure", platform.RoleMeter: "meter", platform.RoleStepper: "stepper",
	platform.RoleColorWell: "color well", platform.RoleAlertDialog: "alert dialog", platform.RoleMenu: "menu",
	platform.RoleMenuBar: "menu bar", platform.RoleMenuItem: "menu item", platform.RoleMenuItemCheckBox: "menu item check box",
	platform.RoleMenuItemRadio: "menu item radio", platform.RoleHeading: "heading",
}

// Accessibility returns the elements of the last frame as assistive
// technology sees them, in tree order, a parent before its children. The
// first call turns on describing frames for it, as a screen reader
// starting does, and draws a frame.
func (t *Tester) Accessibility() []AccessElement {
	if t.h.access == nil {
		t.send(platform.SurfaceEvent{Kind: platform.AccessibilityOn})
	}
	t.Frame()
	tree := t.h.access
	if tree == nil {
		return nil
	}
	out := make([]AccessElement, len(tree.Nodes))
	for i, n := range tree.Nodes {
		e := AccessElement{
			Role: accessRoleNames[n.Role], Name: n.Label, Value: n.Value, Description: n.Description,
			Bounds: Rect{X: float32(n.Bounds.X), Y: float32(n.Bounds.Y), W: float32(n.Bounds.W), H: float32(n.Bounds.H)},
			Index:  n.PosInSet, Count: n.SetSize, Level: n.Level,
			Focused:   n.ID == tree.Focus && n.ID != 0,
			Disabled:  n.States&platform.AccessDisabled != 0,
			Checked:   n.States&platform.AccessChecked != 0,
			Selected:  n.States&platform.AccessSelected != 0,
			Expanded:  n.States&platform.AccessExpanded != 0,
			Offscreen: n.States&platform.AccessOffscreen != 0,
			Pressable: n.Actions&platform.ActionPress != 0,
		}
		if n.Parent >= 0 && n.Parent < i {
			e.Depth = out[n.Parent].Depth + 1
		}
		out[i] = e
	}
	return out
}

// AccessibilityText is Accessibility as text, an element a line indented
// by its depth, as "button “Save”", for a test's failure message or a
// snapshot of what a screen reader reads.
func (t *Tester) AccessibilityText() string {
	var b strings.Builder
	for _, e := range t.Accessibility() {
		b.WriteString(strings.Repeat("  ", e.Depth))
		b.WriteString(e.Role)
		if e.Name != "" {
			b.WriteString(" “" + e.Name + "”")
		}
		if e.Value != "" {
			b.WriteString(" = " + e.Value)
		}
		for _, s := range []struct {
			on   bool
			name string
		}{{e.Focused, "focused"}, {e.Disabled, "disabled"}, {e.Checked, "checked"}, {e.Selected, "selected"}, {e.Expanded, "expanded"}, {e.Offscreen, "offscreen"}} {
			if s.on {
				b.WriteString(" (" + s.name + ")")
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}
