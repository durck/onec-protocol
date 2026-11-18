# Authentication layers and session lifecycle

Evidence: E01, E08 and E09 in the [register](../evidence/README.md).
This chapter states shared ownership boundaries. Exact message algorithms and
orchestration are in [native login](07-native-authentication.md),
[SecPswd](09-secpswd.md) and [web login](11-web-authentication.md).

## Distinct layers

Outer bootstrap negotiates a provider and obtains a remote context. HTTP RPC
can carry a separate security token and completion flag. Application login
creates application session state. These steps must not be collapsed into a
single Boolean named "authenticated".

The native build-18 Windows-credential path has bounded real-server acceptance.
Password schemes and explicit NTLM/Kerberos providers have different synthetic
and observed coverage. HTTP Basic service authentication through a web
publication does not establish a native application session.

## Identifiers

| Value | Scope |
| --- | --- |
| Process GUID | Shared connection-flow identity for manager and worker |
| Application client GUID | Separate `clnId` identity, also used for manager selection |
| Worker bootstrap last GUID | Separate context-construction input |
| `target` | Dynamic remote RPC context/object |
| RPC counter | Native exchange sequence in the client |
| `vrs-session` | Application session header, treated as opaque |
| `Pragma: sync=N` | Application synchronous-call sequence |
| Infobase name/identifier | Selected application database |

Changing one does not substitute for renewing the others. A reconnect must not
silently reuse a stale remote target or infer a valid session from a previously
seen header.

## Application route inventory

The inspected native client uses `/e1cib/users`,
`/e1cib/additionalAuthSettings`, `/e1cib/login` and `/e1cib/logout`.
The login target includes profile/version and client parameters. This is a route
inventory; the complete implemented methods, parameters and authentication
schemas are described in [native authentication](07-native-authentication.md).

## Failure and cleanup

Track application logout, RPC context release and local socket/provider cleanup
separately. Local cleanup success does not prove that the server disposed of
its application session.

A bounded build-18 interrupted-upload study distinguished two outcomes:
the partial second request left 32 KiB staged, while a lost second response hid
an applied second block and left 64 KiB. A separate observer verified and removed
the owned stage. These observations explain why a missing response does not
prove that an operation failed to execute.

A later observation still listed the correlated uploader session through
60 seconds after a socket cut. That observation establishes neither an infinite
session lifetime nor eventual deletion. A listed connection number alone is not
proof of a surviving TCP connection.

Automatic replay of a credential-bearing or mutating application request is
therefore not a protocol rule. The product's one-shot precredential recovery
is a narrowly guarded client policy with bounded acceptance, not permission
to replay arbitrary exchanges.

## Server synchronization: exact-build static evidence

The 8.3.27.2342 server executor was inspected independently of the build-18
client. Let R be the incoming POST sync number, S the completed sync number,
and T the transaction sync number. The recovered branches are:

| Condition | Selected server action |
| --- | --- |
| R=-1 or S=-1 | Invoke the resource handler |
| R<S | Construct status 417 and throw |
| R=S | Read the retained result for R instead of invoking the normal handler |
| R>S and R=T | Construct internal status -5 and throw |
| Otherwise | Invoke the resource handler |

A missing retained result also constructs internal status -5; this is not an
HTTP status code. Normal completion records R, serializes the response into
temporary storage and calls a retained-result service. The transaction sync
update is separate. Exception cancellation follows a different cleanup path
without that normal completed-sync write.

Retention lifetime, eviction, persistence, crash behavior and reassociation
with a fresh TCP context remain unresolved. These static branches do not prove
exactly-once execution or permit automatic retry. The client marks its stream
unsynchronized before a write and restores it only after a complete expected
reply. Partial I/O prevents logout/release on that stream; local disposal still
runs. No generic cancel packet follows from a server function named cancellation.
