package tui

import (
	"mime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtensionByTypeJson(t *testing.T) {
	assert.Equal(t, ".json", extensionByType("application/json", ".dat"))
}

func TestExtensionByTypeTxt(t *testing.T) {
	assert.Equal(t, ".txt", extensionByType("text/plain", ".dat"))
}

func TestExtensionByTypeHtml(t *testing.T) {
	assert.Equal(t, ".html", extensionByType("text/html", ".dat"))
}

func TestExtensionByTypeCss(t *testing.T) {
	assert.Equal(t, ".css", extensionByType("text/css", ".dat"))
}

func TestExtensionByTypeUnknown(t *testing.T) {
	unknown := "made/up"
	exts, err := mime.ExtensionsByType(unknown)
	require.Empty(t, exts, "mime package should not yield results for '%s'", unknown)
	require.NoError(t, err, "mime package should not return an error for '%s'", unknown)
	assert.Equal(t, ".fallback", extensionByType(unknown, "fallback"))
}

func TestExtensionByTypeUnknownEmptyFallback(t *testing.T) {
	assert.Equal(t, "", extensionByType("made/up", ""))
}
