package api

import (
	"log"

	"github.com/mohuishou/scu/library"

	"github.com/mohuishou/scu"

	"github.com/kataras/iris"
	"github.com/mohuishou/scu/jwc"
	"github.com/mohuishou/scuplus-go/cache/bind"
	"github.com/mohuishou/scuplus-go/middleware"
	"github.com/mohuishou/scuplus-go/model"
)

// Login 用户登录
func Login(ctx iris.Context) {
	code := ctx.FormValue("code")

	if code == "" {
		Error(ctx, 10400, "code不能为空", nil)
		return
	}

	// 获取openid
	user := &model.User{}
	if err := user.Wechat.GetOpenid(code); err != nil {
		Error(ctx, 10401, "获取用户信息失败", nil)
		return
	}

	// 登录
	token, err := user.Login()
	if err != nil {
		log.Println(err)
		Error(ctx, 10401, "登录失败", nil)
		return
	}

	userLibrary := model.UserLibrary{}
	model.DB().Model(&user).Related(&userLibrary)
	userConf := model.UserConfig{}
	model.DB().Where("user_id = ?", user.ID).Find(&userConf)

	Success(ctx, "登录成功！", map[string]interface{}{
		"token":          token,
		"jwc_verify":     user.JwcVerify,
		"verify":         user.Verify,
		"library_verify": userLibrary.Verify,
		"user_type":      userConf.UserType,
	})
}

// Remove 注销账号，所有数据会被移除
func Remove(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	// 移除
	model.RemoveAll(uid)
	Success(ctx, "数据清除成功！如需继续使用，请退出小程序之后重新登录，再会", nil)
}

// BindJwc 教务处绑定
func BindJwc(ctx iris.Context) {
	studentID := ctx.FormValue("student_id")
	password := ctx.FormValue("password")

	if studentID == "" || password == "" {
		Error(ctx, 10400, "参数错误", nil)
		return
	}

	uid := middleware.GetUserID(ctx)
	if bind.Get(uid, "jwc") >= 3 {
		Error(ctx, 403, "今日绑定账号切换次数超过限制", nil)
		return
	}

	// 验证教务处账号是否可以登录
	if _, err := jwc.Login(studentID, password); err != nil {
		Error(ctx, 10401, err.Error(), nil)
		return
	}

	user := model.User{
		JwcStudentID: studentID,
		JwcPassword:  password,
		JwcVerify:    1,
	}

	if err := model.DB().Model(&user).Where("id = ?", uid).Updates(&user).Error; err != nil {
		log.Println("用户绑定账号失败: ", err)
		Error(ctx, 30004, "系统错误！", nil)
		return
	}

	// 绑定成功，异步任务获取数据信息
	//sign := &tasks.Signature{
	//	Name: "update_new",
	//	Args: []tasks.Arg{
	//		{
	//			Type:  "uint",
	//			Value: uid,
	//		},
	//	},
	//}
	//_, err := job.Server.SendTask(sign)
	//if err != nil {
	//	log.Println("cron error update all", err)
	//}
	// 删除已有信息
	model.AfterUpdateBindJwc(uid)
	// 添加绑定次数
	bind.Add(uid, "jwc")
	Success(ctx, "绑定成功！", nil)
}

// Bind 绑定统一认证账号
func Bind(ctx iris.Context) {
	studentID := ctx.FormValue("student_id")
	password := ctx.FormValue("password")

	if studentID == "" || password == "" {
		Error(ctx, 10400, "参数错误", nil)
		return
	}

	uid := middleware.GetUserID(ctx)
	if bind.Get(uid, "my") >= 3 {
		Error(ctx, 403, "今日绑定账号切换次数超过限制", nil)
		return
	}

	// 验证统一账号是否可以登录
	if _, err := scu.NewCollector(studentID, password); err != nil {
		Error(ctx, 10401, err.Error(), nil)
		return
	}

	user := model.User{
		StudentID: studentID,
		Password:  password,
	}

	user.Verify = 1
	if err := model.DB().Model(&user).Where("id = ?", uid).Updates(&user).Error; err != nil {
		log.Println("用户绑定账号失败: ", err)
		Error(ctx, 30004, "系统错误！", nil)
		return
	}
	// 删除已有账号信息
	model.AfterUpdateBindMy(uid)
	bind.Add(uid, "my")
	Success(ctx, "绑定成功！", nil)
}

// BindLibrary 绑定图书馆账号
func BindLibrary(ctx iris.Context) {
	studentID := ctx.FormValue("student_id")
	password := ctx.FormValue("password")

	if studentID == "" || password == "" {
		Error(ctx, 10400, "参数错误", nil)
		return
	}

	uid := middleware.GetUserID(ctx)
	if bind.Get(uid, "library") >= 3 {
		Error(ctx, 403, "今日绑定账号切换次数超过限制", nil)
		return
	}

	// 检查账号信息
	_, err := library.NewLibrary(studentID, password)
	if err != nil {
		Error(ctx, 10401, err.Error(), nil)
		return
	}

	// 保存图书馆账号
	userLib := model.UserLibrary{
		UserID:    uid,
		StudentID: studentID,
		Password:  password,
		Verify:    1,
	}
	var oldUserLib model.UserLibrary
	tx := model.DB().Begin()
	if tx.Where("user_id = ?", uid).Find(&oldUserLib).RecordNotFound() {
		if err := tx.Create(&userLib).Error; err != nil {
			Error(ctx, 10401, err.Error(), nil)
			tx.Rollback()
			return
		}
	} else {
		if err := tx.Model(&oldUserLib).Updates(userLib).Error; err != nil {
			Error(ctx, 10401, err.Error(), nil)
			tx.Rollback()
			return
		}
	}
	tx.Commit()

	// 删除已有账号数据
	model.AfterUpdateBindLibrary(uid)
	bind.Add(uid, "library")
	Success(ctx, "绑定成功！", nil)
}
