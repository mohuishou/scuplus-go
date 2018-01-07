package course

import (
	"log"
	"net/url"
	"strconv"

	"github.com/mohuishou/scujwc-go"
	"github.com/mohuishou/scuplus-go/config"
	"github.com/mohuishou/scuplus-go/model"
)

// Task 任务
func Task(conf config.Config) {
	c := conf.CourseTask
	if err := updateCourses(c.StudentID, c.Password, c.PageNO); err != nil {
		log.Println("[Error]:", err)
	}
}

// 更新课程信息
// TODO: 课程更新信息有多次查询的问题，之后可以优化一下
func updateCourses(studentID int, password string, pageNo int) error {
	jwc, err := scujwc.NewJwc(studentID, password)
	if err != nil {
		return err
	}

	params := url.Values{}

	for i := 1; i <= pageNo; i++ {
		params.Set("pageNumber", strconv.Itoa(i))
		courses, err := jwc.GetCourse(params)
		if err != nil {
			return err
		}
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
