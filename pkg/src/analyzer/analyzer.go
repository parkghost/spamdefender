package analyzer

import (
	"bytes"
	"common"
	"fmt"
	"github.com/parkghost/bayesian"
	"goseg"
	"math"
	"sync"
	"time"
)

const (
	Good           = "Good"
	Bad            = "Bad"
	Neutral        = "Neutral"
	ClassIdxOfGood = 0
	ClassIdxOfBad  = 1
	cutset         = ":;=<>"
	Threshold      = 0.01
)

type Analyzer interface {
	Test(text string) string
}

type Learner interface {
	Learn(text, class string)
}

type BayesianAnalyzer struct {
	tokenizer  *goseg.Tokenizer
	classifier *bayesian.Classifier
	updater    Updater
	rwm        *sync.RWMutex
}

func (ba *BayesianAnalyzer) Test(text string) string {
	words := common.Normalize(ba.tokenizer.Cut([]rune(text)), cutset)

	ba.rwm.RLock()
	scores, likely, strict := ba.classifier.LogScores(words)
	ba.rwm.RUnlock()

	if !strict || math.Abs(scores[ClassIdxOfGood]/scores[ClassIdxOfBad]-1) < Threshold {
		return Neutral
	}

	return string(ba.classifier.Classes[likely])
}

func (ba *BayesianAnalyzer) Learn(text, class string) {
	words := common.Normalize(ba.tokenizer.Cut([]rune(text)), cutset)

	ba.rwm.Lock()
	ba.classifier.Learn(words, bayesian.Class(class))
	ba.rwm.Unlock()

	if ba.updater != nil {
		ba.updater.Update()
	}
}

func (ba *BayesianAnalyzer) Explain(text string) WordFreqList {
	words := common.Normalize(ba.tokenizer.Cut([]rune(text)), cutset)

	ba.rwm.RLock()
	freqMatrix := ba.classifier.WordFrequencies(words)
	ba.rwm.RUnlock()

	wordFreqs := make([]WordFreq, len(words))

	for i, _ := range words {
		wordFreq := make([]float64, len(ba.classifier.Classes))
		for j, _ := range ba.classifier.Classes {
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

func NewBayesianAnalyzer(traningDataFilePath string, dictDataFilePath string) (*BayesianAnalyzer, error) {
	tokenizer, err := goseg.NewTokenizerFromFile(dictDataFilePath)
	if err != nil {
		return nil, err
	}

	coordinator := &sync.RWMutex{}
	classifier, err := bayesian.NewClassifierFromFile(traningDataFilePath)
	if err != nil {
		return nil, err
	}

	return &BayesianAnalyzer{tokenizer, classifier, nil, coordinator}, nil
}

func NewBayesianAnalyzerWithUpdater(traningDataFilePath string, dictDataFilePath string, updateDelay time.Duration) (*BayesianAnalyzer, error) {
	tokenizer, err := goseg.NewTokenizerFromFile(dictDataFilePath)
	if err != nil {
		return nil, err
	}

	coordinator := &sync.RWMutex{}
	classifier, err := bayesian.NewClassifierFromFile(traningDataFilePath)
	if err != nil {
		return nil, err
	}

	updater := NewDelayedUpdater(classifier, traningDataFilePath, updateDelay, coordinator)

	return &BayesianAnalyzer{tokenizer, classifier, updater, coordinator}, nil
}
