show the current keg

The {{aka}} command displays the current keg by name, which is resolved as follows:

1.  `KEG_CURRENT` environment variable
2.  `current` var setting
3.  Current directory if `keg` file
4.  `docs` directory if `keg` file

The `docs` directory is included since it is such a common convention for combining documentation with any GitHub repo using GitHub Pages or other static site publishing methods.

***keg [var] set current***

Set the `current` var when you know that you will be moving around a lot from multiple windows and just want to consistently use the same keg for a given session. That way even if you are in the directory of a project with its own `docs/keg` it will not interfere. Then, after your editing session is done, {{aka}} `unset current` (or set it to something else).

***KEG_CURRENT***

The `KEG_CURRENT` environment variable has priority over everything else allowing it to be set on the same line as the call to the {{aka}} command for one-off occasional commands.

     KEG_CURRENT=zet keg titles

(Note that `keg set current keg` would be preferred in this case for multiple commands.)

The use of `KEG_CURRENT` in a small script is recommended to provide a permanent alternative command for that specific keg.

     #!/bin/sh
     KEG_CURRENT=~/Repos/github.com/rwxrob/zet/docs keg "$@"

Scripts are always preferred over aliases or shell functions since they work in *every* context including being called from within editors (common for including keg query output directly into other keg content nodes while editing) or from other programs using {{cmd "exec"}}.

The `KEG_CURRENT` environment variable can be one of the three following:

1. Fully qualified path to keg directory beginning with path separator (`{{pathsep}}`)
2. Relative home directory beginning with tilde (`~`)
3. Name of {{cmd "conf"}} map key with value pointing to keg directory

When using option #3 the directory can either be one that contains a `keg` file or one that contains a `docs` directory with a `keg` file.

    map:
      zet: ~/Repos/github.com/rwxrob/zet
      keg: ~/Repos/github.com/rwxrob/keg

