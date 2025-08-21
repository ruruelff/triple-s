package flagcheck

import (
	"flag"
	"fmt"
	"os"
	"strings"

	s "triple-s/source/structure"
)

func Help() {
	help := `Simple Storage Service.

	**Usage:**
		triple-s [-port <N>] [-dir <S>]  
		triple-s --help
	
	**Options:**
	- --help     Show this screen.
	- --port N   Port number
	- --dir S    Path to the directory`

	fmt.Println(help)
}

func init() {
	flag.Parse()

	if *s.HelpFlag {
		Help()
		os.Exit(0)
	}

	if strings.HasPrefix(*s.DirFlag, "/") || strings.Contains(*s.DirFlag, "..") {
		fmt.Println("Invalid directory path. It must be relative and not escape the base folder.")
		os.Exit(1)
	}

	if *s.DirFlag == "" {
		fmt.Println("Empty dir")
		os.Exit(1)
	}

	if *s.PortFlag == "" {
		fmt.Println("Empty port")
		os.Exit(1)
	}

	finalDir := "base/" + *s.DirFlag
	if _, err := os.Stat(finalDir); err != nil {
		os.MkdirAll(finalDir, os.ModePerm)
	}
}
