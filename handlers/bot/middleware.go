package bot

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/pkg/apperrors"
)

func (h *handler) AuthorizeUser(ctx context.Context, telegramID int64, requiresAdmin bool) (string, error) {
	exists, err := h.user.ExistsUserByTelegramID(ctx, telegramID)
	if err != nil {
		h.log.Error("failed to check user by telegram_id", "error", err)
		return "Произошла ошибка!", err
	}

	if !exists {
		return "Нужно выполнить команду /start", apperrors.ErrNotFoundUser
	}

	if requiresAdmin {
		isAdmin, err := h.user.IsAdmin(ctx, telegramID)
		if err != nil {
			h.log.Error("failed to check admin status", "error", err)
			return "Произошла ошибка!", err
		}

		if !isAdmin {
			return "У вас нет доступа к этой команде", apperrors.ErrAccessDenied
		}
	}

	return "", nil
}
