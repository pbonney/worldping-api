package api

import (
	"encoding/json"
	"os"
	"path"

	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/metrics"
	"github.com/grafana/grafana/pkg/middleware"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/util"
)

func isDasboardStarredByUser(c *middleware.Context, dashId int64) (bool, error) {
	if !c.IsSignedIn {
		return false, nil
	}

	query := m.IsStarredByUserQuery{UserId: c.UserId, DashboardId: dashId}
	if err := bus.Dispatch(&query); err != nil {
		return false, err
	}

	return query.Result, nil
}

func GetDashboard(c *middleware.Context) {
	metrics.M_Api_Dashboard_Get.Inc(1)

	slug := c.Params(":slug")

	query := m.GetDashboardQuery{Slug: slug, OrgId: c.OrgId}
	err := bus.Dispatch(&query)
	if err != nil {
		c.JsonApiErr(404, "Dashboard not found", nil)
		return
	}

	isStarred, err := isDasboardStarredByUser(c, query.Result.Id)
	if err != nil {
		c.JsonApiErr(500, "Error while checking if dashboard was starred by user", err)
		return
	}

	dash := query.Result
	dto := dtos.Dashboard{
		Model: dash.Data,
		Meta:  dtos.DashboardMeta{IsStarred: isStarred, Slug: slug},
	}

	c.JSON(200, dto)
}

func DeleteDashboard(c *middleware.Context) {
	slug := c.Params(":slug")

	query := m.GetDashboardQuery{Slug: slug, OrgId: c.OrgId}
	if err := bus.Dispatch(&query); err != nil {
		c.JsonApiErr(404, "Dashboard not found", nil)
		return
	}

	cmd := m.DeleteDashboardCommand{Slug: slug, OrgId: c.OrgId}
	if err := bus.Dispatch(&cmd); err != nil {
		c.JsonApiErr(500, "Failed to delete dashboard", err)
		return
	}

	var resp = map[string]interface{}{"title": query.Result.Title}

	c.JSON(200, resp)
}

func PostDashboard(c *middleware.Context, cmd m.SaveDashboardCommand) {
	cmd.OrgId = c.OrgId

	err := bus.Dispatch(&cmd)
	if err != nil {
		if err == m.ErrDashboardWithSameNameExists {
			c.JSON(412, util.DynMap{"status": "name-exists", "message": err.Error()})
			return
		}
		if err == m.ErrDashboardVersionMismatch {
			c.JSON(412, util.DynMap{"status": "version-mismatch", "message": err.Error()})
			return
		}
		c.JsonApiErr(500, "Failed to save dashboard", err)
		return
	}

	metrics.M_Api_Dashboard_Post.Inc(1)

	c.JSON(200, util.DynMap{"status": "success", "slug": cmd.Result.Slug, "version": cmd.Result.Version})
}

func GetHomeDashboard(c *middleware.Context) {
	filePath := path.Join(setting.StaticRootPath, "plugins/raintank/dashboards/Litmus-Home.json")
	file, err := os.Open(filePath)
	if err != nil {
		c.JsonApiErr(500, "Failed to load home dashboard", err)
		return
	}

	dash := dtos.Dashboard{}
	dash.Meta.IsHome = true
	jsonParser := json.NewDecoder(file)
	if err := jsonParser.Decode(&dash.Model); err != nil {
		c.JsonApiErr(500, "Failed to load home dashboard", err)
		return
	}

	c.JSON(200, &dash)
}
