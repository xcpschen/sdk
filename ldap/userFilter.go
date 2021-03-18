package ldap

import (
	"encoding/json"
	"fmt"
	"strings"

	ldap "github.com/go-ldap/ldap/v3"
)

type UserReq struct {
	Name       string
	Email      string
	EmployeeID string
	Ou         []string
	DN         string
	IsFound    bool
	f          *FilterOperation

	s *ldap.SearchResult
}
type UserPage struct {
	PageData
	Rows []User
}

const (
	DefaultPageSize uint32 = 1000
	DefaultPage     uint32 = 1
)

func NewUserReq() Filter {
	f := NewAddFilterOperation()
	f.SetVal(map[string]string{"objectclass": "organizationalPerson"})
	return f
}
func (this *UserReq) param() {
	f := &FilterOperation{opt: "&"}
	f.SetVal(map[string]string{"objectclass": "organizationalPerson"})
	tmp := map[string]string{}
	if this.Name != "" {
		tmp["cn"] = fmt.Sprintf("%s", this.Name)
	}
	if this.Email != "" {
		tmp["mail"] = fmt.Sprintf("%s", this.Email)
	}

	if this.EmployeeID != "" {
		// tmp["employeeID"] = fmt.Sprintf("%s", this.EmployeeID)
		tmp["description"] = fmt.Sprintf("%s", this.EmployeeID)
	}
	if len(tmp) > 0 {
		f.SetVal(tmp)
	}

	this.f = f
}
func (this *UserReq) getOuStr() string {
	if this.DN == "" {
		str := []string{}
		if len(this.Ou) > 0 {
			length := len(this.Ou)
			for i := 0; i < length; i++ {
				temp := this.Ou[length-1-i]
				str = append(str, fmt.Sprintf("ou=%s", temp))
			}
		}
		str = append(str, BaseDN)
		return strings.Join(str, ",")
	}
	return this.DN
}
func (this *UserReq) All() error {
	this.param()
	c := &Client{}
	if err := c.SearchPage(ReqSearch(this), 1000); err != nil {
		return err
	}
	return nil
}
func (this *UserReq) List(arc ...uint32) error {
	this.param()
	c := &Client{}
	if err := c.Search(ReqSearch(this)); err != nil {
		return err
	}
	return nil
}

func (this *UserReq) searchReq() *ldap.SearchRequest {
	this.param()
	return ldap.NewSearchRequest(this.getOuStr(), ldap.ScopeWholeSubtree,
		ldap.DerefAlways,
		0,
		0,
		false,
		this.f.ToString(),
		defaultAttributes,
		[]ldap.Control{ldap.NewControlPaging(1000)})
}

func (this *UserReq) setRelut(s *ldap.SearchResult) error {
	if s == nil || len(s.Entries) <= 0 {
		this.IsFound = false
		return nil
	}
	this.IsFound = true
	this.s = s
	return nil
}
func (this *UserReq) GetUsers() (uses []*User) {
	uses = []*User{}
	if this.s == nil {
		return
	}
	for _, entry := range this.s.Entries {
		tmp := map[string]interface{}{}
		for _, attr := range defaultAttributes {
			tmp[attr] = entry.GetAttributeValue(attr)
		}
		tmp["DN"] = entry.DN
		b, _ := json.Marshal(tmp)
		tu := &User{}
		json.Unmarshal(b, tu)
		uses = append(uses, tu)
	}
	return
}
func (this *UserReq) Get(data interface{}) {
	if this.s == nil {
		return
	}
	for _, entry := range this.s.Entries {
		tmp := map[string]interface{}{}
		for _, attr := range defaultAttributes {
			tmp[attr] = entry.GetAttributeValue(attr)
		}
		tmp["DN"] = entry.DN
		b, _ := json.Marshal(tmp)
		json.Unmarshal(b, data)
	}
}
func (this *UserReq) GetByEmploeeID(pid string) map[string]string {
	tmp := map[string]string{}
	if this.s == nil {
		return nil
	}
	for _, entry := range this.s.Entries {
		if entry.GetAttributeValue("description") == pid {
			for _, attr := range defaultAttributes {
				tmp[attr] = entry.GetAttributeValue(attr)
			}
			tmp["DN"] = entry.DN
			break
		}
	}
	return tmp
}
func (this *UserReq) setRelutToSelf(s *ldap.SearchResult) error {
	this.s = s
	return nil
}
