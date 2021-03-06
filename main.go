package main

import (
	"encoding/csv"
	"github.com/satyrius/gonx"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

const nginxLogFormat = "$remote_addr $host [$time_local] \"$request\" " +
	"$status $body_bytes_sent \"$http_referer\" " +
	"\"$http_user_agent\" \"$http_x_forwarded_for\" $request_time" +
	" | OMINI_F: [$http_x_operamini_features]" +
	" | X-Request-Id:[$http_x_request_id] \"$document_root\""

func main() {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	reader := gonx.NewReader(os.Stdin, nginxLogFormat)

	for {
		entry, err := reader.Read()
		if err == io.EOF {
			break
		}

		writeEntryToCsv(entry, writer)
	}
}

func writeEntryToCsv(entry *gonx.Entry, writer *csv.Writer) {
	remoteAddr, _ := entry.Field("remote_addr")
	timeLocal, _ := entry.Field("time_local")
	status, _ := entry.Field("status")
	bodyBytesSent, _ := entry.Field("body_bytes_sent")
	host, _ := entry.Field("host")
	request, _ := entry.Field("request")

	if strings.Count(request, " ") < 2 {
		log.Default().Printf("$request expected to contain 2 spaces, got: '%s'", request)
		return
	}

	firstSpace := strings.Index(request, " ")
	lastSpace := strings.LastIndex(request, " ")
	requestUrl := request[firstSpace+1 : lastSpace]
	uri, uriErr := url.ParseRequestURI(requestUrl)

	if uriErr != nil {
		log.Default().Printf("Unable to parse uri='%s' with err='%s'", requestUrl, uriErr)
		return
	}

	httpVerb := request[:firstSpace]

	parsedDateTime, _ := time.Parse("02/Jan/2006:15:04:05 -0700", timeLocal)

	err := writer.Write([]string{
		parsedDateTime.UTC().Format("2006-01-02"),
		parsedDateTime.UTC().Format("2006-01-02 15:04:05"),
		remoteAddr,
		uri.Path,
		buildRoute(uri),
		httpVerb,
		status,
		bodyBytesSent,
		host,
	})
	checkError("Cannot write to file", err)
}

func buildRoute(uri *url.URL) string {
	return regexp.
		MustCompile(`/\d+(/|$)`).
		ReplaceAllString(uri.Path, `/{id}$1`)
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
