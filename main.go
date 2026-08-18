package hopper

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	srcPath := flag.String("src", "", "Source Grammar File")
	destFolder := flag.String("dest", "", "Destination Folder")
	flag.Parse()
	remaining := flag.Args()
	if len(remaining) > 0 && *srcPath == "" {
		*srcPath = remaining[0]
		remaining = remaining[1:]
	}
	if len(remaining) > 0 && *destFolder == "" {
		*destFolder = remaining[0]
		remaining = remaining[1:]
	}
	if len(remaining) > 0 {
		fmt.Printf("Ignoring %d unused arguments.\n", len(remaining))
	}
	if len(*srcPath) == 0 {
		fmt.Println("Expected a path to a grammar file")
	} else if len(*destFolder) == 0 {
		fmt.Println("Expected a path to a destination folder")
	} else {
		if fs, err := os.Stat(*destFolder); err != nil {
			fmt.Printf("Path \"%s\" does not exist.\n", *destFolder)
		} else if !fs.IsDir() {
			fmt.Printf("Path \"%s\" does not point to a directory.\n", *destFolder)
		} else {
			p := MakeParser()
			statements, err := p.ParseFile(*srcPath)
			if err != nil {
				fmt.Println(err)
			} else {
				Execute(*statements, *destFolder)
			}
		}
	}
	fmt.Printf("Press [ENTER] to exit")
	fmt.Scanf(".")
}
