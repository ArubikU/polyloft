package render

import (
	"sort"

	coreengine "github.com/ArubikU/giocss/core/engine"
)

type ZSortHint struct {
	Count    int
	Sig      uint64
	Needs    bool
	LastSeen int64
}

type ChildTraversalDecision struct {
	Traverse bool
	Order    []int
}

func ChildrenNeedZSort(children []map[string]any) bool {
	for _, child := range children {
		if coreengine.NodeHasExplicitZIndex(mapProps(child)) {
			return true
		}
	}
	return false
}

func BuildChildrenZOrder(children []map[string]any) []int {
	type item struct {
		idx int
		z   int
	}
	items := make([]item, 0, len(children))
	for i, child := range children {
		items = append(items, item{idx: i, z: coreengine.ParseZIndexFromPropsFast(mapProps(child))})
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].z == items[j].z {
			return items[i].idx < items[j].idx
		}
		return items[i].z < items[j].z
	})
	order := make([]int, len(items))
	for i, it := range items {
		order[i] = it.idx
	}
	return order
}

func ChildrenMayContainOutOfFlow(children []map[string]any) bool {
	for _, child := range children {
		if coreengine.PropsMayBeOutOfFlow(mapProps(child)) {
			return true
		}
	}
	return false
}

func ResolveChildrenZSortNeed(children []map[string]any, existing *ZSortHint, childSig uint64, frameNumber int64) (bool, ZSortHint) {
	if existing != nil && existing.Count == len(children) && existing.Sig == childSig {
		out := *existing
		out.LastSeen = frameNumber
		return out.Needs, out
	}
	needs := ChildrenNeedZSort(children)
	return needs, ZSortHint{
		Count:    len(children),
		Sig:      childSig,
		Needs:    needs,
		LastSeen: frameNumber,
	}
}

func PlanChildTraversal(children []map[string]any, isOnScreen bool, needsZSort bool) ChildTraversalDecision {
	if len(children) == 0 {
		return ChildTraversalDecision{}
	}
	if isOnScreen {
		if needsZSort {
			return ChildTraversalDecision{Traverse: true, Order: BuildChildrenZOrder(children)}
		}
		order := make([]int, len(children))
		for i := range children {
			order[i] = i
		}
		return ChildTraversalDecision{Traverse: true, Order: order}
	}
	if !ChildrenMayContainOutOfFlow(children) {
		return ChildTraversalDecision{}
	}
	order := make([]int, len(children))
	for i := range children {
		order[i] = i
	}
	return ChildTraversalDecision{Traverse: true, Order: order}
}

func ChildrenContentBounds(children []map[string]any, fallbackRight int, fallbackBottom int) (int, int) {
	maxRight := fallbackRight
	maxBottom := fallbackBottom
	for _, child := range children {
		layoutMap, ok := child["layout"].(map[string]any)
		if !ok {
			continue
		}
		x := toIntValue(layoutMap["x"], 0)
		y := toIntValue(layoutMap["y"], 0)
		w := toIntValue(layoutMap["width"], 0)
		h := toIntValue(layoutMap["height"], 0)
		if r := x + w; r > maxRight {
			maxRight = r
		}
		if b := y + h; b > maxBottom {
			maxBottom = b
		}
	}
	return maxRight, maxBottom
}

func toIntValue(v any, fallback int) int {
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	case float32:
		return int(n)
	default:
		return fallback
	}
}

func mapProps(child map[string]any) map[string]any {
	if len(child) == 0 {
		return nil
	}
	if v, ok := child["props"].(map[string]any); ok {
		return v
	}
	return nil
}
