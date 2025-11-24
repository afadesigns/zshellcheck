## Description

<!--
Please include a summary of the change and which issue is fixed.
If this PR adds a new Kata, please describe the check briefly.
-->

Fixes # (issue)

## Type of change

- [ ] `feat`: New feature (non-breaking change which adds functionality)
- [ ] `fix`: Bug fix (non-breaking change which fixes an issue)
- [ ] `docs`: Documentation update
- [ ] `chore`: Maintenance (deps, build, etc.)
- [ ] `refactor`: Code restructuring
- [ ] `test`: Adding missing tests

## Checklist

- [ ] I have read the [CONTRIBUTING.md](CONTRIBUTING.md) guide.
- [ ] I have read the [DEVELOPMENT.md](DEVELOPMENT.md) guide.
- [ ] **Linting**: My code passes `go fmt` and `go vet`.
- [ ] **Tests**: I have added tests for my changes (especially for new Katas).
- [ ] **Integration**: I have run `./tests/integration_test.zsh` and it passes.
- [ ] **Documentation**: I have updated relevant documentation (if applicable).

## For New Katas (If Applicable)

- [ ] Added `pkg/katas/zcXXXX.go`
- [ ] Registered in `pkg/katas/katas.go`
- [ ] Added tests in `pkg/katas/katatests/zcXXXX_test.go`
- [ ] Verified that the ID (`ZCXXXX`) is unique and sequential.