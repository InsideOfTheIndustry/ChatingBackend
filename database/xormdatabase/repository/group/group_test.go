package group

import (
	"chatting/database/xormdatabase"
	"chatting/infrastructure/configServer"
	"chatting/infrastructure/logServer"
	"fmt"
	"testing"
)

func TestQueryGroupInfo(t *testing.T) {
	configServer.InitConfig("../../../../config/config.yaml")
	logServer.SetFileLevel("info") // 设置日志等级)
	xormdatabase.InitXormEngine()

	gr := GroupRepository{xormdatabase.DBEngine}
	result, err := gr.QueryGroupInfo(10000)
	fmt.Println(result)
	if err != nil {
		t.Fail()
	}

	if result.Groupid != 10000 {
		t.Fail()
	}
}

func TestQueryGroupMember(t *testing.T) {
	configServer.InitConfig("../../../../config/config.yaml")
	logServer.SetFileLevel("info") // 设置日志等级)
	xormdatabase.InitXormEngine()

	gr := GroupRepository{xormdatabase.DBEngine}
	result, _ := gr.QueryGroupMember(10000)
	fmt.Println(result)
	if len(result.GroupMembers) <= 1 {
		t.Fail()
	}
}
