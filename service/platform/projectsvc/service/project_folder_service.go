package service

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
)

func CreateProjectFolder(orgId, operatorId int64, input vo.CreateProjectFolderReq) (*vo.Void, errs.SystemErrorInfo) {
	//此处的path需要更改
	err := domain.AuthProject(orgId, operatorId, input.ProjectID, consts.RoleOperationPathOrgProFile, consts.RoleOperationCreateFolder)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resourceId, err := domain.CreateProjectFolder(bo.CreateFolderBo{
		ProjectId: input.ProjectID,
		OrgId:     orgId,
		Name:      input.Name,
		ParentId:  input.ParentID,
		FileType:  input.FileType,
		UserId:    operatorId,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &vo.Void{
		ID: resourceId,
	}, nil
}

func UpdateProjectFolder(orgId, operatorId int64, input vo.UpdateProjectFolderReq) (*vo.UpdateProjectFolderResp, errs.SystemErrorInfo) {
	//此处的path需要更改
	err := domain.AuthProject(orgId, operatorId, input.ProjectID, consts.RoleOperationPathOrgProFile, consts.RoleOperationModifyFolder)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	updateBo := &bo.UpdateFolderBo{}
	err1 := copyer.Copy(&input, updateBo)
	if err1 != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
	}
	updateBo.UserId = operatorId
	updateBo.OrgId = orgId
	folderId, err := domain.UpdateProjectFolder(*updateBo)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &vo.UpdateProjectFolderResp{
		FolderID: *folderId,
	}, nil
}

func DeleteProjectFolder(orgId, operatorId int64, input vo.DeleteProjectFolderReq) (*vo.DeleteProjectFolderResp, errs.SystemErrorInfo) {
	//此处的path需要更改
	err := domain.AuthProject(orgId, operatorId, input.ProjectID, consts.RoleOperationPathOrgProFile, consts.RoleOperationDeleteFolder)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	folderIds, err := domain.DeleteProjectFolder(bo.DeleteFolderBo{
		FolderIds: input.FolderIds,
		OrgId:     orgId,
		UserId:    operatorId,
		ProjectId: input.ProjectID,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &vo.DeleteProjectFolderResp{
		FolderIds: folderIds,
	}, nil
}

func GetProjectFolder(orgId, operatorId int64, page, size int, input vo.ProjectFolderReq) (*vo.FolderList, errs.SystemErrorInfo) {
	//此处的path需要更改
	err := domain.AuthProjectWithOutPermission(orgId, operatorId, input.ProjectID, consts.RoleOperationPathOrgProFile, consts.RoleOperationView)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	folderList, err := domain.GetProjectFolder(bo.GetFolderBo{
		ParentId:  input.ParentID,
		OrgId:     orgId,
		UserId:    operatorId,
		ProjectId: input.ProjectID,
		Page:      page,
		Size:      size,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return folderList, nil
}
