# More reason to use commands over functions and aliases

So I have been using the following to get encapsulated support for assigning the `keg` command to a specific path using one of the many ways this can be done. Aliases worked, but failed from within Vim because of all the normal reasons. But I was surprised to discover that so do functions. I'm not sure if it has something to do with the environment variable or the command itself. The solution was to create the following two-line shell script with the setting encapsulated in it.

```
#!/bin/bash
KEG_CURRENT=zet keg "$@"
```

I got the hard link thing to work, but the requirement to map that name to a directory path required putting that into a per `zet` cached var so the benefit seems less than preferable. This means the easiest option for multiple standalone commands is still a simple command.
