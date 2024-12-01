package db

import (
	"database/sql"
	"errors"
	"time"

	"circadian/internal/models"
)

func (d *Db) InsertHealthCheckResult(hcr *models.HealthCheckResult) error {
	query := `
	INSERT INTO health_check_results 
	(target_endpoint, http_method, callout_timestamp, response_code, response, time_taken) 
	VALUES ($1, $2, $3, $4, $5, $6) 
	RETURNING id`

	err := d.conn.QueryRow(
		query,
		hcr.TargetEndpoint,
		hcr.HTTPMethod,
		hcr.CalloutTime,
		hcr.ResponseCode,
		hcr.Response,
		hcr.TimeTaken,
	).Scan(&hcr.ID)

	return err
}

func (d *Db) UpdateHealthCheckResult(hcr *models.HealthCheckResult) error {
	query := `
	UPDATE health_check_results 
	SET target_endpoint = $1, 
		http_method = $2, 
		callout_timestamp = $3, 
		response_code = $4, 
		response = $5, 
		time_taken = $6 
	WHERE id = $7`

	result, err := d.conn.Exec(
		query,
		hcr.TargetEndpoint,
		hcr.HTTPMethod,
		hcr.CalloutTime,
		hcr.ResponseCode,
		hcr.Response,
		hcr.TimeTaken,
		hcr.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}

	return nil
}

func (d *Db) DeleteHealthCheckResult(id string) error {
	query := `DELETE FROM health_check_results WHERE id = $1`

	result, err := d.conn.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows deleted")
	}

	return nil
}

func (d *Db) GetHealthCheckResult(id string) (*models.HealthCheckResult, error) {
	query := `
	SELECT id, target_endpoint, http_method, callout_timestamp, response_code, response, time_taken 
	FROM health_check_results 
	WHERE id = $1`

	hcr := &models.HealthCheckResult{}
	err := d.conn.QueryRow(query, id).Scan(
		&hcr.ID,
		&hcr.TargetEndpoint,
		&hcr.HTTPMethod,
		&hcr.CalloutTime,
		&hcr.ResponseCode,
		&hcr.Response,
		&hcr.TimeTaken,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("health check result not found")
	}

	if err != nil {
		return nil, err
	}

	return hcr, nil
}

func (d *Db) ListHealthCheckResults() ([]models.HealthCheckResult, error) {
	query := `
	SELECT id, target_endpoint, http_method, callout_timestamp, response_code, response, time_taken 
	FROM health_check_results 
	WHERE callout_timestamp > $1
	ORDER BY callout_timestamp`

	rows, err := d.conn.Query(query, (time.Now().AddDate(0, 0, -60)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.HealthCheckResult
	for rows.Next() {
		var hcr models.HealthCheckResult
		err := rows.Scan(
			&hcr.ID,
			&hcr.TargetEndpoint,
			&hcr.HTTPMethod,
			&hcr.CalloutTime,
			&hcr.ResponseCode,
			&hcr.Response,
			&hcr.TimeTaken,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, hcr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
