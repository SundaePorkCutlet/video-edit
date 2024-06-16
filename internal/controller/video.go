package controller

import (
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

	var videos []model.Video

	form, err := c.MultipartForm()
	if err != nil {
		ResponseFailure(c, http.StatusBadRequest)
		return
	}

	appConfig := config.GetConfig()

	files := form.File["files"]

	var wg sync.WaitGroup
	results := make(chan error, len(files))

	for _, file := range files {
		wg.Add(1)

		go func(file *multipart.FileHeader) {
			defer wg.Done()

			splitFile := strings.Split(file.Filename, ".")

			fileExtension := splitFile[len(splitFile)-1]
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

			videos = append(videos, video)

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

	err = service.UploadVideo(videos)
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
		err, errorCode, _ := trimVideo(c, Request.TrimVideoList)
		if err != nil {
			ResponseFailure(c, errorCode)
			return
		}
	}
	// 유저가 concat만 요청한 경우
	if Request.IsConcated && !Request.IsTrimed {
		errorCode, err := concatVideo(c, Request.ConcatVideoIdList)
		if err != nil {
			ResponseFailure(c, errorCode)
			return
		}
	}

	// 유저가 trim, concat 둘다 요청한 경우
	if Request.IsTrimed && Request.IsConcated {
		trimmedVideoList, errorCode, err := trimVideo(c, Request.TrimVideoList)
		if err != nil {
			ResponseFailure(c, errorCode)
			return
		}
		concatVideoIdList := []string{}
		for _, video := range trimmedVideoList {
			concatVideoIdList = append(concatVideoIdList, video.Id)
		}

		_, err = concatVideo(c, concatVideoIdList)
		if err != nil {
			ResponseFailure(c, consts.SuccessTrimFailConcat)
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

func trimVideo(c *gin.Context, trimVideoList []model.TrimVideo) ([]model.Video, int, error) {
	var videos []model.Video

	var wg sync.WaitGroup
	results := make(chan model.OperationError, len(trimVideoList))

	for _, trimVideo := range trimVideoList {
		wg.Add(1)

		go func(trimVideo model.TrimVideo) {
			defer wg.Done()
			video, err := service.GetVideoWithVideoId(trimVideo.VideoId)
			if err != nil {
				results <- model.OperationError{Error: err, ErrorCode: http.StatusBadRequest}
				return
			}

			// trim time 체크
			if trimVideo.StartTime < 0 || trimVideo.EndTime < 0 {
				results <- model.OperationError{Error: fmt.Errorf("invalid trim time"), ErrorCode: http.StatusBadRequest}
				return
			}

			trimedVideo, err := service.TrimVideo(video, trimVideo.StartTime, trimVideo.EndTime)
			if err != nil {
				results <- model.OperationError{Error: err, ErrorCode: http.StatusInternalServerError}
				return
			}
			videos = append(videos, trimedVideo)

			results <- model.OperationError{Error: nil, ErrorCode: 0}
		}(trimVideo)
	}

	wg.Wait()
	close(results)
	// 결과를 확인하여 오류가 발생하면 처리
	for result := range results {
		log.Debug().Msgf("VIDOE: %v", videos)
		removeVideo(c, videos)
		return nil, result.ErrorCode, result.Error

	}

	return videos, 0, nil
}

func concatVideo(c *gin.Context, concatVideoIdList []string) (int, error) {
	var videos []model.Video

	for _, concatVideoId := range concatVideoIdList {
		video, err := service.GetVideoWithVideoId(concatVideoId)
		if err != nil {
			return http.StatusBadRequest, err
		}
		videos = append(videos, video)
	}

	err := service.ConcatVideo(videos, concatVideoIdList)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}
