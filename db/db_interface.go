package db

import (
	"context"

	mocks "github.com/Gorillarock/granitex/db/mocks"
	"github.com/Gorillarock/granitex/model"
)

type DBInteractor interface {
	InsertTx(ctx context.Context, document model.DocumentEntry) (string, error)
	GetRx(ctx context.Context, params model.GetRxParams) model.ResponsePayloadRxHandler
}

func NewDBInteractorMock() *mocks.DBInteractor {
	DBHandler = &mocks.DBInteractor{}
	return DBHandler.(*mocks.DBInteractor)
}
