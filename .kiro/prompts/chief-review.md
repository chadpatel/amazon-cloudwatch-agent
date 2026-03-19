# @chief-review

Evaluate the quality and completeness of the docs-onboarding output and produce a structured review.

## Instructions

1. Read `AGENTS.md` and `docs/arch/overview.md` to understand the project
2. Read every file in `docs/arch/` and `docs/arch/recipes.md`
3. Read `docs/design/_catalog/index.md` and all category files
4. Deeply scan the source code to understand the actual codebase structure

## Evaluation

Analyze the onboarding output against these criteria:

### What Worked
- Are docs consistent in format (Overview → Mermaid → Key Concepts → Patterns → Key Files)?
- Are code examples real (from the actual codebase, not hypothetical)?
- Does each doc cover one concern completely?

### Critical Gap: Pattern Detection
- Do recipes reflect **repeated structures in the codebase**, or are they just summaries of the arch docs?
- Scan all feature/module directories. For each directory structure or file pattern that appears in 2+ features, check if a recipe exists. List the missing ones.

### Tech-Agnostic Gaps
Evaluate whether these universal gaps exist:

1. **No composite workflow recipes** — Are there only atomic recipes, or do composite end-to-end workflows exist too?
2. **No cross-cutting registration points** — Are central files (modified by 3+ features) identified and documented as a checklist?
3. **No recipe from repeated directory structures** — Were feature directory trees compared to find scaffolding patterns?
4. **Design catalog empty** — Are category files populated with existing features, or just placeholder comments?
5. **No recipe coverage validation** — Does every arch doc have a corresponding recipe for its primary use case?
6. **No decision trees** — When multiple approaches exist for the same goal, is there guidance on when to use each?
7. **No boilerplate templates** — For the top 3 most repeated patterns, are there fill-in-the-blanks templates or just step lists?
8. **No failure mode recipes** — Are common error paths, catch blocks, and sentinel values documented as troubleshooting recipes? Search for catch/log.error patterns and check if recipes.md has a "Common Failure Modes" section.
9. **No recipe coverage matrix** — Is there an explicit mapping from each arch doc to its corresponding recipe(s)? Missing mappings indicate undocumented primary use cases.
10. **No pattern-to-recipe traceability** — Does every item in the Phase 2 Pattern Inventory (worklog) have a corresponding recipe in recipes.md? For each inventory subsection, check that every entry above its occurrence threshold maps to a recipe. List any inventory items without recipes.

### Framework-Specific Gaps
- Identify repeated code patterns specific to this repo's language/framework that have no recipe
- List each with: pattern name, how many times it repeats, which files exhibit it

### Test Infrastructure Gaps
- Scan test directories for repeated patterns (mock/shunt subclasses, fake server handlers, fixture layouts, test data factories, test base classes)
- For each pattern appearing 2+ times, check if recipes.md documents it
- List missing test recipes with: pattern name, occurrence count, example files

### What the Script Should Have Done Differently
- List concrete changes to the onboarding methodology (not the repo) that would have prevented the gaps above
- Focus on changes that apply to **any** codebase, not just this one

## Output

Write the review to `chief_onboarding_review.md` in the repository root using this structure:

```markdown
# Review Notes — docs-onboarding Review

## Summary
One paragraph: how many docs, recipes, catalog entries were produced, and the main gaps.

## Current State

### ✅ Completed
- Item that is fully done

### ⏳ Partial
- Item that is partially done with note on what remains

## What Worked
- Strength or positive finding

## Known Gaps

### [Gap Name]
- **Status:** Open | In Progress | Resolved
- **Impact:** High | Medium | Low
- **Priority:** P0 | P1 | P2
- **Recommendation:** What to do about this gap

(Repeat for each gap found)

## Tech-Agnostic Gaps

Evaluate each category and fill in every row — no blank rows allowed:

| # | Category | Status | Evidence |
|---|----------|--------|----------|
| 1 | No composite workflow recipes | Present / Absent | Brief evidence |
| 2 | No cross-cutting registration points | Present / Absent | Brief evidence |
| 3 | No recipe from repeated directory structures | Present / Absent | Brief evidence |
| 4 | Design catalog empty | Present / Absent | Brief evidence |
| 5 | No recipe coverage validation | Present / Absent | Brief evidence |
| 6 | No decision trees | Present / Absent | Brief evidence |
| 7 | No boilerplate templates | Present / Absent | Brief evidence |
| 8 | No failure mode recipes | Present / Absent | Brief evidence |
| 9 | No recipe coverage matrix | Present / Absent | Brief evidence |
| 10 | No pattern-to-recipe traceability | Present / Absent | Brief evidence |

## Recommendations

### High Priority
1. Recommendation

### Medium Priority
1. Recommendation

### Low Priority
1. Recommendation

## Action Items

### Immediate
- [ ] Action item

### Short-term
- [ ] Action item

### Long-term
- [ ] Action item

## Review Cadence
- **Quarterly review** — Re-run chief-review to track gap closure
- **On major changes** — Update immediately when architecture shifts
- **Annual deep dive** — Full re-onboard and comprehensive review
```
