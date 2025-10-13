package errors_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

func TestNolintIsForbidden(t *testing.T) {
	// This test checks that there are no nolint:xxx directives in the comments.
	// This is done specifically to prevent attempts to bypass documentation requirements.
	// Unfortunately, golangci-lint itself does not support this.

	err := filepath.Walk(".", func(path string, info fs.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		// Parse the file, getting all comments.
		var fset token.FileSet
		tree, err := parser.ParseFile(&fset, path, nil, parser.ParseComments|parser.AllErrors)
		if err != nil {
			t.Errorf("parse %s: %v", path, err)
			return nil
		}

		// Require that linters are not "silenced" by prohibiting nolint:* in them.
		for _, group := range tree.Comments {
			for _, cmt := range group.List {
				checkCommentText(t, &fset, cmt)
			}
		}

		return nil
	})
	if err != nil {
		t.Errorf("scan over go code of this repository: %v", err)
	}
}

func checkCommentText(t *testing.T, fset *token.FileSet, cmt *ast.Comment) {
	text := cmt.Text
	switch {
	case strings.HasPrefix(text, "//"):
		text = strings.TrimPrefix(text, "//")
	case strings.HasPrefix(text, "/*"):
		text = strings.TrimPrefix(text, "/*")
	}

	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "nolint:") {
		t.Errorf(
			"\r%s nolint:<linter> directives are not allowed in this project",
			fset.Position(cmt.Pos()),
		)
	}
}
