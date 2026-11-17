# Synthetic wire vectors

[vectors.json](vectors.json) contains constructed examples for the `native-18`
profile. Hex strings have no whitespace. Lengths count bytes; SHA-256 values
are lowercase hex. All examples are synthetic and were never sent to a server.

| Key | Meaning |
| --- | --- |
| `frames` | Payload and its complete escaped record, including End |
| `guids` | Displayed UUID and its 16 binary bytes, without an atom tag |
| `integers` | Canonical uint32 atom examples without structural flags |
| `byte_atoms` | Raw bytes and their atom encodings without structural flags |
| `ordinary_list` | One ordinary list, including begin/end and separators |
| `empty_security_blob_hex` | Complete standalone empty blob, without a parent separator |
| `text_strings` | `BOM + {"S", quoted_value}`, including doubled quotes |
| `introduction` | Inner V8IntroTicketReq payload for a synthetic name |
| `http_files` | Repeated fill-byte payload, chunk lengths, complete HTTPFile length/hash |
| `http_request` | Abbreviated synthetic login HTTP, RPC and framed record |
| `http_continuation` | Same HTTP as the request, next counter, continuation envelope |
| `http_reply` | Synthetic success HTTP, RPC and framed record |

The repeated-byte HTTPFile cases test serialization boundaries, not HTTP syntax.
Their hashes cover the **complete standalone HTTPFile object**, beginning with
`40` and ending with `20`. The HTTP request has no real credentials and omits
normal application headers. It is not a valid login attempt.

The examples were exported offline from the separately maintained product
codec identified in [sources.json](../evidence/sources.json). Export verified
all two-part split positions for the seven framing examples and parsed the
synthetic reply back through that codec. That is implementation consistency,
not an independent native oracle.

The public Go tests independently reconstruct the selected encodings from
the documented rules and compare exact bytes/hashes. They also check HTTP
Content-Length and local Markdown file links. They are deliberately small
specification checks, not a complete streaming parser or client library.

```sh
go test ./...
go vet ./...
```

No test opens a socket, starts 1C, reads captured data or uses authentication
material. Tests require no modules outside the Go standard library.
