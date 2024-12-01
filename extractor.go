package rakeja

const defaultTopNPercent = 33

type IExtractor interface {
	// `text` からキーフレーズを抽出する
	Extract(text *string) (IKeyphraseCollection, error)
}

type extractor struct {
	pos4ContentWords []string
	phraseDelimiters []string
	stopWords        []string
	wordScoring      WordScoring
	// スコアが上位 `topNPercent` パーセントのキーフレーズのみを抽出する
	topNPercent int
}

func (e *extractor) Extract(text *string) (IKeyphraseCollection, error) {
	wm := make(wordMap)

	/*
		- `text` を単語に分割する
	*/
	words, err := parse(text)
	if err != nil {
		return nil, err
	}

	debugf("words: %v\n", words)

	/*
		- キーフレーズの候補を作成する
		- 各ワードの出現回数と共起回数をカウントする
	*/
	candidateCollection := &keyphraseCollection{}
	candidate := keyphrase{}
	for i, word := range words {
		// 内容語として扱わない品詞の場合の処理
		isContentWordPos := false
		for _, pos := range e.pos4ContentWords {
			if pos == word.pos {
				isContentWordPos = true
				break
			}
		}
		if !isContentWordPos {
			debugf("not content word: %s ; pos: %s\n", word.txt, word.pos)

			if !candidate.isEmpty() {
				candidate.incrementDeg()
				candidateCollection.appendIfUnq(candidate)
				candidate = keyphrase{}
			}
			continue
		}

		// ストップワードの場合の処理
		isStopWord := false
		for _, sw := range e.stopWords {
			if sw == word.txt {
				debugf("stop word: %s ; candidate: %s(%+v)\n", word, candidate.GetText(), candidate)
				isStopWord = true
				break
			}
		}
		if isStopWord {
			if !candidate.isEmpty() {
				candidate.incrementDeg()
				candidateCollection.appendIfUnq(candidate)
				candidate = keyphrase{}
			}
			continue
		}

		// 区切り文字の場合の処理
		isPhraseDelimiter := false
		for _, pd := range e.phraseDelimiters {
			if pd == word.txt {
				debugf("phrase delimiter: %s ; candidate: %s(%+v)\n", word, candidate.GetText(), candidate)
				isPhraseDelimiter = true
				break
			}
		}
		if isPhraseDelimiter {
			if !candidate.isEmpty() {
				candidate.incrementDeg()
				candidateCollection.appendIfUnq(candidate)
				candidate = keyphrase{}
			}
			continue
		}

		w := wm.getWord(word.txt)
		w.freq++

		candidate.append(w)

		// 最後の単語の場合の処理
		if i == len(words)-1 {
			debugf("last word: %s ; candidate: %s(%+v)\n", word, candidate.GetText(), candidate)
			if !candidate.isEmpty() {
				candidate.incrementDeg()
				candidateCollection.appendIfUnq(candidate)
				candidate = keyphrase{}
			}
		}
	}

	/*
		- TODO: キーフレーズの連結（adjoining keyphrases）
	*/

	/*
		- 各ワードのスコアを計算する
	*/
	wm.calcScores(e.wordScoring)

	/*
		- キーフレーズの候補をスコアの降順にソートする
	*/
	candidateCollection.sortByScoreInDesc()
	// print
	for _, k := range candidateCollection.List() {
		debugf("keyphrase: %s, score: %f\n", k.GetText(), k.GetScore())
	}

	/*
		- 上位 `topNPercent` パーセントのキーフレーズを抽出する
	*/
	collection := candidateCollection.exportTopNPercent(e.topNPercent)
	// print
	for _, k := range collection.List() {
		debugf("keyphrase: %s, score: %f\n", k.GetText(), k.GetScore())
	}

	return collection, nil
}

func NewDefaultExtractor() IExtractor {
	return &extractor{
		pos4ContentWords: defaultPos4ContentWords,
		phraseDelimiters: defaultPhraseDelimiters,
		stopWords:        defaultStopWords,
		wordScoring:      WordScoringDegToFreq,
		topNPercent:      defaultTopNPercent,
	}
}

type NewExtractorParams struct {
	// 指定された場合、このスライスに含まれる品詞のみを内容語として扱う。
	Pos4ContentWords []string
	PhraseDelimiters []string
	StopWords        []string
	WordScoring      *WordScoring
	TopNPercent      *int
}

func NewExtractor(params NewExtractorParams) IExtractor {
	e := &extractor{}

	if params.Pos4ContentWords == nil {
		e.pos4ContentWords = defaultPos4ContentWords
	} else {
		e.pos4ContentWords = params.Pos4ContentWords
	}

	if params.PhraseDelimiters == nil {
		e.phraseDelimiters = defaultPhraseDelimiters
	} else {
		e.phraseDelimiters = params.PhraseDelimiters
	}

	if params.StopWords == nil {
		e.stopWords = defaultStopWords
	} else {
		e.stopWords = params.StopWords
	}

	if params.WordScoring == nil {
		e.wordScoring = WordScoringDegToFreq
	} else {
		e.wordScoring = *params.WordScoring
	}

	if params.TopNPercent == nil {
		e.topNPercent = defaultTopNPercent
	} else {
		e.topNPercent = *params.TopNPercent
	}

	return e
}
