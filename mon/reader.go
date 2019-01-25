package mon

import "fmt"
import "strings"
import "errors"
import "bytes"
import "io"

import "golang.org/x/net/html"


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

///

func getFirstElementByName(n *html.Node, name string) (*html.Node, error) {
    if n.Type == html.ElementNode && n.Data == name {
        return n, nil
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        r, err := getFirstElementByName(c, name)
        if err==nil && r != nil {
            return r, err
        }
    }
    return nil, errors.New("Not found")
}

func getLastElementByName(doc *html.Node, name string) (*html.Node, error) {
    var b *html.Node
    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == name {
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
    return nil, errors.New("Missing " + name + " in the node tree")
}

func getPriceTable(n *html.Node) (*html.Node, error) {
    if n.Type == html.ElementNode && n.Data == "table" {
        for _, a := range n.Attr {
            if strings.Compare(a.Key, "data-search-table") == 0 {
                return n, nil
            }
        }
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        ptb, err := getPriceTable(c)
        if err==nil && ptb != nil {
            return ptb, err
        }
    }
    return nil, errors.New("Not found")
}

func getRegionHead(nodePrice *html.Node) (*html.Node, error) {
    return getFirstElementByName(nodePrice, "tr")
}

func readRegionInfo(nodeRegion *html.Node) {
    for c := nodeRegion.FirstChild; c != nil; c = c.NextSibling {
        if c.Data == "th" {
            //var title string
            //var namespace string
            var kv string
            for _, a := range c.Attr {
                if a.Key == "title" {
                    //title = a.Val
                    //namespace = a.Namespace
                    kv = kv + "k:" + a.Key + " v:" + a.Val + "\n"
                    break
                }
            }
            //fmt.Println("title: " + title + " namespace: " + namespace)
            fmt.Println(renderNode(c))
            fmt.Println(kv)

        }
    }
}

func Read(htm string) {
    doc, _ := html.Parse(strings.NewReader(htm))
    nodePrice, err := getPriceTable(doc)
    if err != nil {
        return
    }
    fmt.Println(nodePrice)
    //price := renderNode(nodePrice)
    //fmt.Println(price)
    
    nodeRegion, err := getRegionHead(nodePrice)
    if err != nil {
        return
    }
    //fmt.Println(nodeRegion)
    region := renderNode(nodeRegion)
    fmt.Println(region)
    readRegionInfo(nodeRegion)
}