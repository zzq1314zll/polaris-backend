package domain

import (
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"strconv"
)

//添加项目趋势
func AddProjectTrends(projectTrendsBo bo.ProjectTrendsBo) {
	trendsType := projectTrendsBo.PushType
	var trendsBos []bo.TrendsBo = nil
	var err errs.SystemErrorInfo = nil

	if trendsType == consts.PushTypeCreateProject {
		trendsBos = assemblyCreateProjectTrends(projectTrendsBo)
	} else if trendsType == consts.PushTypeUpdateProject {
		trendsBos = assemblyUpdateProjectTrends(projectTrendsBo)
	} else if trendsType == consts.PushTypeUpdateProjectMembers {
		trendsBos, err = assemblyUpdateProjectMemberTrends(projectTrendsBo)
		log.Error(err)
	} else if trendsType == consts.PushTypeStarProject {
		trendsBos = assemblyStarProjectTrends(projectTrendsBo)
	} else if trendsType == consts.PushTypeUnstarProject {
		trendsBos = assemblyUnstarProjectTrends(projectTrendsBo)
	} else if trendsType == consts.PushTypeUnbindProject {
		trendsBos = assemblyUnbindProjectTrends(projectTrendsBo)
	} else if trendsType == consts.PushTypeUpdateProjectStatus {
		trendsBos = assemblyUpdateProjectStatusTrends(projectTrendsBo)
	} else if trendsType == consts.PushTypeCreateIssueBatch {
		trendsBos = assemblyCreateIssueBatchTrends(projectTrendsBo)
	} else if trendsType == consts.PushTypeDeleteResource {
		trendsBos = assemblyDeleteResourceTrends(projectTrendsBo)
	} else if trendsType == consts.PushTypeCreateProjectFile {
		trendsBos = assemblyCreateProjectFileTrends(projectTrendsBo)
	} else if trendsType == consts.PushTypeUpdateProjectFile {
		trendsBos = assemblyUpdateProjectFileTrends(projectTrendsBo)
	} else if trendsType == consts.PushTypeDeleteProjectFile {
		trendsBos = assemblyDeleteProjectFileTrends(projectTrendsBo)
	} else if trendsType == consts.PushTypeCreateProjectFolder {
		trendsBos = assemblyCreateProjectFolderTrends(projectTrendsBo)
	} else if trendsType == consts.PushTypeUpdateProjectFolder {
		trendsBos = assemblyUpdateProjectFolderTrends(projectTrendsBo)
	} else if trendsType == consts.PushTypeDeleteProjectFolder {
		trendsBos = assemblyDeleteProjectFolderTrends(projectTrendsBo)
	}
	//插入操作
	for i, _ := range trendsBos {
		trendsBos[i].CreateTime = types.Time(projectTrendsBo.OperateTime)
	}

	err = CreateTrendsBatch(trendsBos)
	if err != nil {
		log.Error(err)
	}
}

//创建项目
func assemblyCreateProjectTrends(projectTrendsBo bo.ProjectTrendsBo) []bo.TrendsBo {
	trendsBos := []bo.TrendsBo{}

	newValue := projectTrendsBo.NewValue
	trendsBo := bo.TrendsBo{
		OrgId:           projectTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       projectTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       0,
		Module3:         consts.BlankString,
		OperCode:        consts.RoleOperationCreate,
		OperObjId:       projectTrendsBo.ProjectId,
		OperObjType:     consts.TrendsOperObjectTypeProject,
		OperObjProperty: consts.BlankString,
		RelationObjId:   0,
		RelationObjType: consts.BlankString,
		RelationType:    consts.TrendsRelationTypeCreateProject,
		NewValue:        &newValue,
		OldValue:        nil,
		Ext:             json.ToJsonIgnoreError(projectTrendsBo.Ext),
		Creator:         projectTrendsBo.OperatorId,
		CreateTime:      types.NowTime(),
	}

	trendsBos = append(trendsBos, trendsBo)

	return trendsBos
}

//更新项目
func assemblyUpdateProjectTrends(projectTrendsBo bo.ProjectTrendsBo) []bo.TrendsBo {
	newValue := json.ToJsonIgnoreError(projectTrendsBo.NewValue)
	oldValue := json.ToJsonIgnoreError(projectTrendsBo.OldValue)
	operObjProperty := projectTrendsBo.OperateObjProperty
	trendsBo := &bo.TrendsBo{
		OrgId:           projectTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       projectTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       0,
		Module3:         consts.BlankString,
		OperCode:        consts.RoleOperationModify,
		OperObjId:       projectTrendsBo.ProjectId,
		OperObjType:     consts.TrendsOperObjectTypeProject,
		OperObjProperty: operObjProperty,
		RelationObjId:   0,
		RelationObjType: consts.BlankString,
		RelationType:    consts.TrendsRelationTypeUpdateProject,
		NewValue:        &newValue,
		OldValue:        &oldValue,
		Ext:             json.ToJsonIgnoreError(projectTrendsBo.Ext),
		Creator:         projectTrendsBo.OperatorId,
	}
	for _, v := range projectTrendsBo.Ext.ChangeList {
		if v.Field == "notice" {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						log.Errorf("捕获到的错误：%s", r)
					}
				}()
				e := AddNoticeByChangeProjectNotice(projectTrendsBo)
				if e != nil{
					log.Error(e)
				}
			}()
			break
		}
	}

	return []bo.TrendsBo{*trendsBo}
}

//关注项目
func assemblyStarProjectTrends(projectTrendsBo bo.ProjectTrendsBo) []bo.TrendsBo {
	trendsBos := []bo.TrendsBo{}

	trendsBo := bo.TrendsBo{
		OrgId:           projectTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       projectTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       0,
		Module3:         consts.BlankString,
		OperCode:        consts.RoleOperationAttention,
		OperObjId:       projectTrendsBo.ProjectId,
		OperObjType:     consts.TrendsOperObjectTypeProject,
		OperObjProperty: consts.BlankString,
		RelationObjId:   projectTrendsBo.OperatorId,
		RelationObjType: consts.TrendsRelationTypeUser,
		RelationType:    consts.TrendsRelationTypeStarProject,
		NewValue:        nil,
		OldValue:        nil,
		Ext:             json.ToJsonIgnoreError(projectTrendsBo.Ext),
		Creator:         projectTrendsBo.OperatorId,
		CreateTime:      types.NowTime(),
	}

	trendsBos = append(trendsBos, trendsBo)

	return trendsBos
}

//取关项目
func assemblyUnstarProjectTrends(projectTrendsBo bo.ProjectTrendsBo) []bo.TrendsBo {
	trendsBos := []bo.TrendsBo{}

	trendsBo := bo.TrendsBo{
		OrgId:           projectTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       projectTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       0,
		Module3:         consts.BlankString,
		OperCode:        consts.RoleOperationUnAttention,
		OperObjId:       projectTrendsBo.ProjectId,
		OperObjType:     consts.TrendsOperObjectTypeProject,
		OperObjProperty: consts.BlankString,
		RelationObjId:   projectTrendsBo.OperatorId,
		RelationObjType: consts.TrendsRelationTypeUser,
		RelationType:    consts.TrendsRelationTypeUnstarProject,
		NewValue:        nil,
		OldValue:        nil,
		Ext:             json.ToJsonIgnoreError(projectTrendsBo.Ext),
		Creator:         projectTrendsBo.OperatorId,
		CreateTime:      types.NowTime(),
	}

	trendsBos = append(trendsBos, trendsBo)

	return trendsBos
}

//退出项目
func assemblyUnbindProjectTrends(projectTrendsBo bo.ProjectTrendsBo) []bo.TrendsBo {
	trendsBos := []bo.TrendsBo{}

	trendsBo := bo.TrendsBo{
		OrgId:           projectTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       projectTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       projectTrendsBo.OperatorId, //todo
		Module3:         consts.TrendsModuleRole,
		OperCode:        consts.RoleOperationUnbind,
		OperObjId:       projectTrendsBo.ProjectId, //todo
		OperObjType:     consts.TrendsOperObjectTypeRole,
		OperObjProperty: consts.BlankString,
		RelationObjId:   projectTrendsBo.OperatorId,
		RelationObjType: consts.TrendsRelationTypeUser,
		RelationType:    consts.TrendsRelationTypeUnbindProject,
		NewValue:        nil,
		OldValue:        nil,
		Ext:             json.ToJsonIgnoreError(projectTrendsBo.Ext),
		Creator:         projectTrendsBo.OperatorId,
		CreateTime:      types.NowTime(),
	}

	trendsBos = append(trendsBos, trendsBo)

	return trendsBos
}

//更新成员
func assemblyUpdateProjectMemberTrends(projectTrendsBo bo.ProjectTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	deletedFollowerIds, addedFollowerIds := util.GetDifMemberIds(projectTrendsBo.BeforeChangeFollowers, projectTrendsBo.AfterChangeFollowers)
	deletedParticipantIds, addedParticipantIds := util.GetDifMemberIds(projectTrendsBo.BeforeChangeMembers, projectTrendsBo.AfterChangeMembers)

	memberChangeInfos := [][]int64{deletedFollowerIds, addedFollowerIds, deletedParticipantIds, addedParticipantIds}
	trendsBos := &[]bo.TrendsBo{}
	//获取所有相关用户
	allRelationIds := slice.SliceUniqueInt64(append(append(append(append(append(deletedFollowerIds, addedFollowerIds...), deletedParticipantIds...), addedParticipantIds...), projectTrendsBo.AfterOwner), projectTrendsBo.BeforeOwner))
	userInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(consts.AppSourceChannelDingTalk, projectTrendsBo.OrgId, allRelationIds)
	if err != nil {
		log.Error("获取用户相关信息失败" + strs.ObjectToString(err))
		return nil, err
	}
	userInfo := map[int64]bo.SimpleUserInfoBo{}
	for _, v := range userInfos {
		userInfo[v.UserId] = bo.SimpleUserInfoBo{
			Id:     v.UserId,
			Name:   v.Name,
			Avatar: v.Avatar,
		}
	}

	for i, changeInfos := range memberChangeInfos {
		if len(changeInfos) > 0 {
			ext := projectTrendsBo.Ext
			for _, v := range changeInfos {
				ext.MemberInfo = append(ext.MemberInfo, userInfo[v])
			}
			beforeMap := map[string]interface{}{}
			afterMap := map[string]interface{}{}
			operCode := ""
			relationType := ""
			operObjProperty := ""

			handleTrendType(i, projectTrendsBo, &operObjProperty, &operCode, &relationType, &beforeMap, &afterMap)
			newValue := json.ToJsonIgnoreError(afterMap)
			oldValue := json.ToJsonIgnoreError(beforeMap)

			trendsBo := bo.TrendsBo{
				OrgId:           projectTrendsBo.OrgId,
				Module1:         consts.TrendsModuleOrg,
				Module2Id:       projectTrendsBo.ProjectId,
				Module2:         consts.TrendsModuleProject,
				Module3Id:       0,
				Module3:         consts.BlankString,
				OperCode:        operCode,
				OperObjId:       0,
				OperObjType:     consts.TrendsOperObjectTypeProject,
				OperObjProperty: operObjProperty,
				RelationObjId:   0,
				RelationObjType: consts.TrendsRelationTypeUser,
				RelationType:    relationType,
				NewValue:        &newValue,
				OldValue:        &oldValue,
				Ext:             json.ToJsonIgnoreError(ext),
				Creator:         projectTrendsBo.OperatorId,
				CreateTime:      types.NowTime(),
			}
			*trendsBos = append(*trendsBos, trendsBo)
		}
	}

	if projectTrendsBo.BeforeOwner != projectTrendsBo.AfterOwner {
		//负责人动态
		ext := projectTrendsBo.Ext

		oldValue := strconv.FormatInt(projectTrendsBo.BeforeOwner, 10)
		newValue := strconv.FormatInt(projectTrendsBo.AfterOwner, 10)
		if _, ok := userInfo[projectTrendsBo.BeforeOwner]; ok {
			ext.MemberInfo = append(ext.MemberInfo, userInfo[projectTrendsBo.BeforeOwner])
		}
		if _, ok := userInfo[projectTrendsBo.AfterOwner]; ok {
			ext.MemberInfo = append(ext.MemberInfo, userInfo[projectTrendsBo.AfterOwner])
		}

		*trendsBos = append(*trendsBos, bo.TrendsBo{
			OrgId:           projectTrendsBo.OrgId,
			Module1:         consts.TrendsModuleOrg,
			Module2Id:       projectTrendsBo.ProjectId,
			Module2:         consts.TrendsModuleProject,
			Module3Id:       0,
			Module3:         consts.BlankString,
			OperCode:        consts.RoleOperationModify,
			OperObjId:       projectTrendsBo.ProjectId,
			OperObjType:     consts.TrendsOperObjectTypeProject,
			OperObjProperty: consts.TrendsOperObjPropertyNameOwner,
			RelationObjId:   projectTrendsBo.AfterOwner,
			RelationObjType: consts.TrendsRelationTypeUser,
			RelationType:    consts.TrendsRelationTypeUpdateProjectOwner,
			NewValue:        &newValue,
			OldValue:        &oldValue,
			Ext:             json.ToJsonIgnoreError(ext),
			Creator:         projectTrendsBo.OperatorId,
			CreateTime:      types.NowTime(),
		})
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("捕获到的错误：%s", r)
			}
		}()
		noticeErr := AddNoticeByChangeProjectMember(deletedParticipantIds, addedParticipantIds, projectTrendsBo)
		if noticeErr != nil {
			log.Error(noticeErr)
		}
	}()

	return *trendsBos, nil
}

func handleTrendType(i int, projectTrendsBo bo.ProjectTrendsBo, operObjProperty, operCode, relationType *string, beforeMap, afterMap *map[string]interface{}) {
	if i == 0 || i == 1 {
		*operObjProperty = consts.TrendsOperObjPropertyNameFollower
		(*beforeMap)[*operObjProperty] = projectTrendsBo.BeforeChangeFollowers
		(*afterMap)[*operObjProperty] = projectTrendsBo.AfterChangeFollowers
		if i == 0 {
			*operCode = consts.RoleOperationUnbind
			*relationType = consts.TrendsRelationTypeDeleteProjectFollower
		} else if i == 1 {
			*operCode = consts.RoleOperationBind
			*relationType = consts.TrendsRelationTypeAddedProjectFollower
		}
	} else if i == 2 || i == 3 {
		*operObjProperty = consts.TrendsOperObjPropertyNameParticipant
		(*beforeMap)[*operObjProperty] = projectTrendsBo.BeforeChangeMembers
		(*afterMap)[*operObjProperty] = projectTrendsBo.AfterChangeMembers
		if i == 2 {
			*operCode = consts.RoleOperationUnbind
			*relationType = consts.TrendsRelationTypeDeletedProjectParticipant
		} else if i == 3 {
			*operCode = consts.RoleOperationBind
			*relationType = consts.TrendsRelationTypeAddedProjectParticipant
		}
	}

}

//更新项目状态
func assemblyUpdateProjectStatusTrends(projectTrendsBo bo.ProjectTrendsBo) []bo.TrendsBo {
	trendsBos := []bo.TrendsBo{}

	trendsBo := bo.TrendsBo{
		OrgId:           projectTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       projectTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       0,
		Module3:         consts.BlankString,
		OperCode:        consts.RoleOperationModifyStatus,
		OperObjId:       projectTrendsBo.ProjectId,
		OperObjType:     consts.TrendsOperObjectTypeProject,
		OperObjProperty: consts.TrendsOperObjPropertyNameStatus,
		RelationObjId:   0,
		RelationObjType: consts.BlankString,
		RelationType:    consts.TrendsRelationTypeUpdateProjectStatus,
		NewValue:        &projectTrendsBo.NewValue,
		OldValue:        &projectTrendsBo.OldValue,
		Ext:             json.ToJsonIgnoreError(projectTrendsBo.Ext),
		Creator:         projectTrendsBo.OperatorId,
		CreateTime:      types.NowTime(),
	}

	trendsBos = append(trendsBos, trendsBo)

	return trendsBos
}

//批量插入
func assemblyCreateIssueBatchTrends(projectTrendsBo bo.ProjectTrendsBo) []bo.TrendsBo {
	trendsBos := []bo.TrendsBo{}

	trendsBo := bo.TrendsBo{
		OrgId:           projectTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       projectTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       0,
		Module3:         consts.BlankString,
		OperCode:        consts.RoleOperationCreate,
		OperObjId:       projectTrendsBo.ProjectId,
		OperObjType:     consts.TrendsOperObjectTypeProject,
		OperObjProperty: consts.TrendsOperObjPropertyNameStatus,
		RelationObjId:   0,
		RelationObjType: consts.BlankString,
		RelationType:    consts.TrendsRelationTypeCreateIssueBatch,
		NewValue:        &projectTrendsBo.NewValue,
		OldValue:        &projectTrendsBo.OldValue,
		Ext:             json.ToJsonIgnoreError(projectTrendsBo.Ext),
		Creator:         projectTrendsBo.OperatorId,
		CreateTime:      types.NowTime(),
	}

	trendsBos = append(trendsBos, trendsBo)

	return trendsBos
}

//删除附件
func assemblyDeleteResourceTrends(projectTrendsBo bo.ProjectTrendsBo) []bo.TrendsBo {
	trendsBos := []bo.TrendsBo{}

	trendsBo := bo.TrendsBo{
		OrgId:           projectTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       projectTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       0,
		Module3:         consts.BlankString,
		OperCode:        consts.RoleOperationDelete,
		OperObjId:       projectTrendsBo.ProjectId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: consts.BlankString,
		RelationObjId:   0,
		RelationObjType: consts.BlankString,
		RelationType:    consts.TrendsRelationTypeDeleteResource,
		NewValue:        &projectTrendsBo.NewValue,
		OldValue:        &projectTrendsBo.OldValue,
		Ext:             json.ToJsonIgnoreError(projectTrendsBo.Ext),
		Creator:         projectTrendsBo.OperatorId,
		CreateTime:      types.NowTime(),
	}

	trendsBos = append(trendsBos, trendsBo)

	return trendsBos
}

func assemblyCreateProjectFileTrends(projectTrendsBo bo.ProjectTrendsBo) []bo.TrendsBo {
	//上传文件
	trendsBos := []bo.TrendsBo{}

	trendsBo := bo.TrendsBo{
		OrgId:           projectTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       projectTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       0,
		Module3:         consts.BlankString,
		OperCode:        consts.RoleOperationCreate,
		OperObjId:       projectTrendsBo.ProjectId,
		OperObjType:     consts.TrendsOperObjectTypeProject,
		OperObjProperty: consts.BlankString,
		RelationObjId:   0,
		RelationObjType: consts.BlankString,
		RelationType:    consts.TrendsRelationTypeUploadProjectFile,
		NewValue:        &projectTrendsBo.NewValue,
		OldValue:        &projectTrendsBo.OldValue,
		Ext:             json.ToJsonIgnoreError(projectTrendsBo.Ext),
		Creator:         projectTrendsBo.OperatorId,
		CreateTime:      types.NowTime(),
	}

	trendsBos = append(trendsBos, trendsBo)

	return trendsBos
}

func assemblyUpdateProjectFileTrends(projectTrendsBo bo.ProjectTrendsBo) []bo.TrendsBo {
	//更新文件
	trendsBos := []bo.TrendsBo{}

	trendsBo := bo.TrendsBo{
		OrgId:           projectTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       projectTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       0,
		Module3:         consts.BlankString,
		OperCode:        consts.RoleOperationModify,
		OperObjId:       projectTrendsBo.ProjectId,
		OperObjType:     consts.TrendsOperObjectTypeProject,
		OperObjProperty: consts.BlankString,
		RelationObjId:   0,
		RelationObjType: consts.BlankString,
		RelationType:    consts.TrendsRelationTypeUpDateProjectFile,
		NewValue:        &projectTrendsBo.NewValue,
		OldValue:        &projectTrendsBo.OldValue,
		Ext:             json.ToJsonIgnoreError(projectTrendsBo.Ext),
		Creator:         projectTrendsBo.OperatorId,
		CreateTime:      types.NowTime(),
	}

	trendsBos = append(trendsBos, trendsBo)

	return trendsBos
}

func assemblyDeleteProjectFileTrends(projectTrendsBo bo.ProjectTrendsBo) []bo.TrendsBo {
	//删除文件
	trendsBos := []bo.TrendsBo{}

	trendsBo := bo.TrendsBo{
		OrgId:           projectTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       projectTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       0,
		Module3:         consts.BlankString,
		OperCode:        consts.RoleOperationDelete,
		OperObjId:       projectTrendsBo.ProjectId,
		OperObjType:     consts.TrendsOperObjectTypeProject,
		OperObjProperty: consts.BlankString,
		RelationObjId:   0,
		RelationObjType: consts.BlankString,
		RelationType:    consts.TrendsRelationTypeDeleteProjectFile,
		NewValue:        &projectTrendsBo.NewValue,
		OldValue:        &projectTrendsBo.OldValue,
		Ext:             json.ToJsonIgnoreError(projectTrendsBo.Ext),
		Creator:         projectTrendsBo.OperatorId,
		CreateTime:      types.NowTime(),
	}

	trendsBos = append(trendsBos, trendsBo)

	return trendsBos
}

func assemblyCreateProjectFolderTrends(projectTrendsBo bo.ProjectTrendsBo) []bo.TrendsBo {
	//创建文件夹
	trendsBos := []bo.TrendsBo{}

	trendsBo := bo.TrendsBo{
		OrgId:           projectTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       projectTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       0,
		Module3:         consts.BlankString,
		OperCode:        consts.RoleOperationCreate,
		OperObjId:       projectTrendsBo.ProjectId,
		OperObjType:     consts.TrendsOperObjectTypeProject,
		OperObjProperty: consts.BlankString,
		RelationObjId:   0,
		RelationObjType: consts.BlankString,
		RelationType:    consts.TrendsRelationTypeCreateProjectFolder,
		NewValue:        &projectTrendsBo.NewValue,
		OldValue:        &projectTrendsBo.OldValue,
		Ext:             json.ToJsonIgnoreError(projectTrendsBo.Ext),
		Creator:         projectTrendsBo.OperatorId,
		CreateTime:      types.NowTime(),
	}

	trendsBos = append(trendsBos, trendsBo)

	return trendsBos
}

func assemblyUpdateProjectFolderTrends(projectTrendsBo bo.ProjectTrendsBo) []bo.TrendsBo {
	//更新文件夹
	trendsBos := []bo.TrendsBo{}

	trendsBo := bo.TrendsBo{
		OrgId:           projectTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       projectTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       0,
		Module3:         consts.BlankString,
		OperCode:        consts.RoleOperationModify,
		OperObjId:       projectTrendsBo.ProjectId,
		OperObjType:     consts.TrendsOperObjectTypeProject,
		OperObjProperty: consts.BlankString,
		RelationObjId:   0,
		RelationObjType: consts.BlankString,
		RelationType:    consts.TrendsRelationTypeUpdateProjectFolder,
		NewValue:        &projectTrendsBo.NewValue,
		OldValue:        &projectTrendsBo.OldValue,
		Ext:             json.ToJsonIgnoreError(projectTrendsBo.Ext),
		Creator:         projectTrendsBo.OperatorId,
		CreateTime:      types.NowTime(),
	}

	trendsBos = append(trendsBos, trendsBo)

	return trendsBos
}

func assemblyDeleteProjectFolderTrends(projectTrendsBo bo.ProjectTrendsBo) []bo.TrendsBo {
	//删除文件夹
	trendsBos := []bo.TrendsBo{}

	trendsBo := bo.TrendsBo{
		OrgId:           projectTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       projectTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       0,
		Module3:         consts.BlankString,
		OperCode:        consts.RoleOperationDelete,
		OperObjId:       projectTrendsBo.ProjectId,
		OperObjType:     consts.TrendsOperObjectTypeProject,
		OperObjProperty: consts.BlankString,
		RelationObjId:   0,
		RelationObjType: consts.BlankString,
		RelationType:    consts.TrendsRelationTypeDeleteProjectFolder,
		NewValue:        &projectTrendsBo.NewValue,
		OldValue:        &projectTrendsBo.OldValue,
		Ext:             json.ToJsonIgnoreError(projectTrendsBo.Ext),
		Creator:         projectTrendsBo.OperatorId,
		CreateTime:      types.NowTime(),
	}

	trendsBos = append(trendsBos, trendsBo)

	return trendsBos
}
