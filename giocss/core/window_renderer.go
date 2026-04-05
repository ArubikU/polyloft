package core

import (
	"image"
	"time"
)

type WindowFrameSnapshot struct {
	HoverState          map[string]bool
	ActiveState         map[string]bool
	StyleSheet          *StyleSheet
	RenderStore         *RenderStore
	CanUseWindowEffects bool
	FrameNumber         int64
	Debug               bool
	ProfileComponents   bool
	ProfileFull         bool
	RootLayout          map[string]any
	RootCSS             map[string]string
	ScreenW             int
	ScreenH             int
}

type WindowRunnerHooks struct {
	OnDestroy         func()
	BeforeFrame       func(size image.Point, frameTime time.Time) WindowFrameSnapshot
	AfterProcess      func(hoverState map[string]bool, activeState map[string]bool)
	AfterRender       func(renderStore *RenderStore, renderDur time.Duration, frameTime time.Time)
	DrawOverlay       func(rctx *RenderContext, state *GioWindowState, screenW int, screenH int)
	DispatchEvent     func(eventName string, payload map[string]any) error
	EmitRuntimeError  func(err error)
	RecordDebugSample func(path string, kind string, props map[string]any, w int, h int, totalDur time.Duration, selfDur time.Duration)
}

func RunWindowRenderer(win *Window, hooks WindowRunnerHooks) {
	if win == nil {
		return
	}
	state := NewGioWindowState()
	host := &GioRenderHost{}
	host.DispatchEvent = hooks.DispatchEvent
	host.EmitRuntimeError = hooks.EmitRuntimeError
	host.RecordDebugComponent = hooks.RecordDebugSample
	host.Invalidate = win.Invalidate

	frame := WindowFrameSnapshot{}
	RunGioWindowLoop(win, GioWindowLoopHooks{
		OnDestroy: hooks.OnDestroy,
		OnBeforeFrame: func(size image.Point, frameTime time.Time) {
			if hooks.BeforeFrame != nil {
				frame = hooks.BeforeFrame(size, frameTime)
			} else {
				frame = WindowFrameSnapshot{}
			}
			if frame.HoverState == nil {
				frame.HoverState = make(map[string]bool)
			}
			if frame.ActiveState == nil {
				frame.ActiveState = make(map[string]bool)
			}
			host.HoverState = frame.HoverState
			host.ActiveState = frame.ActiveState
			host.StyleSheet = frame.StyleSheet
			host.RenderStore = frame.RenderStore
			host.CanUseWindowEffects = frame.CanUseWindowEffects
			state.BeginFrame(size, frame.FrameNumber, host.HoverState, host.ActiveState, frame.Debug, host.StyleSheet, frame.ProfileComponents, frame.ProfileFull)
		},
		ProcessEvents: func(rctx *RenderContext, frameTime time.Time) bool {
			changed := ProcessGioEvents(rctx, host, state, frameTime)
			if hooks.AfterProcess != nil {
				hooks.AfterProcess(host.HoverState, host.ActiveState)
			}
			return changed
		},
		RenderFrame: func(rctx *RenderContext, frameTime time.Time) {
			renderStart := time.Now()
			if frame.RootLayout != nil && frame.ScreenW > 0 && frame.ScreenH > 0 {
				DrawGioBackground(rctx, 0, 0, frame.ScreenW, frame.ScreenH, frame.RootCSS)
				DrawGioTree(rctx, frame.RootLayout, host, "root", state, "", image.Rect(0, 0, frame.ScreenW, frame.ScreenH), nil, 0, 0, frameTime)
				DrawPickerModal(rctx, state, frame.ScreenW, frame.ScreenH)
				ApplyFrameCursorOverride(rctx, state)
			}
			if hooks.DrawOverlay != nil {
				hooks.DrawOverlay(rctx, state, frame.ScreenW, frame.ScreenH)
			}
			if hooks.AfterRender != nil {
				hooks.AfterRender(host.RenderStore, time.Since(renderStart), frameTime)
			}
		},
	})
}

