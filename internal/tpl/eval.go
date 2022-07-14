package tpl

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/lucasepe/boilr/internal/dir"
	"github.com/lucasepe/boilr/internal/envsubst"
	"github.com/lucasepe/boilr/internal/ignore"
)

type EvalOpts struct {
	OriginDir string
	TargetDir string
}

func EvalDir(indir, outdir string, toSkip []string) (err error) {
	acceptFunc := func(rules []string) dir.FilterFunc {
		if len(rules) == 0 {
			return func(p string) bool {
				return true
			}
		}

		return func(p string) bool {
			object := ignore.CompileIgnoreLines(rules...)
			return !object.MatchesPath(p)
		}
	}

	files := dir.List(indir, acceptFunc(toSkip))

	for _, el := range files {
		fmt.Printf(" > file: %s\n", el)
		src := filepath.Join(indir, el)
		dst, err := EvalString(filepath.Join(outdir, filepath.Dir(el)))
		if err != nil {
			return err
		}

		err = EvalFile(src, dst)
		if err != nil {
			return err
		}
	}

	return nil
}

func EvalFile(filename string, outdir string) (err error) {
	fr, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fclose(fr, &err)

	var res string
	res, err = EvalReader(fr)
	if err != nil {
		return err
	}

	if outdir == "" {
		fmt.Println(res)
		return nil
	}

	name := strings.TrimSuffix(filepath.Base(filename), ".tpl")
	name, err = EvalString(name)
	if err != nil {
		return err
	}

	dst := filepath.Join(outdir, name)
	err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fclose(out, &err)

	_, err = io.Copy(out, strings.NewReader(res))
	return err
}

func EvalReader(r io.Reader) (string, error) {
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, r); err != nil {
		return "", err
	}

	return envsubst.EvalEnv(buf.String())
}

func EvalString(s string) (string, error) {
	return EvalReader(strings.NewReader(s))
}

/*
func feval(src, dst string) error {
	dst = strings.TrimSuffix(dst, ".tpl")
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fclose(out, &err)

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fclose(in, &err)

	scanner := bufio.NewScanner(in)
	// optionally, resize scanner's capacity for lines over 64K...
	for scanner.Scan() {
		var line string
		line, err = envsubst.EvalEnv(scanner.Text())
		if err != nil {
			return err
		}

		if _, err = out.WriteString(line); err != nil {
			return err
		}

		if _, err = out.WriteString("\n"); err != nil {
			return err
		}
	}

	return scanner.Err()
}
*/
// fclose ANYHOW closes file,
// with asiging error raised during Close,
// BUT respecting the error already reported.
func fclose(f *os.File, reported *error) {
	if err := f.Close(); *reported == nil {
		*reported = err
	}
}
