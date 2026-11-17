# Bootstrap and remote context

Profile: `native-18`. Evidence: E01, E05 and E08 in the
[register](../evidence/README.md). Request shapes are source-derived; the native
baseline path has bounded live confirmation. This chapter is not yet a complete
bootstrap implementation recipe.

## Phase order

1. Optionally connect to the cluster manager and obtain a worker selection.
2. Open the worker TCP stream and consume its five-byte greeting.
3. Send Control3 and a client process GUID.
4. Exchange bootstrap records, introduction data and provider tokens.
5. Obtain the dynamic remote context UUID, called `target` here.
6. Query the HTTP interface on that target.
7. Start HTTP RPC exchanges and application authentication.

The manager and worker can use configured ports. Port numbers are not fields
of HTTPFile and must not be hardcoded as protocol constants. Direct access to
a known worker omits the client's manager selection step.

## Introduction payloads

The reference request payload is:

```text
ASCII("V8IntroTicketReq") || u16LE(len(name_utf8)) || name_utf8 || 00
```

The request-side local validation requires nonempty, valid UTF-8 without an
embedded NUL and a byte length at most 65535. The terminating NUL is not counted.

The acknowledgment payload is parsed as:

```text
ASCII("V8IntroTicketAck") || u16BE(len(first_utf8)) || first_utf8
                         || u16BE(len(second_utf8)) || second_utf8
```

The two acknowledgment strings are nonempty in the accepted schema and have
no trailing NUL or trailing bytes. The **length byte order differs** between
request and acknowledgment. These are inner payloads, not standalone TCP
messages. Their complete placement in every provider exchange remains outside
this draft.

## Interface query

After bootstrap returns `target`, the request is a text list:

```text
<BOM>{0,target,counter,1,interfaceGUID}
```

For the HTTP interface, `interfaceGUID` is
`bc15bd01-10bf-413c-a856-ddc907fcd123`. Apply record framing after text encoding.
The current reply reader accepts one three-element list beginning and ending
with integer `1`; it does not establish semantics for every intermediate value.

## Release

The context release request is a binary list:

```text
[0, target, counter, 0, 0]
```

The current accepted release reply is one two-element list beginning with
integer `1`. Complete reply correlation, failure variants and intermediate
field semantics require further specification.

Release refers to the RPC context. Application logout, local socket/provider
cleanup and eventual server-session disposal are distinct outcomes. See
[session lifecycle](05-session-lifecycle.md).

Build **8.3.27.2342** has an additional initial policy-envelope boundary.
Do not apply this phase description to that build as a complete byte recipe;
see the [version matrix](06-version-matrix.md).
