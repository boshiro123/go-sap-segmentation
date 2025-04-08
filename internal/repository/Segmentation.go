package repository

import (
	"github.com/jmoiron/sqlx"
	"go-test/internal/models"
)

type SegmentationRepository struct {
	db *sqlx.DB
}

func NewSegmentationRepository(db *sqlx.DB) *SegmentationRepository {
	return &SegmentationRepository{
		db: db,
	}
}

func (r *SegmentationRepository) InsertOrUpdate(segments []*models.Segmentation) error {
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

func (r *SegmentationRepository) GetByAddressSapID(addressSapID string) (*models.Segmentation, error) {
	var segment models.Segmentation
	err := r.db.Get(&segment, "SELECT * FROM segmentation WHERE address_sap_id = $1", addressSapID)
	return &segment, err
}

func (r *SegmentationRepository) GetAll() ([]*models.Segmentation, error) {
	var segments []*models.Segmentation
	err := r.db.Select(&segments, "SELECT * FROM segmentation")
	return segments, err
}
