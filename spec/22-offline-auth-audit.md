# Offline bootstrap and authentication audit

Date: 2026-11-19. Evidence: E22, S067. This pass compared the published contract
with installed client resources and exact-build native code. It executed no 1C
client/server, sent no authentication request and changed no deployment. The
[stock web control](21-independent-web-auth.md#stock-web-client-control) remains
a separately dated observation; it is not an offline test result.

## Native policy-state correction

The build-27 reply marker is state-changing. Static inspection of rtrsrvc.dll
shows that the reply reader passes the current-level field to the marker decoder,
which writes numeric 0/1/2 for ASCII 0/1/2. It then checks that value against the
minimum. The later request constructor reuses the state and branches on the
current value; initial negotiation does not make that field immutable.

The earlier chapter-10 instruction to accept any marker meeting the minimum
while retaining the original effective level could therefore desynchronize a
client after a permitted transition. The corrected **fixed-c subset requires
marker 2**. It rejects other markers and retires the connection until complete
transition handling is implemented. This is intentionally narrower than the
native parser; it does not assert that native code rejects every marker 0/1.

| Initial effective level / minimum | Later reply marker | Fixed-c subset result |
| --- | --- | --- |
| c / n | ASCII 2 | Continue c with the completed exchange state |
| c / n | ASCII 0 or 1 | Reject unsupported transition, even though the numeric minimum alone permits it |
| c / c | ASCII 0 or 1 | Reject; also below the minimum |
| c / n or c | Any other byte | Reject malformed marker |

These are decision examples, not new executed native exchanges. Historical
original-client n/c query/release checks used stable levels. Live server policy
transitions, y, known-key identifiers and general native compatibility remain
unconfirmed. The greeting, provider introduction and application authentication
remain separate from the record policy.

## Password proof does not establish application readiness

The source comparison found a missing older status contract. The inspected
8.3.22.1709 client reads `passwordExpirationStatus`; the inspected 8.3.27.2342
client reads `passwordRequirementsComplianceStatus`. Chapter 11 previously named
only the newer field and allowed absence to mean zero. That could conceal an
older `passwordExpirationStatus: 2` response and falsely report readiness.

The [corrected table](11-web-authentication.md#credential-post-loop) records both
field names. Assuming a separately verified final credential result and valid
application/session envelope:

| Status fields | Bounded client's readiness decision |
| --- | --- |
| `passwordExpirationStatus: 2` | Change required; not application-ready |
| `passwordExpirationStatus: 1` | Expiry warning; continuation permitted |
| `passwordRequirementsComplianceStatus: 4` | Change required; not application-ready |
| Both fields present | Unsupported mixed contract; precedence is unknown |
| Either field has an unknown or malformed value | Unsupported contract |
| Neither field present | Retained-client compatibility default only; not proof of a universal stock default |

The old client routes nonzero status to a password dialog. Its status-1 warning
permits Ignore; cancelling status 2 performs session cleanup and returns to login.
The inspected stock handlers do not implement a simple absent-field-to-zero
rule. Unknown/malformed/mixed-field rejection in this draft is conservative
client policy, not a claim about every stock-client branch.

The complete 8.3.18.1208 and 8.3.19.1522 loaders contain neither field literal.
That bounded negative search does not prove the absence of server-side password
policy. In particular, both builds 19 and 22 use V1, so the credential codec alone
cannot select a readiness-field schema. Password-change submission remains
outside the specified minimal client.

## Settings state and decoding clarifications

Both inspected 22/27 loaders generate a C1 before settings and generate another
fresh C1 for login when settings did not return S2. The retained implementation
and public synthetic transcript instead preserve C1 byte-for-byte on that branch.
The specification now names that distinction rather than treating the retained
choice as universal stock behavior. No live acceptance of either no-S2 branch
is established here. Once S2 is accepted, preserve its matching state and nonces.

The web chapter now explicitly allows one optional leading UTF-8 BOM for both
successful and exceptional JSON responses. That instruction closes a portability
gap using the already retained build-27 exception observation; it is not newly
observed wire traffic. It also fixes clnId lifetime wording: reuse the identifier
through the login, independently of per-message cryptographic randomness. Stock
loaders can reuse a stored identifier; fresh-per-context is local client policy.

## Verification and remaining work

Four complete installed resource files were hash-checked. Their recorded loader
ranges matched the extracted files byte-for-byte, and fourteen bounded source
excerpts matched their full-loader ranges and hashes. The coordinator repeated
that verification and reviewed the relevant branches. A fresh Ghidra analysis
checked six native procedures and two assembly views on exact build 27; the
state-field write and minimum comparison agree in pseudocode and assembly.

REA parsed selected JavaScript excerpts, but full-loader analysis timed out and
the property-flow trace did not resolve a complete semantic path. The web claims
therefore rest on bounded manual source review and installed-byte equality, not
on a claimed complete call graph. Static branches are not runtime acceptance.

The [sanitized evidence summary](../evidence/offline-auth-audit-2026-11-19.json)
records hashes, byte/address anchors, classifications and tool limits. Original
code, decompilation, debug checks and private path mappings remain unpublished.

The later [native profile audit](23-native-profile-selection.md) establishes
exact-build client-method bindings; complete server-side selection remains open.
Still open: qualified live V1/V2/V3 acceptance; broader native profile selection;
native policy transitions and broader provider selection; full opaque bootstrap
fields; non-ASCII compatibility across Unicode versions; MFA and password-change
flows. This audit corrects specific implementation hazards and does not establish
a complete independent replacement for every client/server version.
