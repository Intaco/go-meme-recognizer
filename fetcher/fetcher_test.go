package fetcher

func Example_test() {
	f := NewFetcher(10, 10)
	urlsList := [][]string{[]string{"https://pp.userapi.com/c841030/v841030005/1826e/Bunv2Om-uv4.jpg"},
		[]string{"https://pp.userapi.com/c638221/v638221662/81171/lLsKjoP3s_E.jpg",
			"https://pp.userapi.com/c638221/v638221602/55c39/TQnoSaS1eVI.jpg",
			"https://pp.userapi.com/c638221/v638221602/55c41/nOmwuC6O2mA.jpg",
			"https://pp.userapi.com/c638221/v638221602/55c48/g0J54ofxdyY.jpg",
			"https://pp.userapi.com/c638221/v638221602/55c4f/HAipw-io3uY.jpg"},
		[]string{"https://pp.userapi.com/c638221/v638221007/58c3c/9A0Tz4d06bc.jpg"},
		[]string{"https://pp.userapi.com/c638221/v638221007/58c32/nuAr6pMJGhs.jpg"},
		[]string{"https://pp.userapi.com/c638221/v638221007/58c04/KMyDz0wwDIc.jpg"},
		[]string{"https://pp.userapi.com/c638221/v638221388/5f8a1/_y7dUsi15b8.jpg"},
		[]string{"https://pp.userapi.com/c837731/v837731337/55cb3/yAsTav_Ap8A.jpg"},
		[]string{"https://pp.userapi.com/c837731/v837731869/5a3ce/0C9xZypRHRo.jpg"}}
	queries := MakeQueryFromUrlsList("result", urlsList)
	f.Download(queries, 1)
}
