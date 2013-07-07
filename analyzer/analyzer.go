package analyzer

import (
	"bytes"
	"fmt"
	"github.com/parkghost/bayesian"
	"github.com/parkghost/spamdefender/analyzer/goseg"
	"github.com/parkghost/spamdefender/common"
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

func (a *Analyzer) Test(text string) (bayesian.Class, map[bayesian.Class]float64) {
	words := common.Normalize(a.tokenizer.Cut([]rune(text)), cutset)
	score, likely, _ := a.classifier.LogScores(words)

	mapping := make(map[bayesian.Class]float64)
	for i, class := range a.classifier.Classes {
		mapping[class] = score[i]
	}

	return a.classifier.Classes[likely], mapping
}

func (a *Analyzer) Explain(text string) WordFreqList {
	words := common.Normalize(a.tokenizer.Cut([]rune(text)), cutset)

	freqMatrix := a.classifier.WordFrequencies(words)
	wordFreqs := make([]WordFreq, len(words))

	for i, _ := range words {
		wordFreq := make([]float64, len(a.classifier.Classes))
		for j, _ := range a.classifier.Classes {
			wordFreq[j] = freqMatrix[j][i]
		}
		wordFreqs[i] = WordFreq{words[i], wordFreq}
	}

	return wordFreqs
}

type WordFreq struct {
	Word       string
	FreqMatrix []float64
}

func (wf *WordFreq) String() string {
	visible := false
	buf := bytes.NewBufferString(wf.Word)
	freqsText := bytes.NewBufferString("(")
	for c, v := range wf.FreqMatrix {
		if v > 0.0000001 {
			visible = true
		}

		if c == 0 {
			freqsText.WriteString(fmt.Sprintf("%f", v))
		} else {
			freqsText.WriteString(fmt.Sprintf(",%f", v))
		}
	}
	freqsText.WriteString(")")

	if visible {
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

func NewAnalyzer(traningDataFilePath string, dictFilePath ...string) (*Analyzer, error) {
	classifier, err := bayesian.NewClassifierFromFile(traningDataFilePath)
	if err != nil {
		return nil, err
	}

	tokenizer, err := goseg.NewTokenizer(dictFilePath...)
	if err != nil {
		return nil, err
	}
	return &Analyzer{classifier, tokenizer}, nil
}
