package model

type Response struct{
	Chinese          string   `json:"chinese"`
	Pinyin           string   `json:"pinyin"`
	PinyinNormalized string   `json:"pinyin_normalized"`
	Meanings         []string `json:"meanings"`
}

type Request struct{
	Search string `json:"word"`
}