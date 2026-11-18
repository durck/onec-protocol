# HTTP(S) publication and web authentication

Evidence: E13. This profile uses a normal HTTP(S) publication directory, not the
native worker socket. Native RPC framing, HTTPFile and Control3 are absent.
Paths below are relative to the publication directory. Keep requests/cookies
within the explicitly selected publication and do not forward credentials on
unreviewed redirects. TLS policy is a client/deployment choice.

## Discover the served contract

GET the publication entry directory and obtain its served `var VERSION` value.
Use it for web version/cache selectors; it is separate from the native client's
advertised version. For automatic selection, GET
`scripts/mod_main_loader.js?sysver=<served-version>` and inspect, without
executing it, the actual login construction/proof-consumption code.

| Recognized contract | Selection |
| --- | --- |
| Legacy login target construction and appended cred query | Legacy |
| Login route and parsed `response.credAns4` access | V1 |
| Login route and parsed `response.credAns6` access | V3 |

Do not select from comments, string mentions or a version-number threshold.
Unknown/mixed contracts require explicit investigation, not password fallback.
The compared clients selected legacy for 8.3.12/17/18, V1 for 19/22 and V3 for
the three studied 27 builds; these are source observations, not all-build rules.

A minimal client can require an explicitly selected legacy/V1/V3 contract and
omit automatic source recognition. To reproduce the narrow recognizer, tokenize
JavaScript as data: skip comments and opaque regex/template contents; distinguish
string literals from executable identifiers. Modern selection requires the
literal route `e1cib/login` and the executable token sequence
`JSON.parse(identifier.responseText).response.credAns4` or `credAns6`.
Legacy requires the literal `e1cib/login?version=` and construction equivalent
to `identifier += "&cred=" + expression`. Conflicting proofs, modern plus legacy,
or an unrecognized executable `credAns...` variant are unsupported. This is a
bounded recognizer, not a general JavaScript evaluator or server negotiation.

All authentication requests in this profile use Accept application/json and
Content-Type application/json; charset=utf-8. Keep cookies in the same session.

## Login query

```text
version=<served-version>&nowa=&nooida=&clnId=<fresh-client-uuid>
```

The HTTP library may serialize flags as empty values. The retained web query
does not add native `vl=ru_RU`. Legacy adds cred containing standard padded
Base64(UTF8(user) || 08 || UTF8(password)), URL-encoded, without the native
1024/64-character wrapping rule. POST `e1cib/login` with an empty body, capture
vrs-session, and require HTTP 200 plus a valid application JSON document.

## Modern prelogin and settings

1. POST login with the common query and empty body, without credentials.
2. Capture any vrs-session. HTTP 200 can mean an anonymously accepted
   application; it does **not** validate the supplied username/password.
3. Otherwise require HTTP 402 and a recognized exception document. Application
   `seanceID` can be supplied directly or through inner exceptions, or a valid
   application session can already have come from vrs-session. Preserve case:
   lowercase `seanceid` is a distinct second-factor identifier.
4. Reject a second-factor challenge for this password-only profile. Establish
   a valid nonzero application routing UUID from the session.
5. Generate SecPswd C1. POST `e1cib/additionalAuthSettings?sysver=<served-version>`
   with body `{"cred":"<C1>"}` and the complete vrs-session header. V3 also
   adds query parameter `version=1`; this is a settings-format selector.
6. Require HTTP 200. If the settings JSON contains cred, process it as S2 and
   obtain C3. If cred is absent, keep the same C1 for the next login request.
   This observed alternative does not select a weaker credential scheme.

For the retained session-selection rule, first capture a nonempty vrs-session
header; an absent header leaves the current session unchanged. Then inspect the
402 exception and its `inner` chain from outermost to innermost, replacing the
full session with each nonempty exact-case `seanceID`. Thus the deepest supplied
application session takes precedence. Require a valid nonzero projected UUID
after this process; a 402 with neither source is not enough. The bounded reader
uses at most 16 exception levels. Stop on reason=-14 or **any nonzero** reason_2fa;
unrecognized factors are unsupported, not ordinary password continuation.

`seanceID` can already be the complete opaque/encoded session, not just a UUID.
Use that complete value in the settings header, then derive sid for login.
After settings and every login response, capture a replacement vrs-session
when supplied and recompute sid for the next request. Unlike the native SecPswd
profile, the retained web flow does not require the intermediate routing UUID
to stay constant. Do not overwrite a retained session with an absent header.

The 19/22 **POST-with-C1** settings contract supersedes an earlier GET description.
Both V1 and V3 use POST here. If settings supplies S2, V1 needs one subsequent
login POST (C3), and V3 needs two (C3,C5). Without S2, the counts are two
(C1,C3) and three (C1,C3,C5), respectively. Do not create a second codec or fresh
C1 when moving from settings to login.

The reference accepts a raw UUID session or Base64 (standard/URL alphabet,
padded or unpadded) whose first decoded 36 bytes are a UUID. Preserve the full
header value for session-bearing requests; use only the extracted UUID in sid.

## Credential POST loop

POST `e1cib/login` with common query plus `sid=<routing-uuid>`, JSON
`{"cred":"<current-client-message>"}`. In the retained web flow these login
POSTs route via sid and cookies, without adding vrs-session; settings/about/logout
requests carry the header. Capture a returned session after each response.

HTTP 200 plus top-level cred continues [SecPswd](09-secpswd.md). A message cannot
contain both cred and response. Final application proof is response.credAns4
for V1 and response.credAns6 for V3. Verify it before accepting the application.
The reference bounds login-loop requests to three after settings; the exact
number depends on whether settings supplied S2.

The accepted application JSON has one non-`#` payload object with user and
config fields and a retained session. Metadata keys beginning with `#` are
separate from the application payload. Exceptions are not application objects.
Password compliance statuses 2/4/6 require password change; 0/1/3/5 allow
application readiness in the observed contract. Preserve the distinction
between authenticated identity and application readiness.

## Failure and cleanup

The studied exception parser recognizes classes
`580392e6-ba49-4280-ac67-fcd6f2180121` and
`17ef6ddf-df7f-481d-b419-86b782720c2b` for **bare** exception objects (case-insensitive
class-ID comparison). A document with a non-null top-level `exception` object
is the wrapped form; the retained parser does not require either class ID in
that wrapper. Its session, reason and factor fields still need validation. reason=-14 or
reason_2fa=1 signals a second factor; reason_2fa=3 an external second factor.
This draft records those signals but does not specify factor submission.

Logout is POST `e1cib/logout`, body `{"root":{}}`, full vrs-session. Accepted
statuses are 200/201/204; a nonempty body must decode without an exception.
On success discard session and cookies. Local Close alone does not perform
logout. Cleanup may be necessary for preauthentication sessions as well.

GET `e1cib/about` with vrs-session is an observed informational operation; HTTP
200 plus response JSON is accepted. It does not initialize application forms.

## HTTP Basic is a separate service profile

Basic authentication of OData/service-description routes does not create an
application vrs-session. Keep Basic credentials separate from web login proof
and cookies. Live Basic evidence exists for builds 18/27; live application login
evidence is build-18 legacy. Web V1/V3 have synthetic acceptance; successful
build-27 V3 application login remains unconfirmed.
