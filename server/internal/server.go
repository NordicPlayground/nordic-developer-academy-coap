package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func loggingMiddleware(next mux.Handler) mux.Handler {
	return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		log.Printf("ClientAddress %v, %v\n", w.Conn().RemoteAddr(), r.String())
		next.ServeCOAP(w, r)
	})
}

func helloResource(w mux.ResponseWriter, r *mux.Message) {
	err := w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader([]byte(fmt.Sprintf("Hello from the cloud! The time is: %s.", time.Now().Format(time.RFC3339)))))
	if err != nil {
		log.Printf("cannot set response: %v", err)
	}
}

func dynamicResource(client *azblob.Client, containerName string) func(mux.ResponseWriter, *mux.Message) {
	return func(w mux.ResponseWriter, r *mux.Message) {
		resp := w.Conn().AcquireMessage(r.Context())
		defer w.Conn().ReleaseMessage(resp)
		resp.SetToken(r.Token())
		resp.SetContentFormat(message.TextPlain)

		path, pErr := r.Path()
		if pErr != nil {
			resp.SetCode(codes.BadRequest)
			w.Conn().WriteMessage(resp)
			return
		}

		key := strings.Split(path[1:], "/")[0]

		switch r.Code() {
		case codes.PUT:
			payloadSize, err := r.BodySize()
			if err != nil {
				log.Fatal(err)
			}
			if (payloadSize > 100000) { // Max size 100 KB
				err := w.SetResponse(codes.RequestEntityTooLarge, message.TextPlain, bytes.NewReader([]byte("Maximum payload size is 100 KB!")))
				if err != nil {
					log.Printf("cannot set response: %v", err)
				}
				return 
			}
			data, err := io.ReadAll(r.Body())
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Stored content: '%s' in '%s'\n", data, path)
			_, err = Store(client, containerName, key, data)
			if err != nil {
				log.Fatal(err)
			}
			resp.SetBody(bytes.NewReader([]byte("OK")))
		case codes.GET:
			stored, err := Retrieve(client, containerName, key)
			if err != nil {
				log.Printf("Not content found at '%s'\n", path)
				resp.SetCode(codes.NotFound)
			} else {
				log.Printf("Loaded content: '%s' from '%s'\n", stored, path)
				resp.SetCode(codes.Content)
				resp.SetBody(bytes.NewReader([]byte(stored)))
			}
		}

		err := w.Conn().WriteMessage(resp)
		if err != nil {
			log.Printf("cannot set response: %v", err)
		}
	}
}

/**
 * Log usage metrics in the format that is supported by Log Analytics agent in Azure Monitor
 *
 * See https://learn.microsoft.com/en-us/azure/azure-monitor/agents/data-sources-custom-logs
*/
func logMetrics (dtls bool, network string) func (next mux.Handler) mux.Handler {
	return func (next mux.Handler) mux.Handler {
		return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
			currentTime := time.Now().Format("2006-01-02T15:04:05Z07:00")
			var protocol string
			if dtls {
				protocol = "dTLS"
			} else {
				protocol = "UDP"
			}
			metricLog := fmt.Sprintf("%s,%s:%s,request", currentTime, protocol, network)
			metricLogsFilePath := fmt.Sprintf("/var/log/academy/%s-coap.log", time.Now().Format("2006-01-02"))
			metricLogsFile, err := os.OpenFile(metricLogsFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Printf("cannot open log file: %v", err)
			}
			defer metricLogsFile.Close()
			_, err = metricLogsFile.WriteString(metricLog + "\n")
			if err != nil {
				log.Printf("cannot write to log file: %v", err)
			}
			log.Printf("Logged metric: %s to %s", metricLog, metricLogsFilePath)
			next.ServeCOAP(w, r)
		})
	}
}

func NewServer(client *azblob.Client, containerName string, dtls bool, network string) *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.Use(logMetrics(dtls, network))
	r.Handle("/static/hello", mux.HandlerFunc(helloResource))
	r.Handle("/{res:[^\\/]+}", mux.HandlerFunc(dynamicResource(client, containerName)))
	return r
}