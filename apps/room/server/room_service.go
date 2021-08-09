package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/pion/ion-log"
	room "github.com/pion/ion/apps/room/proto"
	"github.com/pion/ion/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RoomService struct {
	room.UnimplementedRoomServiceServer
	roomLock sync.RWMutex
	rooms    map[string]*Room
	closed   chan struct{}
}

func NewRoomService() *RoomService {
	s := &RoomService{
		rooms:  make(map[string]*Room),
		closed: make(chan struct{}),
	}
	go s.stat()
	return s
}

func (s *RoomService) Close() {
	close(s.closed)
}

func (s *RoomService) CreateRoom(ctx context.Context, in *room.CreateRoomRequest) (*room.CreateRoomReply, error) {
	sid := in.Sid
	r := s.getRoom(sid)
	if r == nil {
		return nil, fmt.Errorf("room not found: %s", sid)
	}
	//TODO
	return nil, status.Errorf(codes.Unimplemented, "method CreateRoom not implemented")
}

func (s *RoomService) DeleteRoom(ctx context.Context, in *room.DeleteRoomRequest) (*room.DeleteRoomReply, error) {
	sid := in.Sid
	r := s.getRoom(sid)
	if r == nil {
		return nil, fmt.Errorf("room not found: %s", sid)
	}
	//TODO
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRoom not implemented")
}

func (s *RoomService) AddParticipant(ctx context.Context, in *room.AddParticipantRequest) (*room.AddParticipantReply, error) {
	sid := in.Sid
	r := s.getRoom(sid)
	if r == nil {
		return nil, fmt.Errorf("room not found: %s", sid)
	}
	//TODO
	return nil, status.Errorf(codes.Unimplemented, "method AddParticipant not implemented")
}

func (s *RoomService) RemoveParticipant(ctx context.Context, in *room.RemoveParticipantRequest) (*room.RemoveParticipantReply, error) {
	sid := in.Sid
	r := s.getRoom(sid)
	if r == nil {
		return nil, fmt.Errorf("room not found: %s", sid)
	}
	//TODO
	return nil, status.Errorf(codes.Unimplemented, "method RemoveParticipant not implemented")
}

func (s *RoomService) GetParticipants(ctx context.Context, in *room.GetParticipantsRequest) (*room.GetParticipantsReply, error) {
	sid := in.Sid
	r := s.getRoom(sid)
	if r == nil {
		return nil, fmt.Errorf("room not found: %s", sid)
	}
	//TODO
	return nil, nil
}

func (s *RoomService) LockConference(ctx context.Context, in *room.LockConferenceRequest) (*room.LockConferenceReply, error) {
	sid := in.Sid
	r := s.getRoom(sid)
	if r == nil {
		return nil, fmt.Errorf("room not found: %s", sid)
	}
	//TODO
	return nil, status.Errorf(codes.Unimplemented, "method LockConference not implemented")
}

func (s *RoomService) EndConference(ctx context.Context, in *room.EndConferenceRequest) (*room.EndConferenceReply, error) {
	sid := in.Sid
	r := s.getRoom(sid)
	if r == nil {
		return nil, fmt.Errorf("room not found: %s", sid)
	}
	//TODO
	return nil, status.Errorf(codes.Unimplemented, "method EndConference not implemented")
}

func (s *RoomService) SetImportance(ctx context.Context, in *room.SetImportanceRequest) (*room.SetImportanceReply, error) {
	sid := in.Sid
	r := s.getRoom(sid)
	if r == nil {
		return nil, fmt.Errorf("room not found: %s", sid)
	}
	//TODO
	return nil, status.Errorf(codes.Unimplemented, "method SetImportance not implemented")
}

func (s *RoomService) EditParticipantInfo(ctx context.Context, in *room.EditParticipantInfoRequest) (*room.EditParticipantInfoReply, error) {
	sid := in.Sid
	r := s.getRoom(sid)
	if r == nil {
		return nil, fmt.Errorf("room not found: %s", sid)
	}
	//TODO
	return nil, status.Errorf(codes.Unimplemented, "method EditParticipantInfo not implemented")
}

func (s *RoomService) createRoom(sid string, sfuNID string) *Room {
	s.roomLock.Lock()
	defer s.roomLock.Unlock()
	if r := s.rooms[sid]; r != nil {
		return r
	}
	r := newRoom(sid, sfuNID)
	s.rooms[sid] = r
	return r
}

func (s *RoomService) getRoom(id string) *Room {
	s.roomLock.RLock()
	defer s.roomLock.RUnlock()
	return s.rooms[id]
}

func (s *RoomService) delRoom(r *Room) {
	id := r.SID()
	s.roomLock.Lock()
	defer s.roomLock.Unlock()
	if s.rooms[id] == r {
		delete(s.rooms, id)
	}
}

// stat peers
func (s *RoomService) stat() {
	t := time.NewTicker(util.DefaultStatCycle)
	defer t.Stop()
	for {
		select {
		case <-t.C:
		case <-s.closed:
			log.Infof("stop stat")
			return
		}

		var info string
		s.roomLock.RLock()
		for sid, room := range s.rooms {
			info += fmt.Sprintf("room: %s\npeers: %d\n", sid, room.count())
		}
		s.roomLock.RUnlock()
		if len(info) > 0 {
			log.Infof("\n----------------signal-----------------\n" + info)
		}
	}
}
