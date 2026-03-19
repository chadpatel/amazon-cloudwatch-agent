# Design Patterns

Rules for applying design patterns pragmatically in this repository.

## Rules

### Composition Over Inheritance
Favor composing objects from small, focused collaborators rather than building deep inheritance hierarchies.

What to watch for: Inheritance chains deeper than two levels, base classes that accumulate unrelated behavior, subclasses overriding most of the parent's methods, "is-a" relationships that are really "has-a".

### Factory Pattern
Use factories to encapsulate object creation when the concrete type depends on runtime context.

What to watch for: `new` calls with complex conditional logic inlined at the call site, switch statements that pick a class to instantiate, constructors duplicated across multiple callers with the same branching.

### Strategy Pattern
Extract interchangeable algorithms behind a common interface so callers can swap behavior without changing their own code.

What to watch for: `if/else` or `switch` blocks that select an algorithm variant, duplicated functions that differ only in one step, hardcoded behavior that should be configurable or testable in isolation.

### Observer Pattern
Decouple event producers from consumers by letting interested parties subscribe to notifications instead of being called directly.

What to watch for: A class that directly calls multiple unrelated dependents after a state change, tightly coupled notification chains, adding a new listener requires editing the producer, circular dependencies between notifier and listeners.
