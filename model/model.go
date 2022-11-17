package model

// MetaNode wraps a more detailed Node allowing efficient parsing of
// Node content at the highest level when when the specific types of
// content beyond its meta data are not required. This is useful for
// quickly creating indexes based on the title or other meta data (such
// as tags). A MetaNode is also the best schema to serve as an API
// allowing the client to decide if further parsing is needed.
type MetaNode struct {
	Title string         `json:"title,omitempty"`
	Body  string         `json:"body,omitempty"`
	Meta  map[string]any `json:"meta"`
}
