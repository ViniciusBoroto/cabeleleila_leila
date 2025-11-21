package mocks

//go:generate mockgen -source=../repository/service_repository.go -destination=mock_service_repository.go -package=mocks
//go:generate mockgen -source=../service/service_service.go -destination=mock_service_service.go -package=mocks
