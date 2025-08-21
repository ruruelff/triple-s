package router

import (
	"net/http"

	m "triple-s/source/methodhandler"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("PUT /{BucketName}", m.PutHandler)
	mux.HandleFunc("GET /", m.GetHandler)
	mux.HandleFunc("DELETE /{BucketName}", m.DeleteHandler)
	mux.HandleFunc("PUT /{BucketName}/{ObjectKey}", m.PutObjectHandler)
	mux.HandleFunc("GET /{BucketName}/{ObjectKey}", m.GetObjectHandler)
	mux.HandleFunc("DELETE /{BucketName}/{ObjectKey}", m.DeleteObjectHandler)

	return mux
}
