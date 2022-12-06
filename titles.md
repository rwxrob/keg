The {{aka}} command returns a paged list of all titles matching the given regular expression. By default all regular expressions ({{pre "REGEXP"}}) are made case insensitive by adding the prefix {{pre "(?i)"}} which can be explicitly overridden with {{pre "(?-i)"}} for one search or changed as the default by assigning it to the {{pre "regxpre"}} variable:

    keg set regxpre '(?-i)'

Note that if set, {{pre "regxpre"}} applies to *all* searches, which includes the {{cmd "edit"}} and {{cmd "grep"}} commands.
