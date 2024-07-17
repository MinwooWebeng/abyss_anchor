package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

type AnchorHTTPHandler struct {
	identity map[string]string
}

func (h *AnchorHTTPHandler) idHandler(response_writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		response_writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	identity, ok := h.identity[request.Form.Get("name")]
	if !ok {
		response_writer.WriteHeader(http.StatusForbidden)
		return
	}
	response_writer.Write([]byte(identity))
}
func (h *AnchorHTTPHandler) stunHandler(response_writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		response_writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	response_writer.Write([]byte(request.RemoteAddr))
}
func (h *AnchorHTTPHandler) regHandler(response_writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		response_writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if request.ContentLength <= 0 || request.ContentLength > 64 {
		response_writer.WriteHeader(http.StatusBadRequest)
		return
	}
	buf := make([]byte, request.ContentLength)
	pos := 0
	for pos < int(request.ContentLength) {
		n, err := request.Body.Read(buf[pos:])
		pos += n
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
	body := string(buf)
	body_split := strings.SplitN(body, "\n", 2)
	name, ok := strings.CutPrefix(body_split[0], "name:")
	if !ok {
		response_writer.WriteHeader(http.StatusBadRequest)
		return
	}
	h.identity[name] = body_split[1]
}

func (h *AnchorHTTPHandler) ServeHTTP(response_writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	fmt.Println(request.URL.String())

	switch request.URL.Path {
	case "/id":
		h.idHandler(response_writer, request)
	case "/stun":
		h.stunHandler(response_writer, request)
	case "/reg":
		h.regHandler(response_writer, request)
	default:
		response_writer.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	handler := &AnchorHTTPHandler{}
	handler.identity = make(map[string]string)
	err := http3.ListenAndServeQUIC("0.0.0.0:1605", "./credentials/cert.pem", "./credentials/private.pem", handler)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	tr := quic.Transport{Conn: &net.UDPConn{}}
	server := http3.Server{}
	ln, _ := tr.ListenEarly(&tls.Config{}, &quic.Config{})
	for {
		c, _ := ln.Accept(context.Background())
		switch c.ConnectionState().TLS.NegotiatedProtocol {
		case http3.NextProtoH3:
			go server.ServeQUICConn(c)
		}
	}
}
