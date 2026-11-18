# Administration collection schemas

These positional schemas accompany [administration](12-administration.md).
They describe observed decoders and exact profile families. Unknown tails must
be retained or rejected explicitly, never silently interpreted as another build.
Counts are integers and must match the number of entries; limits are client
resource policy, not inferred server maxima.

## Common representation

Each collection is `["#",typeUUID,[count,row...]]`.

| Collection | Type UUID |
| --- | --- |
| Clusters | `41344e9c-168a-4c10-8df4-0a3700dbd815` |
| Administrators | `9d95bd82-b4c3-49d5-9077-bc703bd7dc09` |
| Working servers | `85ae4176-e9f6-4dda-8927-4ebdaea01e3a` |
| Working processes | `ace040f2-4e7e-422a-ac4c-03662bafb6d7` |
| Short infobases | `59654aab-c9d2-49e0-ab3b-32f2a3595c5b` |
| Sessions | `a3b238cc-3fa5-406b-a1eb-5347bcae06e2` |

Type notation: `s` string; `g` UUID (where nullable, integer 0 is also accepted);
`b` integer Boolean 0/1; `h` uint16; `i` int32; `u` uint32; `I` int64;
`U` uint64; `n` a nonnegative integer without an established native width here;
`f` floating-point number. Positional integers use ListStream's
signed-magnitude tags. A timestamp-like integer is not necessarily Unix time.

## Clusters

All profiles begin `[uuid,name,[endpointCount,[host,port],...],...]`, port `h`.

| Interface method count | Row fields | Tail types, starting at index 3 |
| --- | ---: | --- |
| 68 | 12 | `i i i i i u i i b` |
| 79 or 80 | 10 | `n n n n n b b` |
| 89 or 96 | 15 | `i i u i b b b b i i b s` |

For the 10-field schema, indices 3..9 are security level, lifetime, expiration,
fault tolerance, load-balancing mode, kill-problem-process flag, kill-dump flag.
For 15 fields, indices 3..8 are security level, expiration, fault tolerance,
load-balancing mode and the two kill flags; index 14 is the recycling schedule.
Do not assign names to the remaining fields without additional evidence.

## Administrators, servers and infobases

An administrator is `[name,description,opaqueString,osUserName,flags]`. Flags
bit 1 enables password authentication and bit 2 enables OS authentication.
The opaque string is not a public display field or a recovered password.

A working-server row starts `[uuid,name,port,host,signedValue,rangeCount]`,
then `rangeCount` entries `[firstPort,lastPort]`, then tail
`s s u u b b h i i i`, then zero, one or two additional strings. Port ranges
are uint16 pairs. Preserve unnamed fields; do not derive administrator access
from a listening port.

An agent short-infobase row is `[uuid,name,description]`. A worker enumeration
row is `[name,uuid,description]` and uses the different collection UUID in
[worker enumeration](12-administration.md#worker-infobases).

## Working processes

```text
[uuid,host,nullableUUID,port,enabled,running,connections,started,
 statistics,pidText,signedUnknown,licenseCount,license...,BooleanTail]
```

`port` is `h`, `enabled` is `b`, `running` is `i`, connections and started are
`U`; the final row has `13 + licenseCount` fields. The 16 statistics fields
have types `f f f f f u f f u i u u u U f U`.
Some raw dates use 10,000 units per second in synthetic decoding, but a complete
calendar/epoch interpretation is not established. Preserve raw values.

## Session entries

A session entry is `[hasData,data?,hasLicense,license?]`. The covered profile
requires `hasData=1`; the optional license is present only for flag 1.
The data block starts with 22 fields:

```text
0 sessionUUID     1 infobaseUUID    2 infobaseName  3 description
4 displayedID:u   5 started:U      6 active:U      7 host:s
8 user:s          9 s             10 h           11 s
12 application:s 13 g             14 u           15 g
16 g             17 s             18 s           19 g
20 firstStatisticsPresent:b       21 secondStatisticsPresent:b
```

The first two UUIDs are nonzero. UUID slots 13,15,16,19 are nullable. A first
statistics flag of 1 is unresolved and unsupported; do not skip a guessed
length. When the second flag is 1, the next field is a statistics list with
these 48 fixed fields:

```text
u U U U U U u U u b u U u U u U U U U b u U u U U s u s u u
I I I U U U U U U b u U u s I U U i
```

The last `i` is a nonnegative metric-entry count. It is followed by that many
pairs `name:s, values:list`, each values list having types
`u U U u u U U I I U U U U U I u u`.

After the optional statistics list, the data block can contain a prefix of
these extension groups, in order. It must end at a complete group boundary:

```text
i | u b | s | g | u s | s | U | u | i | attributes | timezone |
b | b | b | b u u | b | U | g | s h | u | s | s | u
```

`attributes` is `[count,[s,s,u],...]`. `timezone` is `[b,s]`, with an additional
`I` when the Boolean is 1. Its unit is not established. Unknown semantic names
remain unknown even where their serialized types are recovered.

## License structures

License structure is `[kind,payload]`. Kind 0 has payload types `s i`. Kind 1
starts `i b s i s i b`; when the seventh field is 1, fourteen strings follow.
The remaining tail is `s b i i s b s U U u`, giving 17 or 31 fields total.
This describes structure only and establishes no license entitlement, bypass
or permission to alter server state.
