package mon


func RunMonTask() {
	d, err := FetchPage()
	if err != nil {
		return
	}
	Parse(d)
}