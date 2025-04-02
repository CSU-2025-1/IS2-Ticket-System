package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ticket-service/internal/core"
)

type DB interface {
	CreateTicket(ctx context.Context, ticket core.Ticket) error
	UpdateTicketStatus(ctx context.Context, ticketUUID uuid.UUID, status string) error
	AssignResponsible(ctx context.Context, ticketUUID, responsibleUUID uuid.UUID) error
	GetTickets(ctx context.Context, status string) ([]core.Ticket, error)
}

type Kafka interface {
	SendMessageTicketCreation(ticket core.Ticket) error
}

type Handler struct {
	db    DB
	kafka Kafka
}

func New(db DB, kafka Kafka) Handler {
	return Handler{
		db:    db,
		kafka: kafka,
	}
}

type CreateTicketRequest struct {
	Name          string    `json:"name" binding:"required"`
	Description   string    `json:"description"`
	Status        string    `json:"status" binding:"required"`
	CreatedBy     uuid.UUID `json:"created_by" binding:"required"`
	RecipientType string    `json:"recipient_type" binding:"required"`
	RecipientUUID uuid.UUID `json:"recipient_uuid" binding:"required"`
	Priority      int       `json:"priority"`
}

func (h *Handler) CreateTicket(ctx *gin.Context) {
	var req CreateTicketRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket := core.Ticket{
		UUID:          uuid.New(),
		Name:          req.Name,
		Description:   req.Description,
		Status:        req.Status,
		CreatedBy:     req.CreatedBy,
		CreatedOn:     time.Now(),
		UpdatedOn:     time.Now(),
		RecipientType: req.RecipientType,
		RecipientUUID: req.RecipientUUID,
		Priority:      req.Priority,
	}

	if err := h.db.CreateTicket(ctx.Request.Context(), ticket); err != nil {
		slog.Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create ticket"})
		return
	}

	if err := h.kafka.SendMessageTicketCreation(ticket); err != nil {
	}

	ctx.JSON(http.StatusCreated, ticket)
}

type UpdateTicketStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (h *Handler) UpdateTicketStatus(ctx *gin.Context) {
	ticketUUID, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		slog.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket UUID"})
		return
	}

	var req UpdateTicketStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.UpdateTicketStatus(ctx.Request.Context(), ticketUUID, req.Status); err != nil {
		slog.Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update ticket status"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ticket status updated successfully"})
}

type AssignResponsibleRequest struct {
	ResponsibleUUID uuid.UUID `json:"responsible_uuid" binding:"required"`
}

func (h *Handler) AssignTicketResponsible(ctx *gin.Context) {
	ticketUUID, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		slog.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket UUID"})
		return
	}

	var req AssignResponsibleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.AssignResponsible(ctx.Request.Context(), ticketUUID, req.ResponsibleUUID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign responsible"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "responsible assigned successfully"})
}

func (h *Handler) GetTickets(ctx *gin.Context) {
	status := ctx.Query("status")

	tickets, err := h.db.GetTickets(ctx.Request.Context(), status)
	if err != nil {
		slog.Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tickets"})
		return
	}

	if tickets == nil {
		tickets = []core.Ticket{}
	}

	ctx.JSON(http.StatusOK, tickets)
}

func (h *Handler) Check(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
