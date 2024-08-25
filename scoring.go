package rakeja

type WordScoring string

const (
	// 他のワード（content word）との共起回数をスコアとする方式。短めのキーフレーズを抽出するのに適している。
	WordScoringDeg WordScoring = "deg"
	// 当該ワードの出現回数をスコアとする方式。単純に出現頻度が高いワードを含むキーフレーズを抽出するのに適している。
	WordScoringFreq WordScoring = "freq"
	// 「他のワード（content word）との共起回数 / 当該ワードの出現回数」をスコアとする方式。長いキーフレーズを抽出するのに適している。デフォルトで選択される方式。
	WordScoringDegToFreq WordScoring = "degToFreq"
)

func WordScoringDegPtr() *WordScoring {
	s := WordScoringDeg
	return &s
}

func WordScoringFreqPtr() *WordScoring {
	s := WordScoringFreq
	return &s
}

func WordScoringDegToFreqPtr() *WordScoring {
	s := WordScoringDegToFreq
	return &s
}
