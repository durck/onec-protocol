# Independent public-spec client: bounded live acceptance

Date: 2026-11-18. Evidence: E19. This follows the
[earlier implementation audit](19-bootstrap-auth-audit.md).

A separate agent, started without the research conversation, implemented a new
Go client from public specification commit
`4286bdfe9562e4abf0dba3ccbbaa27657ee8545b`. The first frozen executable then
completed both direct-worker and manager-discovery Windows login/logout against
an existing authorized laboratory server. Both succeeded on their first attempt,
without a runtime source change or protocol correction between those runs.

This is evidence that the documented **bounded native-18 bootstrap and Windows
application authentication path** can be implemented without private research.
It does not establish all password modes, builds, providers or application work.

## Input and implementation boundary

The implementation inputs were the public checkout, its synthetic vectors and
public Microsoft/Go documentation. No private client/probe code, research,
captures, installed 1C binaries or captured authentication material were supplied
to the implementation agent. Windows SSPI supplies current-process outbound
credentials through the OS API. The runtime uses no 1C library or COM.

The agent had a fresh conversation and an explicit allowed-input inventory.
Filesystem isolation was instructional, not enforced by a separate OS sandbox;
this is not a formal legal clean-room certification. The coordinating reviewer
could access research and laboratory metadata, but did not supply private
protocol corrections or source fragments to the implementation. The reviewer
checked code/diagnostic scope, repeated offline tests, verified the frozen hashes
and ran the binary. A password-corpus coverage wording error was corrected before
freeze; it supplied no new protocol rule.

The client is a separate experiment. Its source, executable and tests are
retained separately; they are not copied into this specification repository.
The source/input manifests and source archive identify what was frozen.

| Initial artifact | Identity |
| --- | --- |
| Source commit | `4500da4ec4b76cf70e376b3c4e7c1e89ed99d6c6` |
| Executable SHA-256 | `ed746ebf31b6be9fc134ed38f2b709f89a35d7591a3e2babf4fbd9cff0bd4949` |
| Toolchain | Go 1.25.4, Windows/amd64 |
| Direct nonstandard Go dependency | `golang.org/x/text v0.31.0`, for full Unicode uppercasing |
| Build options | `-trimpath -buildvcs=false` |

The direct dependency uses Unicode 15.0.0. Its module graph also includes public
Go tooling dependencies; no private module or local implementation replacement
was present. The password vectors are ASCII, so this dependency choice does not
prove non-ASCII server compatibility.

## Offline verification before live input

The frozen client matched the public framing, scalar/GUID, HTTPFile and RPC
vectors, all 16 native session exchanges, eight complete SecPswd exchanges,
eight password-key cases and three native password HTTP orchestration flows.
The session fixture was driven through the client's actual state machine over
an in-memory transport with scripted provider tokens.

Additional checks covered malformed/truncated input, provider-final-proof
handling, owned-session cleanup, inner-state failures and stream retirement.
The final initial check reported 54 passing test events and no failures;
`go vet` passed. Separate bounded decoder fuzzing reported 278,075 binary and
239,086 text executions. These were local implementation tests, not this
repository's CI and not native-server behavior on those malformed inputs.

Offline tests found and fixed a stray-delimiter bug in the new text parser.
That was an implementation error, not a missing specification rule. The
implementation agent reported no missing public rule blocking the selected
Windows/legacy construction or the published SecPswd state machines.

## First live results

The laboratory's service configuration identifies build 8.3.18.1208; listener
ownership and process-parent relationships identify manager and worker roles.
The configured on-disk platform inputs match the native-18 hashes in the evidence
register. Direct inspection of the loaded-image paths was denied by the OS, so
those disk hashes must not be described as independently measured loaded images.

Only login and cleanup were requested. No application data operation, external
processing, account/service change, licensing change or OS change was performed.
No infobase username/password was supplied. The two runs used current-process
Windows SSPI and the same frozen binary, with fresh connection/provider/UUID state.

| Observed result | Direct worker | Manager discovery |
| --- | --- | --- |
| Manager bootstrap/query/selection | Not used | Passed |
| Manager release acknowledged | Not used | Yes |
| Worker bootstrap and HTTP-interface query | Passed | Passed |
| Inner provider exchange | Completed | Completed |
| Login HTTP status | 200 | 200 |
| Application XML and session validation | Passed | Passed |
| Logout HTTP status | 200 | 200 |
| Worker release acknowledged | Yes | Yes |
| Local provider/socket cleanup reported successful | Yes | Yes |
| RPC request/reply pairs | 10 | 16 |
| Logical HTTP operations | 2 | 2 |
| Process exit | 0 | 0 |

Logical HTTP operations are login and logout; repeated login PDUs during inner
negotiation count as separate RPC exchanges, not additional login operations.
Greeting and Control3 are outside the RPC-pair counts.

Results are recorded as status/phase/count booleans only. No raw tokens, session
headers, application documents, endpoints, account identities or workstation
names are exported. Logout/release acknowledgments and successful local cleanup
are observed separately; they do not prove that every server-side object has
already disappeared. The actual negotiated security package is not asserted
by the sanitized receipt, so these runs are not a specific Kerberos proof.
The [machine-readable acceptance summary](../evidence/independent-client-2026-11-18.json)
contains only the frozen identities, scope and these outcomes.

## Specification feedback and remaining boundaries

No wire-schema fix was needed to make either first live run succeed. The result
supports the earlier typed-selection and counter corrections instead of revealing
another blocking bootstrap gap. The follow-up makes two boundaries explicit:

1. [Inner security](07-native-authentication.md#inner-security-negotiation):
   application acceptance forbids later negotiation even when no inner provider
   was created at login. Preliminary HTTP completion and application acceptance
   are distinct state transitions. This is a review clarification; the two live
   runs both created inner providers and did not exercise the no-inner branch.
2. [Password casing](09-secpswd.md#password-derived-key): the full Unicode mapping
   version is not established for every shipped build. Record an implementation's
   choice and retain the non-ASCII compatibility limit explicitly.

Legacy construction and SecPswd V1/V2/V3 are implemented, but were **not
authenticated against this server in this experiment**. Offline checks covered
legacy wrapping/XML components and the complete published SecPswd exchanges.
No password credential was supplied for live use. Web publication, protected
native-27 transport, native 8.5, IPv6, general multi-host routing, arbitrary
exceptions, application execution and long-lived sessions remain outside this
acceptance. The two successful sessions are a bounded result, not a reliability
measurement or universal server/client compatibility claim.

## Follow-up state regression and repeated live checks

The implementation agent consumed the published clarification at
`5267ccd63c7a5fa6b65083646da21814013cdea7` before changing the client. Eight
synthetic negative subtests failed against the original runtime: late inner
negotiation and login-after-logout, each for Windows and V1/V2/V3. These failures
were found by state review and synthetic peers, not the laboratory server.

The correction records application acceptance for the lifetime of the connection.
It forbids those later transitions while preserving first negotiation at login
after a completed preliminary HTTP response. The follow-up source is commit
`e586136c7a023aa49435edc945d8593668e291d5`, executable SHA-256
`d789f3ddc8e12ee94ebbd5cfd723d1a846f88b39ad2ba214c1b9e42fb14bf3bc`.
Its 65 test events passed, including the eight previous failures and the allowed
prelogin transition; vet passed. Parser fuzzing was not repeated for a state-only
change. Initial tag, executable, source archive and receipts were preserved.

The coordinator verified this new frozen candidate and repeated the two live
Windows paths once each. Both again passed: 10/16 RPC pairs, login/logout 200,
acknowledged releases and successful local cleanup. There were **four total live
login attempts, all successful**, across the two frozen candidates. The later
two checks verify that the state correction preserved these working paths;
they do not turn the synthetic no-inner or password cases into live evidence.
