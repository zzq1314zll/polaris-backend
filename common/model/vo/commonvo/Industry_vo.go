package commonvo

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
)

type IndustryListRespVo struct {
	vo.Err
	IndustryList *vo.IndustryListResp `json:"data"`
}
