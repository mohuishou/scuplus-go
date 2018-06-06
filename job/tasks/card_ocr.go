package tasks

import (
	"strings"

	"github.com/json-iterator/go"

	"log"

	"github.com/mohuishou/scuplus-go/model"
	"github.com/mohuishou/scuplus-go/util/youtu"
)

// CardOCR 一卡通识别,id为失物招领id
func CardOCR(id uint) error {
	var lost model.LostFind
	if err := model.DB().Find(&lost, id).Error; err != nil {
		log.Println(err)
		return NotifyLostFindStatus(id, "一卡通识别失败", "请尝试修改图片")
	}
	rsp, err := youtu.OCR(lost.Pictures)
	if err != nil {
		log.Println(err)
		return NotifyLostFindStatus(id, "一卡通识别失败", "请尝试修改图片")
	}
	if rsp.ErrorCode != 0 {
		log.Println(rsp)
		return NotifyLostFindStatus(id, "一卡通识别失败", "请尝试修改图片")
	}
	cardInfo := map[string]string{}
	for _, item := range rsp.Items {
		switch item.Item {
		case "电话":
			cardInfo["no"] = item.ItemString
		case "姓名":
			cardInfo["name"] = item.ItemString
		case "公司":
			if strings.Contains(item.ItemString, "学院") {
				cardInfo["college"] = item.ItemString
			}
		}
	}
	if no, ok := cardInfo["no"]; !ok || no == "" {
		log.Println(cardInfo)
		return NotifyLostFindStatus(id, "一卡通识别失败", "请尝试修改图片")
	}
	str, err := jsoniter.MarshalToString(&cardInfo)
	if err != nil {
		log.Println(err)
		return NotifyLostFindStatus(id, "一卡通识别失败", "请尝试修改图片")
	}
	if err := model.DB().Model(&lost).Updates(map[string]interface{}{"card_info": str, "status": model.LostFindShow}).Error; err != nil {
		log.Println(err)
		return NotifyLostFindStatus(id, "一卡通识别失败", "请尝试修改图片")
	}

	// 识别成功，发送通知
	NotifyLostFindStatus(id, "一卡通识别成功", "已尝试给失主发送通知，谢谢！")

	user := model.User{}
	err = model.DB().Where("jwc_student_id = ?", cardInfo["no"]).Find(&user).Error
	if err != nil {
		log.Println(err)
		return err
	}

	// 查询成功之后，查找一卡通所有者，如果为We川大用户尝试发送通知
	return NotifyLostFind(user.ID, id)
}
