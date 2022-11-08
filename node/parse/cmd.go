package parse

import (
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/help"
	"github.com/rwxrob/vars"
)

var Cmd = &Z.Cmd{
	Name:    `parse`,
	Summary: `parse/print semantic node content`,
	Usage:   `[TYPE [FILTER|FILE|DIR]]`,
	Commands: []*Z.Cmd{help.Cmd, conf.Cmd, vars.Cmd,
		yamlCmd, jsonCmd, xmlCmd,
	},
	ConfVars: true,
	VarDefs:  Z.VarVals{`nl`: `KEGNL`},
	Shortcuts: Z.ArgMap{
		`get`:  {`var`, `get`},
		`pegn`: {`emb`, `cat`, `kegml.pegn`},
	},
	Params: []string{
		`title`, `heading`, `block`, `include`, `incfile`, `incnode`,
		`bulleted`, `numbered`, `figure`, `fenced`, `tex`, `quote`, `raw`,
		`ref`, `refs`, `link`, `linkfile`, `linknode`, `tags`, `tag`, `div`,
		`para`, `bullet`, `number`, `span`, `inflect`, `bold`, `verbatim`, `math`,
		`deleted`, `squoted`, `dquoted`, `quoted`, `bracketed`, `parens`,
		`braced`, `angled`, `url`, `longdash`, `shortdash`, `plain`, `ellipsis`,
		`word`,
	},
	Description: `
		The {{ cmd .Name }} command parses and prints different (semantic)
		parts of the KEG node. Matches are printed one per line with any
		line returns replaced with {{ pre KEGNL }} (which can be changed
		with the {{ cmd "set nl" }} command.

		The first parameter indicates the type of parsed content wanted from
		the KEGML file. Type names come from the supported KEGML PEGN
		specification available for reference from the {{ cmd "pegn" }}
		command.

		The second argument indicates the node (or nodes) to parse by KEG
		node identifier or scope filter. See the {{ cmd "keg" }} command
		help for more information about KEG.

		The second argument may also simply be a file system path to a file
		or directory containing a README.md file.

		If the second argument is omitted, the current node is assumed 
		{{ pre "set current" }}. If no current node it set, the parent caller's
		{{ pre "current" }} value is used (if it exists). If even then no
		current node can be resolved, the README.md file within the current
		working directory is assumed.

	`,
}

var yamlCmd = &Z.Cmd{
	Name:    `yaml`,
	Summary: `print/query AST in YAML`,
	Usage:   `[QUERY]`,
}

var jsonCmd = &Z.Cmd{
	Name:    `json`,
	Summary: `print/query AST in JSON`,
	Usage:   `[QUERY]`,
}

var xmlCmd = &Z.Cmd{
	Name:    `xml`,
	Summary: `print AST in XML`,
}
