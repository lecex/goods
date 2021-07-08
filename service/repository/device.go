package repository

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2/util/log"

	"github.com/lecex/core/util"
	pb "github.com/lecex/goods/proto/goods"
)

//Goods 仓库接口
type Goods interface {
	Create(req *pb.Goods) (*pb.Goods, error)
	Delete(req *pb.Goods) (bool, error)
	Update(req *pb.Goods) (bool, error)
	Get(req *pb.Goods) (*pb.Goods, error)
	All(req *pb.Request) ([]*pb.Goods, error)
	List(req *pb.ListQuery) ([]*pb.Goods, error)
	Total(req *pb.ListQuery) (int64, error)
}

// GoodsRepository 用户仓库
type GoodsRepository struct {
	DB *gorm.DB
}

// All 获取所有商品信息
func (repo *GoodsRepository) All(req *pb.Request) (goodss []*pb.Goods, err error) {
	if err := repo.DB.Find(&goodss).Error; err != nil {
		log.Fatal(err)
		return nil, err
	}
	return goodss, nil
}

// List 获取所有商品信息
func (repo *GoodsRepository) List(req *pb.ListQuery) (goodss []*pb.Goods, err error) {
	db := repo.DB
	limit, offset := util.Page(req.Limit, req.Page) // 分页
	sort := util.Sort(req.Sort)                     // 排序 默认 created_at desc
	if req.Where != "" {
		db = db.Where(req.Where)
	}
	if err := db.Order(sort).Limit(limit).Offset(offset).Find(&goodss).Error; err != nil {
		log.Fatal(err)
		return nil, err
	}
	return goodss, nil
}

// Total 获取所有商品查询总量
func (repo *GoodsRepository) Total(req *pb.ListQuery) (total int64, err error) {
	goodss := []pb.Goods{}
	db := repo.DB
	// 查询条件
	if req.Where != "" {
		db = db.Where(req.Where)
	}
	if err := db.Find(&goodss).Count(&total).Error; err != nil {
		log.Fatal(err)
		return total, err
	}
	return total, nil
}

// Get 获取商品信息
func (repo *GoodsRepository) Get(goods *pb.Goods) (*pb.Goods, error) {
	if err := repo.DB.Where(&goods).Find(&goods).Error; err != nil {
		return nil, err
	}
	return goods, nil
}

// Create 创建商品
// bug 无商品名创建商品可能引起 bug
func (repo *GoodsRepository) Create(req *pb.Goods) (*pb.Goods, error) {
	err := repo.DB.Create(req).Error
	if err != nil {
		// 写入数据库未知失败记录
		log.Fatal(err)
		return req, fmt.Errorf("添加商品失败")
	}
	return req, nil
}

// Update 更新商品
func (repo *GoodsRepository) Update(req *pb.Goods) (bool, error) {
	if req.Id == 0 {
		return false, fmt.Errorf("请传入更新id")
	}
	err := repo.DB.Where("id = ?", req.Id).Updates(req).Error
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}

// Delete 删除商品
func (repo *GoodsRepository) Delete(req *pb.Goods) (bool, error) {
	if req.UserId == "" {
		return false, fmt.Errorf("请传入更新id")
	}
	err := repo.DB.Where("user_id = ?", req.UserId).Delete(req).Error
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}
