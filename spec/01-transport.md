# Native transport framing

Profile: `native-18`. Evidence: E01 and E02 in the
[register](../evidence/README.md). The encoder/decoder details below are backed
by source and synthetic fixtures; the baseline native connection path has
bounded live acceptance. Not every malformed-stream behavior was tested natively.

## Greeting and control items

The server greeting is five bytes: `53 F5 C6 1A 7B`. It precedes the framed
exchange and is consumed separately from the record decoder.

| Label | Bytes | Documented interpretation |
| --- | --- | --- |
| Escape | `65 52 B1 A5` | Quotes one following special token |
| End | `66 53 B2 A6` | Ends one record |
| Control2 | `21 47 54 B3` | Separate control event with no payload in this subset |
| Control3 | `22 48 55 B5` | Separate control event followed by 16 raw bytes |

During initialization the client sends `Control3 || processGUID`, where the
GUID uses the [binary GUID order](02-serialization.md#guid-byte-order).
This control item is not an HTTP request and has no appended record terminator.
The full semantics of both control event types remain unspecified.

## Record encoder

Scan the serialized payload left to right. If the next four bytes equal any of
the four tokens above, append `Escape` followed by those four bytes and advance
by four. Otherwise copy one byte and advance by one. Append one unescaped `End`.

This operation covers the entire payload, including scalar values, GUIDs,
length fields and embedded HTTP. It is not restricted to HTTP bodies.

```text
Payload: 41 66 53 B2 A6 42
Wire:    41 65 52 B1 A5 66 53 B2 A6 42 66 53 B2 A6

Payload: 65 52 B1 A5
Wire:    65 52 B1 A5 65 52 B1 A5 66 53 B2 A6
```

For this encoding, `wire_length = payload_length + 4 * quoted_tokens + 4`.
This framing adds no checksum. Other layers may have their own protection.

## Incremental decoding

Preserve partial token prefixes across TCP reads. After `Escape`, wait if the
available following bytes could still become a special token. A complete special
token is appended literally to the record and does not create a control event.
The reference decoder preserves an Escape followed by a non-token as literal
bytes; this is an implementation behavior, not a native malformed-input guarantee.

An unescaped End emits the accumulated record. Control3 consumes its next
16 bytes without applying record-token parsing to them. Control events do not
implicitly finish an accumulated record in the reference decoder.

One TCP read may contain zero, one or many completed events. A record, token,
quoted pair or control payload may span multiple reads. TCP segmentation and
application writes are not message boundaries.

The reference decoder treats EOF with an unfinished record, partial token or
partial control payload as incomplete input. Buffer and record limits must be
explicit local choices. The existing synchronous client expects one response
per exchange; this does not establish a platform-wide ban on multiplexing.
