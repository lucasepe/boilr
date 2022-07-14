package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/lucasepe/boilr/internal/pathutil"
	"github.com/lucasepe/boilr/internal/tpl"
	"github.com/lucasepe/dotenv"
	flag "github.com/spf13/pflag"
)

const (
	appName = "boilr"
	banner  = `
█▀▄ █▀█ █▄█ 
█▄▀ █▀▄  █  {{VERSION}}
`
)

func Run(ver, bld string) error {
	o := &boilrOpts{version: ver, build: bld}
	configureFlags(o)

	if o.help {
		flag.CommandLine.Usage()
		return nil

	}
	err := o.complete(flag.Args())
	if err != nil {
		return err
	}

	return o.run(flag.Args())
}

func configureFlags(o *boilrOpts) {
	flag.CommandLine.Usage = func() {
		version := fmt.Sprintf("[%s / %s]", o.version, o.build)
		cover := strings.ReplaceAll(banner, "{{VERSION}}", version)
		fmt.Print(cover, "\n")
		fmt.Printf("Text templates using ${var} expansion syntax.\n\n")

		fmt.Print("USAGE:\n\n")
		fmt.Printf("  %s [flags] [template]\n\n", appName)

		fmt.Print("EXAMPLE(s):\n\n")
		fmt.Printf("  %s -e 'var=World' 'Hello, ${var}!'\n\n", appName)
		fmt.Printf("  echo 'Hello ${var,,}' | %s -e 'var=world'\n\n", appName)
		fmt.Printf("  %s -s /path/to/my-vars.env -f /path/to/my-template-file.txt\n\n", appName)
		fmt.Printf("  %s -e var=val -o ../my-output-dir -d /path/to/my-template-dir\n", appName)
		fmt.Printf("\n")

		fmt.Print("SUPPORTED FUNCTIONS:\n\n")
		w := tabwriter.NewWriter(os.Stdout, 2, 2, 2, ' ', 0)
		fmt.Fprintln(w, "  ${var}\tValue of `$var`")
		fmt.Fprintln(w, "  ${#var}\tString length of `$var`")
		fmt.Fprintln(w, "  ${var^}\tUppercase first character of `$var`")
		fmt.Fprintln(w, "  ${var^^}\tUppercase all characters in `$var`")
		fmt.Fprintln(w, "  ${var,}\tLowercase first character of `$var`")
		fmt.Fprintln(w, "  ${var,,}\tLowercase all characters in `$var`")
		fmt.Fprintln(w, "  ${var:n}\tOffset `$var` `n` characters from start`")
		fmt.Fprintln(w, "  ${var:n:len}\tOffset `$var` `n` characters with max length of `len`")
		fmt.Fprintln(w, "  ${var#pattern}\tStrip shortest `pattern` match from start")
		fmt.Fprintln(w, "  ${var##pattern}\tStrip longest `pattern` match from start")
		fmt.Fprintln(w, "  ${var*pattern}\tStrip shortest `pattern` match from end")
		fmt.Fprintln(w, "  ${var**pattern}\tStrip longest `pattern` match from end")
		fmt.Fprintln(w, "  ${var-default}\tIf `$var` is not set, evaluate expression as `$default`")
		fmt.Fprintln(w, "  ${var:-default}\tIf `$var` is not set or is empty, evaluate expression as `$default`")
		fmt.Fprintln(w, "  ${var=default}\tIf `$var` is not set, evaluate expression as `$default`")
		fmt.Fprintln(w, "  ${var:=default}\tIf `$var` is not set or is empty, evaluate expression as `$default`")
		fmt.Fprintln(w, "  ${var/pattern/replacement}\tReplace as few `pattern` matches as possible with `replacement`")
		fmt.Fprintln(w, "  ${var//pattern/replacement}\tReplace as many `pattern` matches as possible with `replacement`")
		fmt.Fprintln(w, "  ${var/#pattern/replacement}\tReplace `pattern` match with `replacement` from `$var` start")
		fmt.Fprintln(w, "  ${var/*pattern/replacement}\tReplace `pattern` match with `replacement` from `$var` end")
		w.Flush()
		fmt.Printf("\n")

		fmt.Print("FLAGS:\n\n")
		flag.CommandLine.SetOutput(os.Stdout)
		flag.CommandLine.PrintDefaults()
		flag.CommandLine.SetOutput(ioutil.Discard) // hide flag errors
		fmt.Println()

		fmt.Println("Crafted with passion by Luca Sepe @ https://github.com/lucasepe")
	}

	flag.CommandLine.SetOutput(ioutil.Discard) // hide flag errors
	flag.CommandLine.Init(os.Args[0], flag.ExitOnError)

	flag.BoolVarP(&o.help, "help", "h", false, "prints this help message")
	flag.CommandLine.StringVarP(&o.output, "output", "o", "", "``directory to save results (you can use this when processing a directory)")
	flag.CommandLine.StringVarP(&o.envFile, "env-file", "s", "", "``use a 'dotenv' file with the environment variables")
	flag.CommandLine.StringVarP(&o.inputFile, "file", "f", "", "``process this specific template file")
	flag.CommandLine.StringVarP(&o.inputDir, "dir", "d", "", "``process all templates in a directory (can have unlimited sub-folders)")
	flag.CommandLine.StringSliceVarP(&o.envPairs, "env", "e", []string{}, "``set a key-value environment variable pair")
	flag.CommandLine.StringSliceVarP(&o.ignore, "ignore", "i", []string{".git", ".DS_Store", ".env"}, "``specifies files patterns to ignore")

	flag.CommandLine.Parse(os.Args[1:])
}

type boilrOpts struct {
	envPairs  []string
	envFile   string
	output    string
	inputFile string
	inputDir  string
	ignore    []string
	help      bool
	version   string
	build     string
}

func (o *boilrOpts) complete(args []string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// load '.env' file in the current dir if exists
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	fn := filepath.Join(cwd, ".env")
	if pathutil.Exists(fn) {
		o.envFile = fn
	}
	// load the user specified '.env' file
	if len(o.envFile) > 0 {
		env, err := dotenv.FromFile(o.envFile)
		if err != nil {
			return err
		}
		dotenv.PutInEnv(env, true)
	}
	// override env with key=val pairs (if specified)
	if len(o.envPairs) > 0 {
		lines := strings.Join(o.envPairs, "\n")
		env, err := dotenv.FromReader(strings.NewReader(lines))
		if err != nil {
			return err
		}
		dotenv.PutInEnv(env, true)
	}

	if len(o.inputFile) > 0 {
		o.inputFile = pathutil.ExpandHome(o.inputFile, home)
		if !pathutil.Exists(o.inputFile) {
			return fmt.Errorf("file '%s' does not exists", o.inputFile)
		}
	}

	if len(o.inputDir) > 0 {
		o.inputDir = pathutil.ExpandHome(o.inputDir, home)
		if !pathutil.Exists(o.inputDir) {
			return fmt.Errorf("directory '%s' does not exists", o.inputDir)
		}
	}

	if len(o.output) > 0 {
		o.output = pathutil.ExpandHome(o.output, home)
		if strings.EqualFold(o.inputDir, o.output) {
			return fmt.Errorf("output directory must be different from input directory")
		}
	}

	return nil
}

func (o *boilrOpts) run(args []string) error {
	if len(o.inputFile) > 0 {
		return tpl.EvalFile(o.inputFile, o.output)
	}

	if len(o.inputDir) > 0 {
		fmt.Printf("Processing: %s\n", o.inputDir)

		return tpl.EvalDir(o.inputDir, o.output, o.ignore)
	}

	if len(args) == 0 {
		var err error
		fi, _ := os.Stdin.Stat()
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			res, err := tpl.EvalReader(bufio.NewReader(os.Stdin))
			if err == nil {
				fmt.Println(res)
			}
		}

		return err
	}

	res, err := tpl.EvalString(args[0])
	if err == nil {
		fmt.Println(res)
	}
	return err
}
