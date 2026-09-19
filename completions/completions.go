// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package completions

import _ "embed"

//go:embed zsh/_zshellcheck
var zsh string

// Zsh returns the completion function to install as _zshellcheck in fpath.
func Zsh() string {
	return zsh
}
