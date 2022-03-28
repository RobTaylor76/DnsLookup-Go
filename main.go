package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

var nameToLookup = flag.String("domainName", "www.google.com", "specify the name fof the host to lookup")
var dnsServer = flag.String("dnsIP", "192.168.1.1", "specify the IP of DNS Server")

func main() {
	flag.Parse()
	dnsResponse := contactDNS(*nameToLookup, *dnsServer)

	for _, dnsAnswer := range dnsResponse.Answers {
		fmt.Printf("%s %s %d %d", dnsAnswer.DomainName, dnsAnswer.Answer, dnsAnswer.AnswerType, dnsAnswer.AnswerClass)
		fmt.Println()
	}

}

func bufferToHeader(buffer []byte) dnsHeader {
	header := dnsHeader{}

	header.identifier = binary.BigEndian.Uint16(buffer[:2])
	header.flags = binary.BigEndian.Uint16(buffer[2:4])
	header.questionCount = binary.BigEndian.Uint16(buffer[4:6])
	header.answerCount = binary.BigEndian.Uint16(buffer[6:8])
	header.nsCount = binary.BigEndian.Uint16(buffer[8:10])
	header.additionalRecordsCount = binary.BigEndian.Uint16(buffer[10:12])
	return header
}

type dnsHeader struct {
	identifier             uint16
	flags                  uint16
	questionCount          uint16
	answerCount            uint16
	nsCount                uint16
	additionalRecordsCount uint16
}

//type dnsResponse [512]byte

func headerToBuffer(header dnsHeader) [512]byte {
	out := [512]byte{}
	binary.BigEndian.PutUint16(out[:2], header.identifier)
	binary.BigEndian.PutUint16(out[2:4], header.flags)
	binary.BigEndian.PutUint16(out[4:6], header.questionCount)
	binary.BigEndian.PutUint16(out[6:8], header.answerCount)
	binary.BigEndian.PutUint16(out[8:10], header.nsCount)
	binary.BigEndian.PutUint16(out[10:12], header.additionalRecordsCount)

	return out
}

func contactDNS(query string, dnsServerIP string) DNSResponse {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	header := dnsHeader{
		identifier:             uint16(r1.Intn(64000)),
		flags:                  0x0120,
		questionCount:          1,
		answerCount:            0,
		nsCount:                0,
		additionalRecordsCount: 0,
	}

	fmt.Println(header.identifier)

	buffer := headerToBuffer(header)

	bufferOffset := 12
	for _, label := range strings.Split(query, ".") {
		buffer[bufferOffset] = byte(len(label))
		bufferOffset += 1
		for _, char := range label {
			buffer[bufferOffset] = byte(char)
			bufferOffset += 1
		}
	}

	buffer[bufferOffset] = byte(0)
	bufferOffset++
	binary.BigEndian.PutUint16(buffer[bufferOffset:bufferOffset+2], 1)
	bufferOffset += 2
	binary.BigEndian.PutUint16(buffer[bufferOffset:bufferOffset+2], 1)
	bufferOffset += 2

	//header2 := bufferToHeader(buffer)

	//err := ioutil.WriteFile("./testdata/sample_request.bin", buffer[:], 0644)
	//if err != nil {
	//	fmt.Errorf(err.Error())
	//	return
	//}

	response := make([]byte, 512)

	client, err := net.Dial("udp", dnsServerIP + ":53")
	if err != nil {
		fmt.Errorf(err.Error())
	}
	defer func() { _ = client.Close() }()

	client.Write(buffer[:])

	client.Read(response)

	//err = ioutil.WriteFile("./testdata/sample_response.bin", response[:], 0644)
	//if err != nil {google.com
	//	fmt.Errorf(err.Error())
	//	return
	//}

	return processDNSResponse(response)

}

func processDNSResponse(buffer []byte) DNSResponse {
	header2 := bufferToHeader(buffer)

	questions := make([]Question, 0, 5)
	questionsBuffer := buffer[12:]

	offset := uint16(0)
	for question := uint16(1); question <= header2.questionCount; question++ {
		domainName, offsetBy := ExtractDomainName(questionsBuffer[offset:], buffer)
		offset += offsetBy
		questionType := binary.BigEndian.Uint16(questionsBuffer[offset : offset+2])
		offset += 2
		questionClass := binary.BigEndian.Uint16(questionsBuffer[offset : offset+2])
		offset += 2

		questions = append(questions,
			Question{
				DomainName:    string(domainName),
				QuestionClass: questionClass,
				QuestionType:  questionType,
			})
	}
	// questionsBuffer 1st

	// then answersBuffer
	answers := make([]Answer, 0, 15)
	answersBuffer := questionsBuffer[offset:]

	answersOffset := uint16(0)
	for answer := uint16(1); answer <= header2.answerCount; answer++ {
		isPointer := answersBuffer[answersOffset] > 63

		var domainName string
		var offsetBy uint16

		if isPointer {
			domainNameOffset := binary.BigEndian.Uint16(answersBuffer[answersOffset:(answersOffset+2)]) ^ (192 << 8)
			domainName, _ = ExtractDomainName(buffer[domainNameOffset:], buffer)
			answersOffset += 2
		} else {
			domainName, offsetBy = ExtractDomainName(answersBuffer[answersOffset:], buffer)
			answersOffset += offsetBy
		}

		answerType := binary.BigEndian.Uint16(answersBuffer[answersOffset : answersOffset+2])
		answersOffset += 2
		answerClass := binary.BigEndian.Uint16(answersBuffer[answersOffset : answersOffset+2])
		answersOffset += 2
		ttl := binary.BigEndian.Uint32(answersBuffer[answersOffset : answersOffset+4])
		answersOffset += 4

		rdLength := binary.BigEndian.Uint16(answersBuffer[answersOffset : answersOffset+2])
		answersOffset += 2

		var rdData string
		if answerType == 1 {
			rdData = fmt.Sprintf("%d.%d.%d.%d",
				answersBuffer[answersOffset],
				answersBuffer[answersOffset+1],
				answersBuffer[answersOffset+2],
				answersBuffer[answersOffset+3])
		} else {
			rdData, _ = ExtractDomainName(answersBuffer[answersOffset:], buffer)
		}

		answersOffset += rdLength

		answers = append(answers,
			Answer{DomainName: domainName,
				AnswerClass: answerClass,
				AnswerType:  answerType,
				TTL:         ttl,
				Answer:      rdData,
			})

	}
	return DNSResponse{Questions: questions, Answers: answers}
}

func ExtractDomainName(primaryBuffer []byte, messageBuffer []byte) (string, uint16) {
	offset := uint16(0)
	domainName := make([]byte, 0, 256)
	noOfChars := primaryBuffer[offset]
	offset++
	for {
		isPointer := noOfChars > 63

		if isPointer {
			domainNameOffset := binary.BigEndian.Uint16(primaryBuffer[offset-1:(offset+1)]) ^ (192 << 8)
			partialDomainName, _ := ExtractDomainName(messageBuffer[domainNameOffset:], messageBuffer)
			return string(domainName) + partialDomainName, offset
		}

		for nextChar := byte(0); nextChar < noOfChars; nextChar++ {
			domainName = append(domainName, primaryBuffer[offset])
			offset++
		}
		noOfChars = primaryBuffer[offset]
		offset++
		if noOfChars == 0 {
			break
		} else {
			domainName = append(domainName, '.')
		}
	}

	return string(domainName), offset
}

type DNSResponse struct {
	Questions []Question
	Answers   []Answer
}

type Question struct {
	DomainName    string
	QuestionType  uint16
	QuestionClass uint16
}

type Answer struct {
	DomainName  string
	AnswerType  uint16
	AnswerClass uint16
	TTL         uint32
	Answer      string
}
