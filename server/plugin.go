package main

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// UserID of issueResolver search bot
	issueResolverBotID string

	// set of all skills that are predefined
	allSkills mapset.Set
}

// OnActivate is invoked when the plugin is activated. If an error is returned, the plugin will be terminated.
// The plugin will not receive hooks until after OnActivate returns without error.
// https://developers.mattermost.com/extend/plugins/server/reference/#Hooks.OnActivate
func (p *Plugin) OnActivate() error {

	// Initializa the set of skills with empty set
	p.allSkills = mapset.NewSet()

	// Register the command commandResolve
	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          commandResolve,
		AutoComplete:     true,
		AutoCompleteHint: "[ Your Question/Issue ]",
		AutoCompleteDesc: "Searches stackoverflow for similar/related question. Suggests users to contact for resolving the issue.",
	}); err != nil {
		errorMessage := "failed to register command " + commandResolve
		p.API.LogError(errorMessage, "err", err.Error())
		return errors.Wrapf(err, "failed to register %s command", commandResolve)
	}

	p.API.LogInfo(commandResolve+" command registered", map[string]string{"Info": "Command Registered"})

	// Register the command commandUpdateSkills
	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          commandUpdateSkills,
		AutoComplete:     true,
		AutoCompleteHint: "[ Empty OR Commands ]",
		AutoCompleteDesc: "Empty returns your current skills, Command 'list' returns a list of all skills with given prefix , Command 'add' adds skills against your name, Command 'delete' deletes the skill against your name if it existed.",
	}); err != nil {
		p.API.LogError("Failed ro register command "+commandUpdateSkills, "err", err.Error())
		return errors.Wrapf(err, "failed to register %s command", commandUpdateSkills)
	}

	p.API.LogInfo(commandUpdateSkills+" command registered", map[string]string{"Info": "Command Registered"})

	issueResolverBot := &model.Bot{
		Username:    "issue_resolver",
		DisplayName: "ISSUE RESOLVER",
		Description: "Created by Mattermost Issue Resolver Plugin. Searches stackoverflow for questions posted by Mattermost users using /resolve. Also suggests users to contact for resolving the issue.",
	}

	// Ensure the bot. If not present create an Issue Resolver bot
	issueResolverBotID, err := p.Helpers.EnsureBot(issueResolverBot)
	if err != nil {
		p.API.LogError("Failed to ensure issue resolver bot ", "err", err.Error())
		return errors.Wrap(err, "Failed to ensure issue resolver bot")
	}

	p.API.LogInfo(" Issue Resolver Bot ensured", map[string]string{"Info": "Bot Ensured"})

	// Store created ID in Plugin struct
	p.issueResolverBotID = issueResolverBotID

	// Get the plugin file path
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		p.API.LogError("Failed get bundle path ", "err", err.Error())
		return errors.Wrap(err, "Could not get bundle path")
	}

	// Read the file that lists the predefined skills
	skillsData, err := ioutil.ReadFile(filepath.Join(bundlePath, "assets", "skills.txt"))
	skillsString := strings.ToUpper(string(skillsData))

	p.API.LogInfo("File skills.txt read successfully", map[string]string{"Info": "File read successfully"})

	// Add predefined skills to the set in the plugin
	for _, skill := range strings.Split(skillsString, "\n") {
		p.allSkills.Add(skill)
	}

	return nil
}
