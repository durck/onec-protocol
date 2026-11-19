# Bootstrap and authentication implementation audit

Audit date: 2026-11-18. Evidence: E03, E05, E06, E08, E10 and E13.

Subsequent result: a [fresh public-spec client experiment](20-independent-client.md)
now provides bounded live Windows-login acceptance. This chapter preserves the
earlier audit's own method and limits; its checks were not live-server tests.

The documented native-18 subset is sufficient to construct the covered requests
and implement the illustrated bootstrap/password state machines without private
research inputs. That is a narrower conclusion than proving a new independent
client interoperates with a real server. The executable check described below
has no private implementation dependency, but its author also reviewed the
retained implementation: this was **not a blind clean-room implementation**.

## Findings and corrections

| Finding | Consequence for an implementer | Resolution |
| --- | --- | --- |
| A direct worker-endpoint string was listed as acceptable in selection | A synthetic reply following that alternative failed the retained selection parser | Require `S("host:port")`; URLType's inner URL remains a direct string |
| Initial application sync was described as reserved at 1 | An ordinary first operation would incorrectly send 2 | Start at 0, first ordinary operation sends 1; immediate logout alone sends 2 |
| Request/reply encoding and top-level suffix were scattered | Binary parsing of text replies, or an extra scalar/list terminator | Explicit dispatch table for bootstrap/query/selection/HTTP/release |
| Initial RPC counter and first-HTTP lifetime were implicit | Off-by-one counters or reuse of the initial HTTP envelope during negotiation | Introduction uses N; increment each exchange; first attempted HTTP consumes its envelope |
| SecurityBlob described bytes without explicitly requiring their atom tag | A nonempty provider token could be serialized as untagged data | Both length and bytes are atoms; added a three-byte example |
| Inner negotiation did not explicitly freeze the HTTP PDU or reset request flags | Recreated nonces, changed sid or continued securityFlag=1 after completion | Repeat identical PDU bytes during token legs; new logical HTTP uses empty token and flag 0 |
| Native credential selection could be mistaken for negotiation | Guessed query parameters or automatic version-based downgrades | Explicit profile selection; internal selectors are not HTTP fields |
| Native HTTP identity/header roles were underspecified | Worker address/SPN substituted into application headers; HTTP library additions | Logical ServerName, encapsulated HTTP headers, origin-form target and body termination specified |
| Web prelogin appeared to require a body seanceID | Header-only routing sessions rejected, or lowercase factor seanceid used as application sid | Header capture followed by exact-case outer-to-inner seanceID precedence |
| Settings could be read as GET or a separate codec exchange | Missing C1, mismatched challenge state or wrong number of login requests | V1/V3 POST C1; preserve state after S2; no-S2 byte identity is the retained transcript's policy, qualified by the later offline audit |
| Web session replacement was not distinguished from native routing stability | Stale sid after a web response, or improper rejection of rotation | Retain full replacement session and recompute web sid; native intermediate UUID remains stable |
| XML, native JSON and web JSON acceptance were conflated | Unnecessary required fields or acceptance before final proof | Separate validators and a minimal synthetic XML document |
| V1 key deduplication could be applied to V2; V3 flag relationships were implicit | Missing second encrypted V2 alternative or invented V3 equality requirement | Specify both V2 candidates and independently validated V3 flags |
| Final proof placement and marker meaning could be overgeneralized | HTTP 200/session receipt treated as password acceptance | Require correct in-state proof; V2/V3 suc is a marker, not an independent MAC |

Details are normative only within the bounded profiles in
[bootstrap](03-bootstrap.md), [HTTP RPC](04-http-rpc.md),
[native authentication](07-native-authentication.md),
[SecPswd](09-secpswd.md) and [web authentication](11-web-authentication.md).
Local validator behavior is identified as such; it is not evidence that every
server rejects other encodings.

## Complete native session transcript

[native-session.json](../vectors/native-session.json) contains two consecutive
connections: manager discovery followed by worker login and immediate logout.
All addresses, identities, UUIDs and messages are constructed examples.

Each context supplies endpoint/role, process/client/last UUIDs, remote target,
selector, the separate server greeting and unframed client Control3 item.
Each exchange supplies its operation and counter, semantic token/PDU fields,
exact request/reply payload hex and complete escaped record hex including End.
Text payloads additionally have UTF-8 representations including the real BOM.
No Control3 or greeting is appended to a framed record.

The Boolean `security_flag` and `security_complete` fields describe HTTP and
bootstrap security where applicable; they have no meaning on query/select/release.
Missing `token_hex` or `reply_token_hex` means empty bytes. Missing reply HTTP
means undefined HTTPFile; empty HTTP **body** instead has a full HTTP PDU.
The `continuation` property selects only the HTTP request envelope.

| Context / RPC counter | Request | Expected result / state |
| --- | --- | --- |
| Manager 2000 | Introduction | Incomplete; acknowledgment supplies provider target/package |
| Manager 2001 | First provider token | Incomplete; nonempty challenge |
| Manager 2002 | Second provider token | Both sides complete; retain manager target |
| Manager 2003 | Query SelectInterface | Text success, no trailing scalar |
| Manager 2004 | Select ExampleBase | Binary success, typed endpoint, trailing 2 |
| Manager 2005 | Release | Binary success without trailing scalar; close manager |
| Worker 1000..1002 | Introduction and two provider legs | Fresh provider, last UUID and target; same process UUID |
| Worker 1003 | Query HTTPInterface | Text success; HTTP initial envelope still unused |
| Worker 1004 | Windows login PDU, empty token, flag 0 | Initial HTTP envelope; incomplete security, no initial challenge |
| Worker 1005 | Same PDU, introduction, flag 1 | Continuation envelope; introduction acknowledgment |
| Worker 1006 | Same PDU, first inner token, flag 1 | Continuation envelope; challenge |
| Worker 1007 | Same PDU, second inner token, flag 1 | Complete, empty final token, HTTP 200/XML/session |
| Worker 1008 | Logout, empty token, flag 0, sync=2 | Continuation envelope; complete HTTP 200 |
| Worker 1009 | Release | Release context and close connection/providers |

The provider contract in this file is an explicit **scripted stub**: ASCII
tokens identify first/second client legs and one server challenge, with local
completion on the second client leg. There are three separate states: manager,
worker and inner. These tokens are not valid NTLM/Kerberos messages and cannot
authenticate to a server. Replace the stub with the
[public provider contract and standards](08-security-providers.md) in a client.

The fixture covers the empty-final-token path, where local provider completion
is already required. The nonempty-final-proof path must additionally call the
provider and require completion with no remaining output; it is specified in
chapter 3 but is not exercised by this transcript.

## Password HTTP orchestration transcripts

[login-orchestration.json](../vectors/login-orchestration.json) contains seven
complete login/logout flows, using packets from
[authentication.json](../vectors/authentication.json):

| Profile | Codec | Settings result | Credential login POSTs |
| --- | --- | --- | ---: |
| Native | V1 | No settings call | 2 |
| Native | V2 | No settings call | 2 |
| Native | V3 | No settings call | 3 |
| Web | V1 | No cred | 2 |
| Web | V1 | S2 | 1 |
| Web | V3 | No cred | 3 |
| Web | V3 | S2 | 2 |

The `credential_exchange` field identifies an authentication corpus exchange.
Each HTTP request has method, target, explicit operation headers and body; each
response has status, headers and body. Native records require chapter 7 common
headers and chapter 4 packaging; they are not bare HTTP sent to a worker port.
Web paths are publication-relative. Compute Content-Length from UTF-8 body
bytes; header order and ordinary web-library headers are not prescribed here.
The web version string is illustrative, not a claim that every illustrated
codec is accepted by that exact build. Native V2 also uses the studied
`MobileClient` label for both `User-Agent` and `1C-ApplicationName` when adding
common headers. The [adapter trace](25-http-user-agent.md) establishes that
`User-Agent` drives the inspected mobile classifier; the two equal example
values do not make these headers interchangeable.

The web prelogin examples deliberately give different header, outer-exception
and inner-exception sessions: the last exact-case seanceID wins. A lowercase
seanceid remains unrelated. Settings replaces the full session again, and the
next login uses its projected UUID. Further web replacements exercise sid
updates. Native intermediate routing stays fixed; the final session changes.
Every logout carries the **whole** final encoded session, not its UUID prefix.

For a settings response without cred, C1 is byte-identical in settings and
the first login. When settings supplies S2, the first login is C3. No flow
in this synthetic corpus regenerates a nonce or starts a new codec at that
boundary. The later [installed-client audit](22-offline-auth-audit.md) found
that stock 22/27 loaders generate a fresh C1 on the no-S2 branch. The existing
corpus remains a retained-client profile, not a claim of exact stock byte identity
or server acceptance of either branch. Never reset state after receiving S2.

## Rejection cases

The orchestration corpus names the following deterministic mutations. Apply
one mutation to an otherwise valid flow and reject it at the affected step;
do not treat a later cleanup attempt as authentication acceptance.

| Case ID | Mutation / expected rejection |
| --- | --- |
| missing-final-proof | Remove response.credAns4 from native V1's final reply |
| mixed-proof-names | Add response.credAns6 alongside the V1 final proof |
| mixed-challenge-and-application | Add top-level cred alongside final response |
| top-level-final-proof | Replace final object with top-level credAns4 |
| missing-native-session | Remove vrs-session from first native credential response |
| wrong-settings-method | Replace web settings POST with GET |
| missing-application-session | Remove prelogin header and all exact-case seanceID fields; keep lowercase seanceid only |
| unknown-second-factor | Put reason_2fa=9 in the web prelogin exception |
| changed-native-intermediate-routing | Replace the second native V3 intermediate session with a different projected UUID |
| wrong-login-sid | Change the second native V1 login sid without changing the retained session |

Credential-level checks additionally cover, for each of eight exchanges:
missing final proof (never becomes authenticated), premature final message,
final reuse after completion, a first challenge with its first decoded byte
XORed with 01, and a final packet replaced with Base64(ASCII err).

## What was actually verified

An unpublished standalone Go consumer, using only the standard library and
these public JSON files, passed:

- Eight password hash/key cases and eight complete SecPswd exchanges, comparing
  all 20 outgoing client packets exactly and checking final state and random
  byte consumption; 40 credential rejection/state checks.
- Seven HTTP flows containing 30 request/response exchanges and the ten
  orchestration rejection cases above.
- Reconstruction of all 16 native RPC request/reply pairs and framing from
  semantic fields, counter progression and unchanged inner-negotiation PDU.

A **separate** comparison program using the retained implementation matched
all 16 outgoing RPC payloads exactly and parsed all 16 constructed replies,
including returned targets, tokens, selection and HTTP. That comparison caught
the direct-string selection error before publication. It establishes agreement
with that implementation, not an independent server oracle.

One-time checking programs remain unpublished; the public deliverables are the
specification and synthetic conformance data. The repository's existing
`go test ./...` checks cover its original wire vectors and local document links;
they do **not** run the new cryptographic/orchestration replay. Do not infer
those additional checks from a green CI badge. Publication-time results above
are reported separately and can be reproduced by another implementation using
the public inputs and rules.

## Remaining limits for an independent client

- No server, native component or real account was used in this audit. The later
  [fresh-client experiment](20-independent-client.md) supplies two bounded live
  Windows-login paths; other profiles do not inherit that acceptance.
- NTLM/Kerberos implementations still need their public standards or libraries;
  this repository documents their 1C integration, not every underlying standard.
- The byte-comparison baseline covers ASCII endpoint/infobase identities and
  ASCII passwords. Full Unicode uppercasing is specified but is not validated
  by these password vectors; simple per-character uppercasing is insufficient.
- Native password-profile auto-negotiation, full opaque/extended bootstrap
  fields, general multi-host selection, and all exception variants remain open.
- Build-27 protected transport and 8.5/native or web V3 real-server acceptance
  do not gain new evidence from these baseline synthetic flows.
- Login acceptance does not implement a general application runtime, forms,
  metadata interpretation or every later operation. Follow the separate chapters
  and [version boundaries](06-version-matrix.md) for those features.

An implementer can therefore start with the bounded profile without private
research, while keeping these unsupported states explicit. Complete replacement
of every 1C client/server combination is not established by this specification.
