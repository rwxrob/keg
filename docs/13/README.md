# Add default columns when cannot be determined

Support for the `columns` variable allows users to defined their preferred default column width when it cannot be determined, for example, when not using an interactive terminal such as using `vim` filters with `!!` commands (ex: `!!keg foo`). Usually a `keg` user will have a regular column width they prefer so this allows things to not wrap as strange widths.

The default `keg.DefColumns` is set to 100.

Remember that neither of these will be used if using a terminal interactively (stdout is to a terminal).
