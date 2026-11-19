# SecPswd V1, V2 and V3 credential messages

Evidence: E10. These are the observed credential state machines, separate from
[native](07-native-authentication.md) and [web](11-web-authentication.md)
HTTP orchestration. V1/V3 have original shipped-code comparison; V2 has
schema/synthetic evidence only. Algorithms describe compatibility behavior,
not recommended password-storage design.

The [independent web experiment](21-independent-web-auth.md) later processed one
real V1 challenge and two real V3 challenges. Final application-start requests
failed without a verified final credential result; this is live progression,
not completed SecPswd login acceptance. V2 remains without a live server test.

The [native profile audit](23-native-profile-selection.md) separately binds
installed client-service C1 methods to V1 on exact builds 19/22 and V3 on build
27. The internal server selector is distinct from those client method targets;
the shared phase byte 1 cannot be used to discover the codec.
The [server-selection follow-up](24-server-secpswd-selection.md) traces the
build-27 older-mobile V2 choice and preservation of its selector across
continuations. This adds static selection evidence, not a new runtime check of
V2 cryptographic primitives or successful server authentication.

## Primitives

`B8(x)=u8(byteLength(x)) || x`; `B32(x)=u32LE(byteLength(x)) || x`.
Text is UTF-8. Complete packets use canonical standard padded Base64 with no
whitespace or alternate alphabet.
Lengths count encoded **bytes**, not characters. Validate B8 <=255 and B32
against the available buffer before slicing. Packet literals 01/03 denote
single binary bytes; they are not ASCII digits. Fully consume fields before
accepting a packet; CRC covers the complete preceding plaintext at its stated
layer, not the outer Base64 text or ciphertext.

`CRC(x)=x || u32LE(CRC32-IEEE(x))`, reflected IEEE CRC-32, initial/final XOR
FFFFFFFF; check value for ASCII 123456789 is CBF43926. Verify CRC where the
message grammar places it and require complete field consumption.

`ENC(K,P,IV)` is AES-256-CBC, 32-byte key, 16-byte IV. Padding is random filler
plus a count: n=16-(len(P) mod 16), append n-1 random bytes and byte n. Aligned
plaintext gains a full block. Decrypt nonempty block-aligned ciphertext; final
count must be 1..16. Remove n bytes without checking filler equality. This is
not PKCS#7. Each ENC starts afresh with the stated IV.

Fixed protocol constants, not user credentials:

```text
messageIV  = 9d7b9a206965bb28067a483db26c718e
envelopeIV = 126bafea20772e4211cc2197f85ef32f
envelopeK  = a3aeb390acbc1580aebf118dd5629be42e0e07adb134356a831053c5718109cd
```

Use cryptographically secure fresh random values in real exchanges.

## Password-derived key

If caseSensitive=false, apply Unicode default **full** uppercasing to password;
otherwise preserve it. No extra Unicode normalization is introduced. ASCII-only
or locale-specific casing is not equivalent for all input. Let D=UTF8(password).
For modes 1..3 append ASCII(Base64(salt)); mode 0 ignores salt.

The Unicode data version used by every shipped build has not been established.
An implementation must record the version of its full-case mapping dependency
when claiming non-ASCII compatibility; choosing a library version is not wire
negotiation. All published password vectors currently use ASCII. They establish
none of the version-sensitive non-ASCII cases, and the independent client
experiment's Unicode 15.0.0 choice is not a platform-wide requirement.

| Mode | Digest |
| ---: | --- |
| 0 | SHA-1(D) |
| 1 | SHA-256(D) |
| 2 | SHA-512(D) |
| 3 | PBKDF2-HMAC-SHA-256(password=D, salt=empty, iterations=16000, output=64 bytes) |

`PasswordHash=Base64(digest)`;
`PasswordKey=SHA-256(ASCII(PasswordHash))` (32 binary bytes).
The mode-3 challenge salt is already in D as Base64 text; it is **not** the
PBKDF2 salt parameter. Salt fits B8, at most 255 bytes.

## V1: C1 → S2 → C3 → S4

Generate 32 random bytes Rcc. Derive mode-0 case-sensitive K0 and uppercase K1;
deduplicate equal keys.

```text
C1 = Base64(CRC(01 || B32(UTF8(user)) || B8(Rcc)))
S2 = Base64(ENC(Kselected,
                 CRC(B8(noise) || B8(Rs) || B8(Rcc)), messageIV))
C3 = Base64(ENC(Rs,
                 CRC(B32(UTF8(originalPassword)) || B8(noise)), messageIV))
S4 = Base64(ENC(Rs, CRC(B8(result) || B8(finalNoise)), messageIV))
```

Try each distinct candidate key for S2. Require valid padding, CRC, exact field
layout, 32-byte Rs and exact Rcc echo. Exactly one candidate must validate.
Preserve returned noise for C3; its encoded length is not fixed here. Success
requires result exactly ASCII suc; finalNoise is parsed then ignored.

## V2: C1 → S2 → C3 → S4

Generate 32 random bytes Rc. Derive both mode-0 keys as above and encrypt Rc
under each independently with fresh padding.
V2 always emits **both** encrypted alternatives, even if K0=K1. Only V1's
candidate-decryption loop deduplicates equal keys. Do not apply that V1
optimization to V2's C1 layout.

```text
C1 = Base64(CRC(01 || B8(ENC(K0,Rc,messageIV))
                  || B8(ENC(K1,Rc,messageIV)) || B32(UTF8(user))))
S2 = Base64(CRC(B8(ENC(Rc,Rs,messageIV)) || strongASCII))
C3 = Base64(ENC(Rs,
                 CRC(B8(freshNoise32) || B32(UTF8(originalPassword))), messageIV))
S4 = Base64(ASCII("suc"))
```

strongASCII must be byte 30 or 31 hex, characters 0/1. Rs must decrypt to
exactly 32 bytes. The client validates the flag but sends the **original**
password in C3. Final decoded suc has no CRC or encryption; another decoded
marker is rejection. The studied web flow does not select V2.

## V3: C1 → S2 → C3 → S4 → C5 → S6

```text
C1 = Base64(ENC(envelopeK,
                 CRC(01 || B8(freshNoise16) || B32(UTF8(user))), envelopeIV))
S2 = Base64(ENC(envelopeK,
                 CRC(B8(serverNoise) || hashMode || B8(salt) || sensitiveByte),
                 envelopeIV))
C3 = Base64(CRC(03 || B8(ENC(PasswordKey,Rc,messageIV))))
S4 = Base64(CRC(B8(ENC(Rc,Rs,messageIV)) || strongASCII))
C5 = Base64(ENC(Rs,
                 CRC(B8(freshNoise32) || B32(UTF8(originalPassword))), messageIV))
S6 = Base64(ASCII("suc"))
```

In S2, hashMode is binary 0..3; sensitiveByte is binary 0/1, **not ASCII**.
Parse serverNoise, derive PasswordKey using the received parameters, then
generate 32 fresh bytes Rc. In S4 strongASCII is ASCII 0/1 and Rs is 32 bytes.
V3 S4/C5 share V2 S2/C3's construction. S6 is raw suc without CRC/encryption.
The S4 ASCII flag is separately validated; the retained flow does not require
it to equal S2's binary sensitivity flag or rederive a different key from it.

## State and validation

| Version | Server messages | Final application proof |
| --- | --- | --- |
| V1 | S2, S4 | `response.credAns4` |
| V2 | S2, S4 | `response.credAns4` |
| V3 | S2, S4, S6 | `response.credAns6` |

Reject out-of-state messages; never restart with old nonces. Clear owned secret
material on completion/failure. The reference's 18-second deadline and packet
caps are local policies. Codec text inputs are valid UTF-8 without NUL.

Validate lengths, canonical Base64, CRC, padding count, modes, flags, V1 echo
and key uniqueness, final proof, HTTP placement and routing continuity.
Cryptographic acceptance, application-session creation and cleanup are separate.

V2/V3's final suc is a completion marker, not an independent MAC or proof of
server identity. Its enclosing exchange and provider validation still matter.
The final proof field name does not change those cryptographic properties.
The [audit and transcripts](19-bootstrap-auth-audit.md) exercise missing, early,
repeated and negative final messages independently from HTTP 200/session receipt.
