package models

type Segmentation struct {
	ID           int64  `json:"-" db:"id"`
	AddressSapID string `json:"address_sap_id" db:"address_sap_id"`
	AdrSegment   string `json:"adr_segment" db:"adr_segment"`
	SegmentID    int64  `json:"segment_id" db:"segment_id"`
}
