package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Terry-Mao/goconf"
)

type EFSConfig struct {
	Port         int    `goconf:"server:port"`
	Storage_path string `goconf:"server:storage_path"`
}

type EFSUpload struct {
	EncryptedFile string
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	tmpfile, err := ioutil.TempFile(getConfig().Storage_path, "enc_")
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

	tmpfile.WriteString(efsupload.EncryptedFile)

	io.WriteString(w, tmpfile.Name())
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

	// Read the config
	config := getConfig()

	fmt.Printf("Starting Encrypted File Sharing server at port %d\n", config.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil); err != nil {
		log.Fatal(err)
	}

	// Add uploadHandler

}
