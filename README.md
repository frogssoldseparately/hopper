# hopper - a simple parser generator for simple languages

Give a grammar file, get a parser written in go.

Usage:
```sh
go install github.com/frogssoldseparately/hopper
hopper sourceGrammarFile.txt path/to/destination
```

This will generate a parser written in go that is ready for use in a project. You provide a grammar file (as described in `example.txt`), a destination folder, and your auto generated code will be written. The code will be associated with a package as according to the name of the destination directory. For example:

```sh
hopper source.txt pkg/simpleparser
```

will generate files with a package name of `simpleparser` in the directory `pkg/simpleparser`. A file will get written for each parsing layer (Token, Expression, Statement, etc.) whose names will be `[lowercaseparsegroup]s.go`, and an additional file called `parser.go`. Be wary of modifying these files, as any changes will be reverted the next time hopper is run for that target directory. **It will not ask you before overwriting a file.**