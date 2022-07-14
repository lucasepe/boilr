# Boilr

**Projects templates using ${var} expansion syntax.**

> Something like [cookiecutter](https://github.com/cookiecutter/cookiecutter) but using environment variables expansion syntax.

- supports unlimited levels of directory nesting
- 100% of templating is done with environment variables expansion
- both, directory names and filenames can be templated; for example:
    ```sh
    ${DIR_NAME,,}/${FILE_NAME}.go
    ```
- simply define your template variables in a [_dotenv_](https://www.ibm.com/docs/en/aix/7.2?topic=files-env-file) file; for example:
    ```properties
    DIR_NAME=my-dir
    FILE_NAME=my-file
    ```
- by default the tool looks for a `.env` file in the current working folder
- specify your  `.env` file with the `--env-file` flag
- override or add environment variables using the `--env` flag
- eventually you can specify the output directory with the `--output` flag
- specifies files patterns to ignore with one or more `--ignore` flag(s)

## Supported Functions

| __Expression__                | __Meaning__                                                     |
| -----------------             | --------------                                                  |
| `${var}`                      | Value of `$var`
| `${#var}`                     | String length of `$var`
| `${var^}`                     | Uppercase first character of `$var`
| `${var^^}`                    | Uppercase all characters in `$var`
| `${var,}`                     | Lowercase first character of `$var`
| `${var,,}`                    | Lowercase all characters in `$var`
| `${var:n}`                    | Offset `$var` `n` characters from start
| `${var:n:len}`                | Offset `$var` `n` characters with max length of `len`
| `${var#pattern}`              | Strip shortest `pattern` match from start
| `${var##pattern}`             | Strip longest `pattern` match from start
| `${var*pattern}`              | Strip shortest `pattern` match from end
| `${var**pattern}`             | Strip longest `pattern` match from end
| `${var-default}`              | If `$var` is not set, evaluate expression as `$default`
| `${var:-default}`             | If `$var` is not set or is empty, evaluate expression as `$default`
| `${var=default}`              | If `$var` is not set, evaluate expression as `$default`
| `${var:=default}`             | If `$var` is not set or is empty, evaluate expression as `$default`
| `${var/pattern/replacement}`  | Replace as few `pattern` matches as possible with `replacement`
| `${var//pattern/replacement}` | Replace as many `pattern` matches as possible with `replacement`
| `${var/#pattern/replacement}` | Replace `pattern` match with `replacement` from `$var` start
| `${var/*pattern/replacement}` | Replace `pattern` match with `replacement` from `$var` end

## Install

You can install the pre-compiled binary or compile from source.

### Install the pre-compiled binary

Download the pre-compiled binaries from the [releases page](https://github.com/lucasepe/boilr/releases) and copy them to the desired location.

### Compile from source

If you have [Go installed](https://go.dev/doc/install), just open a terminal and type:

```sh
$ go install github.com/lucasepe/boilr@latest
```

---


### Boilr Templates

- [`Crossplane Provider`](https://github.com/lucasepe/boilr-templates/crossplane-provider) a minimal boilerplate for implementing new [Crossplane](https://crossplane.io) providers 
