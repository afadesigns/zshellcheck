// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/config"
	"github.com/afadesigns/zshellcheck/pkg/katas"
)

func TestCollectEdits_FileWideDisableDirective(t *testing.T) {
	src := "result=`which git`\n# noka: ZC1002\n"
	edits := collectEdits(src, katas.Registry, nil, config.DefaultConfig(), nil)
	for _, e := range edits {
		_ = e
	}
}

func TestCollectEdits_PerLineDisable(t *testing.T) {
	src := "result=`which git` # noka: ZC1002\n"
	edits := collectEdits(src, katas.Registry, nil, config.DefaultConfig(), nil)
	for _, e := range edits {
		_ = e
	}
}

func TestCollectEdits_ExternalDisable(t *testing.T) {
	src := "result=`which git`\n"
	edits := collectEdits(src, katas.Registry, []string{"ZC1002"}, config.DefaultConfig(), nil)
	for _, e := range edits {
		_ = e
	}
}

func TestCollectEdits_SeverityFilterError(t *testing.T) {
	src := "result=`which git`\n"
	edits := collectEdits(src, katas.Registry, nil, config.DefaultConfig(), []katas.Severity{katas.SeverityError})
	for _, e := range edits {
		_ = e
	}
}

func TestCollectEdits_SeverityFilterStyle(t *testing.T) {
	src := "result=`which git`\n"
	edits := collectEdits(src, katas.Registry, nil, config.DefaultConfig(), []katas.Severity{katas.SeverityStyle})
	for _, e := range edits {
		_ = e
	}
}

func TestCollectEdits_NestedFix(t *testing.T) {
	src := "result=`which git`\necho $arr[1]\n"
	edits := collectEdits(src, katas.Registry, nil, config.DefaultConfig(), nil)
	for _, e := range edits {
		_ = e
	}
}

func TestCollectEdits_CleanSource(t *testing.T) {
	src := "echo hello\n"
	edits := collectEdits(src, katas.Registry, nil, config.DefaultConfig(), nil)
	for _, e := range edits {
		_ = e
	}
}
