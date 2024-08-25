package rakeja

import "fmt"

const defaultTopNPercent = 33

type IExtractor interface {
	// `text` からキーフレーズを抽出する
	Extract(text *string) (IKeyphraseCollection, error)
}

type extractor struct {
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

	fmt.Printf("words: %v\n", words)

	/*
		- キーフレーズの候補を作成する
		- 各ワードの出現回数と共起回数をカウントする
	*/
	candidateCollection := &keyphraseCollection{}
	candidate := keyphrase{}
	for i, word := range words {
		// ストップワードの場合の処理
		isStopWord := false
		for _, sw := range e.stopWords {
			if sw == word {
				fmt.Printf("stop word: %s ; candidate: %s(%+v)\n", word, candidate.GetText(), candidate)
				if !candidate.isEmpty() {
					candidate.incrementDeg()
					candidateCollection.appendIfUnq(candidate)
					candidate = keyphrase{}
				}
				isStopWord = true
				break
			}
		}
		if isStopWord {
			continue
		}

		// 区切り文字の場合の処理
		isPhraseDelimiter := false
		for _, pd := range e.phraseDelimiters {
			if pd == word {
				fmt.Printf("phrase delimiter: %s ; candidate: %s(%+v)\n", word, candidate.GetText(), candidate)
				if !candidate.isEmpty() {
					candidate.incrementDeg()
					candidateCollection.appendIfUnq(candidate)
					candidate = keyphrase{}
				}
				isPhraseDelimiter = true
				break
			}
		}
		if isPhraseDelimiter {
			continue
		}

		w := wm.getWord(word)
		w.freq++

		candidate.append(w)

		// 最後の単語の場合の処理
		if i == len(words)-1 {
			fmt.Printf("last word: %s ; candidate: %s(%+v)\n", word, candidate.GetText(), candidate)
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
		fmt.Printf("keyphrase: %s, score: %f\n", k.GetText(), k.GetScore())
	}

	/*
		- 上位 `topNPercent` パーセントのキーフレーズを抽出する
	*/
	collection := candidateCollection.exportTopNPercent(e.topNPercent)
	// print
	for _, k := range collection.List() {
		fmt.Printf("keyphrase: %s, score: %f\n", k.GetText(), k.GetScore())
	}

	return collection, nil
}

func NewDefaultExtractor() IExtractor {
	return &extractor{
		phraseDelimiters: defaultPhraseDelimiters,
		stopWords:        defaultStopWords,
		wordScoring:      WordScoringDegToFreq,
		topNPercent:      defaultTopNPercent,
	}
}

type NewExtractorParams struct {
	PhraseDelimiters []string
	StopWords        []string
	WordScoring      *WordScoring
	TopNPercent      *int
}

func NewExtractor(params NewExtractorParams) IExtractor {
	e := &extractor{}

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
