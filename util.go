package beegoutil

import (
	"github.com/astaxie/beego"
)

// StatusMessageModel defines the StatusMessageModel body
type StatusMessageModel struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}

// StatusMessageModelInterface defines the StatusMessageModeInterface
type StatusMessageModelInterface interface {
	SetStatus(status int)
	SetMessage(message string)
}

// SetStatus sets the response status
func (r *StatusMessageModel) SetStatus(status int) {
	r.Status = status
}

// SetMessage sets the response error message
func (r *StatusMessageModel) SetMessage(message string) {
	r.Message = message
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
