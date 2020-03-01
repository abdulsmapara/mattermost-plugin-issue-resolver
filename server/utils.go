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
func processQuery(userQuery string) (keywords string) {
	candidates := rake.RunRake(userQuery)
	keywords = ""
	for _, candidate := range candidates {
		keywords += candidate.Key + ";"
	}
	return keywords

}

// searchStackoverflow searches stackoverflow to find similar question as userQuery
func searchStackOverflow(userQuery string) map[string]string {

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

// getSkilledUsers tries to find users who are skilled in the domain of the issue
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
		"Message": "User not found ",
	}
}
