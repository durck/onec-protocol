# Text and binary serialization

Profile: `native-18`. Evidence: E01, E03 and E04 in the
[register](../evidence/README.md). This chapter specifies the selected subset;
it does not enumerate every native ListStream type.

## Text records

Selected bootstrap and interface-query records use UTF-8 with the three-byte
BOM `EF BB BF`. Lists use `{` and `}`, and values are comma separated. Strings
are quoted with `"`; an embedded quote is doubled. Unquoted integers and UUIDs
are distinct from quoted strings. Backslash escaping is not the string rule.

```text
<BOM>{0,12345678-1234-5678-9abc-123456789abc,7,1,bc15bd01-10bf-413c-a856-ddc907fcd123}
<BOM>{"S","a""b"}
```

`<BOM>` is explanatory notation, not literal text. The top level can be a
sequence of values, for example a list followed by a scalar. The local reader
accepts spaces, CR, LF and tab between tokens. Its depth/size restrictions are
implementation limits.

## Binary tags

Each atom starts with one tag byte: `code = tag & 0x1F`,
`flags = tag & 0xE0`. All codes in this table are hexadecimal.

| Code | Payload | Interpretation in the selected schema |
| --- | --- | --- |
| `00` | None | Structural flags only; bare zero is not a value |
| `01..0A` | None | Integer `code - 1`, in the range 0..9 |
| `0B` | One byte | Octet; scalar integer or one-byte data according to context |
| `0D` | Two bytes LE | Unsigned integer |
| `0F` | Four bytes LE | Unsigned integer |
| `15` | 16 bytes | GUID in the order defined below |
| `1A` | u8 byte length, then data | Raw bytes |
| `1B` | u16LE byte length, then data | Raw bytes |

An empty byte sequence uses `01`, the same compact form as integer zero. It is
the enclosing schema that distinguishes them. Raw bytes are not intrinsically
UTF-8. This draft does not specify negative integers, floating point, Numeric,
wide strings or long raw atoms; their absence here does not imply native absence.

For the synthetic canonical examples, integers 0..9 use the compact form,
10..255 use `0B`, 256..65535 use `0D`, and larger uint32 values use `0F`.
Raw data uses `01` when empty, `0B` for one byte, `1A` for 2..255 bytes and `1B`
for 256..65535 bytes. These sender choices do not rule out other native encodings.

## Structural flags

| Mask | Logical action |
| --- | --- |
| `20` | End the current list |
| `80` | Separator at a position requiring one |
| `40` | Begin a list |

The native token consumer exposes combined flags in the order **end, separator,
begin, atom**. Thus a combined tag cannot be interpreted as one indivisible
type code. The reference tree builder does not validate every separator; its
ability to reconstruct a tree is weaker than native grammar acceptance.

An ordinary list encoder emits `40`, its first value, later values with `80`,
and `20`. A later nested list can begin with `C0 = 80 | 40`. For example:

```text
[0, 4, 1007] = 40 01 85 8D EF 03 20
```

Typed representations include `["B", 0 or 1]` and `["#", typeGUID, ...]`.
The `0B` atom used to encode the byte `B` is not itself a Boolean type code.

**Exception:** the HTTPFile chunk stream and the SecurityBlob inner pair use
specialized separator rules. See [HTTP RPC](04-http-rpc.md). Mechanically
serializing their logical arrays as ordinary lists produces different bytes.

## GUID byte order

Reverse the first four-byte group, the following two-byte group and the next
two-byte group. Keep the last eight bytes in their displayed order:

```text
12345678-1234-5678-9abc-123456789abc
78 56 34 12 34 12 78 56 9A BC 12 34 56 78 9A BC
```

| Role | GUID |
| --- | --- |
| HTTP interface | `bc15bd01-10bf-413c-a856-ddc907fcd123` |
| HTTPFile type | `bea56fd3-3695-49fc-bbc5-c8bc15f33d41` |
| SecurityBlob type | `671507fd-50a9-4b63-b70e-58b3d364f48f` |

The `15` atom code is separate from the 16 GUID bytes. A separator-prefixed
GUID therefore starts with `95`, not with a byte from the UUID itself.
