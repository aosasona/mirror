# 2.2.0

- Added `OnParseItem` and `OnParseField` hooks to the parser interface
- Implemented both hooks (`OnParseItem` and `OnParseField`) for built-in parser
- Added support for non-struct embedded types (you can now embed fields that are not structs, with full support for struct tags)
- Built-in `parser.Item` types (`Scalar`, `Map`, `List` etc) are now implemented as pointers (i.e. the methods are implemented via pointer receivers)
- Renamed `FlattenEmbeddedStructs` to `FlattenEmbeddedTypes`
- Fully dropped support for legacy `ts` tag

# 2.1.1

- Added new `GenerateItemType` method on generator interface to generate pure types instead of full exported declarations

# v2.0.0

- Full rewrite to support pluggable "backends"
- Deprecated legacy `ts` tag
- Added parser tokens to make parser result more usable outside the library
- Added more tests (still not as extensive as I'd want but it is a start)
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
