package testing

import (
	"fmt"
	"os"
	"testing"

	"github.com/rs/zerolog/log"
	ffmpeg "github.com/u2takey/ffmpeg-go"

	"github.com/google/uuid"
)

func Test(t *testing.T) {
	testNum := 3

	switch testNum {
	case 1:
		convert(t)
	case 2:
		trim(t)
	case 3:
		concat(t)
	}
}

func convert(t *testing.T) {

	videoPath := "/Users/hongjunho/Downloads/workspace/stockfolio-test/cmd/app/stockFoilo/videos"

	inputFilePath := "/Users/hongjunho/Downloads/workspace/stockfolio-test/cmd/app/stockFoilo/videos/4074a540-f9b7-415d-b47c-12c444282e9b.mp4"

	// TestConvert is a test function for Convert function
	u := uuid.New().String()

	outputFilePath := fmt.Sprintf("%s/%s.mp4", videoPath, u)

	err := ffmpeg.Input(inputFilePath).Output(outputFilePath, ffmpeg.KwArgs{"c": "copy"}).OverWriteOutput().Run()

	if err != nil {
		t.Errorf("TestConvert:: error: %s", err.Error())
		return
	}

	t.Logf("TestConvert:: success")
}

func trim(t *testing.T) {
	videoPath := "/Users/hongjunho/Downloads/workspace/stockfolio-test/cmd/app/stockFoilo/videos"

	inputFilePath := "/Users/hongjunho/Downloads/workspace/stockfolio-test/cmd/app/stockFoilo/videos/89c8aae5-523b-44e0-9fc0-886c7816b048.mp4"

	u := uuid.New().String()

	outputFilePath := fmt.Sprintf("%s/%s.mp4", videoPath, u)

	startTime := 2
	endTime := 5

	err := ffmpeg.Input(inputFilePath, ffmpeg.KwArgs{"ss": startTime}).
		Output(outputFilePath, ffmpeg.KwArgs{"c": "copy", "to": endTime - startTime}).OverWriteOutput().Run()
	if err != nil {
		t.Errorf("TestTrim:: error: %s", err.Error())
		return
	}

	t.Logf("TestTrim:: success")
}

func concat(t *testing.T) {
	videoPath := "/Users/hongjunho/Downloads/workspace/stockfolio-test/cmd/app/stockFoilo/videos"
	concatPath := "/Users/hongjunho/Downloads/workspace/stockfolio-test/cmd/app/stockFoilo/concat"

	inputFilePaths := []string{}

	inputFilePath2 := "/Users/hongjunho/Downloads/workspace/stockfolio-test/cmd/app/stockFoilo/videos/a0c18925-5cf3-4fd3-b725-be952c23b23e.mp4"
	inputFilePath1 := "/Users/hongjunho/Downloads/workspace/stockfolio-test/cmd/app/stockFoilo/videos/224ed1c3-77d8-4ed7-add7-081d1126a1b3.mp4"

	inputFilePaths = append(inputFilePaths, inputFilePath2)
	inputFilePaths = append(inputFilePaths, inputFilePath1)

	u := uuid.New().String()

	// concat을 위한 input file list 생성
	inputFileList := fmt.Sprintf("%s/concat_%s.txt", concatPath, u)
	fileContent := ""
	for _, path := range inputFilePaths {
		fileContent += fmt.Sprintf("file '%s'\n", path)
	}

	err := os.MkdirAll(concatPath, 0755)
	if err != nil {
		log.Error().Msgf("failed to ensure directory: %s", err.Error())
		return
	}
	err = os.WriteFile(inputFileList, []byte(fileContent), 0644)
	if err != nil {
		log.Error().Msgf("failed to write input file list: %s", err.Error())
		return
	}

	outputFilePath := fmt.Sprintf("%s/%s.mp4", videoPath, u)

	filterComplex := ""
	for i := range inputFilePaths {
		filterComplex += fmt.Sprintf("[%d:v][%d:a]", i, i)
	}
	filterComplex += fmt.Sprintf("concat=n=%d:v=1:a=1[outv][outa]", len(inputFilePaths))

	err = ffmpeg.Input(inputFileList, ffmpeg.KwArgs{"f": "concat", "safe": "0"}).
		Output(outputFilePath, ffmpeg.KwArgs{"c": "copy"}).
		OverWriteOutput(). // 동일한 해상도, 프레임 레이트, 오디오 샘플 레이트
		Run()
	if err != nil {
		log.Error().Msgf("failed to concat video: %v", err)
		return
	}

	t.Logf("TestConcat:: success")
}
