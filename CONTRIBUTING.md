# Contributing to zshellcheck

We welcome contributions! Whether it's adding new Katas, improving the parser, or fixing bugs, your help is appreciated.

## Development Workflow

1.  **Fork and Clone:** Fork the repo and clone it locally.
2.  **Branching:** Create a feature branch for your changes (e.g., `feat/new-kata-zc1099`).
3.  **Tests:** We use Go's standard testing framework. Run tests with:
    ```bash
    go test ./...
    ```
4.  **Linting:** Ensure code is formatted (`go fmt ./...`) and passes basic linting.

## Adding a New Kata

Katas are the core rules of `zshellcheck`. To add one:

1.  **Define the Kata:** Create a new file `pkg/katas/zcXXXX.go`.
2.  **Register:** In the `init()` function, register the Kata with the `RegisterKata` function, specifying the AST node type it targets.
3.  **Implement Logic:** Write the check function that inspects the node and returns a list of `Violation`s.
4.  **Add Tests:** Create `pkg/katas/katatests/zcXXXX_test.go` with test cases covering valid and invalid Zsh code.

### Example Kata

```go
func init() {
    RegisterKata(ast.SimpleCommandNode, Kata{
        ID: "ZC1099",
        Title: "Avoid foo command",
        Description: "The foo command is deprecated.",
        Check: checkZC1099,
    })
}
```
