# Server SecPswd profile selection

Date: 2026-11-19. Evidence: E24, S069. This chapter traces the Windows x64
8.3.27.2342 HTTP login implementation in installed vendor components. It closes
the missing login-mode-to-credential-state handoff recorded in
[chapter 23](23-native-profile-selection.md). All findings are static: no target
execution, authentication request or new live acceptance is involved.

## When this implementation selects V2

The login-state constructor initializes a 32-bit mode to **2**. The inspected
login reader can change it to **1** for an admitted older mobile request. It does
not assign 0 along this recovered ingress path. The credential dispatchers map
0/1/2 to V1/V2/V3, so the mode-1 assignment selects **V2**, not V1.

The mobile classification is a case-sensitive substring search for
`MobileClient` in the request's user-agent string, returned by interface slot
`+0x30`. The native login-state formatter names this exact field `userAgent`.
The reader tests that getter's field rather than directly looking up
`1C-ApplicationName`; it does not require whole-string equality.
The supplied `version` parameter is read from the request parameter map.
The concrete HTTP adapters feeding the getter were not independently traced;
this is not proof that every transport normalizes or populates it identically.

For the mobile branch, a successfully parsed version `V` must pass all of these
version checks before the selector assignment:

| Check | Rejected values |
| --- | --- |
| Minimum | `V < 8.3.9.1` |
| Maximum | `V > 8.3.27.2342` |
| Excluded range 1 | `8.3.19.1095 <= V < 8.3.19.9999` |
| Excluded range 2 | `8.3.23.0 <= V < 8.3.23.1838` |
| Excluded range 3 | `8.3.24.0 <= V < 8.3.24.1151` |

The exclusions are half-open: their upper endpoints are not rejected by those
particular checks. Comparisons use the platform's parsed version type, not
lexical ordering of version strings. Malformed-version parsing was not rederived.
The rejection helper throws; a decompiler's apparent fallthrough is not admission.

After these checks:

| Reached branch | Mode | Credential handler |
| --- | ---: | --- |
| Mobile, `V <= 8.3.25.9999` | 1 | V2 |
| Mobile, `8.3.25.9999 < V <= 8.3.27.2342` | 2 | V3 |
| Nonmobile, supplied version string exactly `8.3.27.2342` | 2 | V3 |
| Nonmobile, any other supplied version string | No accepted assignment | Version rejection |

The nonmobile comparison is exact string length/content equality; it is not
the mobile parsed-version comparison. These are conditions for this version
gate and profile choice. Earlier request-policy checks and later authentication,
session, licensing and application-start checks can still reject the request.
On the inspected reader's normal path the classifier/version block precedes
the credential and Windows-identity branches, with no surrounding mode bypass.
This table is not a list of supported server/client combinations.

## Assignment chain

The chain is established through concrete call arguments and interface bindings,
not just matching numeric constants:

```mermaid
flowchart LR
    A[Request classification and version checks] --> B[Login mode 1 or 2]
    B --> C[Server service argument 7]
    C --> D[Exchange selector at offset 0x150]
    D --> E[C1 parser dispatch]
    E --> F[S2 builder dispatch]
```

1. `vrsbase.dll` constructs the login state, then reads the request. Its mode
   is a **DWORD at byte offset `+0x4d8`**, not a byte-sized field.
2. For the reached SecPswd path, the credential string extracted from JSON
   `cred` and that mode are passed through server-service virtual slot `+0x188`.
   Both the inspected new-session and attach-session call sites use this mode.
   Empty credentials and other authentication paths need not call it.
3. The server service is requested through process slot `+0x70` with IID
   `f7399bd5-100e-4d0a-a5ce-f97810acfee9`. The backend constructor and interface
   map bind that IID to the subobject whose `+0x188` target is `0x1807f8390`.
4. At the caller, the mode is the seventh Windows x64 ABI argument, counting
   the implicit interface pointer and output-string storage. The next argument
   is a pointer to the phase, not another profile selector.
5. The backend reads that seventh argument and writes it to its local exchange
   state's DWORD at `+0x150` before invoking credential-service slot `+0x50`.
   If C1 parses, the same state reaches the phase-2 routine and slot `+0x58`.
6. That credential service uses IID `b0ff7cc9-64b1-418f-a6b9-39a1098499a1`.
   Its build-27 `backbas.dll` slots `+0x50`/`+0x58` read `state+0x150` and select
   V1/V2/V3 C1 parsing and S2 construction. Other values return failure.

The process-service registration is also connected: backend's descriptor at
`0x18195ff80` pairs class `a6184631-6903-43e8-a83b-df238e43064e` with the server
service IID and factory wrapper `0x180802ca0`. Module registration passes that
class/IID pair to process `+0x60`. On successful registration in `core83.dll`,
that method creates the registered class, queries the requested interface and
stores it under the IID;
process `+0x70` retrieves that stored pointer. This closes the static service
lookup edge rather than assuming the nearest vtable belongs to the getter.
Registration eligibility checks and an already-registered service can prevent
creation. This route does not establish the winner under every runtime
configuration or which modules every process loads.

The credential service has the same registration route: its `backbas.dll`
descriptor `0x1807df7f0` pairs the credential IID with class
`62cc8b00-d8d3-4a95-8b33-a11aa4285682` and wrapper `0x1800ffec0`, which reaches
factory `0x180100920`. The factory/interface map from chapter 23 returns the
credential table containing the C1/S2 dispatchers. Both service lookups are
therefore linked to concrete registered implementations in this static route.

For auditability, the critical exact-build coordinates are:

| Module | Anchor | Role |
| --- | --- | --- |
| vrsbase.dll | `0x18009dc00`, store `0x18009ddbf` | Constructor, default 2 |
| vrsbase.dll | `0x180093230`, store `0x180093cbc` | Login reader, older-mobile value 1 |
| vrsbase.dll | `0x18009b4b1`, `0x18009b951` | Server-service calls |
| backend.dll | `0x180804770`, `0x180805540` | Service factory and interface lookup |
| core83.dll | `0x180126e20`, `0x180127470` | Register and retrieve process service |
| backend.dll | `0x18144e938`, slot `+0x188` | Vtable targeting `0x1807f8390` |
| backend.dll | `0x1807f8525`, `0x1807f852b` | Incoming mode read and selector write |
| backend.dll | `0x1807f8556`, `0x1807f6008` | Credential C1/S2 dispatch calls |
| backbas.dll | `0x1800f88b0`, `0x1800fa3f0` | Selector-based C1/S2 dispatchers |
| backbas.dll | `0x1800f76f0`, `0x1800f9820` | V2 C1 parser and S2 builder |

Addresses use each PE image's preferred base. They are evidence coordinates,
not runtime addresses, wire fields or implementation dependencies. In the
backend, the frame base is `entryRSP-0x1d8` and the exchange begins at `RBP-0x20`;
therefore the write at frame offset `+0x130` is exchange offset `+0x150`.

## Selection persists with the exchange

On the phase-2 save path, backend copies the exchange into a session
context. The copy helper explicitly copies the DWORD selector as well as the
other credential state. A later request first attempts C1 parsing with its
request-derived mode. If that C1 parse fails and a saved context exists, backend
restores its handshake fields before the continuation parsers run. The saved selector
therefore governs that continuation.

The state-copy helper is `0x180810170`; its selector copy is at
`0x180810286`/`0x18081028c`. The phase-2 save calls it at `0x1807f62dc`, and the
continuation restore calls it at `0x1807f86d5`. Both use the same session-property
key. This is a continuation-state lookup, not a fallback from V3 to V2 after
a failed password proof. A fresh C1 can take the initial-request branch; this
finding does not redefine restart, retry or session-lifetime rules.

An independent client should choose its qualified profile before C1, keep
application identity/version consistent through the exchange, and retain the
matching challenges and routing state. Follow the matching message grammar in
[chapter 9](09-secpswd.md) and HTTP orchestration in
[chapter 7](07-native-authentication.md). Do not transmit the internal selector
or reinterpret phase byte 1 as a profile advertisement.

## Remaining boundaries

The static V2 choice here belongs to a server-side older-mobile branch. It does
not turn the fixed desktop client C1 methods in chapter 23 into a negotiated
choice, establish a V2 choice in those clients, or prove live V2 acceptance.
This study does not add a V1 choice to this exact ingress path; the presence of
a V1 dispatcher is not proof that this caller assigns selector 0.

The complete set of request adapters, malformed-version behavior, other
authentication/resume entry points, process configuration variants and runtime
admission remain outside the established scope. No general downgrade algorithm,
first format-changing release or cross-version compatibility matrix follows.
The previous web V1/V3 failures and licensing controls remain separate results.

The [sanitized summary](../evidence/server-secpswd-selection-2026-11-19.json)
records module hashes, native coordinates and evidence limits. Raw vendor code,
disassembly, private paths, debug tools and authentication material are not
included. The E23/S068 summary retains its historical unresolved handoff; E24
supersedes that limitation only for the exact chain described here.
