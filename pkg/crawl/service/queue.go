package service

type Queue struct {
	urls []string
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Push(url string) {
	q.urls = append(q.urls, url)
}

func (q *Queue) Pop() string {
	if len(q.urls) == 0 {
		return ""
	}

	url := q.urls[0]
	q.urls = q.urls[1:]

	return url
}
