# Evidence register

This specification is a curated synthesis of retained research and a separately
maintained implementation. Source files are identified by repository-relative
path, revision and SHA-256 in [sources.json](sources.json). The manifest contains
identifiers and hashes, not source excerpts or raw experimental payloads.
Source hashes cover Git blob content at the recorded revision, independently
of checkout line-ending conversion.

The source repositories and raw evidence are **not bundled and may not be
publicly accessible**. A hash binds an input; it does not independently verify
the interpretation or make a historical live run publicly reproducible.
The [synthetic vectors](../vectors/README.md) are the reproducible public portion.

## Claims

| ID | Claim / source subject | Basis and limit |
| --- | --- | --- |
| E01 | Native layering and build-18 baseline: historical HTTP-over-TCP guide and parity report | Synthetic fixtures and reported bounded live path; old chunk examples superseded by E07 |
| E02 | Framing tokens, escaping and stream handling: `wire/transport.go` | Implementation source and synthetic verification; malformed-stream native behavior incomplete |
| E03 | Text lists and GUID order: `wire/text.go`, `wire/uuid.go` | Implementation source and synthetic verification |
| E04 | Selected binary atoms and specialized encoders: `wire/binary.go` | Implementation source; selected native token semantics corroborated by E07 |
| E05 | Introduction, query and release request shapes: `wire/bootstrap.go` | Source-derived accepted subset, not a complete native state machine |
| E06 | HTTP RPC envelopes and local validation: `wire/http.go` | Implementation source; local caps and permissive parsing are not native protocol rules |
| E07 | HTTPFile continuation without `80`: HTTP file chunk report, 2026-11-12 | Exact native instructions, one retained native reply, synthetic regression and bounded live corrected request |
| E08 | Authentication/version distinctions: parity, adaptive-client and build-27 policy reports | Static, original-client synthetic and bounded live results kept separate |
| E09 | Cleanup and uncertain mutation after lost response: upload-break and context-live reports | Bounded live observations on 8.3.18.1208; no general lifetime or retry guarantee |

## Exact shipped inputs

These are identities recorded by the source reports; no proprietary binary is
included and no installed binary was executed for this publication.

| Build | File | SHA-256 |
| --- | --- | --- |
| 8.3.18.1208 | core83.dll | `1f71a4d95bbff56997024c32840c0395012a02110e3b5d18ed61c9a9ec35094f` |
| 8.3.18.1208 | rtrsrvc.dll | `bb03a507e8c56dd71c874234c54ba6d82c4e223c592998b83cb02db7b14da8d8` |
| 8.3.27.2342 | comcntr.dll | `5ca7d51daba0f250f27ca3623527db3555fd23d86c5c50ea9aef260255e487d9` |
| 8.3.27.2342 | rtrsrvc.dll | `ef8aac04965cafcf35252be1d3bfec237d88e52215f3db3b14c7e59b4d3ba310` |
| 8.3.27.2342 | core83.dll | `83001ee41c5368a9ae6c73fd9a1a58fb2129db462cf882baceac4ce2877ef0d9` |

The E07 native reply observation reported 23,924 HTTP bytes, split into
16,384 and 7,540. The decoded outer payload was 23,996 bytes, SHA-256
`d5eb4e21c975ce56b328838a1d2023c030f7e8f571a79fc72ef02447f2e52cc9`.
Only these metadata and structural observations are published; the payload is
not part of the public test corpus.

## Publication verification

The synthetic vectors were compared offline with the pinned product codec.
The repository's independent checks cover published encodings, lengths and
document links. They do not connect to a server or validate historical native
receipts. See [vector scope](../vectors/README.md) before interpreting a green
test run as a compatibility claim.
