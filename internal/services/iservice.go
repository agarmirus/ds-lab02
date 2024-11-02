package services

type IService interface {
	Prepare(string) error
	Run() error
}
