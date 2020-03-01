package main

import (
	"github.com/afjoseph/RAKE.go"
	"github.com/laktek/Stack-on-Go/stackongo"

	"github.com/abdulsmapara/html2md"
	"github.com/charlesvdv/fuzmatch"
	mapset "github.com/deckarep/golang-set"
	"regexp"
	"strings"
)

// preprocessQuery preprocesses the string to keep only the needed information
func preprocessQuery(userQuery string) (preprocessedString string) {
	reg, err := regexp.Compile("[^a-zA-Z+#_]+")
	if err != nil {
		// on failure return original query
		return userQuery
	}
	preprocessedString = reg.ReplaceAllString(userQuery, " ")

	return preprocessedString
}

// processQuery returns the keywords present in the userQuery
// Uses Rapid Automatic Keyword Extraction (RAKE) algorithm as described in:
// Rose, S., Engel, D., Cramer, N., & Cowley, W. (2010). Automatic Keyword Extraction
// from Individual Documents. In M. W. Berry & J. Kogan (Eds.),
// Text Mining: Theory and Applications: John Wiley & Sons.
func processQuery(userQuery string) (keywords string) {
	candidates := rake.RunRake(userQuery)
	keywords = ""
	for _, candidate := range candidates {
		keywords += candidate.Key + ";"
	}
	return keywords

}

// searchStackoverflow searches stackoverflow to find similar question as userQuery
func (p *Plugin) searchStackOverflow(userQuery string) map[string]string {

	// prepare response if search is unsuccessful
	failedResponse := map[string]string{
		"Found":          "false",
		"Question Title": "NOT FOUND",
		"Question Body":  "NOT FOUND",
		"Solution":       "NOT FOUND",
	}

	// Preprocess query first
	userQuery = preprocessQuery(userQuery)

	// Get keywords from query
	keywords := processQuery(userQuery)

	arrayOfKeywords := strings.Split(keywords, ";")
	wordsOfUserQuery := strings.Split(userQuery, " ")

	// Prepare tags from words of userQuer
	tags := ""
	for _, word := range wordsOfUserQuery {
		tags += word + ";"
	}

	// Use stackongo to call stackoverflow API
	session := stackongo.NewSession("stackoverflow")

	// Call stackoverflow search API with appropriate parameters
	questions, err := session.Search(userQuery, map[string]string{"tagged": tags, "intitle": arrayOfKeywords[0], "filter": "withbody", "page": "1", "pagesize": "10", "sort": "votes"})
	if err != nil {
		p.API.LogError("Failed to search on stackoverflow using stackoverflow API. Returned error.", "err", err.Error())
		return failedResponse
	}

	// Consider only those questions that have an accepted answer
	questionsWithAnswers := []stackongo.Question{}
	atleastOne := false

	for _, question := range questions.Items {

		if question.Is_answered && question.Accepted_answer_id != 0 {
			questionsWithAnswers = append(questionsWithAnswers, question)
			atleastOne = true
		}
	}

	if !atleastOne {
		p.API.LogInfo("No question found with an accepted answer", map[string]string{"Info": "No question found with an accepted answer"})
		return failedResponse
	}

	// Among all questions that have an answer, choose a question whose title best matches with user's question/issue
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

	// Get accepted answer to the question selected
	id := []int{questionsWithAnswers[questionIndexHighestScore].Accepted_answer_id}
	answers, err := session.GetAnswers(id, map[string]string{"filter": "withbody"})
	if err != nil {
		p.API.LogError("Failed to get answer from stackoverflow", "err", err.Error())
		return failedResponse
	}
	answerBody := answers.Items[0].Body

	// Rules for formatting html <code> tag and </code> tag
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

	p.API.LogInfo("Search on stackoverflow successful", map[string]string{"Info": "Search on stackoverflow successful"})
	return map[string]string{
		"Found":          "true",
		"Question Title": questionsWithAnswers[questionIndexHighestScore].Title,
		"Question Body":  html2md.Convert(questionsWithAnswers[questionIndexHighestScore].Body),
		"Solution":       html2md.Convert(answerBody),
	}
}

// getSkilledUsers tries to find users who are skilled in the domain of the issue
func (p *Plugin) getSkilledUsers(userQuery string, selfID string) map[string]string {

	preprocessedString := preprocessQuery(strings.TrimSpace(userQuery))

	setOfWords := mapset.NewSet()

	for _, word := range strings.Split(preprocessedString, " ") {
		setOfWords.Add(strings.ToUpper(strings.TrimSpace(word)))
	}

	// Get skills required to solve the issue (userQuery)
	common := setOfWords.Intersect(p.allSkills)
	commonSkills := common.String()
	commonSkills = commonSkills[4 : len(commonSkills)-1]

	// atleast one skill is required in the issue to search for the user
	atleastOneReq := false

	if commonSkills != "" {
		atleastOneReq = true
	}

	skillsForQuery := strings.Split(strings.TrimSpace(commonSkills), ",")

	// Get list of a few users. Currently considering 500.
	users, err := p.API.KVList(0, 500)

	if err != nil {
		p.API.LogError("Failed to fetch user's list", "err", err.Error())
		return map[string]string{
			"Found":   "false",
			"Error":   "true",
			"Message": "Error",
		}
	}

	// Find user who has all the required skills to resolve the issue
	for _, user := range users {
		if user != selfID {
			skills, err := p.API.KVGet(user)
			if err != nil {
				p.API.LogError("Failed to get user's skills", "err", err.Error())
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
						p.API.LogError("Failed to get user information", "err", err.Error())
						return map[string]string{
							"Found":   "false",
							"Error":   "true",
							"Message": "Error",
						}
					}

					// Skilled user (in the domain of the issue) found.
					// Suggest username to the user who wants to resolve the issue.
					p.API.LogInfo("Skilled user found", map[string]string{"Info": "Skilled user found"})
					return map[string]string{
						"Found":   "true",
						"Error":   "false",
						"Message": "You may contact @" + userInfo.Username + " who is skilled in [ " + string(skills) + "]\n",
					}
				}
			}
		}
	}

	// No skilled user found. Report user not found.
	p.API.LogInfo("Skilled user not found", map[string]string{"Info": "Skilled user not found"})
	return map[string]string{
		"Found":   "false",
		"Error":   "false",
		"Message": "User not found ",
	}
}
