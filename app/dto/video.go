package dto

type VideoReq struct {
	//File    interface{} `form:"file" binding:"required"`
	MinTime uint32 `form:"min_time" binding:"omitempty"`
	MaxTime uint32 `form:"max_time" binding:"omitempty"`
	Level   string `form:"level" binding:"omitempty"`
}
