package service

import (
	log "github.com/sirupsen/logrus"
	"secret-social-network/app/storage"
)

func CommitLog(linkLog storage.ConsensusLinkLog, ok bool) {
	newState := storage.ConsensusLinkLogStateCommit
	if !ok {
		log.Warnf("link log %d mark as [failed]", linkLog.ID)
		newState = storage.ConsensusLinkLogStateFailed
	}
	err := linkLog.ChState(newState)
	if err != nil {
		log.Errorf("commit link log %d state %s failed: %s", linkLog.ID, newState, err.Error())
	}
}
