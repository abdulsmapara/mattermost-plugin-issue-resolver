package main

import (
	// "bytes"
	// "encoding/json"
	"fmt"
	"github.com/laktek/Stack-on-Go/stackongo"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	// "net/http"
	"github.com/afjoseph/RAKE.go"

	"github.com/abdulsmapara/html2md"
	"github.com/charlesvdv/fuzmatch"
	mapset "github.com/deckarep/golang-set"
	"regexp"
	"strings"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// UserID of issueResolver search bot
	issueResolverBotID string
	// set of all skills
	allSkills mapset.Set
}

const (
	commandQuery        = "resolve"
	commandUpdateSkills = "skills"
)

// ExecuteCommand TODO
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	trigger := strings.TrimPrefix(strings.Fields(args.Command)[0], "/")
	switch trigger {
	case commandUpdateSkills:
		return p.executeCommandSkills(args), nil
	case commandQuery:
		return p.searchStackOverflow(args), nil
	default:
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Unknown Command: " + args.Command),
		}, nil
	}

}

func preprocessQuery(userQuery string) (preprocessedString string) {
	reg, err := regexp.Compile("[^a-zA-Z+]+")
	if err != nil {
		// on failure return original query
		return userQuery
	}
	preprocessedString = reg.ReplaceAllString(userQuery, " ")
	return preprocessedString
}
func processQuery(userQuery string) (keywords string) {
	candidates := rake.RunRake(userQuery)
	keywords = ""
	for _, candidate := range candidates {
		keywords += candidate.Key + ";"
	}
	return keywords

}
func searchStackOverflow(userQuery string) map[string]string {

	failedResponse := map[string]string{
		"Found":          "false",
		"Question Title": "NOT FOUND",
		"Question Body":  "NOT FOUND",
		"Solution":       "NOT FOUND",
	}
	//get keywords from query
	userQuery = preprocessQuery(userQuery)
	keywords := processQuery(userQuery)
	fmt.Printf(keywords + "\n")
	arrayOfKeywords := strings.Split(keywords, ";")
	wordsOfUserQuery := strings.Split(userQuery, " ")
	tags := ""
	for _, word := range wordsOfUserQuery {
		tags += word + ";"
	}
	session := stackongo.NewSession("stackoverflow")
	//set the params
	questions, err := session.Search(userQuery, map[string]string{"tagged": tags, "intitle": arrayOfKeywords[0], "filter": "withbody", "page": "1", "pagesize": "10"})
	if err != nil {
		return failedResponse
	}

	questionsWithAnswers := []stackongo.Question{}
	atleastOne := false
	for _, question := range questions.Items {
		// fmt.Printf("HERE " +  question.Body + " ")
		if question.Is_answered && question.Accepted_answer_id != 0 {
			questionsWithAnswers = append(questionsWithAnswers, question)
			atleastOne = true
		}
	}
	if !atleastOne {
		return failedResponse
	}
	maxScore := 0
	questionIndexHighestScore := 0
	index := 0
	for _, question := range questionsWithAnswers {
		score := fuzmatch.TokenSetRatio(strings.ToLower(userQuery), strings.ToLower(question.Title))
		if score > maxScore {
			maxScore = score
			questionIndexHighestScore = index
		}
		index = index + 1
	}

	id := []int{questionsWithAnswers[questionIndexHighestScore].Accepted_answer_id}
	answers, err := session.GetAnswers(id, map[string]string{"filter": "withbody"})
	if err != nil {
		return failedResponse
	}
	answerBody := answers.Items[0].Body

	html2md.AddRule("", &html2md.Rule{
		Patterns: []string{"<code>"},
		Tp:       html2md.Void,
		Replacement: func(innerHTML string, attrs []string) string {
			return "```"
		},
	})
	html2md.AddRule("", &html2md.Rule{
		Patterns: []string{"</code>"},
		Tp:       html2md.Void,
		Replacement: func(innerHTML string, attrs []string) string {
			return "```"
		},
	})

	return map[string]string{
		"Found":          "true",
		"Question Title": questionsWithAnswers[questionIndexHighestScore].Title,
		"Question Body":  html2md.Convert(questionsWithAnswers[questionIndexHighestScore].Body),
		"Solution":       html2md.Convert(answerBody),
	}
}

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

// getSkilledUsers TODO
func (p *Plugin) getSkilledUsers(userQuery string, selfID string) map[string]string {
	preprocessedString := preprocessQuery(strings.TrimSpace(userQuery))
	setOfWords := mapset.NewSet()
	for _, word := range strings.Split(preprocessedString, " ") {
		setOfWords.Add(strings.ToUpper(strings.TrimSpace(word)))
	}
	common := setOfWords.Intersect(p.allSkills)
	commonSkills := common.String()
	commonSkills = commonSkills[4 : len(commonSkills)-1]
	atleastOneReq := false
	if commonSkills != "" {
		atleastOneReq = true
	}
	skillsForQuery := strings.Split(strings.TrimSpace(commonSkills), ",")

	users, err := p.API.KVList(0, 500)
	if err != nil {
		return map[string]string{
			"Found":   "false",
			"Error":   "true",
			"Message": "Error",
		}
	}
	for _, user := range users {
		if user != selfID {
			skills, err := p.API.KVGet(user)
			if err != nil {
				return map[string]string{
					"Found":   "false",
					"Error":   "true",
					"Message": "Error",
				}
			}
			if skills != nil {
				skillsOfUser := strings.Split(string(skills), ",")
				foundAll := true

				for _, skillReq := range skillsForQuery {

					foundThisSkill := false
					for _, skill := range skillsOfUser {
						if strings.ToUpper(strings.TrimSpace(skillReq)) == skill {
							foundThisSkill = true
							break
						}
					}
					if !foundThisSkill {
						foundAll = false
						break
					}
				}
				if foundAll && atleastOneReq {
					userInfo, err := p.API.GetUser(user)
					if err != nil {
						return map[string]string{
							"Found":   "false",
							"Error":   "true",
							"Message": "Error",
						}
					}
					return map[string]string{
						"Found":   "true",
						"Error":   "false",
						"Message": "You may contact @" + userInfo.Username + " who is skilled in [ " + string(skills) + "]\n",
					}
				}
			}
		}
	}

	return map[string]string{
		"Found":   "false",
		"Error":   "false",
		"Message": "User not found " + strings.TrimSpace(commonSkills),
	}
}

// searchStackOverflow TODO
func (p *Plugin) searchStackOverflow(args *model.CommandArgs) *model.CommandResponse {
	userQuery := strings.TrimPrefix(args.Command, "/"+commandQuery)
	response := searchStackOverflow(userQuery)
	msg := "###### QUERY: \n" + userQuery + "\n"
	if response["Found"] == "true" {
		msg += "###### QUESTION FOUND:\n" + response["Question Title"] + "\n" + response["Question Body"] + "\n" + "###### Solution:\n" + response["Solution"]
	} else {
		skilledUserInfo := p.getSkilledUsers(userQuery, args.UserId)
		if skilledUserInfo["Found"] == "true" {
			msg += "###### SUGGESTION: \n" + skilledUserInfo["Message"]
		} else {
			fmt.Printf(skilledUserInfo["Message"])
			msg += "Sorry, the issue could not be resolved.\nTry asking me by including specific keywords. "
		}
	}
	post := &model.Post{
		ChannelId: args.ChannelId,
		RootId:    args.RootId,
		UserId:    p.issueResolverBotID,
		Message:   msg,
	}
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

// OnActivate TODO
func (p *Plugin) OnActivate() error {
	p.allSkills = mapset.NewSet()
	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          commandQuery,
		AutoComplete:     true,
		AutoCompleteHint: "[ Your Question/Issue to be searched on STACKOVERFLOW ]",
		AutoCompleteDesc: "Searches stackoverflow for similar/related question. Suggests users to contact for resolving the issue.",
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandQuery)
	}

	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          commandUpdateSkills,
		AutoComplete:     true,
		AutoCompleteHint: "[ Empty OR Commands ]",
		AutoCompleteDesc: "Empty returns your current skills, Command 'list' returns a list of all skills with given prefix , Command 'add' adds skills against your name, Command 'delete' deletes the skill against your name if it existed.",
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandQuery)
	}

	issueResolverBot := &model.Bot{
		Username:    "issue_resolver",
		DisplayName: "ISSUE RESOLVER",
		Description: "Created by Mattermost Issue Resolver Plugin. Searches stackoverflow for questions posted by Mattermost users using /resolve. Also suggests users to contact for resolving the issue.",
	}

	// if not present create a stackoverflow search bot
	issueResolverBotID, err := p.Helpers.EnsureBot(issueResolverBot)
	if err != nil {
		return errors.Wrap(err, "Failed to ensure stackoverflow search bot")
	}

	p.issueResolverBotID = issueResolverBotID

	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		return errors.Wrap(err, "Could not get bundle path")
	}
	skillsData, err := ioutil.ReadFile(filepath.Join(bundlePath, "assets", "skills.txt"))
	skillsString := strings.ToUpper(string(skillsData))
	for _, skill := range strings.Split(skillsString, "\n") {
		p.allSkills.Add(skill)
	}

	return nil
}

func main() {
	plugin.ClientMain(&Plugin{})
}
