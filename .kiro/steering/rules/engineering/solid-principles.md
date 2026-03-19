# SOLID Principles

Rules for applying SOLID principles in this repository.

## Rules

### Single Responsibility Principle (SRP)
A class or module should have only one reason to change.

What to watch for: Classes handling both business logic and persistence, modules mixing UI rendering with data fetching, functions doing formatting and validation together.

### Open/Closed Principle (OCP)
Software entities should be open for extension but closed for modification.

What to watch for: Long `if/else` or `switch` chains that grow with each new variant, editing existing code instead of adding new implementations, feature additions that require touching unrelated files.

### Liskov Substitution Principle (LSP)
Subtypes must be substitutable for their base types without altering correctness.

What to watch for: Overridden methods that throw unexpected exceptions, subclasses that ignore or violate base class contracts, `instanceof` checks to handle specific subtypes differently.

### Interface Segregation Principle (ISP)
Clients should not be forced to depend on interfaces they do not use.

What to watch for: Large interfaces where implementors stub out methods with no-ops, "god interfaces" that bundle unrelated capabilities, classes implementing methods that always throw `NotImplemented`.

### Dependency Inversion Principle (DIP)
Depend on abstractions, not concretions.

What to watch for: High-level modules importing concrete low-level classes directly, constructors that instantiate their own dependencies, no injection points for testing or swapping implementations.
