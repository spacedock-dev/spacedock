---
title: Bound Claude live failure classification to the phase it actually observes
status: ideation
score: 0.7
source: Captain-approved recovery item 1 on 2026-07-10, after review of merged PR #490 showed classifyBootPreambleFailure scans the full transcript and can label later contract-lookup searches as boot-preamble failures before the scenario assertion is evaluated.
id: p4h6a5wcqe5ddkhnmrac1w9a
started: 2026-07-10T12:56:28Z
---

Separate true boot-orientation failures from later contract-lookup hunts, and preserve the scenario's primary failure evidence so the live harness reports the phase and cause it actually observed.
