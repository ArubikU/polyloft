# giocss

Go CSS engine extracted from Polyloft BVM.

This module currently provides:

- StyleSheet creation and mutation
- CSS text parsing and css file loading
- Style resolution with class, tag, inline, and state maps
- CSS canonical property mapping
- CSS length and float value parsing helpers
- Pseudo-HTML DOM model (`Node`) with style resolution helpers
- Layout helpers (`NodeLayout*`, `ParseGridTrackSpec`, inherited text CSS merge)
- Color/text utility helpers extracted from polyloft-bvm runtime

## Pseudo-HTML model

`giocss` now owns a lightweight HTML-like model that Polyloft can wrap:

- `NewNode(tag string)`
- `Node.SetProp(name, value)`
- `Node.AddClass(className)`
- `Node.AddChild(child)`
- `ResolveNodeStyle(node, stylesheet, viewportW)`

This keeps CSS and pseudo-HTML semantics in one module, while host runtimes
(polyloft-bvm, future adapters) handle app/window/lifecycle orchestration.

## Install

```bash
go get github.com/ArubikU/giocss@v0.1.0
```

## Local development

```bash
go test ./...
```

## Intended consumers

- polyloft-bvm runtime UI
- future external Go render/layout adapters
