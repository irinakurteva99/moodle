package moodle

import (
	"context"
	"github.com/k-yomo/moodle/pkg/urlutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type GradeAPI interface {
	GetGradeItems(ctx context.Context, userID int, courseID int) ([]*UserGrade, error)
}

type gradeAPI struct {
	httpClient *http.Client
	apiURL     *url.URL
}

func newGradeAPI(httpClient *http.Client, apiURL *url.URL) *gradeAPI {
	return &gradeAPI{
		httpClient: httpClient,
		apiURL:     apiURL,
	}
}

type userGradeResponse struct {
	CourseID     int                  `json:"courseid"`
	UserID       int                  `json:"userid"`
	UserFullname string               `json:"userfullname"`
	MaxDepth     int                  `json:"maxdepth"`
	GradeItems   []*gradeItemResponse `json:"gradeitems"`
}

type gradeItemResponse struct {
	ID                     int     `json:"id"`
	ItemName               string  `json:"itemname"`
	ItemType               string  `json:"itemtype"`
	ItemModule             *string `json:"itemmodule"`
	ItemInstance           int     `json:"iteminstance"`
	ItemNumber             *int    `json:"itemnumber"`
	CategoryID             *int    `json:"categoryid"`
	OutcomeID              *int    `json:"outcomeid"`
	ScaleID                *int    `json:"scaleid"`
	Locked                 *bool   `json:"locked"`
	CmID                   int     `json:"cmid"`
	GradeRaw               float64 `json:"graderaw"`
	GradeDateSubmittedUnix *int64  `json:"gradedatesubmitted"`
	GradeDateGradedUnix    int64   `json:"gradedategraded"`
	GradeHiddenByDate      bool    `json:"gradehiddenbydate"`
	GradeNeedsUpdate       bool    `json:"gradeneedsupdate"`
	GradeIsHidden          bool    `json:"gradeishidden"`
	GradeIsLocked          *bool   `json:"gradeislocked"`
	GradeIsOverridden      *bool   `json:"gradeisoverridden"`
	GradeFormatted         string  `json:"gradeformatted"`
	GradeMin               int     `json:"grademin"`
	GradeMax               int     `json:"grademax"`
	RangeFormatted         string  `json:"rangeformatted"`
	Feedback               string  `json:"feedback"`
	FeedbackFormat         int     `json:"feedbackformat"`
}

type getGradeItems struct {
	UserGrades []*userGradeResponse `json:"usergrades"`
	Warnings   Warnings             `json:"warnings"`
}

func (g *gradeAPI) GetGradeItems(ctx context.Context, userID int, courseID int) ([]*UserGrade, error) {
	u := urlutil.CopyWithQueries(g.apiURL, map[string]string{
		"wsfunction": "gradereport_user_get_grade_items",
		"userid":     strconv.Itoa(userID),
		"courseid":   strconv.Itoa(courseID),
	})
	res := getGradeItems{}
	if err := getAndUnmarshal(ctx, g.httpClient, u, &res); err != nil {
		return nil, err
	}
	if len(res.Warnings) > 0 {
		return nil, res.Warnings
	}

	return mapFromUserGradeListResponse(res.UserGrades), nil
}

func mapFromUserGradeListResponse(userGradeResList []*userGradeResponse) []*UserGrade {
	userGrades := make([]*UserGrade, 0, len(userGradeResList))
	for _, gradeItemRes := range userGradeResList {
		userGrades = append(userGrades, mapFromUserGradeResponse(gradeItemRes))
	}
	return userGrades
}

func mapFromUserGradeResponse(userGradeRes *userGradeResponse) *UserGrade {
	return &UserGrade{
		CourseID:     userGradeRes.CourseID,
		UserID:       userGradeRes.UserID,
		UserFullname: userGradeRes.UserFullname,
		MaxDepth:     userGradeRes.MaxDepth,
		GradeItems:   mapFromGradeItemListResponse(userGradeRes.GradeItems),
	}
}

func mapFromGradeItemListResponse(gradeItemResList []*gradeItemResponse) []*GradeItem {
	gradeItems := make([]*GradeItem, 0, len(gradeItemResList))
	for _, gradeItemRes := range gradeItemResList {
		gradeItems = append(gradeItems, mapFromGradeItemResponse(gradeItemRes))
	}
	return gradeItems
}

func mapFromGradeItemResponse(gradeItemRes *gradeItemResponse) *GradeItem {
	var gradeDateSubmitted, gradeDateGraded *time.Time
	if gradeItemRes.GradeDateSubmittedUnix != nil {
		t := time.Unix(*gradeItemRes.GradeDateSubmittedUnix, 0)
		gradeDateSubmitted = &t
	}
	if gradeItemRes.GradeDateGradedUnix > 0 {
		t := time.Unix(gradeItemRes.GradeDateGradedUnix, 0)
		gradeDateGraded = &t
	}
	return &GradeItem{
		ID:                 gradeItemRes.ID,
		ItemName:           gradeItemRes.ItemName,
		ItemType:           gradeItemRes.ItemType,
		ItemModule:         gradeItemRes.ItemModule,
		ItemInstance:       gradeItemRes.ItemInstance,
		ItemNumber:         gradeItemRes.ItemNumber,
		CategoryID:         gradeItemRes.CategoryID,
		OutcomeID:          gradeItemRes.OutcomeID,
		ScaleID:            gradeItemRes.ScaleID,
		Locked:             gradeItemRes.Locked,
		CmID:               gradeItemRes.CmID,
		GradeRaw:           gradeItemRes.GradeRaw,
		GradeDateSubmitted: gradeDateSubmitted,
		GradeDateGraded:    gradeDateGraded,
		GradeHiddenByDate:  gradeItemRes.GradeHiddenByDate,
		GradeNeedsUpdate:   gradeItemRes.GradeNeedsUpdate,
		GradeIsHidden:      gradeItemRes.GradeIsHidden,
		GradeIsLocked:      gradeItemRes.GradeIsLocked,
		GradeIsOverridden:  gradeItemRes.GradeIsOverridden,
		GradeFormatted:     gradeItemRes.GradeFormatted,
		GradeMin:           gradeItemRes.GradeMin,
		GradeMax:           gradeItemRes.GradeMax,
		RangeFormatted:     gradeItemRes.RangeFormatted,
		Feedback:           gradeItemRes.Feedback,
		FeedbackFormat:     gradeItemRes.FeedbackFormat,
	}
}