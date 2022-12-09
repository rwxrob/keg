The {{aka}} command opens a content node README.md file for editing. It is the default command when no other arguments match other commands. Nodes can be identified by integer ID, REGEXP matching the title, or the special `last` (last created) or `same` (last updated) parameters. For REGEXP if more than one match is found the user is prompted to choose between them.

The editor opened depends on the following in order of priority:

1. `VISUAL` environment variable
2. `EDITOR` environment variable
3. `code` command
4. `vim` command
5. `vi` command
7. `emacs` command
7. `nano` command
