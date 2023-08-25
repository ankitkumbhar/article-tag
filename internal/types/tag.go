package types

type StoreTagRequest struct {
	Username    string `json:"username"`
	Publication string `json:"publication"`
	Tags        []Tag  `json:"tags"`
}

type GetTagRequest struct {
	Username    string `json:"username"`
	Publication string `json:"publication"`
	Tags        []Tag  `json:"tags"`
	Order       string `json:"order"`
}

type GetTagResponse struct {
	Tags []Tag `json:"tags"`
}

type DeleteTagRequest struct {
	Username    string `json:"username"`
	Publication string `json:"publication"`
	Tags        []Tag  `json:"tags"`
}

type GetPopularTagRequest struct {
	Username    string `json:"username"`
	Publication string `json:"publication"`
}

type Tag struct {
	TagID   string `json:"tag_id"`
	TagName string `json:"tag_name"`
}
