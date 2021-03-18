package ldap_test

import (
	"fmt"
	"log"
	"regexp"
	"server/ldap"
	"testing"
	// ldap "github.com/go-ldap/ldap/v3"
)

const (
	DialURLStr string = "ldaps://192.168.67.14:636"
	Acc        string = `test\Administrator`
	Pwd        string = "Sz=server@2020"
)

func TestAO(t *testing.T) {
	name := ldap.EscapedSpecialCharacters("北京团2队导师2组(P)")
	// name := "北京团2队导师2组\\28P\\29"
	Orq := &ldap.OUReq{Name: name}
	client := &ldap.Client{}
	defer client.Close()
	client.Search(Orq)
	fmt.Println(Orq.GetOrgs()[0].DN)
}
func TestModityPwd(t *testing.T) {
	client := &ldap.Client{}
	adu := &ldap.User{}
	adu.DN = "CN=王丹,OU=深圳点猫科技有限公司,OU=组织机构,DC=test,DC=local"
	adu.SetPwd("Codemao@202006")
	if err := client.Modify(adu); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func TestUf(t *testing.T) {
	//OU=创作工具设计组,OU=工具中心设计部,OU=工具中心666,OU=深圳市点猫科技有限公司啊,OU=点猫科技,DC=test,DC=local
	client := &ldap.Client{}
	ur := &ldap.UserReq{
		Ou: []string{"点猫科技", "深圳市点猫科技有限公司啊", "工具中心666", "工具中心设计部", "创作工具设计组"},
	}
	if err := client.Search(ur); err != nil {
		log.Printf("moveUser 查询用户失败 Err:%s", err.Error())
		return
	}
	fmt.Println(ur.IsFound)
}

func TestDel(t *testing.T) {
	or := ldap.Organization{}
	or.DN = "OU=默认部门,OU=点猫科技,DC=test,DC=local"
	if err := or.Delete(&ldap.Client{}); err != nil {
		fmt.Println(err)
	}
}

func TestMoveOu(t *testing.T) {
	Orq := &ldap.OUReq{ID: 928497}
	client := &ldap.Client{}
	if err := client.Search(Orq); err != nil {
		fmt.Println(err)
	}
	orgs := Orq.GetOrgs()

	for _, o := range orgs {
		fmt.Println(o.DN)
	}
}

func TestMail(t *testing.T) {
	reg, _ := regexp.Compile(`\w+([-+.]\w+)*@codemao.cn$`)
	fmt.Println(reg.MatchString("chen123@codemao.cn"))
	fmt.Println(reg.MatchString("chen123@qq.cn"))
}
