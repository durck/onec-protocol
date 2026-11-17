# Authentication layers and session lifecycle

Evidence: E01, E08 and E09 in the [register](../evidence/README.md).
This chapter states observed boundaries and implementation cautions. It does
not specify complete authentication algorithms.

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
| Process/client GUID | Client identity for its documented connection role |
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
inventory, not a complete request/response contract; methods, parameters,
authentication schemas and error variants need dedicated chapters.

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
