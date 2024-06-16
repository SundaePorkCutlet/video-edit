package model

type Video struct {
	Id         string `json:"id"`
	Path       string `json:"path"`
	VideoName  string `json:"videoName"`
	Extension  string `json:"extension"`
	Size       int    `json:"size"`
	UploadTime string `json:"uploadTime"`
	IsTrimed   bool   `json:"isTrimed"`
	TrimTime   string `json:"trimTime"`
	IsConcated bool   `json:"isConcated"`
	ConcatTime string `json:"concatTime"`
}

type ModifyVideo struct {
	IsTrimed          bool        `json:"isTrimed"`
	IsConcated        bool        `json:"isConcated"`
	TrimVideoList     []TrimVideo `json:"trimVideoList"`
	ConcatVideoIdList []string    `json:"concatVideoList"`
}

type TrimVideo struct {
	VideoId   string `json:"videoId"`
	StartTime int    `json:"startTime"`
	EndTime   int    `json:"endTime"`
}
