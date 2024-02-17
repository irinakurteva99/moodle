package moodle

import (
	"context"
	"strconv"
	"time"
)

type CourseAPI interface {
	GetEnrolledCoursesByTimelineClassification(ctx context.Context, classification CourseClassification) ([]*Course, error)
	GetEnrolledStudentsByCourseID(ctx context.Context, id int) ([]*Student, error)
}

type courseAPI struct {
	*apiClient
}

func newCourseAPI(apiClient *apiClient) *courseAPI {
	return &courseAPI{apiClient}
}

type courseResponse struct {
	ID              int    `json:"id"`
	FullName        string `json:"fullname"`
	ShortName       string `json:"shortname"`
	Summary         string `json:",omitempty"`
	SummaryFormat   int    `json:"summaryformat"`
	StartDateUnix   int64  `json:"startdate"`
	EndDateUnix     int64  `json:"enddate"`
	Visible         bool   `json:"visible"`
	FullNameDisplay string `json:"fullnamedisplay"`
	ViewURL         string `json:"viewurl"`
	CourseImage     string `json:"courseimage"`
	Progress        int    `json:"progress"`
	HasProgress     bool   `json:"hasprogress"`
	IsSavourite     bool   `json:"isfavourite"`
	Hidden          bool   `json:"hidden"`
	ShowShortName   bool   `json:"showshortname"`
	CourseCategory  string `json:"coursecategory"`
}

type getEnrolledCoursesByTimelineClassificationResponse struct {
	Courses    []*courseResponse `json:"courses"`
	NextOffset int               `json:"nextoffset"`
}

type studentResponse struct {
	ID                   int             `json:"id,omitempty"`
	FirstName            string          `json:"firstname,omitempty"`
	LastName             string          `json:"lastname,omitempty"`
	FullName             string          `json:"fullname,omitempty"`
	Email                string          `json:"email,omitempty"`
	IdNumber             string          `json:"idnumber,omitempty"`
	FirstAccess          int             `json:"firstaccess,omitempty"`
	LastAccess           int             `json:"lastaccess,omitempty"`
	LastCourseAccess     int             `json:"lastcourseaccess,omitempty"`
	Description          string          `json:"description,omitempty"`
	DescriptionFormat    int             `json:"descriptionformat,omitempty"`
	City                 string          `json:"city,omitempty"`
	Country              string          `json:"country,omitempty"`
	ProfileImageURLSmall string          `json:"profileimageurlsmall,omitempty"`
	ProfileImageURL      string          `json:"profileimageurl,omitempty"`
	CustomFields         []interface{}   `json:"customfields,omitempty"`
	Groups               []groupResponse `json:"groups,omitempty"`
	Roles                []roleResponse  `json:"roles,omitempty"`
	EnrolledCourses      []interface{}   `json:"enrolledcourses,omitempty"`
}

type roleResponse struct {
	RoleId    int    `json:"roleid,omitempty"`
	Name      string `json:"name,omitempty"`
	Role      string `json:"shortname,omitempty"`
	SortOrder int    `json:"sortorder,omitempty"`
}

type groupResponse struct {
	RoleId            int    `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	Description       string `json:"description,omitempty"`
	Descriptionformat int    `json:"descriptionformat,omitempty"`
}

type getEnrolledStudentsResponse []*studentResponse

func (c *courseAPI) GetEnrolledCoursesByTimelineClassification(ctx context.Context, classification CourseClassification) ([]*Course, error) {
	res := getEnrolledCoursesByTimelineClassificationResponse{}
	err := c.callMoodleFunction(ctx, &res, map[string]string{
		"wsfunction":     "core_course_get_enrolled_courses_by_timeline_classification",
		"classification": string(classification),
	})
	if err != nil {
		return nil, err
	}
	return mapToCourseList(res.Courses), nil
}

func (c *courseAPI) GetEnrolledStudentsByCourseID(ctx context.Context, id int) ([]*Student, error) {
	res := getEnrolledStudentsResponse{}
	err := c.callMoodleFunction(ctx, &res, map[string]string{
		"wsfunction": "core_enrol_get_enrolled_users",
		"courseid":   strconv.Itoa(id),
	})
	if err != nil {
		return nil, err
	}
	return mapToStudentList(res), nil
}

func mapToCourseList(courseResList []*courseResponse) []*Course {
	courses := make([]*Course, 0, len(courseResList))
	for _, courseRes := range courseResList {
		courses = append(courses, mapToCourse(courseRes))
	}
	return courses
}

func mapToStudentList(studentResList []*studentResponse) []*Student {
	courses := make([]*Student, 0, len(studentResList))
	for _, courseRes := range studentResList {
		courses = append(courses, mapToStudent(courseRes))
	}
	return courses
}

func mapToStudent(studentRes *studentResponse) *Student {
	var group string
	for _, g := range studentRes.Groups {
		group += g.Name + " "
	}
	var role string
	if len(studentRes.Roles) > 0 {
		role = studentRes.Roles[0].Role
	}
	return &Student{
		FN:        studentRes.IdNumber,
		FirstName: studentRes.FirstName,
		LastName:  studentRes.LastName,
		Role:      role,
		Group:     group,
	}
}

func mapToCourse(courseRes *courseResponse) *Course {
	return &Course{
		ID:              courseRes.ID,
		FullName:        courseRes.FullName,
		ShortName:       courseRes.ShortName,
		Summary:         courseRes.Summary,
		SummaryFormat:   courseRes.SummaryFormat,
		StartDate:       time.Unix(courseRes.StartDateUnix, 0),
		EndDate:         time.Unix(courseRes.StartDateUnix, 0),
		Visible:         courseRes.Visible,
		FullNameDisplay: courseRes.FullName,
		ViewURL:         courseRes.ViewURL,
		CourseImage:     courseRes.CourseImage,
		Progress:        courseRes.Progress,
		HasProgress:     courseRes.HasProgress,
		IsSavourite:     courseRes.IsSavourite,
		Hidden:          courseRes.Hidden,
		ShowShortName:   courseRes.ShowShortName,
		CourseCategory:  courseRes.CourseCategory,
	}
}
