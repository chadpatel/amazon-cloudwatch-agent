# @doc-verify

Verify and update architecture docs and READMEs to match actual code.

## Documentation Standard

- **Do document**: Purpose, setup, key concepts, public API, non-obvious patterns
- **Don't document**: Implementation details readable from code, every function/type

## Instructions

1. Read `docs/ralph.md` to understand the plan format
2. Create a new feature directory `docs/design/YYMMDD-doc-verification/` (use today's date)
3. Scan for all existing docs to verify (all files in `docs/arch/*.md`, root `README.md`)
4. Create `plan.md` with one task per doc to verify

Each task should:
1. Read the doc and understand its scope
2. Verify claims against actual code
3. Fill gaps where appropriate
4. Fix any inaccuracies

5. Tell the user the folder path
6. Suggest running `@ralph-plan-review` to validate, then Ralph to execute
