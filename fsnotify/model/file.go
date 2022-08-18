package model

type File struct {
	Name string `json:"name,omitempty"`
	Size int64  `json:"size,omitempty"`
	Url  string `json:"url,omitempty"`
	MD5  string `json:"md_5,omitempty"`
}
