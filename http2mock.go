package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const usageStr = `
HTTP/2 Service Mock -- Emulator of multiple HTTP/2 services.

Usage:

  http2mock <flags>

Flags:
`

const apnsRequestPath = "/3/device/"

func main() {

	// initialize flags
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	v := flagSet.Bool("v", false, "if set, enable DEBUG level logging")
	vv := flagSet.Bool("vv", false, "if set, enable TRACE level logging")
	usage := func() {
		fmt.Fprintf(os.Stderr, "%s", usageStr)
		flagSet.PrintDefaults()
		fmt.Fprintln(os.Stderr)
	}
	flagSet.Usage = usage
	flagSet.Parse(os.Args[1:])

	if *vv {
		log.SetLevel(log.TraceLevel)
	} else if *v {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.SetFormatter(&nested.Formatter{
		TimestampFormat: time.RFC3339,
		HideKeys:        true,
		NoColors:        false,
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc(apnsRequestPath, apnsHandler)
	// default start HTTP/1.1 on :18443
	srv := &http.Server{Addr: ":18443", Handler: mux}

	// start TLS(http/2)
	log.Info("Serving on https://0.0.0.0:18443")
	log.Debug("Loading key and certificate")
	log.Fatal(srv.ListenAndServeTLS("cert/server.crt", "cert/server_plain.key"))
}

// Default handler
func indexHandler(w http.ResponseWriter, r *http.Request) {

	log.Debug("Got connection: ", r.Proto)
	// Default response
	w.Write([]byte("Welcome to DEG."))
}

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
