package mon

import "fmt"
import "strings"
import "errors"
import "bytes"
import "io"

import "golang.org/x/net/html"


func fetchPage() {
	fmt.Println("fetch page")
}

const htm = `<!DOCTYPE html>
<html>
<head>
    <title></title>
</head>
<body>
    body content
    <p>more content</p>
</body>
</html>`

func getBody(doc *html.Node) (*html.Node, error) {
    var b *html.Node
    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "body" {
            b = n
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
    f(doc)
    if b != nil {
        return b, nil
    }
    return nil, errors.New("Missing <body> in the node tree")
}

func renderNode(n *html.Node) string {
    var buf bytes.Buffer
    w := io.Writer(&buf)
    html.Render(w, n)
    return buf.String()
}

func LoadHtml() {
	doc, _ := html.Parse(strings.NewReader(htm))
	bn, err := getBody(doc)
	if err != nil {
		return
	}
	body := renderNode(bn)
	fmt.Println(body)
}