package api

import "github.com/kataras/iris"
import "github.com/mohuishou/scuplus-go/util/classroom"

// GetClassroom 获取自习教室
func GetClassroom(ctx iris.Context) {
	room := ctx.FormValue("classroom")
	res, err := classroom.Get(room)
	if err != nil {
		Error(ctx, 20003, "获取教室信息错误！", err.Error())
		return
	}
	Success(ctx, "获取成功！", string(res))
}
