package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type vuiRequestStruct struct {
	XMLName xml.Name      `xml:"VUI"`
	VER     string        `xml:"ver,attr,omitempty"`
	HDR     hdrStruct     `xml:"HDR,omitempty"`
	Payload payloadStruct `xml:"Payload,omitempty"`
	TrlRec  string        `xml:"TRL>REC,omitempty"`
}

type vuiResponseStruct struct {
	XMLName xml.Name      `xml:"VUI"`
	VER     string        `xml:"ver,attr,omitempty"`
	HDR     hdrStruct     `xml:"HDR,omitempty"`
	Payload payloadStruct `xml:"Payload,omitempty"`
	TrlRec  string        `xml:"TRL>REC,omitempty"`
}

type hdrStruct struct {
	XMLName       xml.Name `xml:"HDR"`
	Acct          string   `xml:"Acct,omitempty"`
	ClientVersion string   `xml:"ClientVersion,omitempty"`
	REC           string   `xml:"REC,omitempty"`
}

type rc1Struct struct {
	XMLName xml.Name `xml:"RC1"`
	RC1     string   `xml:",chardata"`
	Message string   `xml:"message,attr,omitempty"`
}

type aliQueryRequestStruct struct {
	XMLName     xml.Name `xml:"ALIQueryRequest"`
	VER         string   `xml:"ver,attr,omitempty"`
	ExternalKey string   `xml:"ExternalKey,omitempty"`
}

type aliQueryResponseStruct struct {
	XMLName     xml.Name  `xml:"ALIQueryResponse"`
	VER         string    `xml:"ver,attr,omitempty"`
	RC1         rc1Struct `xml:"RC1,omitempty"`
	ExternalKey string    `xml:"ExternalKey,omitempty"`
	HNO         string    `xml:"HNO,omitempty"`
	STN         string    `xml:"STN,omitempty"`
	MCN         string    `xml:"MCN,omitempty"`
	STA         string    `xml:"STA,omitempty"`
	LOC         string    `xml:"LOC,omitempty"`
	NAM         string    `xml:"NAM,omitempty"`
	ClsType     string    `xml:"CLS>TYP,omitempty"`
	TysType     string    `xml:"TYS>TYP,omitempty"`
	ESN         string    `xml:"ESN,omitempty"`
	CPD         string    `xml:"CPD,omitempty"`
	COI         string    `xml:"COI,omitempty"`
	CPF         string    `xml:"CPF,omitempty"`
	ZIP         string    `xml:"ZIP,omitempty"`
	ALT         string    `xml:"ALT,omitempty"`
}

type aliUpdateRequestStruct struct {
	XMLName         xml.Name `xml:"ALIUpdateRequest"`
	FOC             string   `xml:"FOC,attr,omitempty"`
	VER             string   `xml:"ver,attr,omitempty"`
	ExternalKey     string   `xml:"ExternalKey,omitempty"`
	ExternalKeyType string   `xml:"ExternalKeyType,omitempty"`
	HNO             string   `xml:"HNO,omitempty"`
	STN             string   `xml:"STN,omitempty"`
	MCN             string   `xml:"MCN,omitempty"`
	STA             string   `xml:"STA,omitempty"`
	LOC             string   `xml:"LOC,omitempty"`
	NAM             string   `xml:"NAM,omitempty"`
	ClsType         string   `xml:"CLS>TYP,omitempty"`
	TysType         string   `xml:"TYS>TYP,omitempty"`
	COI             string   `xml:"COI,omitempty"`
	CPF             string   `xml:"CPF,omitempty"`
	ZIP             string   `xml:"ZIP,omitempty"`
	SubscriberID    string   `xml:"SubscriberID,omitempty"`
}

type aliUpdateResponseStruct struct {
	XMLName         xml.Name  `xml:"ALIUpdateResponse"`
	VER             string    `xml:"ver,attr,omitempty"`
	ExternalKey     string    `xml:"ExternalKey,omitempty"`
	ExternalKeyType string    `xml:"ExternalKeyType,omitempty"`
	RC1             rc1Struct `xml:"RC1,omitempty"`
}

type payloadStruct struct {
	XMLName           xml.Name                `xml:"Payload"`
	AliQueryRequest   aliQueryRequestStruct   `xml:"ALIQueryRequest"`
	AliQueryResponse  aliQueryResponseStruct  `xml:"ALIQueryResponse"`
	AliUpdateRequest  aliUpdateRequestStruct  `xml:"ALIUpdateRequest,"`
	AliUpdateResponse aliUpdateResponseStruct `xml:"ALIUpdateResponse"`
}

// VUI handler
func vuiHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("VUI Request - Got connection: ", r.Proto)

	// Check Request Method
	log.Debug("VUI Request - Request Method: ", r.Method)
	if strings.ToUpper(r.Method) != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Debug("VUI Request - ", http.StatusMethodNotAllowed, " - ", r.Method, "not allowed")
		return
	}

	// Parse XML body
	r.ParseForm()
	vuiRequestBody := parseXMLRequestBody(r)
	if vuiRequestBody == nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Debug("VUI Request - ", http.StatusBadRequest, " - ", "Failed to parse XML request body")
		return
	}

	vuiResponseBody, err := generateResponse(vuiRequestBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Debug("VUI Request - ", http.StatusInternalServerError, " - ", "Internal Server Error")
		return
	}
	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(http.StatusOK)
	w.Write(vuiResponseBody)
}

func generateResponse(vuiRequestBody *vuiRequestStruct) ([]byte, error) {
	vuiResponseBody := &vuiResponseStruct{}
	if vuiRequestBody.Payload.AliQueryRequest.ExternalKey != "" {
		vuiResponseBody.HDR.REC = "1"
		vuiResponseBody.Payload.AliQueryResponse.ExternalKey = vuiRequestBody.Payload.AliQueryRequest.ExternalKey
		vuiResponseBody.VER = "1.0"
		vuiResponseBody.Payload.AliQueryResponse.VER = "1.0"
		vuiResponseBody.Payload.AliQueryResponse.RC1.RC1 = "0"
		vuiResponseBody.Payload.AliQueryResponse.RC1.Message = "SUCCESS"
		vuiResponseBody.Payload.AliQueryResponse.HNO = "0000006080"
		vuiResponseBody.Payload.AliQueryResponse.STN = "TENNYSON PKWY"
		vuiResponseBody.Payload.AliQueryResponse.MCN = "PLANO"
		vuiResponseBody.Payload.AliQueryResponse.STA = "TX"
		vuiResponseBody.Payload.AliQueryResponse.LOC = "SUITE 400"
		vuiResponseBody.Payload.AliQueryResponse.NAM = "ENTITLEMENT"
		vuiResponseBody.Payload.AliQueryResponse.ClsType = "0"
		vuiResponseBody.Payload.AliQueryResponse.TysType = "0"
		vuiResponseBody.Payload.AliQueryResponse.ESN = "00888"
		vuiResponseBody.Payload.AliQueryResponse.CPD = time.Now().Format("2006-01-02")
		vuiResponseBody.Payload.AliQueryResponse.CPF = "HPE"
		vuiResponseBody.Payload.AliQueryResponse.ZIP = "75024"
		vuiResponseBody.Payload.AliQueryResponse.ALT = "0000000000"
		vuiResponseBody.TrlRec = "1"
	} else if vuiRequestBody.Payload.AliUpdateRequest.ExternalKey != "" {
		vuiResponseBody.HDR.REC = "1"
		vuiResponseBody.Payload.AliUpdateResponse.ExternalKey = vuiRequestBody.Payload.AliUpdateRequest.ExternalKey
		vuiResponseBody.VER = "1.0"
		vuiResponseBody.Payload.AliUpdateResponse.VER = "1.0"
		vuiResponseBody.Payload.AliUpdateResponse.RC1.RC1 = "000"
		vuiResponseBody.Payload.AliUpdateResponse.RC1.Message = "SUCCESS"
		vuiResponseBody.TrlRec = "1"
	}
	return xml.MarshalIndent(vuiResponseBody, "", " ")
}

func parseXMLRequestBody(r *http.Request) *vuiRequestStruct {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Debug("VUI Request - ", "Received body ", string(body))
	requestBody := &vuiRequestStruct{}
	xml.Unmarshal(body, requestBody)
	return requestBody
}
