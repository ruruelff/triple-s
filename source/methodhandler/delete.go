package methodhandler

import (
	"net/http"
	"os"

	s "triple-s/source/structure"
	t "triple-s/source/tools"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("BucketName")
	if bucket == "" {
		xmlErr := s.Error{Code: 405, Message: "Method Not Allowed", Resource: r.URL.Path}
		xmlErr.WriteError(w)
		return
	}
	if bucket == "metadata.csv" {
		xmlErr := s.Error{Code: 400, Message: "Error", Resource: r.URL.Path}
		xmlErr.WriteError(w)
		return
	}

	status, err := t.GetBucketStatus(*s.DirFlag+"/metadata.csv", bucket)

	if status == "" {
		xmlErr := s.Error{Code: 404, Message: "Bucket not found", Resource: r.URL.Path}
		xmlErr.WriteError(w)
		return
	}
	if status != "in-active" {
		xmlErr := s.Error{Code: 409, Message: "Bucket is active, deletion not allowed", Resource: r.URL.Path}
		xmlErr.WriteError(w)
		return
	}

	err = t.RemoveBucketFromCsv(*s.DirFlag+"/metadata.csv", bucket)
	if err != nil {
		xmlErr := s.Error{Code: 400, Message: "Error removing bucket from CSV", Resource: r.URL.Path}
		xmlErr.WriteError(w)
		return
	}

	err = os.RemoveAll(*s.DirFlag + "/" + bucket)

	response := s.SuccessfulResponse{
		Code:   204,
		Status: "No Content",
	}
	response.WriteResponse(w)
}

func DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("BucketName")
	object := r.PathValue("ObjectKey")
	if object == "objectdata.csv" {
		xmlErr := s.Error{Code: 400, Message: "Error", Resource: r.URL.Path}
		xmlErr.WriteError(w)
		return
	}

	if exists, err := t.CheckBucketInMetadata(*s.DirFlag+"/metadata.csv", bucket); err != nil || !exists {

		xmlErr := s.Error{Code: 404, Message: "Bucket not found", Resource: r.URL.Path}
		xmlErr.WriteError(w)

		return
	}

	if exists, err := t.CheckObjectInCSV(*s.DirFlag+"/"+bucket+"/objectdata.csv", object); err != nil || !exists {

		xmlErr := s.Error{Code: 404, Message: "Object not found", Resource: r.URL.Path}
		xmlErr.WriteError(w)

		return
	}

	objectPath := *s.DirFlag + "/" + bucket + "/" + object
	err := os.Remove(objectPath)

	err = t.RemoveObjectFromCsv(*s.DirFlag+"/"+bucket+"/objectdata.csv", object)
	if err != nil {
		xmlErr := s.Error{Code: 500, Message: "Error updating object metadata", Resource: r.URL.Path}
		xmlErr.WriteError(w)
		return
	}
	t.UpdateCSV(*s.DirFlag+"/"+"metadata.csv", bucket, "active")
	asd, _ := t.HasObjectsInBucket(bucket)
	if !asd {
		t.UpdateCSV(*s.DirFlag+"/"+"metadata.csv", bucket, "in-active")
	}
	response := s.SuccessfulResponse{
		Code:   204,
		Status: "No Content",
	}
	response.WriteResponse(w)
}
