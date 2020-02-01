package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	sconsts "github.com/galaxy-book/polaris-backend/service/platform/projectsvc/consts"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func GetIssueSourceById(orgId, sourceId int64) (*bo.IssueSourceBo, errs.SystemErrorInfo) {
	list, err1 := GetIssueSourceList(orgId)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}
	for _, source := range list {
		if source.Id == sourceId {
			return &source, nil
		}
	}
	return nil, errs.BuildSystemErrorInfo(errs.SourceNotExist)
}

func GetIssueSourceListByProjectObjectTypeId(orgId, projectTypeId int64) ([]bo.IssueSourceBo, errs.SystemErrorInfo) {
	list, err1 := GetIssueSourceList(orgId)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}
	result := make([]bo.IssueSourceBo, 0)
	for _, source := range list {
		if source.ProjectObjectTypeId == projectTypeId {
			result = append(result, source)
		}
	}
	return result, nil
}

func GetIssueSourceList(orgId int64) ([]bo.IssueSourceBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheIssueSourceList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	issueSourceListBo := &[]bo.IssueSourceBo{}
	issueSourceListPo := &[]po.PpmPrsIssueSource{}
	issueSourceListJson, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if issueSourceListJson != "" {

		err = json.FromJson(issueSourceListJson, issueSourceListBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return *issueSourceListBo, nil
	} else {
		err := mysql.SelectAllByCond(consts.TableIssueSource, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcStatus:   consts.AppStatusEnable,
		}, issueSourceListPo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		_ = copyer.Copy(issueSourceListPo, issueSourceListBo)
		issueSourceListJson, err := json.ToJson(issueSourceListBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, issueSourceListJson, consts.GetCacheBaseExpire())
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		return *issueSourceListBo, nil
	}
}
