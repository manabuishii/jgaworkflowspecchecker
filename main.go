package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
)

type simpleSchema struct {
	Name string `json:"name"`
	Md5  string `json:"md5,omitempty"`
	Fq1  string `json:"fq1"`
}

func md5File(filePath string) (string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(hash.Sum(nil)[:16]), nil
}

func main() {

	path, err := filepath.Abs("./")

	// MUST must be canonical
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + path + "/" + os.Args[1])
	// MUST must be canonical
	documentLoader := gojsonschema.NewReferenceLoader("file://" + path + "/" + os.Args[2])

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
	fmt.Println("Load sample")
	raw, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var ss simpleSchema

	json.Unmarshal(raw, &ss)
	fmt.Println("Load end")
	fmt.Printf("Name is [%s]\n", ss.Name)
	fmt.Printf("Fq1 is [%s]\n", ss.Fq1)
	fmt.Printf("Md5 is [%s]\n", ss.Md5)

	var fn string
	fn = ss.Fq1
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		fmt.Println("file doesn't exist")
	} else {
		fmt.Println("file exists")
	}
	if ss.Md5 != "" {
		md5, _ := md5File(fn)
		fmt.Printf("expected: [%s]\n", ss.Md5)
		fmt.Printf("actual  : [%s]\n", md5)
		if ss.Md5 == md5 {
			fmt.Println("md5 is match")
		} else {
			fmt.Println("md5 is not match")
		}
	} else {
		fmt.Println("md5 is not specified.")
	}
}
