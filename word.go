package rakeja

type word struct {
	text string
	// 当該ワードの出現回数
	freq int
	// 他のワード（content word）との共起回数
	deg int
	// スコア
	score float64
}

func (w *word) calcScore(scoring WordScoring) {
	if scoring == WordScoringDeg {
		w.score = float64(w.deg)
	} else if scoring == WordScoringFreq {
		w.score = float64(w.freq)
	} else {
		w.score = float64(w.deg) / float64(w.freq)
	}
}
