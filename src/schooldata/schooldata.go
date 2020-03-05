package main

import (
	"bytes"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	belogs "github.com/astaxie/beego/logs"
	_ "github.com/cpusoft/goutil/conf"
	_ "github.com/cpusoft/goutil/logs"
	osutil "github.com/cpusoft/goutil/osutil"
)

type StringList []string

func (s StringList) Len() int {
	return len(s)
}
func (s StringList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s StringList) Less(i, j int) bool {

	iTmp := strings.Replace(s[i], `F:\share\我的坚果云\rongrong\20200220_学校数据抓取\html\华南理工大学研究生指导教师信息`, "", -1)
	iTmp = strings.Replace(iTmp, ".html", "", -1)
	belogs.Debug("s[i]", s[i], " iTmp:", iTmp)
	iInt, _ := strconv.Atoi(iTmp)

	jTmp := strings.Replace(s[j], `F:\share\我的坚果云\rongrong\20200220_学校数据抓取\html\华南理工大学研究生指导教师信息`, "", -1)
	jTmp = strings.Replace(jTmp, ".html", "", -1)
	belogs.Debug("s[j]", s[j], " jTmp:", jTmp)
	jInt, _ := strconv.Atoi(jTmp)

	belogs.Debug("i:", i, " iInt:", iInt, "  j:", j, "  jInt:", jInt)
	return iInt < jInt
}

const (
	//	startStr = `<table class="Grid_Line"`
	startStr = `<tr onmouseover="SetRowBgColor(this,true)"`
	endStr   = `<tr class="Grid_Page" style="height:25px;">`

	hrefStartStr = `title="导师详情信息" class="none" href="`
	hrefEndStr   = `" target="_blank"><img src`

	hrefColor = `#FFE793`
)

func main() {
	files, err := getAllFiles()
	if err != nil {
		return
	}
	var buffer bytes.Buffer

	buffer.WriteString(`<table class="Grid_Line" cellspacing="0" rules="all" border="1" id="dgData" style="width:100%;border-collapse:collapse;">`)
	for _, file := range files {
		belogs.Debug("getAllFiles():readFile file:", file)
		b, err := readFile(file)
		if err != nil {
			belogs.Error("getAllFiles():readFile file fail:", file, err)
			break
		}

		buffer.Write(b)
	}
	buffer.WriteString(`</table>`)

	if err != nil {
		return
	}
	err = ioutil.WriteFile(`E:\Go\go-study\src\schooldata\test.html`, buffer.Bytes(), 0644)
	if err != nil {
		belogs.Error("WriteFile:", err)
		return
	}
	belogs.Debug("getAllFiles():end:")

}
func getAllFiles() (files []string, err error) {

	file := `F:\share\我的坚果云\rongrong\20200220_学校数据抓取\html\`
	belogs.Debug("getAllFiles():input read file or path :", file)
	file = strings.TrimSpace(file)

	// 读取所有文件，加入到fileList列表中
	isDir, err := osutil.IsDir(file)
	if err != nil {
		belogs.Error("getAllFiles():IsDir err:", file, err)
		return nil, err
	}
	files = make([]string, 0)
	if isDir {
		suffixs := make(map[string]string)
		suffixs[".html"] = ".html"

		files, err = osutil.GetAllFilesBySuffixs(file, suffixs)
		if err != nil {
			return nil, err
		}
	} else {
		files = append(files, file)
	}
	belogs.Debug("getAllFiles(): len(files): ", len(files))

	sort.Sort(StringList(files))
	belogs.Debug(files)
	return files, nil
}

func readFile(file string) (b []byte, err error) {
	b, err = ioutil.ReadFile(file)
	if err != nil {
		belogs.Error("readFile():err:", file, err)
		return nil, err
	}
	start := bytes.Index(b, []byte(startStr))
	end := bytes.Index(b, []byte(endStr))
	b2 := b[start:end]

	b3 := bytes.Replace(b2, []byte(hrefStartStr), []byte(">"), -1)
	b4 := bytes.Replace(b3, []byte(hrefEndStr), []byte("<input type='hidden' value"), -1)
	b5 := bytes.Replace(b4, []byte(hrefColor), []byte("#FFFFFF"), -1)
	return b5, nil
}
