package ui

type RenderSnapshot struct {
	CSS    map[string]string
	Layout map[string]int
}

func NewRenderSnapshot(css map[string]string, x int, y int, w int, h int) RenderSnapshot {
	snap := RenderSnapshot{
		CSS:    map[string]string{},
		Layout: map[string]int{"x": x, "y": y, "width": w, "height": h},
	}
	for key, value := range css {
		snap.CSS[key] = value
	}
	return snap
}

func (s RenderSnapshot) Changed(css map[string]string, x int, y int, w int, h int) bool {
	if s.Layout["x"] != x || s.Layout["y"] != y || s.Layout["width"] != w || s.Layout["height"] != h {
		return true
	}
	if len(s.CSS) != len(css) {
		return true
	}
	for key, value := range css {
		if prev, ok := s.CSS[key]; !ok || prev != value {
			return true
		}
	}
	return false
}

type RenderStore struct {
	current map[string]RenderSnapshot
	next    map[string]RenderSnapshot
}

func NewRenderStore() *RenderStore {
	return &RenderStore{
		current: make(map[string]RenderSnapshot),
		next:    make(map[string]RenderSnapshot),
	}
}

func (s *RenderStore) Capture(path string, css map[string]string, x int, y int, w int, h int) RenderSnapshot {
	if s == nil {
		return NewRenderSnapshot(css, x, y, w, h)
	}
	snap := NewRenderSnapshot(css, x, y, w, h)
	s.next[path] = snap
	return snap
}

func (s *RenderStore) Previous(path string) (RenderSnapshot, bool) {
	if s == nil {
		return RenderSnapshot{}, false
	}
	prev, ok := s.current[path]
	return prev, ok
}

func (s *RenderStore) Finalize() {
	if s == nil {
		return
	}
	s.current = s.next
	s.next = make(map[string]RenderSnapshot)
}
