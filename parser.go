package rakeja

import (
	"strings"

	"github.com/shogo82148/go-mecab"
)

func parse(text *string) ([]string, error) {
	tagger, err := mecab.New(map[string]string{"output-format-type": "wakati"})
	if err != nil {
		return nil, err
	}
	defer tagger.Destroy()

	result, err := tagger.Parse(*text)
	if err != nil {
		return nil, err
	}

	return strings.Split(result, " "), nil
}
