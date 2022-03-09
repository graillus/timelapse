package api

import (
	"io"
	"net/http"
	"os"
	"path"

	"github.com/alecthomas/units"
	"github.com/graillus/timelapse/internal/log"
)

const uploadsDir = "uploads"

var StoragePath = "/tmp"

func listFrames(resp http.ResponseWriter, _ *http.Request) {
	resp.WriteHeader(200)
}

func postFrame(resp http.ResponseWriter, req *http.Request) {
	uploadPath := path.Join(StoragePath, uploadsDir)

	err := req.ParseMultipartForm(10 * int64(units.MiB))
	if err != nil {
		http.Error(resp, err.Error(), 500)
		return
	}

	file, fileHeader, err := req.FormFile("frame")
	if err != nil {
		log.Warnf("Cannot retrieve the frame file: %s", err)
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Debugf("Uploaded File: %+v\n", fileHeader.Filename)
	log.Debugf("File Size: %+v\n", fileHeader.Size)
	log.Debugf("MIME Header: %+v\n", fileHeader.Header)

	// Create the uploads folder if it doesn't already exist
	err = os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		log.Errorf("Error creating upload directory: %s", err)
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory
	dest, err := os.Create(path.Join(uploadPath, fileHeader.Filename))
	if err != nil {
		log.Errorf("Error creating file: %s", err)
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dest.Close()

	// Copy the uploaded file to the filesystem at the specified destination
	_, err = io.Copy(dest, file)
	if err != nil {
		log.Errorf("Error copying file: %s", err)
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("Successfully uploaded frame %s", fileHeader.Filename)
	if err != nil {
		log.Errorf("Error writing response: %s", err)
		http.Error(resp, "Error writing response", http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusCreated)
}
