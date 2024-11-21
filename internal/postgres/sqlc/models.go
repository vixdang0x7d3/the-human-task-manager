// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TaskPriority string

const (
	TaskPriorityH    TaskPriority = "H"
	TaskPriorityM    TaskPriority = "M"
	TaskPriorityL    TaskPriority = "L"
	TaskPriorityNone TaskPriority = "none"
)

func (e *TaskPriority) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TaskPriority(s)
	case string:
		*e = TaskPriority(s)
	default:
		return fmt.Errorf("unsupported scan type for TaskPriority: %T", src)
	}
	return nil
}

type NullTaskPriority struct {
	TaskPriority TaskPriority
	Valid        bool // Valid is true if TaskPriority is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTaskPriority) Scan(value interface{}) error {
	if value == nil {
		ns.TaskPriority, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TaskPriority.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTaskPriority) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TaskPriority), nil
}

type TaskStatus string

const (
	TaskStatusStarted   TaskStatus = "started"
	TaskStatusWaiting   TaskStatus = "waiting"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusDeleted   TaskStatus = "deleted"
)

func (e *TaskStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TaskStatus(s)
	case string:
		*e = TaskStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for TaskStatus: %T", src)
	}
	return nil
}

type NullTaskStatus struct {
	TaskStatus TaskStatus
	Valid      bool // Valid is true if TaskStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTaskStatus) Scan(value interface{}) error {
	if value == nil {
		ns.TaskStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TaskStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTaskStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TaskStatus), nil
}

type Project struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	ParentID uuid.NullUUID
	Title    string
}

type Task struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	ProjectID   uuid.NullUUID
	CompletedBy uuid.NullUUID
	Description string
	Priority    TaskPriority
	Status      TaskStatus
	Deadline    time.Time
	Schedule    time.Time
	Wait        time.Time
	Create      time.Time
	End         time.Time
}

type User struct {
	ID        uuid.UUID
	Username  string
	FirstName string
	LastName  string
	Email     string
	Password  string
	SignupAt  time.Time
	LastLogin time.Time
}