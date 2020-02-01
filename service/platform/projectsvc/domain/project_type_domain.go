package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func ProjectTypes(orgId int64) (*[]bo.ProjectTypeBo, errs.SystemErrorInfo) {
	projectType := &[]po.PpmPrsProjectType{}
	err := mysql.SelectAllByCond(consts.TableProjectType, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.AppIsInitStatus,
		consts.TcOrgId:    db.In([]int64{orgId, 0}),
	}, projectType)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	projectTypeBo := &[]bo.ProjectTypeBo{}
	copyErr := copyer.Copy(projectType, projectTypeBo)
	if copyErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, copyErr)
	}

	return projectTypeBo, nil
}

func ProjectType(orgId, id int64) (*bo.ProjectTypeBo, errs.SystemErrorInfo) {

	projectType := &po.PpmPrsProjectType{}

	err := mysql.SelectOneByCond(consts.TableProjectType, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.AppIsInitStatus,
		consts.TcId:       id,
		consts.TcOrgId:    db.In([]int64{orgId, 0}),
	}, projectType)

	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	projectTypeBo := &bo.ProjectTypeBo{}
	copyErr := copyer.Copy(projectType, projectTypeBo)
	if copyErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, copyErr)
	}

	return projectTypeBo, nil

}
