package service

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/zhosyaaa/RoommateTap/internal/domain"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type SessionService struct {
	redisClient *redis.Client
}

func NewSessionService(redisClient *redis.Client) *SessionService {
	return &SessionService{redisClient: redisClient}
}
func (s *SessionService) CreateSession(userID uint, username string, expiresAt time.Time) (string, *domain.Session, error) {
	sessionID := generateSessionID()
	session := &domain.Session{
		Userid:    userID,
		Username:  username,
		ExpiresAt: expiresAt,
	}
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return "", nil, fmt.Errorf("failed to marshal session: %w", err)
	}
	ctx := context.Background()
	err = s.redisClient.Set(ctx, sessionID, sessionJSON, expiresAt.Sub(time.Now())).Err()
	if err != nil {
		return "", nil, fmt.Errorf("failed to set session in Redis: %w", err)
	}
	return sessionID, session, nil
}

func (s SessionService) GetSession(sessionID string) (*domain.Session, error) {
	ctx := context.Background()
	sessionJSON, err := s.redisClient.Get(ctx, sessionID).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session from Redis: %w", err)
	}
	var session domain.Session
	err = json.Unmarshal([]byte(sessionJSON), &session)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}
	return &session, nil
}

func generateSessionID() string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Println("Error generating random bytes:", err)
		return ""
	}
	randomHex := hex.EncodeToString(randomBytes)
	return strconv.FormatInt(timestamp, 10) + randomHex
}
