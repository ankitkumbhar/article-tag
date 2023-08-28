package types

type StoreTagRequest struct {
	Username    string `json:"username" validate:"required"`
	Publication string `json:"publication" validate:"required,oneof=RS AK BC ST"`
	Tags        []Tag  `json:"tags" validate:"required,dive"`
}

type GetTagRequest struct {
	Username    string `json:"username" validate:"required"`
	Publication string `json:"publication" validate:"required,oneof=RS AK ST BC"`
	Order       string `json:"order" validate:"omitempty,oneof=createdatdesc createdatasc tagname"`
}

type GetTagResponse struct {
	Tags []Tag `json:"tags"`
}

type DeleteTagRequest struct {
	Username    string `json:"username" validate:"required"`
	Publication string `json:"publication" validate:"required,oneof=RS AK ST BC"`
	Tags        []Tag  `json:"tags" validate:"required,dive"`
}

type GetPopularTagRequest struct {
	Username    string `json:"username"`
	Publication string `json:"publication" validate:"required,oneof=RS AK ST BC"`
}

type Tag struct {
	TagID   string `json:"tag_id" validate:"required,numeric"`
	TagName string `json:"tag_name" validate:"required"`
}
