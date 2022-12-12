add node to the tags index

The {{aka}} command takes a comma separated list of `TAGS` and single content node to add to the `dex/tags` file for each tag. The node can be specified in the usual ways:

* `same` - last changed node
* `last` - last created node
* NODEID - integer identifier
* REGEXP - regular expression matching title (interactive select if >1 hit)

If the content node parameter is omitted, returns the lines from `dex/tags` for the specified `TAGS`.

The special reserved tag `all` prints everything in the `dex/tags` file. If no arguments are passed, `all` is assumed.

Each line of the `dex/tags` file begins with a tag (which can be anything that does not contain an ASCII space, even though sensible, social-media compatible tags are strongly recommended). Even if there are not node ids on a given line, the tag must be immediately followed by a single space.

Note that the KEG Specification strongly suggests against meta data for individual content nodes arguing that meta data approaches such as "front matter" have proven to be failures by corrupting the actual content with meta content. With KEG the content *is* the meta content by the nature of the semantic syntax required by KEGML. This promotes creation of meta data collections (such as indexes) instead.
