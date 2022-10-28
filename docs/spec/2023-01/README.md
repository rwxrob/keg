# Knowledge Exchange Grid (KEG) Version 1.0

## Terminology

A ***node*** (or "knowledge node") is simply a directory with a
`README.md` or `DATA.*` file in it.

A ***dynamic node*** (or "generated node") is a *node* that requires a
generator be run to produce the `README.md` or `DATA.*` files.

A * and any number of optional source
files from this these other two main files are generated. This
simplicity is by design. A *node* is
considered *static* if it does not have a generator. If it does, it's
considered *dynamic*.

***directory***
***site***
***slug***
***isosec***

A ***text node*** contains only KEGMark text.

A ***collection node*** is a *text node* that includes one or more
*bulleted import lists*.

A ***composite node*** is a *dynamic node* that is generated from other
nodes. A book is an example of a *composite node*.

A ***bulleted import list*** is a KEGMark block type consisting of a
regular bulleted list with nothing but links for items in the list. When
browsing from are rendered Web version or GitHub clicking on any of the
links in the list will open the knowledge source file. But, such lists
can also easily be expanded into larger *composite nodes* by simply
replacing the bullet list item with the content of the link target. Any
headers in the imported content will be incremented by one (prefixed
with a single, additional hashtag `#`).

## Design Principles

> Don't get mad, get busy.

The modern Web makes us mad, real mad. Even Tim Berners-Lee recognized
the significant flaws, which, at this point, can never be fixed without
throwing the whole thing out and starting over. So we did.

The Knowledge Exchange Grid, like Markdown which inspired it, is focused
entirely on the simplest possible method to capture and convey
written knowledge, you know, kinda like the original WorldWideWeb at
CERN, but without all the over-engineering you would expect from a group
of excited, smart, CERN scientists. (Sorry, HTML was a horrible idea
from the very beginning.)

Unlike the modern Web, KEG caters primarily to the needs of *writers* who
wish to control not only their own authenticated content, but also their
view of those they specifically wish to follow, and --- most importantly
--- their *own* search parameters based on predictable semantic content
organization. This is far different from the modern Web where no one can
know for sure who wrote what and "SEO" is dictated by those with the
most money and power to promote their own agenda. KEG promotes informed
written dialog and debate (not advertising, trolling, and violent
disagreement).

KEG Mark, the syntax of knowledge source on KEG, simplifies Markdown
even further. Most can learn it over a single cup of coffee. No tools
other than a means of writing and storage are required. Even paper and
pencil fulfill this need. KEG content creation is, therefore, highly
compatible with most note-taking techniques --- especially the
Zettelkasten method. Need complicated Math notation? The international
MathJax standard is supported making KEG 100% accessible to the blind
community as well. Images are discouraged where a description would
suffice and adding an image in KEG Mark without a full description is
considered invalid, allowing KEG users to blacklist such content.

KEG users can optionally follow others creating local cached copies of
their content available even when the network is inaccessible.

KEG users can also optionally provide a list of those they follow (and
blacklist) enabling rings of social trust to police KEG content without
any heavy-handed moderation needed by any single entity. In fact, it's
impossible to censor content on KEG since KEG does not depend any
specific technology at all. Carrying a KEG collection on a USB stick to
another is a perfectly viable form of KEG exchange --- by design. The
secret *is* the simplicity.

## Node ("KEG node", "knowledge node")

A KEG *node* is simply a directory with a `README.md` file written
strictly in KEG Markdown, a simplification of 2023 Pandoc and
GitHub-Flavored-Markdown with IETF RFC-like limitations by design (72
column width) so that "the source itself is always as readable as how it
is rendered" (Gruber).

A *node* MUST contain a `README.md` file.

The *node* `README.md` file MUST be valid KEG Markdown to be considered a
*node* at all. Invalid `README.md` files MUST disqualify the *node* from
inclusion in any KEG content or indexing, no exceptions, ever.

A *node* directory (not to be confused with a KEG *directory*) MAY
contain any other files.

Files within a *node* MUST be considered a part of the *node* as well.

A *node* MUST NOT contain other nodes.

A *node* MUST belong (be contained within) a *directory*.

A *node* MAY belong to several *directories* through hard or soft
linking.

## Directory ("keg")

A KEG *directory* (also colloquially called a "keg") MUST contain a
collection of one or more *nodes*, each with a unique identifier. 

Any unique identifier type is allowed but the following types are
RECOMMENDED since any human can easily create them:

* slug - title predictably converted
* isosec - node title is uncoupled from name

While slugs are easier to navigate, they can create problems when
interlinking is allowed or encouraged. If the title changes, so must the
slug. This has long been a problem with Web content organization.

A KEG *directory* itself fulfills the requirements for a KEG *node*
except for the obvious exception that it MUST contain other nodes (where
normally a *node* may not).

A KEG *directory* MUST contain a `KEG.yaml` file in addition to the 

## ISO Second ("isosec")

An *isosec* is the GMT current time in ISO8601 (RFC3339) without any
punctuation or the T. This unique identifier is preferred over others
since it can easily be determined with nothing more than a watch or wall
clock with a second hand.

## KEG Mark

KEG Mark is a simplified (and deliberately limited) version of Pandoc
Markdown with standardized extensions for Math notation (MathJax)
fenced, semantic divisions, tables, and bibliographic references. It is
suitable for conversion into most rich publications formats including
academic papers, books, novels, articles, and blog posts.

Although tables are supported, it is strongly RECOMMENDED that *data
nodes* be used instead where possible since the data is immediately
usable in such form. Content creators may opt to use a generator for the
tables in their `README.md` files and store the actual data in the
`DATA.*` file instead. It is also RECOMMENDED to limit tables to one per
*text node*.

## Dynamic Node

Dynamic nodes must be periodically generated. The method of generation
must be available to anyone viewing the node. This method is contained
within the `gen.*` file. Dependencies on public or private libraries and
data is fine. The goal is to inform followers of the method in general.

:::note
The inclusion of the `gen.*` generator is *only* to give insight
into the origin of the node content, not to allow anyone following the
dynamic node to generate their own. Content creators are under no
obligation to include generators that would actually work for anyone
but themselves during the local dynamic node generation process. Some
creators, however, may elect to make sure that others can generate
their own, for example, when *cloning* a dynamic node.
:::

The decision about if and how to communicate the last updated date and
time are completely up to the node content creator, but are usually
included as field within the `DATA.*` file or bit of text in the
`README.md` file.

Whatever process executes the generator MUST update the KEG Directory
`MANIFEST` file as well. Dynamic nodes are marked as "dynamic" ('D')  in
the `MANIFEST` file so that those following can filter out the regularly
updated changes according to their own interests.

Generators MUST be readable files included within the node directory and
begin with `gen` followed by period and suffix indicating the language.
The following suffixes are RECOMMENDED, but not required:

* `sh` - POSIX shell (ash/dash compatible)
* `bash` - Bash 4.*
* `go` - Go 1.18+
* `perl` - Perl ??
* `rb` - Ruby ??
* `py` - Python 3+

Note that suffixes (albeit annoying to some) are required because the
host computer cannot be guaranteed to understand other methods of script
interpreter delegation (i.e. she-bang lines).

Note that generators MUST NEVER be compiled code and MUST NOT contain
secrets of any kind, which will generally be stored in other ways and
pulled in during the generation process.
