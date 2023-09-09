package app

import (
	"github.com/gofrs/uuid"
	"time"
)

type Session struct {
	UUID         *uuid.UUID `bson:"_id"`
	RefreshToken []byte     `bson:"refreshToken"`
	ExpiresAt    time.Time  `bson:"expiresAt"`
}
