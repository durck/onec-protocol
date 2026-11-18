# Independent web password-client experiment

Evidence: E20, S065. Date: 2026-11-18. The same public-input Go client from
[the native experiment](20-independent-client.md) gained a separate ordinary
HTTP publication transport. **Legacy web login/logout passed on 8.3.18.1208.**
V1 on 8.3.22.1709 and V3 on 8.3.27.2342 processed real server challenges but
failed at the final application-start request. Neither completed password
authentication is claimed. V2 was not attempted against a server.

## Implementation and server boundaries

The implementation agent used only the public specification, its synthetic
vectors and public Go/OS documentation. The coordinator alone prepared the
fixtures, supplied their existing credentials through private stdin and ran the
frozen binaries. Private research, server components and account data were not
implementation inputs. The instructed input boundary remains distinct from an
OS-enforced sandbox. No 1C runtime or COM is used by the independent client.

The test publications used IIS Express and disposable file infobases. Two
existing fixtures were reused; the third used a closed, byte-verified copy of
an existing fixture with the same account definitions. A matching web extension
with a valid Authenticode signature was extracted from an existing distribution
without running its installer.
No existing account settings, native service, platform installation or licensing
configuration was changed. File-publication testing does not establish native
cluster bootstrap or imply that a thin-client installation alone is a server.

Loaded web-extension/runtime module identities were inspected; hashes refer to
their corresponding files, not reconstructed in-memory images. The extension
hashes were:

| Build | wsisapi.dll SHA-256 |
| --- | --- |
| 8.3.18.1208 | `752ac230b36786ead9b41344cba4c19b045e574539ad081c564397737bd14431` |
| 8.3.22.1709 | `126421fd61db23528ccafb67c5f0ca301255c2d7fad47d854116b2701a8bb433` |
| 8.3.27.2342 | `ec749117522500fad5b4b909f7d88a24ad171161ccbbefe573655116b2b8e92f` |

No binary, infobase, credential, request log or raw authentication response is
bundled. Publication access/tracing logs were disabled and process request output
was discarded. Only fixed diagnostics, counters, statuses and artifact hashes
are exported in the [structured summary](../evidence/independent-web-auth-2026-11-18.json).

## Frozen candidates and verification

| Artifact | Initial web candidate | Discovery follow-up |
| --- | --- | --- |
| Public specification input | `3009ecfa350b37d5245a8cb610938ae8073b976c` | `db256593b59597d89d7f213f43c408523e88748d` |
| Source commit | `907ec66c68fa1d336409441b8ad92c625d9850f6` | `b9c81be5e666738280fdc622696016a0a0774f37` |
| Executable SHA-256 | `e30e792b35b250a0934c5741c9856f0a9d8a25f0184896eeb8dfcb460f532df1` | `cf428b48268e050c1741dd3c951916a971cd3f0a872731bc44d19c317a80de40` |
| Public-input manifest SHA-256 | `ec88b2a1d9efb5e77c08c08848dd79caa73489011feaba23899d51cfdeafb37e` | `0952772e24933e856a663b70e080eeb760ce4afab06a769a8b4b2210c53a236c` |
| Passing local test events | 150 | 164 |

The coordinator repeated tests and vet and reproduced both executable hashes
using Go 1.25.4 Windows/amd64 with `-trimpath -buildvcs=false`. All preceding
candidates were preserved. The additional tests exercised complete public web
V1/V3 flows with and without settings S2, Legacy, session rotation/precedence,
factor/refusal handling, proof/application failure cleanup, cookies, URL scope,
redirect refusal, body bounds and the discovery regressions. These tests remain
outside this repository; its CI does not run the laboratory.

The client selects Legacy/V1/V3 explicitly, requests identity encoding, validates
TLS normally and follows no redirects. Unexpected response encodings fail.
These are local implementation limits; native raw-DEFLATE behavior does not
establish web response-compression semantics. No alternate scheme is tried after
failure. The existing Unicode 15.0.0/ASCII-vector evidence boundary remains.

## Recorded invocations

Six separate invocations were retained, four of which submitted a credential
packet. The two additional invocations stopped before credential submission.
No credentials were automatically replayed. Later runs used fresh state.

| # | Candidate / profile | Observed outcome | Logical HTTP requests |
| --- | --- | --- | ---: |
| 1 | Initial / Legacy, build 18 | Entry HTTP 200, then VERSION parsing rejected the semicolon-free declaration; no login request | 1 |
| 2 | Initial / V1, build 22 | Entry accepted; credential-free prelogin did not return within the 15-second local deadline | 2 |
| 3 | Initial / V3, build 27 | Client processed two challenges; final C5 HTTP 500 matched the published exact-build license-failure pattern; logout HTTP 500 matched absent-session pattern | 6 |
| 4 | Initial / V1, build 22 | Client processed one challenge; final C3 HTTP 500; logout HTTP 500; cause not classified by that candidate | 5 |
| 5 | Follow-up / Legacy, build 18 | Application JSON/session accepted, login HTTP 200, logout HTTP 200, local cleanup passed, exit 0 | 3 |
| 6 | Follow-up / V1, build 22 | Client processed one challenge; final C3 HTTP 500 and logout HTTP 500 structurally matched the documented license/absent-session shapes | 5 |

Invocations 3-6 used an explicit 60-second per-request and 240-second process
deadline. Invocation 4 followed the credential-free timeout; both the deadline
and server warm-up state differed, so this does not isolate the timeout's cause.
Invocation 6 was one separately declared diagnostic follow-up, not a hidden
retry. Its fixture process was restarted before the run to renew the owned
process deadline, using unchanged publication configuration and account settings.

V1's later structural match is not a new universal reason-code interpretation:
no error descriptions were emitted by the independent client. It supports
comparison with the [known failure structure](11-web-authentication.md#observed-application-start-failure-separate-from-credential-refusal),
not an independently verified licensing diagnosis on every build. Both V1 and
V3 lack a verified final credential result and accepted application document.
Successful challenge processing or HTTP 200 at an intermediate stage does not
close that gap. No non-ASCII password, MFA or password-change flow was tested.

## Specification feedback and cleanup

The first Legacy invocation exposed an overly strict entry parser. The observed
build-18 VERSION literal had no semicolon and was followed by another `var`
declaration across a line boundary. The form was published in chapter 11 before
the agent consumed it. Red regressions reproduced the parser rejection and a
second diagnostic bug: discovery HTTP 200 was attributed to login even though
login never began. The follow-up preserves line boundaries, rejects continued
expressions and records discovery/login statuses separately. Credential bytes,
cryptography and request sequences were unchanged.

The later V1 diagnostic added only fixed exception-shape labels. Matching a
shape does not broaden the initial exact-build/phase cause classifier. Complete
application validation and final credential result are separate from cleanup;
the optional password-compliance field is now named explicitly in chapter 11.

Local client cleanup was reported successful in all six invocations. Only the
Legacy run acknowledged successful application logout. The V1/V3 logout failures
remain failures; stopping an owned fixture does not retroactively make those
operations successful. All task-owned publication processes/supervisors and their
listeners were absent after cleanup; the existing native server remained running
with its original process identity. No application data operation was performed.

Remaining work: completed V1/V3 application acceptance on qualified deployments,
native V2 admission/transport qualification and live acceptance, native password
coverage for this independent client, and broader provider/version behavior.
This experiment establishes a bounded web result, not a complete 1C client.
