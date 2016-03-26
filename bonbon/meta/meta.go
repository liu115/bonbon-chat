package meta

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	// "io/ioutil"
	"golang.org/x/net/html"
	"net/http"
)

// TODO: 改成url=XXXXX&protocal=(http|https)之格式
func HandleURLMeta(c *gin.Context) {
	url := c.Query("url")
	protocal := c.Query("protocal")
	purl := protocal + "://" + url
	println("request url meta " + purl)

	resp, err := http.Get(purl)
	if err != nil {
		fmt.Printf("get html error: %s", err.Error())
		c.String(http.StatusNoContent, "could not get url")
		return
	}
	defer resp.Body.Close()

	meta := make(map[string]string)

	z := html.NewTokenizer(resp.Body)

L:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			break L
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == "meta" {
				switch getValue(t.Attr, "property") {
				case "og:title":
					meta["og:title"] = getValue(t.Attr, "content")
				case "og:description":
					meta["og:description"] = getValue(t.Attr, "content")
				case "og:image":
					meta["image"] = getValue(t.Attr, "content")
				case "og:site_name":
					meta["site_name"] = getValue(t.Attr, "content")
				}
				switch getValue(t.Attr, "name") {
				case "description":
					meta["description"] = getValue(t.Attr, "content")
				}
			} else if t.Data == "title" {
				_ = z.Next()
				t = z.Token()
				meta["title"] = t.Data
			}
		}
	}

	if meta["og:title"] != "" {
		meta["title"] = meta["og:title"]
		delete(meta, "og:title")
	}
	if meta["og:description"] != "" {
		meta["description"] = meta["og:description"]
		delete(meta, "og:description")
	}

	json, err := json.Marshal(meta)
	if err != nil {
		println("unmarshal meta error")
		c.String(http.StatusInternalServerError, "unmarshal error")
		return
	}
	println(string(json))
	c.String(http.StatusOK, string(json))
}

func getValue(attr []html.Attribute, key string) string {
	for _, a := range attr {
		if key == a.Key {
			return a.Val
		}
	}
	return ""
}
