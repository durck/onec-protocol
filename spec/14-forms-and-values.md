# Application forms and typed values

Status: bounded build-18 live operations plus shipped-client static schemas.
This chapter describes server form communication, not a browser renderer or a
complete client-side BSL runtime. Opening a form can execute application hooks.

## JSON transport and type system

Application POST bodies normally wrap the ordered payload as `{"root":payload}`.
Send the raw `vrs-session` token, `Pragma: sync=N`, JSON content type with UTF-8,
and `Accept: application/json`. Advance the sync number before each attempted
request; it is independent of the native RPC counter. Replies select one
non-`#` top-level application object (normally `response`); do not pick the first key from an
ambiguous result. Check application exceptions even when HTTP status is 200.

**Member order matters in the observed JSON-to-XDTO reader.** Do not
alphabetically sort request members. General JSON object order independence
does not describe this implementation boundary.

| Value | Type UUID | Payload |
| --- | --- | --- |
| String | `9b6abf8b-0173-48e5-b0a0-83b21fcf63c5` | JSON string |
| Number | `b0be78f2-0ee6-4d31-a3bb-77dd32ba5bec` | Send decimal string; receive string or exact JSON number |
| Boolean | `5d4125ad-f6e7-4313-be32-f71d0ab60915` | JSON Boolean |
| Date | `aae38c48-a877-411c-a6d3-fbaa1f83c4bd` | `YYYY-MM-DDTHH:mm:ss`, civil time without offset |
| Undefined | `ee8d3e7c-f930-4a76-8aad-4ff9083a6ea6` | Literal `undefined`, outside ordinary JSON |
| Null | `af40a278-63bc-478e-91a8-19e0d16b10b5` | JSON null |
| Binary | `87126200-3e98-44e0-b931-ccb1d7edc497` | Standard base64 |

An individual typed value uses `{"#value":typeUUID,"value":payload}`. Undefined
can omit `#value`; parallel type arrays can use the literal string `undefined`.
An absent type in a declared structure instead selects its descriptor default.
Never collapse absent, Undefined, Null and empty reference. Parse the literal
Undefined extension without evaluating JavaScript.

Number decoding must preserve lexical decimal precision, including integers
above 2^53. Date's empty value is `0001-01-01T00:00:00`; no UTC offset is implied.
Confirmed category-0 reference types serialize a UUID string or `""` for an
empty reference. Enum references instead serialize a metadata value code;
record references use Recorder/LineNumber. A UUID-looking type identifier alone
does not establish the ordinary reference codec.

Qualifiers `digits`, `precision`, `unsigned`, `length`, `fixed`, `dateFraction`
constrain local assignments. String length is measured in UTF-16 units. Native
setters can truncate, round or clamp; directly loading a wire value is a separate
operation. A strict client should explicitly reject unsupported coercions.

## Opening, inspecting and calling

POST `/e1cib/logForm?cmd=query` with payload members in this order:
`key,prms,mw,mh,cw,ch,pres,trdata,cf`. For an attached processing object,
`key=reference+".form"` or `reference+".form."+formName`. `prms` is
`{"prm":[{"name":name,"#val":type,"val":value},...]}`. Empty parameter lists
are real arrays. The other default fields are four numeric zeros and three
false Booleans, in the listed order.

Retain returned `id` (instance), `fid` (form identity), `remotekey`, version
fields and descriptions independently. Fetch
`GET /e1cib/logForm?cmd=dataDescrWMethods&formId=<escaped id>` when method
metadata is not embedded in `moduleIds`. Separate methods align by module
position, not by sorting. Method groups `srvMtd` and `srvWOCMtd` distinguish
contextual and no-context methods; `proc`/`func`, `name`, `exp`, `params` and
`param[].pos/dir` describe callable shape. Metadata is not permission to execute
an arbitrary method.

A no-context call POSTs `/e1cib/logForm?cmd=callWOContext` with
`{"root":{"id":id,"fid":fid,"generic":{"name":method,"#param":[types],"param":[values]}}}`.
Omit both parameter members when absent; preserve them together for an explicit
array. Read `ret` and optional `#ret` separately. The command invokes an existing
server method; it is not an endpoint accepting arbitrary BSL source.

## Context calls and synchronization

POST `/e1cib/logForm?cmd=call` with ordered payload:

```text
rKey, id, fid, fop?, fover?, usver?, prm, pchs, prms, dd
```

`prm` is the method object `name,#param?,param?`. `pchs` is an empty object when
there are no changes, or `sinDesc,sin,seq,structure`. A structure contains
optional `fields` (space-separated descriptor **positions**) followed by `v`.
Ordinary scalar proposals include all supported scalar fields; unchanged nested
objects are omitted. Positions, field IDs and names are distinct identifiers.

`prms` contains `descriptions,type,data`. Descriptions retain ordered
`remoteKey,trackChanges,sinDe,seqDe,sinUo,seqUo,structure`. Nested descriptors use
`id,name,nameRu,object,length,fixed,digits,precision,unsigned,date,type,field,__content`
in that order when present. `dd` is dynamic form state, not a descriptor.

The first local change after finalization saves the previous counters and
increments `sin` once; each change fires `seq`. Constructing a proposal copies
counters; it is not another increment. After a response, clear acknowledged
local dirty state, apply the returned root structure, and validate current
`sin == proposed sin`, `seq <= proposed seq` and descriptor `sin == sinDesc`.
Absence of `pchs` is not an invented acknowledgment.

`pchs` carries values, `pdchs` changes descriptors, `elchs.elemcs` is an ordered
control operation stream, and `hippo`/`lflchs` carries layout state. A descriptor
change can accompany a **full** `props` vector. An element can be removed and
its ID reused within one response: applying only the final map loses meaning.
Stage validation before publishing local changes. Unknown structural changes
require an unsupported-profile result, not a recursive dictionary merge.

An application exception can follow a successful mutation. A bounded
`synchException` recovery can reconcile state without repeating the failed
method; [lifecycle and route limits](17-routes-and-errors.md) apply. Recovery
must validate descriptors, identities and counters before accepting values.
Observed omitted presentation defaults are readOnly=false, visible=true,
userVisible=true, enabled=true; this does not authorize dropping other fields.

## FormDataCollection

A snapshot has `count,lastId,fullChanged,firstAddedLine,ids?,column,item`.
`ids` is a space-separated stable-ID list; unloaded entries are placeholders,
not deletions. Snapshot rows directly include `index,id` and scalar members.
The delta is a different shape:

```text
{count,lastId,fullChanged,firstAddedLine,
 operations?:{"##sequence":[kind,arrayIndex,...],kind:[operation,...],...},
 items?:{item:[{index,id,structure:{fields?,"#v"?,v}},...]},
 column:[{id,total},...]}
```

Operations: add `{index,id}`, addRange `{index1,index2,startId}`, remove `{index}`,
removeRange `{index1,index2}`, swap `{index1,index2}`, move `{index,delta}`,
clear `{}`. Ranges are inclusive; removeRange runs from the end. Preserve
interleaving through `##sequence`. Adjacent additions coalesce only when both
indices and IDs are consecutive. Coalescing does not undo counter increments.

Incoming processing order is flags, ordered operations, optional `dropOld`
placeholder reset, resize/count/lastId, changed rows, then decimal totals.
Replayed operations must not enter the next local proposal log. An ID mismatch
at a position replaces the row. Totals use descriptor IDs and exact decimals.
`firstAddedLine` is historical and need not identify a currently existing row.
Lazy fragments, nested object columns and complete context integration are not
covered by the standalone collection schema.

## Closing

POST `/e1cib/logForm?cmd=processFormsClose` with `{"root":{"id":[fid]}}`,
session and next sync headers. Use the returned **fid**, not the request's id.
The bounded client accepts 200/201/204, including an empty response body.
Track attempted versus acknowledged close.
Transport closure alone is not a form-close acknowledgment. Close forms before
logout and release; a failed or unknown method result does not remove cleanup
ownership.

## Client events versus server release

Static clients 18/27/85 describe events as
`event[].id, handlers[].handler[].handlerName/callType`, with callType
Before, After or Override. Keep the nested groups and chain order. Before
handlers can latch Cancel and suppress main execution; After handlers cannot
change the saved cancellation flag. A canceled Before stage does not skip all
After handling. Executing these events requires a separate client runtime.

| Candidate event | UUID | Arguments established by call placement |
| --- | --- | --- |
| OnOpen | `3ccc650e-f631-4cae-8e33-3eaac610b5f9` | Mutable Cancel=false |
| BeforeClose | `52dbb775-1631-4fd5-8c55-1615b5881dac` | Cancel=false, Exit, WarningText, StandardProcessing=true |
| OnClose | `ca21cd18-35b2-4281-b5c8-016ecc8da8ac` | Exit where compatibility permits; no cancellation |

BeforeClose decisions are 0=run standard closure checks, 1=allow closure while
skipping those checks, 2=return a supported nonempty warning, 3=refuse closure.
Exceptions veto closure. Cancel takes precedence over StandardProcessing.
Decision 0 is not unconditional permission: controller checks can still refuse.
OnClose errors are reported while teardown continues. Builds 27/85 handle
Close during OnOpen differently from 18, so this must be profile-specific.

Native browser disposal can queue form IDs, flush when the queue exceeds 20,
or move them to `prm.pfc`/`generic.pfc` on another request. Queue removal occurs
before the operation's result is known. Build 85 has a conditional sendBeacon
path with `beaconAPI=true`; browser acceptance of a beacon is not server release
acknowledgment. These observations do not replace the explicit close exchange.

Normal `synchronize` sends `{rKey,pchs}`. Recovery sends `{key}` and can replace
descriptors, values, elements, layout and dynamic state. `hasdd:true` without
`dd` creates fresh default dynamic state, not an omitted-update no-op.
