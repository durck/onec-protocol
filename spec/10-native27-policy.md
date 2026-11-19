# Exact-build 8.3.27.2342 record policy

Evidence: E08, E11. Original COM-client n/c query/release was observed with owned
synthetic peers; independent implementation orchestration is synthetic. Real
build-27 native server acceptance is unconfirmed. This specifies the fresh-key,
mode-zero subset, not all native policies.

Layering: `RPC → policy payload → transport escaping → End`. The greeting and
Control3 process identity remain outside the transform.

## Header framing

FH(H) splits a header into 8192-byte chunks, each prefixed by u16LE length.
A short final chunk ends it, including zero for an exact multiple. Bytes after
that chunk are body bytes. This differs from HTTPFile's 16-KiB chunk grammar.

BF(x) is `{#base64:<standard-padded-base64>}`. If Base64 length exceeds 64,
append **CR CR LF** after every full 64-character line, including an exact final
full line. Exactly 64 characters are not wrapped. The implemented reader accepts
CR/LF variations, not arbitrary whitespace or noncanonical trailing bits.

## Initial request

Generate fresh 1024-bit RSA, exponent 65537. N is the unsigned big-endian
128-byte modulus; exponent bytes are 01 00 01.

```text
Hrequest = BOM || '"n"' || ',' || CRLF || BF(N)
               || ',' || CRLF || BF(01 00 01) || ','
requestPayload = FH(Hrequest) || firstRPC
```

Use quoted c instead of n to request c. The final comma denotes the empty
known-key identifier. The header is 215 bytes before its two-byte length prefix.
Nonempty identifiers and #data fields are outside this subset.

Level mapping: n/y/c/u = 0/1/2/4. This implementation selects n or c and rejects
y/u replies. Requested and minimum levels are distinct: request n/minimum n
can accept c; request c/minimum c must reject n. Policy selection is independent
of advertised version.

## Initial response

```text
responsePayload = FH(BOM || BF(rsaCiphertext128)) || encryptedFirstRPCReply
RSA-decrypted bytes = FH(Hmaterial)
Hmaterial = BOM || '"<accepted-level>"'
                 || ',' || CRLF || BF(key24)
                 || ',' || CRLF || BF(iv8)
                 || ',' || CRLF || BF(increment8)
```

Decrypt RSAES-PKCS1-v1_5; require exactly one framed material header with no
remainder. Validate level and 24/8/8-byte field sizes. Decrypt the first RPC
reply with DES-EDE3-CBC using key24/iv8. **The first reply is encrypted even
for n.** Its plaintext is the RPC reply directly, without an ASCII level marker.

DES-EDE3 padding: n=8-(plainLength mod 8), append n-1 random filler bytes and
byte n, always including a final block. The native reader accepts counts 0..8
without checking filler equality; zero is static-only behavior. Normal senders
emit 1..8.

## Subsequent requests and replies

Effective n: payloads are plain RPC. Effective c:

1. Before each request, update
   `u64LE(iv) = (u64LE(iv) + u64LE(increment)) mod 2^64`.
2. Encrypt ASCII 2 followed by RPC, with fresh padding.
3. Decrypt the corresponding response using the **same IV**, no extra increment.
4. For this fixed-c subset, require ASCII marker **2**, remove it, and pass the
   remaining nonempty bytes to RPC. Reject 0/1 as unsupported policy transitions,
   even if they meet the configured minimum; do not continue with stale c state.
5. Keep effective c only after that check. Minimum-level validation remains
   separate from whether a transition is implemented.

No new header accompanies ordinary exchanges. The first c request after
negotiation uses iv+increment, not the initial IV.

### Reply markers can change native state

The [offline audit](22-offline-auth-audit.md) corrected an earlier overbroad
rule that accepted any marker meeting the minimum while keeping the initial
effective level. In exact build 8.3.27.2342, the native reader maps ASCII 0/1/2
to numeric 0/1/2 and **writes that value into the current policy state** before
checking the minimum. The later request builder consults that same state.
Thus a permitted marker change cannot safely be ignored.

The described fixed-c subset deliberately does not implement these transitions.
It rejects markers other than 2 and retires the connection; this is a local
restriction, not a claim that the native reader rejects every 0/1 marker.
The original n/c synthetic-peer observations used stable levels and therefore
do not establish runtime acceptance of level changes. Full transition behavior,
including y, remains outside this subset.

## Synchronization

Allow one pending exchange. Fully validate the initial reply before continuing.
Incomplete I/O, malformed RSA/CBC, unsupported policy or invalid markers terminate
state. Never rewind IV, replay a record or reuse keys on a reconnect. A complete
valid reply advances state even if caller cancellation wins the operation's
result, so cleanup uses current state.

The 1-MiB record/64-KiB header bounds are local policies. AES modes, y selection,
known-key identifiers and alternate providers remain open. Policy encryption
does not itself authenticate the peer; bootstrap/provider/application login
remain separate.
