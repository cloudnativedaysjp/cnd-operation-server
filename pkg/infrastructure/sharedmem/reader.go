package sharedmem

import (
	"fmt"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
)

type ReaderIface interface {
	DisableAutomation(trackId int32) (bool, error)
	Talks(trackId int32) (model.Talks, error)
}

type Reader struct {
	UseStorageForDisableAutomation bool
	UseStorageForTalks             bool
}

func (w Reader) DisableAutomation(trackId int32) (bool, error) {
	if !w.UseStorageForDisableAutomation {
		return false, fmt.Errorf("UseStorageForDisableAutomation was false")
	}
	storageForDisableAutomationMutex.RLock()
	defer storageForDisableAutomationMutex.RUnlock()
	disabled, ok := storageForDisableAutomation[trackId]
	if !ok {
		return false, fmt.Errorf("trackId %d is not found", trackId)
	}
	return disabled, nil
}

func (w Reader) Talks(trackId int32) (model.Talks, error) {
	if !w.UseStorageForTalks {
		return nil, fmt.Errorf("UseStorageForTalks was false")
	}
	storageForTalksMutex.RLock()
	defer storageForTalksMutex.RUnlock()
	talks, ok := storageForTalks[trackId]
	if !ok {
		return nil, fmt.Errorf("trackId %d is not found", trackId)
	}
	return talks, nil
}
