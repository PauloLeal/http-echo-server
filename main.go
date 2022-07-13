package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	serverPort = kingpin.Flag("port", "Web server http port.").Short('p').Default("8008").NoEnvar().Int()
)

type EchoResponse struct {
	Method        string            `json:"method"`
	Path          string            `json:"path"`
	Payload       string            `json:"payload"`
	ParsedPayload any               `json:"parsed_payload"`
	Headers       map[string]string `json:"headers"`
	QueryParams   map[string]string `json:"query_parameters"`
}

func echoHandler(ctx iris.Context) {
	er := EchoResponse{}

	er.Method = ctx.Method()
	er.Path = ctx.Path()

	body, _ := ctx.GetBody()
	er.Payload = string(body)
	json.Unmarshal(body, &er.ParsedPayload)

	er.QueryParams = ctx.URLParams()
	er.Headers = make(map[string]string)

	for k := range ctx.Request().Header {
		er.Headers[strings.ToLower(k)] = ctx.GetHeader(k)
	}

	b, _ := json.MarshalIndent(er, "", "  ")

	fmt.Printf("Method: %s\nPath: %s", er.Method, er.Path)
	fmt.Println(string(b))
	fmt.Println("==========================")

	ctx.StatusCode(200)
	ctx.JSON(er)
}

func main() {
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	app := iris.New()

	app.Handle("ALL", "/*", echoHandler)

	app.Use(recover.New())
	app.Run(iris.Addr(fmt.Sprintf(":%d", *serverPort)), iris.WithoutServerError(iris.ErrServerClosed))
}
