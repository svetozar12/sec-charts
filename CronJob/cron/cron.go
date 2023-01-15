package cron

import (
	Bytes "bytes"
	"encoding/xml"
	model "example/hello/xmlModel"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"golang.org/x/net/html/charset"
)

// https://www.sec.gov/cgi-bin/browse-edgar?action=getcompany&CIK=0001474464&type=&dateb=&owner=exclude&start=0&count=40&output=atom
// https://www.sec.gov/cgi-bin/browse-edgar?company=Microsoft&action=getcompany
func InitCronJob() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(5).Second().Do(func() {

		client := &http.Client{}
		req, _ := http.NewRequest("GET", "https://www.sec.gov/cgi-bin/browse-edgar?company=Microsoft&action=getcompany&output=atom", nil)
		req.Header.Set("User-Agent", "Sample Company Name AdminContact@<sample company domain>.com")
		res, _ := client.Do(req)
		bytes, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		var d model.Feed
		var err error
		reader := Bytes.NewReader(bytes)
		decoder := xml.NewDecoder(reader)
		decoder.CharsetReader = charset.NewReaderLabel
		err = decoder.Decode(&d)
		if err != nil {
			panic(err)
		}
		for i := 0; i < len(d.Entry); i++ {
			fmt.Println(d.Entry[i].Link.Href)

		}
	})
	s.StartBlocking()
}

// TODO: extract getXml login in separate function
