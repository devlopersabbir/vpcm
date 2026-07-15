package notes

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrNoteNotFound = errors.New("note not found")

type sqlNoteRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) NoteRepository {
	return &sqlNoteRepository{db: db}
}

func (r *sqlNoteRepository) Create(ctx context.Context, note *Note) error {
	return r.db.WithContext(ctx).Create(note).Error
}

func (r *sqlNoteRepository) GetByID(ctx context.Context, id uint) (*Note, error) {
	var note Note
	err := r.db.WithContext(ctx).First(&note, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoteNotFound
		}
		return nil, err
	}
	return &note, nil
}

func (r *sqlNoteRepository) ListByServer(ctx context.Context, serverID uint) ([]Note, error) {
	var notes []Note
	err := r.db.WithContext(ctx).Where("server_id = ?", serverID).Find(&notes).Error
	return notes, err
}

func (r *sqlNoteRepository) Update(ctx context.Context, note *Note) error {
	return r.db.WithContext(ctx).Save(note).Error
}

func (r *sqlNoteRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&Note{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoteNotFound
	}
	return nil
}
