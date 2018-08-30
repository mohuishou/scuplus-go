package evaluate

type Answer struct {
	AnswerStr string `json:"answer"`
}

type Answers []Answer

var answer2starTeacher = map[string]float64{
	"10_1":   5,
	"10_0.8": 4,
	"10_0.7": 3,
	"10_0.6": 2,
	"10_0.2": 1,
}

var answer2starZJ = map[string]float64{
	"10_1":   5,
	"10_0.8": 4,
	"10_0.6": 3,
	"10_0.3": 2,
	"10_0":   1,
}

func (answers Answers) average(teacherType int) float64 {
	var sum float64
	for _, v := range answers {
		if teacherType == 1 {
			sum = sum + answer2starZJ[v.AnswerStr]
		} else {
			sum = sum + answer2starTeacher[v.AnswerStr]
		}
	}
	return sum / float64(len(answers))
}
