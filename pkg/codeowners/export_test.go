package codeowners

import "github.com/spf13/afero"

func SetFS(newFs afero.Fs) func() {
	oldFS := fs
	fs = newFs

	revert := func() {
		fs = oldFS
	}

	return revert
}



