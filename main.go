package main

import (
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/encoder"
	"github.com/sorenmat/ranky/playerservice"
)

var rootdir = "/Users/soren/go/src/github.com/sorenmat/ranky/static/app"

func main() {
	m := martini.Classic()
	// map json encoder
	m.Use(martini.Static("static/app"))
	m.Use(func(c martini.Context, w http.ResponseWriter, r *http.Request) {
		pretty, _ := strconv.ParseBool(r.FormValue("pretty"))

		c.MapTo(encoder.JsonEncoder{PrettyPrint: pretty}, (*encoder.Encoder)(nil))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	})

	m.Get("/", func() string {
		return "Hello world!"
	})

	m.Get("/players", func(enc encoder.Encoder) (int, []byte) {
		players := playerservice.FindAllUsersInDb()
		return http.StatusOK, encoder.Must(enc.Encode(players))

	})
	m.Run()
	/*
		restful.Add(playerservice.New())
		restful.Add(matchservice.New())

		// add static resources..
		ws := new(restful.WebService)
		ws.Route(ws.GET("/static/{resource}").To(staticFromPathParam))
		ws.Route(ws.GET("/static").To(staticFromQueryParam))
		restful.Add(ws)
		log.Fatal(http.ListenAndServe(":8080", nil))
	*/
}
func staticFromQueryParam(req *restful.Request, resp *restful.Response) {
	fmt.Printf("staticFromQueryParam\n")
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		path.Join(rootdir, req.QueryParameter("resource")))
}

func staticFromPathParam(req *restful.Request, resp *restful.Response) {
	fmt.Println(req.PathParameters)
	for k, v := range req.PathParameters() {
		fmt.Printf("%s == %s\n", k, v)
	}
	resource := req.PathParameter("resource")
	joinedPath := path.Join(rootdir, resource)
	fmt.Printf("RootDir %s\n", rootdir)
	fmt.Printf("JoinPath %s\n", joinedPath)
	fmt.Printf("Resource %s\n", resource)
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		joinedPath)
}
