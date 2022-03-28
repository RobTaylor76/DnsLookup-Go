package main

import (
	"os"
	"testing"
)

func Test_request(t *testing.T) {

}

func Test_response(t *testing.T) {

	file, err := os.Open("testdata/bbc_response.bin")
	if err != err {
		t.Fatal(err.Error())
	}

	buffer := make([]byte, 512)

	file.Read(buffer)

	// fmt.Print(buffer)

	response := processDNSResponse(buffer)


	if len(response.Answers) != 4 {
		t.Fatalf("Not enough answers")
	}

	if response.Answers[0].DomainName != "www.bbc.co.uk" || response.Answers[0].Answer != "www.bbc.co.uk.pri.bbc.co.uk" {
		t.Fatalf("Expected different answers www.bbc.co.uk.pri.bbc.co.uk")
	}

	if response.Answers[1].DomainName != "www.bbc.co.uk.pri.bbc.co.uk" || response.Answers[1].Answer != "uk.www.bbc.co.uk.pri.bbc.co.uk" {
		t.Fatalf("Expected different answers uk.www.bbc.co.uk.pri.bbc.co.uk")
	}

	if response.Answers[2].DomainName != "uk.www.bbc.co.uk.pri.bbc.co.uk" || response.Answers[2].Answer != "212.58.233.253" {
		t.Fatalf("Expected different answers 212.58.233.253")
	}

	if response.Answers[3].DomainName != "uk.www.bbc.co.uk.pri.bbc.co.uk" || response.Answers[3].Answer != "212.58.237.253" {
		t.Fatalf("Expected different answers 212.58.237.253")
	}
}
