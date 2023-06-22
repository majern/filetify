package server

import (
	"context"
	protov1 "github.com/msoft-dev/filetify/pkg/proto/v1"
	"github.com/msoft-dev/filetify/pkg/shared"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"io"
	"sync"
)

type synchronizationService struct {
	lock           sync.Mutex
	syncInProgress bool
}

func (s *synchronizationService) StartSynchronization(ctx context.Context, request *protov1.StartSynchronizationRequest) (*protov1.StartSynchronizationResponse, error) {
	p, _ := peer.FromContext(ctx)
	logrus.WithFields(logrus.Fields{"request": request, "client": p.Addr}).Debug("Incomming synchronization request received")

	s.lock.Lock()
	if s.syncInProgress {
		s.lock.Unlock()
		logrus.WithField("client", p.Addr).Warning("Synchronization rejected for client as another synchronization in progress")
		return nil, status.Error(codes.Aborted, "Another synchronization in progress")
	}
	s.syncInProgress = true
	s.lock.Unlock()

	logrus.WithField("client", p.Addr).Info("Starting synchronization")
	syncFiles := SyncFiles(getFileEntriesFromReq(request)) //TODO: return files to synchronize
	var syncFileInfos []*protov1.ServerFileSyncInfo
	for _, file := range syncFiles {
		syncFileInfos = append(syncFileInfos, file.ToProto())
	}

	response := &protov1.StartSynchronizationResponse{Files: syncFileInfos}
	return response, nil
}

func (s *synchronizationService) FinishSynchronization(ctx context.Context, request *protov1.FinishSynchronizationRequest) (*protov1.FinishSynchronizationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *synchronizationService) UploadFile(server protov1.SynchronizationService_UploadFileServer) error {
	var receivedParts uint64
	var receivedLength uint64

	for {
		value, err := server.Recv()
		receivedParts += 1
		receivedLength += uint64(len(value.Data))

		if err == io.EOF {
			if receivedLength == value.TotalLength && receivedParts == value.TotalParts {
				logrus.WithField("filePart", value).Info("File received successfully")
				return server.SendAndClose(&protov1.UploadFileResponse{
					Status: protov1.FileTransferStatus_FILE_TRANSFER_STATUS_OK})
			}

			logrus.WithField("filePart", value).Warning("File not received correctly")
			return server.SendAndClose(&protov1.UploadFileResponse{
				Status: protov1.FileTransferStatus_FILE_TRANSFER_STATUS_FAILED})
		}
		if err != nil {
			logrus.WithError(err).WithField("filePart", value).Error("Error during receiving file part")
			return err
		}

		logrus.WithField("filePart", value).Debug("File part received")

	}
}

func (s *synchronizationService) DownloadFile(request *protov1.DownloadFileRequest, server protov1.SynchronizationService_DownloadFileServer) error {
	//TODO implement me
	panic("implement me")
}

func getFileEntriesFromReq(startSyncReq *protov1.StartSynchronizationRequest) []*shared.FileEntry {
	var result []*shared.FileEntry
	for _, file := range startSyncReq.Files {
		result = append(result, shared.ToFileEntry(file))
	}

	return result
}
