# HTTP carried by native RPC

Profile: `native-18`. Evidence: E01, E04, E06 and E07 in the
[register](../evidence/README.md). The multi-chunk correction is supported by
static code, an observed native reply and a bounded live request.

## HTTP bytes

Construct an HTTP/1.1 message with CRLF line endings, a blank line and body bytes.
`Content-Length` counts only body bytes. Reassemble all HTTPFile chunks before
interpreting headers, text encodings or compressed content. Chunk boundaries
can split a header, CRLF pair or UTF-8 sequence.

The examples use one Content-Length and no Transfer-Encoding. The reference
parser rejects Transfer-Encoding; that restriction is not a claim about every
platform HTTP implementation. The native file chunk format below is separate
from HTTP chunked transfer coding.

## Request envelope

The first HTTP call uses the following logical binary list:

```text
[
  0, target, counter, 4, HTTPInterface, 0, 1, 3, 0,
  HTTPFile(httpBytes), 2, SecurityBlob(token), 2, ["B", securityFlag]
]
```

The continuation form inserts an extra zero **after counter**, before `4`.
The reference connection selects that form for subsequent HTTP exchanges on
the same context. This envelope continuation is unrelated to HTTPFile chunks.
`securityFlag` is 0 or 1. Fixed fields retain their observed values; a complete
official IDL signature has not been recovered.

`target` identifies the remote RPC object. `counter` is the RPC exchange counter,
not the HTTP Content-Length or the application synchronization number.

## HTTPFile

```text
["#", bea56fd3-3695-49fc-bbc5-c8bc15f33d41, 1, 0, 0, 0, chunks]
```

The sender uses `1,0,0,0`; the full native meanings of these fields remain open.
The reader currently accepts a nonzero first numeric field and zero in the
following three positions. That is an accepted local schema, not a proof of
all valid native variants. The object does not imply creating a disk file.

Split the **entire HTTP message** into chunks of 16,384 bytes. The terminal
chunk must be shorter. An exact multiple therefore requires an extra empty
chunk.

| Total HTTP bytes | Chunk lengths |
| ---: | --- |
| 97 | 97 |
| 16383 | 16383 |
| 16384 | 16384, 0 |
| 16385 | 16384, 1 |
| 32768 | 16384, 16384, 0 |

The chunks are raw byte atoms inside a specialized list. The list is preceded
by the parent separator (`C0` when the begin marker is separate). **No chunk
continuation has an `80` separator.** The full-chunk atom is `1B 00 40`, followed
by 16,384 data bytes. A zero terminal chunk is `01`; a one-byte terminal chunk
is `0B xx`. Close the chunk list with `20`.

```text
HTTP length 16384:
C0 1B 00 40 <16384 bytes> 01 20

HTTP length 16385:
C0 1B 00 40 <16384 bytes> 0B <one byte> 20
```

The observed native reply fused the first list begin and atom into
`DB = 80 | 40 | 1B`. Separate `C0 1B` and combined `DB` expose the same logical
tokens in this position. The correction concerns continuation separators.

An observed 23,924-byte HTTP reply had chunks 16,384 and 7,540; its second atom
was `1B 74 1D`, without `80`. The empty terminal chunk rule has static and
synthetic evidence; that particular capture did not exercise an empty tail.

Older examples with `9B` continuations or an `81` empty tail describe the
superseded generic-list encoder and must not be used for native HTTPFile.

## SecurityBlob

```text
["#", 671507fd-50a9-4b63-b70e-58b3d364f48f, [tokenLength, tokenBytes]]
```

The inner pair is specialized: emit its begin marker with the parent separator,
then length and bytes **without a separator between them**. Close both lists.
For an empty blob in a later parent-list position:

```text
C0 0B 23 95 FD 07 15 67 A9 50 63 4B B7 0E 58 B3 D3 64 F4 8F
C0 01 01 20 20
```

The length must match the token byte count. An empty token still has its typed
object; omitting the field changes the envelope. Synthetic vectors use only
an empty blob and contain no authentication token.

## Reply envelope

The reference success reader accepts a top-level sequence consisting of a list
and trailing integer `2`:

```text
[1, field1, 1, fileOrUndefined, 3, 0, 2, SecurityBlob(token),
 2, ["B", securityComplete], 0], 2
```

`field1` is deliberately unnamed: this reader does not validate its semantics.
`fileOrUndefined` is HTTPFile or `["U"]`. The Boolean flag must be 0 or 1.
Native exception records are a separate branch, not ordinary HTTP error PDUs.
The complete exception grammar and correlation rules remain outside this draft.

The source also accepts a bounded empty response with a final CRLF instead of
the normal blank-line separator, only when Content-Length is zero. This
compatibility branch is a local accepted schema, not a requirement for senders.

## Complete synthetic example

[vectors.json](../vectors/vectors.json) includes a complete HTTP request, binary
RPC and framed record using the target
`12345678-1234-5678-9abc-123456789abc`, counter 1007, an empty security blob and
security disabled. The body `{"cred":"demo"}` is explanatory text, not usable
authentication material. The example checks packaging only and is never sent
to a server. Read its recorded byte lengths rather than counting displayed
characters or TCP segments.
