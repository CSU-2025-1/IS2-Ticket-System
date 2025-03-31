package usecase

import (
	"context"
	"github.com/google/uuid"
	"notification-service/internal/model"
)

type (
	receiverRepository interface {
		GetAllMailReceiversByUUIDs(ctx context.Context, ids []uuid.UUID) ([]model.Receiver, error)
		CreateMailReceiver(ctx context.Context, userUUID uuid.UUID, mail string) error
	}
	userRepository interface {
		GetAllUserIDsByGroupID(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error)
	}

	mailService interface {
		Notify(ctx context.Context, message, title string, receiver model.Receiver) error
	}

	logger interface {
		Infof(message string, args ...interface{})
		Warnf(message string, args ...interface{})
	}
)
