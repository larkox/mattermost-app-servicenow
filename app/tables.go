package app

import (
	"github.com/mattermost/mattermost-plugin-apps/apps"

	"github.com/mattermost/mattermost-app-servicenow/config"
	"github.com/mattermost/mattermost-app-servicenow/constants"
	"github.com/mattermost/mattermost-app-servicenow/utils"
)

func GetTablesBindings(cc *apps.Context) (post, command, header *apps.Binding) {
	siteURL := cc.MattermostSiteURL
	appID := cc.AppID

	pt, ct, ht := filterTables(config.GetTables(cc))
	pb := baseBinding(siteURL, "Create Ticket", appID)
	cb := baseBinding(siteURL, "create-ticket", appID)
	hb := baseBinding(siteURL, "Create Ticket", appID)
	post = subBindings(siteURL, pt, pb, false, appID)
	command = subBindings(siteURL, ct, cb, true, appID)
	header = subBindings(siteURL, ht, hb, false, appID)

	return
}

func baseBinding(siteURL, label string, appID apps.AppID) *apps.Binding {
	return &apps.Binding{
		Location: constants.LocationCreate,
		Label:    label,
		Icon:     utils.GetIconURL(siteURL, "now-mobile-icon.png", appID),
		Bindings: []*apps.Binding{},
	}
}

func subBindings(siteURL string, tt config.TablesConfig, base *apps.Binding, useLocationLabel bool, appID apps.AppID) *apps.Binding {
	switch len(tt) {
	case 0:
		return nil
	case 1:
		for _, t := range tt {
			base.Call = &apps.Call{
				Path: t.ID,
			}

			return base
		}
	}

	for _, t := range tt {
		label := t.DisplayName
		if useLocationLabel {
			label = t.ID
		}

		base.Bindings = append(base.Bindings, &apps.Binding{
			Location: apps.Location(t.ID),
			Label:    label,
			Icon:     utils.GetIconURL(siteURL, "now-mobile-icon.png", appID),
			Call: &apps.Call{
				Path: t.ID,
			},
		})
	}

	return base
}

func filterTables(tt config.TablesConfig) (post, command, header config.TablesConfig) {
	post = config.TablesConfig{}
	command = config.TablesConfig{}
	header = config.TablesConfig{}

	for _, t := range tt {
		if t.Post {
			post[t.ID] = t
		}

		if t.Command {
			command[t.ID] = t
		}

		if t.Header {
			header[t.ID] = t
		}
	}

	return
}
