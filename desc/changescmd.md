The {{aka}} command displays COUNT content nodes in reverse chronological order that were most recently updated. If no COUNT is specified, the number displayed is {{changesdef}} by default, but this can be changed by setting the {{pre "changes"}} variable to something else:

    keg changes set default 10

When interactive, output is colored and sent to pager if detected.

When not interactive, renders as plain text KEGML include block with node links.

Note that this is different than the {{cmd "last"}} command.
