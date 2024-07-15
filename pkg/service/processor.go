package processor

import (
	"errors"

	"vk/pkg/model"
	"vk/pkg/repository"
)

type Processor interface {
	Process(d *model.Document) (*model.Document, error)
}

type processorImpl struct {
	repo repository.Repository
}

func NewProcessor(repo repository.Repository) Processor {
	return &processorImpl{repo: repo}
}

func (p *processorImpl) Process(d *model.Document) (*model.Document, error) {

	// Lock the document
	if err := p.repo.LockDocument(d.Url); err != nil {
		return nil, err
	}
	defer p.repo.UnlockDocument(d.Url)

	existingDoc, err := p.repo.GetDocument(d.Url)
	if err != nil && !errors.Is(err, repository.ErrDocumentNotFound) {
		return nil, err
	}

	updatedDoc := mergeDocuments(existingDoc, d)

	if err := p.repo.SaveDocument(updatedDoc); err != nil {
		return nil, err
	}

	return updatedDoc, nil
}

func mergeDocuments(existingDoc, newDoc *model.Document) *model.Document {
	if existingDoc == nil {
		existingDoc = &model.Document{
			Url:            newDoc.Url,
			PubDate:        newDoc.PubDate,
			FetchTime:      newDoc.FetchTime,
			Text:           newDoc.Text,
			FirstFetchTime: newDoc.FetchTime,
		}
		return existingDoc
	}

	if newDoc.FetchTime > existingDoc.FetchTime {
		existingDoc.Text = newDoc.Text
		existingDoc.FetchTime = newDoc.FetchTime
	}

	if newDoc.FetchTime < existingDoc.FirstFetchTime {
		existingDoc.PubDate = newDoc.PubDate
		existingDoc.FirstFetchTime = newDoc.FetchTime
	}

	return existingDoc
}
