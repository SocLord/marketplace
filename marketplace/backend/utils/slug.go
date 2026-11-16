package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var translitMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "shch",
	'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",

	'ғ': "gh", 'ӣ': "i", 'ӯ': "u", 'ҳ': "h", 'ҷ': "j", 'қ': "q",
}

func Slugify(input string) string {
	lower := strings.ToLower(input)
	var sb strings.Builder
	for _, r := range lower {
		if val, ok := translitMap[r]; ok {
			sb.WriteString(val)
		} else if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
		} else {
			sb.WriteRune('-')
		}
	}
	re := regexp.MustCompile("-+")
	clean := re.ReplaceAllString(sb.String(), "-")
	clean = strings.Trim(clean, "-")
	if clean == "" {
		clean = fmt.Sprintf("item-%d", time.Now().Unix())
	}
	return clean
}
