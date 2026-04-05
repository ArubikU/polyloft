package giocss_test

import (
	"testing"

	giocss "github.com/ArubikU/giocss"
)

func TestReconcileTreesAddsNewRoot(t *testing.T) {
	newRoot := giocss.NewNode("button")

	patches, fx := giocss.ReconcileTrees(nil, newRoot)

	if len(patches) != 1 {
		t.Fatalf("len(patches) = %d, want 1", len(patches))
	}
	if len(fx) != 1 {
		t.Fatalf("len(fx) = %d, want 1", len(fx))
	}
	if got := patches[0]["op"]; got != "add" {
		t.Fatalf("patch op = %v, want add", got)
	}
	if got := patches[0]["path"]; got != "root" {
		t.Fatalf("patch path = %v, want root", got)
	}
	if got := patches[0]["kind"]; got != "button" {
		t.Fatalf("patch kind = %v, want button", got)
	}
	if got := fx[0]["type"]; got != "fade-in" {
		t.Fatalf("fx type = %v, want fade-in", got)
	}
}

func TestReconcileTreesDetectsKeyedMoves(t *testing.T) {
	oldRoot := giocss.NewNode("div")
	oldA := giocss.NewNode("label")
	oldA.SetProp("key", "a")
	oldB := giocss.NewNode("label")
	oldB.SetProp("key", "b")
	oldRoot.AddChild(oldA)
	oldRoot.AddChild(oldB)

	newRoot := giocss.NewNode("div")
	newB := giocss.NewNode("label")
	newB.SetProp("key", "b")
	newA := giocss.NewNode("label")
	newA.SetProp("key", "a")
	newRoot.AddChild(newB)
	newRoot.AddChild(newA)

	patches, fx := giocss.ReconcileTrees(oldRoot, newRoot)

	if len(patches) != 2 {
		t.Fatalf("len(patches) = %d, want 2", len(patches))
	}
	assertMovePatch(t, patches[0], "root/b", 1, 0)
	assertMovePatch(t, patches[1], "root/a", 0, 1)
	if len(fx) != 2 {
		t.Fatalf("len(fx) = %d, want 2", len(fx))
	}
	if got := fx[0]["type"]; got != "move" {
		t.Fatalf("fx[0] type = %v, want move", got)
	}
	if got := fx[1]["type"]; got != "move" {
		t.Fatalf("fx[1] type = %v, want move", got)
	}
}

func TestLayoutNodeToNativeUsesStyleSheetFallbackDimensions(t *testing.T) {
	ss := giocss.NewStyleSheet()
	ss.SetRule("panel", "width", "320px")
	ss.SetRule("panel", "height", "180px")

	root := giocss.NewNode("div")
	root.AddClass("panel")

	native := giocss.LayoutNodeToNative(root, 0, 0, ss)
	layout := nativeLayout(t, native)

	if got := layout["width"]; got != 320 {
		t.Fatalf("layout width = %v, want 320", got)
	}
	if got := layout["height"]; got != 180 {
		t.Fatalf("layout height = %v, want 180", got)
	}
}

func TestLayoutNodeToNativeSerializesChildrenAndProps(t *testing.T) {
	root := giocss.NewNode("div")
	root.SetProp("id", "root")
	root.SetProp("width", 200)
	root.SetProp("height", 100)

	child := giocss.NewNode("button")
	child.SetProp("id", "cta")
	child.SetProp("width", 50)
	child.SetProp("height", 30)
	root.AddChild(child)

	native := giocss.LayoutNodeToNative(root, 0, 0, nil)

	if got := native["kind"]; got != "div" {
		t.Fatalf("native kind = %v, want div", got)
	}
	layout := nativeLayout(t, native)
	if got := layout["width"]; got != 200 {
		t.Fatalf("layout width = %v, want 200", got)
	}
	if got := layout["height"]; got != 100 {
		t.Fatalf("layout height = %v, want 100", got)
	}

	children, ok := native["children"].([]any)
	if !ok {
		t.Fatalf("children type = %T, want []any", native["children"])
	}
	if len(children) != 1 {
		t.Fatalf("len(children) = %d, want 1", len(children))
	}
	childNative, ok := children[0].(map[string]any)
	if !ok {
		t.Fatalf("child type = %T, want map[string]any", children[0])
	}
	if got := childNative["kind"]; got != "button" {
		t.Fatalf("child kind = %v, want button", got)
	}
	childProps, ok := childNative["props"].(map[string]any)
	if !ok {
		t.Fatalf("child props type = %T, want map[string]any", childNative["props"])
	}
	if got := childProps["id"]; got != "cta" {
		t.Fatalf("child props id = %v, want cta", got)
	}
	childLayout := nativeLayout(t, childNative)
	if got := childLayout["width"]; got != 50 {
		t.Fatalf("child layout width = %v, want 50", got)
	}
	if got := childLayout["height"]; got != 30 {
		t.Fatalf("child layout height = %v, want 30", got)
	}
}

func assertMovePatch(t *testing.T, patch map[string]any, path string, from int, to int) {
	t.Helper()
	if got := patch["op"]; got != "move" {
		t.Fatalf("patch op = %v, want move", got)
	}
	if got := patch["path"]; got != path {
		t.Fatalf("patch path = %v, want %s", got, path)
	}
	if got := patch["from"]; got != from {
		t.Fatalf("patch from = %v, want %d", got, from)
	}
	if got := patch["to"]; got != to {
		t.Fatalf("patch to = %v, want %d", got, to)
	}
}

func nativeLayout(t *testing.T, node map[string]any) map[string]any {
	t.Helper()
	layout, ok := node["layout"].(map[string]any)
	if !ok {
		t.Fatalf("layout type = %T, want map[string]any", node["layout"])
	}
	return layout
}
