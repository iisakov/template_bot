package model

import (
	"fmt"
	"strings"
)

type Question struct {
	Text    string   `json:"text"`
	Answers []string `json:"answers"`
	Options []string `json:"options"`
}

func (q Question) String() string {
	return fmt.Sprintf(
		"text: %s, answers: %s, options: %s",
		q.Text,
		strings.Join(q.Answers, ", "),
		strings.Join(q.Options, ", "),
	)

}

type Questions struct {
	Questions          []Question `json:"questions"`
	CurrentQuestionNum int        `json:"currentQuestionNum"`
}

func (qs Questions) unpack() (result []string) {
	for _, str := range qs.Questions {
		result = append(result, str.String())
	}
	return
}

func (qs Questions) String() string {
	return fmt.Sprintf(
		"questions:\n%s,\ncurrentQuestionNum: %d",
		strings.Join(qs.unpack(), "\n"),
		qs.CurrentQuestionNum,
	)
}

func (qs Questions) GetQuestion() Question {
	return qs.Questions[qs.CurrentQuestionNum]
}

func (qs *Questions) Back() bool {
	if qs.CurrentQuestionNum <= 0 {
		return false
	} else {
		qs.CurrentQuestionNum -= 1
		return true
	}

}

func (qs *Questions) Next() bool {
	if qs.CurrentQuestionNum < len(qs.Questions)-1 {
		qs.CurrentQuestionNum += 1
	} else {
		return false
	}
	return true
}

func (qs Questions) FindQuestionIndex(text string) (int, bool) {
	for i, q := range qs.Questions {
		if q.Text == text {
			return i, true
		}
	}
	return -1, false
}

func (qs Questions) FindQuestion(text string) *Question {
	for _, q := range qs.Questions {
		if q.Text == text {
			return &q
		}
	}
	return nil
}

func (qs Questions) GetCurrentAnswers() ([]string, bool) {
	return qs.Questions[qs.CurrentQuestionNum].Answers, true
}

func (qs Questions) GetAnswers(text string) ([]string, bool) {
	i, ok := qs.FindQuestionIndex(text)
	if !ok {
		return nil, false
	}
	return qs.Questions[i].Answers, true
}

func (qs Questions) CreateBackup() {}
func (qs Questions) ReadBackup()   {}
