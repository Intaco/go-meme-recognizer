package get_mempedia_url

import "testing"

func TestGet_mempedia_url(t *testing.T) {
	link, name, _ := Get_mempedia_url("С блэкджеком и шлюхами")
	print(link, " ", name, "\n")
	link, name, _ = Get_mempedia_url("кот качок")
	print(link, " ", name, "\n")
	link, name, _ = Get_mempedia_url("нельзя просто так взять и")
	print(link, " ", name, "\n")
	link, name, _ = Get_mempedia_url("давай досвидания")
	print(link, " ", name, "\n")
	link, name, _ = Get_mempedia_url("бессмысленный запрос")
	print(link, " ", name, "\n")
}