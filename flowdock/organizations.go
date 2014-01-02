package flowdock

// Organization represents a Flowdock organization to which members belong
type Organization struct {
	Id                int64  `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	ParameterizedName string `json:"paramererized_name,omitempty"`
	UserLimit         int64  `json:"user_limit,omitempty"`
	UserCount         int64  `json:"user_count,omitempty"`
	Active            bool   `json:"active,omitempty"`
	Url               string `json:"url,omitempty"`
}
