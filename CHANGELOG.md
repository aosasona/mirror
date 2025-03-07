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

# 2.1.1

- Added new `GenerateItemType` method on generator interface to generate pure types instead of full exported declarations

# 2.2.0

- Added `OnParseItem` and `OnParseField` hooks to the parser interface
- Implemented both hooks (`OnParseItem` and `OnParseField`) for built-in parser
- Added support for non-struct embedded types (you can now embed fields that are not structs, with full support for struct tags)
- Built-in `parser.Item` types (`Scalar`, `Map`, `List` etc) are now implemented as pointers (i.e. the methods are implemented via pointer receivers)
- Renamed `FlattenEmbeddedStructs` to `FlattenEmbeddedTypes`
- Fully dropped support for legacy `ts` tag

# 2.3.0

- Added basic field lookup methods for struct items

# 2.4.0

- Expose parser's internal parse method for one-off parsing

# 2.4.1

- `fix`: Typescript generator will now check if inlining is enabled before enforcing the reference rule

# 2.5.0

- Add array generic control and allow empty headers

# 2.6.0

- Return a reference to the generator in the Typescript target
  > This is so that you can actually make changes to the underlying generator used by the target, there was previously no way to do that.

# 2.6.1

- Added no-op methods (`nil` and `void`)

# 2.7.0

- Added support for custom types (for item overrides) in parser
- Only test public interfaces/methods

# 2.7.1

- Fixed indentation level for structs in lists and maps with inlining enabled

# 2.7.2

- Caching key is only generated if caching is enabled in the parser
- Refined `unimplemented` error message for unsupported types in the parser
- Fixed `interface`/`any` detection in the parser
  > This was causing an `unimplemented` error to be returned whenever any was used in a struct or any other anonymous context
- Updated examples and tests to include `any`

# 2.7.3

- Fixed the naive parsing of mirror tags
  > `mirror:"type:{ 'foo': bar}"` will now be parsed correctly instead of skipping due to the presence of another `:` that we were naively splitting on

# 2.8.0

- Added a new `meta.Optional` type to represent the `optional` attribute.
  > Instead of just reporting false negatives in cases where the attribute was not even set at all in a mirror tag, we can now accurately represent the state with one of `OptionalNone`, `OptionalTrue` or `OptionalFalse`. This was crucial to deciding when a user was actively setting a field to be optional or not, the previous approach would have ALWAYS been false by default and cause every single field with `Nullable` set to `true` to be overriden.
- Override field nullability based on the `optional` property in the mirror tag
  > If your type has "native" nullibility (e.g `name *string`), i.e. `Nullable` is set to `true`, and the `optional` tag is set to `false` as in `optional:false`, your type will be generated without being nullable (e.g. `name: string` instead of `name: string | undefined`). This was previously not properly used or implemented in the past and required using the `type` override.

> [!WARNING]
> In the future, we might have to split `nullability` and `optionality` since they kind of mix too uncomfortably and might change this behaviour; making it too difficult to do something like `foo?: string | null` at the moment (as far as I can remember now anyway).
