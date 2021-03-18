package ldap

import (
	"crypto/tls"
	"fmt"
	"log"
	"regexp"

	"github.com/astaxie/beego"
	ldap "github.com/go-ldap/ldap/v3"
)

type Client struct {
	HostURL string
	Acc     string
	pwd     string
	conn    *ldap.Conn
}

type LdapReq interface {
	// searchReq(f Filter) *ldap.SearchRequest
}
type ReqAdd interface {
	AddReq() (*ldap.AddRequest, error)
}
type ReqSearch interface {
	searchReq() *ldap.SearchRequest
	setRelut(*ldap.SearchResult) error
	setRelutToSelf(*ldap.SearchResult) error
}

type PageData struct {
	PageSize uint32
	Page     uint32
}

const (
	AttrAdd       int = 0
	AttrDel       int = 1
	AttrReplace   int = 2
	AttrIncrement int = 3
)

var (
	DialURLStr       string = ""
	Acc              string = ""
	Pwd              string = ""
	BaseDN           string
	BasePrincipaName string
)

func init() {
	appConfig := beego.AppConfig
	BaseDN = appConfig.String("BaseDN")
	BasePrincipaName = appConfig.String("BasePrincipaName")
	DialURLStr = appConfig.String("DialURLStr")
	Acc = appConfig.String("Acc")
	Pwd = appConfig.String("Pwd")
	//BaseDN = "DC=codemao,DC=local"
	//BasePrincipaName = "codemao.local"
	//DialURLStr = "ldaps://192.168.67.253:636"
	//Acc = `codemao\dongwulin`
	//Pwd = "Codemao@2020"
}
func NewClient(host, acc, pwd string) (*Client, error) {
	return &Client{
		HostURL: host,
		Acc:     acc,
		pwd:     pwd,
	}, nil
}

func (c *Client) Close() {
	if c != nil {
		if c.conn != nil {
			c.conn.Close()
		}
	}
}
func (c *Client) init() error {
	if c.conn == nil {
		c.HostURL = DialURLStr
		conn, err := ldap.DialURL(DialURLStr, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
		if err != nil {
			log.Println(err)
			return err
		}
		// defer conn.Close()

		err = conn.Bind(Acc, Pwd)
		if err != nil {
			log.Println(err)
			return err
		}
		c.conn = conn
	}
	return nil
}
func (c *Client) Search(s ReqSearch) error {
	if err := c.init(); err != nil {
		return err
	}
	sr, err := c.conn.Search(s.searchReq())
	if err != nil {
		fmt.Println(err)
		reg, _ := regexp.Compile("LDAP Result Code 32 \"No Such Object\"")
		if reg.FindString(err.Error()) != "" {
			s.setRelut(nil)
			return nil
		}

		return err
	}
	if err := s.setRelut(sr); err != nil {
		return err
	}
	return nil
}

func (c *Client) SearchOne(s ReqSearch) error {
	if err := c.init(); err != nil {
		return err
	}

	sr, err := c.conn.Search(s.searchReq())
	if err != nil {
		reg, _ := regexp.Compile("LDAP Result Code 32 \"No Such Object\"")
		if reg.FindString(err.Error()) != "" {
			s.setRelut(nil)
			return nil
		}
		return err
	}
	if err := s.setRelutToSelf(sr); err != nil {
		return err
	}
	return nil
}

func (c *Client) SearchPage(s ReqSearch, pageSize uint32) error {
	if err := c.init(); err != nil {
		return err
	}
	sr, err := c.conn.SearchWithPaging(s.searchReq(), pageSize)
	if err != nil {
		return err
	}
	if err := s.setRelut(sr); err != nil {
		return err
	}
	return nil
}

func (c *Client) AddReq(ra ReqAdd) error {
	if err := c.init(); err != nil {
		return err
	}
	r, err := ra.AddReq()
	beego.Warn(r, err)
	if err != nil {
		return err
	}
	return c.conn.Add(r)
}

type IModifyReq interface {
	ModifyReq() (*ldap.ModifyRequest, error)
	ModifyDNReq() *ldap.ModifyDNRequest
	isChangeDn() bool
}

func (c *Client) Modify(mr IModifyReq) error {
	if err := c.init(); err != nil {
		return err
	}
	r, err := mr.ModifyReq()
	if err != nil {
		return err
	}
	if err = c.conn.Modify(r); err != nil {
		return err
	}
	if mr.isChangeDn() {
		if err := c.conn.ModifyDN(mr.ModifyDNReq()); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ModifyDN(req *ldap.ModifyDNRequest) error {
	if err := c.init(); err != nil {
		return err
	}
	return c.conn.ModifyDN(req)
}

func (c *Client) ModifyPwd(req *ldap.PasswordModifyRequest) error {
	if err := c.init(); err != nil {
		return err
	}
	respone, err := c.conn.PasswordModify(req)
	if err != nil {
		return err
	}
	fmt.Println(respone.GeneratedPassword)
	return nil
}

func (c *Client) Delete(req *ldap.DelRequest) error {
	if err := c.init(); err != nil {
		return err
	}
	return c.conn.Del(req)
}
