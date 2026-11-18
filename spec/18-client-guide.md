# Building an independent client

Implement the exact native-18 subset first. No proprietary library is required
to encode its records or password messages. Windows Negotiate can use the OS
API; portable providers can use established NTLM/Kerberos libraries with the
profile constraints in [security providers](08-security-providers.md).

## Minimal implementation sequence

1. Implement bounded incremental [record framing](01-transport.md). Separate
   stream reads from message boundaries and preserve unread bytes.
2. Implement [text and binary codecs](02-serialization.md), typed UUIDs,
   specialized HTTPFile chunks and SecurityBlob. Compare the
   [public vectors](../vectors/README.md) byte for byte.
3. Implement [bootstrap](03-bootstrap.md) as a state machine. Keep RPC counter,
   process UUID, client UUID, worker last UUID, target and provider target separate.
4. Implement [HTTP RPC](04-http-rpc.md), HTTP parsing and bounded raw-DEFLATE
   decoding. Distinguish incomplete security exchange from a complete HTTP reply.
5. Implement exactly one [native login mode](07-native-authentication.md), then
   add others. For password proof use [SecPswd](09-secpswd.md) and its complete
   synthetic transcripts; do not substitute a conventional CBC padding scheme.
6. Track raw session, routing UUID, application readiness, sync counter and
   cleanup ownership. Add one known application operation before broadening
   forms/data support. Close forms, logout, release and close providers/socket.

If starting with HTTP publication instead, use [web authentication](11-web-authentication.md)
directly. Native TCP bootstrap and the web settings sequence are not combined.

## Native connection algorithm

```text
processID = fresh UUID
clientID = fresh UUID

if endpoint is a manager:
    manager = connectContext(ManagerContext, managerSelector, last=zero,
                             processID, freshProvider())
    query manager for SelectInterface
    selected = selectWorker(manager, base, application, clientID)
    release manager; close its provider and socket
    workerEndpoint = validate selected endpoint under caller routing policy
else:
    workerEndpoint = configured endpoint

worker = connectContext(WorkerContext, base+":"+base, last=fresh UUID,
                        processID, freshProvider())
query worker for HTTPInterface
choose explicit login profile
perform login HTTP requests through rpcHTTP(worker, pdu)
verify final provider/password proof and application document
retain raw session; begin application sync state
perform supported application requests
close owned forms; logout if a session is owned
release synchronized worker context
close inner provider, outer provider and socket
```

`connectContext` implements the greeting/control/introduction/token loop in
chapter 3. `rpcHTTP` uses the initial envelope exactly once per worker context;
every later invocation includes the continuation zero, including inner security
legs. A repeated PDU during provider negotiation is part of that handshake,
not permission to retry an uncertain completed operation.

## Required state boundaries

| State | Can become ready when | Failure consequence |
| --- | --- | --- |
| Outer provider | Both local and remote completion, final proof validated | Dispose the connection context |
| Remote RPC context | Bootstrap returns usable target | No application RPC without target |
| Inner provider | Independent HTTP security loop completes | Do not reuse outer context as a substitute |
| Application routing | Valid intermediate session yields sid | Does not imply authenticated application |
| Application session | Proof, HTTP result and application document all valid | Still retain created session for cleanup |
| Form context | Descriptors, values, counters and versions understood | Retire unknown state; retain close ownership |

Do not issue multiple in-flight requests on the same synchronous context unless
a separately established multiplexing contract is implemented. Enforce size and
time limits as local policy; identify them separately from protocol limits.

## Independent verification

### Qualify a server before password acceptance

Implementing a credential codec does not identify a compatible server. Select
the transport and application profile before submitting a password, and record
the server build, selected profile and frozen client revision separately.

| Selected profile | Established prerequisite / remaining qualification |
| --- | --- |
| Native Legacy | The studied native-18 handler consumes the credential query after the preliminary 402. It does not establish native SecPswd support. |
| Native V1 | The published baseline uses `1CV8C`; a successful native V1 server/build combination has not been established by this independent experiment. |
| Native V2 | The studied label is `MobileClient`. That label alone does not establish that a particular server admits the profile or selects V2. |
| Native V3 | The baseline uses `1CV8C`. A newer advertised version does not implement a different record policy or establish server acceptance. |
| Web Legacy / V1 / V3 | Use the separate [publication contract](11-web-authentication.md). Loader observations identify request construction, not successful password acceptance. No web V2 orchestration is specified. |

The [native-27 policy](10-native27-policy.md) is a separate transport layer,
not a SecPswd selector. A client implementing only native-18 framing cannot
claim native-27 readiness by changing its version string. The specification
currently provides no native on-wire negotiation that chooses V1/V2/V3 and no
general matrix of live-confirmed server/profile combinations for these modes.
Keep this missing qualification explicit; do not invent a query selector or
try the same credential under successive schemes to discover it.

For each qualified test, distinguish these outcomes:

1. Transport/bootstrap completed and the selected HTTP operation was reached.
2. Each expected credential message was processed; any required challenge
   binding and final credential result were verified.
3. The returned application document and session were accepted. An intermediate
   routing session, HTTP 200, or an anonymous prelogin success is insufficient.
4. Logout, remote context release and local cleanup completed, recorded as
   separate outcomes.

Record an HTTP status only for the operation that actually received it. If
logout fails before receiving a response, its status is unknown; the preceding
login status is not a logout result. Do not collapse transport errors, unsupported
ingress, password refusal, additional-factor requirements, application-start
failures or cleanup failures into a single "wrong password" outcome. A license
error before final acceptance does not by itself prove or disprove the password
exchange. Record the last validated phase and stop without automatic fallback.

### Public corpus and live evidence

The public corpus contains synthetic framing, GUID, integer, raw-byte, text,
HTTPFile and HTTP RPC encodings, one worker introduction, eight password-key
vectors and eight complete V1/V2/V3 exchanges. The new authentication transcripts
were constructed by an independent encoder and compared with the retained
implementation's entire state machine. They contain deliberately public toy
credentials, deterministic randomness and decoded packet bytes.

The [bootstrap/authentication audit](19-bootstrap-auth-audit.md) adds a complete
16-exchange native session and seven native/web HTTP flows. It records concrete
corrections, including typed selection endpoints and first-operation sync=1,
plus the limits of the source-independent executable check. That audit did not
include a fresh-agent implementation or a live server.

It now has a [separate follow-up](20-independent-client.md): a fresh
public-spec Go consumer completed the direct-worker and manager-discovery Windows
login/logout paths on its first live attempts. Other profiles retain their
previous evidence limits; the initial synthetic audit is not reclassified as live.

The [web continuation](21-independent-web-auth.md) separately passed Legacy
login/logout on build 18 and processed V1/V3 challenges on builds 22/27 before
final application-start failures. Keep its web-only scope and failed cleanup
outcomes distinct from the successful native Windows path.

Conformance data can establish encoder/decoder agreement. It cannot establish
server acceptance, negotiated permissions or a full 1C replacement. Use the
[version matrix](06-version-matrix.md) to choose which statements your client
can make, and reject unsupported states explicitly.
