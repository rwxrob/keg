work with index (`dex`) files

The {{aka}} command is a command branch containing commands related to indexing a keg in the standard and customized ways. Most kegs have a `dex` directory in which the following files are kept up to date every time there is a change to any keg node:

* `dex/changes.md` - last changes in reverse chronological order in markdown
* `dex/nodes.tsv` - all nodes in tab-separated format ordered by integer id

These files are updated every time any command is executed successfully that changes the state of the keg itself.
