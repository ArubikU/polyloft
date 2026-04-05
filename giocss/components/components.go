// Package components provides ergonomic constructors for giocss Node trees.
// Usage:
//
//	import "github.com/ArubikU/giocss/components"
//
//	btn   := components.Button("Click me", "btn", "btn-primary")
//	input := components.Input("email", "Email address", "email")
//	card  := components.Card(btn, input)
package components

import (
	"fmt"

	giocss "github.com/ArubikU/giocss"
)

// Button creates a <button> node with optional CSS classes.
func Button(text string, classes ...string) *giocss.Node {
	n := giocss.NewNode("button")
	n.Text = text
	for _, c := range classes {
		n.AddClass(c)
	}
	return n
}

// Input creates an <input> node with id, placeholder and type props.
func Input(id, placeholder, inputType string) *giocss.Node {
	n := giocss.NewNode("input")
	n.SetProp("id", id)
	if id != "" {
		n.SetProp("name", id)
	}
	n.SetProp("placeholder", placeholder)
	n.SetProp("type", inputType)
	return n
}

// Form creates a <form> node with id, optional CSS classes, and children.
func Form(id string, classes []string, children ...*giocss.Node) *giocss.Node {
	n := giocss.NewNode("form")
	if id != "" {
		n.SetProp("id", id)
	}
	for _, c := range classes {
		n.AddClass(c)
	}
	for _, ch := range children {
		n.AddChild(ch)
	}
	return n
}

// SubmitButton creates a button wired for form submission.
// If formID is set, it behaves like HTML's "form" attribute.
func SubmitButton(text, formID string, classes ...string) *giocss.Node {
	n := Button(text, classes...)
	n.SetProp("type", "submit")
	if formID != "" {
		n.SetProp("form", formID)
	}
	return n
}

// Label creates a <label> node with optional classes.
func Label(text string, classes ...string) *giocss.Node {
	n := giocss.NewNode("label")
	n.Text = text
	for _, c := range classes {
		n.AddClass(c)
	}
	return n
}

// Text creates a <text> node (inline text run) with optional classes.
func Text(content string, classes ...string) *giocss.Node {
	n := giocss.NewNode("text")
	n.Text = content
	for _, c := range classes {
		n.AddClass(c)
	}
	return n
}

// Heading creates a heading node h1–h6. Level is clamped to 1–6.
func Heading(level int, text string, classes ...string) *giocss.Node {
	if level < 1 {
		level = 1
	}
	if level > 6 {
		level = 6
	}
	n := giocss.NewNode(fmt.Sprintf("h%d", level))
	n.Text = text
	for _, c := range classes {
		n.AddClass(c)
	}
	return n
}

// Div creates a generic <div> container with optional CSS classes and children.
func Div(classes []string, children ...*giocss.Node) *giocss.Node {
	n := giocss.NewNode("div")
	for _, c := range classes {
		n.AddClass(c)
	}
	for _, ch := range children {
		n.AddChild(ch)
	}
	return n
}

// Span creates an inline <span> node with optional classes.
func Span(text string, classes ...string) *giocss.Node {
	n := giocss.NewNode("span")
	n.Text = text
	for _, c := range classes {
		n.AddClass(c)
	}
	return n
}

// Card wraps children in a <div class="card">.
func Card(children ...*giocss.Node) *giocss.Node {
	n := giocss.NewNode("div")
	n.AddClass("card")
	for _, ch := range children {
		n.AddChild(ch)
	}
	return n
}

// Badge creates a <span class="badge"> node with optional extra classes.
func Badge(text string, classes ...string) *giocss.Node {
	n := giocss.NewNode("span")
	n.Text = text
	n.AddClass("badge")
	for _, c := range classes {
		n.AddClass(c)
	}
	return n
}

// NavBar creates a <nav> node with a brand text node and a list of item nodes.
func NavBar(brand string, items []string) *giocss.Node {
	nav := giocss.NewNode("nav")
	nav.AddClass("navbar")

	brandNode := giocss.NewNode("span")
	brandNode.Text = brand
	brandNode.AddClass("navbar-brand")
	nav.AddChild(brandNode)

	itemsRow := giocss.NewNode("div")
	itemsRow.AddClass("navbar-items")
	for _, item := range items {
		a := giocss.NewNode("span")
		a.Text = item
		a.AddClass("navbar-item")
		itemsRow.AddChild(a)
	}
	nav.AddChild(itemsRow)
	return nav
}

// Row creates a flex-row <div> (class "row") containing the given children.
func Row(children ...*giocss.Node) *giocss.Node {
	n := giocss.NewNode("div")
	n.AddClass("row")
	for _, ch := range children {
		n.AddChild(ch)
	}
	return n
}

// Column creates a flex-column <div> (class "col") containing the given children.
func Column(children ...*giocss.Node) *giocss.Node {
	n := giocss.NewNode("div")
	n.AddClass("col")
	for _, ch := range children {
		n.AddChild(ch)
	}
	return n
}
