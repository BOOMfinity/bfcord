package slash

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/api"
	"github.com/BOOMfinity/golog"
	"github.com/andersfylling/snowflake/v5"
	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"
	"strings"
)

type apiClient struct {
	client *api.Client
	appID  snowflake.ID
}

func (a *apiClient) Global() GlobalQuery {
	return apiGlobalQuery{api: a}
}

func (a *apiClient) Guild(id snowflake.ID) GuildQuery {
	return apiGlobalQuery{api: a, id: id}
}

type apiGlobalQuery struct {
	api *apiClient
	id  snowflake.ID
}

func (a apiGlobalQuery) Permissions() (list PermissionList, err error) {
	req := a.api.client.New(true)
	req.SetRequestURI(a.queryURL("permissions"))
	err = a.api.client.DoResult(req, &list)
	if err != nil {
		return
	}
	return
}

func (a apiGlobalQuery) CommandPermissions(id snowflake.ID) (list PermissionList, err error) {
	req := a.api.client.New(true)
	req.SetRequestURI(a.queryURL(id.String(), "permissions"))
	err = a.api.client.DoResult(req, &list)
	if err != nil {
		return
	}
	return
}

func (a apiGlobalQuery) EditPermissions(id snowflake.ID, perms PermissionList) (err error) {
	data, err := json.Marshal(perms)
	if err != nil {
		return
	}
	req := a.api.client.New(true)
	req.SetBody(data)
	req.SetRequestURI(a.queryURL(id.String(), "permissions"))
	req.Header.SetMethod(fasthttp.MethodPut)
	return a.api.client.DoNoResp(req)
}

func (a apiGlobalQuery) queryURL(path ...string) string {
	if a.id.IsZero() {
		return fmt.Sprintf("%v/applications/%v/commands/%v", api.FullApiUrl, a.api.appID, strings.Join(path, "/"))
	} else {
		return fmt.Sprintf("%v/applications/%v/guilds/%v/commands/%v", api.FullApiUrl, a.api.appID, a.id, strings.Join(path, "/"))
	}
}

func (a apiGlobalQuery) List() (list []Command, err error) {
	req := a.api.client.New(true)
	req.SetRequestURI(a.queryURL(""))
	err = a.api.client.DoResult(req, &list)
	if err != nil {
		return
	}
	return
}

func (a apiGlobalQuery) Create(name, desc string) CommandBuilder {
	return newCommandBuilder(a.api.client, a.api.appID, a.id, name, desc)
}

func (a apiGlobalQuery) Get(id snowflake.ID) (Command, error) {
	var cmd Command
	req := a.api.client.New(true)
	req.SetRequestURI(a.queryURL(id.String()))
	err := a.api.client.DoResult(req, &cmd)
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

func (a apiGlobalQuery) Edit(id snowflake.ID) EditCommandBuilder {
	return newEditCommandBuilder(a.api.client, a.api.appID, a.id, id)
}

func (a apiGlobalQuery) Delete(id snowflake.ID) error {
	req := a.api.client.New(true)
	req.SetRequestURI(a.queryURL(id.String()))
	req.Header.SetMethod(fasthttp.MethodDelete)
	return a.api.client.DoNoResp(req)
}

func NewClient(token string) Query {
	cl := api.NewClient(token, api.WithLogger(golog.New("slash-api")))
	app, err := cl.CurrentUser()
	if err != nil {
		panic(fmt.Errorf("could not fetch current user: %w", err))
	}
	return &apiClient{
		client: cl,
		appID:  app.ID,
	}
}

func NewClientWithAppID(token string, appID snowflake.ID) Query {
	return &apiClient{
		client: api.NewClient(token, api.WithLogger(golog.New("slash-api"))),
		appID:  appID,
	}
}
