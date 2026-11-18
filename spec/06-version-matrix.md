# Version evidence and remaining work

Evidence snapshot: 2026-11-18. The reports referenced here retain their original
observation dates. No new live experiment was performed for this publication.

| Build / profile | Established evidence | Not established |
| --- | --- | --- |
| Windows x64 8.3.18.1208 native | Live baseline connection/Windows login, bounded application operation and cleanup; corrected multi-chunk upload; retained native reply; static serializer studies | Complete protocol, all authentication providers, arbitrary operations or platforms |
| Windows x64 8.3.27.2342 native | Static policy envelope; original COM client query/release against owned synthetic n-policy and c-policy peers | Real build-27 native server acceptance and complete application compatibility |
| 8.5.1.1529 | Selected shipped-code authentication paths and 56 matching synthetic V3 exchanges reported in research | Complete TCP envelope or native live application compatibility |
| Other installed 8.3 clients | Selected administration calls accepted by original clients against a synthetic peer | End-to-end support for servers of those builds |
| HTTP(S) web publication | Separate route and authentication studies | Interchangeability with the native worker transport |

The eight-client administration comparison covers 8.3.12.1529, 8.3.17.1549,
8.3.18.1208, 8.3.19.1522, 8.3.22.1709, 8.3.27.1859, 8.3.27.2130 and
8.3.27.2342. Its earlier build-27.2342 bootstrap failure was followed by separate
successful synthetic policy studies. Preserve the operation and chronology
instead of combining these into a blanket "supported" label.

## Build-27 policy boundary

The exact 27.2342 native client adds an initial policy exchange. The observed
original COM client requested `n`; separate owned peers demonstrated an `n`
reply and acceptance of a `c` reply with encrypted query/release. The mapped
policy codes are `n/y/c/u = 0/1/2/4`. One studied default cipher mode used
DES-EDE3-CBC with a 24-byte key and eight-byte IV.

The [exact-build policy chapter](10-native27-policy.md) now specifies the initial
envelope, asymmetric key exchange, padding and IV progression for n/c. This is
implementable synthetic-peer coverage, not proof of a real build-27 server or
every policy/provider combination.

## Next work

1. Obtain separately authorized real-server evidence for build-27 policy and
   native SecPswd/provider combinations, without broadening synthetic claims.
2. Recover the complete build-85 native envelope and policies y/u.
3. Extend public synthetic conformance data for unknown parser variants and
   schema tails; keep research debugging programs and captured tokens private.
4. Recover general form descriptors, control/layout changes, lazy collections,
   tree lists and asynchronous continuation; do not infer them from route names.
5. Specify MFA, external authentication and password changes only after their
   message contracts are established. See the [gap register](17-routes-and-errors.md#explicit-gaps).

See the [evidence register](../evidence/README.md) for provenance and
[contribution requirements](../CONTRIBUTING.md) for extending these claims.
