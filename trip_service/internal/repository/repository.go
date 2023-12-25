package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"strconv"
	"trip_service/internal/model"
)

type Repository struct {
	db *sqlx.DB
}

func (repo *Repository) GetStarted(ctx context.Context, params *model.ParamsStarted) ([]model.TripStarted, error) {
	query := sq.Select("*").
		From("started_trips").
		Offset(params.Offset).
		Limit(params.Limit).
		PlaceholderFormat(sq.Dollar)
	idInt, err := strconv.Atoi(params.Id)
	if err != nil {
		return nil, err
	}
	if params.Id != "" {
		query = query.Where(sq.Eq{"id": idInt})
	}
	if params.OfferId != "" {
		query = query.Where(sq.Eq{"offer_id": params.OfferId})
	}
	if params.DriverId != "" {
		query = query.Where(sq.Eq{"driver_id": params.DriverId})
	}
	if params.CurrentStage != "" {
		query = query.Where(sq.Eq{"current_stage": params.CurrentStage})
	}

	sql, args, err := query.ToSql()
	rows, err := repo.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	trips := make([]model.TripStarted, 0)
	for rows.Next() {
		trip := model.TripStarted{}

		if err = rows.StructScan(&trip); err != nil {
			return nil, err
		}

		trips = append(trips, trip)
	}

	return trips, nil
}
func (repo *Repository) CreateStarted(ctx context.Context, trip *model.TripStarted) (string, error) {
	sql, args, err := sq.
		Insert("started_trips").Columns("offer_id", "driver_id", "current_stage").
		Values(trip.OfferId, trip.DriverId, trip.CurrentStage).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	var id string
	row := repo.db.QueryRowContext(ctx, sql, args...)
	if err = row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}
func (repo *Repository) DeleteStarted(ctx context.Context, params *model.ParamsStarted) error {
	query := sq.Delete("started_trips")
	if params.Id != "" {
		query = query.Where(sq.Eq{"id": params.Id})
	}
	if params.OfferId != "" {
		query = query.Where(sq.Eq{"offer_id": params.OfferId})
	}
	if params.DriverId != "" {
		query = query.Where(sq.Eq{"driver_id": params.DriverId})
	}
	if params.CurrentStage != "" {
		query = query.Where(sq.Eq{"current_stage": params.CurrentStage})
	}
	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	rows, err := repo.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}
func (repo *Repository) UpdateStarted(ctx context.Context, params *model.ParamsStarted, value *model.TripStarted) error {
	query := sq.Update("started_trips")

	if value.OfferId != "" {
		query = query.Set("offer_id", value.OfferId)
	}
	if value.DriverId != "" {
		query = query.Set("driver_id", value.DriverId)
	}
	if value.CurrentStage != "" {
		query = query.Set("current_stage", value.CurrentStage)
	}

	idInt, err := strconv.Atoi(params.Id)
	if err != nil {
		return err
	}
	if params.Id != "" {
		query = query.Where(sq.Eq{"id": idInt})
	}
	if params.OfferId != "" {
		query = query.Where(sq.Eq{"offer_id": params.OfferId})
	}
	if params.DriverId != "" {
		query = query.Where(sq.Eq{"driver_id": params.DriverId})
	}
	if params.CurrentStage != "" {
		query = query.Where(sq.Eq{"current_stage": params.CurrentStage})
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	rows, err := repo.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (repo *Repository) GetFinished(ctx context.Context, params *model.ParamsFinished) ([]model.TripFinished, error) {
	query := sq.Select("*").
		From("finished_trip").
		Offset(params.Offset).
		Limit(params.Limit).
		PlaceholderFormat(sq.Dollar)
	idInt, err := strconv.Atoi(params.Id)
	if err != nil {
		return nil, err
	}
	if params.Id != "" {
		query = query.Where(sq.Eq{"id": idInt})
	}
	if params.OfferId != "" {
		query = query.Where(sq.Eq{"offer_id": params.OfferId})
	}
	if params.DriverId != "" {
		query = query.Where(sq.Eq{"driver_id": params.DriverId})
	}
	if params.CurrentStage != "" {
		query = query.Where(sq.Eq{"current_stage": params.CurrentStage})
	}

	sql, args, err := query.ToSql()
	rows, err := repo.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	trips := make([]model.TripFinished, 0)
	for rows.Next() {
		trip := model.TripFinished{}

		if err = rows.StructScan(&trip); err != nil {
			return nil, err
		}

		trips = append(trips, trip)
	}

	return trips, nil
}
func (repo *Repository) CreateFinished(ctx context.Context, trip *model.TripFinished) (string, error) {
	sql, args, err := sq.
		Insert("finished_trip").Columns("offer_id", "driver_id", "current_stage").
		Values(trip.OfferId, trip.DriverId, trip.CurrentStage).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	var id string
	row := repo.db.QueryRowContext(ctx, sql, args...)
	if err = row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}
func (repo *Repository) UpdateFinished(ctx context.Context, params *model.ParamsFinished, value *model.TripFinished) error {
	query := sq.Update("finished_trip")

	if value.OfferId != "" {
		query = query.Set("offer_id", value.OfferId)
	}
	if value.DriverId != "" {
		query = query.Set("offer_id", value.DriverId)
	}
	if value.CurrentStage != "" {
		query = query.Set("offer_id", value.CurrentStage)
	}

	idInt, err := strconv.Atoi(params.Id)
	if err != nil {
		return err
	}
	if params.Id != "" {
		query = query.Where(sq.Eq{"id": idInt})
	}
	if params.OfferId != "" {
		query = query.Where(sq.Eq{"offer_id": params.OfferId})
	}
	if params.DriverId != "" {
		query = query.Where(sq.Eq{"driver_id": params.DriverId})
	}
	if params.CurrentStage != "" {
		query = query.Where(sq.Eq{"current_stage": params.CurrentStage})
	}

	sql, args, err := query.ToSql()
	rows, err := repo.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}
