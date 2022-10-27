package mark

// Combine will recursively combine all README.md files starting with
// the target directory into a single string with all headers adjusted
// based on the composition level used according to the following
// algorithm:
//
// Does the README.md file within the Knowledge Node contain only
// a single title header (always on the first line)?
//
// If not, include as is.
//
// If so, for each composition block detected (a simple bulleted block
// consisting of nothing but Markdown links) fetch the README.md file
// targeted and recursively call Combine on it incrementing the level
// each time.
//
// Combine will fail if the level passed is greater than five (since
// only six HTML header tags are supported). This requires some content
// organization in advance to meet these limitations, which are by
// design, content nested more than six levels usually needs a separate
// document all together.
//
// Only standard rooted Node references are supported as targets
// (`/some-thing`). Consider importing a remote node (with full
// attribution) when remote content is desired in the Combine
// composition.
//
// No local/anchor link collision checking is done and it is considered
// bad practice to add a local link to a Node that is planned to be
// composed into another where a local link would resolve. In other
// words, all Nodes should fully resolve before composition.
//
func Combine(target string, level int) string {
	// TODO
	return ""
}
