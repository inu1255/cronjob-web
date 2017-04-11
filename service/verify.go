package service

import (
	"github.com/inu1255/gev2/models"
)

type Verify struct {
	models.VerifyModel `xorm:"extends"`
}
