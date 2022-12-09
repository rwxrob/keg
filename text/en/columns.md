print the number of columns resolved

The {{aka}} command first looks for `term.WinSize.Col` which is set by many UNIX-like operating systems. If not found, the columns variable (from `vars`) is checked and used if found. Finally, the default hard-coded value ({{columns}}) is used.
