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

func ParseRegion(nodePrice *html.Node) {
    // the region header
    nodeRegion, err := getFirstElementByName(nodePrice, "tr")
    if err != nil {
        return
    }
    //region := renderNode(nodeRegion)
    //fmt.Println(region)
    for c := nodeRegion.FirstChild; c != nil; c = c.NextSibling {
        if c.Data == "th" {
            var title string
            for _, a := range c.Attr {
                if a.Key == "title" {
                    title = a.Val
                    break
                }
            }
            if len(title) == 0 {
                continue
            }
            rgn := renderNode(c)
            a := strings.LastIndex(rgn, "/>")
            b := strings.LastIndex(rgn, "<")
            abbr := rgn[a+2:b]
            abbr = strings.Trim(abbr, " ")
            fmt.Println(abbr, title)
        }
    }
}

func ParseGamePrice(nodePrice *html.Node) {
    tbody, err := getFirstElementByName(nodePrice, "tbody")
    if err != nil {
        return
    }
    for row:=tbody.FirstChild; row!=nil; row=row.NextSibling {
        var name string
        var price string
        var lrgn, hrgn int
        var lp, hp string
        i := 0
        for c:=row.FirstChild; c!=nil; c=c.NextSibling {
            if c.Data == "th" {
                gname := renderNode(c)
                a := strings.LastIndex(gname, "\">")
                b := strings.LastIndex(gname, "</a")
                name = gname[a+2:b]
            } else if c.Data == "td" {
                p := renderNode(c)
                a := strings.Index(p, ">")
                b := strings.LastIndex(p, "</")
                np := p[a+1:b]
                np = strings.Trim(np, "¥")
                if len(price) == 0 {
                    price = np
                } else {
                    price = price + "," + np
                }

                for _, a := range c.Attr {
                    if a.Key == "class" {
                        cls := a.Val
                        if cls == "l" {
                            lrgn = i
                            lp = np
                        } else if cls == "h" {
                            hrgn = i
                            hp = np
                        }
                        break
                    }
                }
                i++
            }
        }
        if len(name) > 0 && len(price) > 0 {
            fmt.Println(name, price)
            fmt.Println("lrgn:", lrgn, " lprice:", lp, " hrgn:", hrgn, " hprice:", hp)
        }
    }
}

func Parse(htm string) {
    doc, _ := html.Parse(strings.NewReader(htm))
    nodePrice, err := getPriceTable(doc)
    if err != nil {
        return
    }
    
    ParseRegion(nodePrice)

    ParseGamePrice(nodePrice)
}