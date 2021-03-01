package notice

import (
	"fmt"
	"strings"
)

type Markdown struct {
	text []string
}

// 获取内容
func (this *Markdown) GetText() string {
	return strings.Join(this.text, "\n\n")
}

// 添加内容项
func (this *Markdown) AddItem(text ...string) {
	this.text = append(this.text, text...)
}

// 添加文本内容
func (this *Markdown) AddText(text ...string) {
	this.AddItem(strings.Join(text, "<br/>\n"))
}

// 添加加粗文本
func (this *Markdown) AddBoldText(text string) {
	this.AddItem(fmt.Sprintf("**%s**", text))
}

// 添加斜体文本
func (this *Markdown) AddItalicText(text string) {
	this.AddItem(fmt.Sprintf("*%s*", text))
}

// 添加加粗斜体文本
func (this *Markdown) AddBoldItalicText(text string) {
	this.AddItem(fmt.Sprintf("***%s***", text))
}

// 添加列表
func (this *Markdown) AddList(is_order bool, item ...string) {
	if is_order == true {
		for idx, val := range item {
			item[idx] = fmt.Sprintf("%d. %s", idx+1, val)
		}
	} else {
		for idx, val := range item {
			item[idx] = fmt.Sprintf("- %s", val)
		}
	}

	this.AddItem(strings.Join(item, "\n"))
}

// 添加头标题
func (this *Markdown) AddTitle(text string, level int) {
	lv := StrPad("#", level)
	this.AddItem(fmt.Sprintf("%s %s", lv, text))
}

// 添加链接
func (this *Markdown) AddLink(title, url string) {
	this.AddItem(fmt.Sprintf(`[%s](%s "%s")`, title, url, title))
}

// 添加图片
func (this *Markdown) AddImage(title, src string) {
	this.AddItem(fmt.Sprintf(`![%s](%s "%s")`, title, src, title))
}

// 添加引用
func (this *Markdown) AddQuote(text string) {
	this.AddItem(fmt.Sprintf("> %s", text))
}

// 添加分割线
func (this *Markdown) AddSplitline() {
	this.AddItem("---")
}
