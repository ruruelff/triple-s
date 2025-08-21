package tools

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	s "triple-s/source/structure"
)

func CreateCSV() {
	if _, err := os.Stat(*s.DirFlag + "/" + "metadata.csv"); err != nil {
		file, err := os.Create(*s.DirFlag + "/" + "metadata.csv")
		if err != nil {
			return
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		headers := []string{"Name", "CreationTime", "LastModifiedTime", "Status"}
		err = writer.Write(headers)
		if err != nil {
			return
		}
	}
}

func CreateCSVinbucket(bucket string) {
	if _, err := os.Stat(*s.DirFlag + "/" + bucket + "/" + "objectdata.csv"); err != nil {
		file, err := os.Create(*s.DirFlag + "/" + bucket + "/" + "objectdata.csv")
		if err != nil {
			return
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		headers := []string{"ObjectKey:", "Size:", "ContentType:", "LastModified:"}
		err = writer.Write(headers)
		if err != nil {
			return
		}
	}
}

func RemoveBucketFromCsv(filename string, bucketname string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	var updatedrecords [][]string

	for _, record := range records {
		if record[0] != bucketname {
			updatedrecords = append(updatedrecords, record)
		}
	}

	f, err = os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return nil
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	for _, record := range updatedrecords {
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func RemoveObjectFromCsv(filename string, object string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	var updatedrecords [][]string

	for _, record := range records {
		if record[0] != object {
			updatedrecords = append(updatedrecords, record)
		}
	}

	fw, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer fw.Close()

	writer := csv.NewWriter(fw)
	defer writer.Flush()

	for _, record := range updatedrecords {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func ReadCSVfile(filename string, bucketname string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	for _, line := range records {
		if line[0] == bucketname {
			return errors.New("bucket exist")
		}
	}

	return nil
}

func ReadObject(filename string, objectname string, content int64, contenttype string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	for _, line := range records {
		if line[0] == objectname {
			line[0] = objectname
			line[1] = strconv.FormatInt(content, 10)
			line[2] = contenttype
			time.Now().Format(time.RFC3339)
		}
	}
	return nil
}

func UpdateCSV(filename string, bucketName string, status string) error {
	file, err := os.OpenFile(filename, os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	updated := false
	for i, record := range records {
		if record[0] == bucketName {
			records[i][2] = time.Now().Format(time.RFC3339)
			records[i][3] = status
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("bucket %s not found", bucketName)
	}

	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(records)
	if err != nil {
		return err
	}

	return nil
}

func UpdateOrWriteCSVRecord(BucketName, ObjectKey string, content int64, contenttype string) error {
	filePath := *s.DirFlag + "/" + BucketName + "/objectdata.csv"

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	var records [][]string
	found := false

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if record[0] == ObjectKey {
			record[1] = strconv.FormatInt(content, 10)
			record[2] = contenttype
			record[3] = time.Now().Format(time.RFC3339)
			found = true
		}
		records = append(records, record)
	}

	if !found {
		record := []string{
			ObjectKey,
			strconv.FormatInt(content, 10),
			contenttype,
			time.Now().Format(time.RFC3339),
		}
		records = append(records, record)
	}

	f.Truncate(0)
	f.Seek(0, io.SeekStart)
	writer := csv.NewWriter(f)
	err = writer.WriteAll(records)
	if err != nil {
		return err
	}

	return nil
}

func CheckBucketInMetadata(filename, bucketName string) (bool, error) {
	f, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	_, _ = reader.Read()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, err
		}
		if record[0] == bucketName {
			return true, nil
		}
	}

	return false, nil
}

func CheckObjectInCSV(filename, objectName string) (bool, error) {
	f, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	_, _ = reader.Read()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, err
		}
		if record[0] == objectName {
			return true, nil
		}
	}

	return false, nil
}

func WriteCSVRecord(bucketName, status string) error {
	f, err := os.OpenFile(*s.DirFlag+"/metadata.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	record := []string{
		bucketName,
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339),
		status,
	}

	err = writer.Write(record)
	if err != nil {
		return err
	}

	return nil
}

func WriteCSVRecordinBucket(BucketName string, ObjectKey string, content int64, contenttype string) error {
	f, err := os.OpenFile(*s.DirFlag+"/"+BucketName+"/"+"objectdata.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	record := []string{
		ObjectKey,
		strconv.FormatInt(content, 10),
		contenttype,
		time.Now().Format(time.RFC3339),
	}
	err = writer.Write(record)
	if err != nil {
		return err
	}

	return nil
}

func GetBucketStatus(filePath, bucketName string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	for _, record := range records[1:] {
		if record[0] == bucketName {
			return record[3], nil
		}
	}
	return "", nil
}

func HasObjectsInBucket(bucket string) (bool, error) {
	filePath := *s.DirFlag + "/" + bucket + "/objectdata.csv"
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read()
	if err != nil {
		return false, err
	}

	_, err = reader.Read()
	if err == io.EOF {
		return false, nil
	}
	return true, nil
}

func InitSSS(baseDIR string) error {
	metadataPath := baseDIR + "/metadata.csv"

	// Create base directory if it doesn't exist
	if _, err := os.Stat(baseDIR); os.IsNotExist(err) {
		if err := os.MkdirAll(baseDIR, os.ModePerm); err != nil {
			return err
		}
	}

	// Create metadata.csv if it doesn't exist
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		file, err := os.Create(metadataPath)
		if err != nil {
			return err
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()
		writer.Write([]string{"Name", "CreationTime", "LastModifiedTime", "Status"})
	}

	return nil
}

func VerifyBuckets(csvPath string, baseDir string) error {
	file, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil || len(headers) != 4 {
		return fmt.Errorf("invalid or missing CSV header")
	}

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		if len(record) != 4 || record[0] == "" {
			return fmt.Errorf("invalid bucket record: %+v", record)
		}

		bucketPath := filepath.Join(baseDir, record[0])
		if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
			return fmt.Errorf("bucket directory missing: %s", record[0])
		}
	}

	return nil
}

func CheckBucket(bucketName, baseDir string) error {
	objectPath := filepath.Join(baseDir, bucketName, "objectdata.csv")
	file, err := os.Open(objectPath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil || len(header) != 4 {
		return fmt.Errorf("invalid objectdata.csv header in %s", bucketName)
	}

	return nil
}

func VerifyAll(baseDir string) error {
	csvPath := filepath.Join(baseDir, "metadata.csv")

	if err := VerifyBuckets(csvPath, baseDir); err != nil {
		return fmt.Errorf("bucket validation failed: %v", err)
	}

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if err := CheckBucket(entry.Name(), baseDir); err != nil {
				return fmt.Errorf("bucket integrity error in %s: %v", entry.Name(), err)
			}
		}
	}

	return nil
}
