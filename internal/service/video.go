package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"stockfoilo_test/internal/config"
	"stockfoilo_test/internal/db"
	"stockfoilo_test/internal/model"
	"stockfoilo_test/internal/utils"

	"stockfoilo_test/internal/repo"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func UploadVideo(videos []model.Video) error {

	log.Debug().Msgf("videos: %v", videos)

	dbCtx := db.GetDbConnection(context.Background())

	err := dbCtx.BeginTxn()
	if err != nil {
		log.Debug().Msgf("failed to begin transaction: %s", err.Error())
		return err
	}

	defer func() {
		if err != nil {
			dbCtx.Rollback()
			return
		}
	}()

	for _, video := range videos {
		err := repo.InsertVideo(dbCtx, video)

		if err != nil {
			return err
		}
	}

	err = dbCtx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func TrimVideo(video model.Video, startTrim, endTrim int) (model.Video, error) {

	log.Debug().Msgf("video: %v, startTrim: %d, endTrim: %d", video, startTrim, endTrim)

	config := config.GetConfig()
	u := uuid.New().String()

	outputFilePath := fmt.Sprintf("%s/%s.%s", config.FileConfig.VideoPath, u, video.Extension)
	// FFmpeg 명령어를 구성
	cmd := exec.Command("ffmpeg",
		"-i", video.Path,
		"-ss", fmt.Sprintf("00:00:%02d", startTrim),
		"-to", fmt.Sprintf("00:00:%02d", endTrim),
		"-c", "copy",
		outputFilePath,
	)

	// 구성된 명령어 실행
	err := cmd.Run()
	if err != nil {
		os.Remove(outputFilePath)
		return model.Video{}, fmt.Errorf("failed to trim video: %w", err)
	}

	trimVideo := model.Video{
		Id:        u,
		VideoName: "trim_" + video.VideoName,
		Extension: video.Extension,
		IsTrimed:  true,
		TrimTime:  utils.GetTime(),
		Path:      outputFilePath,
	}

	trimHistory := model.TrimVideo{
		VideoId:   u,
		StartTime: startTrim,
		EndTime:   endTrim,
	}

	// DB에 trim 정보 업데이트
	dbCtx := db.GetDbConnection(context.Background())
	err = repo.InsertVideo(dbCtx, trimVideo)
	if err != nil {
		os.Remove(outputFilePath)
		return model.Video{}, err
	}

	// DB에 trim history 정보 업데이트
	err = repo.InsertTrimHistory(dbCtx, video.Id, trimHistory)
	if err != nil {
		os.Remove(outputFilePath)
		// DB에 trim history 정보 업데이트 실패 시, DB에 저장된 trim video 정보 삭제
		repo.DeleteVideo(dbCtx, trimVideo.Id)
		return model.Video{}, err
	}

	return trimVideo, nil
}

func ConcatVideo(videos []model.Video, concatVideoList []string) error {

	config := config.GetConfig()
	u := uuid.New().String()
	inputFilePathsToMp4 := []string{}

	// 확장자 통일을 위해 확장자가 mp4가 아닌 경우 mp4로 변환, 해상도도 통일
	for _, video := range videos {

		MP4path, err := ConvertToMp4(video.Path)
		if err != nil {
			return err
		}

		inputFilePathsToMp4 = append(inputFilePathsToMp4, MP4path)
	}

	// concat을 위한 input file list 생성
	inputFileList := fmt.Sprintf("%s/concat_%s.txt", config.FileConfig.ConcatPath, u)
	fileContent := ""
	for _, path := range inputFilePathsToMp4 {
		fileContent += fmt.Sprintf("file '%s'\n", path)
	}

	err := os.MkdirAll(config.FileConfig.ConcatPath, 0755)
	if err != nil {
		log.Error().Msgf("failed to ensure directory: %s", err.Error())
		return err
	}
	err = os.WriteFile(inputFileList, []byte(fileContent), 0644)
	if err != nil {
		log.Error().Msgf("failed to write input file list: %s", err.Error())
		return err
	}

	outputFilePath := fmt.Sprintf("%s/%s.mp4", config.FileConfig.VideoPath, u)

	//파일방식은 재인코딩을 하지 않아서 속도가 빠르지만, 코덱이 다를 경우 문제가 발생할 수 있음
	// args := []string{"-f", "concat", "-safe", "0", "-i", inputFileList, "-c", "copy", outputFilePath}
	filterComplex := ""
	for i := range inputFilePathsToMp4 {
		filterComplex += fmt.Sprintf("[%d:v] [%d:a] ", i, i)
	}
	args := []string{}
	filterComplex += fmt.Sprintf("concat=n=%d:v=1:a=1 [outv] [outa]", len(inputFilePathsToMp4))

	for _, v := range inputFilePathsToMp4 {
		args = append(args, "-i", v)
	}
	args = append(args, "-filter_complex", filterComplex, "-map", "[outv]", "-map", "[outa]", outputFilePath)

	log.Debug().Msgf("args: %v", args)
	cmd := exec.Command("ffmpeg", args...)

	// 구성된 명령어 실행
	output, err := cmd.CombinedOutput()
	log.Debug().Msgf("output: %s", output)

	if err != nil {
		log.Error().Msgf("failed to concat video: %s", err.Error())
		for _, outputFilePath := range inputFilePathsToMp4 {
			log.Debug().Msgf("remove video: %s", outputFilePath)
			err := removeVideo(outputFilePath)
			if err != nil {
				log.Error().Msgf("failed to remove video: %s", err.Error())
				return err
			}
		}
		os.Remove(inputFileList)

		return nil
	}

	concatVideo := model.Video{
		Id:         u,
		VideoName:  "concat_" + utils.GetTime() + "." + "mp4",
		Extension:  "mp4",
		IsConcated: true,
		ConcatTime: utils.GetTime(),
		Path:       outputFilePath,
	}

	dbCtx := db.GetDbConnection(context.Background())
	err = repo.InsertVideo(dbCtx, concatVideo)
	if err != nil {
		err := removeVideo(outputFilePath)
		if err != nil {
			return err
		}
		os.Remove(inputFileList)
		return err
	}
	err = repo.InsertConcatHistory(dbCtx, concatVideo.Id, inputFileList)
	if err != nil {
		for _, outputFilePath := range inputFilePathsToMp4 {
			err := removeVideo(outputFilePath)
			if err != nil {
				return err
			}
			os.Remove(inputFileList)
		}
		return err
	}

	return nil

}

func GetVideoWithVideoId(videoId string) (model.Video, error) {
	dbCtx := db.GetDbConnection(context.Background())
	video, err := repo.FetchVideoWithVideoId(dbCtx, videoId)
	if err != nil {
		return model.Video{}, err
	}

	return video, nil
}

func GetVideoInfoList() ([]model.Video, error) {
	dbCtx := db.GetDbConnection(context.Background())
	videos, err := repo.FetchVideoInfoList(dbCtx)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func ConvertToMp4(inputFilePath string) (string, error) {
	config := config.GetConfig()
	u := uuid.New().String()

	outputFilePath := fmt.Sprintf("%s/%s.mp4", config.FileConfig.VideoPath, u)
	// 비디오 코덱이랑 해상도 맞춰주기....
	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-c:a", "copy", "-c:v", "libx264", "-filter:v", "scale=1080:1920", outputFilePath)

	output, err := cmd.CombinedOutput()
	log.Debug().Msgf("output: %s", output)
	if err != nil {
		return "", fmt.Errorf("failed to convert video: %w", err)
	}

	return outputFilePath, nil
}

func removeVideo(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		log.Error().Msgf("failed to remove video: %s", err.Error())
		return err
	}
	return nil
}
