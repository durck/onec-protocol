# Native bootstrap, discovery and context release

Profile: `native-18`; exact build 8.3.18.1208. Evidence: E01, E05, E08.
This chapter specifies the complete implemented successful bootstrap path,
including manager discovery. Unknown native alternatives remain explicit.
Positions are zero-based. Logical lists use text or binary serialization as
marked, followed by transport framing. Text UUIDs are unquoted tokens.

## Constants and notation

| Symbol | UUID |
| --- | --- |
| Root | `e23134a2-14ff-4160-ba5f-ccef04e3786f` |
| ConnectInterface | `7f58f27d-5ad8-43a1-aa1e-c982f41bed5c` |
| SelectInterface | `73b7d3a3-fe0b-4fdf-ba70-b74b3589ffc3` |
| ManagerContext | `4ecea39c-b4c9-453e-aed0-ce8a24cc71c2` |
| WorkerContext | `1e01900d-82fd-413a-bfe5-571a4c30a685` |
| UUIDType | `ae135932-4f94-44df-92c1-c91f15a92848` |
| BlobType | `671507fd-50a9-4b63-b70e-58b3d364f48f` |
| RemoteType | `3ace6d91-51bb-4344-9388-c8105ad4ad11` |
| ScalarType | `ba5edc33-8936-4212-9ed2-485e759203aa` |
| URLType | `dd5dcb0e-62b5-4d01-8305-6f87bbc772c6` |
| HTTPInterface | `bc15bd01-10bf-413c-a856-ddc907fcd123` |

Define `S(x) = ["S", x]`, `U = ["U"]`,
`UUID(x) = ["#", UUIDType, [1, x]]`.
Text bootstrap blobs use `TB(b) = ["#", BlobType, [B64(b)]]`, where B64 is
standard padded Base64 emitted as an **unquoted token**. For empty bytes the
last list is empty. This differs from binary SecurityBlob encoding.

## Inputs and identities

An application client needs an endpoint, infobase, advertised platform version,
application name and fresh security-provider state. The established baseline
advertises `8.3.18.1208` and `1CV8C`. These identify the client; they do not
discover the server version. Native V2 password login was studied with a mobile
application profile; changing identity does not establish compatibility.

The endpoint must include its role and explicit TCP port. A manager endpoint
and a worker endpoint are different inputs; a greeting alone does not identify
the role. The conventional manager port is 1541, but configuration can change
it. Do not infer a worker port by adding to that number. The bounded profile
uses IPv4 host:port endpoints and ASCII infobase/header names. Unicode names,
IPv6 routing and arbitrary cluster failover are not established by this subset.

Generate separate process, application client (`clnId`) and worker bootstrap
(`last`) UUIDs. Manager and worker in one discovery flow share the process UUID.
A manager uses zero `last`; a worker uses a fresh UUID. The computer-name string
is client input; documentation examples use a synthetic `CLIENT-EXAMPLE`.

| Connection | Context kind | Selector string |
| --- | --- | --- |
| Manager | ManagerContext | `RegMngrCntxt:RegMngrCntxt` |
| Worker | WorkerContext | `<infobase>:<infobase>` |

## Connect one context

1. Create provider state and obtain its introduction identity string.
2. Open TCP and read the exact greeting `53 F5 C6 1A 7B`.
3. Send `22 48 55 B5 || GUID_LE(processUUID)`, with no End marker.
4. Choose a request counter. The reference chooses 1000..30999 randomly; this
   is a client convention, not a native requirement.
5. Send the introduction inside the bootstrap request below.
6. Parse the incomplete reply and decode its introduction acknowledgment.
7. Drive the provider, sending each output token in a new bootstrap request
   and incrementing the RPC counter for every leg.
8. Require a completed reply, a remote context UUID and verified local provider
   completion. Retain the UUID as `target` for later RPCs.

The introduction uses the chosen initial counter N **without** first incrementing
it. Each provider leg, interface query, selection, HTTP RPC and release on that
connection uses the next counter. A new manager/worker connection starts its
own counter. Do not reset the worker counter after bootstrap or at login; do
not wrap uint32 after exhaustion. Reply field1 is not an echo requirement.

## Encoding dispatch

After removing record framing (and policy protection, when selected), use the
encoding of the outstanding operation, not a guess that all replies are binary:

| Operation | Request | Successful reply | Top-level reply suffix |
| --- | --- | --- | --- |
| Bootstrap, including introduction/provider legs | BOM + text | BOM + text | `,0` after the reply list |
| Interface query | BOM + text | BOM + text | None |
| Manager selection | Binary ListStream | Binary ListStream | Integer 2 after the reply list |
| HTTP RPC | Binary ListStream | Binary ListStream | Integer 2 after the reply list |
| Release | Binary ListStream | Binary ListStream | None |

Top-level binary values need the ordinary separator: a canonical suffix after
an already closed list is `83` (separator | integer 2), not bare `03`.
Do not append a second list-end marker for an implicit top-level container.
Exceptional replies are a separate parse result; never proceed using a target
from a failed bootstrap. Undecoded exceptions remain failures.

## Introduction bytes

```text
request = ASCII("V8IntroTicketReq") || u16LE(len(identity_utf8))
          || identity_utf8 || 00

ack = ASCII("V8IntroTicketAck") || u16BE(len(target_utf8)) || target_utf8
      || u16BE(len(package_utf8)) || package_utf8
```

The request's terminating NUL is not counted. Identity is nonempty UTF-8 without
NUL and fits the 16-bit byte length. Both acknowledgment strings are nonempty,
without NUL or trailing bytes. Length byte order differs between request and
acknowledgment. Here target is the **provider target string**, not the RPC UUID.
The implemented path accepts package `Negotiate`; static newer-build research
also identifies `Kerberos` selection without complete live coverage.

## Bootstrap request: text

Serialize this 28-element list with a UTF-8 BOM:

```text
[
  0, Root, counter, 4, ConnectInterface, 0, 1, 11,
  0, S(selector), 0, UUID(contextKind),
  0, S(application), 0, S(computerName),
  2, TB(token), 1, 0, S(advertisedVersion),
  0, S("1ru_RU"), 0, S("1ru_RU"), 0, UUID(last), 1
]
```

Keep context kind, selector, application and last fixed through the exchange.
Only token and counter change. The locale strings and unexplained numeric
positions retain observed values; no unsupported field semantics are inferred.

## Bootstrap reply: text

The top level is `R, 0`; R has 20 or 21 fields:

| Position | Required / interpreted value |
| ---: | --- |
| 0 | `1` |
| 1 | Opaque reply field; not correlated by the retained reader |
| 2 | `1` |
| 3 | `U`, or `["#", RemoteType, [remoteUUID]]` when complete |
| 4..9 | `11, 0, 0, 0, 0, 2` |
| 10 | `TB(returnedToken)` |
| 11 | `1` |
| 12 | `["B", completed]`, completed is 0 or 1 |
| 13..17 | `0, 0, 0, 0, 1` |
| 18 | `U` when incomplete; `["#", ScalarType, integer]` when complete |
| 19 | `0` for 20 fields; `1` for 21 fields |
| 20 | Extended-form field; semantics/type not established here |

The blob list contains zero or one unquoted Base64 token. Strictly decode it
to bytes. Before completion, positions 3 and 18 must both be U. The generic
reader permits U at position 3 on completion, but an application bootstrap
cannot proceed without a remote UUID. Completion on the introduction response
is rejected by the implemented path.

## Provider loop

Start with `Step(securityTarget, emptyInput)`. Send its output token even if
the provider already reports local completion. If the server stays incomplete,
require a nonempty returned token and feed it to the next Step. The reference
permits three provider legs; this is a local bound.

On server completion, feed any nonempty final token to the provider. It must
verify it, finish, and return **no further output token**. If the final token
is empty, local completion must already be established. The server flag alone
is insufficient, particularly for Kerberos AP-REP verification. See
[security providers](08-security-providers.md).

## Interface query: text

```text
request: <BOM>{0,target,counter,1,interfaceUUID}
```

The application HTTP/selection query reader accepts one three-element reply
list beginning and ending with `1`; its middle field is uninterpreted.
For example, a synthetic successful reply is `<BOM>{1,0,1}` followed by End.
Administration uses the last field as an interface method count under the
separate [administration contract](12-administration.md).

## Manager selection: binary

Bootstrap the manager, query SelectInterface, then send:

```text
[0, target, counter, 4, SelectInterface, 0, 3, 7,
 0, S(""), 0, S(""), 0, S(infobase), 0, S(application),
 0, UUID(clientUUID), 1, 1]
```

The reply is a 14-field list followed by top-level integer 2:

```text
[1, field1, 3, 7, 0, 0, 0, 0, 0, 1, workerEndpoint,
 1, ["#", URLType, [1, urlText]], 0], 2
```

workerEndpoint is the typed string `S("host:port")`, not a bare byte/string atom;
urlText inside URLType is a direct string atom. field1 is
not interpreted by the retained reader. Validate the returned endpoint. The
reference pins the manager host and uses only the selected port, so it does
not establish general multi-host routing. Release and close the manager, then
bootstrap the worker. Direct worker connections skip discovery. Query
HTTPInterface before the first HTTP RPC.

The selected endpoint has form `host:decimalPort`, with port 1..65535. The
returned URL field is parsed for structure but is not used as an HTTP URL,
provider target or session token by this profile. Binary selection strings in
the covered ASCII-name profile use byte-string atoms; the `S` discriminator is
also a byte-string atom. This does not resolve non-ASCII name compatibility.

## Release: binary

```text
request: [0, target, counter, 0, 0]
reply:   [1, field1]
```

The reply is one list, with no trailing scalar. field1 remains uninterpreted.
Release requires synchronized framing; never send it into a partially consumed
record. Application logout precedes release when a session is owned. Release,
socket closure and server-session disposal are distinct outcomes.

For exact build 8.3.27.2342, apply the [policy transform](10-native27-policy.md)
to each record payload before escaping. It does not replace bootstrap or login.
