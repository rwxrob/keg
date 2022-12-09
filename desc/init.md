The {{aka}} command creates a `keg` YAML file in the current working directory and opens it up for editing.

{{aka}} also creates a **zero node** (`../0`) typically used for linking to planned content from other content nodes.

Finally, {{aka}} creates the `dex/changes.md` and `dex/nodes.tsv` index files and updates the `keg` file `updated` field to match the latest update (effectively the same as calling {{cmd "dex update"}}).

Also see the `Getting Started` docs in the main {{cmd "help"}} command.
