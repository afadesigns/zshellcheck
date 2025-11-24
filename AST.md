# ZShellCheck AST Reference

This document serves as a reference for the Abstract Syntax Tree (AST) nodes used in ZShellCheck. Understanding these nodes is crucial for writing new Katas (checks).

## Node Types

### `ProgramNode`
The root node of the AST. Contains a list of `Statements`.

### `SimpleCommandNode`
Represents a basic command execution.
- **Example:** `ls -la /tmp`
- **Fields:**
  - `Name` (Expression): The command name (e.g., `ls`).
  - `Arguments` ([]Expression): List of arguments (e.g., `-la`, `/tmp`).

### `ExpressionStatementNode`
A statement that consists of a single expression.
- **Example:** `echo "hello"` (technically a SimpleCommand, but in some contexts can be treated as expression)
- **Used for:** Wrappers around expressions that act as statements.

### `IfStatementNode`
Represents an `if` control structure.
- **Example:**
  ```zsh
  if [[ -f file ]]; then
    echo "exists"
  else
    echo "missing"
  fi
  ```
- **Fields:**
  - `Condition` (BlockStatement)
  - `Consequence` (BlockStatement)
  - `Alternative` (BlockStatement)

### `ForLoopStatementNode`
Represents a `for` loop (both C-style and for-each).
- **Example:** `for i in 1 2 3; do echo $i; done`
- **Example:** `for ((i=0; i<10; i++)); do echo $i; done`
- **Fields:**
  - `Init`, `Condition`, `Post` (Expression): For C-style loops.
  - `Items` ([]Expression): For for-each loops.
  - `Body` (BlockStatement)

### `WhileLoopStatementNode`
Represents a `while` loop.
- **Example:** `while true; do echo "loop"; done`
- **Fields:**
  - `Condition` (BlockStatement)
  - `Body` (BlockStatement)

### `FunctionDefinitionNode`
Represents a function definition.
- **Example:** `my_func() { echo "hi"; }` or `function my_func { ... }`
- **Fields:**
  - `Name` (Identifier)
  - `Body` (BlockStatement)

### `Variable Assignment` (`LetStatementNode`)
Represents a variable assignment (currently parsed as `LetStatement` but Zsh assignments are complex).
- **Example:** `x=10`
- **Fields:**
  - `Name` (Identifier)
  - `Value` (Expression)

### `CommandSubstitutionNode`
Represents backtick command substitution.
- **Example:** `` `date` ``
- **Fields:**
  - `Command` (Expression)

### `DollarParenExpressionNode`
Represents `$()` command substitution.
- **Example:** `$(date)`
- **Fields:**
  - `Command` (Expression)

### `ArithmeticCommandNode`
Represents `(( ... ))` arithmetic evaluation.
- **Example:** `(( x + 1 ))`
- **Fields:**
  - `Expression` (Expression)

### `BracketExpressionNode`
Represents single bracket test `[ ... ]`.
- **Example:** `[ -f file ]`
- **Fields:**
  - `Expressions` ([]Expression)

### `DoubleBracketExpressionNode`
Represents double bracket test `[[ ... ]]`.
- **Example:** `[[ -f file ]]`
- **Fields:**
  - `Expressions` ([]Expression)

### `RedirectionNode`
Represents IO redirection.
- **Example:** `> file.txt`
- **Fields:**
  - `Left` (Expression): Source (optional)
  - `Operator` (string): `>`, `>>`, `<`, etc.
  - `Right` (Expression): Target

### `ArrayAccessNode`
Represents proper array access `${arr[i]}`.
- **Fields:**
  - `Left` (Expression): Array variable.
  - `Index` (Expression): Index.

### `InvalidArrayAccessNode`
Represents improper array access `$arr[i]` (which ZShellCheck flags as an error).

## Visitor Pattern

To inspect the AST, use the `ast.Walk` function with a visitor callback.

```go
ast.Walk(rootNode, func(node ast.Node) bool {
    if cmd, ok := node.(*ast.SimpleCommand); ok {
        // Inspect command...
    }
    return true // Continue traversal
})
```
