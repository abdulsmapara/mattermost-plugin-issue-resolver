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
		p.API.LogInfo("executing command skills", map[string]string{"Info": "executing command skills"})
		return p.executeCommandSkills(args), nil

	case commandResolve:
		p.API.LogInfo("executing command resolve", map[string]string{"Info": "executing command resolve"})
		return p.resolveIssue(args), nil

	default:
		p.API.LogInfo("Slash command not recognised", map[string]string{"Info": "Slash command not recognised"})
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
	response := p.searchStackOverflow(userQuery)

	// Prepare message to be returned
	msg := "###### QUERY: \n" + userQuery + "\n"

	// Check if stackoverflow search succeeded
	if response["Found"] == "true" {

		// Return the question found and the accepted solution
		msg += "###### QUESTION FOUND:\n" + response["Question Title"] + "\n" + response["Question Body"] + "\n" + "###### Solution:\n" + response["Solution"]
		p.API.LogInfo("Search on stackoverflow successful", map[string]string{"Info": "Search on stackoverflow successful"})
	} else {

		p.API.LogInfo("Search on stackoverflow failed", map[string]string{"Info": "Search on stackoverflow failed"})

		// Get skilled users if stackoverflow search fails
		skilledUserInfo := p.getSkilledUsers(userQuery, args.UserId)

		// Check if skilled user is found
		if skilledUserInfo["Found"] == "true" {

			// Return the suggestion of the skilled User
			msg += "###### SUGGESTION: \n" + skilledUserInfo["Message"]

			p.API.LogInfo("Found a skilled user", map[string]string{"Info": "Found a skilled user"})
		} else {

			// Skilled user not found. Issue could not be resolved.
			msg += "Sorry, the issue could not be resolved.\n\nA few tips that will help me:\nTry asking me by including specific keywords.\nFor multiword skills, try to seperate words with underscore. Example: Java_Servlets."
			p.API.LogInfo("User issue could not be resolved", map[string]string{"Info": "User issue could not be resolved"})
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
	p.API.LogInfo("Result of /resolve posted by bot", map[string]string{"Info": "Result of /resolve posted by bot"})
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
		// Empty command - returns the skills of the user
		userSkills, err := p.API.KVGet(args.UserId)
		if err != nil {
			errorMessage := "Failed to fetch your skills !"
			p.API.LogError(errorMessage, "err", err.Error())
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         errorMessage,
			}
		}

		// prepare return message
		yourSkills := "Your skills: \n [ "
		if userSkills != nil {
			yourSkills += string(userSkills)
		}
		yourSkills += " ]"
		p.API.LogInfo("User skills reported successfully", map[string]string{"Info": "User skills reported successfully"})

		// return user's skills as stored in Database(KV)
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         yourSkills,
		}
	case "list":
		// list command returns the list of all skills with given prefix (if no prefix given return all)

		// prepare return message
		returnMessage := ""

		// get pattern (prefix) mentioned by the user
		pattern := strings.TrimPrefix(args.Command, "/"+commandUpdateSkills+" list")
		pattern = strings.TrimSpace(pattern)
		for _, skill := range strings.Split(arrayOfSkills, ",") {
			if pattern == "" {

				// if no pattern, include all
				returnMessage += skill + ","
			} else {

				// include all skills that have given pattern
				if strings.HasPrefix(strings.TrimSpace(skill), strings.ToUpper(pattern)) {
					returnMessage += skill + ","
				}
			}
		}

		p.API.LogInfo("Skills listed successfully", map[string]string{"Info": "Skills listed successfully"})

		// return the available skills for the given prefix
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Available skills for the prefix: " + pattern + " are [ " + returnMessage + " ]",
		}

	case "add":

		// add command adds the mentioned skills for the user
		if len(fields) < 3 {
			// no skill mentioned along with add command
			p.API.LogInfo("/skills add not used properly", map[string]string{"Info": "/skills add not used properly"})
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         fmt.Sprintf("No skill mentioned. Mention one/more skill out of skills listed by list command"),
			}

		}

		// Get the skills mentioned by the user
		skillsMentioned := strings.TrimPrefix(args.Command, "/"+commandUpdateSkills+" add ")
		skillsMentionedArray := strings.Split(skillsMentioned, ",")

		// Get skills of the user first
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
			// check if mentioned skill is present in the predefined set or not
			skill = strings.ToUpper(strings.TrimSpace(skill))
			if p.allSkills.Contains(skill) {
				// check if skill to be added is already present
				alreadyExists := false
				for _, existingSkill := range strings.Split(userSkills, ",") {
					if existingSkill == skill {
						alreadyExists = true
						break
					}
				}

				// if skill does not exist beforehand, add it
				if !alreadyExists {
					userSkills += skill + ","
				}
			}

		}

		// Update the skill of the user
		err = p.API.KVSet(args.UserId, []byte(userSkills))
		if err != nil {
			errorMessage := "Failed to add skills !"
			p.API.LogError(errorMessage, "err", err.Error())
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         errorMessage,
			}
		}

		// prepare return Message
		returnMessage := "Skills Updated Successfully !"
		p.API.LogInfo("Skills Updated successfully", map[string]string{"Info": "Skills Updated successfully"})

		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         returnMessage,
		}
	case "delete":
		// delete command deletes the mentioned skills for the user

		if len(fields) < 3 {
			// No skill mentioned with delete command
			p.API.LogInfo("/skills delete not used properly", map[string]string{"Info": "/skills delete not used properly"})
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         fmt.Sprintf("No skill mentioned. Mention one/more skills that you want to delete"),
			}
		}

		// Get the skills mentioned with delete command
		skillsMentioned := strings.TrimPrefix(args.Command, "/"+commandUpdateSkills+" delete ")
		skillsMentionedArray := strings.Split(skillsMentioned, ",")

		// Get skills of the user first
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

			// Check if mentioned skills exist for the user.
			// Keep all skills that exist for the user and not mentioned for deleting
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

			// Delete the user's instance from KV. Will be set again.
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
				// Updated skills are set for the user
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
		p.API.LogInfo("/skills deleted successfully", map[string]string{"Info": "/skills deleted successfully"})
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "If you previously held skills that you mentioned, they have been deleted !",
		}

	default:
		// skills used with inappropriate command
		p.API.LogInfo("/skills not used properly", map[string]string{"Info": "/skills not used properly"})
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Incorrect usage of /skills. Command not recognised: " + command),
		}
	}

}
