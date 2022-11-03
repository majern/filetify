package client

import (
	"context"
	proto_v1 "github.com/msoft-dev/filetify/pkg/proto/v1"
	"github.com/msoft-dev/filetify/pkg/shared"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func getConn() (conn *grpc.ClientConn, ok bool) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(GetConfiguration().Rcp.Endpoint, opts...)
	st, _ := status.FromError(err)
	if err != nil {
		switch st.Code() {
		case codes.Unavailable:
			logrus.WithField("message", st.Message()).Warning("Failed to initialize gRPC connection")
		default:
			shared.HandleErrorWithMsg(err, true, "Failed to initialize gRPC connection")
		}
		return nil, false
	}

	return conn, true
}

func StartSynchronization() (entries []*shared.SyncEntry, ok bool) {
	conn, ok := getConn()
	if !ok {
		return nil, false
	}

	defer conn.Close()

	client := proto_v1.NewSynchronizationServiceClient(conn)

	res, err := client.StartSynchronization(context.Background(), getStartSyncReq())
	st, _ := status.FromError(err)
	if err != nil {
		switch st.Code() {
		case codes.Aborted:
		case codes.Unavailable:
			logrus.WithField("message", st.Message()).Warning("Failed to perform synchronization with server")
		default:
			shared.HandleErrorWithMsg(err, true, "Failed to perform synchronization with server")
		}
		return nil, false
	}

	logrus.WithField("response", res).Info("Synchronization started")

	return getEntryFromStartSyncRes(res), true
}

func getStartSyncReq() *proto_v1.StartSynchronizationRequest {
	filesFromCache := shared.GetAllFromCache()
	var filesToSynchronize []*proto_v1.FileSyncInfo

	for _, entry := range filesFromCache {
		fileSyncInfo := &proto_v1.FileSyncInfo{
			Path:      entry.Path,
			Status:    proto_v1.FileStatus(entry.Operation + 1),
			Timestamp: timestamppb.New(entry.Modified),
			IsDir:     entry.IsDir,
		}

		filesToSynchronize = append(filesToSynchronize, fileSyncInfo)
	}

	return &proto_v1.StartSynchronizationRequest{Files: filesToSynchronize}
}

func getEntryFromStartSyncRes(res *proto_v1.StartSynchronizationResponse) []*shared.SyncEntry {
	var syncEntries []*shared.SyncEntry

	for _, file := range res.Files {
		entry := &shared.SyncEntry{
			IsDir:     file.IsDir,
			Path:      file.Path,
			Action:    shared.SyncAction(file.Action - 1),
			Timestamp: file.Timestamp.AsTime(),
		}

		syncEntries = append(syncEntries, entry)
	}

	return syncEntries
}
