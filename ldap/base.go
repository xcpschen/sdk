package ldap

import (
	"fmt"

	ldap "github.com/go-ldap/ldap/v3"
)

type Base struct {
	DN          string
	ResetDn     string
	modifyAttrs []ModifyAttr `json:"-"`
}

//ModifyReq 修改属性
func (this *Base) ModifyReq() (*ldap.ModifyRequest, error) {
	if len(this.modifyAttrs) <= 0 {
		return nil, fmt.Errorf("修改属性参数不能为空")
	}
	if this.DN == "" {
		return nil, fmt.Errorf("修改DN不能为空")
	}
	return modifyReq(this.DN, this.modifyAttrs...)
}
func (this *Base) SetModifyAttr(ho ...ModifyAttr) {
	this.modifyAttrs = append(this.modifyAttrs, ho...)
}

//ModifyDNReq 修改目录
func (this *Base) ModifyDNReq() *ldap.ModifyDNRequest {
	return modifyDnReq(this.DN, this.ResetDn)
}
func (this *Base) isChangeDn() bool {
	if this.ResetDn == "" {
		return false
	}
	return true
}

func modifyDnReq(dn string, newDn string) *ldap.ModifyDNRequest {
	return ldap.NewModifyDNRequest(dn, newDn, true, "")
}
func modifyReq(dn string, ho ...ModifyAttr) (*ldap.ModifyRequest, error) {
	req := ldap.NewModifyRequest(dn, nil)
	for i, h := range ho {
		switch h.Opt {
		case AttrAdd:
			for attrType, attr := range h.Attrs {
				req.Add(attrType, attr.([]string))
			}
		case AttrDel:
			for attrType, attr := range h.Attrs {
				req.Delete(attrType, attr.([]string))
			}
		case AttrIncrement:
			for attrType, attr := range h.Attrs {
				req.Increment(attrType, attr.(string))
			}
		case AttrReplace:
			for attrType, attr := range h.Attrs {
				switch n := attr.(type) {
				case []string:
					req.Replace(attrType, n)
				case string:
					req.Replace(attrType, []string{n})
				}
			}
		default:
			return nil, fmt.Errorf("第%d操作属性类型错误,操作码%d", i, h.Opt)
		}
	}
	return req, nil
}

func (this *Base) Delete(client *Client) error {
	req := ldap.NewDelRequest(this.DN, nil)
	return client.Delete(req)
}
