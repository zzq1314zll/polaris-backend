package dao

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertIssueRelation(po po.PpmPriIssueRelation, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("IssueRelation dao Insert err %v", err)
	}
	return nil
}

func InsertIssueRelationBatch(pos []po.PpmPriIssueRelation, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
		conn, err := mysql.GetConnect()
		defer func() {
			if conn != nil {
				if err := conn.Close(); err != nil {
					logger.GetDefaultLogger().Info(strs.ObjectToString(err))
				}
			}
		}()
		if err != nil {
			return err
		}
		batch = conn.InsertInto(consts.TableIssueRelation).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableIssueRelation).Batch(len(pos))
	}
	go func() {
		defer batch.Done()
		for i := range pos {
			batch.Values(pos[i])
		}
	}()
	err = batch.Wait()
	if err != nil {
		log.Errorf("Iteration dao InsertBatch err %v", err)
		return err
	}
	return nil
}

func UpdateIssueRelation(po po.PpmPriIssueRelation, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("IssueRelation dao Update err %v", err)
	}
	return err
}

func UpdateIssueRelationById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateIssueRelationByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateIssueRelationByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateIssueRelationByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateIssueRelationByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableIssueRelation, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableIssueRelation, cond, upd)
	}
	if err != nil {
		log.Errorf("IssueRelation dao Update err %v", err)
	}
	return mod, err
}

func DeleteIssueRelationById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateIssueRelationByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteIssueRelationByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateIssueRelationByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectIssueRelationById(id int64) (*po.PpmPriIssueRelation, error) {
	po := &po.PpmPriIssueRelation{}
	err := mysql.SelectById(consts.TableIssueRelation, id, po)
	if err != nil {
		log.Errorf("IssueRelation dao SelectById err %v", err)
	}
	return po, err
}

func SelectIssueRelationByIdAndOrg(id int64, orgId int64) (*po.PpmPriIssueRelation, error) {
	po := &po.PpmPriIssueRelation{}
	err := mysql.SelectOneByCond(consts.TableIssueRelation, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("IssueRelation dao Select err %v", err)
	}
	return po, err
}

func SelectIssueRelation(cond db.Cond) (*[]po.PpmPriIssueRelation, error) {
	pos := &[]po.PpmPriIssueRelation{}
	err := mysql.SelectAllByCond(consts.TableIssueRelation, cond, pos)
	if err != nil {
		log.Errorf("IssueRelation dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneIssueRelation(cond db.Cond) (*po.PpmPriIssueRelation, error) {
	po := &po.PpmPriIssueRelation{}
	err := mysql.SelectOneByCond(consts.TableIssueRelation, cond, po)
	if err != nil {
		log.Errorf("IssueRelation dao Select err %v", err)
	}
	return po, err
}

func SelectIssueRelationByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmPriIssueRelation, uint64, error) {
	pos := &[]po.PpmPriIssueRelation{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableIssueRelation, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("IssueRelation dao SelectPage err %v", err)
	}
	return pos, total, err
}
