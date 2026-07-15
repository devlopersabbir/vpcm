package inventory

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type sqlServerRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) ServerRepository {
	return &sqlServerRepository{db: db}
}

func (r *sqlServerRepository) Create(ctx context.Context, server *Server) error {
	return r.db.WithContext(ctx).Create(server).Error
}

func (r *sqlServerRepository) GetByID(ctx context.Context, id uint) (*Server, error) {
	var server Server
	err := r.db.WithContext(ctx).Preload("Tags").Preload("Software").First(&server, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrServerNotFound
		}
		return nil, err
	}
	return &server, nil
}

func (r *sqlServerRepository) GetByUUID(ctx context.Context, uuid string) (*Server, error) {
	var server Server
	err := r.db.WithContext(ctx).Preload("Tags").Preload("Software").Where("uuid = ?", uuid).First(&server).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrServerNotFound
		}
		return nil, err
	}
	return &server, nil
}

func (r *sqlServerRepository) List(ctx context.Context) ([]Server, error) {
	var servers []Server
	err := r.db.WithContext(ctx).Preload("Tags").Find(&servers).Error
	return servers, err
}

func (r *sqlServerRepository) Update(ctx context.Context, server *Server) error {
	return r.db.WithContext(ctx).Save(server).Error
}

func (r *sqlServerRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&Server{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrServerNotFound
	}
	return nil
}

func (r *sqlServerRepository) AddTag(ctx context.Context, serverID uint, tagName string) error {
	var server Server
	if err := r.db.WithContext(ctx).First(&server, serverID).Error; err != nil {
		return err
	}
	var tag Tag
	err := r.db.WithContext(ctx).Where("name = ?", tagName).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tag = Tag{Name: tagName}
			if err := r.db.WithContext(ctx).Create(&tag).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return r.db.WithContext(ctx).Model(&server).Association("Tags").Append(&tag)
}

func (r *sqlServerRepository) RemoveTag(ctx context.Context, serverID uint, tagName string) error {
	var server Server
	if err := r.db.WithContext(ctx).First(&server, serverID).Error; err != nil {
		return err
	}
	var tag Tag
	if err := r.db.WithContext(ctx).Where("name = ?", tagName).First(&tag).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&server).Association("Tags").Delete(&tag)
}
