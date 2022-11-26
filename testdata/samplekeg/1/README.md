# Sample content node

Hello KEG friend!

This sample **content node** document serves as a KEGML[^kegml] cheat-sheet and primer when creating new nodes to help remember how. It also includes reminders about best practices that aren't necessarily required by the specification[^spec]. This document is intended to be easier to read than the specification and should take about 30 minutes.

Rather than describe what KEGML is *not* (for those that already know Markdown) let's specifically describe what it is.

***A KEG content node is the fundamental unit of KEG content written in KEGML.*** A node is a directory inside a **keg directory** with an integer name[^nodeid] that contains (1) a `README.md` file; (2) an optional `meta` YAML file with structured data; and (3) one or more optional content files (including images). Within the `README.md` file (like this one) other nodes and files are linked and included producing a **knowledge graph** structure that can be rendered as a document, static Web site, outline, or mind map. This design is inspired by the original WorldWideWeb and Luhmann's very successful Zettelkasten method[^luhmann].

***KEGML is Markdown with limitations.*** If you already know any flavor of Markdown[^md] then you already know KEGML. Only the terms and identifiers have changed to provide more semantic meaning.  KEGML limitations are by design to promote good content that is easy to create, read, search, and maintain. KEGML is 100% compatible with all major Markdown versions including CommonMark, GitHub Flavored Markdown, Myst, and Pandoc Markdown. Compatibility means KEGML content can be hosting most anywhere without fear of later compatibility issues.

***A KEGML document is composed of blocks separated by a single blank line.*** Each block has a specific type and syntax and begins with a unique token. (Those developing parsers will be keenly interested in these tokens.) Blocks that allow line returns also have end tokens (that always match the beginning token). All blocks must end with a blank line (effectively two line returns `\n\n`, no carriage returns) or the end of the document. Blocks may contain **spans** and **links** (defined by KEGML) or just raw **runes**[^unicode] (undefined). Some blocks are limited in their placement in relation to one other and the body as a whole. Here is a summary of block types, their token identifiers, what they are allowed to contain, whether or not the block allows line returns within it, and a short description:

|Block          | Token       |  Lines | Allowed | Description
|              -|-            |    -   |   -   |
| Title         | `# `        | No     | Inflect, Math, Code | Must be first line, 72 runes max
| Bulleted List | `* ` `- ` `+ ` | No  | All but Lede | One line/item, never follows another list
| Numbered List | `1. `       | No     | All but Lede | One line/item, never follows another list
| Include List  | `* [`       | No     | Inflect, Math, Code | Node/file target, never follows list
| Footnotes     | `[^`        | No     | Inflect, Math, Code, URL, Link | Citations, comments, etc.
| Fenced        | `` ``` `` `~~~`  | Yes | Runes | 3+ backticks or squiggles, syntax class
| Quote         | `> `        | No     | All but URL, Link, Lede | Single paragraph of quotation
| Math          | `$$`       | Yes     | Runes | MathJax, LaTeX math notation
| Figure        | `![`       | No      | Inflect, Math, Code | Meaningful images and graphics
| Separator     | `----`     | No      | None | Subdivisions, not always "horizontal rule"
| Table         | `|`        | No      | All but Lede | Compatible with GFM
| Paragraph     | None       | No      | All but URL | Usually the most common

***Spans are composed of an inline run of runes and begin and end with a token.*** All spans begin and end with the same token. Begin tokens must be preceded by a space (or begin block) and must not be followed by a space (ex: ` **foo`). End tokens must be followed by a space (or end block) and not be preceded by space (ex: `foo** `). The special double-backtick Code span (which exists primarily to allow backticks themselves to be included in Code spans) is the only exception and must be surrounded by spaces. Spans must never contain line returns.

|  Span   | Tokens    | Description |
|    -    |   -       |     -       |
| Inflect | `*`       | Alter voice, tone, or mood
| Beacon  | `**`      | Draw attention to a term or phrase without inflection
| Lede    | `***`     | Introductory text meant to summarize, provoke, and entice
| Math    | `$`       | Inline MathJax markup notation
| Code    | `` ` ``   | Code and content other than terms and words, monospace, preformatted
| Code    | `` `` ``  | (special)
| URL     | `<` `>`   | Web and other universal resource locators
| Plain   | (none)    | Anything not in another span type

***Links have two parts: link text, and a link target.*** The **link text** begins with a left bracket (`[`) and ends with a right bracket (`]`). The text itself is composed of Inflect, Beacon, Math, Code or Plain spans. The link text is immediately followed by the **link target** which begins with left parenthesis (`(`) and ends with right parenthesis (`)`) and contains either a **node link target** (beginning with `../`) or a **file link target** (matching name of file in local directory). Both link target types may have a **query code** suffix that begins with a question mark (`?`).


***Inflect, Beacon, Lede spans may included Math and Code spans.*** They may now, however, include themselves. This prevents the common stylistic problems but more importantly ensures that these spans maintain their semantic meaning over any possible stylistic formatting.

***First line is the title block.*** The **title** must begin with a single hashtag (#) followed by a single space and be less than 72 total **runes**[^unicode]. No other headings are allowed in a KEGML document. Titles should usually be in sentence case for easy reading and writing depending on the topic. The text of a title may contain one or more Inflect, Math, or Code spans. Span tokens count against the 72 total max. Plain text should usually be preferred.

***Paragraphs are composed of one or more spans.***

***Bulleted lists begin with star (*), dash (-), or plus (+).***

***Numbered lists begin with number one (1) and a dot (.).***

***Include list is a bulleted list but containing node and file links to be included.***

***Footnotes begin with left square bracket ([) and a caret (^).***

***Math block begins and ends with two dollar signs ($).***

***Lede span begins and ends with three stars and must be first span.***

***One list block after another is not allowed.*** Because of the confusion associated with Markdown "long" lists, KEGML does not allow two list blocks to appear in the body next to one another. Separate lists with any other block type to ensure your documents remain valid KEGML.

***Consider using beacon for terms that first appear.*** A **beacon** is specifically designed to draw attention without change in voice or tone to a span of text primarily for the purpose of introducing new terms and language. By following a simple convention of making the first occurrence of a term in a document into a beacon you can automate linking to a node with that same case-insensitive title. Editors can be customized to add the link automatically, or renderers can create the links when generating content in other forms.

[^spec]: <https://github.com/rwxrob/keg-spec>
[^md]: "Writing mathematical expressions". *GitHub Docs*. https://docs.github.com/en/get-started/writing-on-github/working-with-advanced-formatting/writing-mathematical-expressions
[^luhmann]: https://luhmann.surge.sh
[^unicode]: A Unicode code point can be from one to four bytes long and allows the use of emojis and character sets for all of the world's languages. The Go programming language (which was created by the main creators of Unicode itself, Rob Pike) was the first to call them "runes" for short. In contrast, the term "char" retains its one-byte meaning.
[^kegml]: Knowledge Exchange Graph Markup Language
[^nodeid]: All node IDs must be integers. However, an **index** qualifies as being a node even though it has a non-integer ID. This is to prevent indexes from being indexed themselves. But for the purposes of linking, an index *is* a node and therefore a node link target may include a non-integer after its identifying prefix (ex: `../2` or `../dex`).
