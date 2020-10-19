package apiBuilder

// 模型信息
type structInfo struct {
	Name         string `json:"name"`
	Types        string `json:"types"`
	MapName      string `json:"map_name"`
	XormTags     string `json:"xorm_tags"`
	ValidateTags string `json:"validate_tags"`
	CommentTags  string `json:"comment_tags"`
	AttrTags     string `json:"attr_tags"`
}

type TableFieldsResp struct {
	Fields        []structInfo    `json:"fields"`
	AutoIncrement string          `json:"autoincr"`
	Updated       string          `json:"updated"`
	Deleted       string          `json:"deleted"`
	Created       map[string]bool `json:"created"`
	Version       string          `json:"version"`
}
