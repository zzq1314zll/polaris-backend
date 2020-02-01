package consts

import (
	"github.com/galaxy-book/polaris-backend/common/core/consts"
)

var (
	//对象id缓存key前缀
	CacheObjectIdPreKey = consts.CacheKeyPrefix + consts.IdsvcApplicationName + consts.CacheKeyOfSys + "object_id:"
)
