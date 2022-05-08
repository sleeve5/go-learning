package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type DictReqCaiyun struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
	UserID    string `json:"user_id"`
}

type DictReqHuoshan struct {
	//SourceLanguage string   `json:"source_language"`
	//TargetLanguage string   `json:"target_language"`
	//HomeLanguage   string   `json:"home_language"`
	//Category       string   `json:"category"`
	//GlossaryList   []string `json:"glossary_list"`
	Text     string `json:"text"`
	Language string `json:"language"`
}

type DictRespCaiyun struct {
	Rc   int `json:"rc"`
	Wiki struct {
		KnownInLaguages int `json:"known_in_laguages"`
		Description     struct {
			Source string      `json:"source"`
			Target interface{} `json:"target"`
		} `json:"description"`
		ID   string `json:"id"`
		Item struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"item"`
		ImageURL  string `json:"image_url"`
		IsSubject string `json:"is_subject"`
		Sitelink  string `json:"sitelink"`
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string      `json:"explanations"`
		Synonym      []string      `json:"synonym"`
		Antonym      []string      `json:"antonym"`
		WqxExample   [][]string    `json:"wqx_example"`
		Entry        string        `json:"entry"`
		Type         string        `json:"type"`
		Related      []interface{} `json:"related"`
		Source       string        `json:"source"`
	} `json:"dictionary"`
}

type DictRespHuoshan struct {
	Words []struct {
		Source  int    `json:"source"`
		Text    string `json:"text"`
		PosList []struct {
			Type      int `json:"type"`
			Phonetics []struct {
				Type int    `json:"type"`
				Text string `json:"text"`
			} `json:"phonetics"`
			Explanations []struct {
				Text     string `json:"text"`
				Examples []struct {
					Type      int `json:"type"`
					Sentences []struct {
						Text      string `json:"text"`
						TransText string `json:"trans_text"`
					} `json:"sentences"`
				} `json:"examples"`
				Synonyms []interface{} `json:"synonyms"`
			} `json:"explanations"`
			Relevancys []interface{} `json:"relevancys"`
		} `json:"pos_list"`
	} `json:"words"`
	Phrases  []interface{} `json:"phrases"`
	BaseResp struct {
		StatusCode    int    `json:"status_code"`
		StatusMessage string `json:"status_message"`
	} `json:"base_resp"`
}

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, `usage: dict WORD
example: go run main.go good
`)
		os.Exit(1)
	}
	word := os.Args[1]
	tik := time.Now()
	queryCaiyun(word)
	queryHuoshan(word)
	tok := time.Since(tik)
	fmt.Println("用时:", tok)
}

func queryHuoshan(word string) {
	client := &http.Client{}
	//var data = strings.NewReader(`{"source_language":"detect","target_language":"zh","text":"good\n","home_language":"en","category":"","glossary_list":["ailab/menu"]}`)
	// request := DictReqHuoshan{SourceLanguage: "detect", TargetLanguage: "zh", Text: word, HomeLanguage: "en", Category: "", GlossaryList: []string{"ailab/menu"}}
	// var data = strings.NewReader(`{"text":"good","language":"en"}`)
	request := DictReqHuoshan{Text: word, Language: "en"}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://translate.volcengine.com/web/dict/match/v1/?msToken=&X-Bogus=DFSzswVLQDc7EiQrSW0Wk2UClLHg&_signature=_02B4Z6wo00001QX9qMQAAIDALnfzARghKFkF.axAACMLLr9Y.U-vgGXmkokwsubIjGF1lGEYJ7L8p5wR1vTd0cq.WWjc32r53a7oWjceMyfOAt5eJEjlbN5yVADvEZX5BoPWdLoxwADozjJr49", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "translate.volcengine.com")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "en,zh;q=0.9,zh-CN;q=0.8")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", "x-jupiter-uuid=16520014523241606; s_v_web_id=verify_e185db78b90ce93dc6709c6538e1037c; _tea_utm_cache_2018=undefined; ttcid=087ca460006a40878678cf1eb6f4f16729; i18next=translate; referrer_title=%E6%9C%BA%E5%99%A8%E7%BF%BB%E8%AF%91-%E7%81%AB%E5%B1%B1%E5%BC%95%E6%93%8E; csrfToken=8b13c7ddf17a645320fff792b453c01c; __tea_cookie_tokens_3569=%257B%2522web_id%2522%253A%25227095306200476714499%2522%252C%2522ssid%2522%253A%25222218f3be-6c6d-436e-aea1-212c1f58d44b%2522%252C%2522user_unique_id%2522%253A%25227095306200476714499%2522%252C%2522timestamp%2522%253A1652004721452%257D; isIntranet=-1; tt_scid=yBA2mHgR8.2XBxBITjvUdzb1gb3kYYUHFmh.dCX41jYgd6i8p5A4v0GTmISXlPzp1496")
	req.Header.Set("origin", "https://translate.volcengine.com")
	req.Header.Set("referer", "https://translate.volcengine.com/translate?category=&home_language=zh&source_language=detect&target_language=zh&text=good")
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="101", "Google Chrome";v="101"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36")
	resp, err := client.Do(req)
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("Bad StatusCode:", resp.StatusCode, "body", string(bodyText))
	}
	var dictResponse DictRespHuoshan
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("释义:", dictResponse.Words[0].PosList[0].Explanations[0].Text)
	fmt.Printf("例句:")
	fmt.Println(dictResponse.Words[0].PosList[0].Explanations[0].Examples[0].Sentences[0].Text, dictResponse.Words[0].PosList[0].Explanations[0].Examples[0].Sentences[0].TransText)
}

func queryCaiyun(word string) {
	client := &http.Client{}
	// var data = strings.NewReader(`{"trans_type":"en2zh","source":"good"}`)
	request := DictReqCaiyun{TransType: "en2zh", Source: word}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en,zh;q=0.9,zh-CN;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("Referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36")
	req.Header.Set("X-Authorization", "token:qgemv4jr1y38jyq6vhvi")
	req.Header.Set("app-name", "xy")
	req.Header.Set("os-type", "web")
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="101", "Google Chrome";v="101"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	resp, err := client.Do(req)
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("Bad StatusCode:", resp.StatusCode, "body", string(bodyText))
	}
	var dictResponse DictRespCaiyun
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("单词:", word)
	fmt.Println("读音: UK:", dictResponse.Dictionary.Prons.En, "US:", dictResponse.Dictionary.Prons.EnUs)
	//for _, item := range dictResponse.Dictionary.Explanations {
	//	fmt.Println(item)
	//}
}
