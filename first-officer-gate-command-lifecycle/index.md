---
id: 6yyyyemkqwsett3g1c991w9f
title: Make First Officers operate the recorded gate lifecycle
status: ideation
source: "Durable-decisions dogfood audit: PRs #557/#560 shipped gate commands without the planned FO operating contract, 2026-07-23"
started: 2026-07-23T02:01:56Z
completed:
verdict:
score:
worktree:
issue:
sprint: durable-decisions
---

Make the normal First Officer gate path use the landed 3k/h1 commands so a presented decision is durably recorded, validated, checked for eligibility, and consumed before workflow advancement or dispatch.

## Problem

The recorder and application commands exist, but the shipped First Officer contract still assembles and presents gates only in prose. PR #557 changed no First Officer skill files despite budgeting that integration; PR #560 likewise added eligibility and consume without an FO caller. This sprint used the commands only through manual dogfood directives, so an ordinary agent can still approve and advance without the durable lifecycle.

## Minimum value demonstration seed

In one fixture-backed live First Officer journey, package an exact validation Briefing, record and validate it, record an evidence-bearing delegated approval, observe eligibility, consume the application, and only then advance and dispatch. Two controls must fail closed without status mutation: a revise decision routes through feedback rather than consume, and a hold or ineligible approval does not advance. Removing any recorder command from the FO procedure must make the journey fail.

## Boundary

This task owns the First Officer invocation contract and its behavioral proof. It reuses the landed 3k recorder and h1 eligibility/consume commands; it does not add a recorder schema, duplicate gate judgment, change presentation UI, or implement advisory correction rounds.
