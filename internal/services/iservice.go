package services

type IService interface {
	Prepare() error
	Run() error
}
