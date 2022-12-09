package keg

import _ "embed"

// Keep all documentation in this file as embeds in order to easily
// support compiles to other target languages by simply changing the
// language identifier before compilation.

//go:embed text/en/keg.md
var _keg string

//go:embed text/en/current.md
var _current string

//go:embed text/en/directory.md
var _directory string

//go:embed text/en/titles.md
var _titles string

//go:embed text/en/delete.md
var _delete string

//go:embed text/en/index.md
var _index string

//go:embed text/en/index-update.md
var _index_update string

//go:embed text/en/last.md
var _last string

//go:embed text/en/changes.md
var _changes string

//go:embed text/en/init.md
var _init string

//go:embed text/en/edit.md
var _edit string

//go:embed text/en/random.md
var _random string

//go:embed text/en/import.md
var _import string

//go:embed text/en/columns.md
var _columns string

//go:embed text/en/grep.md
var _grep string

//go:embed text/en/view.md
var _view string

//go:embed text/en/create.md
var _create string

//go:embed text/en/keg
var _kegyaml string

//go:embed text/en/zero-node.md
var _zero_node string

const (
	_NoKegsFound     = `no kegs found`
	_NodeNotFound    = `node not found: %v`
	_InvalidNodeID   = `invalid node id: %q`
	_FileNotFound    = `file not found: %v`
	_ChooseTitleFail = `unable to choose a title`
	_AbsPathFail     = `unable to determine absolute path to current directory`
	_BadChangesLine  = `bad line in changes.md: %v`
	_NoRemoteRepo    = `%vNo remote repo has been setup.%v First create it and git push to it.`
	_NotDirNotExist  = `not a directory or does not exist: %v`
	_CantGetNextNode = `could not determine next node id: %v`
)
