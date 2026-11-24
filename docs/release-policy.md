# Versioning and Release Policy

## Current stability

Krewire TUI is a pre-1.0 library. Until maintainers approve a stable API contract and publish v1.0.0, consumers should treat minor releases as potentially containing source-breaking changes. Any breaking change must still be deliberate, documented, and accompanied by migration guidance where practical.

No stable release is implied by the version string used in examples or by the existence of a successful CI run.

## Version tags

When releases are approved, tags use Semantic Versioning with the Go module tag form:

- v0.x.y for pre-1.0 development releases.
- v1.0.0 for the first explicitly approved stable API.
- vMAJOR.MINOR.PATCH for subsequent stable releases.

Release notes must describe user-visible changes, compatibility impact, required Go version, and migration steps. Do not create a stable tag solely to satisfy a downstream dependency or because CI is green.

## Stable v1.0.0 entry criteria

Before the first stable tag, maintainers must review and record:

1. The supported public packages, types, functions, and behavioral contracts.
2. The minimum supported Go version and supported build targets.
3. Backward-compatibility expectations and the project's deprecation approach.
4. Consumer integration evidence from active Krewire repositories.
5. Passing CI evidence for formatting, static analysis, unit tests, race tests, and pure-Go builds on the release commit.
6. Release notes and migration guidance for any known compatibility changes.

The v1.0.0 decision is an explicit maintainer decision. Until then, the project remains pre-1.0.

## Release integrity

A release tag must point to a reviewed commit on the default branch whose required CI checks passed. Tags are immutable release identifiers and must not be moved or reused. If a release is withdrawn, publish a follow-up release or advisory rather than rewriting an existing tag.

This repository does not gain an automated publishing workflow from this policy alone. Any future release automation must be reviewed separately, use least-privilege permissions, and publish only from an approved tag after the release checks pass.
