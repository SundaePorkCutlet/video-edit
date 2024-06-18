package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"stockfoilo_test/internal/config"
	"stockfoilo_test/internal/db"
	"stockfoilo_test/internal/model"
	"stockfoilo_test/internal/utils"
	"strings"
	"sync"

	"stockfoilo_test/internal/repo"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func UploadVideo(ctx context.Context, videos []model.Video) error {

	dbCtx := db.GetDbConnection(ctx)

	return dbCtx.Transactional(func(dbCtx *db.DbCtx) error {
		for _, video := range videos {
			err := repo.InsertVideo(dbCtx.Tx, dbCtx.Ctx, video)

			if err != nil {
				return err
			}
		}
		return nil
	})

}

func TrimVideo(ctx context.Context, trimVideoList []model.TrimVideo) ([]model.TrimHistory, error) {

	config := config.GetConfig()

	dbCtx := db.GetDbConnection(ctx)

	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make(chan error, len(trimVideoList))

	var trimHistoryList []model.TrimHistory

	for _, trimVideo := range trimVideoList {
		wg.Add(1)

		// 고루틴으로 비디오 트림 처리
		go func(trimVideo model.TrimVideo) {
			defer wg.Done()
			video, err := repo.FetchVideoWithVideoId(dbCtx, trimVideo.VideoId)
			if err != nil {
				results <- err
				return
			}

			// trim time 체크
			if trimVideo.StartTime < 0 || trimVideo.EndTime < 0 {
				results <- fmt.Errorf("invalid trim time")
				return
			}

			u := uuid.New().String()

			outputFilePath := fmt.Sprintf("%s/%s.%s", config.FileConfig.VideoPath, u, video.Extension)
			// FFmpeg 명령어를 구성
			cmd := exec.Command("ffmpeg",
				"-i", video.Path,
				"-ss", fmt.Sprintf("00:00:%02d", trimVideo.StartTime),
				"-to", fmt.Sprintf("00:00:%02d", trimVideo.EndTime),
				"-c", "copy",
				outputFilePath,
			)

			// 구성된 명령어 실행
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Error().Msgf("failed to trim video: %s , %s", err.Error(), output)
				results <- err
				return
			}
			mu.Lock()

			defer mu.Unlock()
			// DB에 trim 정보 업데이트
			trimmedVideo := model.Video{
				Id:        u,
				VideoName: "tr-" + video.VideoName,
				Extension: video.Extension,
				IsTrimed:  true,
				TrimTime:  utils.GetTime(),
				Path:      outputFilePath,
			}

			trimInfo := model.TrimVideo{
				VideoId:   video.Id,
				StartTime: trimVideo.StartTime,
				EndTime:   trimVideo.EndTime,
			}

			trimHistory := model.TrimHistory{
				Video:    trimmedVideo,
				TrimInfo: trimInfo,
			}

			trimHistoryList = append(trimHistoryList, trimHistory)

			results <- nil

		}(trimVideo)
	}

	wg.Wait()
	close(results)

	// 결과를 확인하여 오류가 발생하면 생성된 비디오 삭제
	for err := range results {
		if err != nil {
			for _, trimmedVideo := range trimHistoryList {
				err := removeVideo(trimmedVideo.Video.Path)
				if err != nil {
					return trimHistoryList, err
				}
			}
			return trimHistoryList, err
		}
	}

	//트랜잭션 시작
	err := dbCtx.Transactional(func(dc *db.DbCtx) error {

		for _, trimmedVideo := range trimHistoryList {
			err := repo.InsertVideo(dbCtx.Tx, dbCtx.Ctx, trimmedVideo.Video)
			if err != nil {
				return err
			}
			// DB에 trim history 정보 업데이트
			err = repo.InsertTrimHistory(dbCtx.Tx, dbCtx.Ctx, trimmedVideo.Video.Id, trimmedVideo.TrimInfo)
			if err != nil {
				return err
			}
		}
		return nil
	})

	//트랜잭션 에러 발생시 생성된 비디오 삭제
	if err != nil {
		for _, trimmedVideo := range trimHistoryList {
			log.Debug().Msgf("remove video: %s", trimmedVideo.Video.Path)
			err := removeVideo(trimmedVideo.Video.Path)
			if err != nil {
				return trimHistoryList, err
			}
		}
		return trimHistoryList, err
	}

	return trimHistoryList, nil
}

func ConcatVideo(ctx context.Context, videos []model.Video, concatVideoList []string) error {

	config := config.GetConfig()
	u := uuid.New().String()
	inputFilePathsToMp4 := []string{}

	dbCtx := db.GetDbConnection(context.Background())

	// 확장자 통일을 위해 확장자가 mp4가 아닌 경우 mp4로 변환, 해상도도 통일
	err := dbCtx.Transactional(func(dbCtx *db.DbCtx) error {
		for _, video := range videos {
			MP4path, err := ConvertToMp4(dbCtx.Tx, dbCtx.Ctx, video.Path, video.Id, video.VideoName)
			if err != nil {
				return err
			}
			inputFilePathsToMp4 = append(inputFilePathsToMp4, MP4path)
		}
		return nil
	})
	if err != nil {
		for _, outputFilePath := range inputFilePathsToMp4 {
			err := removeVideo(outputFilePath)
			if err != nil {
				return err
			}
		}
		return err
	}

	// concat정보 저장을 위해서 file 생성
	inputFileList := fmt.Sprintf("%s/concat_%s.txt", config.FileConfig.ConcatPath, u)
	fileContent := ""
	for _, path := range inputFilePathsToMp4 {
		fileContent += fmt.Sprintf("file '%s'\n", path)
	}

	err = os.MkdirAll(config.FileConfig.ConcatPath, 0755)
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

	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error().Msgf("failed to concat video: %s, %s", err.Error(), output)
		for _, outputFilePath := range inputFilePathsToMp4 {
			log.Debug().Msgf("remove video: %s", outputFilePath)
			err := removeVideo(outputFilePath)
			if err != nil {
				return err
			}
		}
		os.Remove(inputFileList)

		return err
	}

	concatVideo := model.Video{
		Id:         u,
		VideoName:  "cc-" + utils.GetTime() + "." + "mp4",
		Extension:  "mp4",
		IsConcated: true,
		ConcatTime: utils.GetTime(),
		Path:       outputFilePath,
	}

	//트랜잭션 시작
	err = dbCtx.Transactional(func(dbCtx *db.DbCtx) error {
		err = repo.InsertVideo(dbCtx.Tx, dbCtx.Ctx, concatVideo)
		if err != nil {
			err := removeVideo(outputFilePath)
			if err != nil {
				return err
			}
			err = removeVideo(inputFileList)
			if err != nil {
				return err
			}
			return err
		}
		err = repo.InsertConcatHistory(dbCtx.Tx, dbCtx.Ctx, concatVideo.Id, inputFileList)
		if err != nil {
			for _, outputFilePath := range inputFilePathsToMp4 {
				err := removeVideo(outputFilePath)
				if err != nil {
					return err
				}
				err = removeVideo(inputFileList)
				if err != nil {
					return err
				}
			}
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Msgf("failed to concat video: %s", err.Error())
		return err
	}

	return nil

}

func GetVideoWithVideoId(ctx context.Context, videoId string) (model.Video, error) {
	dbCtx := db.GetDbConnection(ctx)
	video, err := repo.FetchVideoWithVideoId(dbCtx, videoId)
	if err != nil {
		return model.Video{}, err
	}

	return video, nil
}

func GetVideoInfoList() ([]model.VideoInfo, error) {
	dbCtx := db.GetDbConnection(context.Background())
	videoInfo := []model.VideoInfo{}

	videos, err := repo.FetchVideoInfoList(dbCtx)
	if err != nil {
		return nil, err
	}

	for _, video := range videos {
		var trimInfo model.TrimVideo
		var concatInfo string
		var encodeInfo string
		var err error

		// 각 비디오 정보에 대한 trim, concat, encode 정보를 가져옴
		if video.IsTrimed {
			trimInfo, err = repo.FetchTrimInfo(dbCtx, video.Id)
			if err != nil {
				return nil, err
			}
		}
		if video.IsConcated {
			concatInfo, err = repo.FetchConcatInfo(dbCtx, video.Id)
			if err != nil {
				return nil, err
			}
		}
		if video.IsEncoded {
			encodeInfo, err = repo.FetchEncodeInfo(dbCtx, video.Id)
			if err != nil {
				return nil, err
			}
		}

		tempVideoInfo := model.VideoInfo{
			Video:          video,
			TrimInfo:       trimInfo,
			ConcatInfoPath: concatInfo,
			EncodeInfoPath: encodeInfo,
		}

		videoInfo = append(videoInfo, tempVideoInfo)
	}

	return videoInfo, nil
}

func ConvertToMp4(tx *sql.Tx, ctx context.Context, inputFilePath string, originVideoId string, originVideoName string) (string, error) {
	config := config.GetConfig()
	u := uuid.New().String()

	outputFilePath := fmt.Sprintf("%s/%s.mp4", config.FileConfig.VideoPath, u)
	// 비디오 코덱이랑 해상도 맞춰주기....
	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-c:a", "copy", "-c:v", "libx264", "-filter:v", "scale=-1:1080", "-threads", "4", outputFilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Msgf("failed to convert video: %s, %s", err.Error(), output)
		return "", err
	}

	videoName := strings.Split(originVideoName, ".")

	video := model.Video{
		Id:         u,
		VideoName:  "ecd-" + videoName[0] + ".mp4",
		Extension:  "mp4",
		IsTrimed:   false,
		IsConcated: false,
		IsEncoded:  true,
		EncodeTime: utils.GetTime(),
		Path:       outputFilePath,
	}

	err = repo.InsertVideo(tx, ctx, video)
	if err != nil {
		return "", err
	}

	err = repo.InsertEncodeHistory(tx, ctx, u, originVideoId)
	if err != nil {
		return "", err
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
