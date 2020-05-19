package rbac

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/vickydk/utl/constants"
	rds "github.com/vickydk/utl/dbhandler/redis"
	"github.com/vickydk/utl/helper"
	"github.com/vickydk/utl/model"
	"github.com/vickydk/utl/rbac"
	"github.com/vickydk/utl/structs"
	"net/http"
	"strings"
)

type Route struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Name   string `json:"name"`
}

func New() *Service {
	return &Service{rb: rbac.GetRBAC()}
}

type Service struct {
	rb *rbac.RBAC
}

func (j *Service) MWFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var permission string
			var buf strings.Builder

			if _, ok := j.rb.Roles[c.Get("role").(string)]; !ok {
				if rds.ValidateData(constants.RoleMRdsPfx + c.Get("role").(string)) {
					j.rb.AddRole(c.Get("role").(string))
				} else {
					resp := helper.Respond(nil, nil, http.StatusForbidden)
					return c.JSON(resp.Code, resp)
				}
			}
			if j.rb.IsGranted(c.Get("role").(string), j.rb.Permissions["all"], nil) {
				return next(c)
			} else {
				for _, eachRoute := range c.Echo().Routes() {
					routes := new(Route)
					structs.Merge(routes, eachRoute)
					if routes.Method == c.Request().Method && routes.Path == c.Path() {
						spl := strings.Split(routes.Name, ".")
						splP := strings.Split(c.Path(), "/")
						if splP[1] == "v1" {
							buf.WriteString(splP[2])
						} else {
							buf.WriteString(splP[1])
						}
						buf.WriteString("_")
						buf.WriteString(strings.Split(spl[len(spl)-1], "-")[0])
						permission = buf.String()
						fmt.Println(permission)
						break
					}
				}

				if !j.rb.IsGranted(c.Get("role").(string), j.rb.Permissions[permission], nil) {
					if rds.ValidateData(constants.RoleMRdsPfx + c.Get("role").(string)) {
						var r model.Rbac
						rds.GetKeyValueByField(constants.RoleMRdsPfx+c.Get("role").(string), &r)
						p := strings.Split(r.Permission, ",")
						for _, pid := range p {
							if permission == pid {
								j.rb.UpdateRolePermission(c.Get("role").(string), permission, "add")
								return next(c)
							}
						}
					}
					resp := helper.Respond(nil, nil, http.StatusForbidden)
					return c.JSON(resp.Code, resp)
				}

				return next(c)
			}
		}
	}
}
