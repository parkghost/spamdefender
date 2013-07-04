package analyzer

import (
	"bytes"
	"fmt"
	"github.com/jbrukh/bayesian"
	"spamdefender/analyzer/goseg"
	"spamdefender/common"
)

const (
	Good   bayesian.Class = "Good"
	Bad    bayesian.Class = "Bad"
	cutset                = ":;=<>"
)

type Analyzer struct {
	classifier *bayesian.Classifier
	tokenizer  *goseg.Tokenizer
}

func (a *Analyzer) Test(text string) ([]float64, bool) {
	score, likely, _ := a.classifier.LogScores(common.Normalize(a.tokenizer.Cut([]rune(text)), cutset))
	if a.classifier.Classes[likely] == Bad {
		return score, false
	}
	return score, true
}

type WordFreq struct {
	Word       string
	FreqMatrix []float64
}

func (wf *WordFreq) String() string {
	buf := bytes.NewBufferString(wf.Word)

	outputFreqsText := false
	freqsText := bytes.NewBufferString("(")
	for c, v := range wf.FreqMatrix {
		if v > 0.0000001 {
			outputFreqsText = true
		}

		if c == 0 {
			freqsText.WriteString(fmt.Sprintf("%f", v))
		} else {
			freqsText.WriteString(fmt.Sprintf(",%f", v))
		}
	}
	freqsText.WriteString(")")

	if outputFreqsText {
		freqsText.WriteTo(buf)
	}

	return buf.String()
}

type WordFreqList []WordFreq

func (wfl WordFreqList) String() string {
	buf := bytes.NewBufferString("")
	buf.WriteString("[")
	for i, wf := range wfl {
		if i == 0 {
			buf.WriteString(wf.String())
		} else {
			buf.WriteString(", " + wf.String())
		}
	}
	buf.WriteString("]")
	return buf.String()
}

func (a *Analyzer) Explain(text string) WordFreqList {
	words := a.tokenizer.Cut([]rune(text))
	normalizedWords := common.Normalize(words, cutset)
	freqMatrix := a.classifier.WordFrequencies(normalizedWords)

	wordFreqs := make([]WordFreq, len(normalizedWords))

	for i, _ := range normalizedWords {
		wordFreqMatrix := make([]float64, len(a.classifier.Classes))

		for j, _ := range a.classifier.Classes {
			wordFreqMatrix[j] = freqMatrix[j][i]
		}

		wordFreqs[i] = WordFreq{normalizedWords[i], wordFreqMatrix}
	}

	return wordFreqs
}

func NewAnalyzer(traningDataFilePath, dictFilePath string) (*Analyzer, error) {
	classifier, err := bayesian.NewClassifierFromFile(traningDataFilePath)
	if err != nil {
		return nil, err
	}

	tokenizer, err := goseg.NewTokenizer(dictFilePath)
	if err != nil {
		return nil, err
	}

	return &Analyzer{classifier, tokenizer}, nil
}
