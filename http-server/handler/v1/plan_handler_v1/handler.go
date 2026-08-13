package plan_handler_v1

import "paylist.server/handler"

type Handler struct {
	*handler.BaseHandler
}

func NewHandler(base *handler.BaseHandler) *Handler {
	return &Handler{BaseHandler: base}
}
