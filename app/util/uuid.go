package util

import (
	"github.com/bwmarrin/snowflake"
	"log"
)

var defaultNode *snowflake.Node

func init() {
	InitNode(0)
}

func InitNode(nodeID int64) {
	var err error
	defaultNode, err = snowflake.NewNode(nodeID)
	if err != nil {
		log.Fatalf("[snowflake] init node failed: %s", err.Error())
		return
	}
}

func Generate() int64 {
	return int64(defaultNode.Generate())
}
