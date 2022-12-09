show info about last created node

The {{aka}} command shows information about the last content node that was created, which is assumed to be the one with the highest integer identifier within the current keg directory. By default the colorized form is displayed to interactive terminals and a KEGML include link when non-interactive (assuming `!!` from vim, for example).

* `dir` shows only the full directory path
* `id` shows only the node ID
* `title` shows only the title
* `time` shows only the time of last change

Note that this is different than the latest {{cmd "changes"}} command.
