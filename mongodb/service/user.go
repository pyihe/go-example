package service

type UserAPI interface {
}

type UserService struct {
	api UserAPI
}

func NewUserService(api UserAPI) *UserService {
	return &UserService{api: api}
}

func (u *UserService) Work() {

}

func (u *UserService) InsertOne() {

}
