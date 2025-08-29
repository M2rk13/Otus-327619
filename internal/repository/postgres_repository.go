package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/M2rk13/Otus-327619/internal/config"
	"github.com/M2rk13/Otus-327619/internal/model/api"
	logmodel "github.com/M2rk13/Otus-327619/internal/model/log"

	"github.com/google/uuid"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(ctx context.Context, cfg config.PostgresConfig) (*PostgresStore, error) {
	db, err := sql.Open("pgx", cfg.PostgresUri)

	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()

		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	store := &PostgresStore{db: db}

	log.Println("PostgreSQL repository initialized successfully.")

	return store, nil
}

func (s *PostgresStore) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *PostgresStore) executeInTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

func (s *PostgresStore) CreateRequest(req *api.Request) {
	req.Id = uuid.New().String()
	query := `INSERT INTO requests (id, "from", "to", amount) VALUES ($1, $2, $3, $4)`

	err := s.executeInTransaction(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.Exec(query, req.Id, req.From, req.To, req.Amount)

		return err
	})

	if err != nil {
		log.Printf("ERROR: failed to create request: %v", err)
	}
}

func (s *PostgresStore) GetRequestByID(id string) *api.Request {
	query := `SELECT id, "from", "to", amount FROM requests WHERE id = $1`
	var req api.Request
	err := s.db.QueryRow(query, id).Scan(&req.Id, &req.From, &req.To, &req.Amount)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("ERROR: failed to get request by id: %v", err)
		}

		return nil
	}

	return &req
}

func (s *PostgresStore) GetAllRequests() []*api.Request {
	query := `SELECT id, "from", "to", amount FROM requests`
	rows, err := s.db.Query(query)

	if err != nil {
		log.Printf("ERROR: failed to get all requests: %v", err)

		return nil
	}

	defer rows.Close()

	var requests []*api.Request

	for rows.Next() {
		var req api.Request

		if err := rows.Scan(&req.Id, &req.From, &req.To, &req.Amount); err != nil {
			log.Printf("ERROR: failed to scan request: %v", err)

			continue
		}

		requests = append(requests, &req)
	}

	return requests
}

func (s *PostgresStore) UpdateRequest(req *api.Request) bool {
	query := `UPDATE requests SET "from" = $2, "to" = $3, amount = $4 WHERE id = $1`

	err := s.executeInTransaction(context.Background(), func(tx *sql.Tx) error {
		res, err := tx.Exec(query, req.Id, req.From, req.To, req.Amount)

		if err != nil {
			return err
		}

		rowsAffected, _ := res.RowsAffected()

		if rowsAffected == 0 {
			return sql.ErrNoRows
		}

		return nil
	})

	return err == nil
}

func (s *PostgresStore) DeleteRequest(id string) bool {
	query := `DELETE FROM requests WHERE id = $1`

	err := s.executeInTransaction(context.Background(), func(tx *sql.Tx) error {
		res, err := tx.Exec(query, id)

		if err != nil {
			return err
		}

		rowsAffected, _ := res.RowsAffected()

		if rowsAffected == 0 {
			return sql.ErrNoRows
		}

		return nil
	})

	return err == nil
}

func (s *PostgresStore) CreateResponse(resp *api.Response) {
	resp.Id = uuid.New().String()
	query := `INSERT INTO responses
    	(id, success, terms, privacy, query_id, query_from, query_to, query_amount, info_timestamp, info_quote, result)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	err := s.executeInTransaction(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.Exec(query,
			resp.Id,
			resp.Success,
			resp.Terms,
			resp.Privacy,
			resp.Query.Id,
			resp.Query.From,
			resp.Query.To,
			resp.Query.Amount,
			resp.Info.Timestamp,
			resp.Info.Quote,
			resp.Result)

		return err
	})

	if err != nil {
		log.Printf("ERROR: failed to create response: %v", err)
	}
}

func (s *PostgresStore) GetResponseByID(id string) *api.Response {
	query := `
		SELECT
		    id,
		    success,
		    terms,
		    privacy,
		    query_id,
		    query_from,
		    query_to,
		    query_amount,
		    info_timestamp,
		    info_quote,
		    result
		FROM responses
		WHERE id = $1`

	var resp api.Response
	err := s.db.QueryRow(query, id).Scan(&resp.Id,
		&resp.Success,
		&resp.Terms,
		&resp.Privacy,
		&resp.Query.Id,
		&resp.Query.From,
		&resp.Query.To,
		&resp.Query.Amount,
		&resp.Info.Timestamp,
		&resp.Info.Quote,
		&resp.Result)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("ERROR: failed to get response by id: %v", err)
		}

		return nil
	}

	return &resp
}

func (s *PostgresStore) GetAllResponses() []*api.Response {
	query := `
		SELECT
			id,
			success,
			terms,
			privacy,
			query_id,
			query_from,
			query_to,
			query_amount,
			info_timestamp,
			info_quote,
			result
		FROM responses`

	rows, err := s.db.Query(query)

	if err != nil {
		log.Printf("ERROR: failed to get all responses: %v", err)

		return nil
	}

	defer rows.Close()

	var responses []*api.Response

	for rows.Next() {
		var resp api.Response

		err := rows.Scan(&resp.Id,
			&resp.Success,
			&resp.Terms,
			&resp.Privacy,
			&resp.Query.Id,
			&resp.Query.From,
			&resp.Query.To,
			&resp.Query.Amount,
			&resp.Info.Timestamp,
			&resp.Info.Quote,
			&resp.Result)

		if err != nil {
			log.Printf("ERROR: failed to scan response: %v", err)

			continue
		}

		responses = append(responses, &resp)
	}

	return responses
}

func (s *PostgresStore) UpdateResponse(resp *api.Response) bool {
	query := `
		UPDATE responses
		SET success = $2,
		    terms = $3,
		    privacy = $4,
		    query_id = $5,
		    query_from = $6,
		    query_to = $7,
		    query_amount = $8,
		    info_timestamp = $9,
		    info_quote = $10,
		    result = $11
		WHERE id = $1`

	err := s.executeInTransaction(context.Background(), func(tx *sql.Tx) error {
		res, err := tx.Exec(query,
			resp.Id,
			resp.Success,
			resp.Terms,
			resp.Privacy,
			resp.Query.Id,
			resp.Query.From,
			resp.Query.To,
			resp.Query.Amount,
			resp.Info.Timestamp,
			resp.Info.Quote,
			resp.Result)

		if err != nil {
			return err
		}

		rowsAffected, _ := res.RowsAffected()

		if rowsAffected == 0 {
			return sql.ErrNoRows
		}

		return nil
	})

	return err == nil
}

func (s *PostgresStore) DeleteResponse(id string) bool {
	query := `DELETE FROM responses WHERE id = $1`

	err := s.executeInTransaction(context.Background(), func(tx *sql.Tx) error {
		res, err := tx.Exec(query, id)

		if err != nil {
			return err
		}

		rowsAffected, _ := res.RowsAffected()

		if rowsAffected == 0 {
			return sql.ErrNoRows
		}

		return nil
	})

	return err == nil
}

func (s *PostgresStore) CreateConversionLog(logItem *logmodel.ConversionLog) {
	logItem.Id = uuid.New().String()
	requestJSON, err := json.Marshal(logItem.Request)

	if err != nil {
		log.Printf("ERROR: failed to marshal request for log: %v", err)

		return
	}

	responseJSON, err := json.Marshal(logItem.Response)

	if err != nil {
		log.Printf("ERROR: failed to marshal response for log: %v", err)

		return
	}

	query := `INSERT INTO conversion_logs (id, timestamp, request, response) VALUES ($1, $2, $3, $4)`

	err = s.executeInTransaction(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.Exec(query, logItem.Id, logItem.Timestamp, requestJSON, responseJSON)

		return err
	})

	if err != nil {
		log.Printf("ERROR: failed to create conversion log: %v", err)
	}
}

func (s *PostgresStore) GetConversionLogByID(id string) *logmodel.ConversionLog {
	query := `SELECT id, timestamp, request, response FROM conversion_logs WHERE id = $1`

	var logItem logmodel.ConversionLog
	var requestJSON, responseJSON []byte

	err := s.db.QueryRow(query, id).Scan(&logItem.Id, &logItem.Timestamp, &requestJSON, &responseJSON)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("ERROR: failed to get conversion log by id: %v", err)
		}

		return nil
	}

	if err := json.Unmarshal(requestJSON, &logItem.Request); err != nil {
		log.Printf("ERROR: failed to unmarshal request from log: %v", err)

		return nil
	}

	if err := json.Unmarshal(responseJSON, &logItem.Response); err != nil {
		log.Printf("ERROR: failed to unmarshal response from log: %v", err)

		return nil
	}

	return &logItem
}

func (s *PostgresStore) GetAllConversionLogs() []*logmodel.ConversionLog {
	query := `SELECT id, timestamp, request, response FROM conversion_logs`
	rows, err := s.db.Query(query)

	if err != nil {
		log.Printf("ERROR: failed to get all conversion logs: %v", err)

		return nil
	}

	defer rows.Close()

	var logs []*logmodel.ConversionLog

	for rows.Next() {
		var logItem logmodel.ConversionLog
		var requestJSON, responseJSON []byte

		if err := rows.Scan(&logItem.Id, &logItem.Timestamp, &requestJSON, &responseJSON); err != nil {
			log.Printf("ERROR: failed to scan conversion log: %v", err)

			continue
		}

		if err := json.Unmarshal(requestJSON, &logItem.Request); err != nil {
			log.Printf("ERROR: failed to unmarshal request from log: %v", err)

			continue
		}

		if err := json.Unmarshal(responseJSON, &logItem.Response); err != nil {
			log.Printf("ERROR: failed to unmarshal response from log: %v", err)

			continue
		}

		logs = append(logs, &logItem)
	}

	return logs
}

func (s *PostgresStore) UpdateConversionLog(logItem *logmodel.ConversionLog) bool {
	requestJSON, err := json.Marshal(logItem.Request)

	if err != nil {
		log.Printf("ERROR: failed to marshal request for log update: %v", err)

		return false
	}

	responseJSON, err := json.Marshal(logItem.Response)

	if err != nil {
		log.Printf("ERROR: failed to marshal response for log update: %v", err)

		return false
	}

	query := `UPDATE conversion_logs SET timestamp = $2, request = $3, response = $4 WHERE id = $1`

	err = s.executeInTransaction(context.Background(), func(tx *sql.Tx) error {
		res, err := tx.Exec(query, logItem.Id, logItem.Timestamp, requestJSON, responseJSON)
		if err != nil {
			return err
		}
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			return sql.ErrNoRows
		}
		return nil
	})

	return err == nil
}

func (s *PostgresStore) DeleteConversionLog(id string) bool {
	query := `DELETE FROM conversion_logs WHERE id = $1`

	err := s.executeInTransaction(context.Background(), func(tx *sql.Tx) error {
		res, err := tx.Exec(query, id)
		if err != nil {
			return err
		}
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			return sql.ErrNoRows
		}
		return nil
	})

	return err == nil
}

func (s *PostgresStore) GetNewConversionRequests() []*api.Request        { return nil }
func (s *PostgresStore) GetNewConversionResponses() []*api.Response      { return nil }
func (s *PostgresStore) GetNewConversionLogs() []*logmodel.ConversionLog { return nil }
