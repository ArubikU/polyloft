package layout

import (
	"image"

	layoutcore "github.com/ArubikU/giocss/layout"
	uicore "github.com/ArubikU/giocss/ui"
)

type PositionLayoutResult = layoutcore.PositionLayoutResult

type OverflowLayoutResult = layoutcore.OverflowLayoutResult

type RenderVisibilityResult = layoutcore.RenderVisibilityResult

type ChildTraversalPlan = layoutcore.ChildTraversalPlan

func toLayoutPositionInfo(pos uicore.PositionInfo) layoutcore.PositionInfo {
	return layoutcore.PositionInfo{
		Position:  pos.Position,
		OffsetX:   pos.OffsetX,
		OffsetY:   pos.OffsetY,
		IsLeft:    pos.IsLeft,
		IsTop:     pos.IsTop,
		HasLeft:   pos.HasLeft,
		HasRight:  pos.HasRight,
		HasTop:    pos.HasTop,
		HasBottom: pos.HasBottom,
	}
}

func ApplyPositionLayout(x, y, w, h int, pos uicore.PositionInfo, containingBox image.Rectangle, viewport image.Rectangle) PositionLayoutResult {
	return layoutcore.ApplyPositionLayout(x, y, w, h, toLayoutPositionInfo(pos), containingBox, viewport)
}

func ShouldQuickCullNode(path string, x, y, w, h, inheritedShiftX, inheritedShiftY int, frameViewport image.Rectangle, viewW, viewH, childCount int, kind string) bool {
	return layoutcore.ShouldQuickCullNode(path, x, y, w, h, inheritedShiftX, inheritedShiftY, frameViewport, viewW, viewH, childCount, kind)
}

func ShouldProfileNode(profileComponents, profileFull bool, path string, profileSampleFrame bool) bool {
	return layoutcore.ShouldProfileNode(profileComponents, profileFull, path, profileSampleFrame)
}

func ResolveOverflowLayout(css map[string]string) OverflowLayoutResult {
	return layoutcore.ResolveOverflowLayout(css)
}

func ComputeRenderVisibility(path string, x, y, w, h int, inheritedShiftX, inheritedShiftY int, translateX, translateY, scale float64, frameViewport image.Rectangle, viewW, viewH int) RenderVisibilityResult {
	return layoutcore.ComputeRenderVisibility(path, x, y, w, h, inheritedShiftX, inheritedShiftY, translateX, translateY, scale, frameViewport, viewW, viewH)
}

func ShouldRegisterContainerEvents(kind string, hasHover, hasActive bool, overflowX, overflowY string) bool {
	return layoutcore.ShouldRegisterContainerEvents(kind, hasHover, hasActive, overflowX, overflowY)
}

func ShouldApplyNodeCursor(cursorRaw string) bool {
	return layoutcore.ShouldApplyNodeCursor(cursorRaw)
}

func ShouldPromoteFrameCursor(path, currentPath string, pointerKnown bool, pointerPos image.Point, renderRect image.Rectangle) bool {
	return layoutcore.ShouldPromoteFrameCursor(path, currentPath, pointerKnown, pointerPos, renderRect)
}

func BuildChildTraversalPlan(visible, isOnScreen, clipChildren bool, inheritedShiftX, inheritedShiftY int, translateX, translateY float64, scroll image.Point) ChildTraversalPlan {
	return layoutcore.BuildChildTraversalPlan(visible, isOnScreen, clipChildren, inheritedShiftX, inheritedShiftY, translateX, translateY, scroll)
}

func IsDisplayNone(displayValue string) bool {
	return layoutcore.IsDisplayNone(displayValue)
}

func IsVisibilityHidden(visibilityValue string) bool {
	return layoutcore.IsVisibilityHidden(visibilityValue)
}

func IsFilterEnabled(filterRaw string) bool {
	return layoutcore.IsFilterEnabled(filterRaw)
}

func BuildFilterFactorVars(filterRaw string) map[string]string {
	return layoutcore.BuildFilterFactorVars(filterRaw)
}
