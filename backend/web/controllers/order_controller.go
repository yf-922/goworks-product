package controllers

import (
	"errors"
	"product/common"
	"product/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type OrderController struct {
	Ctx          iris.Context
	OrderService services.IOrderService
}

func (o *OrderController) Get() mvc.View {
	orderArray, err := o.OrderService.GetAllOrderInfo()
	if err != nil {
		o.Ctx.Application().Logger().Errorf(err.Error())
	}
	return mvc.View{
		Name: "order/view.html",
		Data: iris.Map{
			"order": orderArray,
		},
	}
}

func (o *OrderController) PostCreate() {
	order, err := common.BuildOrderForCreateFromContext(o.Ctx)
	if err != nil {
		o.Ctx.StatusCode(iris.StatusBadRequest)
		_ = o.Ctx.JSON(common.NewFailResult(err.Error()))
		return
	}
	orderID, err := o.OrderService.InsertOrder(order)
	if err != nil {
		if errors.Is(err, services.ErrOrderUserNotFound) {
			o.Ctx.StatusCode(iris.StatusBadRequest)
			_ = o.Ctx.JSON(common.NewFailResult("user not found"))
			return
		}
		if errors.Is(err, services.ErrOrderProductNotFound) {
			o.Ctx.StatusCode(iris.StatusBadRequest)
			_ = o.Ctx.JSON(common.NewFailResult("product not found"))
			return
		}
		o.Ctx.StatusCode(iris.StatusInternalServerError)
		_ = o.Ctx.JSON(common.NewFailResult("订单创建失败"))
		return
	}
	_ = o.Ctx.JSON(iris.Map{
		"orderID":     orderID,
		"userID":      order.UserID,
		"productId":   order.ProductId,
		"orderStatus": order.OrderStatus,
	})
}
