delete node from current keg

The {{aka}} command will delete a specified content node from the current keg. A node may be specified in any of the following ways:

1. By integer node identifier
2. `same` indicating most recently changed node
3. `last`  indicating most recently created node

In addition to deleting the content node directory and everything within it recursively the node entry is removed from the current index files within `dex` and the entire keg is published with these changes.

If the specified content node does not exist the command is ignored.
