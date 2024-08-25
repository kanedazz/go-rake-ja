package rakeja

type wordMap map[string]*word

// 当該 `text` を包含する `word` を返す。存在しない場合は新たに作成する。
func (wm *wordMap) getWord(text string) *word {
	if _, ok := (*wm)[text]; !ok {
		(*wm)[text] = &word{text: text}
	}
	return (*wm)[text]
}

func (wm *wordMap) calcScores(scoring WordScoring) {
	for _, w := range *wm {
		w.calcScore(scoring)
	}
}
