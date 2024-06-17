package model

type Video struct {
	Id         string `json:"id"`
	Path       string `json:"path"`
	VideoName  string `json:"videoName"`
	Extension  string `json:"extension"`
	UploadTime string `json:"uploadTime"`
	IsTrimed   bool   `json:"isTrimed"`
	TrimTime   string `json:"trimTime"`
	IsConcated bool   `json:"isConcated"`
	ConcatTime string `json:"concatTime"`
	IsEncoded  bool   `json:"isEncoded"`
	EncodeTime string `json:"encodeTime"`
}

type ModifyVideo struct {
	IsTrimed          bool        `json:"isTrimed"`
	IsConcated        bool        `json:"isConcated"`
	TrimVideoList     []TrimVideo `json:"trimVideoList"`
	ConcatVideoIdList []string    `json:"concatVideoList"`
}

type TrimVideo struct {
	VideoId   string `json:"videoId"`
	VideoPath string `json:"videoPath"`
	StartTime int    `json:"startTime"`
	EndTime   int    `json:"endTime"`
}

type TrimHistory struct {
	Video    Video     `json:"videoList"`
	TrimInfo TrimVideo `json:"trimInfo"`
}

type VideoInfo struct {
	Video          Video     `json:"video"`
	TrimInfo       TrimVideo `json:"trimInfo"`
	ConcatInfoPath string    `json:"concatInfoPath"`
}
