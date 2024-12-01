package rakeja

import (
	"strings"

	"github.com/shogo82148/go-mecab"
)

type parsedWord struct {
	txt string
	pos string
}

func parse(text *string) ([]parsedWord, error) {
	tagger, err := mecab.New(map[string]string{
		"node-format": "%m\t%H\n",
	})
	if err != nil {
		return nil, err
	}
	defer tagger.Destroy()

	result, err := tagger.Parse(*text)
	if err != nil {
		return nil, err
	}

	var words []parsedWord
	for _, line := range strings.Split(result, "\n") {
		debugf("line: %s\n", line)
		if line == "EOS" {
			break
		}

		parts := strings.Split(line, "\t")
		debugf("parts[0] %v\n", parts[0])
		debugf("parts[1] %v\n", parts[1])

		words = append(words, parsedWord{
			txt: parts[0],
			pos: strings.Split(parts[1], ",")[0],
		})
	}

	return words, nil
}
