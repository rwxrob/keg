# KEG Quick Start Guide

***Take control of your personal knowledge management quickly by starting
your own knowledge exchange graph with the `keg`[^kn] command.***

Hello friend. Welcome to the KEG community!

We are obsessed with lightning-fast knowledge management from the terminal and created the `keg` (aka `kn`) command to help. We hope it will be the first fully compliant implementation of the KEG Specification[^spec]. This guide is designed to help you learn to create KEGML[^kegml] content quickly so you can learn the rest by using it.

* [What is knowledge managment any why should you care?](../1?L)
* [How does `keg` compare to other knowlege management apps?](../3?L)
* [A KEG content node is the fundamental unit of KEG content written in KEGML.](../2?L)
* [KEGML is Markdown with limitations.](../4?L)


***KEGML documents must be encoded in UTF-8 and contain only printable[^printable] runes[^unicode] and line returns (`\n`).*** Tabs are not allowed. Carriage returns are not allowed. Terminal ANSI escapes are not allowed. Line returns are only allowed in Fenced and Math blocks. Emojis are allowed and encouraged.

***A KEGML `README.md` is composed of blocks separated by a single blank line.*** Each block has a specific type and syntax and begins with a unique token. Those developing parsers will be keenly interested in these identifying tokens. Blocks that allow line returns also have ending tokens that always match the beginning token. Except for the last block, all blocks must end with a blank line (two line returns `\n\n`). Depending on the block type a block may contain spans, links, or raw runes[^unicode].

|Block          | Token            | Contains
|              -|-                 |   -
| Title         | `# `             | Inflect, Math, Code
| Bulleted List | `* ` `- ` `+ `   | All but Lede
| Numbered List | `1. `            | All but Lede
| Include List  | `* [`            | Inflect, Math, Code
| Footnotes     | `[^`             | All buf Lede
| Fenced        | `` ``` `` `~~~`  | Runes
| Quote         | `> `             | All but URL, Link, Lede
| Math          | `$$`             | Runes
| Figure        | `![`             | Inflect, Math, Code
| Separator     | `----`           | None
| Table         | `|`              | All but Lede
| Paragraph     | None             | All but URL

Take a moment to identify these different block types in this document. A full description and example of each is included below after the summaries of spans and links.

***Some blocks are limited in their placement in relation to one other and the body as a whole.*** These limitations provide stronger semantic meaning and avoid well-known editing, parsing, and rendering pitfalls.

* Only a single Title or Footnotes block is allowed
* Title must be first line and not exceed 72 total runes
* Footnotes must be the last block
* Lists must never follow other Lists of any type
* Math must never follow another Math block
* Separator must never follow another Separator block

***Spans are composed of printable[^printable] runes[^unicode] and begin and end with a token.*** Beginning tokens must be preceded by a space and must not be followed by a space (ex: ` **foo`). Ending tokens must be followed by a space and not be preceded by space (ex: `foo** `). Exceptions apply when span is at the beginning or ending of a block. The special double-backtick Code span exists primarily to allow backticks themselves to be included in Code spans and is the only exception and must be surrounded by spaces.

|  Span       | Tokens    | Description                               |
|    -        |   -       |     -                                     |
| Inflect[^i] | `*`       | Alter voice, tone, or mood                |
| Beacon[^b]  | `**`      | Draw attention, terminology, phrases      |
| Lede        | `***`     | Introductory, summarize, provoke, entice  |
| Math        | `$`       | Inline MathJax/LaTeX markup notation      |
| Code        | `` ` ``   | Code, monospace, preformatted             |
| Code        | `` `` ``  | (same as Code, but allows backtick)       |
| URL         | `<` `>`   | Universal resource locator                |
| Plain       | (none)    | Anything not in another span type         |

Some Spans are limited in where they may appear:

* Lede must be first (and possibly only) span in paragraph block

Take a moment to identify use of each span type within this document. Full descriptions and examples of each span type are included after the description of block types below.

***There are three types of links: node, file, and foot.*** Node links target another nodes including indexes[^dexnode]. File links target files in the local node directory. Foot links target a footnote in the same `README.md` file containing the foot link.

***Node and file links have two parts: link text, and a link target.*** The **link text** begins with a left bracket (`[`) and ends with a right bracket (`]`). The text itself is composed of Inflect, Beacon, Math, Code or Plain spans. The link text is immediately followed by the **link target** which begins with left parenthesis (`(`) and ends with right parenthesis (`)`) and contains either a **node link target** (beginning with `../`) or a **file link target** (matching name of file in local directory). Both link target types may have a **query code** suffix that begins with a question mark (`?`).

***Inflect, Beacon, Lede spans may included Math and Code spans.*** They may now, however, include themselves. This prevents the common stylistic problems but more importantly ensures that these spans maintain their semantic meaning over any possible stylistic formatting.

***First line is the title block.*** The **title** must begin with a single hashtag (#) followed by a single space and be less than 72 total **runes**[^unicode]. No other headings are allowed in a KEGML document. Titles should usually be in sentence case for easy reading and writing depending on the topic. The text of a title may contain one or more Inflect, Math, or Code spans. Span tokens count against the 72 total max. Plain text should usually be preferred.

***Paragraphs are composed of one or more spans.***

***Bulleted lists begin with star (*), dash (-), or plus (+).***

***Numbered lists begin with number one (1) and a dot (.).***

***Include list is a bulleted list but containing node and file links to be included.***

***Footnotes begin with left square bracket ([) and a caret (^).***

***Math block begins and ends with two dollar signs ($).***

***Lede span begins and ends with three stars and must be first span.***

***Consider using beacon for terms that first appear.*** A **beacon** is specifically designed to draw attention without change in voice or tone to a span of text primarily for the purpose of introducing new terms and language. By following a simple convention of making the first occurrence of a term in a document into a beacon you can automate linking to a node with that same case-insensitive title. Editors can be customized to add the link automatically, or renderers can create the links when generating content in other forms.

***URL in paragraph is just considered text.*** URLs are notoriously bad at destroying the readability of any text --- especially paragraphs. Therefore, if they are included in a paragraph they must be considered just like any other text. If and when a Web address is required consider strongly just using the domain name without the preceding `http` schema portion or create a list item containing the URL in long form surrounded by angle brackets so that it is visible even if the document were printed. Best of all, create a proper bibliographic citation in a footnote.

[^spec]: <https://github.com/rwxrob/keg-spec>
[^md]: "Writing mathematical expressions". *GitHub Docs*. https://docs.github.com/en/get-started/writing-on-github/working-with-advanced-formatting/writing-mathematical-expressions
[^luhmann]: https://luhmann.surge.sh
[^unicode]: A Unicode code point can be from one to four bytes long and allows the use of emojis and character sets for all of the world's languages. The Go programming language (which was created by the main creators of Unicode itself, Rob Pike) was the first to call them "runes" for short. In contrast, the term "char" retains its one-byte meaning.
[^kegml]: Knowledge Exchange Graph Markup Language
[^nodeid]: All node IDs must be integers. However, an **index** qualifies as being a node even though it has a non-integer ID. This is to prevent indexes from being indexed themselves. But for the purposes of linking, an index *is* a node and therefore a node link target may include a non-integer after its identifying prefix (ex: `../2` or `../dex`).
[^printable]: Defined by Unicode to be letters, numbers, punctuation, symbols, and ASCII space.
[^i]: Reuse of the "i" from "italic" but defined semantically same as HTML5.
[^b]: Reuse of the "b" from "bold" but defined semantically same as HTML5.
[^dexnode]: An index directory *is* a node but without the integer name making them **unindexed nodes**.
[^sample]: This document is actually a lot longer than a **content node** would normally be. Usually, most of this document would be broken up into smaller nodes and aggregated back together again with **includes**.
[^kn]: The `keg` command is frequently renamed to `kn` for easier typing.
