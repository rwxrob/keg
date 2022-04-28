// Copyright 2022 bonzai-example Authors
// SPDX-License-Identifier: Apache-2.0

// Package example provides the Bonzai command branch of the same name.
package example

import (
	"log"
	"text/template"

	Z "github.com/rwxrob/bonzai/z"
	compfile "github.com/rwxrob/compfile"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/help"
	"github.com/rwxrob/vars"
)

// Most Cmds that make use of the conf and vars branches will want to
// call SoftInit in order to create the persistence layers or whatever
// else is needed to initialize their use. This cannot be done
// automatically from these imported modules because Cmd authors may
// with to change the default values before calling SoftInit and
// committing them.

func init() {
	Z.Conf.SoftInit()
	Z.Vars.SoftInit()
}

// Cmd provides a Bonzai branch command that can be composed into Bonzai
// trees or used as a standalone with light wrapper (see cmd/).
var Cmd = &Z.Cmd{

	Name:      `example`,
	Summary:   `an example of Bonzai composite command tree`,
	Version:   `v0.4.0`,
	Copyright: `Copyright 2021 Robert S Muhlestein`,
	License:   `Apache-2.0`,
	Site:      `rwxrob.tv`,
	Source:    `git@github.com:rwxrob/bonzai-example.git`,
	Issues:    `github.com/rwxrob/bonzai-example/issues`,

	// Composite commands, local and external, all have their own names
	// that are added to the command tree depending on where they are
	// composed.

	Commands: []*Z.Cmd{

		// standard external branch imports (see rwxrob/{help,conf,vars})
		help.Cmd, conf.Cmd, vars.Cmd,

		// local commands (in this module)
		BarCmd, ownCmd, pkgexampleCmd, BazCmd,
	},

	// Add custom BonzaiMark template extensions (or overwrite existing ones).
	Dynamic: template.FuncMap{
		"uname": func(_ *Z.Cmd) string { return Z.Out("uname", "-a") },
		"dir":   func() string { return Z.Out("dir") },
	},

	Description: `
		The **{{.Name}}** command branch is a well-documented example to get
		you started.  You can start the description here and wrap it to look
		nice and it will just work.  Descriptions are written in BonzaiMark,
		a simplified combination of CommonMark, "go doc", and text/template
		that uses the Cmd itself as a data source and has a rich set of
		builtin template functions ({{pre "pre"}}, {{pre "exename"}}, {{pre
		"indent"}}, etc.). There are four block types and four span types in
		BonzaiMark:

		Spans

		    Plain
		    *Italic*
		    **Bold**
		    ***BoldItalic***
		    <Under> (brackets remain)

		Note that on most terminals italic is rendered as underlining and
		depending on how old the terminal, other formatting might not appear
		as expected. If you know how to set LESS_TERMCAP_* variables they
		will be observed when output is to the terminal.

		Blocks

		1. Paragraph
		2. Verbatim (block begins with '    ', never first)
		3. Numbered (block begins with '* ')
		4. Bulleted (block begins with '1. ')

		Currently, a verbatim block must never be first because of the
		stripping of initial white space.

		Templates

		Anything from Cmd that fulfills the requirement to be included in
		a Go text/template may be used. This includes {{ "{{ .Name }}" }}
		and the rest. A number of builtin template functions have also been
		added (such as {{ "indent" }}) which can receive piped input. You
		can add your own functions (or overwrite existing ones) by adding
		your own Dynamic template.FuncMap (see text/template for more about
		Go templates). Note that verbatim blocks will need to indented to work:

		    {{ "{{ dir | indent 4 }}" }}

		Produces a nice verbatim block:

		{{ dir | indent 4 }}

		Note this is different for every user and their specific system. The
		ability to incorporate dynamic data into any help documentation is
		a game-changer not only for creating very consumable tools, but
		creating intelligent, interactive training and education materials
	 	as well.

		Templates Within Templates

		Sometimes you will need more text than can easily fit within
		a single action. (Actions may not span new lines.) For such things
		defining a template with that text is required and they you can
		include it with the {{pre "template"}} tag.

		    {{define "long" -}}
		    Here is something
		    that spans multiple
		    lines that would otherwise be too long for a single action.
		    {{- end}}

		    The {{ "**{{.Name}}**" }} branch is for everything to help with
		    development, use, and discovery of Bonzai branches and leaf
		    commands ({{ "{{- template \"long\" \"\" | pre -}}" }}).

		The help documentation can scan the state of the system and give
		specific pointers and instruction based on elements of the host
		system that are missing or misconfigured.  Such was *never* possible
		with simple "man" pages and still is not possible with Cobra,
		urfave/cli, or any other commander framework in use today. In fact,
		Bonzai branch commands can be considered portable, dynamic web
		servers (once the planned support for embedded fs assets is
		added).`,

	Other: []Z.Section{
		{`Custom Sections`, `
			Additional sections can be added to the Other field.

			A Z.Section is just a Title and Body and can be assigned using
			composite notation (without the key names) for cleaner, in-code
			documentation.

			The Title will be capitalized for terminal output if using the
			common help.Cmd, but should use a suitable case for appearing in
			a book for other output renderers later (HTML, PDF, etc.)`,
		},
	},

	// no Call since has Commands, if had Call would only call if
	// commands didn't match
}

// Commands can be grouped into the same file or separately, public or
// private. Public let's others compose specific subcommands (example.Bar),
// private just keeps it composed and only available within this Bonzai
// command.

// exported branch
var BarCmd = &Z.Cmd{
	Name: `bar`,

	// Aliases are not commands but will be replaced by their target names
	// during bash tab completion. Aliases show up in the COMMANDS section
	// of help, but do not display during tab completion so as to keep the
	// list shorter.
	Aliases: []string{"B", "notbar"},

	// Commands are the main way to compose other commands into your
	// branch. When in doubt, add a command, even if it is in the same
	// file.
	Commands: []*Z.Cmd{help.Cmd, fileCmd},

	// Call first-class functions can be highly detailed, refer to an
	// existing function someplace else, or can call high-level package
	// library functions. Developers are encouraged to consider well where
	// they maintain the core logic of their applications. Often, it will
	// not be here within the Z.Cmd definition. One use case for
	// decoupled first-class Call functions is when creating multiple
	// binaries for different target languages. In such cases this
	// Z.Cmd definition is essentially just a wrapper for
	// documentation and other language-specific embedded assets.
	Call: func(_ *Z.Cmd, _ ...string) error { // note conventional _
		log.Printf("would bar stuff")
		return nil
	},
}

// Different completion methods are be set including the expected
// standard ones from bash and other shells. Not that these completion
// methods only work if the shell supports completion (including
// the Bonzai Shell, which can be set as the default Cmd to provide rich
// shell interactivity where normally no shell is available, such as in
// FROM SCRATCH containers that use a Bonzai tree as the core binary).

// private leaf
var fileCmd = &Z.Cmd{
	Name:     `file`,
	Commands: []*Z.Cmd{help.Cmd},
	Comp:     compfile.New(),
	Call: func(x *Z.Cmd, args ...string) error {
		if len(args) == 0 {
			return x.UsageError()
		}
		// always use "log" and not "fmt" for errors and debugging
		log.Printf("would show file information about %v", args[0])
		return nil
	},
}

// When combining a high-level package library with a Bonzai command it
// is customary to create a pkg directory to avoid cyclical package
// import dependencies. This also makes Bonzai modules usable in
// multiple ways:
//
// * As a command branch in a composite monolith or multicall
// * As a standalone command (see "bonzai-example" directory)
// * As a high-level library unrelated to Bonzai at all (see "pkg")

// private leaf
var pkgexampleCmd = &Z.Cmd{
	Name: `pkgexample`,

	// Several argument checks are available to keep your Call functions
	// minimal and free of redundant if statements.
	NumArgs: 1,

	// Params are specifically to help with completion and seeing regular
	// possibilities. Be careful not to use a param as a Command. It the
	// param is effectively a completely different leaf command then make
	// its own, even if it is tiny and just as a Name, Summary, and Call.
	Params: []string{"bar", "baz"},

	Call: func(_ *Z.Cmd, args ...string) error {
		Foo(args[1]) // calls high-level pkg library function
		return nil
	},
}
