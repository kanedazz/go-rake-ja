package rakeja

import "testing"

func TestExtractor_Extract(t *testing.T) {
	type testCase struct {
		title            string
		extractorBuilder func() IExtractor
		text             string
		expectedTexts    []string
		expectedScores   []float64
	}

	var testCases = []testCase{
		{ // default extractor
			title: "default extractor",
			extractorBuilder: func() IExtractor {
				return NewDefaultExtractor()
			},
			text:           "昨日の最高気温は15度で涼しく、今日の最高気温も20度で暑くないですが、明日の最高気温は25度でやや暑く、明後日の最高気温は30度でとても暑くなると、週間天気予報が言っています。",
			expectedTexts:  []string{"週間天気予報", "25度", "最高気温", "15度"},
			expectedScores: []float64{9, 4, 4, 4},
		},
		{ // custom part-of-speech list for content words
			title: "custom part-of-speech list for content words",
			extractorBuilder: func() IExtractor {
				return NewExtractor(NewExtractorParams{
					Pos4ContentWords: []string{Noun}, // noun only
				})
			},
			text:           "昨日の最高気温は15度で涼しく、今日の最高気温も20度で暑くないですが、明日の最高気温は25度でやや暑く、明後日の最高気温は30度でとても暑くなると、週間天気予報が言っています。",
			expectedTexts:  []string{"週間天気予報", "最高気温", "15度"},
			expectedScores: []float64{9, 4, 4},
		},
		{ // custom phrase delimiters
			title: "custom phrase delimiters",
			extractorBuilder: func() IExtractor {
				return NewExtractor(NewExtractorParams{
					PhraseDelimiters: []string{"、", "。", "度"},
				})
			},
			text:           "昨日の最高気温は15度で涼しく、今日の最高気温も20度で暑くないですが、明日の最高気温は25度でやや暑く、明後日の最高気温は30度でとても暑くなると、週間天気予報が言っています。",
			expectedTexts:  []string{"週間天気予報", "最高気温", "暑く", "涼しく"},
			expectedScores: []float64{9, 4, 1, 1},
		},
		{ // empty phrase delimiters
			title: "empty phrase delimiters",
			extractorBuilder: func() IExtractor {
				return NewExtractor(NewExtractorParams{
					Pos4ContentWords: []string{Noun, Symbol},
					PhraseDelimiters: []string{},
					TopNPercent:      IntPtr(100),
				})
			},
			text:           "今日、明日",
			expectedTexts:  []string{"今日、明日"},
			expectedScores: []float64{9},
		},
		{ // custom stop words
			title: "custom stop words",
			extractorBuilder: func() IExtractor {
				return NewExtractor(NewExtractorParams{
					StopWords: []string{"最高"},
				})
			},
			text:           "昨日の最高気温は15度で涼しく、今日の最高気温も20度で暑くないですが、明日の最高気温は25度でやや暑く、明後日の最高気温は30度でとても暑くなると、週間天気予報が言っています。",
			expectedTexts:  []string{"週間天気予報", "25度", "15度", "30度"},
			expectedScores: []float64{9, 4, 4, 4},
		},
		{ // empty stop words
			title: "empty stop words",
			extractorBuilder: func() IExtractor {
				return NewExtractor(NewExtractorParams{
					Pos4ContentWords: []string{Noun, Particle, Adjective},
					StopWords:        []string{},
					TopNPercent:      IntPtr(100),
				})
			},
			text:           "今日は晴れ",
			expectedTexts:  []string{"今日は晴れ"},
			expectedScores: []float64{9},
		},
		{ // WordScoringDeg
			title: "WordScoringDeg",
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
			title: "WordScoringFreq",
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
			title: "top 0%",
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
			title: "top 50%",
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
			title: "top 100%",
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
			title: "empty text",
			extractorBuilder: func() IExtractor {
				return NewDefaultExtractor()
			},
			text:           "",
			expectedTexts:  []string{},
			expectedScores: []float64{},
		},
	}

	for _, test := range testCases {
		if len(test.expectedTexts) != len(test.expectedScores) {
			t.Errorf("Testing: %v\nlen(test.expectedTexts): %d does not match len(test.expectedScores): %d", test.title, len(test.expectedTexts), len(test.expectedScores))
			continue
		}

		keyphrases, err := test.extractorBuilder().Extract(&test.text)
		if err != nil {
			t.Errorf("Error: %v", err)
			continue
		}

		if len(keyphrases.List()) != len(test.expectedTexts) {
			t.Errorf("Testing: %v\nExpected: %v\nGot: %v", test.title, test.expectedTexts, keyphrases.ListTexts())
			continue
		}

		for i, k := range keyphrases.List() {
			if k.GetText() != test.expectedTexts[i] {
				t.Errorf("Testing: %v\nExpected: %v\nGot: %v", test.title, test.expectedTexts, keyphrases.ListTexts())
			}
			if k.GetScore() != test.expectedScores[i] {
				t.Errorf("Testing: %v\nExpected: %v\nGot: %v", test.title, test.expectedScores, keyphrases.ListScores())
			}
		}
	}
}
