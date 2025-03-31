package usecase

import (
	"context"
	"github.com/google/uuid"
)

type RegisterMailReceiverUseCase struct {
	receiverRepository receiverRepository
	logger             logger
}

func NewRegisterMailReceiverUseCase(
	logger logger,
	receiverRepository receiverRepository,
) *RegisterMailReceiverUseCase {
	return &RegisterMailReceiverUseCase{
		logger:             logger,
		receiverRepository: receiverRepository,
	}
}

func (r *RegisterMailReceiverUseCase) Execute(ctx context.Context, userUUID uuid.UUID, mail string) error {
	err := r.receiverRepository.CreateMailReceiver(ctx, userUUID, mail)
	if err != nil {
		r.logger.Infof("не удалось зарегистрировать получателя: id: %s, mail: %s", userUUID.String(), mail)
		return err
	}

	r.logger.Infof("зарегистрирован получатель: id: %s, mail: %s", userUUID.String(), mail)
	return nil
}
