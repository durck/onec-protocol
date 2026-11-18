# Dynamic-list windows

Status: bounded build-18 live rows and source-derived request construction.
Not a general query-language endpoint, tree engine or arbitrary paging API.

## Ordered query

POST `/e1cib/dlist?cmd=query` using the application JSON/session/sync envelope.
Preserve this member order; alphabetically sorted requests failed native XDTO
conversion before the form binding was consumed:

```text
remoteKey, dataPath:{id:[pathComponents]}, tableId, cacheID,
list:{ChangeState:0,SettingsComposer:{}}, choiceMode:false, formUUID,
refresh:true, forward:true, autopos, page, "#current"?, current?,
source:0, showTree:false, "#filter":[], filter:[], level:-1,
period:{startDate,endDate,"#variant":type,variant:"Custom"},
presents:true, search?, allowAsync:false, waitTime:200, needDataDescr:true
```

The form's `remotekey`, `fid`, bound field path, table ID and period descriptor
provide these identities. They cannot be replaced with a table name. `cacheID`
is `00000000-0000-0000-0000-000000000000` in this profile. First page uses autopos=true. A
continuation uses autopos=false and the returned row key with its actual type
in `#current,current`; do not use a visible-column index as a bookmark.
`page` is requested window size. A typed search entry is
`{path:fieldName,"#value":type,value:payload}`.

## Response envelope

The bounded flat profile has level=-1, no hierarchy or normalizeKeys, one
`forward` chunk and an absent/empty backward list (or a single
`{"intentionallyEmpty":true}` entry). A forward intentionallyEmpty object
contains no row data. Otherwise `fif` and `lil` mark first/last boundaries,
`dataName` is `DynamicListRow`, and `key` names the key field.

An `inProgress:true` response carries `backgroundSearchInfo` and no completed
rows. It is a distinct state, not an empty successful page; general asynchronous
continuation remains outside this bounded client contract.

Rows can be JSON `dataDescr,rows` or compact `dataDescrData,rowsData` base64
fields. Do not accept a mixed representation. The compact fields in the
observed profile are binary ListStream, not compressed XML. An absent rows
member can mean an empty row sequence. Match descriptor fields by identity
and name, then project visible columns; hidden key fields still occupy cells.

## Compact grammar

The observed descriptor structure is:

```text
[0,1,[0,[0,fieldCount,
 [[fieldID],name,name,[0,0,["Pattern",pattern]],marker,0,0],...
 ],"DynamicListRow"],0,0,0,0]
```

The key field has marker 2; observed visible fields have 0. Covered patterns
are `["N"]` and `["S",6,1]`; do not generalize the latter two numbers to
unobserved qualifiers. Names may use UTF-16LE string tags. The two name copies
must agree, IDs must be unique, and exactly one name must match envelope `key`.

Rows are `[count,[genericCell,...],...]`. A synthetic three-column row could
be `[["N",1],["N",1],["S","ROW_1"]]`: the leading key is additional to two
visible values. Validate full byte consumption, list boundaries, count, row
width and type correspondence. A tag such as `4B 4E` begins a list with a
one-byte string `N`; `4E` there is payload, not a negative integer tag.

Unknown headers, patterns, hierarchical rows, reference keys, changed descriptors,
appearance payloads and cross-session bookmarks require further evidence.
Preserve the difference between successful decode of a captured payload and
live acceptance of a newly constructed query.
