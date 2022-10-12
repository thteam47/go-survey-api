package models

type TokenInfo struct {
	AuthenticationDone bool          `json:"authentication_done,omitempty" bson:"authentication_done,omitempty"`
	UserId             string        `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Subject            string        `json:"subject,omitempty" bson:"subject,omitempty"`
	Exp                int32         `json:"exp,omitempty" bson:"exp,omitempty"`
	Role               string        `json:"role,omitempty" bson:"role,omitempty"`
	PermissionAll      bool          `json:"permission_all,omitempty" bson:"permission_all,omitempty"`
	Permissions        []*Permission `json:"permissions,omitempty" bson:"permissions,omitempty"`
}
