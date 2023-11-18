package service

import (
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/dto"
)

type MoodService struct {
}

func NewMoodService() *MoodService {
	return &MoodService{}
}

func (_service *MoodService) AddMood(content, urls string) *utils.Error {
	_, err := Post("/api/mood/add", map[string]interface{}{"content": content, "urls": urls})
	if err != nil {
		return log.WithError(err)
	}
	return nil
}
func (_service *MoodService) DeleteMood(id uint64) *utils.Error {
	_, err := Post("/api/mood/delete", map[string]uint64{"id": id})
	if err != nil {
		return log.WithError(err)
	}
	return nil
}

// SelectOneMood 获取单个动态用于回复时更新
func (_service *MoodService) SelectOneMood(id uint64) (string, *utils.Error) {
	resultDTO, err := Post("/api/mood/selectOne", map[string]uint64{"id": id})
	if err != nil {
		return "", log.WithError(err)
	}
	if resultDTO.Data == nil {
		return "", nil
	}
	return resultDTO.Data.(string), nil
}

// PagingMood 分页获取动态 userId可以只看某人
func (_service *MoodService) PagingMood(page, pageSize int, userId uint64) (string, *utils.Error) {
	var resultDTO *dto.ResultDTO
	var err *utils.Error
	if userId == 0 {
		resultDTO, err = Post("/api/mood/paging", map[string]interface{}{"page": page, "pageSize": pageSize})
	} else {
		resultDTO, err = Post("/api/mood/paging", map[string]interface{}{"page": page, "pageSize": pageSize, "userId": userId})
	}
	if err != nil {
		return "", log.WithError(err)
	}
	return resultDTO.Data.(string), nil
}

// AddReply 添加动态回复
func (_service *MoodService) AddReply(moodId, replyUserId uint64, content string) *utils.Error {
	var err *utils.Error
	if replyUserId == 0 {
		_, err = Post("/api/reply/add", map[string]interface{}{"moodId": moodId, "content": content})
	} else {
		_, err = Post("/api/reply/add", map[string]interface{}{"moodId": moodId, "content": content, "replyUserId": replyUserId})
	}
	if err != nil {
		return log.WithError(err)
	}
	return nil
}
