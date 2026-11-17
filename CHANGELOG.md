# Changelog

## [Unreleased]

### Added

- Initial native-18 draft covering framing, selected serialization, bootstrap
  phases, HTTP RPC, session boundaries and exact-build evidence.
- Synthetic vectors, offline checks and a hash-bound source manifest.
- Separate 8.3.27 and 8.5 evidence boundaries and a specification roadmap.

### Corrected from historical research

- HTTPFile continuation chunks omit the ordinary-list `0x80` separator; an
  empty terminal chunk is `01`, following the 2026-11-12 native study.
