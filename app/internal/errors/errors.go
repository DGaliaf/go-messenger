package custom_error

import "encoding/json"

var (
	// --- Adapters Layer ---
	ErrEntityNotFound    = NewCustomError(nil, "entity not found", "")
	ErrNoRowsAffected    = NewCustomError(nil, "no rows affected", "")
	ErrSQLExecution      = NewCustomError(nil, "sql execution error", "")
	ErrCreateTransaction = NewCustomError(nil, "failed to create transaction", "")
	ErrTransactionClosed = NewCustomError(nil, "action was done in closed transaction", "")
	ErrCommitTransaction = NewCustomError(nil, "failed to commit changes in transaction", "")
	ErrAcquireConnection = NewCustomError(nil, "failed to acquire connection to database", "")
	// ----------------------

	// --- Service Layer ---
	ErrUserNotExist      = NewCustomError(nil, "user does not exists", "")
	ErrUserNotInChat     = NewCustomError(nil, "user does not participated in chat", "")
	ErrChatDuplicate     = NewCustomError(nil, "chat already exists", "")
	ErrUserDuplicate     = NewCustomError(nil, "user already exists", "")
	ErrNotEnoughUsers    = NewCustomError(nil, "not enough users", "")
	ErrUserAlreadyInChat = NewCustomError(nil, "chat contain users with the same ID", "")
	ErrChatNotExist      = NewCustomError(nil, "chat does not exists", "")
	// ----------------------
)

type CustomError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
}

func NewCustomError(err error, message, developerMessage string) *CustomError {
	return &CustomError{
		Err:              err,
		Message:          message,
		DeveloperMessage: developerMessage,
	}
}

func (c CustomError) Unwrap() error {
	return c.Err
}

func (c CustomError) Error() string {
	return c.Message
}

func (c CustomError) Marshal() []byte {
	marshal, err := json.Marshal(c)
	if err != nil {
		return nil
	}

	return marshal
}

func systemError(err error) *CustomError {
	return NewCustomError(err, "internal system error", err.Error())
}
