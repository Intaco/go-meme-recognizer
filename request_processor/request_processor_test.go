package request_processor

import "testing"

func TestIsUrl(t *testing.T) {
	print(is_link("http://windows-school.ru/images/top-lg.png") == true)
	print(is_link("https://pp.userapi.com/c837123/v837123072/56fc6/cxY7uc-kZB8.jpg") == true)
	print(is_link("https://open-file.ru/types/pictures/") == false)
	print(is_link("имя мема") == false)
}
