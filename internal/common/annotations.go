package common

import "strings"

// AnnotationFlags captures semantic toggles associated with annotations.
type AnnotationFlags struct {
	IsOverride bool // marks methods that should override a parent implementation
}

// Merge combines two sets of annotation flags.
func (f AnnotationFlags) Merge(other AnnotationFlags) AnnotationFlags {
	return AnnotationFlags{
		IsOverride: f.IsOverride || other.IsOverride,
	}
}

// AnnotationInfo describes a registered annotation and its semantic flags.
type AnnotationInfo struct {
	Name  string
	Flags AnnotationFlags
}

var annotationRegistry = map[string]AnnotationInfo{}

// NormalizeAnnotationName converts an annotation identifier into its canonical form.
func NormalizeAnnotationName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// RegisterAnnotation stores metadata for an annotation in the registry.
func RegisterAnnotation(name string, flags AnnotationFlags) {
	normalized := NormalizeAnnotationName(name)
	annotationRegistry[normalized] = AnnotationInfo{
		Name:  normalized,
		Flags: flags,
	}
}

// LookupAnnotation resolves annotation metadata by name, returning whether it is known.
func LookupAnnotation(name string) (AnnotationInfo, bool) {
	normalized := NormalizeAnnotationName(name)
	info, ok := annotationRegistry[normalized]
	if ok {
		return info, true
	}
	return AnnotationInfo{Name: normalized}, false
}

func init() {
	RegisterAnnotation("override", AnnotationFlags{IsOverride: true})
}
