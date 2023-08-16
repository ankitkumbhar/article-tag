package types

type StoreTagRequest struct {
	Username    string   `json:"username"`
	Publication string   `json:"publication"`
	Tags        []string `json:"tags"`
}

type GetTagRequest struct {
	Username    string   `json:"username"`
	Publication string   `json:"publication"`
	Tags        []string `json:"tags"`
}

type GetTagResponse struct {
	Tags []string `json:"tags"`
}

type DeleteTagRequest struct {
	Username    string   `json:"username"`
	Publication string   `json:"publication"`
	Tags        []string `json:"tags"`
}

type GetPopularTagRequest struct {
	Username    string `json:"username"`
	Publication string `json:"publication"`
}
