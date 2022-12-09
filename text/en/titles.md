find titles containing regular expressions

The {{aka}} command returns a paged list of all titles matching the given regular expression. By default all regular expressions (REGEXP) are made case insensitive by adding the prefix `(?i)` which can be explicitly overridden with `(?-i)` for one search or changed as the default by assigning it to the `regxpre` variable:

    keg set regxpre '(?-i)'

Note that if set, `regxpre` applies to *all* searches, which includes the {{cmd "edit"}} and {{cmd "grep"}} commands.
