package ldap

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	ldap "github.com/go-ldap/ldap/v3"
)

type (
	ModifyAttr struct {
		Opt   int
		Attrs map[string]interface{}
	}
	Organization struct {
		Base
		OU          string
		Name        string
		Description string
	}
	OUReq struct {
		Ou         []string
		Name       string
		PDN        string
		ID         int64
		IsNotFound bool
		f          *FilterOperation
		s          *ldap.SearchResult
	}
)

const (
	OrgnObjetClass string = "top organizationalUnit"
)

var OrganizationKey = []string{"name", "ou", "Description", "dn"}

func EscapedSpecialCharacters(key string) string {
	return strings.Replace(
		strings.Replace(
			strings.Replace(
				strings.Replace(key, "*", "\\2A", -1),
				"(", "\\28", -1),
			")", "\\29", -1),
		"Nul", "\\00", -1)
}
func GenderOu(Ou ...string) string {
	var ous []string
	length := len(Ou)
	for i := 0; i < length; i++ {
		temp := Ou[length-1-i]
		ous = append(ous, fmt.Sprintf("OU=%s", EscapedSpecialCharacters(temp)))
	}
	ous = append(ous, BaseDN)
	return strings.Join(ous, ",")
}
func NewOrganrization(DN, name, Description string) *Organization {
	o := &Organization{
		Name:        EscapedSpecialCharacters(name),
		Description: EscapedSpecialCharacters(Description),
	}
	o.DN = EscapedSpecialCharacters(DN)
	return o
}
func (this *Organization) AddReq() (*ldap.AddRequest, error) {
	if this.Name == "" {
		return nil, fmt.Errorf("ou name不能为空")
	}
	add := &ldap.AddRequest{
		DN: this.DN,
		Attributes: []ldap.Attribute{
			ldap.Attribute{
				Type: "objectclass",
				Vals: []string{"top", "organizationalUnit"},
			},
			ldap.Attribute{
				Type: "name",
				Vals: []string{this.Name},
			},
			ldap.Attribute{
				Type: "description",
				Vals: []string{this.Description},
			},
		},
	}
	return add, nil
}

func (this *OUReq) param() Filter {
	this.f = &FilterOperation{opt: "&"}
	this.f.SetVal(map[string]string{"objectclass": "organizationalUnit"})
	tmp := map[string]string{}
	if this.Name != "" {
		tmp["name"] = EscapedSpecialCharacters(this.Name)
	}
	/*if this.ID > 0 {
	 	tmp["description"] = fmt.Sprintf("%d", this.ID)

	}*/
	if len(tmp) > 0 {
		this.f.SetVal(tmp)
	}
	beego.Info(this.f.ToString())
	return this.f
}
func (this *OUReq) getOuStr() string {
	return GenderOu(this.Ou...)

}
func (this *OUReq) List() error {
	c := &Client{}
	if err := c.SearchPage(ReqSearch(this), DefaultPageSize); err != nil {
		return err
	}
	return nil
}
func (this *OUReq) searchReq() *ldap.SearchRequest {
	filter := this.param()
<<<<<<< HEAD
=======
	beego.Warn(this.PDN)
>>>>>>> 5de25e465a37fcd77c5d044bfe0e710f858dca98
	return ldap.NewSearchRequest(this.PDN, ldap.ScopeWholeSubtree,
		ldap.DerefAlways,
		0,
		0,
		false,
		filter.ToString(),
		OrganizationKey,
		[]ldap.Control{})
}

func (this *OUReq) setRelut(s *ldap.SearchResult) error {
	if s == nil {
		this.IsNotFound = true
		return nil
	}
	if len(s.Entries) <= 0 {
		this.IsNotFound = true
	}
	this.s = s
	return nil
}
func (this *OUReq) setRelutToSelf(s *ldap.SearchResult) error {
	if s == nil {
		this.IsNotFound = true
	}
	this.s = s
	fmt.Println(this.s.Entries)
	return nil
}

func (this *OUReq) GetResponeData() []string {
	var ou []string
	for _, entry := range this.s.Entries {
		ou = append(ou, entry.GetAttributeValue("name"))
	}
	return ou
}

func (this *OUReq) GetOrgs() []*Organization {
	orgs := []*Organization{}
	if this.s == nil {
		return nil
	}
	for _, entry := range this.s.Entries {
		org := &Organization{
			Name:        EscapedSpecialCharacters(entry.GetAttributeValue("name")),
			Description: EscapedSpecialCharacters(entry.GetAttributeValue("Description")),
		}
		org.DN = EscapedSpecialCharacters(entry.DN)
		orgs = append(orgs, org)
	}
	return orgs
}

func (this *Organization) DiffOuAndMove(name, newDn string, client *Client) error {
<<<<<<< HEAD
	name = EscapedSpecialCharacters(name)
=======
	// newDn := EscapedSpecialCharacters(fmt.Sprintf("OU=%s,%s", name, ou))
>>>>>>> 5de25e465a37fcd77c5d044bfe0e710f858dca98
	ou := strings.Replace(newDn, fmt.Sprintf("OU=%s,", name), "", -1)
	if newDn != this.DN {
		beego.Warn(newDn, "<==", this.DN)
		if err := client.ModifyDN(ldap.NewModifyDNRequest(this.DN, fmt.Sprintf("OU=%s", name), true, ou)); err != nil {
			return err
		}
		this.DN = newDn
	}
	return nil
}
