package server

import (
	"context"
	proto_v1 "github.com/msoft-dev/filetify/pkg/proto/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

type synchronizationService struct {
	lock           sync.Mutex
	syncInProgress bool
}

func (s *synchronizationService) StartSynchronization(_ context.Context, request *proto_v1.StartSynchronizationRequest) (*proto_v1.StartSynchronizationResponse, error) {
	s.lock.Lock()
	if s.syncInProgress {
		return nil, status.Error(codes.Aborted, "Another synchronization in progress")
	}
	s.syncInProgress = true
	s.lock.Unlock()

	logrus.WithField("request", request).Info("Starting synchronization")

	response := &proto_v1.StartSynchronizationResponse{Files: []*proto_v1.ServerFileSyncInfo{}}
	return response, nil
}

func (s *synchronizationService) FinishSynchronization(ctx context.Context, request *proto_v1.FinishSynchronizationRequest) (*proto_v1.FinishSynchronizationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *synchronizationService) UploadFile(server proto_v1.SynchronizationService_UploadFileServer) error {
	//TODO implement me
	panic("implement me")
}

func (s *synchronizationService) DownloadFile(request *proto_v1.DownloadFileRequest, server proto_v1.SynchronizationService_DownloadFileServer) error {
	//TODO implement me
	panic("implement me")
}
