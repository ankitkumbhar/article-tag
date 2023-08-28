package constant

var TagError = map[string]interface{}{
	"Username":    "field is required",
	"Publication": "field is required, and must be a valid publications",
	"Tags":        "atleast one tag is required",
	"TagID":       "field is required and must have a numeric format",
	"TagName":     "field is required",
	"Order":       "invalid order field, should be either createdatdesc, createdatasc or tagname",
}
