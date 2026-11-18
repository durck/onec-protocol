# Route catalog, exceptions and unresolved surfaces

This catalog preserves source-established capabilities without inventing full
request grammars. A route's existence does not establish anonymous access,
authorization, live acceptance or compatibility across builds. Full grammars
are linked where recovered; other rows are discovery information.

| Method and route | Known purpose / boundary |
| --- | --- |
| GET publication root | Served VERSION discovery; separate HTTP publication profile |
| GET `scripts/mod_main_loader.js?sysver=...` | Select authentication flow from shipped source, without evaluating it |
| GET/POST `e1cib/login`, GET `e1cib/users` | [Native](07-native-authentication.md) and [web](11-web-authentication.md) orchestration differ |
| POST `e1cib/additionalAuthSettings` | Web credential challenge; native settings discovery also has a GET path |
| GET `e1cib/about`, POST `e1cib/logout` | Session-bound information and cleanup |
| GET `e1cib/ping` | Exact path and session/state guards; not anonymous health discovery |
| POST `e1cib/refresh` | Refreshes reusable values; mutates state, not a keepalive |
| POST `e1cib/messages` | 204 when absent, otherwise 200 UserMessages; may consume/replace message state |
| GET/POST `e1cib/desktop` | DesktopInfo reading/applying; presence flags applyUserSettings, onlyVisible, fullInfo |
| POST `e1cib/clearpool` | Base64 ticket, separate guarded operation; final response flags unresolved |
| `e1cib/start` | Authentication/start and redirect workflow; not a passive metadata route |
| GET `e1cib/types`, `sysenums`, `modules/defs`, `syspics` | Version/cache selectors and session permissions apply; complete metadata grammar not recovered |
| GET `e1cib/misc?cmd=localeInfo&locale=...&sysver=...` | Regional formatting, not infobase/user enumeration |
| GET `e1cib/logForm/id/<metadata-type>/<kind>?confver=...` | Form definition/cache selection |
| POST `e1cib/logForm?cmd=query` | Open form by id/key/url and typed named parameters |
| GET `e1cib/logForm?cmd=dataDescrWMethods&formId=...` | Form descriptors and server method metadata |
| POST `e1cib/logForm?cmd=callWOContext` | Existing no-context method |
| POST `e1cib/logForm?cmd=call` | Existing method with tracked form context |
| POST `e1cib/logForm?cmd=synchException` | Authoritative reconciliation after a method exception |
| POST `e1cib/logForm?cmd=synchronize` | Tracked proposal `{rKey,pchs}`, distinct from recovery |
| POST `e1cib/logForm?cmd=processFormsClose` | Form instance release |
| POST `e1cib/dlist?cmd=query` | [Bounded dynamic-list window](15-dynamic-lists.md) |
| GET `e1cib/dlist?cmd=getEmptyASSParam` | Typed list settings form-creation parameters |
| POST `e1cib/logForm?cmd=querySearchResult` | Type/choice parameters/text/count; values and search completion state |
| POST `e1cib/logForm?cmd=queryValuesList` | Typed values with choice/additional parameters |
| POST `e1cib/logForm?cmd=queryEnumSearchResult` | Enum choices, standard-list or generated mode |
| POST `e1cib/misc?cmd=post` | Typed current-item array, post/operational flags, optional form-close data |
| POST `e1cib/misc?cmd=delete` | Form/table, typed current item(s), mark/newMark, blocking and close data |
| POST `e1cib/misc?cmd=move` | Table, typed current array and typed parent |
| GET `e1cib/cmi/handler/<id>` or `handlerEx/<id>` | Form branch or command module |
| GET/POST `e1cib/cmi/command/<id>` | Command metadata / typed command invocation |
| `e1cib/uobjects?cmd=getFragment` | Source loader sends key/path/known row IDs for lazy fragments; full server grammar unverified |
| POST `e1cib/openDocument`, `extmd/processing`, `upload`, `opersec/insert` | [File and processing operations](16-files-and-storage.md) |

Navigation-link families DATA, APP and LIST (also SYSTEM in later clients)
are interpreted by the client and can become a form query's `url`. They are not
evidence for a direct `GET e1cib/data/...` endpoint. The selected sources do
not establish a general SQL execution or arbitrary BSL evaluation endpoint.
Basic-authenticated services, metadata resources and application sessions have
different access boundaries; one successful service request proves only that
surface.

Exact build-27 native handler inspection additionally registers `data`,
`metadata`, `tempstorage`, `pictureCollection`, `userSettings`, `favorites`,
`history` and `funcoption` families. Registration alone is not a usable request
schema. The data service identifier is `5e9726a8-ed76-4bb2-bee9-c37051cfe1b6`;
its general object/query contract remains unresolved. The broader source route
catalog includes static assets and templates, not hundreds of confirmed APIs.

`start` consumes fields including usr, pwd, authfailhandling, authfailredirectURL,
launchparameter, systemlanguage, localecode, zone, trmes and trmsnc. Its branches
use 1C-SafeSession, vrs-session2, VRSSESSION2 and Location; the complete workflow
is unresolved. `logout?cleanRMT` has a saved-authentication-state branch whose
full cookie contract is not decoded. A normal logout can return 200 when its
local session UUID is already null.

## Native exception records

Application native exceptions occur before embedded HTTP as a top-level list
followed by integer 2. Known shapes are:

```text
detail = [classUUID,message,[0],stack,source,0]
[2,classUUID,detail,0],2
[2,classUUID,[detail,extraText],1,extraText],2
```

The administration numeric-code variant is in
[administration](12-administration.md#errors-and-limits). Do not turn every
unknown shape into an authentication failure. Preserve classification while
keeping stacks, source text, raw session tokens and credentials out of public
diagnostics.

The released-context class `05a83ab3-14ea-4d74-aef1-3128b035c350` has a narrowly
recognized preliminary-login recovery path in build 18. It is not a generic
reconnect signal. Only a known pre-credential state, no application session,
no inner security context and complete old-context cleanup permit the bounded
single fresh attempt. No failed application method or password-bearing request
is automatically replayed by this rule.

## Form reconciliation

After a known application exception, the observed recovery sends
`{"root":{"key":remoteKey}}` to `logForm?cmd=synchException` with the session
and next sync number. Its returned authoritative descriptions, props, dynamic
state and presentation must fit the same supported schema. Recovery reports
the original failure; it does not convert it into method success. An unknown
transport outcome instead retires the context. Session cleanup remains owned.
The bounded application-exception classifier recognizes
`f08d92f8-9eb2-4e19-9dd9-977d907cec2d`,
`8d366056-4d5a-4d88-a207-0ae535b7d28e` and
`afd2b550-0d52-4e03-b448-15ed68d87ac5` in JSON exceptions. This is a selected
recovery policy, not a catalog of every server exception.

## Transport controls with incomplete semantics

The four escaped transport tokens are specified in [transport](01-transport.md).
Native inspection corroborates preservation of an Escape followed by an
unrecognized token. Control 2's payload can update a millisecond-valued field
or be ignored; no heartbeat interval follows. Control 3 can notify a callback
or set a distinguished flag; it is not a universally specified cancel action.

A conditional outer control header begins `23 49 56 B6`. The selected reader
requires at least eight bytes, a zero u16BE at offset 4, and a u16BE total length
at offset 6 not exceeding 32 or the available buffer. Its minimum valid total
length, payload semantics and enabling state remain unresolved. It is not a
fifth escaped-record delimiter and must not be recognized inside arbitrary data.

## Explicit gaps

The published bootstrap state machine and password message formats are complete
for their stated profiles. This is not a claim that every native service,
provider, challenge, platform build or application method is recovered. Still
open: full build-85 native envelope, build-27 live server proof, native policies
y/u, all NTLM/SPNEGO variants, MFA/external providers/password change, unknown
administration schema tails, general TypeDescription, complete form control and
layout grammar, client event execution, compiled BSL runtime, tree lists,
background search and interrupted storage reclamation. Native library internals
without an identified wire role are outside this protocol specification.
