package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

// See https://developers.mattermost.com/extend/plugins/server/reference/

func (p *Plugin) OnActivate() error {
	_, err := p.Helpers.EnsureBot(&model.Bot{
		Username:    "somebot",
		Description: "A bot to test MessageHasBeenPosted",
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	ok, err := p.Helpers.ShouldProcessMessage(
		post,
		//plugin.AllowSystemMessages(),
		//plugin.AllowBots(),
		//plugin.OnlyBotDMs(),
		// The following two are not so easy to enable because you have to provide a list of real user/channel ids
		//plugin.FilterChannelIDs(),
		//plugin.FilterUserIDs(),
	)

	if err != nil {
		p.API.LogError("failed to check message in MessageHasBeenPosted", "err", err)
		return
	}

	if !ok {
		return
	}

	p.API.LogWarn("A user posted a message", "user id", post.UserId)
}
