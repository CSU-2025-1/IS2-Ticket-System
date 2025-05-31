package queries

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/resty.v1"
	"log/slog"
	"math/rand"
	"net/http"
)

var (
	groups = make([]string, 0)
)

type PreparedQuery struct {
	req      *resty.Request
	method   string
	path     string
	prepare  func(query PreparedQuery) (PreparedQuery, error)
	callback func(response *resty.Response)
}

var (
	ReadQueries = []func(token string) *PreparedQuery{
		func(token string) *PreparedQuery {
			return &PreparedQuery{
				req:    resty.R().SetAuthToken(token),
				path:   "/api/users",
				method: http.MethodGet,
			}
		},
		func(token string) *PreparedQuery {
			return &PreparedQuery{
				req:    resty.R().SetAuthToken(token),
				path:   "/api/users/groups",
				method: http.MethodGet,
			}
		},
		func(token string) *PreparedQuery {
			return &PreparedQuery{
				req:  resty.R().SetAuthToken(token),
				path: "/api/users/groups/%s/list",
				prepare: func(query PreparedQuery) (PreparedQuery, error) {
					if len(groups) == 0 {
						return PreparedQuery{}, fmt.Errorf("no groups")
					}

					i := rand.Intn(len(groups)) - 1
					if i < 0 {
						i = 0
					}
					query.path = fmt.Sprintf(query.path, groups[i])
					return query, nil
				},
				method: http.MethodGet,
			}
		},
		func(token string) *PreparedQuery {
			return &PreparedQuery{
				req:    resty.R().SetAuthToken(token),
				path:   "/api/ticket",
				method: http.MethodGet,
			}
		},
	}

	WriteQueries = []func(token string) *PreparedQuery{
		func(token string) *PreparedQuery {
			return &PreparedQuery{
				req: resty.R().SetAuthToken(token),
				prepare: func(query PreparedQuery) (PreparedQuery, error) {
					query.req.SetBody(
						map[string]interface{}{
							"name":           "Имя",
							"description":    "Описания",
							"status":         "Новый",
							"created_by":     uuid.New().String(),
							"recipient_type": "user",
							"recipient_uuid": uuid.New().String(),
							"priority":       0,
						})
					return query, nil
				},
				path:   "/api/ticket",
				method: http.MethodPost,
			}
		},
		func(token string) *PreparedQuery {
			return &PreparedQuery{
				req: resty.R().SetAuthToken(token),
				prepare: func(query PreparedQuery) (PreparedQuery, error) {
					query.req.SetBody(
						map[string]interface{}{
							"name": "Группа внедрения цифровых продуктов",
						})
					return query, nil
				},
				path: "/api/users/groups",
				callback: func(response *resty.Response) {
					var res struct {
						Uuid string `json:"uuid"`
					}
					err := json.Unmarshal(response.Body(), &res)
					if err != nil {
						slog.Error(err.Error())
						return
					}

					groups = append(groups, res.Uuid)
				},
				method: http.MethodPost,
			}
		},
		func(token string) *PreparedQuery {
			return &PreparedQuery{
				req:    resty.R().SetAuthToken(token),
				path:   "/api/users",
				method: http.MethodPost,
				prepare: func(query PreparedQuery) (PreparedQuery, error) {
					query.req.SetBody(map[string]string{
						"login":    uuid.New().String(),
						"password": "123",
					})
					return query, nil
				},
			}
		},
	}
)

func DoRequest(target, token string, writeRatio float32) error {
	n := rand.Float32()

	var pquery *PreparedQuery
	var err error

	if n < writeRatio {
		pquery = WriteQueries[rand.Intn(len(WriteQueries))](token)
	} else {
		pquery = ReadQueries[rand.Intn(len(ReadQueries))](token)
	}

	var query PreparedQuery

	query = *pquery

	if pquery.prepare != nil {
		query, err = pquery.prepare(query)
		if err != nil {
			return nil
		}
	}

	var resp *resty.Response

	switch query.method {
	case http.MethodGet:
		resp, err = query.req.Get(target + query.path)
	case http.MethodPost:
		resp, err = query.req.Post(target + query.path)
	}

	if err != nil {
		return err
	}

	fmt.Printf("%s %s%s %d\n", query.method, target, query.path, resp.StatusCode())

	if query.callback != nil {
		query.callback(resp)
	}

	return nil
}
