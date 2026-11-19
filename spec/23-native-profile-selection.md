# Native password-profile selection

Date: 2026-11-19. Evidence: E23, S068. This is an offline study of exact installed
Windows x64 client components. It identifies client-side credential construction,
separately from server admission and live login acceptance. No 1C executable was
run and no authentication request was sent during this study.

## Client construction and server dispatch are different decisions

The native credential service exposes both client and server methods. In the
inspected build-27 service table, client slot `+0x48` directly references the V3
C1 builder. The adjacent server slots `+0x50` and `+0x58` instead dispatch C1
parsing and S2 construction using a 32-bit field in the supplied exchange state.
That field is at byte offset `+0x150` in this exact build:

| Internal server-state value | Handler |
| ---: | --- |
| 0 | SecPswd V1 |
| 1 | SecPswd V2 |
| 2 | SecPswd V3 |
| Other | The two inspected dispatchers reject the operation |

These values do **not** choose the ordinary client's C1 builder through that
client slot. Nor are they a field to append to HTTP or to a credential message.
The first decoded credential byte `01` is a phase marker shared by the variants,
not a protocol-version advertisement. Legacy encoding and Windows identity are
separate login paths, not additional values of this three-way switch.

## Exact installed client observations

The following classification comes from native instructions and interface-table
bindings, not from the web loaders shipped with the same versions:

| Installed build | Credential service's client C1 method | Construction |
| --- | --- | --- |
| 8.3.19.1522 | V1 | Phase 1, UTF-8 username length/data, 32-byte echo challenge, four-byte checksum and direct Base64 |
| 8.3.22.1709 | V1 | Same V1 field order and direct Base64; no pair of wrapped client-key alternatives |
| 8.3.27.2342 | V3 | Fixed-envelope V3 C1 builder, followed by the V3 client continuation methods |

These are exact-build observations, not ranges of compatible server versions.
They do not establish the first release where the format changed or select a
profile for an uninspected patch. The existing native-18 **independent-client**
Legacy acceptance is a separate result; this pass did not reclassify every
stock native-18 login branch.

The older native login routines actually call these C1 methods, put their output
into JSON `cred`, and send it to `e1cib/login`. This closes the distinction between
an exposed but possibly unused service method and its inspected password-login
caller. The native continuations also show why interface offsets must be bound
to an exact build:

| Build | C1 method slot | S2/client continuation slot | Final-check slot |
| --- | --- | --- | --- |
| 19.1522 | `+0x48` | `+0x60` | `+0x78` |
| 22.1709 | `+0x48` | `+0x60` | `+0x80` |
| 27.2342 | `+0x48` | `+0x60`, then additional `+0x80` | `+0xa0` |

These are static analysis coordinates. The new older-build classification uses
field order, checksum placement and resolved Base64 calls; it is not a new
byte-exact execution oracle or an independent rederivation of every primitive.
The credential algorithms retain the evidence boundaries of chapter 9.

For build 27, the native login routine obtains the credential service and calls
its client methods at `+0x48`, `+0x60`, `+0x80` and `+0xa0`. The fresh-password
branch sends C1/C3/C5 and checks the final result. Its inspected body contains no
V1/V2 selection or failed-proof downgrade loop. Other identity/resume branches
exist; this is not a claim that every invocation enters the password branch.

The same service IID is used by the inspected older client components. A matched
IID alone is insufficient: the constructors, interface maps and actual method
targets were checked, and the older C1 bodies were classified from field order
and resolved encoding calls. A method name or the presence of all three server
parsers would not establish this client behavior.

## What this means for an independent client

Select a supported password profile before constructing C1. Exact client
observations can inform an explicitly named compatibility profile, but cannot
replace qualification of the target server and application mode. Keep the
[transport policy](10-native27-policy.md), application identity and credential
codec as separate settings.

Changing the advertised HTTP version does not change a compiled service-table
entry. A successful bootstrap, preliminary 402 or settings response likewise
does not prove that a particular SecPswd codec will be accepted. Do not send a
guessed selector, infer a variant from phase byte 1, or automatically try another
variant after an authentication or application-start failure.

After choosing the profile, follow [native HTTP orchestration](07-native-authentication.md)
and the corresponding [credential state machine](09-secpswd.md). Preserve the
matching challenge state, verify the final result and application document, and
record cleanup separately. This audit changes the evidence for profile choice;
it adds no new wire fields, credential packets or live acceptance.

## V2 and server-side selection remain bounded

V2 server parsers/builders are present, but this pass did not establish a V2 C1
choice in the three inspected desktop client-service bindings. That is not proof
that V2 is unused by all native or mobile clients. `MobileClient` remains a
studied compatibility label, not a universal command to force V2.

The retained build-27 HTTP handler initializes a separate login-mode field to 2
and assigns 1 in an admitted older MobileClient-related branch at version
`<= 8.3.25.9999`. This is not yet a complete traced assignment to the credential
state's `+0x150`: equal numeric values and a version comparison do not prove the
cross-module handoff. Admission checks also remain applicable. Do not turn this
partial server observation into an unconditional client/server compatibility
matrix or a new query parameter.

Still unresolved: the complete server-side origin of the dispatcher field,
native V2 client selection and live admission, every upstream login/resume
branch, the precise release transition and general cross-version interoperability.
The earlier V1/V3 web failures and licensing controls remain separate evidence.

## Provenance and limits

The [sanitized summary](../evidence/native-profile-selection-2026-11-19.json)
records exact module hashes, native addresses and evidence boundaries. Ghidra
pseudocode and independent GNU objdump instructions corroborate the inspected
factory/interface bindings. Retained build-27 login and credential observations
were reused only after checking their installed-artifact identities.

Initial MCP analysis calls exceeded the tool's 300-second limit. An owned local
REA session with a longer wait completed the static analysis; target execution
was never substituted for it. Analysis sessions were closed after evidence export.
Raw disassembly, decompilation, local verification tools and private path mappings
are not bundled. Static construction is not successful authentication.
