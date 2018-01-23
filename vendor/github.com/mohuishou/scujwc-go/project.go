package scujwc

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

//ProjectCourse 培养方案课程
type ProjectCourse struct {
	CourseID   string `json:"course_id"`
	CourseName string `json:"course_name"`
	Credit     string `json:"credit"`
	CourseType string `json:"course_type"`
	Grade      string `json:"grade"`
	Time       string `json:"time"`
	ChooseType string `json:"choose_type"` // 0:未选课程 1:已选及格课程 2:已选不及格课程
}

//Project 培养方案
type Project struct {
	Name         string `json:"name"`
	Credit       string `json:"credit"`
	CreditPass   string `json:"credit_pass"`
	CourseChoose string `json:"course_choose"`
	CoursePass   string `json:"course_pass"`
	CourseFail   string `json:"course_fail"`
}

//ProjectResult 1
type ProjectResult struct {
	Project
	Required Projects
	Choose   Projects
}

//Projects 1
type Projects struct {
	Project
	Courses []ProjectCourse
}

//Project 获取培养方案
func (j *Jwc) Project() (data []ProjectResult, err error) {

	//获取goquery.Document 对象，以便解析需要的数据
	url := DOMAIN + "/gradeLnAllAction.do"
	doc, err := j.jPost(url, "type=ln&oper=lnfaqk&flag=zx")
	if err != nil {
		return data, err
	}

	//由于该网页由js动态渲染得到,所以通过html字符串，正则匹配得到
	html, err := doc.Html()
	if err != nil {
		return data, err
	}

	//首先去除所有的换行符
	re, err := regexp.Compile(`\s`)
	if err != nil {
		return data, err
	}

	html = re.ReplaceAllString(html, "")

	//匹配到调用js函数当中所需的参数值
	re, err = regexp.Compile(`add\((.*?)\);`)
	if err != nil {
		return data, err
	}
	s := re.FindAllStringSubmatch(html, -1)

	//去除所有的引号
	re, err = regexp.Compile(`"|'`)
	if err != nil {
		return data, err
	}

	/*params
	用来保存抓取到的参数值
	参数值说明：
	[
		0  // id
		-1 // pid
		公共课(最低修读学分:35 通过学分:33.0 已修课程门数:24 已通过课程门数:24 未通过课程门数:0) //值
		null //未知，无用数据
		ifra //未知，无用数据
		img/kzxx.gif //选择类型：无用数据
		img/kzxx.gif //选择类型： kzxx.gif：根节点，
								img/kcxx.gif：未选课程
								img/yxjg.gif: 已选及格课程
								img/yxbjg.gif:已选不及格课程

	]
	*/
	params := make([][]string, 10)
	for i := range s {
		//去除引号
		tmp := re.ReplaceAllString(s[i][1], "")

		//分割字符串保存参数值
		if i > 9 {
			params = append(params, strings.Split(tmp, ","))

		} else {
			params[i] = strings.Split(tmp, ",")
		}

	}

	data = make([]ProjectResult, 0)
	var projectRes ProjectResult
	var projects Projects

	//第一次获取根节点
	node0, err := getProNode(params, "-1")
	if err != nil {
		return data, err
	}
	for k, v := range node0 {
		projectRes.Project = v
		t, err := strconv.Atoi(k)
		if err != nil {
			return data, err
		}
		//不对中华文化和校任选做进一步分析
		if t < 10 {
			node1, err := getProNode(params, k)
			if err != nil {
				return data, err
			}
			i := 0
			for key, val := range node1 {
				projects.Project = val
				projects.Courses, err = getProcNode(params, key)
				if err != nil {
					return data, err
				}
				if i == 0 {
					projectRes.Required = projects
				} else {
					projectRes.Choose = projects
				}
				i++
			}
			data = append(data, projectRes)
		} else {
			//中华文化与校任选
			projectRes.Required = Projects{}
			projectRes.Choose = Projects{}
			data = append(data, projectRes)
		}
	}

	return data, err
}

func getProcNode(params [][]string, pid string) (data []ProjectCourse, err error) {
	data = make([]ProjectCourse, 1, 10)
	i := 0
	for k := range params {
		//找到节点
		if params[k][1] == pid {
			chooseType := "0"
			switch params[k][len(params[k])-1] {
			case "img/kcxx.gif":
				chooseType = "0"
			case "img/yxjg.gif":
				chooseType = "1"
			case "img/yxbjg.gif":
				chooseType = "2"
			}
			str := params[k][2]
			if chooseType != "0" {
				str = params[k][2] + params[k][3]
			}
			tmp, err := str2proc(str, chooseType)
			if err != nil {
				return nil, err
			}
			if i < 1 {
				data[i] = tmp
				i++
			} else {
				data = append(data, tmp)
			}
		}
	}
	return data, nil
}

//getProNode 获取project数据类型的节点
func getProNode(params [][]string, pid string) (data map[string]Project, err error) {
	data = make(map[string]Project)
	for k := range params {

		if pid == "-1" {
			if len(params[k]) != 12 {
				continue
			}
		}

		//找到节点
		if params[k][1] == pid {
			str := params[k][2] + params[k][3] + params[k][4] + params[k][5] + params[k][6]
			p, err := str2pro(str)
			if err != nil {
				return data, err
			}
			data[params[k][0]] = p
		}
	}
	return data, nil
}

//str2pro 通过字符串转换为project数据
func str2pro(s string) (p Project, err error) {

	//提取数字
	re, err := regexp.Compile(`(\d*\.)?\d+`)
	if err != nil {
		return p, err
	}

	tmp := re.FindAllString(s, -1)

	if len(tmp) < 5 {
		return p, errors.New("数据提取失败:" + s)
	}

	p.Credit = tmp[0]
	p.CreditPass = tmp[1]
	p.CourseChoose = tmp[2]
	p.CoursePass = tmp[3]
	p.CourseFail = tmp[4]

	//提取Name
	re, err = regexp.Compile(`(.*)\(`)
	if err != nil {
		return p, err
	}
	t := re.FindAllStringSubmatch(s, -1)
	p.Name = t[0][1]
	return p, nil
}

//str2proc  通过字符串转换为ProjectCourse数据
func str2proc(s string, chooseType string) (p ProjectCourse, err error) {
	p.ChooseType = chooseType
	//提取课程号和学分
	re, err := regexp.Compile(`\[(.*?)\]`)
	if err != nil {
		return p, err
	}
	tmp := re.FindAllStringSubmatch(s, -1)
	p.CourseID = tmp[0][1]
	if chooseType != "0" {
		p.Credit = tmp[1][1]
	}

	//提取课程名
	t := re.Split(s, -1)
	p.CourseName = t[1]
	if chooseType == "0" {
		return p, nil
	}
	s = t[2]

	//获取日期
	re, err = regexp.Compile(`\d{4}-\d{1,2}-\d{1,2}`)
	if err != nil {
		return p, err
	}
	t = re.FindAllString(s, 1)
	if t != nil {
		p.Time = t[0]
	}

	//获取成绩
	re, err = regexp.Compile(`(\d*\.)?\d+`)
	if err != nil {
		return p, err
	}
	temp := re.FindString(s)
	p.Grade = temp

	//获取类型
	re, err = regexp.Compile(`[^\x00-\xff]+`)
	if err != nil {
		return p, err
	}
	temp = re.FindString(s)
	p.CourseType = temp

	return p, nil
}
