package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

const addr = "https://localhost:1605/"

func main() {
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{http3.NextProtoH3},
		},
		QUICConfig: &quic.Config{},
	}
	defer roundTripper.Close()
	client := &http.Client{
		Transport: roundTripper,
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		method, _ := reader.ReadString(' ')
		switch method {
		case "GET ":
			path, _ := reader.ReadString('\n')
			path = strings.TrimSuffix(path, "\r\n")
			response, err := client.Get(addr + path)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			fmt.Println(response.Status)
			buf := make([]byte, response.ContentLength)
			pos := 0
			for pos < int(response.ContentLength) {
				n, err := response.Body.Read(buf[pos:])
				pos += n
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
			fmt.Println(string(buf))
		case "POST ":
			path, _ := reader.ReadString('\n')
			path = strings.TrimSuffix(path, "\r\n")
			body := []byte{}
			for {
				line, _ := reader.ReadString('\n')
				line = strings.Replace(line, "\r", "", -1)
				if line == "#end\n" {
					break
				}
				body = append(body, []byte(line)...)
			}

			response, err := client.Post(addr+path, "any", bytes.NewReader(body))
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			fmt.Println(response.Status)
			buf := make([]byte, response.ContentLength)
			pos := 0
			for pos < int(response.ContentLength) {
				n, err := response.Body.Read(buf[pos:])
				pos += n
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
			fmt.Println(string(buf))
		default:
			fmt.Println("unsupported method")
			os.Exit(1)
		}
	}
}
