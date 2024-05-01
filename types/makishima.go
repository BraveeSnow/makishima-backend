package types

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type MakishimaData struct {
	Logger   *zerolog.Logger
	Database *gorm.DB
}
