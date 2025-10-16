package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TicketType string

const (
	KitchenTicket  TicketType = "kitchen"
	CashierTicket  TicketType = "cashier"
	BarTicket      TicketType = "bar"
	CustomerTicket TicketType = "customer"
)

type PrintPriority int

const (
	LowPriority    PrintPriority = 0
	NormalPriority PrintPriority = 1
	HighPriority   PrintPriority = 2
	UrgentPriority PrintPriority = 3
)

type TicketStatus string

const (
	PendingStatus   TicketStatus = "pending"
	PrintingStatus  TicketStatus = "printing"
	CompletedStatus TicketStatus = "completed"
	FailedStatus    TicketStatus = "failed"
)

type TicketRequest struct {
	Content     string     `json:"content" binding:"required"`
	Type        TicketType `json:"type" binding:"required"`
	PrinterName string     `json:"printer_name" binding:"required"`
	Copies      int        `json:"copies" binding:"min=1,max=5"`

	OrderID  string `json:"order_id" binding:"required"`
	TableID  string `json:"table_id"`
	ServerID string `json:"server_id"`

	Priority PrintPriority `json:"priority"`

	Metadata map[string]interface{} `json:"metadata"`
}

type TicketJob struct {
	ID        string       `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Status    TicketStatus `json:"status"`

	FilePath     string `json:"file_path,omiempty"`
	PrinterName  string `json:"printer_name"`
	Copies       int    `json:"copies"`
	AttemptCount int    `json:"attempt_count"`
	LastError    string `json:"last_error,omitempty"`

	Request *TicketRequest `json:"request"`
}

func NewTicketJob(request *TicketRequest) *TicketJob {
	return &TicketJob{
		ID:           uuid.New().String(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Status:       PendingStatus,
		PrinterName:  request.PrinterName,
		Copies:       request.Copies,
		Request:      request,
		AttemptCount: 0,
	}
}

func (ticket *TicketJob) String() string {
	data, _ := json.MarshalIndent(ticket, "", "  ")

	return string(data)
}

func (ticket *TicketJob) ToArchiveData() map[string]interface{} {
	return map[string]interface{}{
		"job_id":     ticket.ID,
		"order_id":   ticket.Request.OrderID,
		"type":       ticket.Request.Type,
		"created_at": ticket.CreatedAt,
		"updated_at": ticket.UpdatedAt,
		"status":     ticket.Status,
		"printer":    ticket.PrinterName,
		"copies":     ticket.Copies,
		"attempts":   ticket.AttemptCount,
		"table_id":   ticket.Request.TableID,
		"server_id":  ticket.Request.ServerID,
		"metadata":   ticket.Request.Metadata,
		"last_error": ticket.LastError,
	}
}

func (ticket *TicketJob) ToArchiveMetadata() map[string]interface{} {
	return map[string]interface{}{
		"job_id":     ticket.ID,
		"order_id":   ticket.Request.OrderID,
		"type":       ticket.Request.Type,
		"created_at": ticket.CreatedAt,
		"updated_at": ticket.UpdatedAt,
		"status":     ticket.Status,
		"printer":    ticket.PrinterName,
		"copies":     ticket.Copies,
		"attempts":   ticket.AttemptCount,
		"table_id":   ticket.Request.TableID,
		"server_id":  ticket.Request.ServerID,
		"metadata":   ticket.Request.Metadata,
		"last_error": ticket.LastError,
		"content":    ticket.Request.Content,
	}
}

func (ticket *TicketJob) UpdateStatus(status TicketStatus, err error) {
	ticket.Status = status
	ticket.UpdatedAt = time.Now()

	if err != nil {
		ticket.LastError = err.Error()
		ticket.AttemptCount++
	}
}
