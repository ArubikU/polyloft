package ui

import (
	"sync"

	"gioui.org/font/gofont"
	"gioui.org/text"
)

type ResolvedCSSCacheEntry struct {
	Hovered   bool
	Active    bool
	ViewW     int
	LowPower  bool
	ParentSig string
	LastSeen  int64
	CSS       map[string]string
}

type ZChildrenHintCacheEntry struct {
	Count    int
	Sig      uint64
	Needs    bool
	LastSeen int64
}

var RenderConfigOnce sync.Once

func NewGioShaper() *text.Shaper {
	return text.NewShaper(text.WithCollection(gofont.Collection()))
}

func PurgeStaleRenderCaches(resolvedCSS map[string]ResolvedCSSCacheEntry, zChildrenHint map[string]ZChildrenHintCacheEntry, frame int64) {
	if frame <= 0 {
		return
	}
	const ttl = int64(240)
	for key, entry := range resolvedCSS {
		if frame-entry.LastSeen > ttl {
			delete(resolvedCSS, key)
		}
	}
	for key, entry := range zChildrenHint {
		if frame-entry.LastSeen > ttl {
			delete(zChildrenHint, key)
		}
	}
}
