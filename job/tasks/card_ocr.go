package tasks

import (
	"errors"
	"strings"

	"github.com/json-iterator/go"

	"github.com/mohuishou/scuplus-go/model"
	"github.com/mohuishou/scuplus-go/util/youtu"
)

// CardOCR 一卡通识别,id为失物招领id
func CardOCR(id uint) error {
	var lost model.LostFind
	if err := model.DB().Find(&lost, id).Error; err != nil {
		return err
	}
	rsp, err := youtu.OCR(lost.Pictures)
	if err != nil {
		return err
	}
	if rsp.ErrorCode != 0 {
		return errors.New(rsp.ErrorMsg)
	}
	var cardInfo map[string]string
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
	str, err := jsoniter.MarshalToString(&cardInfo)
	if err != nil {
		return err
	}
	if err := model.DB().Model(&lost).Update("card_info", str).Error; err != nil {
		return err
	}
	return nil
}
