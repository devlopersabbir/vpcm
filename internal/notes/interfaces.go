package notes

import "context"

type NoteRepository interface {
	Create(ctx context.Context, note *Note) error
	GetByID(ctx context.Context, id uint) (*Note, error)
	ListByServer(ctx context.Context, serverID uint) ([]Note, error)
	Update(ctx context.Context, note *Note) error
	Delete(ctx context.Context, id uint) error
}

type NoteService interface {
	AddNote(ctx context.Context, note *Note) error
	GetNote(ctx context.Context, id uint) (*Note, error)
	GetServerNotes(ctx context.Context, serverID uint) ([]Note, error)
	UpdateNote(ctx context.Context, note *Note) error
	RemoveNote(ctx context.Context, id uint) error
}
