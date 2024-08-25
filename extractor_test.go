package rakeja

import "testing"

func TestExtractor_Extract(t *testing.T) {
	type testCase struct {
		extractorBuilder func() IExtractor
		text             string
		expectedTexts    []string
		expectedScores   []float64
	}

	var testCases = []testCase{
		{ // default extractor
			extractorBuilder: func() IExtractor {
				return NewDefaultExtractor()
			},
			text:           "今日の最高気温は20度でしたが、明日の最高気温は25度で、明後日の最高気温は30度にもなると、週間天気予報が言っています。",
			expectedTexts:  []string{"週間天気予報", "最高気温"},
			expectedScores: []float64{9, 4},
		},
		{ // custom phrase delimiters
			extractorBuilder: func() IExtractor {
				return NewExtractor(NewExtractorParams{
					PhraseDelimiters: []string{"、", "。", "度"},
				})
			},
			text:           "今日の最高気温は20度でしたが、明日の最高気温は25度で、明後日の最高気温は30度にもなると、週間天気予報が言っています。",
			expectedTexts:  []string{"週間天気予報", "最高気温", "今日"},
			expectedScores: []float64{9, 4, 1},
		},
		{ // empty phrase delimiters
			extractorBuilder: func() IExtractor {
				return NewExtractor(NewExtractorParams{
					PhraseDelimiters: []string{},
					TopNPercent:      IntPtr(100),
				})
			},
			text:           "今日、明日",
			expectedTexts:  []string{"今日、明日\n"}, // NOTE: 外部ライブラリでパースすると何故か必ず最後に改行文字が入ってくる
			expectedScores: []float64{16},       // NOTE: 改行文字があるので9ではなく16
		},
		{ // custom stop words
			extractorBuilder: func() IExtractor {
				return NewExtractor(NewExtractorParams{
					StopWords: []string{"の", "最高", "気温", "は", "でし", "た", "が", "で", "に", "も", "なる", "と", "が", "い", "て", "い", "ます"},
				})
			},
			text:           "今日の最高気温は20度でしたが、明日の最高気温は25度で、明後日の最高気温は30度にもなると、週間天気予報が言っています。",
			expectedTexts:  []string{"週間天気予報", "20度"},
			expectedScores: []float64{9, 4},
		},
		{ // empty stop words
			extractorBuilder: func() IExtractor {
				return NewExtractor(NewExtractorParams{
					StopWords:   []string{},
					TopNPercent: IntPtr(100),
				})
			},
			text:           "今日は晴れ",
			expectedTexts:  []string{"今日は晴れ"},
			expectedScores: []float64{9},
		},
		{ // WordScoringDeg
			extractorBuilder: func() IExtractor {
				return NewExtractor(NewExtractorParams{
					WordScoring: WordScoringDegPtr(),
				})
			},
			text:           "今日の最高気温は20度でしたが、明日の最高気温は25度で、明後日の最高気温は30度にもなると、週間天気予報が言っています。",
			expectedTexts:  []string{"最高気温", "週間天気予報"},
			expectedScores: []float64{12, 9},
		},
		{ // WordScoringFreq
			extractorBuilder: func() IExtractor {
				return NewExtractor(NewExtractorParams{
					WordScoring: WordScoringFreqPtr(),
				})
			},
			text:           "今日の最高気温は20度でしたが、明日の最高気温は25度で、明後日の最高気温は30度にもなると、週間天気予報が言っています。",
			expectedTexts:  []string{"最高気温", "20度"},
			expectedScores: []float64{6, 4},
		},
		{ // top 0%
			extractorBuilder: func() IExtractor {
				return NewExtractor(NewExtractorParams{
					TopNPercent: IntPtr(0),
				})
			},
			text:           "今日の最高気温は20度でしたが、明日の最高気温は25度で、明後日の最高気温は30度にもなると、週間天気予報が言っています。",
			expectedTexts:  []string{},
			expectedScores: []float64{},
		},
		{ // top 50%
			extractorBuilder: func() IExtractor {
				return NewExtractor(NewExtractorParams{
					TopNPercent: IntPtr(50),
				})
			},
			text:           "今日の最高気温は20度でしたが、明日の最高気温は25度で、明後日の最高気温は30度にもなると、週間天気予報が言っています。",
			expectedTexts:  []string{"週間天気予報", "最高気温", "20度", "25度"},
			expectedScores: []float64{9, 4, 4, 4},
		},
		{ // top 100%
			extractorBuilder: func() IExtractor {
				return NewExtractor(NewExtractorParams{
					TopNPercent: IntPtr(100),
				})
			},
			text:           "今日の最高気温は20度でしたが、明日の最高気温は25度で、明後日の最高気温は30度にもなると、週間天気予報が言っています。",
			expectedTexts:  []string{"週間天気予報", "最高気温", "20度", "25度", "30度", "今日", "明日", "明後日", "言っ"},
			expectedScores: []float64{9, 4, 4, 4, 4, 1, 1, 1, 1},
		},
		{ // empty text
			extractorBuilder: func() IExtractor {
				return NewDefaultExtractor()
			},
			text:           "",
			expectedTexts:  []string{},
			expectedScores: []float64{},
		},
	}

	for _, test := range testCases {
		keyphrases, err := test.extractorBuilder().Extract(&test.text)
		if err != nil {
			t.Errorf("Error: %v", err)
			continue
		}

		if len(keyphrases.List()) != len(test.expectedTexts) {
			t.Errorf("Expected: %v\nGot: %v", test.expectedTexts, keyphrases.ListTexts())
			continue
		}

		for i, k := range keyphrases.List() {
			if k.GetText() != test.expectedTexts[i] {
				t.Errorf("Expected: %v\nGot: %v", test.expectedTexts, keyphrases.ListTexts())
			}
			if k.GetScore() != test.expectedScores[i] {
				t.Errorf("Expected: %v\nGot: %v", test.expectedScores, keyphrases.ListScores())
			}
		}
	}
}
