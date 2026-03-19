# Clean Code

Rules for writing clean, readable, maintainable code in this repository.

## Rules

### Meaningful Names
Use intention-revealing names for variables, functions, and classes.

What to watch for: Single-letter variables outside tiny loops, generic names like `data`, `temp`, `result`, `info`, abbreviations that obscure meaning, names that require a comment to explain.

### Small Functions
Keep functions short and focused on a single task.

What to watch for: Functions longer than ~20 lines, multiple levels of abstraction in one function, functions that need "and" in their description, deeply nested control flow.

### No Magic Numbers
Replace raw literals with named constants that explain their purpose.

What to watch for: Numeric literals in conditions or calculations, string literals repeated across files, unexplained timeout or threshold values, status codes without named constants.

### DRY (Don't Repeat Yourself)
Every piece of knowledge should have a single authoritative representation.

What to watch for: Copy-pasted logic across functions or files, identical validation in multiple places, repeated error-handling blocks, configuration values hardcoded in several locations.

### No Dead Code
Remove unused code instead of commenting it out or leaving it in place.

What to watch for: Commented-out blocks of code, unreachable branches, unused imports or variables, functions that nothing calls, TODO-marked stubs that are never completed.

### Early Returns
Return early to avoid deep nesting and clarify preconditions.

What to watch for: Deeply nested `if` blocks, long `else` branches that could be flipped, functions where the main logic is indented 3+ levels, error handling buried at the bottom.

### Guard Clauses
Validate inputs and preconditions at the top of a function, then proceed with the happy path.

What to watch for: Validation mixed into the middle of business logic, repeated null/empty checks throughout a function, defensive checks scattered instead of consolidated at entry.
