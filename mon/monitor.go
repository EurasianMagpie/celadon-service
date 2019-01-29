package mon


func RunMonTask(deep bool) {
	d, err := FetchPage()
	if err != nil {
		return
	}
	result, err := Parse(d, deep)
	if err != nil {
		return
	}
	UpdateResult(result)
}