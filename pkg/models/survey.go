package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Survey struct {
	ID           primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	UserIdCreate string             `json:"user_id_create,omitempty" bson:"user_id_create,omitempty"`
	UserIdVerify string             `json:"user_id_verify,omitempty" bson:"user_id_verify,omitempty"`
	UserIdJoin   string             `json:"user_id_join,omitempty" bson:"user_id_join,omitempty"`
	Contents     []*Content         `json:"contents,omitempty" bson:"contents,omitempty"`
	Status       string             `json:"status,omitempty" bson:"status,omitempty"`
	SurveyId     string             `json:"survey_id,omitempty" bson:"survey_id,omitempty"`
	CreateTime   int32              `json:"create_time" bson:"create_time,omitempty"`
	UpdateTime   int32              `json:"update_time" bson:"update_time,omitempty"`
}

type Content struct {
	Question string   `json:"question,omitempty" bson:"question,omitempty"`
	Answers  []string `json:"answers,omitempty" bson:"answers,omitempty"`
	Choose   string     `json:"choose,omitempty" bson:"choose,omitempty"`
}
