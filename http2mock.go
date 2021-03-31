package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

const usageStr = `
HTTP/2 Service Mock -- Emulator of multiple HTTP/2 services.

Usage:

  http2mock <flags>

Flags:
`

// APNS Mock Request Path
const apnsRequestPath = "/3/device/"

// VUI Mock Request Path
const vuiRequestPath = "/vui/VuiServlet"

func main() {

	// initialize flags
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	v := flagSet.Bool("v", false, "if set, enable DEBUG level logging")
	vv := flagSet.Bool("vv", false, "if set, enable TRACE level logging")
	port := flagSet.String("p", "18443", "port to listen")
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
	// Router
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc(apnsRequestPath, apnsHandler)
	mux.HandleFunc(vuiRequestPath, vuiHandler)
	// default start HTTP/1.1 on :18443
	srv := &http.Server{Addr: ":" + *port, Handler: mux}

	// start TLS(http/2)
	log.Info("Serving on https://0.0.0.0:" + *port)
	log.Debug("Loading key and certificate")
	log.Fatal(srv.ListenAndServeTLS("cert/server.crt", "cert/server_plain.key"))
}

// Default handler
func indexHandler(w http.ResponseWriter, r *http.Request) {

	log.Debug("Got connection: ", r.Proto)
	// Default response
	w.Write([]byte("Welcome to DEG."))
}
