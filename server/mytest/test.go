package test

// import (
// 	// "bytes"
// 	// "encoding/json"
// 	"fmt"
// 	"github.com/laktek/Stack-on-Go/stackongo"
// 	// "io/ioutil"
// 	// "net/http"
// 	"github.com/afjoseph/RAKE.go"

// 	"github.com/abdulsmapara/html2md"
// 	"github.com/charlesvdv/fuzmatch"
// 	// "github.com/deckarep/golang-set"
// 	"regexp"
// 	"strings"
// )

// func preprocessQuery(userQuery string) (processedString string) {
// 	reg, err := regexp.Compile("[^a-zA-Z]+")
// 	if err != nil {
// 		// on failure return original query
// 		return userQuery
// 	}
// 	processedString = reg.ReplaceAllString(userQuery, "")
// 	return processedString
// }
// func processQuery(userQuery string) (keywords string) {
// 	candidates := rake.RunRake(userQuery)
// 	keywords = ""
// 	for _, candidate := range candidates {
// 		keywords += candidate.Key + ";"
// 	}
// 	return keywords

// }
// func searchStackOverflow(userQuery string) map[string]string {

// 	//get keywords from query
// 	userQuery = preprocessQuery(userQuery)
// 	keywords := processQuery(userQuery)
// 	fmt.Printf(keywords + "\n")
// 	arrayOfKeywords := strings.Split(keywords, ";")
// 	wordsOfUserQuery := strings.Split(userQuery, " ")
// 	tags := ""
// 	for _, word := range wordsOfUserQuery {
// 		tags += word + ";"
// 	}
// 	session := stackongo.NewSession("stackoverflow")
// 	//set the params
// 	questions, err := session.Search(userQuery, map[string]string{"tagged": tags, "intitle": arrayOfKeywords[0], "filter": "withbody", "page": "1", "pagesize": "10"})
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	questionsWithAnswers := []stackongo.Question{}
// 	for _, question := range questions.Items {
// 		// fmt.Printf("HERE " +  question.Body + " ")
// 		if question.Is_answered && question.Accepted_answer_id != 0 {
// 			questionsWithAnswers = append(questionsWithAnswers, question)
// 		}
// 	}
// 	maxScore := 0
// 	questionIndexHighestScore := 0
// 	index := 0
// 	for _, question := range questionsWithAnswers {
// 		score := fuzmatch.TokenSetRatio(strings.ToLower(userQuery), strings.ToLower(question.Title))
// 		if score > maxScore {
// 			maxScore = score
// 			questionIndexHighestScore = index
// 		}
// 		index = index + 1
// 	}
// 	ids := []int{questionsWithAnswers[questionIndexHighestScore].Accepted_answer_id}
// 	answers, err := session.GetAnswers(ids, map[string]string{"filter": "withbody"})
// 	answerBody := answers.Items[0].Body

// 	answerMarkdown := html2md.Convert(answerBody)
// 	return map[string]string{
// 		"Found":          "true",
// 		"Question Title": questionsWithAnswers[questionIndexHighestScore].Title,
// 		"Question Body":  questionsWithAnswers[questionIndexHighestScore].Body,
// 		"Solution":       answerMarkdown,
// 	}
// }

// func main() {
// 	// var allSkills2 mapset.Set
// 	// allSkills2 = mapset.NewSet()
// 	// dat, _ := ioutil.ReadFile("../../assets/skills.txt")
// 	// skillsString := strings.ToLower(string(dat))

// 	// for _, skill := range strings.Split(skillsString,"\n"){
// 	// 	allSkills2.Add(strings.ToUpper(skill))
// 	// }

// 	// fmt.Println(allSkills2.Contains("JAVA "))

// 	// skillstr := allSkills2.String()

// 	// fmt.Println(skillstr[4: len(skillstr)-1])
// 	str := "JAVA, PYTHON"
// 	fmt.Println(strings.Split(str,","))
// }
