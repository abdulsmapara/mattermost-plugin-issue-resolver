package main

import (
	"fmt"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"strings"
)

// ExecuteCommand executes the commands registered via RegisterCommand hook
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	trigger := strings.TrimPrefix(strings.Fields(args.Command)[0], "/")

	switch trigger {

	case commandUpdateSkills:
		return p.executeCommandSkills(args), nil

	case commandResolve:
		return p.resolveIssue(args), nil

	default:
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Unknown Command: " + args.Command),
		}, nil
	}

}

// resolveIssue first searches on stackoverflow. If failure encountered, a user is searched to be suggested.
func (p *Plugin) resolveIssue(args *model.CommandArgs) *model.CommandResponse {
	// Get the user's query
	userQuery := strings.TrimPrefix(args.Command, "/"+commandResolve)

	// Search on Stackoverflow
	response := searchStackOverflow(userQuery)

	// Prepare message to be returned
	msg := "###### QUERY: \n" + userQuery + "\n"

	// Check if stackoverflow search succeeded
	if response["Found"] == "true" {

		// Return the question found and the accepted solution
		msg += "###### QUESTION FOUND:\n" + response["Question Title"] + "\n" + response["Question Body"] + "\n" + "###### Solution:\n" + response["Solution"]
	} else {

		// Get skilled users if stackoverflow search fails
		skilledUserInfo := p.getSkilledUsers(userQuery, args.UserId)

		// Check if skilled user is found
		if skilledUserInfo["Found"] == "true" {

			// Return the suggestion of the skilled User
			msg += "###### SUGGESTION: \n" + skilledUserInfo["Message"]
		} else {

			// Skilled user not found. Issue could not be resolved.
			msg += "Sorry, the issue could not be resolved.\nA few tips that will help me:\nTry asking me by including specific keywords.\nFor multiword skills, try to seperate words with underscore. Example: Java_Servlets."
		}
	}

	// Create a post to be returned
	post := &model.Post{
		ChannelId: args.ChannelId,
		RootId:    args.RootId,
		UserId:    p.issueResolverBotID,
		Message:   msg,
	}

	// Call CreatePost API to create post in the channel
	_, err := p.API.CreatePost(post)
	if err != nil {
		errorMessage := "Failed to create post"
		p.API.LogError(errorMessage, "err", err.Error())
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         errorMessage,
		}
	}

	return &model.CommandResponse{}
}

// executeCommandSkills handles the functionality of /skills slash command
func (p *Plugin) executeCommandSkills(args *model.CommandArgs) *model.CommandResponse {

	fields := strings.Split(args.Command, " ")
	command := ""
	if len(fields) >= 2 {
		command = fields[1]
		command = strings.TrimSpace(command)
	}

	skillStr := p.allSkills.String()
	arrayOfSkills := skillStr[4 : len(skillStr)-1]

	switch command {

	case "":
		userSkills, err := p.API.KVGet(args.UserId)
		if err != nil {
			errorMessage := "Failed to fetch your skills !"
			p.API.LogError(errorMessage, "err", err.Error())
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         errorMessage,
			}
		}
		yourSkills := "Your skills: \n [ "
		if userSkills != nil {
			yourSkills += string(userSkills)
		}
		yourSkills += " ]"
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         yourSkills,
		}
	case "list":
		returnMessage := ""
		pattern := strings.TrimPrefix(args.Command, "/"+commandUpdateSkills+" list")
		pattern = strings.TrimSpace(pattern)
		for _, skill := range strings.Split(arrayOfSkills, ",") {
			if pattern == "" {
				returnMessage += skill + ","
			} else {
				// fmt.Printf("[ " + skill + ", " + pattern + " ]")
				// fmt.Println(strings.HasPrefix(strings.TrimSpace(skill), strings.ToUpper(pattern)))
				if strings.HasPrefix(strings.TrimSpace(skill), strings.ToUpper(pattern)) {
					returnMessage += skill + ","
				}
			}
		}
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Available skills for the pattern " + pattern + " [ " + returnMessage + " ]",
		}

	case "add":

		if len(fields) < 3 {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         fmt.Sprintf("No skill mentioned. Mention one/more skill out of skills listed by list command"),
			}
		}

		skillsMentioned := strings.TrimPrefix(args.Command, "/"+commandUpdateSkills+" add ")
		skillsMentionedArray := strings.Split(skillsMentioned, ",")

		skillsForUser, err := p.API.KVGet(args.UserId)
		if err != nil {
			errorMessage := "Failed to add skills !"
			p.API.LogError(errorMessage, "err", err.Error())
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         errorMessage,
			}
		}
		userSkills := ""
		if skillsForUser != nil {
			userSkills = string(skillsForUser)
		}
		for _, skill := range skillsMentionedArray {
			skill = strings.ToUpper(strings.TrimSpace(skill))
			if p.allSkills.Contains(skill) {
				alreadyExists := false
				for _, existingSkill := range strings.Split(userSkills, ",") {
					if existingSkill == skill {
						alreadyExists = true
						break
					}
				}
				if !alreadyExists {
					userSkills += skill + ","
				}
			}

		}

		err = p.API.KVSet(args.UserId, []byte(userSkills))
		if err != nil {
			errorMessage := "Failed to add skills !"
			p.API.LogError(errorMessage, "err", err.Error())
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         errorMessage,
			}
		}

		returnMessage := "Skills Updated Successfully !"
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         returnMessage,
		}
	case "delete":
		if len(fields) < 3 {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         fmt.Sprintf("No skill mentioned. Mention one/more skills that you want to delete"),
			}
		}
		skillsMentioned := strings.TrimPrefix(args.Command, "/"+commandUpdateSkills+" delete ")
		skillsMentionedArray := strings.Split(skillsMentioned, ",")
		skillsForUser, err := p.API.KVGet(args.UserId)
		if err != nil {
			errorMessage := "Failed to delete skills !"
			p.API.LogError(errorMessage, "err", err.Error())
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         errorMessage,
			}
		}
		if skillsForUser != nil {
			userSkills := string(skillsForUser)
			updatedSkills := ""
			userSkillsArray := strings.Split(userSkills, ",")
			for _, existingSkill := range userSkillsArray {
				if existingSkill != "" {
					skillExists := false
					for _, skill := range skillsMentionedArray {
						skill = strings.ToUpper(strings.TrimSpace(skill))
						if skill == existingSkill {
							skillExists = true
							break
						}
					}
					if !skillExists {
						updatedSkills += existingSkill + ","
					}
				}
			}
			err = p.API.KVDelete(args.UserId)
			if err != nil {
				errorMessage := "Failed to delete skills !"
				p.API.LogError(errorMessage, "err", err.Error())
				return &model.CommandResponse{
					ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
					Text:         errorMessage,
				}
			}
			if updatedSkills != "" {
				err = p.API.KVSet(args.UserId, []byte(updatedSkills))
				if err != nil {
					errorMessage := "Failed to delete skills !"
					p.API.LogError(errorMessage, "err", err.Error())
					return &model.CommandResponse{
						ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
						Text:         errorMessage,
					}
				}
			}
		}
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "If you previously held skills that you mentioned, they have been deleted !",
		}

	default:
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Incorrect usage of /skills. Command not recognised: " + command),
		}
	}

}
