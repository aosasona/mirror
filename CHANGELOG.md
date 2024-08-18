# v2.0.0

- Added more tests (still not as extensive as I'd want but it is a start)
- Deprecated `ts` tag (full removal in a future release)
- Added support for multiple languages
  > Currently, it only includes the Typescript backend but it can now be easily extended in the future to support other languages like Gleam, Rust etc by simply adding more "backends"
- Added support for embedded structs
- Added support for common built-in types (`time.Time`, `time.Duration`, `sql.NullX` types)
- Validate the presence of all referenced types when inlining is disabled
- Implemented new "parser": the parser now returns tokens that can be manipulated or used by several generators
- Added support for user extensions
  - Custom parsers can now be used by implementing `ParserInterface`
  - Custom targets can be added by implementing `TargetInterface`
- Exposed more internal modules (parser, default backends etc)
