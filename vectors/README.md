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
material from a real account. Tests require no modules outside the Go standard library.

## Authentication and bootstrap corpus

[authentication.json](authentication.json) contains eight password hash/key
vectors and eight complete credential exchanges: two each for V1/V2 and four
V3 exchanges covering all hash modes. Packet names C1/S2/... indicate direction
and sequence. Each gives standard Base64, decoded packet hex and, for encrypted
whole packets, plaintext before padding (including CRC where applicable).
The published toy user/password, keys and deterministic random streams are
synthetic inputs, never captured authentication material. Use fresh secure
randomness for actual clients.

The local publication audit constructed every packet independently from the
documented primitives and compared every client message and final state against
the retained implementation. It also checked all eight hash/key outputs and
rejected final-message reuse after completion. Those generators and debugging
programs remain outside this repository. The existing public Go conformance
suite checks the original wire corpus and document links; it does not replay
the new credential state machines. Publication-time comparison is a separate
verification result, not an additional live-server claim.

[bootstrap.json](bootstrap.json) gives a synthetic worker introduction request
as UTF-8 text, exact bytes and a framed record. It includes the BOM and unquoted
Base64 token, a neutral computer label and fixed example UUID. It is the first
bootstrap message, not an authenticated connection or reusable provider token.

## Complete session and orchestration traces

[native-session.json](native-session.json) supplies 16 exact RPC request/reply
pairs, their framed records and semantic fields: manager discovery, worker
bootstrap, inner security, Windows application login, logout and release.
Its provider tokens are explicitly scripted ASCII stubs, not real NTLM/Kerberos.

[login-orchestration.json](login-orchestration.json) supplies seven native/web
credential flows (30 HTTP exchanges) using the published authentication packets.
It covers V1/V2/V3 native placement, web V1/V3 settings with/without S2, full
session retention, sid replacement and logout. It names ten rejection mutations.
These are conformance data, not executable debugging tests or captured traffic.

The [audit chapter](../spec/19-bootstrap-auth-audit.md) defines the file fields,
expected transitions, mutations and publication-time validation. The existing
public tests do not replay these new files; their separate offline check used
only public inputs and the Go standard library, followed by a separate comparison
with the retained codec. No actual client/server interoperability is implied.
