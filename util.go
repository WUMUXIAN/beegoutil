package beegoutil

import (
	"github.com/astaxie/beego"
)

// StatusMessageModelInterface defines the StatusMessageModeInterface
type StatusMessageModelInterface interface {
	SetStatus(status int)
	SetMessage(message string)
}

// ResponseIfError responds the client with given status and message if error is encountered
func ResponseIfError(err error, controller beego.Controller, response StatusMessageModelInterface, status, httpStatus int) bool {
	if err == nil {
		return false
	}
	response.SetStatus(status)
	response.SetMessage(err.Error())
	controller.Ctx.Output.Status = httpStatus
	controller.Data["json"] = response
	controller.ServeJSON()
	return true
}

// ResponseError responds the client with given status and message for error.
func ResponseError(errMsg string, controller beego.Controller, response StatusMessageModelInterface, status, httpStatus int) {
	response.SetStatus(status)
	response.SetMessage(errMsg)
	controller.Ctx.Output.Status = httpStatus
	controller.Data["json"] = response
	controller.ServeJSON()
}
