package info

import (
	info3 "github.com/codestates/WBABEProject-05/config/info"
	"github.com/codestates/WBABEProject-05/protocol"
	"github.com/codestates/WBABEProject-05/util"
	"github.com/gin-gonic/gin"
)

var instance *infoControl

type infoControl struct {
}

func GetInstance() *infoControl {
	if instance != nil {
		return instance
	}
	instance = &infoControl{}
	return instance
}

func (h *infoControl) GetInformation(g *gin.Context) {
	path := util.Flags[util.InformationFlag.Name]
	info := info3.NewInfo(*path)
	protocol.SuccessData(info).Response(g)
}
