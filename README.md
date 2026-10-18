# Screenshots for egoist/mygo#488

Before (mygo v0.5.5) and after (v0.5.5 with #488) of `repro/`, a standalone
mygo app: an `editor.DiffView` of 400 lines in a root laid out at 1/z of
the window and scaled by z, as an app zooms. Images and scripts only; not
code to merge.

- `repro/`: the app. `go.mod` builds the fix (a replace with HeroSizy/mygo
  v0.5.5-herosizy.1, v0.5.5 plus the patch); `go mod edit -dropreplace
  github.com/egoist/mygo` builds stock v0.5.5. The zoom is `$ZOOM` on the
  desktop, `?zoom=` on the web.
- `linux/`: an X11 window under Xvfb, Ubuntu 24.04 in Docker
  (`repro/linux/`), at 200% and 50%.
- `web/`: built for js/wasm (`repro/index.html`), in Chromium, at 200% and 50%.
- `macos/`: `repro/` run natively (`ZOOM=2` and `ZOOM=0.5`) on macOS 27.0.1, Apple silicon, Go 1.27.1, the 800x600 window captured with `screencapture -l` (Retina, so 1600x1200).
- `windows/`: `repro/` run natively (`ZOOM=2` and `ZOOM=0.5`) on Windows, the window captured.
