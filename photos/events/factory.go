package events

import (
	"github.com/graeme-hill/gnet/photos/events/pb"
	"github.com/graeme-hill/gnet/sys/eventstore"
)

const (
	EventPhotoUploaded    = "PhotoUploaded"
	EventPhotoRearrange   = "PhotoRearrange"
	EventPhotoDelete      = "PhotoDelete"
	EventAlbumRemovePhoto = "AlbumRemovePhoto"
	EventAlbumCreate      = "AlbumCreate"
	EventAlbumRename      = "AlbumRename"
	EventAlbumAddPhoto    = "AlbumAddPhoto"
)

func NewPhotoUploadedEvent(path string) (eventstore.DomainEvent, error) {
	return eventstore.NewDomainEvent(EventPhotoUploaded, &pb.PhotoUploaded{
		Path: path,
	})
}

func NewPhotoRearrangeEvent(albumID string, path string, after string) (eventstore.DomainEvent, error) {
	return eventstore.NewDomainEvent(EventPhotoRearrange, &pb.PhotoRearrange{
		AlbumId: albumID,
		Path:    path,
		After:   after,
	})
}

func NewPhotoDeleteEvent(path string) (eventstore.DomainEvent, error) {
	return eventstore.NewDomainEvent(EventPhotoDelete, &pb.PhotoDelete{
		Path: path,
	})
}

func NewAlbumRemovePhotoEvent(albumID string, path string) (eventstore.DomainEvent, error) {
	return eventstore.NewDomainEvent(EventAlbumRemovePhoto, &pb.AlbumRemovePhoto{
		AlbumId: albumID,
		Path:    path,
	})
}

func NewAlbumCreateEvent(id string, name string) (eventstore.DomainEvent, error) {
	return eventstore.NewDomainEvent(EventAlbumCreate, &pb.AlbumCreate{
		Id:   id,
		Name: name,
	})
}

func NewAlbumRenameEvent(id string, newName string) (eventstore.DomainEvent, error) {
	return eventstore.NewDomainEvent(EventAlbumRename, &pb.AlbumRename{
		Id:      id,
		NewName: newName,
	})
}

func NewAlbumAddPhotoEvent(albumID string, path string, after string) (eventstore.DomainEvent, error) {
	return eventstore.NewDomainEvent(EventAlbumAddPhoto, &pb.AlbumAddPhoto{
		AlbumId: albumID,
		Path:    path,
		After:   after,
	})
}
