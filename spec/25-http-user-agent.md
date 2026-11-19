# HTTP headers and the mobile profile classifier

Date: 2026-11-19. Evidence: E25, S070. This is an offline static study of
Windows x64 8.3.27.2342 components. It follows the adapter gap in
[server SecPswd selection](24-server-secpswd-selection.md). No target was
executed and no authentication request or new live acceptance is involved.

## Which header matters

The inspected request implementation fills the login reader's `userAgent`
from the standard **`User-Agent`** header. The classifier looks for the
case-sensitive substring **`MobileClient`** in that value. It does not use
`1C-ApplicationName` as a substitute or fallback for this field.

Header **names** are matched without regard to case. Header **values** are not
folded to lowercase for the mobile test. These are two different comparisons.

The following are static consequences of this classifier, using a fresh request
with at most one occurrence of each header. They are not authentication tests:

| User-Agent value | 1C-ApplicationName value | Mobile classifier |
| --- | --- | --- |
| `MobileClient` | `1CV8C` | Match |
| `1CV8C` | `MobileClient` | No match |
| `Example/MobileClient` | Absent | Match: substring, not whole-value equality |
| `mobileclient` | `MobileClient` | No match: value comparison is case-sensitive |
| Empty or absent | `MobileClient` | No match: no application-name fallback |

Matching is only the first part of profile selection. For the inspected login
reader, an admitted mobile `version` through `8.3.25.9999` selects mode 1/V2;
newer admitted mobile versions select mode 2/V3. Apply the minimum, maximum and
all three excluded ranges in [chapter 24](24-server-secpswd-selection.md#when-this-implementation-selects-v2).
A nonmobile request instead uses that reader's exact-version check. Neither
changing one header nor reaching the V2 branch proves server admission,
successful password proof, licensing or application readiness.

## Common request adapter

The concrete chain is:

1. Header names pass through `vrscore.dll`'s name-to-enum lookup. Its standard
   header table maps `User-Agent` to enum **39** (`0x27`), using a case-insensitive
   name comparison. The header container also retains named headers separately.
2. The request's `+0x30` getter reads enum 39 from the container at object
   `+0x58`. It widens the stored bytes using the container's UTF-8 converter and
   stores the result in its UTF-16 field at `+0x190`.
3. The getter returns this field. It reads the container when the cached value
   is empty; a missing header leaves the initially empty field unchanged.
   No read of `1C-ApplicationName` occurs in this getter or its enum lookup.
4. QueryInterface for request IID `fc2c0e5c-f5a8-4e25-8c12-05a416d107cb`
   returns the subobject containing this getter. Two inspected concrete
   request tables bind that IID and slot to the same implementation.
5. `vrsbase.dll` queries that request interface, retains it in its request
   owner and passes it into the request context. Context `+0x28` returns this
   pointer. The login reader calls its `+0x30`, copies the result to the
   `userAgent` field and applies the `MobileClient` test from chapter 24.

`1C-ApplicationName` has separate named-header consumers. Its use for application
identity is not evidence that it fills `userAgent`. The earlier native examples
set both headers to the same label, which did not distinguish their roles.

## IIS / ISAPI ingress in the same build

The retained vendor `wsisapi.dll` from 8.3.27.2342 was inspected separately from
the installed `vrscore.dll`. The IIS wrapper is a different interface; its
vtable offsets must not be interpreted as common-request getter offsets.

The wrapper's enumeration method obtains the IIS server variable `ALL_RAW`,
splits header lines at a colon, trims surrounding whitespace and calls a header
callback with byte-string names and values. For ordinary `User-Agent`, the
callback forwards the pair to the common request's header sink. The concrete
sink maps the name to enum 39 and inserts the value into the container used by
the getter above. `1C-ApplicationName` is also passed as its own named header.

This is a concrete interface binding: the IIS request factory uses class
`7f7156dc-0567-42de-b99e-86fc76ec993b`, interface
`b5daf25f-3825-453a-9aa6-41e6f0cde20f`, slot `+0x48`. Its zero fourth ABI
argument selects the common request variant whose table is `0x180131fa0`.
The header callback queries sink interface
`df585fc9-0872-47dd-a840-40920cad90a6`; sink `+0x18` reaches that variant's
name/value insertion. These are internal evidence coordinates, not headers
or parameters an independent client needs to send.

Individual header lookup follows a separate path: it prefixes the requested
name with `HEADER_` and invokes the IIS control block's server-variable
callback. Thus `HEADER_User-Agent`, when requested through this method, is an
IIS lookup key, **not another HTTP header to send**. The demonstrated bulk
transfer to the common request uses `ALL_RAW`; a literal `HTTP_USER_AGENT`
lookup is not required for this route.

The publication-property bridge can supply `1C-ApplicationName: WebClient`
when the incoming application-name lookup fails and the relevant publication
flags permit it. A present-but-empty value does not itself trigger that
missing-header fallback. This branch writes the application-name header; it
does not fill or replace `User-Agent`.

This establishes the static header handoff on the inspected route. It does
not rederive IIS's own parsing, reverse-proxy changes, every optional request
delegate or the adapter selected by a live deployment.

## Independent-client requirements and limits

For the studied mobile profile, explicitly set `User-Agent: MobileClient` and
retain the profile's `1C-ApplicationName: MobileClient` header from chapter 7.
The former is the demonstrated input to this classifier; removing or changing
the latter is not qualified by this study. Keep the identity headers and
advertised version consistent throughout the exchange. Set the header explicitly
instead of depending on the HTTP library's default user agent.

Send a single unambiguous `User-Agent`. Repeated headers, intermediary merging,
post-read mutation of the cached request field, malformed UTF-8 and malformed
HTTP are outside the portable behavior established here. UTF-8 conversion does
not establish how invalid byte sequences are handled.

The header distinction applies to this exact common request implementation and
the adapter routes documented here. It does not establish every web server,
proxy, platform or release. Web `settings` scheme selection remains a separate
flow; this study does not replace it with the native login reader's version
rules or establish live V2 support.

## Evidence coordinates

Addresses use the PE preferred image base and are audit coordinates only.
Independent clients must implement the wire behavior, not depend on addresses,
internal enums, object offsets or interface layouts.

| Module | Coordinate | Role |
| --- | --- | --- |
| vrscore.dll | `0x18002bdf0`, `0x180037c30` | Header-name mapping and case-insensitive table search |
| vrscore.dll | `0x1801874a0` | User-Agent record: enum 39, literal and length |
| vrscore.dll | `0x18002cac0`, `0x18002d7e0` | UTF-8 container construction and enum-value getter |
| vrscore.dll | `0x1800a68f0` | Cached user-agent getter |
| vrscore.dll | `0x18012ef38`, `0x180131fa0` | Request vtables with that getter at +0x30 |
| vrscore.dll | `0x1800a72b0`, `0x1800a66d0`, `0x18002d770` | Header-sink forwarding and name/value insertion |
| vrscore.dll | `0x180131d60`, `0x1800aa2d0` | Factory table and request constructor dispatch |
| wsisapi.dll | `0x180004f20`, `0x18002b410` | IIS wrapper constructor and vtable |
| wsisapi.dll | `0x1800051a0`, `0x1800070f0`, `0x180003850` | Individual HEADER_ lookup and IIS callback |
| wsisapi.dll | `0x180005e30`, `0x18001be70` | ALL_RAW enumeration and header forwarding |
| wsisapi.dll | `0x180013de0`, `0x180012300` | Dispatcher and publication-property bridge |
| wsisapi.dll | `0x180015a49`, `0x180015a53` | Zero fourth ABI argument and factory call |
| vrsbase.dll | `0x180075310`, `0x180075840`, `0x180077b30` | Request interface acquisition, context construction and context getter |
| vrsbase.dll | `0x180093230`, `0x1801925f0` | Login reader and case-sensitive substring matcher |

The [sanitized summary](../evidence/http-user-agent-2026-11-19.json) records exact
module hashes and claim boundaries. Raw vendor code, disassembly, research
tools, private paths and authentication material are excluded. E24/S069 remains
the historical selector-chain summary; E25 resolves its header-source gap only
for the concrete routes established here.

The [Apache comparison and safety review](26-apache-user-agent.md) extends the
adapter study to three available Apache builds and traces the optional build-27
header validator. E25/S070 remains the original IIS/common-request evidence;
the follow-up does not turn these static observations into live acceptance.
