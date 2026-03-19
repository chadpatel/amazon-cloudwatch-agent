# @catalog-reconcile

Audit and reconcile the design catalog to keep it accurate and up-to-date.

## Instructions

1. Read `docs/ralph.md` to understand the plan format
2. Create a new feature directory `docs/design/YYMMDD-catalog-reconcile/` (use today's date)
3. Create `plan.md` with tasks to:

### Phase 1: Find Uncatalogued Designs
- List all design directories in `docs/design/` (YYMMDD-name format)
- For each directory, check if it appears in any `docs/design/_catalog/*.md` file
- Create tasks to add missing entries to appropriate category files

### Phase 2: Update Design Statuses
- For each catalogued design, read its `plan.md`
- Count completed (✅) vs incomplete tasks
- Create tasks to update statuses that differ from actual state:
  - All tasks done → Status: ✅ Done
  - Any task in progress → Status: 🟡 In Progress
  - All tasks not started → Status: ⬜ Planned

### Phase 3: Verify Cross-References
- Check each design's "Related" field references exist
- Create tasks to fix invalid references

### Phase 4: Category Placement
- Identify designs that touch multiple domains but only appear in one category
- Create tasks to add entries to additional relevant categories

## Entry Format for New Catalog Entries

```markdown
### YYMMDD - design-name
Brief description of what this design covers.

- **Tags**: `tag1`, `tag2`, `tag3`
- **Related**: [YYMMDD-other-design](../YYMMDD-other-design/)
- **Status**: ✅ Done | 🟡 In Progress | ⬜ Planned
```

4. Tell the user the folder path
5. Suggest running `@ralph-plan-review` to validate, then Ralph to execute
