package mo2search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	_ "github.com/leopku/bleve-gse-tokenizer"
)

var searchHost = os.Getenv("MO2_SEARCH_HOST")
var createIndexURL = path.Join(searchHost, "api/index?name=")
var searchURL = path.Join(searchHost, "api/search?name=")
var indexURL = path.Join(searchHost, "api/%s?id=%s")
var deleteURL = path.Join(searchHost, "api/%s?id=%s")

// CreateOrLoadIndex as name
func CreateOrLoadIndex(name string) {
	http.Post(createIndexURL+name, "", nil)
}

func JsonRPC(url string, mthd string, body interface{}, resp interface{}) error {
	sb, err := json.Marshal(body)
	if err != nil {
		return err
	}
	var buf *bytes.Buffer = nil
	if sb != nil {
		buf = bytes.NewBuffer(sb)
	}
	req, err := http.NewRequest(mthd, url, buf)
	if err != nil {
		return err
	}
	re, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp == nil {
		return nil
	}
	bs, err := ioutil.ReadAll(re.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, resp)
}

// Query query
func Query(index string, query query.Query, page int, pagesize int, fields []string) *bleve.SearchResult {
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = fields
	searchRequest.From = page * pagesize
	searchRequest.Size = pagesize
	searchRequest.Highlight = bleve.NewHighlightWithStyle("html")
	searchResult := &bleve.SearchResult{}
	_ = JsonRPC(searchURL+index, http.MethodPost, searchRequest, searchResult)
	return searchResult
}

func Index(index string, id string, document interface{}) {
	JsonRPC(fmt.Sprintf(indexURL, index, id), http.MethodPut, document, nil)
}
func Delete(index string, id string) {
	JsonRPC(fmt.Sprintf(deleteURL, index, id), http.MethodDelete, nil, nil)
}
