package domain

import "time"

type ScmCommit struct {
	Hash    ScmCommitHash
	Author  ScmCommitAuthor
	Date    time.Time
	Message string
}

func NewScmCommit(hash ScmCommitHash, author ScmCommitAuthor, date time.Time, message string) *ScmCommit {
	return &ScmCommit{
		Hash:    hash,
		Author:  author,
		Date:    date,
		Message: message,
	}
}
