package otp

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"
)

// MockOTPStore represents a stored OTP
type MockOTPStore struct {
	Code      string
	ExpiresAt time.Time
}

// MockProvider implements the Provider interface for development
type MockProvider struct {
	mu     sync.RWMutex
	otps   map[string]*MockOTPStore
	useFixedOTP bool
}

// NewMockProvider creates a new mock OTP provider
func NewMockProvider(useFixedOTP bool) *MockProvider {
	provider := &MockProvider{
		otps:   make(map[string]*MockOTPStore),
		useFixedOTP: useFixedOTP,
	}
	// Start cleanup goroutine
	go provider.cleanupExpired()
	return provider
}

// SendOTP sends a mock OTP (stores it in memory)
func (m *MockProvider) SendOTP(ctx context.Context, phone string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var code string
	if m.useFixedOTP {
		// Fixed OTP for testing
		code = "123456"
	} else {
		// Generate random 6-digit OTP
		code = generateOTP()
	}

	m.otps[phone] = &MockOTPStore{
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	log.Printf("[MOCK OTP] Sent OTP %s to phone %s", code, phone)
	return code, nil
}

// VerifyOTP verifies the mock OTP
func (m *MockProvider) VerifyOTP(ctx context.Context, phone string, code string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	stored, exists := m.otps[phone]
	if !exists {
		return false, nil
	}

	if time.Now().After(stored.ExpiresAt) {
		delete(m.otps, phone)
		return false, nil
	}

	if stored.Code != code {
		return false, nil
	}

	// OTP verified, remove it
	delete(m.otps, phone)
	log.Printf("[MOCK OTP] Verified OTP for phone %s", phone)
	return true, nil
}

// cleanupExpired periodically removes expired OTPs
func (m *MockProvider) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for phone, stored := range m.otps {
			if now.After(stored.ExpiresAt) {
				delete(m.otps, phone)
			}
		}
		m.mu.Unlock()
	}
}

// generateOTP generates a random 6-digit OTP
func generateOTP() string {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "123456" // Fallback
	}
	return fmt.Sprintf("%06d", n.Int64())
}
