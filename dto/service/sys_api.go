package service

import (
	"anew-server/dto/request"
	"anew-server/models"
	"anew-server/pkg/utils"
	"gorm.io/gorm"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)


func (s *MysqlService) GetApis(req *request.ApiListReq) ([]models.SysApi, error) {
	var err error
	list := make([]models.SysApi, 0)
	query := s.db.Table(new(models.SysApi).TableName())
	name := strings.TrimSpace(req.Name)
	if name != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	method := strings.TrimSpace(req.Method)
	if method != "" {
		query = query.Where("method LIKE ?", fmt.Sprintf("%%%s%%", method))
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		query = query.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	category := strings.TrimSpace(req.Category)
	if category != "" {
		query = query.Where("category LIKE ?", fmt.Sprintf("%%%s%%", category))
	}

	// 查询条数
	err = query.Find(&list).Count(&req.PageInfo.Total).Error
	if err == nil {
		if req.PageInfo.All {
			// 不使用分页
			err = query.Find(&list).Error
		} else {
			// 获取分页参数
			limit, offset := req.GetLimit()
			err = query.Limit(limit).Offset(offset).Find(&list).Error
		}
	}

	return list, err
}

// 创建接口
func (s *MysqlService) CreateApi(req *request.CreateApiReq) (err error) {
	var api models.SysApi
	utils.Struct2StructByJson(req, &api)
	// 创建数据
	err = s.db.Create(&api).Error
	return
}

// 更新接口
func (s *MysqlService) UpdateApiById(id uint, req gin.H) (err error) {
	var oldApi models.SysApi
	query := s.db.Table(oldApi.TableName()).Where("id = ?", id).First(&oldApi)
	if query.Error == gorm.ErrRecordNotFound {
		return errors.New("记录不存在")
	}

	// 比对增量字段
	var m models.SysApi
	utils.CompareDifferenceStructByJson(oldApi, req, &m)
	// 更新指定列
	err = query.Updates(m).Error
	return
}

// 批量删除接口
func (s *MysqlService) DeleteApiByIds(ids []uint) (err error) {

	return s.db.Where("id IN (?)", ids).Delete(models.SysApi{}).Error
}