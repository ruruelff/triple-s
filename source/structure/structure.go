package structure

import (
	"encoding/xml"
	"flag"
	"net/http"
)

var (
	HelpFlag = flag.Bool("help", false, "show help message")
	PortFlag = flag.String("port", "8080", "port number")
	DirFlag  = flag.String("dir", "data", "path to the directory")
)

type Bucket struct {
	Name             string `xml:"Name"`
	CreationTime     string `xml:"CreationDate"`
	LastModifiedTime string `xml:"LastModified"`
	Status           string `xml:"Status"`
}

type BucketList struct {
	Buckets []Bucket `xml:"Buckets>Bucket"`
}

type Error struct {
	Code     int    `xml:"Code"`
	Message  string `xml:"Message"`
	Resource string `xml:"Resource"`
}

type SuccessfulResponse struct {
	Code   int    `xml:"Code"`
	Status string `xml:"Status"`
}

func (e *Error) WriteError(w http.ResponseWriter) {
	x, err := xml.MarshalIndent(e, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(e.Code)
	w.Write(x)
}

func (s *SuccessfulResponse) WriteResponse(w http.ResponseWriter) {
	x, err := xml.MarshalIndent(s, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(s.Code)
	w.Write(x)
}
