# Changelog

## [Unreleased]

### Added

- Independent web continuation: Legacy login/logout passed on build 18; V1/V3
  processed real challenges on builds 22/27 but final application-start requests
  failed. Six invocations, preserved first failures, separate cleanup outcomes
  and reproducible frozen candidates are reported without publishing debug tests.
- Observed semicolon-free VERSION declaration in the build-18 publication entry,
  with line-boundary and discovery-versus-login diagnostic requirements.
- Bounded web application-start/cleanup exception patterns from a retained
  build-27 licensing failure, distinguished from password refusal and successful
  logout; descriptions and private runtime details are excluded.
- Password-mode server qualification and phase-specific acceptance criteria:
  native and web profiles, transport policy, credential results, application
  readiness and cleanup remain distinct; no new live password result is claimed.
- Independent password-readiness audit: complete synthetic Legacy and SecPswd
  cleanup scenarios, a corrected stale-status diagnostic and a reproducible
  frozen candidate; 93 local test events passed, without a live password claim.
- Independent public-spec Go implementation acceptance: the same frozen first
  candidate passed direct-worker and manager-discovery Windows login/logout
  against the native-18 laboratory, with separate cleanup acknowledgments.
- Eight synthetic state regressions reproduced late negotiation/context reuse;
  the client consumed the published clarification, and its corrected candidate
  passed both live Windows paths again. Initial and follow-up evidence remain separate.
- Bootstrap/authentication implementation audit, a complete synthetic native
  session and seven native/web password HTTP flows with rejection cases.
- Complete implemented manager/worker bootstrap and native login contracts;
  provider integration, SecPswd V1/V2/V3 and separate web orchestration.
- Exact build-27 n/c policy, administration methods/collections, scalar and
  tracked-form schemas, dynamic-list windows, file/storage contracts and routes.
- Independent-client implementation guide, eight synthetic credential exchanges,
  eight password-key vectors and a synthetic worker bootstrap request.
- Neutral source IDs and a 61-source provenance register; private location
  mappings and research debugging programs remain outside the publication.
- Initial native-18 draft covering framing, selected serialization, bootstrap
  phases, HTTP RPC, session boundaries and exact-build evidence.
- Synthetic vectors, offline checks and a hash-bound source manifest.
- Separate 8.3.27 and 8.5 evidence boundaries and a specification roadmap.

### Corrected from historical research

- Clarified that accepted application login forbids new inner negotiation even
  if no inner provider was created, and made unspecified Unicode casing versions
  explicit for non-ASCII password compatibility.
- Worker selection requires a typed `S` endpoint; the previously documented
  direct-string alternative failed the retained parser.
- First ordinary application sync is 1, while immediate logout uses 2.
- Clarified reply encoding/suffixes, counters, first-HTTP lifetime, nonempty
  SecurityBlob atoms, repeated inner-security PDUs and post-negotiation flags.
- Specified native profile/header selection, web session precedence/rotation,
  POST settings continuation, application validators and credential edge cases.
- HTTPFile continuation chunks omit the ordinary-list `0x80` separator; an
  empty terminal chunk is `01`, following the 2026-11-12 native study.
