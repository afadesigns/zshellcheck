# Glossary

This document defines common terms and jargon used within the ZShellCheck project and codebase.

## Core Concepts

### Kata
A **Kata** is a specific rule or check implemented in ZShellCheck.
*   Derived from martial arts "katas" (forms/patterns), implying a pattern of code that should be perfected.
*   Each Kata has a unique ID (e.g., `ZC1001`).
*   Example: "Prefer `[[ ... ]]` over `[ ... ]`".

### Violation
A **Violation** is an instance where the code fails a Kata check. It includes:
*   The ID of the failed Kata.
*   A message explaining the issue.
*   The line and column number where it occurred.

## Technical Terms

### AST (Abstract Syntax Tree)
A tree representation of the abstract syntactic structure of the source code.
*   **Node**: An element in the tree (e.g., a command, a loop, an if-statement).
*   ZShellCheck parses Zsh scripts into an AST to analyze the logical structure rather than just the text.

### Lexer / Tokenizer
The component that breaks the raw source code string into meaningful chunks called **Tokens**.
*   **Token**: A pair of type and literal value (e.g., `TOKEN_IF` ("if"), `TOKEN_IDENT` ("my_var")).

### Parser
The component that takes the stream of Tokens from the Lexer and constructs the AST. It handles the grammar rules of Zsh (e.g., ensuring `if` is followed by `then` and closed by `fi`).

### Walker / Visitor
A pattern used to traverse the AST. The **Walker** visits every node in the tree and allows the registered Katas to inspect them.

### Registry
The global store where all available Katas are registered. When the application starts, it loads all Katas into the Registry so they can be looked up by the Walker.

### .zshellcheckrc
The configuration file (YAML format) used to customize ZShellCheck behavior, primarily for disabling specific Katas.
