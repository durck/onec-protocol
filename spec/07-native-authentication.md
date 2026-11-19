# Native application authentication

Evidence: E08, E10. This starts after worker bootstrap and HTTP-interface query.
Application and transport-provider credentials are separate. Build-18 Windows
and legacy login have bounded live acceptance; native SecPswd V1/V3 have
original-code synthetic evidence, and V2 has schema/synthetic evidence.

## Common HTTP construction

Select Windows, legacy, V1, V2 or V3 explicitly before submitting credentials.
The [installed native-client audit](23-native-profile-selection.md) identifies
V1 C1 methods in exact builds 19.1522/22.1709 and V3 in 27.2342. These client
bindings are separate from the server's three-way credential dispatcher and do
not establish live server compatibility or an automatic negotiation algorithm.
No native on-wire scheme-negotiation parameter has been established here.
An advertised build, successful settings request, or bootstrap ScalarType value
does not by itself select SecPswd. The inspected native **server-state** selectors
0/1/2 mean V1/V2/V3 internally; they are not extra HTTP fields or a switch in the
ordinary build-27 client's C1 method. Do not send a guessed
`auth`, `code` or `version=1` parameter to force a native codec. The last form
is specifically a **web settings** selector. V2's studied application label
is `MobileClient`; other baseline examples use `1CV8C`.

These labels specify the studied request profiles, not a guarantee that a
server admits them. A live-confirmed native server/build selection for each
SecPswd variant is still missing from this independent-client experiment.
See [password acceptance prerequisites](18-client-guide.md#qualify-a-server-before-password-acceptance)
before treating codec coverage as server compatibility.

```text
/e1cib/login?vl=ru_RU&version=<advertised-version>&clnId=<client-uuid>
```

Percent-encode parameter values; encode spaces as `%20`. Password modes append
valueless `&nowa`. Later SecPswd requests append `&sid=<routing-uuid>` before
nowa. The common native headers are illustrated with synthetic inputs:

```http
1C-BaseLocation: e1c://server/server.example/ExampleBase
1C-ConnectString: Srvr="server.example";Ref="ExampleBase";
User-Agent: 1CV8C
1C-ApplicationName: 1CV8C
Accept-Encoding: deflate
Accept-Language: ru-RU
Accept-Charset: utf-8
Content-Length: <body-byte-count>
```

The profile/locale are the observed baseline, not universal requirements. Add
Content-Type and Accept per operation. Carry the entire HTTP message in
[HTTP RPC](04-http-rpc.md). Responses with Content-Encoding deflate use **raw
DEFLATE** in this subset, without a zlib header. Bound expansion and require a
complete stream. An optional UTF-8 BOM can precede XML, JSON or user-list text.

`ServerName` in the two 1C headers is a configured logical server identity,
normally the host without the TCP port. It is not the selected worker endpoint
or the provider's target/SPN. The embedded native profile sends an origin-form
request target, no Host header, no Connection header and no cookie jar. This
is encapsulated HTTP; an ordinary HTTP library may add headers or reject a
request that the native profile accepts. Web publication uses ordinary HTTP
Host/TLS handling instead. Empty POSTs still send `Content-Length: 0` and a
complete CRLF CRLF terminator. Header names are case-insensitive; body lengths
count bytes. No HTTP chunked transfer coding is used in this subset.

## Inner security negotiation

Start HTTP RPC with empty blob and securityFlag=0. A completed reply with empty
token can supply HTTP immediately. An incomplete reply with empty token asks
for inner security:

1. Create fresh provider state, independent of outer bootstrap.
2. Repeat the same logical HTTP PDU with V8IntroTicketReq and securityFlag=1.
3. Parse V8IntroTicketAck from the reply blob; completion here is premature.
4. Exchange provider tokens with securityFlag=1, repeating that PDU and
   incrementing the RPC counter every time.
5. Apply bootstrap's final-token verification; only then consume the HTTP PDU.

All later RPCs use the continuation zero after counter, including negotiation
legs. The first incomplete reply must not already contain a provider challenge
in this implemented path. A completed preliminary reply permits first inner
negotiation at later login, but established inner state is not silently
renegotiated. Final completion requires an HTTP PDU.

Repeat the **same bytes**, including credential body and query, during the
inner negotiation. Do not regenerate C1/nonces, change sid or increment
application sync for a token leg. Only the enclosing RPC counter, token and
security flag change. Ignore any intermediate PDU for application purposes;
only the completed reply's PDU can satisfy the HTTP operation.

For every subsequent logical HTTP operation after negotiation, start again
with an **empty SecurityBlob and securityFlag=0**, retaining the established
provider. Flag 1 belongs only to the introduction/token legs. If established
security receives incomplete status or a nonempty new token, fail instead of
restarting it. A provider returning localComplete=true still sends its output
once; the enclosing reply and any final peer proof must also complete.

Track application acceptance separately from whether an inner provider was
created. A completed **prelogin** HTTP response with no inner provider can be
followed by the first negotiation at login. Once application login is accepted,
every later HTTP reply must be security-complete with an empty token, even if
login itself completed without creating an inner provider. An incomplete logout
reply must not start a new provider exchange. Keep this boundary until the
connection is retired; logout does not authorize another login on that context.

## Windows/OS identity

Send an empty-body POST to login **without nowa**, Accept application/xml.
Complete inner negotiation if requested. Require HTTP 200, valid vrs-session
and valid application XML. No infobase username/password is submitted. The
server must map the proven OS identity to an infobase user.

## Legacy infobase password

1. Empty POST login with nowa and Accept application/xml; require HTTP **402**.
2. Construct `UTF8(username) || 08 || UTF8(password)`, standard padded Base64.
3. Apply native wrapping below, then percent-encode it.
4. Empty POST to the same target plus `&cred=<encoded-field>`.
5. Require HTTP 200, session header and application XML.

Wrapping: split Base64 into blocks of at most 1024 characters. Copy a block
of at most 64 unchanged. For a longer block, write chunks of at most 64 and
append CRLF after **every full chunk**, including an exact final full chunk.
A short final chunk has no CRLF. Thus a whole 64-character value has no newline;
a 128-character value has two CRLF pairs. URL escaping happens after wrapping.
This is reversible credential encoding, not hashing.

The local input contract excludes NUL, backspace, CR and LF, requires nonempty
username, and limits username/password to 256/1024 Unicode scalar values.
These size caps are local choices.

A complete 401/402 credential refusal without session or Connection header can
permit an explicitly requested fresh legacy attempt on synchronized state.
This bounded client policy does not authorize replay after partial writes, lost
replies, provider failures or SecPswd failures.

## SecPswd HTTP placement

Construct C1 using [SecPswd](09-secpswd.md). The native flow sends C1 directly,
without the legacy preliminary POST or web settings exchange. For every client
message, POST login with nowa, Content-Type application/json; charset=utf-8,
Accept application/json, and `{"cred":"<base64-packet>"}`. No vrs-session or
Pragma sync header is sent in this native login flow; sid carries routing.

Each response must be HTTP 200 with valid vrs-session. Intermediate responses
have a nonempty top-level cred string and no response/credAns4/credAns6. Feed
cred to the state machine, which must return another client message. Extract
sid from the header and require intermediate routing UUID stability.

Final JSON uses `response.credAns4` for V1/V2 or `response.credAns6` for V3,
inside the returned application object. Reject mixed challenge/application
documents, both proof names, top-level final proofs or early completion.
V1/V2 require two client POSTs, V3 three. Verify the final proof before declaring
authentication. Retain the final session for cleanup even if proof/document
validation then fails. An intermediate routing session is not an authenticated
application session. The accepted flow allows final session identity to differ
from intermediate routing identity.

## Session and application documents

Retain the **complete** vrs-session for later headers. To derive sid, the native
reader accepts a UUID or Base64 whose decoded first 36 bytes form a UUID.
Its Base64 normalization maps `-`→`+`, `_`→`/`, `.`→`=`, removes spaces/tabs and
adds padding where possible. Reject malformed input and the zero UUID.
The extracted UUID does not replace the complete session header.

XML root: `{http://v8.1c.ru/8.2/mngsrv/ws}response`. Required direct children
in namespace `http://v8.1c.ru/8.2/managed-application/core` are
infoBaseInstanceID, config, user and seance in the baseline validator. Do not
enable external entities or accept trailing documents. Optional scalar
assemblyVersion, config and rolesID select metadata caches; structured config
does not supply a scalar selector. JSON has corresponding application fields;
rolesID can be a string or number. Selectors are not credentials.

These validators have different scopes: native XML requires the four named
children but does not recover their entire nested schemas. Native SecPswd JSON
requires `response` and the correct proof; its metadata reader extracts optional
cache selectors and does **not** require all four XML fields. Web application
acceptance separately requires `user` and `config`. A synthetic minimal XML is:

```xml
<response xmlns="http://v8.1c.ru/8.2/mngsrv/ws"
          xmlns:c="http://v8.1c.ru/8.2/managed-application/core">
  <c:infoBaseInstanceID/><c:config/><c:user/><c:seance/>
</response>
```

It demonstrates the bounded validator only, not a real server's full metadata.
Do not infer session identity from `<seance>`: use the session header.

## Other endpoints

| Operation | Request | Expected result |
| --- | --- | --- |
| Optional users | GET common login query at `/e1cib/users`; Accept text/plain | 200, UTF-8 names separated by CRLF |
| Settings | Empty POST login with nowa, expect 402; GET common query at `/e1cib/additionalAuthSettings` | 200, decode body using HTTP headers |
| Logout | POST `/e1cib/logout`, empty body, full vrs-session, `Pragma: sync=N` | 200 |

The user-list flow requires completed security with an empty token and does
not initiate inner negotiation. User/settings operations precede login in the
retained implementation; availability depends on server configuration.

Application sync is distinct from RPC counter. Increment it for synchronous
session operations. Start the local sync counter at **0**; the first ordinary
synchronous application operation sends **1**. Login and its token legs send
no sync header and do not advance it. Logout alone uses `max(lastSync,1)+1`, so
immediate logout sends 2, and logout after operation 1 also sends 2. Never wrap
after 2147483647. These are the retained client's bounded sequence rules.
After logout, release RPC context and close both providers and the socket.
