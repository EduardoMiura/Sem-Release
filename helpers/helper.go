package helpers

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

func ErrorWriter(ctx context.Context, w http.ResponseWriter, err error, logger *zap.Logger) {
	// var status int
	// var jerr *jobad.Error

	// switch terr := err.(type) {
	// case *httpError:
	// 	status = terr.status
	// 	jerr = jobad.NewError(fmt.Sprintf("HTTP%d", terr.status), terr.message)
	// case *jobad.Error:
	// 	status = getErrorHTTPCode(terr)
	// 	jerr = terr
	// default:
	// 	status = http.StatusInternalServerError
	// 	jerr = jobad.NewUnknownError(err.Error())
	// }

	// logger.Error(
	// 	jerr.Message,
	// 	zap.String("error-code", string(jerr.ErrCode)),
	// 	zap.String("transaction-id", jobad.GetTransactionID(ctx)),
	// 	zap.String("internal-id", jobad.GetInternalID(ctx)),
	// )

	// responseWriter(ctx, w, status, jerr)
}
