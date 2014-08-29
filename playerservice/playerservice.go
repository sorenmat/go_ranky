package playerservice

import (
	"errors"
	"net/http"

	"code.google.com/p/go-uuid/uuid"

	"github.com/emicklei/go-restful"
)

type User struct {
	Id   string
	Name string
}

var (
	repository PlayerRepository
)

func New(repo PlayerRepository) *restful.WebService {
	repository = repo // set global state.. ohh noo
	service := new(restful.WebService)
	service.
		Path("/players").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	service.Route(service.GET("/").To(FindAllUsers))
	service.Route(service.GET("/{user-id}").To(FindUserService))
	service.Route(service.POST("").To(UpdateUser))
	service.Route(service.PUT("").To(CreateUser))
	service.Route(service.DELETE("/{user-id}").To(RemoveUser))

	return service
}

func FindUserService(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")

	result, err := repository.FindUser(id)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	} else {

		response.WriteEntity(result)
	}
}

func FindAllUsers(request *restful.Request, response *restful.Response) {
	response.WriteEntity(repository.FindAllUsers())
}

func UpdateUser(request *restful.Request, response *restful.Response) {
	usr := new(User)
	err := request.ReadEntity(&usr)
	// here you would update the user with some persistence system
	if err == nil {
		response.WriteEntity(usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

func CreateUser(request *restful.Request, response *restful.Response) {
	uuid := uuid.New()
	//	fmt.Println("Trying to create user ", uuid)
	usr := User{Id: uuid}

	err := request.ReadEntity(&usr)
	if repository.IsUserUnique(usr) {
		usr.Id = uuid
		//	fmt.Println("Usr: ", usr.Id)
		repository.SaveUser(usr)
		// here you would create the user with some persistence system
		if err == nil {
			response.WriteEntity(usr)
		} else {
			response.WriteError(http.StatusInternalServerError, err)
		}
	} else {
		response.WriteError(http.StatusConflict, errors.New("Player name not unique"))

	}
}

func RemoveUser(request *restful.Request, response *restful.Response) {
	// here you would delete the user from some persistence system
}
