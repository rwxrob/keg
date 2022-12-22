# Sample content node

Hello KEG friend! For full reference try KEG Quick Start Guide[^start] or official KEG Specification[^spec]. If you already know some form of Markdown this sample summary should be enough to help you understand KEGML constraints to regular Markdown for the sake of simplicity, clarity, and efficiency.

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

* Only a single Title or Footnotes block is allowed
* Title must be first line and not exceed 72 total runes[^unicode]
* Footnotes must be the last block
* Lists must never follow other Lists of any type
* Separator must never follow another Separator block

|  Span       | Tokens    | Description                               |
|    -        |   -       |     -                                     |
| Inflect[^i] | `*`       | Alter voice, tone, or mood                |
| Beacon[^b]  | `**`      | Draw attention, terminology, phrases      |
| Lede        | `***`     | Introductory, summarize, provoke, entice  |
| Math        | `$`       | Inline MathJax/LaTeX markup notation      |
| Code        | `` ` ``   | Code, monospace, preformatted             |
| URL         | `<` `>`   | Universal resource locator                |
| Plain       | (none)    | Anything not in another span type         |

* Lede must be first (and possibly only) span in paragraph block

***There are three types of links: node, file, and footnote.*** Node links target another nodes including index nodes[^dexnode]. Both node and file nodes can have a **query code** that begins with question mark (`?`) and specifies how to include and expand the linked node or file.

|Query Code | Behavior
|          -|-
| (none)    | Link text becomes relative heading (beginning with hashtags `#`)
| T         | Target title becomes relative heading
| L         | Link text becomes lede
| 0         | Just include target body


```md
This is [node link to zero node](../0) that always has `../` in front. If linking to [a file](somefile) must be local[^1].

* [Include node 1 as a heading (this text)](../1)
* [Include node 2 as a heading (target title)](../2?T)
* [Include node 3 as a lede paragraph](../3?L)
* [Include node 4 with no heading](../4?0)

[^1]: This footnote explains that local is same directory as README.md (no slash).
```

[^start]: <https://rwxrob.github.com/keg>
[^spec]: <https://github.com/rwxrob/keg-spec>
[^unicode]: A Unicode code point can be from one to four bytes long and allows the use of emojis and character sets for all of the world's languages. The Go programming language (which was created by the main creators of Unicode itself, Rob Pike) was the first to call them "runes" for short. In contrast, the term "char" retains its one-byte meaning.
[^nodeid]: All node IDs must be integers. However, an **index** qualifies as being a node even though it has a non-integer ID. This is to prevent indexes from being indexed themselves. But for the purposes of linking, an index *is* a node and therefore a node link target may include a non-integer after its identifying prefix (ex: `../2` or `../dex`).
[^i]: Reuse of the "i" from "italic" but defined semantically same as HTML5.
[^b]: Reuse of the "b" from "bold" but defined semantically same as HTML5.
[^dexnode]: An index directory *is* a node but without the integer name making them **unindexed nodes**.
