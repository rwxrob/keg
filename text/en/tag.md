tag one or more nodes

The {{aka}} command tags one or more content nodes by adding them to the `dex/tags` file. Each line of the file begins with a tag (which can be anything that does not contain an ASCII space, even though sensible, social-media compatible tags are strongly recommended).

Note that the KEG Specification strongly suggests against meta data for individual content nodes arguing that meta data approaches such as "front matter" have proven to be failures by corrupting the actual content with meta content. With KEG the content *is* the meta content by the nature of the semantic syntax required by KEGML. This promotes creation of meta data collections (such as indexes) instead.
