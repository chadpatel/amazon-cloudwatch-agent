# @code-audit

Audit the codebase against documented patterns and generate a fix plan.

## Instructions

1. **Discover patterns and rules** — Read all files in `docs/arch/` and `AGENTS.md`. Extract:
   - Architecture patterns (component structure, API clients, state management, data fetching, routing, etc.)
   - Coding standards (type safety, import conventions, comment style, naming, linter rules)
   - Any project-specific conventions or constraints

   Compile these into a single list of patterns and rules to check against.

2. **Identify source directories** — Scan the repository root and source tree to find all major directories containing source code. Exclude `node_modules`, `dist`, `build`, `.git`, `docs`, and other non-source directories. Each directory becomes a subagent assignment.

3. **Spawn audit subagents in parallel** — For each source directory, spawn a `ralph-auditor` subagent that will:
   - Receive the full list of discovered patterns and rules
   - Receive the finding format (below)
   - Read all files in its assigned directory
   - Check each file against all four audit categories:
     - **Architecture Compliance** — violations of patterns from `docs/arch/`
     - **Coding Standards** — violations of rules from `AGENTS.md`
     - **Performance & Memory** — re-execution from unstable inputs, missing cleanup in lifecycle hooks, stale captures, inefficient data structures, missing async error handling, eagerly loaded large modules
     - **Abstraction Opportunities** — patterns repeated 3+ times, missing shared helpers, inconsistent approaches, hardcoded values that should be constants
   - Return structured findings using the finding format

4. **Aggregate findings** — Collect all subagent findings. Deduplicate and group by category, then sort by severity (Critical → High → Medium → Low). Compute summary stats: total findings by severity, most-violated patterns, directories with most issues.

5. **Write code-audit-report.md** — Write the aggregated report to `code-audit-report.md` in the repository root using this structure:

   ```markdown
   # Code Audit Report

   ## Summary
   Total findings, breakdown by severity, most-violated patterns, most-affected directories.

   ## Architecture Compliance
   (Findings sorted by severity)

   ## Coding Standards
   (Findings sorted by severity)

   ## Performance & Memory
   (Findings sorted by severity)

   ## Abstraction Opportunities
   (Findings sorted by severity)
   ```

6. **Create fix plan** — Create `docs/design/YYMMDD-code-audit-fixes/plan.md` (where YYMMDD is today's date) as a standard Ralph plan with:
   - Phases ordered by severity (Critical first, then High, Medium, Low)
   - Atomic fix tasks — one per finding or group of closely related findings
   - A final phase that verifies docs affected by fixes are updated

7. **Report to user** — Tell the user:
   - The report is at `code-audit-report.md`
   - The fix plan is at `docs/design/YYMMDD-code-audit-fixes/plan.md`
   - Suggest running `@ralph-plan-review` to validate the fix plan, then `./ralph.sh code-audit-fixes` to apply fixes

## Finding Format

Each subagent must return findings in this format:

```markdown
## Finding: [Category] - Brief Description

**Location**: `path/to/file.ext:line`
**Severity**: Critical | High | Medium | Low
**Pattern Violated**: Name of pattern from docs/arch/ or AGENTS.md
**Issue**: What's wrong
**Recommendation**: How to fix
```
