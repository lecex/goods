package handler

import (
	"github.com/micro/go-micro/v2"

	goodsPB "github.com/lecex/goods/proto/goods"

	db "github.com/lecex/goods/providers/database"
	service "github.com/lecex/goods/service/repository"
)

const topic = "event"

// Register 注册
func Register(srv micro.Service) {
	server := srv.Server()
	publisher := micro.NewPublisher(topic, srv.Client())
	// 获取 broker 实例
	goodsPB.RegisterGoodssHandler(server, &Goods{&service.GoodsRepository{db.DB}, publisher})
}
