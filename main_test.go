package main

import (
	"os"
	"testing"
)

func Test_request(t *testing.T) {

}

func Test_response(t *testing.T) {

	file, err := os.Open("testdata/sample_response.bin")
	if err != err {
		t.Fatal(err.Error())
	}

	buffer := make([]byte, 512)

	file.Read(buffer)

	// fmt.Print(buffer)

	response := processDNSResponse(buffer)


	if len(response.Answers) != 4 {
		t.Fatalf("Not enough questions")
	}
}
