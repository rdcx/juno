package crawler

import (
	"testing"

	"github.com/h2non/gock"
)

const clothesPage = `<html>
	<head>
		<title>Shop</title>
	</head>
	<body>
		<h1>Clothes</h1>
	</body>
</html>`

func TestFetchPage(t *testing.T) {
	defer gock.Off()

	gock.New("https://shop.com").
		Get("/clothes").
		Reply(200).
		BodyString(clothesPage)

	page, err := FetchPage("https://shop.com/clothes")

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if string(page) != clothesPage {
		t.Errorf("Expected page to be %s, got %s", clothesPage, page)
	}
}
