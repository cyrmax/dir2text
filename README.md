# Dir2text

This small utility transforms any directory and file tree into a single text file.
This may be useful when feeding the entire project into an AI via chat or for other tasks.

## Features

* Written in Go (runs anywhere);
* Has no external dependencies and uses only standard library (another pro for "runs anywhere);
* Just does its job and nothing more;
* Does not include empty directories (and never will do so by design);

## Usage

The `-help` command line flag tells it all.

```
Usage of C:\Users\cyrmax\projects\dir2text\dir2text.exe:
  -dir string
        Specify custom working directory. Current working directory is used by default
  -output string
        Specify custom output file name. If empty, output file will have the same name as the working directory
```

## Building and installation

The utility is built like any other simple Go project:

```bash
cd dir2text
go build cmd/dir2text/main.go -o dir2text.exe
```

Installation is simple too: just put the executable somewhere into `PATH` and you are good to go.

## Contributing

Dir2text aims to be as much cross-platform as possible and as lightweight as possible, so avoid adding external dependencies and check your code on all possible platforms and architectures before submitting a pull request.
