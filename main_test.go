package main

import (
	"fmt"
	"testing"

	_ "github.com/lecex/goods/providers/migrations" // 执行数据迁移
)

func TestGoodsCreate(t *testing.T) {
	fmt.Println(t)
}
