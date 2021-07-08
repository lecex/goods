package handler

import (
	"context"
	"encoding/json"
	"fmt"

	eventPB "github.com/lecex/core/proto/event"
	"github.com/micro/go-micro/v2"

	pb "github.com/lecex/goods/proto/goods"
	service "github.com/lecex/goods/service/repository"
)

// Goods 盘点
type Goods struct {
	Repo      service.Goods
	Publisher micro.Publisher
}

// List 获取所有设备
func (srv *Goods) List(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	goodss, err := srv.Repo.List(req.ListQuery)
	total, err := srv.Repo.Total(req.ListQuery)
	if err != nil {
		return err
	}
	res.Goodss = goodss
	res.Total = total
	return err
}

// Get 获取设备
func (srv *Goods) Get(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	goods, err := srv.Repo.Get(req.Goods)
	if err != nil {
		return err
	}
	res.Goods = goods
	return err
}

// Create 创建设备
func (srv *Goods) Create(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	_, err = srv.Repo.Create(req.Goods)
	if err != nil {
		res.Valid = false
		return fmt.Errorf("添加设备失败")
	}
	res.Valid = true
	return err
}

// Update 更新设备
func (srv *Goods) Update(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	valid, err := srv.Repo.Update(req.Goods)
	if err != nil {
		res.Valid = false
		return fmt.Errorf("更新设备失败:%s", err.Error())
	}
	res.Valid = valid
	if valid {
		goods, err := srv.Repo.Get(req.Goods)
		if err != nil {
			return err
		}
		if err := srv.publish(ctx, goods, "goods.Goodss.Update"); err != nil {
			return err
		}
	}
	return err
}

// Delete 删除设备
func (srv *Goods) Delete(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	valid, err := srv.Repo.Delete(req.Goods)
	if err != nil {
		res.Valid = false
		return fmt.Errorf("删除设备失败")
	}
	res.Valid = valid
	return err
}

// publish 消息发布
func (srv *Goods) publish(ctx context.Context, goods *pb.Goods, topic string) (err error) {
	data, _ := json.Marshal(&goods)
	event := &eventPB.Event{
		UserId:    "",
		GoodsInfo: goods.Info,
		GroupId:   "",
		Topic:     topic,
		Data:      data,
	}
	return srv.Publisher.Publish(ctx, event)
}
