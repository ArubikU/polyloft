package layout

import (
	"encoding/json"
	"fmt"
	"strings"

	pkgcore "github.com/ArubikU/giocss/pkg"
	uicore "github.com/ArubikU/giocss/ui"
)

func nodeKey(node *Node, path string) string {
	if node == nil {
		return path
	}
	if key := strings.TrimSpace(nodePropString(node, "key", "")); key != "" {
		return key
	}
	if id := strings.TrimSpace(nodePropString(node, "id", "")); id != "" {
		return id
	}
	return path
}

func nodeSignature(node *Node) string {
	if node == nil {
		return "nil"
	}
	payload := map[string]any{
		"kind":     node.Tag,
		"props":    node.Props,
		"children": len(node.Children),
	}
	b, _ := json.Marshal(payload)
	return string(b)
}

func nodeSelfSignature(node *Node) string {
	if node == nil {
		return "nil"
	}
	payload := map[string]any{
		"kind":  node.Tag,
		"props": node.Props,
	}
	b, _ := json.Marshal(payload)
	return string(b)
}

func childStableKey(node *Node, index int) string {
	if node == nil {
		return fmt.Sprintf("#%d", index)
	}
	if key := strings.TrimSpace(nodePropString(node, "key", "")); key != "" {
		return key
	}
	if id := strings.TrimSpace(nodePropString(node, "id", "")); id != "" {
		return id
	}
	return fmt.Sprintf("#%d", index)
}

func flattenNodes(node *Node, path string, out map[string]string, kinds map[string]string) {
	if node == nil {
		return
	}
	key := nodeKey(node, path)
	out[key] = nodeSignature(node)
	kinds[key] = node.Tag
	for index, child := range node.Children {
		flattenNodes(child, fmt.Sprintf("%s/%d", key, index), out, kinds)
	}
}

func ReconcileTrees(oldRoot *Node, newRoot *Node) ([]map[string]any, []map[string]any) {
	patches := make([]map[string]any, 0, 16)
	fx := make([]map[string]any, 0, 16)
	var walk func(oldNode *Node, newNode *Node, path string)
	walk = func(oldNode *Node, newNode *Node, path string) {
		if oldNode == nil && newNode == nil {
			return
		}
		if oldNode == nil {
			patches = append(patches, map[string]any{"op": "add", "path": path, "kind": newNode.Tag})
			fx = append(fx, map[string]any{"type": "fade-in", "path": path, "durationMs": 140})
			return
		}
		if newNode == nil {
			patches = append(patches, map[string]any{"op": "remove", "path": path, "kind": oldNode.Tag})
			fx = append(fx, map[string]any{"type": "fade-out", "path": path, "durationMs": 120})
			return
		}

		if oldNode.Tag != newNode.Tag || nodeSelfSignature(oldNode) != nodeSelfSignature(newNode) {
			patches = append(patches, map[string]any{"op": "update", "path": path, "kind": newNode.Tag})
			fx = append(fx, map[string]any{"type": "morph", "path": path, "durationMs": 120})
		}

		oldByKey := make(map[string]*Node, len(oldNode.Children))
		newByKey := make(map[string]*Node, len(newNode.Children))
		oldIndex := make(map[string]int, len(oldNode.Children))
		for i, child := range oldNode.Children {
			key := childStableKey(child, i)
			oldByKey[key] = child
			oldIndex[key] = i
		}
		for i, child := range newNode.Children {
			key := childStableKey(child, i)
			newByKey[key] = child
		}

		for i, child := range newNode.Children {
			key := childStableKey(child, i)
			nextPath := path + "/" + key
			oldChild := oldByKey[key]
			if oldChild == nil {
				walk(nil, child, nextPath)
				continue
			}
			if oldIndex[key] != i {
				patches = append(patches, map[string]any{
					"op":     "move",
					"path":   nextPath,
					"kind":   child.Tag,
					"from":   oldIndex[key],
					"to":     i,
					"parent": path,
				})
				fx = append(fx, map[string]any{"type": "move", "path": nextPath, "from": oldIndex[key], "to": i, "durationMs": 180})
			}
			walk(oldChild, child, nextPath)
		}

		for i, child := range oldNode.Children {
			key := childStableKey(child, i)
			if _, ok := newByKey[key]; !ok {
				walk(child, nil, path+"/"+key)
			}
		}
	}

	walk(oldRoot, newRoot, "root")
	return patches, fx
}

func LayoutNodeToNative(node *Node, width int, height int, ss *uicore.StyleSheet) map[string]any {
	if node == nil {
		return map[string]any{}
	}
	rootCSS := ResolveNodeStyle(node, ss, width)
	if width <= 0 {
		width = NodeLayoutLength(node, rootCSS, "width", "width", 800, 800, 600, 800)
	}
	if height <= 0 {
		height = NodeLayoutLength(node, rootCSS, "height", "height", 600, width, 600, 600)
	}
	return layoutNodeToNativeWithBox(node, 0, 0, width, height, width, height, ss, nil)
}

func nativeNodeProps(node *Node) map[string]any {
	if node == nil {
		return map[string]any{}
	}
	out := make(map[string]any, len(node.Props)+2)
	for k, v := range node.Props {
		out[k] = v
	}
	if strings.TrimSpace(anyToString(out["tag"], "")) == "" && strings.TrimSpace(node.Tag) != "" {
		out["tag"] = node.Tag
	}
	if strings.TrimSpace(anyToString(out["text"], "")) == "" {
		if strings.TrimSpace(node.Text) != "" {
			out["text"] = node.Text
		}
	}
	return out
}

func layoutNodeToNativeWithBox(node *Node, x int, y int, width int, height int, viewportW int, viewportH int, ss *uicore.StyleSheet, parentCSS map[string]string) map[string]any {
	if node == nil {
		return map[string]any{}
	}
	props := nativeNodeProps(node)
	css := ResolveNodeStyle(node, ss, viewportW)
	css = MergeInheritedTextCSS(css, parentCSS)
	paddingX := nodePropInt(node, "padx", nodePropInt(node, "padding", 0))
	paddingY := nodePropInt(node, "pady", nodePropInt(node, "padding", 0))
	paddingTop := NodeLayoutLength(node, css, "paddingTop", "padding-top", height, viewportW, viewportH, paddingY)
	paddingRight := NodeLayoutLength(node, css, "paddingRight", "padding-right", width, viewportW, viewportH, paddingX)
	paddingBottom := NodeLayoutLength(node, css, "paddingBottom", "padding-bottom", height, viewportW, viewportH, paddingY)
	paddingLeft := NodeLayoutLength(node, css, "paddingLeft", "padding-left", width, viewportW, viewportH, paddingX)
	direction := NodeLayoutString(node, css, "direction", "direction", strings.ToLower(nodePropString(node, "layout", "column")))
	if direction != "row" {
		direction = "column"
	}
	justify := NodeLayoutString(node, css, "justify", "justify", "start")
	align := NodeLayoutString(node, css, "align", "align", "stretch")
	gap := NodeLayoutLength(node, css, "gap", "gap", max(width, height), viewportW, viewportH, 0)
	rowGap := NodeLayoutLength(node, css, "rowGap", "row-gap", max(width, height), viewportW, viewportH, gap)
	columnGap := NodeLayoutLength(node, css, "columnGap", "column-gap", max(width, height), viewportW, viewportH, gap)

	contentWidth := width - paddingLeft - paddingRight
	contentHeight := height - paddingTop - paddingBottom
	if contentWidth < 0 {
		contentWidth = 0
	}
	if contentHeight < 0 {
		contentHeight = 0
	}

	children := make([]any, 0, len(node.Children))
	if len(node.Children) > 0 {
		displayMode := strings.ToLower(strings.TrimSpace(css["display"]))
		if displayMode == "grid" {
			if strings.TrimSpace(css["justify"]) == "" {
				justify = "stretch"
			}
			if strings.TrimSpace(css["align"]) == "" {
				align = "stretch"
			}
			containerJustify := strings.ToLower(strings.TrimSpace(css["justify-content"]))
			if containerJustify == "" {
				containerJustify = strings.ToLower(strings.TrimSpace(css["justify"]))
			}
			if containerJustify == "" {
				containerJustify = "start"
			}
			itemJustifyDefault := strings.ToLower(strings.TrimSpace(css["justify-items"]))
			if itemJustifyDefault == "" {
				itemJustifyDefault = "stretch"
			}
			if place := strings.Fields(strings.ToLower(strings.TrimSpace(css["place-items"]))); len(place) > 0 {
				align = place[0]
				if len(place) > 1 {
					itemJustifyDefault = place[1]
				} else {
					itemJustifyDefault = place[0]
				}
			}
			if itemJustifyDefault == "start" {
				itemJustifyDefault = "left"
			}
			if itemJustifyDefault == "end" {
				itemJustifyDefault = "right"
			}
			gridColGap := columnGap
			gridRowGap := rowGap
			colTracks := ParseGridTrackSpec(css["grid-template-columns"], contentWidth, viewportW, viewportH)
			autoRows := strings.TrimSpace(css["grid-template-rows"]) == ""
			rowTracks := ParseGridTrackSpec(css["grid-template-rows"], contentHeight, viewportW, viewportH)
			{
				hint := len(colTracks)
				if hint > 1 {
					adj := max(1, contentWidth-(hint-1)*gridColGap)
					colTracks = ParseGridTrackSpec(css["grid-template-columns"], adj, viewportW, viewportH)
				}
			}
			if autoRows {
				rowTracks = make([]int, 0)
			}
			colStarts := make([]int, len(colTracks))
			cursorX := x + paddingLeft
			for i, cw := range colTracks {
				colStarts[i] = cursorX
				cursorX += cw + gridColGap
			}
			gridTracksWidth := 0
			for _, cw := range colTracks {
				gridTracksWidth += cw
			}
			if len(colTracks) > 1 {
				gridTracksWidth += (len(colTracks) - 1) * gridColGap
			}
			if gridTracksWidth < contentWidth {
				extra := contentWidth - gridTracksWidth
				offX := 0
				if containerJustify == "center" {
					offX = extra / 2
				} else if containerJustify == "end" || containerJustify == "right" {
					offX = extra
				}
				if offX > 0 {
					for i := range colStarts {
						colStarts[i] += offX
					}
				}
			}
			rowCursor := y + paddingTop
			rowIndex := 0
			colIndex := 0
			if len(rowTracks) == 0 {
				if autoRows {
					rowTracks = append(rowTracks, 0)
				} else {
					rowTracks = append(rowTracks, max(48, contentHeight))
				}
			}
			for _, child := range node.Children {
				childCSS := ResolveNodeStyle(child, ss, viewportW)
				childCSS = MergeInheritedTextCSS(childCSS, css)
				if strings.EqualFold(strings.TrimSpace(childCSS["display"]), "none") {
					continue
				}
				intrinsicW, intrinsicH := IntrinsicNodeSize(child, ss, viewportW, viewportH, css)
				span := pkgcore.CSSGridSpan(childCSS["grid-column"])
				if span < 1 {
					span = 1
				}
				if span > len(colTracks) {
					span = len(colTracks)
				}
				if colIndex+span > len(colTracks) {
					rowIndex++
					colIndex = 0
				}
				if rowIndex > 0 && rowIndex > len(rowTracks)-1 {
					if autoRows {
						rowTracks = append(rowTracks, 0)
					} else {
						rowTracks = append(rowTracks, max(48, contentHeight/max(1, (len(node.Children)+len(colTracks)-1)/len(colTracks))))
					}
				}
				childX := colStarts[colIndex]
				childW := 0
				for ci := 0; ci < span && colIndex+ci < len(colTracks); ci++ {
					childW += colTracks[colIndex+ci]
					if ci < span-1 {
						childW += gridColGap
					}
				}
				childH := rowTracks[rowIndex]
				if autoRows {
					if intrinsicH > 0 {
						childH = max(childH, intrinsicH)
					}
					if childH <= 0 {
						childH = 48
					}
					rowTracks[rowIndex] = childH
				} else if childH <= 0 {
					childH = intrinsicH
				}
				if childW <= 0 {
					childW = intrinsicW
				}
				childY := rowCursor
				for ri := 0; ri < rowIndex; ri++ {
					childY += rowTracks[ri] + gridRowGap
				}
				slotW := childW
				slotH := childH
				itemJustify := strings.ToLower(strings.TrimSpace(childCSS["justify-self"]))
				if itemJustify == "" || itemJustify == "auto" {
					itemJustify = itemJustifyDefault
				}
				itemAlign := strings.ToLower(strings.TrimSpace(childCSS["align-self"]))
				if itemAlign == "" || itemAlign == "auto" {
					itemAlign = align
				}
				renderW := slotW
				renderH := slotH
				explicitW := NodeLayoutLength(child, childCSS, "width", "width", slotW, viewportW, viewportH, -1)
				explicitH := NodeLayoutLength(child, childCSS, "height", "height", slotH, viewportW, viewportH, -1)
				if itemJustify != "stretch" && intrinsicW > 0 {
					renderW = min(slotW, intrinsicW)
				}
				if itemAlign != "stretch" && intrinsicH > 0 {
					renderH = min(slotH, intrinsicH)
				}
				if itemJustify == "stretch" && explicitW > 0 && explicitW < slotW {
					renderW = explicitW
				}
				if itemAlign == "stretch" && explicitH > 0 && explicitH < slotH {
					renderH = explicitH
				}
				renderX := childX
				renderY := childY
				if itemJustify == "center" {
					renderX += max(0, (slotW-renderW)/2)
				} else if itemJustify == "right" || itemJustify == "end" {
					renderX += max(0, slotW-renderW)
				} else if itemJustify == "stretch" && renderW < slotW {
					renderX += max(0, (slotW-renderW)/2)
				}
				if itemAlign == "center" {
					renderY += max(0, (slotH-renderH)/2)
				} else if itemAlign == "end" || itemAlign == "bottom" {
					renderY += max(0, slotH-renderH)
				} else if itemAlign == "stretch" && renderH < slotH {
					renderY += max(0, (slotH-renderH)/2)
				}
				children = append(children, layoutNodeToNativeWithBox(child, renderX, renderY, renderW, renderH, viewportW, viewportH, ss, css))
				colIndex += span
				if colIndex >= len(colTracks) {
					rowIndex++
					colIndex = 0
				}
			}
			return map[string]any{
				"kind":     node.Tag,
				"props":    props,
				"layout":   map[string]any{"x": x, "y": y, "width": width, "height": height},
				"children": children,
			}
		}

		mainSize := contentHeight
		crossSize := contentWidth
		mainGap := rowGap
		crossGap := columnGap
		if direction == "row" {
			mainSize = contentWidth
			crossSize = contentHeight
			mainGap = columnGap
			crossGap = rowGap
		}

		mainLens := make([]int, len(node.Children))
		crossLens := make([]int, len(node.Children))
		mainMarginBefore := make([]int, len(node.Children))
		mainMarginAfter := make([]int, len(node.Children))
		crossMarginBefore := make([]int, len(node.Children))
		crossMarginAfter := make([]int, len(node.Children))
		fixedMain := 0
		totalFlex := 0.0
		for i, child := range node.Children {
			childCSS := ResolveNodeStyle(child, ss, viewportW)
			childCSS = MergeInheritedTextCSS(childCSS, css)
			if strings.EqualFold(strings.TrimSpace(childCSS["display"]), "none") {
				continue
			}
			intrinsicW, intrinsicH := IntrinsicNodeSize(child, ss, viewportW, viewportH, css)
			itemAlignPref := strings.ToLower(strings.TrimSpace(childCSS["align-self"]))
			if itemAlignPref == "" || itemAlignPref == "auto" {
				itemAlignPref = align
			}
			if direction == "row" {
				mainMarginBefore[i] = NodeLayoutLength(child, childCSS, "marginLeft", "margin-left", mainSize, viewportW, viewportH, 0)
				mainMarginAfter[i] = NodeLayoutLength(child, childCSS, "marginRight", "margin-right", mainSize, viewportW, viewportH, 0)
				crossMarginBefore[i] = NodeLayoutLength(child, childCSS, "marginTop", "margin-top", crossSize, viewportW, viewportH, 0)
				crossMarginAfter[i] = NodeLayoutLength(child, childCSS, "marginBottom", "margin-bottom", crossSize, viewportW, viewportH, 0)
			} else {
				mainMarginBefore[i] = NodeLayoutLength(child, childCSS, "marginTop", "margin-top", mainSize, viewportW, viewportH, 0)
				mainMarginAfter[i] = NodeLayoutLength(child, childCSS, "marginBottom", "margin-bottom", mainSize, viewportW, viewportH, 0)
				crossMarginBefore[i] = NodeLayoutLength(child, childCSS, "marginLeft", "margin-left", crossSize, viewportW, viewportH, 0)
				crossMarginAfter[i] = NodeLayoutLength(child, childCSS, "marginRight", "margin-right", crossSize, viewportW, viewportH, 0)
			}
			flex := NodeLayoutFloat(child, childCSS, "flex", "flex", 0)
			totalFlex += flex
			if direction == "row" {
				cw := NodeLayoutLength(child, childCSS, "width", "width", mainSize, viewportW, viewportH, -1)
				hasExplicitMain := cw >= 0
				if cw < 0 && flex <= 0 {
					cw = intrinsicW
				}
				if cw >= 0 && (hasExplicitMain || flex <= 0) {
					minW := NodeLayoutLength(child, childCSS, "minWidth", "min-width", mainSize, viewportW, viewportH, -1)
					maxW := NodeLayoutLength(child, childCSS, "maxWidth", "max-width", mainSize, viewportW, viewportH, -1)
					if minW >= 0 && cw < minW {
						cw = minW
					}
					if maxW >= 0 && cw > maxW {
						cw = maxW
					}
					mainLens[i] = cw
					fixedMain += cw + mainMarginBefore[i] + mainMarginAfter[i]
				}
				ch := NodeLayoutLength(child, childCSS, "height", "height", crossSize, viewportW, viewportH, -1)
				if ch < 0 && itemAlignPref != "stretch" {
					ch = intrinsicH
				}
				if ch >= 0 {
					minH := NodeLayoutLength(child, childCSS, "minHeight", "min-height", crossSize, viewportW, viewportH, -1)
					maxH := NodeLayoutLength(child, childCSS, "maxHeight", "max-height", crossSize, viewportW, viewportH, -1)
					if minH >= 0 && ch < minH {
						ch = minH
					}
					if maxH >= 0 && ch > maxH {
						ch = maxH
					}
					crossLens[i] = ch
				}
			} else {
				ch := NodeLayoutLength(child, childCSS, "height", "height", mainSize, viewportW, viewportH, -1)
				hasExplicitMain := ch >= 0
				if ch < 0 && flex <= 0 {
					ch = intrinsicH
				}
				if ch >= 0 && (hasExplicitMain || flex <= 0) {
					minH := NodeLayoutLength(child, childCSS, "minHeight", "min-height", mainSize, viewportW, viewportH, -1)
					maxH := NodeLayoutLength(child, childCSS, "maxHeight", "max-height", mainSize, viewportW, viewportH, -1)
					if minH >= 0 && ch < minH {
						ch = minH
					}
					if maxH >= 0 && ch > maxH {
						ch = maxH
					}
					mainLens[i] = ch
					fixedMain += ch + mainMarginBefore[i] + mainMarginAfter[i]
				}
				cw := NodeLayoutLength(child, childCSS, "width", "width", crossSize, viewportW, viewportH, -1)
				if cw < 0 && itemAlignPref != "stretch" {
					cw = intrinsicW
				}
				if cw >= 0 {
					minW := NodeLayoutLength(child, childCSS, "minWidth", "min-width", crossSize, viewportW, viewportH, -1)
					maxW := NodeLayoutLength(child, childCSS, "maxWidth", "max-width", crossSize, viewportW, viewportH, -1)
					if minW >= 0 && cw < minW {
						cw = minW
					}
					if maxW >= 0 && cw > maxW {
						cw = maxW
					}
					crossLens[i] = cw
				}
			}
		}

		remaining := mainSize - fixedMain - (max(0, len(node.Children)-1) * mainGap)
		if remaining < 0 {
			remaining = 0
		}
		if totalFlex > 0 {
			for i, child := range node.Children {
				if mainLens[i] > 0 {
					continue
				}
				childCSS := ResolveNodeStyle(child, ss, viewportW)
				flex := NodeLayoutFloat(child, childCSS, "flex", "flex", 0)
				if flex > 0 {
					mainLens[i] = int((float64(remaining) * flex) / totalFlex)
				}
			}
		}

		totalMainUsed := 0
		for i, size := range mainLens {
			totalMainUsed += size + mainMarginBefore[i] + mainMarginAfter[i]
		}
		totalMainUsed += max(0, len(node.Children)-1) * mainGap
		extra := mainSize - totalMainUsed
		if extra < 0 {
			extra = 0
		}

		cursor := 0
		crossLineOffset := 0
		lineCrossMax := 0
		effectiveGap := mainGap
		wrap := strings.EqualFold(strings.TrimSpace(css["flex-wrap"]), "wrap")
		if justify == "center" {
			cursor = extra / 2
		} else if justify == "end" {
			cursor = extra
		} else if justify == "space-between" && len(node.Children) > 1 {
			effectiveGap = mainGap + (extra / (len(node.Children) - 1))
		} else if justify == "space-around" && len(node.Children) > 0 {
			effectiveGap = mainGap + (extra / len(node.Children))
			cursor = effectiveGap / 2
		}

		for i, child := range node.Children {
			childCSS := ResolveNodeStyle(child, ss, viewportW)
			childCSS = MergeInheritedTextCSS(childCSS, css)
			if strings.EqualFold(strings.TrimSpace(childCSS["display"]), "none") {
				continue
			}
			intrinsicW, intrinsicH := IntrinsicNodeSize(child, ss, viewportW, viewportH, css)
			aspectRatio, hasAspectRatio := uicore.ParseAspectRatio(childCSS)
			hasExplicitWidth := strings.TrimSpace(childCSS["width"]) != ""
			hasExplicitHeight := strings.TrimSpace(childCSS["height"]) != ""
			mainLen := mainLens[i]
			if mainLen <= 0 {
				if direction == "row" {
					mainLen = max(120, intrinsicW)
				} else {
					mainLen = max(32, intrinsicH)
				}
			}
			if direction == "row" {
				minW := NodeLayoutLength(child, childCSS, "minWidth", "min-width", mainSize, viewportW, viewportH, -1)
				maxW := NodeLayoutLength(child, childCSS, "maxWidth", "max-width", mainSize, viewportW, viewportH, -1)
				if minW >= 0 && mainLen < minW {
					mainLen = minW
				}
				if maxW >= 0 && mainLen > maxW {
					mainLen = maxW
				}
			} else {
				minH := NodeLayoutLength(child, childCSS, "minHeight", "min-height", mainSize, viewportW, viewportH, -1)
				maxH := NodeLayoutLength(child, childCSS, "maxHeight", "max-height", mainSize, viewportW, viewportH, -1)
				if minH >= 0 && mainLen < minH {
					mainLen = minH
				}
				if maxH >= 0 && mainLen > maxH {
					mainLen = maxH
				}
			}
			crossLen := crossLens[i]
			if crossLen <= 0 {
				if align == "stretch" {
					crossLen = crossSize
				} else if direction == "column" {
					crossLen = min(crossSize, max(intrinsicW, 32))
				} else {
					crossLen = max(intrinsicH, min(140, crossSize))
				}
			}
			if direction == "row" {
				minH := NodeLayoutLength(child, childCSS, "minHeight", "min-height", crossSize, viewportW, viewportH, -1)
				maxH := NodeLayoutLength(child, childCSS, "maxHeight", "max-height", crossSize, viewportW, viewportH, -1)
				if minH >= 0 && crossLen < minH {
					crossLen = minH
				}
				if maxH >= 0 && crossLen > maxH {
					crossLen = maxH
				}
			} else {
				minW := NodeLayoutLength(child, childCSS, "minWidth", "min-width", crossSize, viewportW, viewportH, -1)
				maxW := NodeLayoutLength(child, childCSS, "maxWidth", "max-width", crossSize, viewportW, viewportH, -1)
				if minW >= 0 && crossLen < minW {
					crossLen = minW
				}
				if maxW >= 0 && crossLen > maxW {
					crossLen = maxW
				}
			}
			itemAlign := strings.ToLower(strings.TrimSpace(childCSS["align-self"]))
			if itemAlign == "" || itemAlign == "auto" {
				itemAlign = align
			}
			if hasAspectRatio && aspectRatio > 0.001 {
				if direction == "row" {
					if !hasExplicitHeight && itemAlign != "stretch" {
						crossLen = max(1, int(float64(mainLen)/aspectRatio))
					}
				} else {
					if !hasExplicitWidth && itemAlign != "stretch" {
						crossLen = max(1, int(float64(mainLen)*aspectRatio))
					}
				}
			}
			crossLen = max(1, crossLen-crossMarginBefore[i]-crossMarginAfter[i])
			crossPos := 0
			if itemAlign == "center" {
				crossPos = (crossSize - crossLen) / 2
			} else if itemAlign == "end" || itemAlign == "bottom" {
				crossPos = crossSize - crossLen
			}
			if crossPos < 0 {
				crossPos = 0
			}

			childX := x + paddingLeft
			childY := y + paddingTop
			childW := crossLen
			childH := mainLen
			if wrap && direction == "row" && cursor+mainMarginBefore[i]+mainLen+mainMarginAfter[i] > mainSize && cursor > 0 {
				cursor = 0
				crossLineOffset += lineCrossMax + crossGap
				lineCrossMax = 0
			}
			if wrap && direction == "column" && cursor+mainMarginBefore[i]+mainLen+mainMarginAfter[i] > mainSize && cursor > 0 {
				cursor = 0
				crossLineOffset += lineCrossMax + crossGap
				lineCrossMax = 0
			}
			if direction == "row" {
				childX += cursor + mainMarginBefore[i]
				childY += crossPos + crossMarginBefore[i] + crossLineOffset
				childW = mainLen
				childH = crossLen
			} else {
				childX += crossPos + crossMarginBefore[i] + crossLineOffset
				childY += cursor + mainMarginBefore[i]
				childW = crossLen
				childH = mainLen
			}
			children = append(children, layoutNodeToNativeWithBox(child, childX, childY, childW, childH, viewportW, viewportH, ss, css))
			if direction == "row" {
				lineCrossMax = max(lineCrossMax, childH+crossMarginBefore[i]+crossMarginAfter[i])
			} else {
				lineCrossMax = max(lineCrossMax, childW+crossMarginBefore[i]+crossMarginAfter[i])
			}
			cursor += mainMarginBefore[i] + mainLen + mainMarginAfter[i] + effectiveGap
		}
	}

	return map[string]any{
		"kind":  node.Tag,
		"props": props,
		"layout": map[string]any{
			"x":      x,
			"y":      y,
			"width":  width,
			"height": height,
		},
		"children": children,
	}
}
