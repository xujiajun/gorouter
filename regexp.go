package gorouter

import (
	"strings"
	"regexp"
)

var (
	defaultPattern = `[\w]+`
	idPattern      = `[\d]+`
)

func Match(requestUrl string, path string) (bool, map[string]string) {
	res := strings.Split(path, "/")
	if res == nil {
		return false, nil
	}

	var (
		matchName   []string
		matchParams map[string]string
		sTemp       string
	)

	matchParams = make(map[string]string)

	for _, str := range res {

		if str != "" {
			r := []byte(str)

			if string(r[0]) == "{" && string(r[len(r)-1]) == "}" {
				matchStr := string(r[1:len(r)-1])
				res := strings.Split(matchStr, ":")

				matchName = append(matchName, res[0])

				sTemp = sTemp + "/" + "(" + res[1] + ")"
			} else if string(r[0]) == ":" {
				matchStr := string(r)
				res := strings.Split(matchStr, ":")
				matchName = append(matchName, res[1])

				if res[1] == "id" {
					sTemp = sTemp + "/" + "(" + idPattern + ")"
				} else {
					sTemp = sTemp + "/" + "(" + defaultPattern + ")"
				}
			} else {
				sTemp = sTemp + "/" + str
			}
		}
	}

	pattern := sTemp

	re := regexp.MustCompile(pattern)
	submatch := re.FindSubmatch([]byte(requestUrl))

	if submatch != nil {
		if string(submatch[0]) == requestUrl {
			submatch = submatch[1:]
			for k, v := range submatch {
				matchParams[matchName[k]] = string(v)
			}
			return true, matchParams
		}
	}

	return false, nil
}
