# Windows with native UI

A window with `Content` takes the [window options](../windows.md#options)
of any window, such as its size, `StateKey`, `Frameless` and `Parent`, as
well as menus, dialogs and the other [desktop APIs](../native.md). It has no
page: `Page()` is nil, and it ignores `URL` and `WindowOptions.Page`. Its Go
code needs no bindings: the view calls it directly. `CapturePage` returns a
PNG of what it shows.

```go
mygo.NewWindow(mygo.WindowOptions{
	Title:    "Notes",
	Width:    900,
	Height:   600,
	StateKey: "main",
	Content:  ui.View(app.view),
})
```

## Title bars

With `TitleBarStyle: mygo.TitleBarHidden`, the view draws the title bar
under the window controls, as a page does with the `--mygo-titlebar-*` CSS
variables: `c.TitleBar()` returns the room the controls take, zero in full
screen, and `DragWindow` makes elements drag the window, which
double-clicking them zooms or minimizes as a title bar would:

```go
bar := c.TitleBar()
ui.Row(c).Height(max(bar.Height, 32)).Padding(0, bar.Right+12, 0, bar.Left+12).DragWindow().Children(func() {
	ui.Text(c, "Inbox").Bold()
})
```

In a `Frameless` window, `DragWindow` makes an element move the window, and
a double click on it maximizes the window.

## Zoom

`c.SetZoom(z)` zooms the window's content, as a browser zooms a page: the
view is laid out in a window `1/z` its size and shown `z` times larger. It
is not a picture of the window scaled up: frames draw at the screen's
resolution times `z`, so text, icons, vectors, images and their edges
stay as sharp as at that size natively, and the window's frames keep
their size in pixels. A zoom of 2 on a screen of scale 2 draws the frame
a screen of scale 4 would, glyphs and all; unzoomed, nothing changes, and
text keeps its subpixel antialiasing. The zoom stays until set again,
between 0.25 and 5; `c.Zoom()` returns it, 1 for none.

```go
c.SetZoom(app.zoom) // 1.25 for 125%, every frame or once
width, height := c.Size() // the content's: the window's divided by the zoom
```

What the view sees is zoomed: `c.Size()`, `c.TitleBar()`, `Bounds`
and every length are in the content's DIPs, `z` of them to the window's.
A frame that sets another zoom builds again at the new size before it
shows. What leaves the content goes to the window in its DIPs: the
pointer hits what the zoom shows under it, input methods compose at the
caret, assistive technology outlines elements where they show, context
menus and drags start under the pointer and web views take the box
the content gives them. A touchpad scrolls the content as far as the
fingers move; a wheel's lines scroll as far as unzoomed. The window
controls keep their size: `c.TitleBar()` is the room they take, zoomed,
so that a title bar laid out clear of them stays clear.

A transform's `Scale`, by contrast, scales what is drawn at the content's
own size: its text is drawn as bitmaps of their untransformed size,
sampled bilinearly ([Transforms](transforms.md#rendering)), right for
motion, blurry for zooming a window.

## Vibrancy

A window's `Vibrancy` shows wherever its native UI draws no background, on
macOS and, in a window created with one, on Windows 11 22H2 and later.
`c.Vibrancy()` tells whether it shows: elsewhere, what the view leaves
transparent shows no material, so draw a background there instead. The
root draws the theme's by default: make it transparent and give
backgrounds to the parts that need one, as a sidebar beside opaque content
does:

```go
if c.Vibrancy() {
	c.Root().Background(ui.Transparent)
}
ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
	sidebar := ui.Column(c).Width(240)
	if !c.Vibrancy() {
		sidebar.Background(c.Theme().Surface)
	}
	sidebar.Children(func() { app.sidebar(c) }) // over the material
	ui.Column(c).Grow(1).Background(c.Theme().Background).Children(func() { app.content(c) })
})
```

Text over an opaque background keeps the system's subpixel antialiasing
(ClearType) under a transparent root; over the material, it is
antialiased in grayscale. `examples/vibrancy` picks the window's material
in such a sidebar.

## Closing

A view can close its own window, from a button's click or its callback.
The window goes at once, by `Destroy` or by a `Close` no `OnClose`
listener prevents, and the rest of the view builds as it would, with a
`Context` that still has the theme and the window's size. The frame ends
there: nothing more is laid out or painted, and no other callback runs.
Open the window that takes over first, since an app [quits](../app.md)
when its last window closes:

```go
if ui.Button(c, "Done").Clicked() {
	app.showMain()
	win.Close()
}
```

## No webview

An app whose windows all show native UI needs no webview: on Linux it needs
GTK 3 alone, not WebKitGTK, and on Windows no WebView2 Runtime.
