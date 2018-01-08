package jwc

import "github.com/kataras/iris"

// GetGrade 获取成绩
func GetGrade(ctx iris.Context) {
	uid := uint(ctx.Values().Get("user_id").(float64))

	ctx.JSON(map[string]interface{}{
		"status": 0,
		"uid":    uid,
	})
}

// UpdateGrade 更新成绩
func UpdateGrade(ctx iris.Context) {

}
