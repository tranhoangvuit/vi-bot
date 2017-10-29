package vibot

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"golang.org/x/net/html"
)

// to pull the href attribute from a token
func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

func getContain(t html.Token) (ok bool, href string) {
	for _, div := range t.Attr {
		if div.Key == "id" && div.Val == "Definition" {
			fmt.Println(div.Val)
			fmt.Println(div.Key)
			fmt.Println(div.Namespace)
		}
	}
	return
}

func crawl(url string, ch chan string, chFinished chan bool) {
	resp, err := http.Get(url)

	defer func() {
		chFinished <- true
	}()

	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
		return
	}

	b := resp.Body
	defer b.Close()

	z := html.NewTokenizer(b)
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			// End of document
			return
		case tt == html.StartTagToken:
			t := z.Token()

			isAnChor := t.Data == "div"
			if !isAnChor {
				continue
			}

			ok, url := getContain(t)
			if !ok {
				continue
			}

			// Make sure the url begines in http
			hasProto := strings.Index(url, "http") == 0
			if hasProto {
				ch <- url
			}
		}
	}
}

//func (vb *ViBot) GetAllHref(url string) {
//foundUrls := make(map[string]bool)

//// Channels
//chUrls := make(chan string)
//chFinished := make(chan bool)

//// Kick off the crawl process
////for _, url := range seedUrls {
//go crawl(url, chUrls, chFinished)
////}

//// Subcribe to both channel
//for c := 0; c < 1; {
//select {
//case url := <-chUrls:
//foundUrls[url] = true
//case <-chFinished:
//c++
//}
//}

//for url, _ := range foundUrls {
//fmt.Println(" - " + url + " \n")
//}

//close(chUrls)
//}

func (vb *ViBot) GetDefinition(url, searchWord string) {
	doc, err := goquery.NewDocument(url + searchWord)
	if err != nil {
		log.Fatal(err)
	}

	divDefi := doc.Find("#Definition")
	pronoun := divDefi.Find("section[data-src|=\"hc_dict\"] > span.pron").Text()
	fmt.Printf("%s: /%s/ - ", searchWord, pronoun)
	divDefi.Find("section[data-src|=\"hc_dict\"] > div:not(.etyseg, .runseg) > i").Each(func(i int, s *goquery.Selection) {
		kindOfWord := s.Text()
		if kindOfWord != "tr" && kindOfWord != "or" {
			fmt.Printf("%s ,", strings.ToUpper(s.Text()))
		}
	})
	fmt.Printf("\n")
	divDefi.Find("section[data-src|=\"hc_dict\"] > div > div.ds-list").Each(func(i int, s *goquery.Selection) {
		fmt.Printf("   %s\n", s.Text())
	})
	fmt.Printf("Relate: ")
	divDefi.Find("section[data-src|=\"hc_dict\"] > div.runseg").Each(func(i int, s *goquery.Selection) {
		fmt.Printf("%s (%s) ", s.Find("b").Text(), s.Find("i").Text())
	})

	fmt.Printf("\n")
	//fmt.Println("Idioms: ")

	//fmt.Printf("\n")
	//fmt.Printf("Verb")
	//divDefi.Find("section[data-src|=\"hcVerbTblEn\"]").Each(func(i int, s *goquery.Selection) {
	//fmt.Printf("%s - %s\n", s.Find("b > span.hvr").Text(), s.Find("span.hvr").Text())
	//})
	sound := doc.Find("#content > div > span:nth-child(3)")
	if sound != nil {
		data, _ := sound.Attr("data-snd")
		for i := 0; i < 2; i++ {
			vb.PlaySound("http://img2.tfd.com/pron/mp3/" + data + ".mp3")
			time.Sleep(3 * time.Second)
		}
	}
}
