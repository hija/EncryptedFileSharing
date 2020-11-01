package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Terry-Mao/goconf"
)

type EFSConfig struct {
	Port          int    `goconf:"server:port"`
	Storage_path  string `goconf:"server:storage_path"`
	ServeFrontend bool   `goconf:"server:serve_frontend"`
	FrontendPath  string `goconf:"server:frontend_path"`
}

type EFSUpload struct {
	// This is what we get from the client as JSON
	EncryptedFile     string `json:"encryptedFile"`     // The encrypted file itself
	EncryptedFileName string `json:"encryptedFileName"` // The encrypted file name
}

type EFSConfiguration struct {
	// This is how we save on the server to save some space
	EncryptedFileName string // The encrypted file name
	EncryptedFilePath string // The encrypted file itself, stored seperately
}

type EFSDownload struct {
	// This is how the client asks for a download as JSON
	// Additionally this is our response to an upload request
	EFSConfigurationName string // The efsconfiguration name
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	efsdownload := EFSDownload{}
	json.NewDecoder(r.Body).Decode(&efsdownload)

	efsConfigurationBytes, err := ioutil.ReadFile(filepath.Join(getConfig().Storage_path, efsdownload.EFSConfigurationName+".json"))
	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}
	efsConfiguration := EFSConfiguration{}
	json.Unmarshal(efsConfigurationBytes, &efsConfiguration)

	// Next we need to read in the file again and encode it to Base64 again
	binFileBytes, err := ioutil.ReadFile(filepath.Join(getConfig().Storage_path, efsdownload.EFSConfigurationName))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	base64EncodedFile := base64.StdEncoding.EncodeToString(binFileBytes)

	// Now we create the efsupload again
	efsupload := EFSUpload{}
	efsupload.EncryptedFile = base64EncodedFile
	efsupload.EncryptedFileName = efsConfiguration.EncryptedFileName

	// Send it back to the client
	efsuploadBytes, err := json.Marshal(efsupload)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Write(efsuploadBytes)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	tmpfile, err := ioutil.TempFile(getConfig().Storage_path, "*")
	if err != nil {
		log.Fatal(err)
	}

	defer tmpfile.Close()

	efsupload := EFSUpload{}

	err = json.NewDecoder(r.Body).Decode(&efsupload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// We decode the base64 file, so we can save some disk space
	base64DecodedFile, _ := base64.StdEncoding.DecodeString(efsupload.EncryptedFile)
	tmpfile.Write(base64DecodedFile)
	// io.WriteString(w, tmpfile.Name())

	// Create overview file
	efsconfiguration := EFSConfiguration{}
	efsconfiguration.EncryptedFileName = efsupload.EncryptedFileName
	efsconfiguration.EncryptedFilePath = tmpfile.Name()

	// Write configuration to jsonfile
	jsonFileName := tmpfile.Name() + ".json"
	jsonEFSConfiguration, err := json.Marshal(efsconfiguration)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(jsonFileName, jsonEFSConfiguration, 0644)

	// Write response to client
	efsdownload := EFSDownload{}
	efsdownload.EFSConfigurationName = filepath.Base(tmpfile.Name())
	efsdownloadBytes, err := json.Marshal(efsdownload)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(efsdownloadBytes)
}

func getConfig() *EFSConfig {

	conf := goconf.New()
	if err := conf.Parse("config.conf"); err != nil {
		log.Fatal(err)
	}

	efsconfig := &EFSConfig{}
	if err := conf.Unmarshal(efsconfig); err != nil {
		panic(err)
	}
	return efsconfig
}

func main() {

	// Register Handlers
	http.HandleFunc("/api/upload", uploadHandler)
	http.HandleFunc("/api/download", downloadHandler)

	// Read the config
	config := getConfig()

	// If we shall serve frontend we do it :)
	if config.ServeFrontend {
		fs := http.FileServer(http.Dir(config.FrontendPath))
		http.Handle("/", fs)
	}

	fmt.Printf("Starting Encrypted File Sharing server at port %d\n", config.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil); err != nil {
		log.Fatal(err)
	}

	// Add uploadHandler

}
