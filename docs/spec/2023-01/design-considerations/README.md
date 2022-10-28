# Design Considerations

**Use of 4-letter suffix for YAML** is preferred by the creator (instead
of `.yml`).

**Mixing meta data into the data itself as "front matter"** will never
be considered since it violates fundamental principles of separation of
concerns.

**Simplified YAML is the KEG default structured data format.** DATA
Nodes, which fundamentally require multiple formats for data
representation, and not included in this constraint. Any communication
or configuration between any KEG tooling MUST use simplified YAML
(which includes JSON).
