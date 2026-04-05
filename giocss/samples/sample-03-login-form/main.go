package main

import (
	"image"
	"strings"
	"sync"

	giocss "github.com/ArubikU/giocss"
	"github.com/ArubikU/giocss/components"
)

const css = `
body {
  background-color: #f3f4f6;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.form-card {
  background-color: #ffffff;
  border-radius: 16px;
  padding: 40px 48px;
  width: 400px;
  display: flex;
  flex-direction: column;
}

.form-title {
  font-size: 24px;
  font-weight: bold;
  color: #111827;
  text-align: center;
  margin-bottom: 8px;
}

.form-subtitle {
  font-size: 14px;
  color: #6b7280;
  text-align: center;
  margin-bottom: 28px;
}

.form-label {
  font-size: 13px;
  font-weight: bold;
  color: #374151;
  margin-bottom: 6px;
}

.form-input {
  background-color: #f9fafb;
  border: 1.5px solid #d1d5db;
  border-radius: 8px;
  padding: 10px 14px;
  font-size: 14px;
  color: #111827;
  margin-bottom: 18px;
  width: 100%;
}

.form-input:hover {
  border-color: #9ca3af;
}

.form-submit {
  background-color: #4f46e5;
  color: #ffffff;
  border-radius: 8px;
  padding: 12px 0;
  font-size: 15px;
  font-weight: bold;
  text-align: center;
  cursor: pointer;
  margin-top: 8px;
}

.form-submit:hover {
  background-color: #4338ca;
}

.form-submit:active {
  background-color: #3730a3;
}

.form-status {
  font-size: 12px;
  color: #475569;
  text-align: center;
  margin-top: 14px;
}

.form-status-ok {
  color: #166534;
  font-weight: bold;
}
`

type appState struct {
	mu       sync.Mutex
	email    string
	password string
	focused  string
	status   string
	success  bool
}

func newAppState() *appState {
	return &appState{status: "Waiting for submit...", success: false}
}

func (s *appState) onEvent(eventName string, payload map[string]any) error {
	eventName = strings.TrimSpace(eventName)
	value, _ := payload["value"].(string)

	s.mu.Lock()
	defer s.mu.Unlock()

	switch eventName {
	case "auth.login.email.input":
		s.email = value
		if strings.TrimSpace(s.email) == "" {
			s.status = "Type your email to start testing the input."
			s.success = false
		} else if strings.TrimSpace(s.password) == "" {
			s.status = "Email captured. Now type your password."
			s.success = false
		} else {
			s.status = "Ready to submit."
			s.success = false
		}
	case "auth.login.password.input":
		s.password = value
		if strings.TrimSpace(s.password) == "" {
			s.status = "Password cleared."
			s.success = false
		} else if strings.TrimSpace(s.email) == "" {
			s.status = "Password captured. Email is still missing."
			s.success = false
		} else {
			s.status = "Ready to submit."
			s.success = false
		}
	case "auth.login.email.focus":
		s.focused = "email"
		if strings.TrimSpace(s.email) == "" {
			s.status = "Editing email..."
		}
	case "auth.login.password.focus":
		s.focused = "password"
		if strings.TrimSpace(s.password) == "" {
			s.status = "Editing password..."
		}
	case "auth.login.email.blur", "auth.login.password.blur":
		s.focused = ""
		if strings.TrimSpace(s.email) != "" && strings.TrimSpace(s.password) != "" && !s.success {
			s.status = "Ready to submit."
		}
	case "auth.login.submit":
		values, _ := payload["values"].(map[string]any)
		email, _ := values["email"].(string)
		password, _ := values["password"].(string)
		s.email = email
		s.password = password
		if strings.TrimSpace(email) == "" {
			s.status = "Email is required."
			s.success = false
			return nil
		}
		if !strings.Contains(email, "@") {
			s.status = "Enter a valid email address."
			s.success = false
			return nil
		}
		if len(strings.TrimSpace(password)) < 6 {
			s.status = "Password must contain at least 6 characters."
			s.success = false
			return nil
		}
		s.status = "Submit trigger received for " + email
		s.success = true
	}
	return nil
}

func (s *appState) snapshot() (string, string, string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.email, s.password, s.status, s.success
}

func buildUI(email string, password string, status string, success bool) *giocss.Node {
	root := giocss.NewNode("body")

	emailInput := components.Input("email", "you@example.com", "email")
	emailInput.AddClass("form-input")
	emailInput.SetProp("value", email)
	emailInput.SetProp("oninput", "auth.login.email.input")
	emailInput.SetProp("onfocus", "auth.login.email.focus")
	emailInput.SetProp("onblur", "auth.login.email.blur")
	passwordInput := components.Input("password", "••••••••", "password")
	passwordInput.AddClass("form-input")
	passwordInput.SetProp("value", password)
	passwordInput.SetProp("oninput", "auth.login.password.input")
	passwordInput.SetProp("onfocus", "auth.login.password.focus")
	passwordInput.SetProp("onblur", "auth.login.password.blur")

	statusNode := components.Text(status, "form-status")
	if success {
		statusNode.AddClass("form-status-ok")
	}

	form := components.Form("login-form", []string{"form-card"},
		components.Heading(2, "Sign in", "form-title"),
		components.Text("Welcome back! Please enter your details.", "form-subtitle"),
		components.Label("Email", "form-label"),
		emailInput,
		components.Label("Password", "form-label"),
		passwordInput,
		components.SubmitButton("Sign in", "login-form", "form-submit"),
		statusNode,
	)
	form.SetProp("onsubmit", "auth.login.submit")
	root.AddChild(form)

	return root
}

func main() {
	ss := giocss.NewStyleSheet()
	ss.ParseCSSText(css)
	app := newAppState()

	rt := giocss.NewWindowRuntime(
		giocss.WindowOptions{Title: "Sample 03 – Login Form", Width: 600, Height: 520},
		giocss.WindowRuntimeHooks{
			DispatchEvent: app.onEvent,
			Snapshot: func(size image.Point) giocss.WindowRuntimeSnapshot {
				email, password, status, success := app.snapshot()
				root := buildUI(email, password, status, success)
				return giocss.WindowRuntimeSnapshot{
					RootLayout:   giocss.LayoutNodeToNative(root, size.X, size.Y, ss),
					RootCSS:      giocss.ResolveNodeStyle(root, ss, size.X),
					StyleSheet:   ss,
					ScreenWidth:  size.X,
					ScreenHeight: size.Y,
				}
			},
		},
	)
	rt.Run()
}
