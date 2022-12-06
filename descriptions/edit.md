The {{aka}} command opens a content node README.md file for editing. It is the default command when no other arguments match other commands. Nodes can be identified by integer ID, TITLEWORD contained in the title, or the special {{pre "last"}} (last created) or {{pre "same"}} (last updated) parameters.

For TITLEWORD if more than one match is found the user is prompted to choose between them. Otherwise, the match is opened in the EDITOR. See rwxrob/fs.file.Edit for more about how editor is resolved.
