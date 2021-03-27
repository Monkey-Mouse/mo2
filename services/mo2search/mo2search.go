package mo2search

import (
	"os"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	_ "github.com/leopku/bleve-gse-tokenizer"
)

var Indexes map[string]bleve.Index = make(map[string]bleve.Index, 0)

// CreateOrLoadIndex as name
func CreateOrLoadIndex(name string) {

	dir := ".bleve/" + name
	var index bleve.Index
	var err error
	if _, err = os.Stat(dir); !os.IsNotExist(err) {
		index, err = bleve.Open(dir)
	} else {

		mapping := bleve.NewIndexMapping()
		if err := mapping.AddCustomTokenizer("gse", map[string]interface{}{
			"type":       "gse",
			"user_dicts": "./dict.txt", // <-- MUST specified, otherwise panic would occurred.
		}); err != nil {
			panic(err)
		}
		if err := mapping.AddCustomAnalyzer("gse", map[string]interface{}{
			"type":      "gse",
			"tokenizer": "gse",
		}); err != nil {
			panic(err)
		}
		mapping.DefaultAnalyzer = "gse"

		index, err = bleve.New(dir, mapping)
	}
	if err != nil {
		panic(err)
	}
	Indexes[name] = index
}

// Query query
func Query(index string, query query.Query) *bleve.SearchResult {
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Highlight = bleve.NewHighlightWithStyle("html")
	searchRequest.Fields = []string{"*"}
	searchResult, _ := Indexes[index].Search(searchRequest)
	return searchResult
}
