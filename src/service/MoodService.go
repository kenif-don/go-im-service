package service

import (
	"go-im-service/src/configs/conf"
	utils "go-im-service/src/configs/err"
	"go-im-service/src/configs/log"
	"go-im-service/src/dto"
)

type MoodService struct {
}

func NewMoodService() *MoodService {
	return &MoodService{}
}

func (_service *MoodService) AddMood(tp int, content, urls string) *utils.Error {
	_, err := Post("/api/mood/add", map[string]interface{}{"type": tp, "content": content, "urls": urls})
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
	if resultDTO.Data == nil {
		return "", nil
	}
	return resultDTO.Data.(string), nil
}

// AddReply 添加动态回复
func (_service *MoodService) AddReply(moodId, replyUserId uint64, content string) *utils.Error {
	var err *utils.Error
	//没有回复ID 或者回复的自己都算直接评论
	if replyUserId == 0 || replyUserId == conf.GetLoginInfo().User.Id {
		_, err = Post("/api/reply/add", map[string]interface{}{"moodId": moodId, "content": content})
	} else {
		_, err = Post("/api/reply/add", map[string]interface{}{"moodId": moodId, "content": content, "replyUserId": replyUserId})
	}
	if err != nil {
		return log.WithError(err)
	}
	return nil
}
