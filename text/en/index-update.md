update dex files

The {{aka}} command forces a rescan and update of the current files in the `dex` index directory. Normally, these files are updated every time any command is executed successfully that changes the state of the keg itself. But, sometimes things might get out of sync, say after editing directories or files directly without using this command. In such cases running the {{aka}} command is needed.

While (re)making the index files, this command ensures that any "empty" content nodes are removed. An empty node is one that recursively contains no file of any length greater than zero. This means that a content author can effectively force the deletion of a content node just by zeroing out the `README.md` file during an editing session and saving it (in most cases).
