package main

import (
	"net/http"
	"regexp"
	"strings"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// APNS handler
func apnsHandler(w http.ResponseWriter, r *http.Request) {

	log.Debug("Got connection: ", r.Proto)

	// Set APNS ID into response header
	setApnsId(w, r)

	// Check Request Method
	log.Debug("Request Method: ", r.Method)
	if strings.ToUpper(r.Method) != "POST" {
		apnsRespError(w, http.StatusMethodNotAllowed, "MethodNotAllowed")
		return
	}

	// Check Device Token

	regExDeviceToken := regexp.MustCompile("^[[:xdigit:]]+$")
	deviceToken := r.URL.Path[len(apnsRequestPath):]
	log.Debug("Device Token: ", deviceToken)
	matched := regExDeviceToken.MatchString(deviceToken)
	log.Debug("Does Device Token match hexadecimal digit? ", matched)
	if !matched {
		apnsRespError(w, http.StatusBadRequest, "BadDeviceToken")
		return
	}

	// Check APNS Topic
	apnsTopic := r.Header.Get("apns-topic")
	log.Debug("APNS Topic: ", apnsTopic)
	if apnsTopic == "" {
		apnsRespError(w, http.StatusBadRequest, "MissingTopic")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func setApnsId(w http.ResponseWriter, r *http.Request) {
	apnsId := r.Header.Get("apns-id")
	if apnsId == "" {
		apnsId = uuid.NewV4().String()
	}
	w.Header().Set("apns-id", apnsId)
}

func apnsRespError(w http.ResponseWriter, statusCode int, desc string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(desc))
	log.Info(statusCode, " - ", desc)
}
