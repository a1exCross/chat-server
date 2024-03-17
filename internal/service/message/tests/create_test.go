package tests

import (
	"context"
	"testing"
	"time"

	"github.com/a1exCross/chat-server/internal/model"
	logrepository "github.com/a1exCross/chat-server/internal/repository/log"
	messagerepository "github.com/a1exCross/chat-server/internal/repository/message"
	messageservice "github.com/a1exCross/chat-server/internal/service/message"

	"github.com/a1exCross/common/pkg/client/db"
	dbmocks "github.com/a1exCross/common/pkg/client/db/mocks"
	"github.com/a1exCross/common/pkg/client/db/transaction"

	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	type dbClientMock func(mc *minimock.Controller) db.Client
	type txManagerMock func(mc *minimock.Controller) db.TxManager

	ctx := context.Background()
	mc := minimock.NewController(t)
	id := int64(1)

	timeNow := time.Now()

	msgDTO := model.MessageDTO{
		Author:    "author",
		Content:   "content",
		CreatedAt: timeNow,
	}

	tests := []struct {
		name      string
		err       error
		dbClient  dbClientMock
		txManager txManagerMock
	}{
		{
			name: "successfull test",
			dbClient: func(mc *minimock.Controller) db.Client {
				client := dbmocks.NewClientMock(mc)
				dbb := dbmocks.NewDBMock(mc)
				row := dbmocks.NewRowMock(mc)

				row.ScanMock.Set(func(dest ...interface{}) (err error) {
					res, ok := dest[0].(*int64)
					if ok {
						*res = id
					}

					_ = res

					return nil
				})

				dbb.QueryRowContextMock.Return(row)
				dbb.ExecContextMock.Return(pgconn.CommandTag{}, nil)

				client.DBMock.Return(dbb)

				return client
			},
			txManager: func(mc *minimock.Controller) db.TxManager {
				tx := dbmocks.NewTxMock(mc)
				transactor := dbmocks.NewTransactorMock(mc)

				txOptions := pgx.TxOptions{
					IsoLevel: pgx.ReadCommitted,
				}

				tx.CommitMock.Return(nil)
				transactor.BeginTxMock.Expect(ctx, txOptions).Return(tx, nil)

				txManager := transaction.NewTransactionManager(transactor)

				return txManager
			},
			err: nil,
		},
		{
			name: "error at create",
			dbClient: func(mc *minimock.Controller) db.Client {
				client := dbmocks.NewClientMock(mc)
				dbb := dbmocks.NewDBMock(mc)
				row := dbmocks.NewRowMock(mc)

				row.ScanMock.Set(func(dest ...interface{}) (err error) {
					res, ok := dest[0].(*int64)
					if ok {
						*res = id
					}

					_ = res

					return errors.New("failed to scan")
				})

				dbb.QueryRowContextMock.Return(row)

				client.DBMock.Return(dbb)

				return client
			},
			txManager: func(mc *minimock.Controller) db.TxManager {
				tx := dbmocks.NewTxMock(mc)
				transactor := dbmocks.NewTransactorMock(mc)

				txOptions := pgx.TxOptions{
					IsoLevel: pgx.ReadCommitted,
				}

				tx.RollbackMock.Return(nil)

				transactor.BeginTxMock.Expect(ctx, txOptions).Return(tx, nil)

				txManager := transaction.NewTransactionManager(transactor)

				return txManager
			},
			err: errors.New("failed executing code inside transaction: error at query to database: failed to scan"),
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			dbMockClient := test.dbClient(mc)

			txManager := test.txManager(mc)

			msgRepo := messagerepository.NewRepository(dbMockClient)
			logRepo := logrepository.NewRepository(dbMockClient)

			msgServ := messageservice.NewService(msgRepo, txManager, logRepo)

			err := msgServ.SendMessage(ctx, msgDTO)

			if err != nil && test.err != nil {
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.Equal(t, test.err, err)
			}
		})
	}
}
