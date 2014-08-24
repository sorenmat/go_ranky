package playerservice

import (
	"net/http"

	"code.google.com/p/go-uuid/uuid"

	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id, Name string
}

func New() *restful.WebService {
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

func FindUser(id string) User {
	result := User{}

	c, session := getUserCollection()
	defer session.Close()

	err := c.Find(bson.M{"id": id}).One(&result)
	if err == nil {
		panic("Unable to find user with id: " + id)
	}
	return result
}
func FindUserService(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")

	result := User{}
	c, session := getUserCollection()
	defer session.Close()

	err := c.Find(bson.M{"id": id}).One(&result)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	} else {

		response.WriteEntity(result)
	}
}

func FindAllUsers(request *restful.Request, response *restful.Response) {
	var (
		results []User
	)

	c, session := getUserCollection()
	defer session.Close()
	c.Find(bson.M{}).All(&results)
	response.WriteEntity(results)
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

func getUserCollection() (*mgo.Collection, *mgo.Session) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	//	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("ranky").C("players")
	return c, session
}

func SaveUser(usr User) {
	c, session := getUserCollection()
	defer session.Close()

	err := c.Insert(usr)
	if err != nil {
		panic(err)
	}
}

func CreateUser(request *restful.Request, response *restful.Response) {
	uuid := uuid.New()
	//	fmt.Println("Trying to create user ", uuid)
	usr := User{Id: uuid}

	err := request.ReadEntity(&usr)
	//	fmt.Println("Usr: ", usr.Id)
	SaveUser(usr)
	// here you would create the user with some persistence system
	if err == nil {
		response.WriteEntity(usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

func RemoveUser(request *restful.Request, response *restful.Response) {
	// here you would delete the user from some persistence system
}
