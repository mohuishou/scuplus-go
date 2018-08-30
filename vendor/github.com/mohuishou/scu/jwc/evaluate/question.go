package evaluate

type Question struct {
	ID struct {
		Code string `json:"questionsCode"`
	} `json:"id"`
}

type Questions []Question

var star2answerTeacher = map[float64]string{
	5: "10_1",
	4: "10_0.8",
	3: "10_0.7",
	2: "10_0.6",
	1: "10_0.2",
}

var star2answerZJ = map[float64]string{
	5: "10_1",
	4: "10_0.8",
	3: "10_0.6",
	2: "10_0.3",
	1: "10_0",
}

func (questions Questions) params(star float64, teacherType int) map[string]string {
	param := map[string]string{}
	for _, q := range questions {
		param[q.ID.Code] = star2answerTeacher[star]
		if teacherType == 1 {
			param[q.ID.Code] = star2answerZJ[star]
		}
	}
	return param
}
