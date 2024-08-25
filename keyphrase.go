package rakeja

import (
	"fmt"
	"sort"
)

type IKeyphrase interface {
	GetText() string
	GetScore() float64
}

type keyphrase struct {
	words []*word
}

func (k keyphrase) GetText() string {
	var text string
	for _, w := range k.words {
		text += w.text
	}
	return text
}

func (k keyphrase) GetScore() float64 {
	var score float64
	for _, w := range k.words {
		score += w.score
	}
	return score
}

func (k *keyphrase) isEmpty() bool {
	return k.words == nil || len(k.words) == 0
}

func (k *keyphrase) append(w *word) {
	k.words = append(k.words, w)
}

// 包含するワードの共起回数をインクリメントする
func (k *keyphrase) incrementDeg() {
	for _, w := range k.words {
		w.deg += len(k.words) // NOTE: 自らも含むため、-1 は不要
	}
}

type IKeyphraseCollection interface {
	List() []IKeyphrase
	ListTexts() []string
	ListScores() []float64
}

type keyphraseCollection struct {
	keyphrases []IKeyphrase
}

// キーフレーズのリストを返す
func (kc *keyphraseCollection) List() []IKeyphrase {
	return kc.keyphrases
}

// キーフレーズのテキストのみを返す
func (kc *keyphraseCollection) ListTexts() []string {
	var texts []string
	for _, k := range kc.keyphrases {
		texts = append(texts, k.GetText())
	}
	return texts
}

// キーフレーズのスコアのみを返す
func (kc *keyphraseCollection) ListScores() []float64 {
	var scores []float64
	for _, k := range kc.keyphrases {
		scores = append(scores, k.GetScore())
	}
	return scores
}

// まだ存在しないキーフレーズの場合のみ追加する
func (kc *keyphraseCollection) appendIfUnq(k IKeyphrase) {
	for _, existingK := range kc.keyphrases {
		if k.GetText() == existingK.GetText() {
			return
		}
	}

	kc.keyphrases = append(kc.keyphrases, k)
}

// スコアの降順に in-place ソートする
func (kc *keyphraseCollection) sortByScoreInDesc() {
	sort.Slice(kc.keyphrases, func(i, j int) bool {
		return kc.keyphrases[i].GetScore() > kc.keyphrases[j].GetScore()
	})
}

func (kc *keyphraseCollection) exportTopNPercent(topNPercent int) *keyphraseCollection {
	fmt.Printf("len(kc.keyphrases): %d\n", len(kc.keyphrases))
	for _, k := range kc.keyphrases {
		fmt.Printf("keyphrase: %s, score: %f\n", k.GetText(), k.GetScore())
	}

	if len(kc.keyphrases) == 0 {
		return &keyphraseCollection{}
	}

	if topNPercent <= 0 {
		return &keyphraseCollection{}
	}

	if topNPercent >= 100 {
		return kc
	}

	n := len(kc.keyphrases) * topNPercent / 100
	fmt.Printf("n: %d\n", n)

	return &keyphraseCollection{keyphrases: kc.keyphrases[:n]}
}
