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
)

func GetConn() (conn *grpc.ClientConn, ok bool) {
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

func StartSynchronization(conn *grpc.ClientConn) (entries []*shared.SyncEntry) {

	client := proto_v1.NewSynchronizationServiceClient(conn)

	res, err := client.StartSynchronization(context.Background(), getStartSyncReq())

	if err != nil {
		handleSyncError(err)
		return nil
	}

	logrus.WithField("response", res).Info("Synchronization started")

	return getEntryFromStartSyncRes(res)
}

func UploadFile(conn *grpc.ClientConn, entry *shared.SyncEntry) {
	client := proto_v1.NewSynchronizationServiceClient(conn)

	stream, err := client.UploadFile(context.Background())

	if err != nil {
		handleSyncError(err)
	}

	//TODO: Get file from hard drive
	parts, _ := shared.LoadFile(entry.Path)
	logrus.WithField("parts-count", len(parts))
	//TODO: Get parts from file and foreach do stream send.

	err = stream.Send(&proto_v1.UploadFileRequest{})

	if err != nil {
		shared.HandleErrorWithMsg(err, true, "Failed to upload file")
	}
}

func handleSyncError(err error) {
	st, _ := status.FromError(err)
	switch st.Code() {
	case codes.Aborted, codes.Unavailable:
		logrus.WithField("message", st.Message()).Warning("Failed to perform synchronization with server")
	default:
		shared.HandleErrorWithMsg(err, true, "Failed to perform synchronization with server")
	}
}

func getStartSyncReq() *proto_v1.StartSynchronizationRequest {
	filesFromCache := shared.GetAllFromCache()
	var filesToSynchronize []*proto_v1.FileSyncInfo

	for _, entry := range filesFromCache {
		fileSyncInfo := entry.ToProto()
		filesToSynchronize = append(filesToSynchronize, fileSyncInfo)
	}

	return &proto_v1.StartSynchronizationRequest{Files: filesToSynchronize}
}

func getEntryFromStartSyncRes(res *proto_v1.StartSynchronizationResponse) []*shared.SyncEntry {
	var syncEntries []*shared.SyncEntry

	for _, file := range res.Files {
		entry := &shared.SyncEntry{
			Key:       file.Key,
			IsDir:     file.IsDir,
			Path:      file.Path,
			Action:    shared.SyncAction(file.Action - 1),
			Timestamp: file.Timestamp.AsTime(),
		}

		syncEntries = append(syncEntries, entry)
	}

	return syncEntries
}
