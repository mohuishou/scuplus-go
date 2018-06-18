package main

import (
	"flag"
	"log"
	"net/url"
	"strconv"

	"github.com/mohuishou/scu/jwc/course"

	"github.com/gocolly/colly"

	"github.com/mohuishou/scu/jwc"

	"time"

	"github.com/mohuishou/scuplus-go/model"
)

func main() {
	studentID := flag.String("u", "", "请输入用户学号")
	password := flag.String("p", "", "请输入密码")
	pages := flag.Int("page", 100, "请输入每页50条数据时，总页数")
	flag.Parse()
	log.Println(*studentID, *password)
	c, err := jwc.Login(*studentID, *password)
	if err != nil {
		panic(err)
	}
	updateCourses(c, *pages)
}

// 更新课程信息
// TODO: 课程更新信息有多次查询的问题，之后可以优化一下
func updateCourses(c *colly.Collector, pageNo int) error {

	params := url.Values{}
	params.Set("pageSize", "50")
	for i := 1; i <= pageNo; i++ {
		params.Set("pageNumber", strconv.Itoa(i))
		courses := course.Get(c, params)

		for _, course := range courses {
			c := model.Course{
				College:     course.College,
				CourseID:    course.CourseID,
				Name:        course.Name,
				LessonID:    course.LessonID,
				Credit:      course.Credit,
				ExamType:    course.ExamType,
				AllWeek:     course.AllWeek,
				Day:         course.Day,
				Session:     course.Session,
				Campus:      course.Campus,
				Building:    course.Building,
				Classroom:   course.Classroom,
				Max:         course.Max,
				StudentNo:   course.StudentNo,
				CourseLimit: course.CourseLimit,
			}

			tx := model.DB().Begin()

			// 获取教师信息,不存在就新建
			teachers := []model.Teacher{}
			for _, name := range course.Teachers {
				teacher := model.Teacher{
					College: course.College,
					Name:    name,
				}

				if err := tx.Where(teacher).FirstOrCreate(&teacher).Error; err != nil {
					tx.Rollback()
					return err
				}
				teachers = append(teachers, teacher)
			}
			c.Teachers = teachers

			// 如果数据库有10天之前的相同课程号课序号的课程，则删除
			tx.Where(model.Course{
				CourseID: c.CourseID,
				LessonID: c.LessonID,
			}).Where("updated_at < ?",
				time.Now().AddDate(0, 0, -10),
			).Delete(model.Course{})

			// 获取数据库已有信息, 存在则更新不存在则新建
			oldCourse := model.Course{}
			tx.Where(model.Course{
				CourseID: c.CourseID,
				LessonID: c.LessonID,
				Day:      c.Day,
				Session:  c.Session}).First(&oldCourse)
			if oldCourse.ID != 0 {
				tx.Model(&oldCourse).Related(&oldCourse.Teachers, "Teachers")

				// 查找失效的course_teacher关联关系,并且删除
			OLD:
				for _, oldTeacher := range oldCourse.Teachers {
					for _, t := range teachers {
						if t.ID == oldTeacher.ID {
							continue OLD
						}
					}
					del := []model.CourseTeacher{}
					tx.Where(model.CourseTeacher{CourseID: oldCourse.ID, TeacherID: oldTeacher.ID}).Find(&del).Delete(&del)
				}

				if err := tx.Model(&oldCourse).Update(c).Error; err != nil {
					tx.Rollback()
					return err
				}
			} else {
				if err := tx.Create(&c).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
			tx.Commit()
		}
	}
	return nil
}
