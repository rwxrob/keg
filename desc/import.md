The {{aka}} command imports a specific NODEDIR or all the apparent node directories within DIR into the current node. If no argument is passed, imports the current working directory into the current keg. If any of the arguments end in an integer they are assumed to be node directories. Arguments without a base integer are assumed to be directories containing node directories with integer identifiers.

This command is useful when indirectly migrating nodes from one keg into another by way of an intermediary directory (like `tmp`)

Currently, there is no resolution of links within any of the imported nodes. Each node to be imported should be checked individually to ensure that any dependencies are met or adjusted.
