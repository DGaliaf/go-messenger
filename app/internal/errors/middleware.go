package custom_error

import (
	"errors"
	"net/http"
)

type handler func(w http.ResponseWriter, r *http.Request) error

func Middleware(h handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		var customError *CustomError
		err := h(w, r)
		if err != nil {
			if errors.As(err, &customError) {
				if errors.Is(err, ErrEntityNotFound) {
					w.WriteHeader(http.StatusNoContent)
					w.Write(ErrEntityNotFound.Marshal())
					return
				}

				if errors.Is(err, ErrNoRowsAffected) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(ErrNoRowsAffected.Marshal())
					return
				}

				if errors.Is(err, ErrSQLExecution) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(ErrSQLExecution.Marshal())
					return
				}

				if errors.Is(err, ErrCreateTransaction) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(ErrCreateTransaction.Marshal())
					return
				}

				if errors.Is(err, ErrTransactionClosed) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(ErrTransactionClosed.Marshal())
					return
				}

				if errors.Is(err, ErrCommitTransaction) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(ErrCommitTransaction.Marshal())
					return
				}

				if errors.Is(err, ErrAcquireConnection) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(ErrAcquireConnection.Marshal())
					return
				}

				if errors.Is(err, ErrUserNotExist) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(ErrUserNotExist.Marshal())
					return
				}

				if errors.Is(err, ErrChatDuplicate) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(ErrChatDuplicate.Marshal())
					return
				}

				if errors.Is(err, ErrUserDuplicate) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(ErrUserDuplicate.Marshal())
					return
				}

				if errors.Is(err, ErrNotEnoughUsers) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(ErrNotEnoughUsers.Marshal())
					return
				}

				if errors.Is(err, ErrUserAlreadyInChat) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(ErrUserAlreadyInChat.Marshal())
					return
				}

				if errors.Is(err, ErrChatNotExist) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(ErrChatNotExist.Marshal())
					return
				}

				if errors.Is(err, ErrUserNotInChat) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(ErrUserNotInChat.Marshal())
					return
				}

				err = err.(*CustomError)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(customError.Marshal())
				return
			}

			w.WriteHeader(http.StatusTeapot)
			w.Write(systemError(err).Marshal())
		}
	}
}
