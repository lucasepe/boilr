package dir

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

type FilterFunc func(p string) bool

func ListAll(dir string) []string {
	return List(dir, func(p string) bool {
		return true
	})
}

func List(dir string, accept FilterFunc) []string {
	res := []string{}

	dir = filepath.Clean(dir)
	fs.WalkDir(os.DirFS(dir), ".",
		func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			fi, err := d.Info()
			if err != nil {
				return err
			}

			if !fi.Mode().IsRegular() {
				return nil
			}

			if accept == nil {
				res = append(res, p)
			} else if accept(p) {
				res = append(res, p)
			}

			return nil
		})

	sort.Strings(res)

	return res
}
