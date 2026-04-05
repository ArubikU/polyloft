# Files Organization

This map describes the current migration state toward domain subpackages.

## Root compatibility facade

- exports.go

## Extracted domain subpackages (phase in progress)

### events

- events/renderer.go

### common

- common/misc_utils.go

### input

- input/cursor_runtime.go

### interaction

- interaction/state_runtime.go

### layout

- layout/renderer.go

### render

- render/font.go
- render/gio_state.go
- render/store.go

### scroll

- scroll/logic.go

### style

- style/color_helpers.go
- style/state_helpers.go
- style/stylesheet.go
- style/text_helpers.go
- style/transition.go

## Core implementation package

- core/*.go (migrated legacy implementation from src)

## Compatibility source package

- src/compat_types.go (type aliases to core)
- src/compat_api.go (function wrappers to core)

## Tests

- tests/helpers_test.go
- tests/input_normalization_test.go
- tests/layout_reconcile_test.go
