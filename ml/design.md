# Design Considerations

Here are the considerations that resulted in the current design.

### No smart quotes

Smart quotes are great, if your language uses them. But KEG has a far
greater scope than just those languages. Therefore, smart-quote parsing
and rendering will have to be done outside of anything in the KEG
specification.

### Expose the internal scan.R

By exporting the scan.R (S) developers can use it to assist with
debugging by enabling S.Trace. This also allows skipping sections of
buffered data that are not Mark syntax ("front matter", etc.)
