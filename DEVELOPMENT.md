# Developer Guide

This document is intended for developers who want to contribute to the ZShellCheck codebase or understand its internal workings.

## Prerequisites

- **Go**: Version 1.18 or higher.
- **Git**: For version control.
- **Make** (Optional): For running build scripts if available.

## Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/afadesigns/zshellcheck.git
    cd zshellcheck
    ```

2.  **Install dependencies:**
    ```bash
    go mod download
    ```

## Building

To build the project from source:

```bash
go build -o zshellcheck cmd/zshellcheck/main.go
```

## Running Tests

We use the standard Go testing framework.

### Run all tests
```bash
go test ./...
```

### Run specific tests
```bash
go test -v pkg/parser/parser_test.go
```

### Integration Tests
We have a script for running integration tests against real Zsh scripts.
```bash
./tests/integration_test.zsh
```

## Debugging

### dumping AST
Currently, there isn't a direct CLI flag to dump the AST (Planned). However, you can use `fmt.Printf("%+v\n", node)` within a temporary walker or test to inspect the structure.

## Creating a New Kata

1.  **Identify the Anti-Pattern**: What Zsh issue do you want to catch?
2.  **Determine the AST Node**: Use [AST.md](AST.md) to find the relevant node type (e.g., `SimpleCommandNode` for command usage, `ForLoopStatementNode` for loops).
3.  **Create the File**: Add `pkg/katas/zcXXXX.go` (next available ID).
4.  **Implement**:
    ```go
    package katas
    import "github.com/afadesigns/zshellcheck/pkg/ast"

    func init() {
        RegisterKata(ast.SimpleCommandNode, Kata{
            ID: "ZCXXXX",
            Title: "Title of your check",
            Description: "Description of what is wrong.",
            Check: checkZCXXXX,
        })
    }

    func checkZCXXXX(node ast.Node) []Violation {
        // Cast node to specific type
        cmd := node.(*ast.SimpleCommand)
        // Check logic...
        return nil
    }
    ```
5.  **Test**: Add `pkg/katas/katatests/zcXXXX_test.go` with test cases.

## Project Structure

- `cmd/`: Entry point.
- `pkg/ast/`: AST definitions.
- `pkg/lexer/`: Tokenizer.
- `pkg/parser/`: Parser.
- `pkg/katas/`: The individual checks.
