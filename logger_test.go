package bjaeger

import (
	"testing"

	mock_log "github.com/go-masonry/mortar/interfaces/log/mock"
	"github.com/golang/mock/gomock"
)

func TestLogger(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockLogger := mock_log.NewMockLogger(ctrl)
	logger := &logWrapper{inner: mockLogger}
	mockLogger.EXPECT().Debug(gomock.Not(nil), "msg-debug")
	logger.Debugf("msg-debug")
	mockLogger.EXPECT().Info(gomock.Not(nil), "msg-info")
	logger.Infof("msg-info")
	mockLogger.EXPECT().Error(gomock.Not(nil), "error")
	logger.Error("error")
}
