package controllers

import (
	"product/common"
	"product/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type UserController struct {
	Ctx         iris.Context
	UserService services.IUserService
}

func (u *UserController) Get() mvc.View {
	return mvc.View{
		Name: "user/view.html",
	}
}

func (u *UserController) PostLogin() {
	userName := u.Ctx.FormValue("userName")
	password := u.Ctx.FormValue("password")
	if userName == "" {
		u.Ctx.StatusCode(iris.StatusBadRequest)
		_ = u.Ctx.JSON(common.NewFailResult("userName is required"))
		return
	}
	if password == "" {
		u.Ctx.StatusCode(iris.StatusBadRequest)
		_ = u.Ctx.JSON(common.NewFailResult("password is required"))
		return
	}
	user, ok := u.UserService.IsPwdSuccess(userName, password)
	if !ok {
		u.Ctx.StatusCode(iris.StatusUnauthorized)

		_ = u.Ctx.JSON(common.NewFailResult("用户名或者密码错误"))
		return
	}
	_ = u.Ctx.JSON(common.NewSuccessResult("登陆成功", iris.Map{
		"id":       user.ID,
		"nickName": user.NickName,
		"userName": user.Username,
	}))

}

func (u *UserController) PostCreate() {
	user, err := common.BuildUserForCreateFromContext(u.Ctx)
	if err != nil {
		u.Ctx.StatusCode(iris.StatusBadRequest)

		_ = u.Ctx.JSON(common.NewFailResult(err.Error()))
		return
	}

	userID, err := u.UserService.AddUser(user)
	if err != nil {
		if err.Error() == "userName already exists" {
			u.Ctx.StatusCode(iris.StatusBadRequest)
			_ = u.Ctx.JSON(common.NewFailResult(err.Error()))
			return
		}
		u.Ctx.StatusCode(iris.StatusInternalServerError)
		_ = u.Ctx.JSON(common.NewFailResult("用户创建失败"))
		return
	}

	_ = u.Ctx.JSON(common.NewSuccessResult("用户创建成功", iris.Map{
		"id":       userID,
		"nickName": user.NickName,
		"userName": user.Username,
	}))

}

func (u *UserController) GetAll() {
	users, err := u.UserService.GetAllUsers()
	if err != nil {
		u.Ctx.StatusCode(iris.StatusInternalServerError)
		_ = u.Ctx.JSON(common.NewFailResult("查询失败"))
		return
	}
	_ = u.Ctx.JSON(common.NewSuccessResult("查询成功", users))
}

func (u *UserController) GetBy(id int64) {
	user, err := u.UserService.GetUserById(id)
	if err != nil {
		u.Ctx.StatusCode(404)
		_ = u.Ctx.JSON(common.NewFailResult("没有得到"))
		return
	}
	_ = u.Ctx.JSON(common.NewSuccessResult("成功得到", user))

}

func (u *UserController) PostDelete() {
	idString := u.Ctx.FormValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		u.Ctx.StatusCode(400)
		_ = u.Ctx.JSON(common.NewFailResult("用户输入不合法"))
		return
	}
	ok := u.UserService.DeleteUserById(id)
	if !ok {
		u.Ctx.StatusCode(iris.StatusInternalServerError)
		_ = u.Ctx.JSON(common.NewFailResult("删除失败"))
		return
	}
	_ = u.Ctx.JSON(common.NewSuccessResult("删除成功", nil))
}

func (u *UserController) PostUpdate() {
	user, err := common.BuildUserForUpdateFromContext(u.Ctx)
	if err != nil {
		u.Ctx.StatusCode(iris.StatusBadRequest)

		_ = u.Ctx.JSON(common.NewFailResult(err.Error()))
		return
	}
	err = u.UserService.UpdateUser(user)
	if err != nil {
		u.Ctx.StatusCode(iris.StatusInternalServerError)

		_ = u.Ctx.JSON(common.NewFailResult("更新失败"))
		return
	}

	_ = u.Ctx.JSON(common.NewSuccessResult("更新成功", nil))
}
