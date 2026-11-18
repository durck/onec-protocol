# Files, temporary storage and processing attachments

Status: selected build-18 live operations; chunk construction and ownership
rules traced in shipped clients 18, 27 and 85. These mechanisms are separate
from native HTTPFile's 16-KiB transport chunks.

## Open-document processing upload

POST `/e1cib/openDocument?seanceId=<escaped full session token>` as multipart
form data, with fields in order `storageId` (fresh UUID), `docType` (`processing`),
`file` (filename, original bytes, application/octet-stream). This path carries
the session in the query, without an additional session header or sync pragma.
Accept 200/201 only after parsing the result.

The returned text can begin with a UTF-8 BOM. An encoded result starts
`processing%08`; decode with JavaScript decodeURI semantics (reserved URI
characters `;/?:@&=+$,#` stay escaped). A raw result starts `processing` followed
by byte 08; split it directly and preserve literal percent escapes. Mixed or
invalid representations are errors. Split on backspace into two to four fields;
the second is the opaque processing reference. It is not a filesystem path.

## Metadata attachment alternative

An explicitly selected alternative is POST `/e1cib/extmd/processing` with one
multipart `file` part and the session header, without sync. A 201 reply must
contain a single Location of form
`/e1cib/extmd/processing/<kind.name>/<40 hex characters>`. Supported kinds are
`ExternalDataProcessor` and its localized equivalent. The location is an
identity result, not a redirect to follow. The processing reference is kind.name.
Reject absolute origins, query or fragment components in this profile.

Do not switch between these two operations after failure: a lost response may
follow successful server acceptance. No claim of disk-free server execution
follows from uploading bytes rather than a server path.

## Interactive security challenge

The observed challenge is HTTP 400, application/json, exception class
`330cdbe8-ffef-47e5-b574-38a1ba2368a6`, numeric operationType=2,
Boolean commandCouldNotRestart, and moduleSignature holding canonical Base64
for 20 bytes. Other exceptions are not this challenge.

The explicit confirmation operation POSTs `/e1cib/opersec/insert` with ordered,
**unwrapped** `{"sign":opaqueSignature,"oper":2}`, session header and no sync.
Accepted HTTP codes are 200/201/204; a nonempty JSON body must not contain an
exception. Approval can persist beyond the connection. Bind approval to the
exact file bytes and originating session; do not automatically confirm an
arbitrary file. Confirmation does not open a form or perform the upload again.
Its lost reply is an unknown outcome, not permission to replay the signature.

## Chunk upload

The inspected browser clients slice a file into 524288-byte parts and send
them **sequentially** to `/e1cib/upload`. Query fields are `part` (zero-based),
`isLast`, `storageId`, `seanceParamPrefix`, and optional replacement `url`.
The multipart `file` part contains only that slice; send the session header.
`storageId` is the supplied/default owner UUID; `seanceParamPrefix` is a fresh
per-file UUID, independent of filename and returned storage address.

There are `max(1,ceil(size/524288))` parts. Empty input sends one empty final
part. Exact multiples have no extra empty tail. Only the final response text
becomes the storage address. Validate status and final address even though the
inspected browser chunk callback does not perform its own status check.

Cancellation prevents later parts but does not establish deletion of already
accepted bytes. No resume, duplicate-part semantics, checksum field or explicit
partial-transfer rollback was established. A lost final response is an unknown
outcome. Native transport limits must accommodate multipart and HTTP overhead
in addition to the 524288-byte payload; a 512-KiB complete-PDU cap is too small.

## Storage and download boundaries

Typed-value temporary storage, uploaded file storage and open-document references
have different APIs. A storage address is not an OS path. Inspected typed-value
operations select ownership from a supplied UUID, an existing recognized
storage address, or the default client owner. Unowned documented temporary
values can be removed on a subsequent call; form/session association changes
lifetime. This does not specify interrupted-upload assembly lifetime.

Known download construction is `/e1cib/files?url=<address>&targetFileName=<name>`;
archive mode first requests `/e1cib/archive`. Form cleanup can issue
`DELETE files/<storageID>?cmd=RemoveFromTempStorage` relative to the application
service. A controlled small build-18 download was verified. General large-file
streaming, ranges, resumability, secure erasure, cross-worker continuation and
replacement atomicity remain unknown.

Client upload-size ceilings are resource policy, not server limits. Network
break experiments showed that an upload can mutate server state before its
reply is lost and that a session can survive a broken socket. Always distinguish
attempted, acknowledged and unknown cleanup.
