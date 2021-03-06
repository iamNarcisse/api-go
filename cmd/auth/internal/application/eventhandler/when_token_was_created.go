package eventhandler

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/vardius/go-api-boilerplate/cmd/auth/internal/domain/token"
	"github.com/vardius/go-api-boilerplate/cmd/auth/internal/infrastructure/persistence"
	"github.com/vardius/go-api-boilerplate/pkg/domain"
	apperrors "github.com/vardius/go-api-boilerplate/pkg/errors"
	"github.com/vardius/go-api-boilerplate/pkg/eventbus"
)

// WhenTokenWasCreated handles event
func WhenTokenWasCreated(db *sql.DB, repository persistence.TokenRepository) eventbus.EventHandler {
	fn := func(parentCtx context.Context, event domain.Event) error {
		ctx, cancel := context.WithTimeout(parentCtx, time.Second*120)
		defer cancel()

		e := token.WasCreated{}
		if err := json.Unmarshal(event.Payload, &e); err != nil {
			return apperrors.Wrap(err)
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return apperrors.Wrap(err)
		}
		defer tx.Rollback()

		if err := repository.Add(ctx, e); err != nil {
			return apperrors.Wrap(err)
		}

		if err := tx.Commit(); err != nil {
			return apperrors.Wrap(err)
		}

		return nil
	}

	return fn
}
