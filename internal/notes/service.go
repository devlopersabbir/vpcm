package notes

import (
	"context"
)

type noteService struct {
	repo NoteRepository
}

func NewService(repo NoteRepository) NoteService {
	return &noteService{repo: repo}
}

func (s *noteService) AddNote(ctx context.Context, note *Note) error {
	return s.repo.Create(ctx, note)
}

func (s *noteService) GetNote(ctx context.Context, id uint) (*Note, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *noteService) GetServerNotes(ctx context.Context, serverID uint) ([]Note, error) {
	return s.repo.ListByServer(ctx, serverID)
}

func (s *noteService) UpdateNote(ctx context.Context, note *Note) error {
	return s.repo.Update(ctx, note)
}

func (s *noteService) RemoveNote(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
