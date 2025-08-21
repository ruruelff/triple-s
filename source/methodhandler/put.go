package methodhandler

import (
	"io"
	"net/http"
	"os"

	h "triple-s/source/flagcheck"
	s "triple-s/source/structure"
	t "triple-s/source/tools"
)

func PutHandler(w http.ResponseWriter, r *http.Request) {
	t.CreateCSV()

	bucket := r.PathValue("BucketName")

	if bucket == "" {
		err := s.Error{Code: 409, Message: "Method not allowe", Resource: r.URL.Path}
		err.WriteError(w)
		return
	}

	if bucket == "metadata.csv" {
		err := s.Error{Code: 400, Message: "Error", Resource: r.URL.Path}
		err.WriteError(w)
		return
	}
	passed := h.IsValidBucketName(bucket)
	if !passed {
		err := s.Error{Code: 400, Message: "Bad Request for invalid names:" + bucket, Resource: r.URL.Path}
		err.WriteError(w)
		return
	}

	if err := t.ReadCSVfile(*s.DirFlag+"/metadata.csv", bucket); err != nil {
		err := s.Error{Code: 409, Message: "Conflict for duplicate names:" + bucket, Resource: r.URL.Path}
		err.WriteError(w)
		return
	}

	err := t.WriteCSVRecord(bucket, "in-active")
	if err != nil {
		http.Error(w, "Error writing to CSV file", http.StatusInternalServerError)
		return
	}

	if err := os.MkdirAll(*s.DirFlag+"/"+bucket, os.ModePerm); err != nil {
		http.Error(w, "Error creating directory for bucket", http.StatusInternalServerError)
		return
	}

	response := s.SuccessfulResponse{
		Code:   200,
		Status: "OK",
	}
	response.WriteResponse(w)
}

func PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	contentsize := r.ContentLength
	contenttype := r.Header.Get("Content-Type")
	bucket := r.PathValue("BucketName")
	object := r.PathValue("ObjectKey")
	t.CreateCSVinbucket(bucket)
	passed := h.IsValidBucketName(bucket)
	if !passed {
		err := s.Error{Code: 400, Message: "Bad Request for invalid bucket name: " + bucket, Resource: r.URL.Path}
		err.WriteError(w)
		return
	}

	passed = h.ValidateObjectName(object)
	if !passed {
		err := s.Error{Code: 400, Message: "Bad Request for invalid object name: " + object, Resource: r.URL.Path}
		err.WriteError(w)
		return
	}

	if object == "objectdata.csv" {
		err := s.Error{Code: 400, Message: "Invalid object name: " + object, Resource: r.URL.Path}
		err.WriteError(w)
		return
	}

	if exists, err := t.CheckBucketInMetadata(*s.DirFlag+"/metadata.csv", bucket); err != nil || !exists {

		err := s.Error{Code: 404, Message: "Bucket not found: " + bucket, Resource: r.URL.Path}
		err.WriteError(w)

		return
	}

	if err := t.UpdateOrWriteCSVRecord(bucket, object, contentsize, contenttype); err != nil {
		err := s.Error{Code: 500, Message: "Error updating/writing object metadata to CSV", Resource: r.URL.Path}
		err.WriteError(w)
		return
	}

	file, err := os.Create(*s.DirFlag + "/" + bucket + "/" + object)
	if err != nil {
		err := s.Error{Code: 500, Message: "Error creating object file", Resource: r.URL.Path}
		err.WriteError(w)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		err := s.Error{Code: 500, Message: "Error writing object data", Resource: r.URL.Path}
		err.WriteError(w)
		return
	}

	if err := t.UpdateCSV(*s.DirFlag+"/metadata.csv", bucket, "active"); err != nil {
		err := s.Error{Code: 500, Message: "Failed to update metadata", Resource: r.URL.Path}
		err.WriteError(w)
		return
	}

	response := s.SuccessfulResponse{
		Code:   200,
		Status: "OK",
	}
	response.WriteResponse(w)
}
