package git

/*
#include <git2.h>
#include <git2/errors.h>
*/
import "C"

import (
	"runtime"
)

type Status int

const (
	StatusCurrent         Status = C.GIT_STATUS_CURRENT
	StatusIndexNew               = C.GIT_STATUS_INDEX_NEW
	StatusIndexModified          = C.GIT_STATUS_INDEX_MODIFIED
	StatusIndexDeleted           = C.GIT_STATUS_INDEX_DELETED
	StatusIndexRenamed           = C.GIT_STATUS_INDEX_RENAMED
	StatusIndexTypeChange        = C.GIT_STATUS_INDEX_TYPECHANGE
	StatusWtNew                  = C.GIT_STATUS_WT_NEW
	StatusWtModified             = C.GIT_STATUS_WT_NEW
	StatusWtDeleted              = C.GIT_STATUS_WT_DELETED
	StatusWtTypeChange           = C.GIT_STATUS_WT_TYPECHANGE
	StatusWtRenamed              = C.GIT_STATUS_WT_RENAMED
	StatusIgnored                = C.GIT_STATUS_IGNORED
)

type StatusEntry struct {
	Status         Status
	HeadToIndex    DiffDelta
	IndexToWorkdir DiffDelta
}

func statusEntryFromC(statusEntry *C.git_status_entry) StatusEntry {
	return StatusEntry {
		Status:         Status(statusEntry.status),
		HeadToIndex:    diffDeltaFromC(statusEntry.head_to_index),
		IndexToWorkdir: diffDeltaFromC(statusEntry.index_to_workdir),
	}
}

type StatusList struct {
	ptr *C.git_status_list
}

func newStatusListFromC(ptr *C.git_status_list) *StatusList {
	if ptr == nil {
		return nil
	}

	statusList := &StatusList{
		ptr: ptr,
	}

	runtime.SetFinalizer(statusList, (*StatusList).Free)
	return statusList
}

func (statusList *StatusList) Free() error {
	if statusList.ptr == nil {
		return ErrInvalid
	}
	runtime.SetFinalizer(statusList, nil)
	C.git_status_list_free(statusList.ptr)
	statusList.ptr = nil
	return nil
}

func (statusList *StatusList) ByIndex(index int) (StatusEntry, error) {
	if statusList.ptr == nil {
		return StatusEntry{}, ErrInvalid
	}
	ptr := C.git_status_byindex(statusList.ptr, C.size_t(index))
	return statusEntryFromC(ptr), nil
}

func (statusList *StatusList) EntryCount() (int, error) {
	if statusList.ptr == nil {
		return -1, ErrInvalid
	}
	return int(C.git_status_list_entrycount(statusList.ptr)), nil
}
