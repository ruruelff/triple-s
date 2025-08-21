package methodhandler

import (
	"encoding/csv"
	"encoding/xml"
	"io"
	"net/http"
	"os"

	s "triple-s/source/structure"
	t "triple-s/source/tools"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	var buckets []s.Bucket

	if r.URL.Path != "/" {
		xmlErr := s.Error{Code: 400, Message: "Not correct Request", Resource: r.URL.Path}
		xmlErr.WriteError(w)
		return
	}
	f, err := os.Open(*s.DirFlag + "/" + "metadata.csv")
	if err != nil {
		xmlErr := s.Error{Code: 404, Message: "metadata.csv not found", Resource: r.URL.Path}
		xmlErr.WriteError(w)

		return
	}
	defer f.Close()

	reader := csv.NewReader(f)

	_, err = reader.Read()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if len(record) < 4 {
			xmlErr := s.Error{Code: 500, Message: "Invalid csv", Resource: r.URL.Path}
			xmlErr.WriteError(w)
			return

		}
		bucket := s.Bucket{
			Name:             record[0],
			CreationTime:     record[1],
			LastModifiedTime: record[2],
			Status:           record[3],
		}
		buckets = append(buckets, bucket)
	}

	bucketList := s.BucketList{
		Buckets: buckets,
	}

	w.Header().Set("Content-Type", "application/xml")

	err = xml.NewEncoder(w).Encode(bucketList)
	if err != nil {
		xmlErr := s.Error{Code: 500, Message: "Error encoding XML response", Resource: r.URL.Path}
		xmlErr.WriteError(w)
		return

	}
}

func GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("BucketName")
	object := r.PathValue("ObjectKey")
	if object == "objectdata.csv" {
		xmlErr := s.Error{Code: 400, Message: "You cannot get data about objects", Resource: r.URL.Path}
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
	file, err := os.Open(objectPath)
	if err != nil {
		xmlErr := s.Error{Code: 500, Message: "Failed to open object", Resource: r.URL.Path}
		xmlErr.WriteError(w)
		return
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	if err != nil {
		xmlErr := s.Error{Code: 500, Message: "Failed to send object data", Resource: r.URL.Path}
		xmlErr.WriteError(w)
		return
	}
}
