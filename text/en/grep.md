grep regular expression out of all nodes

The {{aka}} command performs a simple regular expression search of all node `README.md` files. It does not depend on any {{cmd "grep"}} command being installed.

By default, all regular expressions are case-sensitive (unlike default {{cmd "title"}} or {{cmd "edit"}} commands). This can be explicitly overridden for a given search by adding `(?i)` or for all searches by setting the global `regxpre` variable to the same:

    keg grep set regxpre '(?i)'

Note that if set, `regxpre` applies to *all* searches, which includes the {{cmd "titles"}} and {{cmd "edit"}} commands.

When run interactively the choice of hits is presented so that the user may select. The selection is then delegated to the {{cmd "edit"}} command.

When run non-interactively the matching files are listed with their titles as a node include list.
