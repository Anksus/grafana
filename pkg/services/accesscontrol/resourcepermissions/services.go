package resourcepermissions

import (
	"context"
	"errors"
	"strconv"

	"github.com/grafana/grafana/pkg/api/routing"
	ac "github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/sqlstore"
)

var dashboardsView = []string{ac.ActionDashboardsRead}
var dashboardsEdit = append(dashboardsView, []string{ac.ActionDashboardsWrite, ac.ActionDashboardsDelete, ac.ActionDashboardsEdit}...)
var dashboardsAdmin = append(dashboardsEdit, []string{ac.ActionDashboardsPermissionsRead, ac.ActionDashboardsPermissionsWrite}...)
var foldersView = []string{ac.ActionFoldersRead}
var foldersEdit = append(foldersView, []string{ac.ActionFoldersWrite, ac.ActionFoldersDelete, ac.ActionFoldersEdit, ac.ActionDashboardsCreate}...)
var foldersAdmin = append(foldersEdit, []string{ac.ActionFoldersPermissionsRead, ac.ActionFoldersPermissionsWrite}...)

func ProvideServices(sql *sqlstore.SQLStore, router routing.RouteRegister, ac ac.AccessControl, store Store) (*Services, error) {
	dashboardsService, err := provideDashboardService(sql, router, ac, store)
	if err != nil {
		return nil, err
	}

	folderService, err := provideFolderService(sql, router, ac, store)
	if err != nil {
		return nil, err
	}

	return &Services{services: map[string]*Service{
		"folders":    folderService,
		"dashboards": dashboardsService,
	}}, nil
}

type Services struct {
	services map[string]*Service
}

func (s *Services) GetDashboardService() *Service {
	return s.services["dashboards"]
}

func (s *Services) GetFolderService() *Service {
	return s.services["folders"]
}

func provideDashboardService(sql *sqlstore.SQLStore, router routing.RouteRegister, accesscontrol ac.AccessControl, store Store) (*Service, error) {
	options := Options{
		Resource: "dashboards",
		ResourceValidator: func(ctx context.Context, orgID int64, resourceID string) error {
			id, err := strconv.ParseInt(resourceID, 10, 64)
			if err != nil {
				return err
			}

			if _, err := sql.GetDashboard(id, orgID, "", ""); err != nil {
				return err
			}
			return nil
		},
		Assignments: Assignments{
			Users:        true,
			Teams:        true,
			BuiltInRoles: true,
		},
		PermissionsToActions: map[string][]string{
			"View":  dashboardsView,
			"Edit":  dashboardsEdit,
			"Admin": dashboardsAdmin,
		},
		ReaderRoleName: "Dashboard permission reader",
		WriterRoleName: "Dashboard permission writer",
		RoleGroup:      "Dashboards",
	}

	return New(options, router, accesscontrol, store, sql)
}

func provideFolderService(sql *sqlstore.SQLStore, router routing.RouteRegister, accesscontrol ac.AccessControl, store Store) (*Service, error) {
	options := Options{
		Resource: "folders",
		ResourceValidator: func(ctx context.Context, orgID int64, resourceID string) error {
			id, err := strconv.ParseInt(resourceID, 10, 64)
			if err != nil {
				return err
			}
			if dashboard, err := sql.GetDashboard(id, orgID, "", ""); err != nil {
				return err
			} else if !dashboard.IsFolder {
				return errors.New("not found")
			}

			return nil
		},
		Assignments: Assignments{
			Users:        true,
			Teams:        true,
			BuiltInRoles: true,
		},
		PermissionsToActions: map[string][]string{
			"View":  append(dashboardsView, foldersView...),
			"Edit":  append(dashboardsEdit, foldersEdit...),
			"Admin": append(dashboardsAdmin, foldersAdmin...),
		},
		ReaderRoleName: "Folder permission reader",
		WriterRoleName: "Folder permission writer",
		RoleGroup:      "Folders",
	}

	return New(options, router, accesscontrol, store, sql)
}
