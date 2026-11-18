# Evidence register

This specification is a curated synthesis of retained research and a separately
maintained implementation. Source files have neutral IDs, subject titles,
claim IDs and SHA-256 in [sources.json](sources.json). The public manifest does
not expose local paths, private repository names, machine identifiers or source
excerpts. A separate unpublished receipt retains the source-location mapping.
Each entry states its hash basis: committed source bytes before checkout
line-ending conversion, or retained observation bytes where identified.

The source repositories and raw evidence are **not bundled and may not be
publicly accessible**. A hash binds an input; it does not independently verify
the interpretation or make a historical live run publicly reproducible.
The [synthetic vectors](../vectors/README.md) are the reproducible public portion.

## Claims

| ID | Claim / source subject | Basis and limit |
| --- | --- | --- |
| E01 | Native layering and build-18 baseline: historical HTTP-over-TCP guide and parity report | Synthetic fixtures and reported bounded live path; old chunk examples superseded by E07 |
| E02 | Framing tokens, escaping and stream handling | Implementation and native consumer study; malformed-stream behavior incomplete |
| E03 | Text lists and GUID order | Implementation source and synthetic verification |
| E04 | Binary atoms, strings, Numeric and specialized encoders | Implementation plus exact-build native scalar/string studies |
| E05 | Complete implemented bootstrap, manager selection, query and release | Source-derived state machine, bounded native-18 acceptance; unknown alternatives explicit |
| E06 | HTTP RPC envelopes and local validation | Implementation source; local caps and permissive parsing are not native protocol rules |
| E07 | HTTPFile continuation without `80`: HTTP file chunk report, 2026-11-12 | Exact native instructions, one retained native reply, synthetic regression and bounded live corrected request |
| E08 | Authentication/version distinctions: parity, adaptive-client and build-27 policy reports | Static, original-client synthetic and bounded live results kept separate |
| E09 | Cleanup and uncertain mutation after lost response: upload-break and context-live reports | Bounded live observations on 8.3.18.1208; no general lifetime or retry guarantee |
| E10 | Native login and SecPswd V1/V2/V3 | V1/V3 original-code comparisons; V2 schema/synthetic evidence; new independent synthetic transcript comparison |
| E11 | Build-27 n/c policy, RSA header, 3DES and IV progression | Original client with owned synthetic peers; no real build-27 native server proof |
| E12 | SSPI, explicit NTLMv2 and Kerberos integration | OS-provider live baseline; portable providers synthetic, with standard protocol references |
| E13 | Web prelogin/settings/login, served contract and access boundaries | Shipped-source and synthetic evidence; build-18 legacy and 18/27 Basic live subsets |
| E14 | Administration RPC, collection profiles and worker enumeration | Build-18 bounded live subset; later original clients against synthetic peers |
| E15 | Typed values, form calls, collections, presentation and synchronization | Build-18 source-derived schemas and bounded live scalar/form operations; general renderer/runtime unimplemented |
| E16 | Dynamic-list windows and compact payloads | Retained native-18 payloads, independent decoding and bounded live windows |
| E17 | Processing uploads, confirmation, chunks and storage | Build-18 bounded live subset; 18/27/85 source-derived chunk construction and ownership limits |
| E18 | Route inventory, server roles and preauthentication surfaces | Static route contracts and bounded observations; no universal anonymous access or arbitrary execution |
| E19 | [Independent public-spec Go client](../spec/20-independent-client.md) | Fresh-context public-input implementation, frozen before feedback; direct-worker and manager-discovery Windows login/logout passed on the authorized native-18 laboratory; separate password-readiness audit is synthetic only, with no new live password/provider-family coverage |
| E20 | [Independent web password continuation](../spec/21-independent-web-auth.md) | Public-input client passed web Legacy login/logout on build 18; web V1/V3 challenges progressed on builds 22/27 but final application-start requests failed; no live V2 or completed SecPswd acceptance |

## Exact shipped inputs

These are identities recorded by the source reports; no proprietary binary is
included. The original publication was offline; subsequent live experiments
identify their server-runtime and implementation boundaries separately.

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

The original synthetic vectors were compared offline with the retained codec.
The expansion independently constructed eight complete SecPswd exchanges and
eight password-key cases, compared them against the retained state machine,
and constructed a bootstrap example. All identities, credentials, nonces and
message bytes in that corpus are synthetic. The one-time generators and
debugging/research tests are not published.
The repository's independent checks cover published encodings, lengths and
document links. They do not connect to a server or validate historical native
receipts. See [vector scope](../vectors/README.md) before interpreting a green
test run as a compatibility claim.

The subsequent [bootstrap/authentication audit](../spec/19-bootstrap-auth-audit.md)
checked the published credential corpus with a standalone consumer, constructed
16 complete native RPC pairs and seven HTTP credential flows, and separately
compared native bytes/replies with the retained implementation. Its author had
access to that implementation; it is not a blind clean-room or live-server proof.
No additional historical acceptance claim or private source export was added.
