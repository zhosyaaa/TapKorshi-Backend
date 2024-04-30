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
func (s *SessionService) CreateSession(session *domain.Session) (string, error) {
	sessionID := generateSessionID()

	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return "", fmt.Errorf("failed to marshal session: %w", err)
	}
	ctx := context.Background()
	err = s.redisClient.Set(ctx, sessionID, sessionJSON, session.ExpiresAt.Sub(time.Now())).Err()
	if err != nil {
		return "", fmt.Errorf("failed to set session in Redis: %w", err)
	}
	return sessionID, nil
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
func (s *SessionService) DeleteSession(sessionID string) error {
	ctx := context.Background()
	err := s.redisClient.Del(ctx, sessionID).Err()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("session not found")
		}
		return fmt.Errorf("failed to delete session from Redis: %w", err)
	}
	return nil
}
