package service

import (
	"anew-server/dto/request"
	"anew-server/dto/response"
	"anew-server/models"
	"anew-server/pkg/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sort"
)

// 获取用户菜单的切片
func (s *MysqlService) GetUserMenuList(roleId uint) ([]models.SysMenu, error) {
	//tree := make([]models.SysMenu, 0)
	var role models.SysRole
	err := s.db.Table(new(models.SysRole).TableName()).Preload("Menus", "status = ?", true).Where("id = ?", roleId).Find(&role).Error
	menus := make([]models.SysMenu, 0)
	if err != nil {
		return menus, err
	}
	// 生成菜单树
	//tree = GenMenuTree(nil, role.Menus)
	return role.Menus, nil
}

// 获取所有菜单
func (s *MysqlService) GetMenus() []models.SysMenu {
	//tree := make([]models.SysMenu, 0)
	menus := s.getAllMenu()
	// 生成菜单树
	//tree = GenMenuTree(nil, menus)
	return menus
}


// 生成菜单树
func GenMenuTree(parent *response.MenuTreeResp, menus []models.SysMenu) []response.MenuTreeResp {
	tree := make(response.MenuTreeRespList, 0)
	// 转为MenuTreeResponseStruct
	var resp []response.MenuTreeResp
	utils.Struct2StructByJson(menus, &resp)
	// parentId默认为0, 表示根菜单
	var parentId uint
	if parent != nil {
		parentId = parent.Id
	}
	for _, menu := range resp {
		// 父菜单编号一致
		if menu.ParentId == parentId {
			// 递归获取子菜单
			menu.Children = GenMenuTree(&menu, menus)
			// 加入菜单树
			tree = append(tree, menu)
		}
	}
	// 排序
	sort.Sort(tree)
	return tree
}



// 创建菜单
func (s *MysqlService) CreateMenu(req *request.CreateMenuReq) (err error) {
	var menu models.SysMenu
	utils.Struct2StructByJson(req, &menu)
	// 创建数据
	err = s.db.Create(&menu).Error
	return
}

// 更新菜单
func (s *MysqlService) UpdateMenuById(id uint, req gin.H) (err error) {
	var oldMenu models.SysMenu
	query := s.db.Table(oldMenu.TableName()).Where("id = ?", id).First(&oldMenu)
	if query.Error == gorm.ErrRecordNotFound {
		return errors.New("记录不存在")
	}
	// 比对增量字段,使用map确保gorm可更新零值
	var m map[string]interface{}
	utils.CompareDifferenceStructByJson(oldMenu, req, &m)
	// 更新指定列
	err = query.Updates(m).Error
	return
}

// 批量删除菜单
func (s *MysqlService) DeleteMenuByIds(ids []uint) (err error) {
	var menu models.SysMenu
	// 先解除父级关联
	err = s.db.Table(menu.TableName()).Where("parent_id IN (?)", ids).Update("parent_id",0).Error
	if err != nil{
		return err
	}
	// 再删除
	err = s.db.Where("id IN (?)", ids).Delete(&menu).Error
	if err != nil{
		return err
	}
	return
}



// 获取全部菜单, 非菜单树
func (s *MysqlService) getAllMenu() []models.SysMenu {
	menus := make([]models.SysMenu, 0)
	// 查询所有菜单
	s.db.Order("sort").Find(&menus)
	return menus
}
