package controller

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"sync"

	"stockfoilo_test/internal/config"
	"stockfoilo_test/internal/consts"
	"stockfoilo_test/internal/model"
	"stockfoilo_test/internal/service"
	"stockfoilo_test/internal/utils"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

func UploadVideo(c *gin.Context) {

	ctx := c.Request.Context()

	var videos []model.Video

	form, err := c.MultipartForm()
	if err != nil {
		ResponseFailure(c, http.StatusBadRequest)
		return
	}

	appConfig := config.GetConfig()

	files := form.File["files"]

	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make(chan error, len(files))

	for _, file := range files {
		wg.Add(1)

		// 고루틴으로 파일 업로드 처리
		go func(file *multipart.FileHeader) {
			defer wg.Done()

			// 파일 확장자 확인
			splitFile := strings.Split(file.Filename, ".")

			fileExtension := splitFile[len(splitFile)-1]
			// 비디오 파일이 아닌 경우 처리
			if !isVideoFile(fileExtension) {
				log.Error().Msgf("not video file: %s", file.Filename)
				results <- fmt.Errorf("not video file: %s", file.Filename)
				return
			}

			uuid := uuid.New().String()
			filePath := fmt.Sprintf("%s/%s.%s", appConfig.FileConfig.VideoPath, uuid, fileExtension)
			video := model.Video{
				Id:         uuid,
				VideoName:  file.Filename,
				Extension:  fileExtension,
				UploadTime: utils.GetTime(),
				Path:       filePath,
			}
			mu.Lock()
			videos = append(videos, video)
			mu.Unlock()

			if err := c.SaveUploadedFile(file, filePath); err != nil {
				log.Error().Msgf("upload file error: %s", err.Error())
				results <- err
				return
			}
			results <- nil
		}(file)
	}

	wg.Wait()
	close(results)

	// 결과를 확인하여 오류가 발생하면 처리
	for err := range results {
		if err != nil {
			removeVideo(c, videos)
			ResponseFailure(c, http.StatusInternalServerError)
			return
		}
	}

	err = service.UploadVideo(ctx, videos)
	if err != nil {
		removeVideo(c, videos)
		ResponseFailure(c, http.StatusInternalServerError)
		return
	}

	ResponseSuccess(c)
}

func ModifyVideo(c *gin.Context) {
	var Request model.ModifyVideo
	if err := c.ShouldBindJSON(&Request); err != nil {
		log.Error().Msgf("bind json error: %s", err.Error())
		ResponseFailure(c, http.StatusBadRequest)
		return
	}

	// 유저가 trim만 요청한 경우
	if Request.IsTrimed && !Request.IsConcated {
		_, err := trimVideo(c, Request.TrimVideoList)
		if err != nil {
			ResponseFailure(c, consts.TrimFail)
			return
		}
	}
	// 유저가 concat만 요청한 경우
	if Request.IsConcated && !Request.IsTrimed {
		err := concatVideo(c, Request.ConcatVideoIdList)
		if err != nil {
			ResponseFailure(c, consts.ConcatFail)
			return
		}
	}

	// 유저가 trim, concat 둘다 요청한 경우
	if Request.IsTrimed && Request.IsConcated {
		trimmedHistoryList, err := trimVideo(c, Request.TrimVideoList)
		if err != nil {
			ResponseFailure(c, consts.TrimFail)
			return
		}
		concatVideoIdList := []string{}
		for _, trimHistory := range trimmedHistoryList {
			concatVideoIdList = append(concatVideoIdList, trimHistory.Video.Id)
		}

		err = concatVideo(c, concatVideoIdList)
		if err != nil {
			ResponseFailure(c, consts.TrimSuccessConcatFail)
			return
		}
	}

	ResponseSuccess(c)
}

func GetVideoInfoList(c *gin.Context) {
	videos, err := service.GetVideoInfoList()
	if err != nil {
		ResponseFailure(c, http.StatusInternalServerError)
		return
	}

	ResponseWithData(c, videos)
}

func GetDownloadVideo(c *gin.Context) {
	ctx := c.Request.Context()

	videoId := c.Param("uuid")

	log.Debug().Msgf("videoId: %s", videoId)

	video, err := service.GetVideoWithVideoId(ctx, videoId)
	if err != nil {
		ResponseFailure(c, http.StatusInternalServerError)
		return
	}

	filePath := video.Path
	// 파일 다운로드 , 파일명은 video.VideoName으로 설정
	c.Header("Content-Disposition", "attachment; filename="+video.VideoName)
	c.File(filePath)

}

func isVideoFile(filename string) bool {
	videoExtensions := []string{"mp4", "avi", "mov", "mkv", "flv", "wmv"}
	for _, ext := range videoExtensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}

func removeVideo(c *gin.Context, videos []model.Video) {

	if len(videos) > 0 {
		for _, video := range videos {
			err := os.Remove(video.Path)
			if err != nil {
				log.Error().Msgf("remove file error: %s", err.Error())
			}
		}
	}
}

func trimVideo(ctx context.Context, trimVideoList []model.TrimVideo) ([]model.TrimHistory, error) {

	trimHistory, err := service.TrimVideo(ctx, trimVideoList)
	if err != nil {
		return trimHistory, err
	}

	return trimHistory, nil
}

func concatVideo(c *gin.Context, concatVideoIdList []string) error {
	ctx := c.Request.Context()

	var videos []model.Video

	for _, concatVideoId := range concatVideoIdList {
		video, err := service.GetVideoWithVideoId(ctx, concatVideoId)
		if err != nil {
			return err
		}
		videos = append(videos, video)
	}

	err := service.ConcatVideo(ctx, videos, concatVideoIdList)
	if err != nil {
		return err
	}

	return nil
}
