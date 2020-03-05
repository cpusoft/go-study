package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type TeacherDetail struct {
	Name     string
	Url      string
	PhotoUrl string
	Gender   string
	Birth    string
	//籍贯
	NativePlace string
	//民族
	Nation string
	//政治面貌
	Political string
	//最后学历
	FinalEducation string
	//最后学位
	FinalDegree string
	//技术职称
	TechnicalTitle string
	//导师类别
	TutorType string
	//行政职务
	AdministrativePost string
	Email              string
	//工作单位
	WorkUnit string
	//邮政编码
	PostalCode string
	//通讯地址
	PostalAddress string
	//单位电话
	WorkTelephone string
	//个人主页
	PersonalHomepage string
	//个人简介
	PersonalProfile string
	//工作经历
	WorkExperience string
	//教育经历
	EducationExperience string
	//获奖、荣誉称号
	AwardsHonoraryTitle string
	//社会、学会及学术兼职
	SocialAcademicPartTimeJob string
	//研究领域
	ResearchField string
	//科研项目
	ResearchProject string
	//发表论文
	PublishThesis string
	//出版专著和教材
	PublishingTextbooks string
	//科研创新
	ScientificResearchInnovation string
	//教学活动
	TeachingActivity string
	//指导学生情况
	Guidestudent string
	//我的团队
	MyTeam string
}

func (c *TeacherDetail) HtmlString() string {
	buffer := bytes.NewBufferString("<tr><td>")
	buffer.WriteString(c.Name)
	buffer.WriteString("</td><td>")
	buffer.WriteString(c.Url)
	buffer.WriteString("</td><td>")
	buffer.WriteString(c.PhotoUrl)
	buffer.WriteString("</td><td>")
	buffer.WriteString(c.Gender)
	buffer.WriteString("</td><td>")
	buffer.WriteString(c.Birth)
	buffer.WriteString("</td><td>")
	buffer.WriteString(c.NativePlace)
	buffer.WriteString("</td><td>")
	buffer.WriteString(c.Nation)
	buffer.WriteString("</td><td>")
	buffer.WriteString(c.Political)
	buffer.WriteString("</td><td>")
	buffer.WriteString(c.FinalEducation)
	buffer.WriteString("</td><td>")
	buffer.WriteString(c.FinalDegree)
	buffer.WriteString("</td><td>")
	buffer.WriteString(c.TechnicalTitle)
	buffer.WriteString("</td><td>")
	buffer.WriteString(c.TutorType)
	buffer.WriteString("</td><td>")
	buffer.WriteString(c.AdministrativePost)
	buffer.WriteString("</td><td>")
	buffer.WriteString(c.Email)
	buffer.WriteString("</td><td>")
	buffer.WriteString(c.WorkUnit)
	buffer.WriteString("</td><td>")
	buffer.WriteString(c.PostalCode)
	buffer.WriteString("</td><td>")
	buffer.WriteString(c.PostalAddress)
	buffer.WriteString("</td><td>")
	buffer.WriteString(c.WorkTelephone)
	buffer.WriteString("</td><td>")
	buffer.WriteString(c.PersonalHomepage)
	buffer.WriteString("</td><td >")
	buffer.WriteString(replaceHtml(c.PersonalProfile))
	buffer.WriteString("</td><td >")
	buffer.WriteString(replaceHtml(c.WorkExperience))
	buffer.WriteString("</td><td >")
	buffer.WriteString(replaceHtml(c.EducationExperience))
	buffer.WriteString("</td><td >")
	buffer.WriteString(replaceHtml(c.AwardsHonoraryTitle))
	buffer.WriteString("</td><td >")
	buffer.WriteString(replaceHtml(c.SocialAcademicPartTimeJob))
	buffer.WriteString("</td><td >")
	buffer.WriteString(replaceHtml(c.ResearchField))
	buffer.WriteString("</td><td >")
	buffer.WriteString(replaceHtml(c.ResearchProject))
	buffer.WriteString("</td><td >")
	buffer.WriteString(replaceHtml(c.PublishThesis))
	buffer.WriteString("</td><td >")
	buffer.WriteString(replaceHtml(c.PublishingTextbooks))
	buffer.WriteString("</td><td >")
	buffer.WriteString(replaceHtml(c.ScientificResearchInnovation))
	buffer.WriteString("</td><td >")
	buffer.WriteString(replaceHtml(c.TeachingActivity))
	buffer.WriteString("</td><td >")
	buffer.WriteString(replaceHtml(c.Guidestudent))
	buffer.WriteString("</td><td >")
	buffer.WriteString(replaceHtml(c.MyTeam))
	buffer.WriteString("</td></tr>")
	return buffer.String()

}
func (c *TeacherDetail) ToStrings() []string {
	s := make([]string, 0)
	s = append(s, c.Name)
	s = append(s, c.Url)
	s = append(s, c.PhotoUrl)
	s = append(s, c.Gender)
	s = append(s, c.Birth)
	s = append(s, c.NativePlace)
	s = append(s, c.Nation)
	s = append(s, c.Political)
	s = append(s, c.FinalEducation)
	s = append(s, c.FinalDegree)
	s = append(s, c.TechnicalTitle)
	s = append(s, c.TutorType)
	s = append(s, c.AdministrativePost)
	s = append(s, c.Email)
	s = append(s, c.WorkUnit)
	s = append(s, c.PostalCode)
	s = append(s, c.PostalAddress)
	s = append(s, c.WorkTelephone)
	s = append(s, c.PersonalHomepage)
	s = append(s, replaceCsv(c.PersonalProfile))
	s = append(s, replaceCsv(c.WorkExperience))
	s = append(s, replaceCsv(c.EducationExperience))
	s = append(s, replaceCsv(c.AwardsHonoraryTitle))
	s = append(s, replaceCsv(c.SocialAcademicPartTimeJob))
	s = append(s, replaceCsv(c.ResearchField))
	s = append(s, replaceCsv(c.ResearchProject))
	s = append(s, replaceCsv(c.PublishThesis))
	s = append(s, replaceCsv(c.PublishingTextbooks))
	s = append(s, replaceCsv(c.ScientificResearchInnovation))
	s = append(s, replaceCsv(c.TeachingActivity))
	s = append(s, replaceCsv(c.Guidestudent))
	s = append(s, replaceCsv(c.MyTeam))
	return s

}

func main() {
	urlFileName := `E:\Go\go-study\src\schooldata2\urls.txt`
	htmlFileName := `E:\Go\go-study\src\schooldata2\teachers.html`
	csvFileName := `E:\Go\go-study\src\schooldata2\teachers_2.csv`
	fmt.Println(urlFileName, htmlFileName, csvFileName)

	urls, err := readFile(urlFileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	teacherDetails, err := getTeacherDetails(urls)
	if err != nil {
		fmt.Println(err)
		return
	}
	//err = saveHtmlFile(teacherDetails, htmlFileName)
	err = saveCsvFile(teacherDetails, csvFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func getTeacherDetails(urls []string) (teacherDetails []TeacherDetail, err error) {
	for i, url := range urls {
		fmt.Println(i, url)
		if len(url) == 0 {
			continue
		}
		detail, err := getTeacherDetail(url)
		if err != nil {
			return nil, err
		}
		teacherDetails = append(teacherDetails, detail)
		time.Sleep(time.Duration(1) * time.Second)
	}
	return teacherDetails, nil
}

func getTeacherDetail(url string) (teacherDetail TeacherDetail, err error) {
	//url := `http://yanzhao.scut.edu.cn/open/ExpertInfo.aspx?zjbh=2bYKzGSqoltzzj5aH2r-Zg==`
	//url := `https://yanzhao.scut.edu.cn/open/ExpertInfo.aspx?zjbh=FDt5poe5RTYiSlSRAMcWnQ==`
	res, err := http.Get(url)

	if err != nil {
		fmt.Println(url, err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println(url, res.StatusCode)
		fmt.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(url, err)
		return
	}

	teacherDetail.Url = url
	photoUrl, exist := doc.Find(".tbline img").Attr("src")
	fmt.Println(photoUrl, exist)
	if strings.Compare("../images/noimage.gif", photoUrl) != 0 {
		photoUrl = strings.Replace(photoUrl, `..`, "http://yanzhao.scut.edu.cn", -1)
		fmt.Println(photoUrl)
		teacherDetail.PhotoUrl = photoUrl
	}
	// Find the review items
	//个人简介
	doc.Find(".tbline .tbline td").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		v1, _ := s.Html()
		v2 := strings.TrimSpace(v1)
		//fmt.Printf("%s\r\n", v2)
		switch i {
		case 1:
			teacherDetail.Name = v2
		case 3:
			teacherDetail.Gender = v2
		case 5:
			teacherDetail.Birth = v2
		case 7:
			teacherDetail.NativePlace = v2
		case 9:
			teacherDetail.Nation = v2
		case 11:
			teacherDetail.Political = v2
		case 13:
			teacherDetail.FinalEducation = v2
		case 15:
			teacherDetail.FinalDegree = v2
		case 17:
			teacherDetail.TechnicalTitle = v2
		case 19:
			teacherDetail.TutorType = v2
		case 21:
			teacherDetail.AdministrativePost = v2
		case 23:
			teacherDetail.Email = v2
		case 25:
			teacherDetail.WorkUnit = v2
		case 27:
			teacherDetail.PostalCode = v2
		case 29:
			teacherDetail.PostalAddress = v2
		case 31:
			teacherDetail.WorkTelephone = v2
		case 33:
			teacherDetail.PersonalHomepage = v2
		}
	})
	value1, value2 := "", ""
	fmt.Println(value1, value2)

	//value1 := strings.TrimSpace(doc.Find("#contentParent_divGrjj1").Text())
	value2, _ = doc.Find("#contentParent_divGrjj2").Html()
	//fmt.Println(value2)
	teacherDetail.PersonalProfile = strings.TrimSpace(value2) //strings.Replace(value2, "\n", "<br/>", -1)

	//value1 = strings.TrimSpace(doc.Find("#contentParent_divGzjl1").Text())
	value2, _ = doc.Find("#contentParent_divGzjl2").Html()
	//fmt.Println(value2)
	teacherDetail.WorkExperience = strings.TrimSpace(value2) //strings.Replace(value2, "\n", "<br/>", -1)

	//value1 = strings.TrimSpace(doc.Find("#contentParent_divJyjl1").Text())
	value2, _ = doc.Find("#contentParent_divJyjl2").Html()
	//fmt.Println(value2)
	teacherDetail.EducationExperience = strings.TrimSpace(value2) //strings.Replace(value2, "\n", "<br/>", -1)

	//value1 = strings.TrimSpace(doc.Find("#contentParent_divHjry1").Text())
	value2, _ = doc.Find("#contentParent_divHjry2").Html()
	//fmt.Println(value2)
	teacherDetail.AwardsHonoraryTitle = strings.TrimSpace(value2) //strings.Replace(value2, "\n", "<br/>", -1)

	//value1 = strings.TrimSpace(doc.Find("#contentParent_divShjz1").Text())
	value2, _ = doc.Find("#contentParent_divShjz2").Html()
	//fmt.Println(value2)
	teacherDetail.SocialAcademicPartTimeJob = strings.TrimSpace(value2) //strings.Replace(value2, "\n", "<br/>", -1)

	//value1 = strings.TrimSpace(doc.Find("#contentParent_divYjly1").Text())
	value2, _ = doc.Find("#contentParent_divYjly2").Html()
	//fmt.Println(value2)
	teacherDetail.ResearchField = strings.TrimSpace(value2) //strings.Replace(value2, "\n", "<br/>", -1)

	//value1 = strings.TrimSpace(doc.Find("#contentParent_divKyxm1").Text())
	value2, _ = doc.Find("#contentParent_divKyxm2").Html()
	//fmt.Println(value2)
	teacherDetail.ResearchProject = strings.TrimSpace(value2) //strings.Replace(value2, "\n", "<br/>", -1)

	//value1 = strings.TrimSpace(doc.Find("#contentParent_divFblw1").Text())
	value2, _ = doc.Find("#contentParent_divFblw2").Html()
	//fmt.Println(value2)
	teacherDetail.PublishThesis = strings.TrimSpace(value2) //strings.Replace(value2, "\n", "<br/>", -1)

	//value1 = strings.TrimSpace(doc.Find("#contentParent_divZzjc1").Text())
	value2, _ = doc.Find("#contentParent_divZzjc2").Html()
	//fmt.Println(value2)
	teacherDetail.PublishingTextbooks = strings.TrimSpace(value2) //strings.Replace(value2, "\n", "<br/>", -1)

	//value1 = strings.TrimSpace(doc.Find("#contentParent_divKycx1").Text())
	value2, _ = doc.Find("#contentParent_divKycx2").Html()
	//fmt.Println(value2)
	teacherDetail.ScientificResearchInnovation = strings.TrimSpace(value2) //strings.Replace(value2, "\n", "<br/>", -1)

	//value1 = strings.TrimSpace(doc.Find("#contentParent_divJxhd1").Text())
	value2, _ = doc.Find("#contentParent_divJxhd2").Html()
	//fmt.Println(value2)
	teacherDetail.TeachingActivity = strings.TrimSpace(value2) // strings.Replace(value2, "\n", "<br/>", -1)

	//value1 = strings.TrimSpace(doc.Find("#contentParent_divZdxs1").Text())
	value2, _ = doc.Find("#contentParent_divZdxs2").Html()
	//fmt.Println(value2)
	teacherDetail.Guidestudent = strings.TrimSpace(value2) //strings.Replace(value2, "\n", "<br/>", -1)

	//value1 = strings.TrimSpace(doc.Find("#contentParent_divWdtd1").Text())
	value2, _ = doc.Find("#contentParent_divWdtd2").Html()
	//fmt.Println(value2)
	teacherDetail.MyTeam = strings.TrimSpace(value2) //strings.Replace(value2, "\n", "<br/>", -1)
	return
}

func readFile(fileName string) (lines []string, err error) {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return nil, err
	}
	defer file.Close()
	lines = make([]string, 0)
	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		lines = append(lines, line)

		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				return lines, nil
			} else {
				fmt.Println("Read file error!", err)
				return nil, err
			}
		}
	}
}

func saveCsvFile(teacherDetails []TeacherDetail, fileName string) (err error) {
	nfs, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("can not create file, err is %+v", err)
		return err
	}
	defer nfs.Close()
	nfs.Seek(0, io.SeekEnd)
	nfs.WriteString("\xEF\xBB\xBF")

	w := csv.NewWriter(nfs)
	//设置属性
	w.Comma = ','
	w.UseCRLF = true
	row := []string{"姓名", "URL", "PhotoUrl", "性别", "出生年月", "籍贯", "民族", "政治面貌", "最后学历", "最后学位", "技术职称", "导师类别", "行政职务", "Email", "工作单位", "邮政编码", "通讯地址", "单位电话", "个人主页", "个人简介", "工作经历", "教育经历", "获奖、荣誉称号", "社会、学会及学术兼职", "研究领域", "科研项目", "发表论文", "出版专著和教材", "科研创新", "教学活动", "指导学生情况", "我的团队"}
	err = w.Write(row)
	for i, _ := range teacherDetails {
		if err = w.Write(teacherDetails[i].ToStrings()); err != nil {
			fmt.Println(err)
			return err
		}
		w.Flush()
	}
	w.Flush()
	return nil
}

func saveHtmlFile(teacherDetails []TeacherDetail, fileName string) (err error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()
	newWriter := bufio.NewWriterSize(file, 1024*10)
	newWriter.WriteString("<table border='1' cellspacing='0'>")
	newWriter.WriteString(`<tr><td>姓名</td>
	<td>URL</td>
	<td>PhotoUrl</td>
	<td>性别</td>
	<td>出生年月</td>
	<td>籍贯</td>
	<td>民族</td>
	<td>政治面貌</td>
	<td>最后学历</td>
	<td>最后学位</td>
	<td>技术职称</td>
	<td>导师类别</td>
	<td>行政职务</td>
	<td>Email</td>
	<td>工作单位</td>
	<td>邮政编码</td>
	<td>通讯地址</td>
	<td>单位电话</td>
	<td>个人主页</td>
	<td>个人简介</td>
	<td>工作经历</td>
	<td>教育经历</td>
	<td>获奖、荣誉称号</td>
	<td>社会、学会及学术兼职</td>
	<td>研究领域</td>
	<td>科研项目</td>
	<td>发表论文</td>
	<td>出版专著和教材</td>
	<td>科研创新</td>
	<td>教学活动</td>
	<td>指导学生情况</td>
	<td>我的团队</td></tr>`)

	for i, _ := range teacherDetails {
		if _, err = newWriter.WriteString(teacherDetails[i].HtmlString()); err != nil {
			return err
		}
	}
	newWriter.WriteString("</table>")
	if err = newWriter.Flush(); err != nil {
		return err
	}
	fmt.Println("write file successful")
	return nil

}

func replaceHtml(s string) string {

	return "=\"" + strings.Replace(s, "<br/>", "\"&CHAR(10)&\"", -1) + "\""
}
func replaceCsv(s string) string {

	return strings.Replace(s, "<br/>", "$$$$", -1)
}
