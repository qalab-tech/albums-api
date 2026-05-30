package repository

import (
	"context"
	"fmt"

	"albums-api/internal/db"
	"albums-api/internal/models"
)

// AlbumRepository — интерфейс репозитория
type AlbumRepository interface {
	GetAll(ctx context.Context) ([]models.Album, error)
	GetByID(ctx context.Context, id int) (models.Album, error)
	Create(ctx context.Context, album models.NewAlbum) (models.Album, error)
	Update(ctx context.Context, id int, update models.UpdateAlbum) (models.Album, error)
	Delete(ctx context.Context, id int) error
}

// albumRepo — реализация
type albumRepo struct {
	db *db.DB
}

// NewAlbumRepository — конструктор
func NewAlbumRepository(database *db.DB) AlbumRepository {
	return &albumRepo{db: database}
}

func (r *albumRepo) GetAll(ctx context.Context) ([]models.Album, error) {
	query := `SELECT id, title, artist, price FROM albums ORDER BY id`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения альбомов: %w", err)
	}
	defer rows.Close()

	var albums []models.Album
	for rows.Next() {
		var a models.Album
		if err := rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Price); err != nil {
			return nil, fmt.Errorf("ошибка сканирования: %w", err)
		}
		albums = append(albums, a)
	}

	return albums, nil
}

func (r *albumRepo) GetByID(ctx context.Context, id int) (models.Album, error) {
	query := `SELECT id, title, artist, price FROM albums WHERE id = $1`

	var album models.Album
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
	if err != nil {
		return models.Album{}, fmt.Errorf("альбом с id %d не найден: %w", id, err)
	}

	return album, nil
}

func (r *albumRepo) Create(ctx context.Context, newAlbum models.NewAlbum) (models.Album, error) {
	query := `INSERT INTO albums (title, artist, price) 
			  VALUES ($1, $2, $3) 
			  RETURNING id, title, artist, price`

	var album models.Album
	err := r.db.Pool.QueryRow(ctx, query, newAlbum.Title, newAlbum.Artist, newAlbum.Price).
		Scan(&album.ID, &album.Title, &album.Artist, &album.Price)

	if err != nil {
		return models.Album{}, fmt.Errorf("ошибка создания альбома: %w", err)
	}

	return album, nil
}

func (r *albumRepo) Update(ctx context.Context, id int, update models.UpdateAlbum) (models.Album, error) {
	query := `UPDATE albums SET 
				title = COALESCE($1, title),
				artist = COALESCE($2, artist),
				price = COALESCE($3, price)
			  WHERE id = $4 
			  RETURNING id, title, artist, price`

	var album models.Album
	err := r.db.Pool.QueryRow(ctx, query,
		update.Title,
		update.Artist,
		update.Price,
		id,
	).Scan(&album.ID, &album.Title, &album.Artist, &album.Price)

	if err != nil {
		return models.Album{}, fmt.Errorf("ошибка обновления альбома: %w", err)
	}

	return album, nil
}

func (r *albumRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM albums WHERE id = $1`
	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления альбома: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("альбом с id %d не найден", id)
	}

	return nil
}
