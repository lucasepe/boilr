//go:build windows
// +build windows

package pathutil_test

import (
	"path/filepath"
	"testing"

	"github.com/lucasepe/dry/internal/pathutil"
	"github.com/stretchr/testify/require"
)

func TestKnownFolder(t *testing.T) {
	expected := `C:\ProgramData`
	require.Equal(t, expected, pathutil.KnownFolder(windows.FOLDERID_ProgramData, nil, nil))
	require.Equal(t, expected, pathutil.KnownFolder(nil, []string{"ProgramData"}, nil))
	require.Equal(t, expected, pathutil.KnownFolder(nil, nil, []string{expected}))
	require.Equal(t, "", pathutil.KnownFolder(nil, nil, nil))
}

func TestExpandHome(t *testing.T) {
	home := `C:\Users\test`

	require.Equal(t, home, pathutil.ExpandHome(`%USERPROFILE%`, home))
	require.Equal(t, filepath.Join(home, "appname"), pathutil.ExpandHome(`%USERPROFILE%\appname`, home))

	require.Equal(t, "", pathutil.ExpandHome("", home))
	require.Equal(t, home, pathutil.ExpandHome(home, ""))
	require.Equal(t, "", pathutil.ExpandHome("", ""))

	require.Equal(t, home, pathutil.ExpandHome(home, home))
}
