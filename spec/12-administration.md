# Administration and worker enumeration

Status: implemented and bounded live-confirmed for the build-18 administration
profile; later collection profiles have static/synthetic coverage. This is a
wire-method subset, not the complete COM administration API.

## Contexts

Use the [bootstrap exchange](03-bootstrap.md) with the following selectors.
Application authentication, agent authentication and cluster authentication are
separate scopes. Successful Windows negotiation does not grant all three.

| Purpose | Context kind | Selector | Application | Interface |
| --- | --- | --- | --- | --- |
| Agent administration | `bc9113bc-5b88-4c9f-9fc0-ff7ade5450f0` | `##AdminProcess##:##AdminProcess##` | `COMConsole` | `73cdaa77-2afc-436c-b738-58c772e3a0bf` |
| Worker infobase enumeration | `a2c6af52-7d8b-46ef-9a69-cb1fa650e4f5` | `##AdminProcess##:##AdminProcess##` | `COMConsole` | `7607aa1f-f87e-47a2-8a6b-701c7105ca63` |

For worker enumeration the final bootstrap UUID is zero. Query the interface
with `[0,target,sequence,1,interface]`; a successful reply is
`[1,marker,methodCount]`. Agent profiles observed are 68, 79, 80, 89 and 96
methods. Worker enumeration has 15. The method count selects a schema; it is
not the platform version or an invitation to invoke every ordinal.

## Marker and invocation grammar

Let `S(g)` add the numeric UUID fields Data1, Data2, Data3 and each of the last
eight individual bytes, modulo 2^32. Read the first three fields as the usual
UUID numeric values, not as a sum of all serialized bytes. The marker is:

```text
S(target) + uint16(sequence) + uint16(operation)
  + (operation != 0 ? S(interface) : 0)
  + (operation != 0 && operation != 1 ? uint16(method) : 0)
```

The final result wraps modulo 2^32. Operations are release=0, query=1, call=4.
This is a correlation value, not a cryptographic authenticator. This formula
is established for the agent interface; the worker enumeration reader retains
its middle reply field as opaque instead of imposing this correlation rule.

Calls are binary ListStream records:

```text
[0,target,sequence,4,interface,method,returnMode,parameterCount,
 direction,argument?, ...]
```

An input parameter has direction 0 followed by its generic value; an output
parameter has direction 1 and no request value. Return mode 1 requests a value;
mode 3 is void. Generic values used here are `S(s)=["S",wideString(s)]`,
`B(b)=["B",0|1]`, `G(g)=["#",ae135932-4f94-44df-92c1-c91f15a92848,[1,g]]`
and `I(n)=["#",ba5edc33-8936-4212-9ed2-485e759203aa,n]`.

Replies contain a list followed by top-level integer 2:

```text
getter: [1,marker,1,value,N,directionAndOutputValues...,0],2
void:   [1,marker,3,N,directionAndOutputValues...,0],2
```

There is one reply direction for each argument: 0 for input, or 1 followed by
the returned value for output. Mode 3 can therefore carry output parameters.
Combined generic-return plus output-parameter methods are not covered.

## Known methods

All methods use the agent interface above. `cluster`, `base` and `session`
are generic UUID values. Empty strings are actual input strings, not omitted
arguments. No authentication method is automatically inserted by this table.

| Ordinal | Operation | Ordered arguments | Return mode |
| ---: | --- | --- | ---: |
| 0 | Clusters | none | 1 |
| 19 | Authenticate agent | `S(user),S(password)` | 3 |
| 20 | Agent administrators | none | 1 |
| 9 | Authenticate cluster | `G(cluster),S(user),S(password)` | 3 |
| 10 | Cluster administrators | `G(cluster)` | 1 |
| 13 | Working servers | `G(cluster)` | 1 |
| 3 | Working processes | `G(cluster)` | 1 |
| 5 | Short infobases | `G(cluster)` | 1 |
| 31 | Sessions | `G(cluster),G(zero),OUT(collection)` | 3 |
| 32 | Terminate session | `G(cluster),G(base),G(session),S(message)` | 3 |
| 23 | Set security level | `G(cluster),I(level)` | 3 |
| 30 | Set recycling parameters | profile-dependent, below | 3 |
| 39 | Set fault tolerance | `G(cluster),I(level)` | 3 |

Profile 68 omits the termination message. A session's internal UUID is distinct
from its displayed numeric session ID. Method 31's reply has three directions
`0,0,1,collection`. Collection schemas are in [the next chapter](13-collections.md).

Method 30 for the 10-field cluster profile takes cluster, lifetime, expiration,
load-balancing mode, kill-problem-process flag, kill-dump flag, in that order
(three `I` values followed by two `B` values). The 15-field profile replaces
lifetime with a schedule string. Profile 68 takes cluster plus the seven
settings at row positions 4,5,6,7,9,10,11: six integers then a Boolean. These
are grouped settings: read and preserve the other values before changing one.
This is not a compare-and-swap operation.

Named administration passwords are generic strings in this RPC, not SecPswd
messages. Negotiating SSPI alone does not seal every subsequent record. Empty
authentication and omitted authentication are distinct choices; do not imitate
native implicit retries with alternate credentials. A lost mutation reply has
an unknown outcome and must not trigger automatic replay.

## Worker infobases

After querying the worker interface, send
`[0,target,sequence,4,interface,1,1,0]`. The reply is:

```text
[1,marker,1,["#",33f42b5a-58eb-482c-a6b0-9f92260e5c2b,
  [count,[name,uuid,description],...]],0,0],2
```

Entries are materialized values, not objects requiring additional property
RPCs. Worker row order differs from the agent short-infobase collection.
Anonymous enumeration was confirmed in a bounded loopback configuration;
default Windows credentials are not anonymous and remote access is not implied.

## Errors and limits

The observed administration error shape is
`[2,class,[[class,message,[0],stack,source,0],code],0],2`. Codes 22 and 4
were observed for authentication failure and denied access respectively; they
are not a complete error catalog. Preserve unknown classes and codes.

Later native clients can insert an extra zero before operation 4 in some
states. Its complete transition rule is unresolved. The initial envelope
above is the implemented, bounded accepted profile. Administration creation,
deletion, arbitrary properties and the full 108-method COM inventory cannot
be inferred from these ordinals. Release uses `[0,target,sequence,0,0]` and
expects `[1,marker]`.
