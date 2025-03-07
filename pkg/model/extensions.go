package model

import "fmt"

func (r *RatingResult) GetKey() string {
	return fmt.Sprintf("%s%s", r.SubjectName, r.InsightName)
}

func (r *RatingPrompt) GetKey() string {
	return fmt.Sprintf("%s%s", r.SubjectName, r.InsightName)
}
