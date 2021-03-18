package ldap

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	ldap "github.com/go-ldap/ldap/v3"
	"golang.org/x/text/encoding/unicode"
)

type User struct {
	Base
	Name                       string `json:"displayName"`
	GivenName                  string
	FirstName                  string `json:"sn"`
	Pwd                        string `json:"-"`
	Phone                      string `json:"telephoneNumber"`
	Email                      string `json:"mail"`
	EmployeeID                 string `json:"EmployeeID"`
	UserAccountControl         string
	Department                 string
	Status                     string
	UserPrincipalName          string `json:"UserPrincipalName"`
	SAMAccountName             string `json:"sAMAccountName"`
	Job                        string `json:"title"`
	PhysicalDeliveryOfficeName string `json:"physicalDeliveryOfficeName"`
	Ou                         string `json:"ou"`

	F Filter `json:"-"`

	sr *ldap.SearchResult `json:"-"`
}

type UserRespone struct {
	F        Filter
	PageSize uint64
	Page     uint64
}

const (
	SCRIPT int = 1

	ACCOUNTDISABLE int = 2

	HOMEDIR_REQUIRED int = 8

	LOCKOUT int = 16

	PASSWD_NOTREQD int = 32

	PASSWD_CANT_CHANGE int = 64

	ENCRYPTED_TEXT_PWD_ALLOWED int = 128

	TEMP_DUPLICATE_ACCOUNT int = 256

	NORMAL_ACCOUNT int = 512

	INTERDOMAIN_TRUST_ACCOUNT int = 2048

	WORKSTATION_TRUST_ACCOUNT int = 4096

	SERVER_TRUST_ACCOUNT int = 8192

	DONT_EXPIRE_PASSWORD int = 65536

	MNS_LOGON_ACCOUNT int = 131072

	SMARTCARD_REQUIRED int = 262144

	TRUSTED_FOR_DELEGATION int = 524288

	NOT_DELEGATED int = 1048576

	USE_DES_KEY_ONLY int = 2097152

	DONT_REQ_PREAUTH int = 4194304

	PASSWORD_EXPIRED int = 8388608

	TRUSTED_TO_AUTH_FOR_DELEGATION int = 16777216
)

const (
	defaultPwd string = "!@#$123456qwer"
)

var defaultAttributes = []string{"cn", "ou", "title", "UserPrincipalName", "description", "sAMAccountName", "displayName", "mail", "mobile", "employeeID", "sn", "givenName", "department", "userAccountControl"}

func NewUser() *User {
	u := &User{}
	u.DN = BaseDN
	return u
}
func (this *User) objectClass() []string {
	return []string{"top", "person", "organizationalPerson", "user"}
}

func (this *User) SetUACCode(code string) {
	this.UserAccountControl = code
}

func (this *User) pwdEncoded(pwd string) string {
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	pwdEncoded, _ := utf16.NewEncoder().String("\"" + pwd + "\"")
	return pwdEncoded
}

func (this *User) data() *ldap.AddRequest {
	data := ldap.NewAddRequest(this.DN, nil)

	data.Attribute("objectClass", this.objectClass())

	return data
}
func (this *User) searchReq() *ldap.SearchRequest {
	return ldap.NewSearchRequest(this.DN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0,
		0,
		false,
		this.F.ToString(),
		defaultAttributes,
		[]ldap.Control{})
}

func (this *User) setRelut(s *ldap.SearchResult) error {
	this.sr = s
	for _, entry := range this.sr.Entries {
		if this.DN == entry.DN {

		}
	}
	return nil
}
func (this *User) setRelutToSelf(s *ldap.SearchResult) error {

	this.sr = s

	return nil
}

func (this *User) GetResultMap() []map[string]string {
	data := []map[string]string{}
	for _, entry := range this.sr.Entries {
		tmp := map[string]string{}
		for _, atr := range defaultAttributes {
			tmp[atr] = entry.GetAttributeValue(atr)
		}
		data = append(data, tmp)
	}
	return data
}

func (this *User) SetPwd(newPwd string) {

	this.SetModifyAttr(ModifyAttr{
		Opt: AttrReplace,
		Attrs: map[string]interface{}{
			"unicodePwd": this.pwdEncoded(newPwd),
		},
	})
}

func (this *User) Forbidden(client *Client) error {
	this.SetModifyAttr(ModifyAttr{
		Opt: AttrReplace,
		Attrs: map[string]interface{}{
			"userAccountControl": fmt.Sprintf("%d", ACCOUNTDISABLE),
		},
	})
	if err := client.Modify(this); err != nil {
		return err
	}
	return nil
}

func (this *User) Active(client *Client) error {
	this.SetModifyAttr(ModifyAttr{
		Opt: AttrReplace,
		Attrs: map[string]interface{}{
			"userAccountControl": fmt.Sprintf("%d", NORMAL_ACCOUNT),
		},
	})
	if err := client.Modify(this); err != nil {
		return err
	}
	return nil
}

func (this *User) AddReq() (*ldap.AddRequest, error) {
	if this.Name == "" {
		return nil, fmt.Errorf("名字")
	}
	if this.DN == "" {
		if this.Ou == "" {
			return nil, fmt.Errorf("请设置DN或ou")
		}
		this.DN = fmt.Sprintf("CN=%s,%s", this.Name, this.Ou)
	}
	if this.Email == "" {
		return nil, fmt.Errorf("请设置邮箱")
	}
	// if this.EmployeeID == "" {
	// 	return nil, fmt.Errorf("请设置工号")
	// }

	arr := strings.Split(this.Email, "@")
	this.UserPrincipalName = fmt.Sprintf("%s@%s", arr[0], BasePrincipaName)
	this.SAMAccountName = arr[0]
	name := []rune(this.Name)
	this.FirstName = string(name[:1])
	this.GivenName = string(name[1:])
	if this.GivenName == "" {
		this.GivenName = this.FirstName
	}
	ar := &ldap.AddRequest{
		DN: this.DN,
		Attributes: []ldap.Attribute{
			ldap.Attribute{
				Type: "objectclass",
				Vals: []string{"top", "person", "organizationalPerson", "user"},
			},
			ldap.Attribute{
				Type: "cn",
				Vals: []string{this.Name},
			},
			ldap.Attribute{
				Type: "sn",
				Vals: []string{this.FirstName},
			},
			ldap.Attribute{
				Type: "mail",
				Vals: []string{this.Email},
			},
		},
	}
	if this.Pwd != "" {
		ar.Attributes = append(ar.Attributes, ldap.Attribute{
			Type: "unicodePwd",
			Vals: []string{this.pwdEncoded(this.Pwd)},
		})
	} else {
		ar.Attributes = append(ar.Attributes, ldap.Attribute{
			Type: "unicodePwd",
			Vals: []string{this.pwdEncoded(defaultPwd)},
		})
	}
	if arr[0] != "" {
		ar.Attributes = append(ar.Attributes, ldap.Attribute{
			Type: "userPrincipalName",
			Vals: []string{this.UserPrincipalName},
		})
	}
	if this.SAMAccountName != "" {
		ar.Attributes = append(ar.Attributes, ldap.Attribute{
			Type: "sAMAccountName",
			Vals: []string{this.SAMAccountName},
		})
	}
	if this.Name != "" {
		ar.Attributes = append(ar.Attributes, ldap.Attribute{
			Type: "displayName",
			Vals: []string{this.Name},
		})
	}
	if this.GivenName != "" {
		ar.Attributes = append(ar.Attributes, ldap.Attribute{
			Type: "givenName",
			Vals: []string{this.GivenName},
		})
	}
	if this.EmployeeID != "" {
		ar.Attributes = append(ar.Attributes, ldap.Attribute{
			Type: "description",
			Vals: []string{this.EmployeeID},
		})
	}
	if this.Department != "" {
		ar.Attributes = append(ar.Attributes, ldap.Attribute{
			Type: "department",
			Vals: []string{this.Department},
		})
	}
	if this.Job != "" {
		ar.Attributes = append(ar.Attributes, ldap.Attribute{
			Type: "title",
			Vals: []string{this.Job},
		})
	}
	if this.PhysicalDeliveryOfficeName != "" {
		ar.Attributes = append(ar.Attributes, ldap.Attribute{
			Type: "physicalDeliveryOfficeName",
			Vals: []string{this.PhysicalDeliveryOfficeName},
		})
	}
	return ar, nil
}

func (this *User) AddUser(client *Client) error {
	if err := client.AddReq(this); err != nil {
		return fmt.Errorf("%s  userPrincipalName:%s", err.Error(), this.UserPrincipalName)
	}
	return this.Active(client)
}
func (this *User) IsForbidden() bool {
	if this.UserAccountControl == fmt.Sprintf("%d", NORMAL_ACCOUNT|ACCOUNTDISABLE) {
		return true
	}
	return false
}
func (this *User) ExistWithEmail(client *Client) (bool, error) {
	ur := &UserReq{
		Email: this.Email,
		// EmployeeID: this.EmployeeID,
	}
	if err := client.Search(ur); err != nil {
		return false, err
	}
	if !ur.IsFound {
		return false, nil
	}
	tmp := ur.GetByEmploeeID(this.DN)
	this.Email = tmp["email"]
	this.FirstName = tmp["sn"]
	this.Name = tmp["cn"]
	this.UserPrincipalName = tmp["userPrincipalName"]
	this.SAMAccountName = tmp["sAMAccountName"]
	this.GivenName = tmp["givenName"]
	this.EmployeeID = tmp["description"]
	this.Department = tmp["department"]
	this.Job = tmp["title"]
	this.PhysicalDeliveryOfficeName = tmp["physicalDeliveryOfficeName"]
	return true, nil
}
func (this *User) CheckWithEmployeeID(client *Client) (bool, error) {
	ur := &UserReq{
		// Email: this.Email,
		EmployeeID: this.EmployeeID,
	}
	if err := client.SearchPage(ur, 1000); err != nil {
		return false, err
	}
	if !ur.IsFound {
		return false, nil
	}
	tmp := ur.GetByEmploeeID(this.EmployeeID)
	//this.Email = tmp["email"]
	//this.FirstName = tmp["sn"]
	//this.Name = tmp["cn"]
	//this.UserPrincipalName = tmp["userPrincipalName"]
	//this.SAMAccountName = tmp["sAMAccountName"]
	//this.GivenName = tmp["givenName"]
	//this.EmployeeID = tmp["description"]
	fmt.Println(tmp)
	this.DN = tmp["DN"]
	this.Department = tmp["department"]
	this.Job = tmp["title"]
	this.PhysicalDeliveryOfficeName = tmp["physicalDeliveryOfficeName"]
	return true, nil
}

func (this *User) DiffOuAndMove(ou string, client *Client) error {
	newDn := fmt.Sprintf("CN=%s,%s", this.Name, ou)
	beego.Info(newDn)
	beego.Info(this.DN)
	if newDn != this.DN {
		if err := client.ModifyDN(ldap.NewModifyDNRequest(this.DN, fmt.Sprintf("CN=%s", this.Name), true, ou)); err != nil {
			return err
		}
		this.DN = newDn
	}
	return nil
}

func (this *User) DiffAndModify(dep, job, OfficeName string, client *Client) error {
	isEmpty := true
	if this.Department != dep && dep != "" {
		this.SetModifyAttr(ModifyAttr{
			Opt: AttrReplace,
			Attrs: map[string]interface{}{
				"department": dep,
			},
		})
		isEmpty = false
	}
	if this.Job != job && job != "" {
		this.SetModifyAttr(ModifyAttr{
			Opt: AttrReplace,
			Attrs: map[string]interface{}{
				"title": job,
			},
		})
		isEmpty = false
	}
	if this.PhysicalDeliveryOfficeName != OfficeName && OfficeName != "" {
		this.SetModifyAttr(ModifyAttr{
			Opt: AttrReplace,
			Attrs: map[string]interface{}{
				"physicalDeliveryOfficeName": job,
			},
		})
		isEmpty = false
	}
	if isEmpty {
		return nil
	}
	// if this.Pwd != "" {
	// 	this.SetPwd(this.Pwd)
	// }
	if err := client.Modify(this); err != nil {
		return err
	}
	this.Job = job
	this.Department = dep
	this.PhysicalDeliveryOfficeName = OfficeName
	return nil
}
