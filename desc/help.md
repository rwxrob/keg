The {{aka}} command is for personal and public knowledge management as a Knowledge Exchange Graph (sometimes called "personal knowledge graph" or "zettelkasten"). Using {{aka}} you can create, update, search, and organize everything that passes through your brain that you may want to recall later, for whatever reason: school, training, team knowledge, or publishing a paper, article, blog, or book.

***Getting Started***

The steps to create your first KEG directory are below, but first a little about the structure of this directory, which does not necessarily require this {{aka}} command to create and maintain.

A KEG directory (aka "keg") is just a directory containing a `keg` YAML file and a number of directories that have a `README.md` file called **content nodes**. A node directory must have an incrementing integer name as would be used in a database table. A keg usually also has a `dex` directory containing at least two files:

1. Latest changes `dex/changes.md`
2. All nodes by ID `dex/nodes.tsv`

The {{aka}} command keeps these files up to date.

A special **zero node** is used by convention as a target for links to nodes that have yet to be created.

Okay, here are the specific steps to get started by creating your first keg directory. If you plan on using Git or GitHub hold off on doing anything with git for now.

1. Create a directory and change into it
2. Run the {{aka}} {{cmd "init"}} command
3. Update the YAML file it opens
4. Exit your editor
5. List contents of directory to see what was created
6. Run the {{aka}} {{cmd "create sample"}} command to create your first node
7. Read and understand the sample
8. Exit your editor
9. Check your index with {{aka}} {{cmd "changes"}} or {{aka}} {{cmd "titles"}}
10. Repeat 6-9 creating several nodes (optionally omitting {{cmd "sample"}})
11. Search titles with the {{aka}} {{cmd "titles"}} command
12. Edit node with title keywords with {{aka}} {{cmd "edit WORD"}} command
13. Edit node with grep regexp matches with {{aka}} {{cmd "grep WORD"}} command
14. Notice that {{aka}} {{cmd "edit"}} is the default (ex: {{aka}} WORD)

***Git and GitHub***

It's important when using Git that either the remote git repo has been fully created (so that {{cmd "git pull"}} will work) or that {{cmd "git"}} has not been run at all (no `.git`}} directory). Otherwise, {{aka}} will attempt to pull and fail. These instructions assume the reader understands {{cmd "git"}} and the {{cmd "gh"}} commands.

Here are the steps to follow when Git and GitHub are wanted. They are essentially the same as `Getting Started` but include creating a GitHub repo with the {{cmd "gh"}} command afterward.

1. Create a directory and change into it
2. Run the {{aka}} {{cmd "init"}} command
3. Update the YAML file it opens
4. Exit your editor
5. Create and push as Git repo with {{cmd "gh repo create"}}
6. Continue with steps 9+ from `Getting Started`

Alternatively, one can simply create a GitHub repo from the web site and {{cmd "git clone"}} it down to the local machine and then run {{aka}} {{cmd "init"}} from within it.

***Learning KEG Markup Language***

Use the {{aka}} {{cmd "create sample"}} command to automatically create a new content node sample that introduces the KEG Markup Language (KEGML). You can delete it later after reading it. Or, you can use it instead of just {{aka}} {{cmd "create"}} (which gives you a blank) to help you remember how to write KEGML until you get proficient enough not to have to look it up every time.

For more about the emerging KEG 2023-01 specification and how to create content that complies for knowledge exchange and publication (while we work more on linting and validation within the {{aka}} command) have a look at <https://github.com/rwxrob/keg-spec>
