# Scope and evidence

## Profile

The `native-18` profile describes the observed native TCP transport and HTTP RPC
subset associated with Windows x64 1C:Enterprise **8.3.18.1208**. This is a draft
interoperability description. It is neither a vendor specification nor a claim
that all operations, platforms or builds use these layouts.

In this draft, **must** describes a requirement of the documented subset or a
clearly identified local validation rule. **Unknown** means that the evidence
does not establish a contract. An unknown field retains its numeric position
and observed value; a plausible name does not establish its meaning.

`HTTPFile`, `SecurityBlob`, `target`, `Control2` and `Control3` are research and
implementation labels. They are not asserted to be official IDL names.

## Evidence labels

| Label | What it establishes |
| --- | --- |
| Static | A path or layout in identified shipped code; not execution of that path |
| Synthetic | Behavior against constructed data or an owned model peer |
| Observed | Behavior of an original component or retained capture, within the recorded conditions |
| Live-confirmed | A bounded exchange against an actual configured server of the stated build |

These are dimensions, not interchangeable scores. An original native client
talking to a synthetic server is **observed + synthetic**, not real-server
acceptance. A model agreeing with itself is not independent evidence.

Each chapter cites entries in the [evidence register](../evidence/README.md).
The source manifest pins report/source hashes and revisions. Those sources are
retained separately and are not necessarily publicly accessible. The public
synthetic vectors are reproducible; the historical native/live claims require
the retained evidence or an independently authorized replication.

## Boundaries

Included: transport framing, selected scalar encodings, bootstrap phase order,
query/release request shapes, HTTP RPC envelopes, native file chunks and session
ownership distinctions.

Not yet specified end to end: all bootstrap response variants, password proof
algorithms, administration interfaces, managed forms, application metadata,
all errors, compression variants, multiplexing and complete newer-build policy
negotiation. [The roadmap](06-version-matrix.md#next-work) records these gaps.

Byte examples are synthetic. Dynamic GUIDs, account names and application state
from captures are not published. Fixed interface/type GUIDs are protocol
identifiers. Exact shipped-image hashes identify inputs; they do not distribute
those images.

Implementation size caps, timeouts, nesting limits and retry policies are local
choices unless the evidence explicitly identifies them as native requirements.
For example, the research client's 512-KiB HTTPFile cap is **not** a claimed
maximum size of the 1C protocol.
