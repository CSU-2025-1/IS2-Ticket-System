package repository

import (
	"context"
	"errors"
	"fmt"
	"gateway/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
	"net"
)

const (
	serviceRepositoryCreateOp    = "ServiceRepository.Create"
	serviceRepositoryGetWithType = "ServiceRepository.GetWithType"
	serviceRepositoryGetAll      = "ServiceRepository.GetAll"
)

var (
	ErrIncorrectIP = errors.New("incorrect IP address")
)

type ServiceRepository struct {
	db *pgxpool.Pool
}

func NewServiceRepository(db *pgxpool.Pool) ServiceRepository {
	return ServiceRepository{
		db: db,
	}
}

func (s *ServiceRepository) Create(ctx context.Context, service entity.Service) (err error) {
	defer func() {
		err = fmt.Errorf("%s: %w", serviceRepositoryCreateOp, err)
	}()

	sql := `INSERT INTO gateway.services (ip, port, type) VALUES ($1, $2, $3)`
	_, err = s.db.Exec(ctx, sql, service.IP, service.Port, service.Type)
	return err
}

func (s *ServiceRepository) GetWithType(ctx context.Context, typeToSearch string) (services []entity.Service, err error) {
	defer func() {
		err = fmt.Errorf("%s: %w", serviceRepositoryGetWithType, err)
	}()

	sql := `SELECT ip, port, type FROM gateway.service WHERE type = $1`
	rows, err := s.db.Query(ctx, sql, typeToSearch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		ip          string
		port        uint16
		serviceType string
	)
	for rows.Next() {
		err = rows.Scan(&ip, &port, &serviceType)
		if err != nil {
			return nil, err
		}

		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			return nil, ErrIncorrectIP
		}

		services = append(services, entity.Service{
			Type: serviceType,
			IP:   parsedIP,
			Port: port,
		})
	}

	return services, err
}

func (s *ServiceRepository) GetAll(ctx context.Context) (services []entity.Service, err error) {
	defer func() {
		err = fmt.Errorf("%s: %w", serviceRepositoryGetAll, err)
	}()

	sql := `SELECT ip, port, type FROM gateway.service`
	rows, err := s.db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		ip          string
		port        uint16
		serviceType string
	)
	for rows.Next() {
		err = rows.Scan(&ip, &port, &serviceType)
		if err != nil {
			return nil, err
		}

		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			return nil, ErrIncorrectIP
		}

		services = append(services, entity.Service{
			Type: serviceType,
			IP:   parsedIP,
			Port: port,
		})
	}

	return services, nil
}
