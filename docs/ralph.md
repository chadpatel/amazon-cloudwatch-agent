# Ralph

Ralph is an AI agent that implements features iteratively from structured plans.

## Creating Work for Ralph

Ralph reads a `plan.md` file and executes tasks one iteration at a time. Each feature gets a subdirectory:
```
docs/design/<YYMMDD-feature-name>/
├── design.md         # Feature design document (optional)
├── plan.md           # Implementation plan
├── .ralph/
│   └── .iteration    # Current iteration number (starts at 0)
└── worklog.md        # Ralph's progress notes (created by Ralph)
```

### Catalog Integration

After creating a new feature directory and plan.md, add an entry to the design catalog:

1. Use today's date as prefix in YYMMDD format (e.g., `260211-my-feature`)
2. Name your directory `YYMMDD-feature-name`
3. Add entries to relevant category files in `docs/design/_catalog/`:
   - A design may appear in multiple categories if it spans domains
   - Follow the existing entry format with description, tags, related designs, and status
   - Set status to `🟡 In Progress` when starting, update to `✅ Done` when complete
4. **Every plan must include a final task** to update the catalog entry status from `🟡 In Progress` to `✅ Done` — Ralph won't do it unless it's in the plan

See `docs/design/_catalog/index.md` for the list of categories and entry format.

## Writing Effective Plans for Ralph

### Task Structure

Tasks must be **atomic** - one task = one file change. Ralph increments `.ralph/.iteration` and works through tasks sequentially.

```markdown
## Tasks

### Phase 1: Component Implementation

- [ ] **1.1** Create `path/to/new-file.tsx` - brief description
- [ ] **1.2** Update `path/to/existing-file.ts` - what to add/change
```

### Reference Files Section

Always include files Ralph should examine for patterns:

```markdown
## Reference Files

- `path/to/similar-component` - Pattern for X
- `path/to/existing-hook` - How Y is done
```

### Code Patterns

Include concrete code examples for patterns Ralph needs. Ralph has no memory between iterations, so each plan must be self-contained:

````markdown
### Pattern Name

```language
// Actual code Ralph can adapt
import { Thing } from 'package';

export function example() {
  // ...
}
```
````

### File Changes Table

Map files to tasks so Ralph knows what each iteration touches:

```markdown
## File Changes

| File | Change | Task |
|------|--------|------|
| `path/to/file` | New | 1.1 |
| `path/to/other` | Add X | 1.2 |
```

## Commit Behavior

Ralph does NOT commit unless explicitly instructed. Rules:

1. **Do not commit** unless the plan or user explicitly asks for it
2. **Never use `--no-verify`** - pre-commit hooks must run
3. **If tests fail, fix them** - even if unrelated to the current task, fix failing tests before committing

## Common Patterns in This Codebase

Put codebase-specific patterns in the feature plan itself, not here. This keeps Ralph's plans self-contained.
