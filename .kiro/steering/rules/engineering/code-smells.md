# Code Smells

Rules for recognizing and fixing common code smells in this repository.

## Rules

### Long Methods
Keep methods short and focused — a method should do one thing at one level of abstraction.

What to watch for: Methods exceeding ~20 lines, multiple levels of nesting, comments separating "sections" within a single method, scrolling required to read the full body. Fix by extracting helper methods with descriptive names.

### Large Classes
A class should represent a single concept — split classes that accumulate unrelated responsibilities.

What to watch for: Files exceeding ~300 lines, classes with many unrelated fields, groups of methods that only use a subset of the class's state. Fix by extracting cohesive groups into their own classes.

### Feature Envy
A method should primarily use data from its own class, not reach into another object's internals.

What to watch for: Methods that call multiple getters on the same foreign object, logic that would read more naturally as a method on the other class. Fix by moving the method to the class whose data it uses most.

### Data Clumps
Groups of values that always travel together should be bundled into their own type.

What to watch for: The same three-or-more parameters repeated across multiple method signatures, parallel arrays or maps that represent a single concept. Fix by introducing a class, struct, or record to hold the group.

### Primitive Obsession
Use domain types instead of raw primitives for values that carry meaning or constraints.

What to watch for: Strings used for emails, IDs, currencies, or statuses; integers used for temperatures, ports, or quantities with validation scattered across callers. Fix by wrapping the primitive in a value object that enforces its rules.

### God Objects
No single object should know or do everything — break apart objects that centralize too much knowledge.

What to watch for: A class that most other classes depend on, objects passed everywhere as a "context" bag, classes whose name includes "Manager", "Handler", or "Processor" with hundreds of lines. Fix by distributing responsibilities to focused collaborators.
