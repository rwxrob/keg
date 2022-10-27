# Knowledge Exchange Grid (KEG) Version 1.0

## Terminology

* *node*
* *directory*
* *site*
* *slug*
* *isosec*

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
strictly in KEG Mark, a simplification of 2023 Pandoc and
GitHub-Flavored-Markdown with IETF RFC-like limitations by design (72
column width) so that "the source itself is always as readable as how it
is rendered" (Gruber).

A *node* MUST contain a `README.md` file.

The *node* `README.md` file MUST be valid KEG Mark to be considered a
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
