# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Plan mode default

- Enter plan mode for ANY non-trivial task (3+ steps or achitectural decisions).
- If something goes sideways, STOP and re-plan immediately. Don't keep pushing.
- Use plan mode for verification steps, not just building.
- Always asking questions to get more clarity
- Write detailed specs upfront to reduce ambiguity
- Always breakdown tasks with its dependencies

## Subagent and Agent teams  strategy

- Use subagents liberally to keep main context window clean
- Offload research, exploration, and parallel analysis to subagents
- For complex problems, throw more compute at it via subagents
- One task per subagents for focused execution


## Git Workflow

- ALWAYS develop based on `master` branch
