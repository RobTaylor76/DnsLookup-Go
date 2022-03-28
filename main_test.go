package main

import (
	"os"
	"testing"
)

func Test_request(t *testing.T) {

}

func Test_A_response(t *testing.T) {

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

func Test_MX_response(t *testing.T) {

	file, err := os.Open("testdata/mx_response.bin")
	if err != err {
		t.Fatal(err.Error())
	}

	buffer := make([]byte, 512)

	file.Read(buffer)

	// fmt.Print(buffer)

	response := processDNSResponse(buffer)


	if len(response.Answers) != 5 {
		t.Fatalf("Not enough answers")
	}

	// fullstackconsultancy.com: type MX, class IN, preference 10, mx alt3.aspmx.l.google.com

	if response.Answers[0].DomainName != "fullstackconsultancy.com" ||
		response.Answers[0].Answer != "alt3.aspmx.l.google.com" ||
		response.Answers[0].Preference != 10 {
		t.Fatalf("Expected different answers alt3.aspmx.l.google.com 10")
	}

	// fullstackconsultancy.com: type MX, class IN, preference 1, mx aspmx.l.google.com

	if response.Answers[1].DomainName != "fullstackconsultancy.com" ||
		response.Answers[1].Answer != "aspmx.l.google.com" ||
		response.Answers[1].Preference != 1 {
		t.Fatalf("Expected different answers aspmx.l.google.com 1")
	}

	// fullstackconsultancy.com: type MX, class IN, preference 5, mx alt1.aspmx.l.google.com
	if response.Answers[2].DomainName != "fullstackconsultancy.com" ||
		response.Answers[2].Answer != "alt1.aspmx.l.google.com" ||
		response.Answers[3].Preference != 5 {
		t.Fatalf("Expected different answers alt1.aspmx.l.google.com 5")
	}
//fullstackconsultancy.com: type MX, class IN, preference 5, mx alt2.aspmx.l.google.com
	if response.Answers[3].DomainName != "fullstackconsultancy.com" ||
		response.Answers[3].Answer != "alt2.aspmx.l.google.com" ||
		response.Answers[3].Preference != 5 {
		t.Fatalf("Expected different answers alt2.aspmx.l.google.com 10")
	}
// fullstackconsultancy.com: type MX, class IN, preference 10, mx alt4.aspmx.l.google.com
	if response.Answers[4].DomainName != "fullstackconsultancy.com" ||
		response.Answers[4].Answer != "alt4.aspmx.l.google.com" ||
		response.Answers[4].Preference != 10 {
		t.Fatalf("Expected different answers alt4.aspmx.l.google.com 10")
	}
}