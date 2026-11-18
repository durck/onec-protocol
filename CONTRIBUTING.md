# Contributing

Prefer a small correction with reproducible evidence over a broad compatibility
claim. Write maintained documentation, code and comments in English.

For each protocol change, include:

1. Exact platform build, relevant OS/architecture and transport profile.
2. Evidence type: static, synthetic, observed and/or live-confirmed.
3. Input identities/hashes, the relevant source or receipt, and what another
   contributor can actually reproduce from public materials.
4. A field layout or state transition, including byte order, units, valid
   values, unknown fields and failure cases.
5. A synthetic example and a description of its verification. State whether it
   came from an independent original component or the implementation under test.

Keep a later correction visible in the affected chapter and CHANGELOG. Update
the version matrix when the strength or scope of evidence changes. Do not
silently promote synthetic acceptance to real-server support.

Publish minimal synthetic fixtures. Do not submit real credentials, captured authentication
tokens, cookies, tickets, captured session state, private host inventories,
proprietary binaries or bulk decompiler/capture output. Fixed protocol constants
and input hashes can identify a finding without redistributing its raw inputs.
Clearly label public toy credentials and deterministically constructed tokens.
Do not copy debugging tests, live probes, research harnesses, local absolute
paths, private directory/repository names or workstation identifiers. Use neutral
source IDs and hashes for provenance; keep location mappings outside this project.
Review every staged file before publication, including JSON and generated metadata.

Run from the repository root before submitting a change:

```sh
go test ./...
go vet ./...
```

Default checks must remain offline. This repository owns specifications and
small conformance examples; full client implementations and live probes belong
in separate projects. List every maintained chapter in the root README.

For an independent implementation experiment, record the allowed public inputs
and begin with a fresh implementation context. Freeze source and binary hashes
before live feedback. Keep laboratory setup and credentials with the operator;
publish any missing protocol rule before the implementation consumes it.
Distinguish instructional input isolation from an OS-enforced sandbox, initial
attempts from later corrections, and implementation tests from specification CI.
Preserve unsuccessful attempts and report authentication, logout, release and
local cleanup separately. Only sanitized results and minimal synthetic examples
belong here; keep the experimental client and debugging tests separate.
