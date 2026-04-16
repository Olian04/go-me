---
name: go-quality-reviewer-subagent
description: Delegate read-only Go quality checks to `.cursor/agents/go-quality-reviewer.md`.
---

# Go quality reviewer subagent

Use this skill when:

- The user asks for a Go quality or standards review.
- You completed substantive Go changes and want verification.

Delegate scope explicitly (diff, paths, or branch), then return findings. The subagent only reports; it does not edit files.
