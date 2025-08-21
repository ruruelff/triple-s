package main

import (
	"fmt"
	"log"
	"net/http"

	r "triple-s/source/router"
	s "triple-s/source/structure"
	t "triple-s/source/tools"
)

const RootBaseDir = "base"

func main() {
	*s.DirFlag = RootBaseDir + "/" + *s.DirFlag
	fmt.Println("Creating server on port", *s.PortFlag)
	fmt.Println("http://localhost:" + *s.PortFlag)

	if err := t.InitSSS(*s.DirFlag); err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	mux := r.Router()
	log.Fatal(http.ListenAndServe(":"+*s.PortFlag, mux))
}
