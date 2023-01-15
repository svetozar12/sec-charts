package cron

import (
	Bytes "bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"golang.org/x/net/html/charset"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	Author  struct {
		Text  string `xml:",chardata"`
		Email string `xml:"email"`
		Name  string `xml:"name"`
	} `xml:"author"`
	Entry []struct {
		Text    string `xml:",chardata"`
		Content struct {
			Text        string `xml:",chardata"`
			Type        string `xml:"type,attr"`
			CompanyInfo struct {
				Text      string `xml:",chardata"`
				Addresses struct {
					Text    string `xml:",chardata"`
					Address []struct {
						Text    string `xml:",chardata"`
						Type    string `xml:"type,attr"`
						City    string `xml:"city"`
						Phone   string `xml:"phone"`
						State   string `xml:"state"`
						Street1 string `xml:"street1"`
						Zip     string `xml:"zip"`
					} `xml:"address"`
				} `xml:"addresses"`
				CancelledMaFlag      string `xml:"cancelled-ma-flag"`
				Cik                  string `xml:"cik"`
				FiscalYearEnd        string `xml:"fiscal-year-end"`
				LastDate             string `xml:"last-date"`
				Name                 string `xml:"name"`
				RevokeFlag           string `xml:"revoke-flag"`
				RevokedMaFlag        string `xml:"revoked-ma-flag"`
				State                string `xml:"state"`
				StateOfIncorporation string `xml:"state-of-incorporation"`
				IrsNumber            string `xml:"irs-number"`
				Sic                  string `xml:"sic"`
			} `xml:"company-info"`
		} `xml:"content"`
		ID   string `xml:"id"`
		Link struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Summary struct {
			Text string `xml:",chardata"`
			Type string `xml:"type,attr"`
		} `xml:"summary"`
		Title   string `xml:"title"`
		Updated string `xml:"updated"`
	} `xml:"entry"`
	ID   string `xml:"id"`
	Link []struct {
		Text string `xml:",chardata"`
		Href string `xml:"href,attr"`
		Rel  string `xml:"rel,attr"`
		Type string `xml:"type,attr"`
	} `xml:"link"`
	Title   string `xml:"title"`
	Updated string `xml:"updated"`
}

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
		var d Feed
		var err error
		reader := Bytes.NewReader(bytes)
		decoder := xml.NewDecoder(reader)
		decoder.CharsetReader = charset.NewReaderLabel
		err = decoder.Decode(&d)
		if err != nil {
			panic(err)
		}
		fmt.Println(d.Author)
	})
	s.StartBlocking()
}
