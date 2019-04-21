package aggregator

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

var client = &http.Client{}
var logglyUrl = "https://logs-01.loggly.com/inputs/%s/tag/http/"
var url string

func Start(bulkQueue chan []byte, authToken string) {
	url = fmt.Sprintf(logglyUrl, authToken)
	buf := new(bytes.Buffer)
	for {
		msg := <-bulkQueue // blocks until a message arrives on the channel
		buf.Write(msg)
		buf.WriteString("\n")

		size := buf.Len()

		if size > 1024 {
			sendBulk(*buf)
			buf.Reset()
		}
	}
}

func sendBulk(buffer bytes.Buffer) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(buffer.Bytes()))
	if err != nil {
		logrus.Errorf("Error creating bulk upload HTTP request: %s", err.Error())
		return
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		logrus.Errorf("Error sending bulk: %s", err.Error())
		return
	}

	logrus.Debugf("Successfully sent batch of %v bytes to Loggly\n", buffer.Len())

}
