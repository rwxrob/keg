The {{aka}} command displays the current keg by name, which is resolved as follows:

1.  The `KEG_CURRENT` environment variable
2.  The current working directory if `keg` file found
3.  The `docs` directory in current working if found
4.  The `current` var setting (see {{cmd “var”}})

Note that setting the var forces {{aka}} to always use that setting until it is explicitly changed or temporarily overridden with `KEG_CURRENT` environment variable.

     keg()
     {
       KEG_CURRENT=zet keg "$@"
     }

It is often useful to have `current` set to the most frequently used keg and then change into the working directory of another, less updated, keg when needed.
