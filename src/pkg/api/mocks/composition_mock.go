// Code generated by MockGen. DO NOT EDIT.
// Source: composition.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	model "hms/gateway/pkg/docs/model"
	processing "hms/gateway/pkg/docs/service/processing"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockCompositionService is a mock of CompositionService interface.
type MockCompositionService struct {
	ctrl     *gomock.Controller
	recorder *MockCompositionServiceMockRecorder
}

// MockCompositionServiceMockRecorder is the mock recorder for MockCompositionService.
type MockCompositionServiceMockRecorder struct {
	mock *MockCompositionService
}

// NewMockCompositionService creates a new mock instance.
func NewMockCompositionService(ctrl *gomock.Controller) *MockCompositionService {
	mock := &MockCompositionService{ctrl: ctrl}
	mock.recorder = &MockCompositionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCompositionService) EXPECT() *MockCompositionServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCompositionService) Create(ctx context.Context, userID, systemID string, ehrUUID, groupAccessUUID *uuid.UUID, composition *model.Composition, procRequest *processing.Request) (*model.Composition, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, userID, systemID, ehrUUID, groupAccessUUID, composition, procRequest)
	ret0, _ := ret[0].(*model.Composition)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockCompositionServiceMockRecorder) Create(ctx, userID, systemID, ehrUUID, groupAccessUUID, composition, procRequest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCompositionService)(nil).Create), ctx, userID, systemID, ehrUUID, groupAccessUUID, composition, procRequest)
}

// DefaultGroupAccess mocks base method.
func (m *MockCompositionService) DefaultGroupAccess() *model.GroupAccess {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DefaultGroupAccess")
	ret0, _ := ret[0].(*model.GroupAccess)
	return ret0
}

// DefaultGroupAccess indicates an expected call of DefaultGroupAccess.
func (mr *MockCompositionServiceMockRecorder) DefaultGroupAccess() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DefaultGroupAccess", reflect.TypeOf((*MockCompositionService)(nil).DefaultGroupAccess))
}

// DeleteByID mocks base method.
func (m *MockCompositionService) DeleteByID(ctx context.Context, procRequest *processing.Request, ehrUUID *uuid.UUID, versionUID, userID, systemID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", ctx, procRequest, ehrUUID, versionUID, userID, systemID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteByID indicates an expected call of DeleteByID.
func (mr *MockCompositionServiceMockRecorder) DeleteByID(ctx, procRequest, ehrUUID, versionUID, userID, systemID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockCompositionService)(nil).DeleteByID), ctx, procRequest, ehrUUID, versionUID, userID, systemID)
}

// GetByID mocks base method.
func (m *MockCompositionService) GetByID(ctx context.Context, userID, systemID string, ehrUUID *uuid.UUID, versionUID string) (*model.Composition, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, userID, systemID, ehrUUID, versionUID)
	ret0, _ := ret[0].(*model.Composition)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockCompositionServiceMockRecorder) GetByID(ctx, userID, systemID, ehrUUID, versionUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockCompositionService)(nil).GetByID), ctx, userID, systemID, ehrUUID, versionUID)
}

// GetLastByBaseID mocks base method.
func (m *MockCompositionService) GetLastByBaseID(ctx context.Context, userID, systemID string, ehrUUID *uuid.UUID, versionUID string) (*model.Composition, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastByBaseID", ctx, userID, systemID, ehrUUID, versionUID)
	ret0, _ := ret[0].(*model.Composition)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastByBaseID indicates an expected call of GetLastByBaseID.
func (mr *MockCompositionServiceMockRecorder) GetLastByBaseID(ctx, userID, systemID, ehrUUID, versionUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastByBaseID", reflect.TypeOf((*MockCompositionService)(nil).GetLastByBaseID), ctx, userID, systemID, ehrUUID, versionUID)
}

// GetList mocks base method.
func (m *MockCompositionService) GetList(ctx context.Context, userID, systemID string, ehrUUID *uuid.UUID) ([]*model.EhrDocumentItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetList", ctx, userID, systemID, ehrUUID)
	ret0, _ := ret[0].([]*model.EhrDocumentItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetList indicates an expected call of GetList.
func (mr *MockCompositionServiceMockRecorder) GetList(ctx, userID, systemID, ehrUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetList", reflect.TypeOf((*MockCompositionService)(nil).GetList), ctx, userID, systemID, ehrUUID)
}

// IsExist mocks base method.
func (m *MockCompositionService) IsExist(ctx context.Context, userID, systemID, ehrUUID, ID string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsExist", ctx, userID, systemID, ehrUUID, ID)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsExist indicates an expected call of IsExist.
func (mr *MockCompositionServiceMockRecorder) IsExist(ctx, userID, systemID, ehrUUID, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsExist", reflect.TypeOf((*MockCompositionService)(nil).IsExist), ctx, userID, systemID, ehrUUID, ID)
}

// Update mocks base method.
func (m *MockCompositionService) Update(ctx context.Context, procRequest *processing.Request, userID, systemID string, ehrUUID, groupAccessUUID *uuid.UUID, composition *model.Composition) (*model.Composition, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, procRequest, userID, systemID, ehrUUID, groupAccessUUID, composition)
	ret0, _ := ret[0].(*model.Composition)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockCompositionServiceMockRecorder) Update(ctx, procRequest, userID, systemID, ehrUUID, groupAccessUUID, composition interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCompositionService)(nil).Update), ctx, procRequest, userID, systemID, ehrUUID, groupAccessUUID, composition)
}

// MockIndexer is a mock of Indexer interface.
type MockIndexer struct {
	ctrl     *gomock.Controller
	recorder *MockIndexerMockRecorder
}

// MockIndexerMockRecorder is the mock recorder for MockIndexer.
type MockIndexerMockRecorder struct {
	mock *MockIndexer
}

// NewMockIndexer creates a new mock instance.
func NewMockIndexer(ctrl *gomock.Controller) *MockIndexer {
	mock := &MockIndexer{ctrl: ctrl}
	mock.recorder = &MockIndexerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIndexer) EXPECT() *MockIndexerMockRecorder {
	return m.recorder
}

// GetEhrUUIDByUserID mocks base method.
func (m *MockIndexer) GetEhrUUIDByUserID(ctx context.Context, userID, systemID string) (*uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEhrUUIDByUserID", ctx, userID, systemID)
	ret0, _ := ret[0].(*uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEhrUUIDByUserID indicates an expected call of GetEhrUUIDByUserID.
func (mr *MockIndexerMockRecorder) GetEhrUUIDByUserID(ctx, userID, systemID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEhrUUIDByUserID", reflect.TypeOf((*MockIndexer)(nil).GetEhrUUIDByUserID), ctx, userID, systemID)
}

// MockProcessingService is a mock of ProcessingService interface.
type MockProcessingService struct {
	ctrl     *gomock.Controller
	recorder *MockProcessingServiceMockRecorder
}

// MockProcessingServiceMockRecorder is the mock recorder for MockProcessingService.
type MockProcessingServiceMockRecorder struct {
	mock *MockProcessingService
}

// NewMockProcessingService creates a new mock instance.
func NewMockProcessingService(ctrl *gomock.Controller) *MockProcessingService {
	mock := &MockProcessingService{ctrl: ctrl}
	mock.recorder = &MockProcessingServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProcessingService) EXPECT() *MockProcessingServiceMockRecorder {
	return m.recorder
}

// NewRequest mocks base method.
func (m *MockProcessingService) NewRequest(reqID, userID, ehrUUID string, kind processing.RequestKind) (*processing.Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRequest", reqID, userID, ehrUUID, kind)
	ret0, _ := ret[0].(*processing.Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewRequest indicates an expected call of NewRequest.
func (mr *MockProcessingServiceMockRecorder) NewRequest(reqID, userID, ehrUUID, kind interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRequest", reflect.TypeOf((*MockProcessingService)(nil).NewRequest), reqID, userID, ehrUUID, kind)
}
