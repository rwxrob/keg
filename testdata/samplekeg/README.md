# Sample Knowledge Exchange Graph

This is the main `README.md` file of a "keg" (a **knowledge exchange graph**). Change it to be whatever you want. This is just a boilerplate sample.

On the **KEG Web** (a combination of WorldWideWeb and KEG) this file is people would see visiting your **KEG Site.**

## Content nodes

```
kn create
kn create sample
```

A keg contains one or more **content nodes** in the form of directories that have an integer number and each contain their own required `README.md` file along with an optional **meta matter** YAML file (`meta`), and any other content files that make sense to add (just not other content nodes). This produces a flat directory structure where all the content lives at the root level and can easily be organized by changing the content of the nodes themselves (instead of moving them around and breaking links).

Content nodes are written in KEGML. Have a look at the [sample content node](1) for a full description of KEGML

## Expanded KEGML

Being the main entry point, KEGMLX is allowed which includes the following:

* Multiple headings
* Special node include links (no `../`, but have **query code**)

