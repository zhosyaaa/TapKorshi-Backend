package service

import (
	"github.com/zhosyaaa/RoommateTap/pkg/cache"
	"github.com/zhosyaaa/RoommateTap/pkg/email"
)

type EmailService struct {
	sender email.Sender
	config config.EmailConfig
	cache  cache.Cache

	sendpulseClients map[primitive.ObjectID]*sendpulse.Client
}
