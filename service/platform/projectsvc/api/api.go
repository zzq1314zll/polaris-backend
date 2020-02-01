package api

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/polaris-backend/common/extra/gin/mvc"
)

var log = logger.GetDefaultLogger()

var postGreeter = PostGreeter{}

var getGreeter = GetGreeter{}

type PostGreeter struct {
	mvc.Greeter
}

type GetGreeter struct {
	mvc.Greeter
}

func (GetGreeter) Health() string {
	return "ok"
}
