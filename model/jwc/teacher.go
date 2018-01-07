package jwc

import "github.com/mohuishou/scuplus-go/model"

// Teacher 教师表
type Teacher struct {
	model.Model
	Name    string
	College string
}
