package main

import (
	"log"
	"net/http"
)

type LogWriter struct {
	responseWriter http.ResponseWriter
}

func NewLogWriter(writer http.ResponseWriter) *LogWriter {
	return &LogWriter{
		responseWriter: writer,
	}
}

func (l *LogWriter) Write(p []byte) (n int, err error) {
	go func() {
		log.Println(string(p))
	}()

	return l.responseWriter.Write(p)
}

func (l *LogWriter) Header() http.Header {
	return l.responseWriter.Header()
}

func (l *LogWriter) WriteHeader(statusCode int) {
	l.responseWriter.WriteHeader(statusCode)
}
