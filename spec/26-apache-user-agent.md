# Apache adapters and User-Agent safety boundaries

Date: 2026-11-19. Evidence: E26, S071. This offline study covers the Windows x64
`wsap22.dll`, `wsap24.dll` and `wsapch2.dll` artifacts available for
8.3.18.1208, 8.3.27.2342 and 8.5.1.1529. It compares their header ingress with
the [exact build-27 IIS route](25-http-user-agent.md). No web server or vendor
target was executed, and no authentication or malformed request was sent.

## Header acquisition and forwarding

All nine inspected Apache adapters obtain incoming headers from the Apache/APR
request header table. They do not parse an IIS `ALL_RAW` string:

1. Individual lookup calls imported `apr_table_get`. A null pointer means the
   lookup failed; a nonnull pointer to an empty string is a found empty value.
2. Enumeration calls imported `apr_table_do` with a local callback and no
   selected-key list. The bridge forwards to the supplied callback and returns
   1 to continue enumeration.
3. The next callback skips a null or empty name **or value**. For remaining
   entries it scans the host-provided C strings to their first NUL, constructs
   counted byte strings and invokes the header filter/forwarder.
4. Ordinary `User-Agent` passes through that forwarder as a name/value pair.
   It is not lowercased, trimmed or parsed as commands by this adapter path.
   Host parsing and any earlier normalization remain separate.
5. `1C-ApplicationName` is forwarded separately. It is an allowed application
   header, not a replacement source for `User-Agent`.

The request field offsets differ: the inspected ap24 wrappers read the header
table at request offset `+0xe8`, while ap22/apch2 read `+0xf0`. These are
build-specific ABI observations, not fields of the wire protocol. Adapter and
host compatibility cannot be inferred from their similar names or code.

In build 27, the Apache dispatcher requests the same factory class/interface
and header-sink interface established for IIS in chapter 25. Its zero fourth
Windows x64 ABI argument selects the same common request variant. The sink
therefore reaches the same enum-39 container and `User-Agent` getter used by
the login reader. The case-sensitive `MobileClient` classifier and the
[version gates for V2](24-server-secpswd-selection.md) apply to that recovered
build-27 path. They are not automatically asserted for builds 18 or 85.

The exact build-85 common request was also traced independently: its getter
reads enum 39/User-Agent, with UTF-8 widening and no application-name fallback.
That establishes the header source, not build-85 authentication version gates.
For build 18 this pass stops at the common sink interface and does not substitute
the newer common implementation as evidence.

## Comparison with IIS

These comparisons describe the inspected adapter layers, not what arbitrary
hosts or proxies accept on the network:

| Property | Apache artifacts studied here | IIS 8.3.27.2342 |
| --- | --- | --- |
| Bulk header source | APR request header table | IIS server variable `ALL_RAW` |
| Individual lookup | `apr_table_get` | `HEADER_` key through the IIS callback |
| Local delimiter parsing | No reparse of raw header lines | Split lines at newline/colon |
| Local surrounding-whitespace handling | No additional trim in the ordinary forwarding path | Trim name and value during `ALL_RAW` parsing |
| Empty value during enumeration | Skipped before counted-string forwarding | Empty trimmed value can reach the header callback |
| Name/value representation at sink | Counted byte strings made from host C strings | Counted byte strings made from parsed spans |
| User-Agent-specific rewrite | None in the recovered ordinary branch | None in the recovered ordinary branch |

An individual lookup and enumeration need not expose empty headers identically.
Do not infer application-name defaulting from what remains in the common
container after enumeration. Likewise, adapter forwarding is not proof of
first-wins, last-wins or comma-joining semantics for repeated wire headers;
upstream table construction and intermediaries are outside this comparison.

In all three studied builds, the conditional `WebClient` application-name
default uses the individual lookup result. A present empty application-name
header therefore suppresses that missing-header default, although the empty
entry is skipped during Apache enumeration. It does not populate User-Agent.

## Bounded security review

**No exploitable vulnerability was established in the inspected User-Agent
path. This is not a guarantee that the adapters or platform are vulnerability
free.** The review traced acquisition, conversion, insertion, classification
and selected diagnostic output. It did not execute the target, fuzz a web
server or analyze every downstream consumer.

### Control characters and header injection

The common build-27 request constructor attempts to acquire a header-validation
service implemented by `inet.dll`. When that object is present, insertion
checks the name and value before storing either. A failed check throws a
exception carrying status 400. The inspected byte-value predicate rejects C0 control bytes
other than horizontal TAB, including raw CR, LF and NUL; it accepts an empty
counted string. Its accepted byte range also includes DEL and high bytes;
this predicate alone is not Unicode validation or log/HTML escaping.
This establishes an additional check downstream of the Apache
callback, not a universal promise about every process configuration.

The service lookup can return null, and the insertion method skips this check
when the validator is absent. This study does not prove that every deployment
instantiates the validator. Apache also supplies C strings: an embedded NUL in
a hypothetical host buffer would terminate its scan before insertion, rather
than arrive as a counted NUL for that validator. Neither observation establishes
that a corresponding malformed HTTP request reaches the adapter.

Build 85 also has an optional policy acquisition/check branch, but its exact
validator byte rules were not recovered in this study. The build-27 predicate
must not be substituted for build-18 or build-85 behavior.

The predicate was checked with a private synthetic model against the original
table bytes: all 256 single-byte values and eight sequence cases. This model
is not execution of vendor code. Malformed UTF-8 conversion, Unicode control
characters, obsolete line folding, proxy disagreement and request smuggling
remain unqualified.

### Execution and formatting sinks

The recovered ordinary header path copies data, compares header names and
performs the mobile substring test. No command, SQL, template or HTML execution
sink was established on that path. This is a scoped negative finding, not a
search result for all platform uses of User-Agent.

The inspected login-state formatter appends `userAgent:` followed by the value
using a counted string append. The value is not used as a printf format string.
The result is passed as data in a logger property list. At this intermediate
layer, commas and colons within the value are not escaped from the surrounding
diagnostic text. Consequently, that text must not be treated as an unambiguous
trusted map of login-state fields. Final log serialization, log consumers and
any exploitable log-forgery effect were not established.

The inspected invalid-header error path likewise takes its format template
from a resource and supplies the header strings as arguments. It does not use
the incoming value as the format template.

### Memory, limits and profile selection

The inspected forwarding paths construct dynamic counted strings, not a fixed
User-Agent stack buffer. Their C-string scans depend on the Apache/APR host's
valid-pointer and NUL-termination contract. This review did not establish a
network-triggerable out-of-bounds access, overflow or crash. It also did not
establish one universal header-length or allocation limit; host limits and
resource-exhaustion behavior require separate evidence.

User-Agent is client-controlled classification data. Matching `MobileClient`
can select the older-mobile branch only with the accompanying admitted version;
it does not authenticate the caller or bypass the subsequent password proof.
The recovered profile choice is not itself evidence of an authentication bypass
or a successful downgrade attack.

## Independent-client guidance

Send one explicit User-Agent and keep it stable with the selected application
identity and version. Reject CR, LF and NUL when constructing header values;
do not rely on server-side validation being present. Avoid repeated or ambiguous
identity headers. Treat received diagnostics as untrusted text and apply
escaping appropriate to the log or UI destination when displaying them.

The [sanitized evidence summary](../evidence/apache-user-agent-2026-11-19.json)
contains exact module hashes, coordinates, distinctions and remaining gaps.
Debug tools/tests, synthetic probe programs, original vendor code, internal
paths and authentication material are excluded from the publication.

Selected reproducible coordinates below are RVAs in the exact **wsap24.dll**
images; they are not runtime addresses or offsets shared with every adapter:

| Build | Lookup / enumeration | C-string callback / counted forwarder |
| --- | --- | --- |
| 8.3.18.1208 | `0x5710` / `0x5850` | `0x192d0` / `0x18d20` |
| 8.3.27.2342 | `0x24d0` / `0x2560` | `0x17e00` / `0x17ed0` |
| 8.5.1.1529 | `0x2570` / `0x2600` | `0x18040` / `0x18110` |
