package snowflake

import (
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	node *snowflake.Node
	once sync.Once
)

// Init 初始化雪花节点
func Init(nodeID int64) {
	once.Do(func() {
		n, err := snowflake.NewNode(nodeID)
		if err != nil {
			panic(err)
		}
		node = n
	})
}

// GenerateID 生成雪花ID
func GenerateID() uint64 {
	if node == nil {
		Init(1)
	}
	return uint64(node.Generate().Int64())
}
