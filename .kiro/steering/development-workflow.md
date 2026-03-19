# Development Workflow

Steering for choosing the right workflow, offering pre-commit reviews,
and keeping documentation in sync with code changes.

## Workflow Selection

Choose the right approach based on task characteristics.

### Routing Table

| Task Type | Recommended Workflow | Why |
|-----------|---------------------|-----|
| Unclear scope or multi-file change | `@ralph-plan` | Plan first to avoid rework |
| Estimated effort > 1 day | `@ralph-plan` | Needs decomposition and tracking |
| Exploratory or design-heavy work | `@ralph-plan` | Captures decisions in design.md |
| Single-file bug fix with clear root cause | Direct implementation | Plan overhead not justified |
| Small refactor (rename, extract, inline) | Direct implementation | Scope is already clear |
| Config or dependency update | Direct implementation | Mechanical change |

### Decision Criteria

Suggest `@ralph-plan` when:
- Scope is unclear or spans multiple files
- The change touches public APIs or architecture
- Estimated effort exceeds one day

Suggest direct implementation when:
- The fix is a single-file change
- The bug has a clear, identified root cause
- The refactor is small and self-contained

## Pre-Commit Review

After the developer makes code changes, suggest:

> Would you like me to review these changes before committing?
> I can run `@chief-review`.

Do NOT auto-run the review — only suggest it. The developer decides.

## Documentation Update Triggers

After significant code changes, check whether documentation needs
updating and suggest it to the developer.

### Triggers

- **New file or module created** — check if the architecture docs
  cover it. Suggest: "I noticed you added a new module. Should I
  update the relevant arch doc?"

- **Public function or API changed** — check if recipes or usage
  docs reference it. Suggest: "This API changed — should I update
  recipes.md?"

- **New directory matching a known pattern** — suggest updating the
  pattern inventory in recipes or arch docs.

Do NOT auto-update docs. Suggest the update and let the developer
decide.

## Guardrails

- Never auto-commit — always let the developer commit explicitly
- Never suppress errors — surface them clearly
- Always explain failures — include what went wrong and why
- Prioritize the developer's explicit requests over these guidelines
