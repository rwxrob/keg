view a specific node

The {{aka}} command renders a specific node for viewing in the terminal suitable for being cutting and pasting into other text documents and description fields. The argument passed may be an integer ID or a regular expression to be matched in the title text (as with {{cmd "edit"}} and {{cmd "title"}} commands. When matting a REGEXP case insensitive matching is assumed (prefix `(?i)` is added. (See {{cmd "grep"}} for how this default an be changed.)

The {{aka}} command uses the <https://github.com/charmbracelet/glamour> package for rendering markdown directly to the terminal and therefore can be customized by setting the GLAMOUR_STYLE environment variable for those who wish. Since the popular GitHub command line utility uses this as well the same customization can be applied to both {{cmd "keg"}} and {{cmd "gh"}}.  By default, a variation on the `dark` style is used with line wrapping and margins disabled (for better cutting and pasting). To get a full copy of the style JSON used see the {{cmd "style"}} command.

If the output is not to a terminal then the `notty` Glamour theme is used automatically.
