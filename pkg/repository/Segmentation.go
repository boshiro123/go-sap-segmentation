package repository

import (
	"github.com/jmoiron/sqlx"
)

type Segmentation struct {
	ID           int64  `json:"-" db:"id"`
	AddressSapID string `json:"address_sap_id" db:"address_sap_id"`
	AdrSegment   string `json:"adr_segment" db:"adr_segment"`
	SegmentID    int64  `json:"segment_id" db:"segment_id"`
}

type SegmentationRepository struct {
	db *sqlx.DB
}

func NewSegmentationRepository(db *sqlx.DB) *SegmentationRepository {
	return &SegmentationRepository{
		db: db,
	}
}

func (r *SegmentationRepository) InsertOrUpdate(segments []*Segmentation) error {
	if len(segments) == 0 {
		return nil
	}

	query := `
		INSERT INTO segmentation (address_sap_id, adr_segment, segment_id)
		VALUES (:address_sap_id, :adr_segment, :segment_id)
		ON CONFLICT (address_sap_id) DO UPDATE 
		SET adr_segment = EXCLUDED.adr_segment, 
			segment_id = EXCLUDED.segment_id
	`

	_, err := r.db.NamedExec(query, segments)
	return err
}

func (r *SegmentationRepository) GetByAddressSapID(addressSapID string) (*Segmentation, error) {
	var segment Segmentation
	err := r.db.Get(&segment, "SELECT * FROM segmentation WHERE address_sap_id = $1", addressSapID)
	return &segment, err
}

func (r *SegmentationRepository) GetAll() ([]*Segmentation, error) {
	var segments []*Segmentation
	err := r.db.Select(&segments, "SELECT * FROM segmentation")
	return segments, err
}
