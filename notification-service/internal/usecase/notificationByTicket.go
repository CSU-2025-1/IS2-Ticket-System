package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

type Ticket struct {
	ID              uuid.UUID
	Type            string
	Title           string
	Priority        int
	ResponsibleUUID uuid.UUID
}

type NotificationByTicketUseCase struct {
	userRepository     userRepository
	receiverRepository receiverRepository
	mailService        mailService
	logger             logger
}

func NewNotificationByTicketUseCase(
	logger logger,
	receiverRepository receiverRepository,
	userRepository userRepository,
	mailService mailService,
) *NotificationByTicketUseCase {
	return &NotificationByTicketUseCase{
		logger:             logger,
		userRepository:     userRepository,
		receiverRepository: receiverRepository,
		mailService:        mailService,
	}
}

func (c *NotificationByTicketUseCase) Execute(ctx context.Context, ticket Ticket) error {
	c.logger.Infof("получен тикет: %#v для рассылки уведомления разосланы", ticket)

	var receiversIds []uuid.UUID
	switch ticket.Type {
	case "user":
		receiversIds = []uuid.UUID{ticket.ResponsibleUUID}
	case "group":
		ids, err := c.userRepository.GetAllUserIDsByGroupID(ctx, ticket.ResponsibleUUID)
		if err != nil {
			return err
		}
		receiversIds = ids
	}

	receivers, err := c.receiverRepository.GetAllMailReceiversByUUIDs(ctx, receiversIds)
	if err != nil {
		return err
	}

	for i := range receivers {
		err := c.mailService.Notify(
			ctx,
			fmt.Sprintf(
				"Вам поступил Тикет: %s с ID: %s и приоритетом: %d",
				ticket.Title,
				ticket.ID.String(),
				ticket.Priority,
			),
			fmt.Sprintf(
				"Тикет ID: %s",
				ticket.ID.String(),
			),
			receivers[i],
		)
		if err != nil {
			c.logger.Warnf("ошибка отправки уведомления: %s", err.Error())
		}

		c.logger.Infof("отправлено уведомление по тикету: %s на почту: %s", ticket.ID, receivers[i].Mail)
	}

	c.logger.Infof("уведомления по тикету разосланы %s", ticket.ID)
	return nil
}
