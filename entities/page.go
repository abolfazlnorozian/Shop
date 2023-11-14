package entities

import (
	"time"
)

type Pages struct {
	Id        int       `json:"_id" bson:"_id"`
	Meta      Metas     `json:"meta" bson:"meta"`
	Mode      string    `json:"mode" form:"mode" bson:"mode"`
	Rows      []int     `json:"rows" bson:"rows"`
	Url       string    `json:"url" bson:"url"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	V         int       `json:"__v" bson:"__v"`
}

type Metas struct {
	Keywords    []string `json:"keywords" bson:"keywords"`
	Title       string   `json:"title,omitempty" bson:"title,omitempty"`
	Description string   `json:"description,omitempty" bson:"description,omitempty"`
}

// type PagePresentation struct {
// 	ID        int               `json:"_id" bson:"_id"`
// 	Meta      Metas             `json:"meta" bson:"meta"`
// 	Mode      string            `json:"mode" bson:"mode"`
// 	URL       string            `json:"url" bson:"url"`
// 	CreatedAt time.Time         `json:"createdAt" bson:"createdAt"`
// 	UpdatedAt time.Time         `json:"updatedAt" bson:"updatedAt"`
// 	V         int               `json:"__v" bson:"__v"`
// 	Rows      []RowPresentation `json:"rows" bson:"rows"`
// }

// type RowPresentation struct {
// 	Fluid           bool                 `json:"fluid" bson:"fluid"`
// 	BackgroundColor string               `json:"backgroundColor" bson:"backgroundColor"`
// 	CreatedAt       time.Time            `json:"createdAt" bson:"createdAt"`
// 	UpdatedAt       time.Time            `json:"updatedAt" bson:"updatedAt"`
// 	PageID          int                  `json:"pageId" bson:"pageId"`
// 	V               int                  `json:"__v" bson:"__v"`
// 	ID              int                  `json:"_id" bson:"_id"`
// 	Cols            []ColumnPresentation `json:"cols" bson:"cols"`
// }

// type ColumnPresentation struct {
// 	ID              int                 `json:"_id" bson:"_id"`
// 	Size            Size                `json:"size" bson:"size"`
// 	Elevation       int                 `json:"elevation" bson:"elevation"`
// 	Padding         string              `json:"padding" bson:"padding"`
// 	Radius          string              `json:"radius" bson:"radius"`
// 	Margin          string              `json:"margin" bson:"margin"`
// 	BackgroundColor string              `json:"backgroundColor" bson:"backgroundColor"`
// 	DataUrl         string              `json:"dataUrl" bson:"dataUrl"`
// 	IsMore          bool                `json:"isMore" bson:"isMore"`
// 	LayoutType      string              `json:"layoutType" bson:"layoutType"`
// 	DataType        string              `json:"dataType" bson:"dataType"`
// 	MoreUrl         string              `json:"moreUrl" bson:"moreUrl"`
// 	Content         ContentPresentation `json:"content" bson:"content"`
// 	Name            string              `json:"name" bson:"name"`
// 	RowID           int                 `json:"rowId" bson:"rowId"`
// 	CreatedAt       time.Time           `json:"createdAt" bson:"createdAt"`
// 	UpdatedAt       time.Time           `json:"updatedAt" bson:"updatedAt"`
// 	V               int                 `json:"__v" bson:"__v"`
// }

// type ContentPresentation struct {
// 	Alt   string    `json:"alt" bson:"alt"`
// 	Link  string    `json:"link" bson:"link"`
// 	Image ImageCont `json:"image" bson:"image"`
// }
