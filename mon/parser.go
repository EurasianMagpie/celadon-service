package mon

import "fmt"
import "strings"
import "errors"
import "bytes"
import "io"
import "regexp"
import "time"

import "golang.org/x/net/html"

import "github.com/EurasianMagpie/celadon-service/db"
import "github.com/EurasianMagpie/celadon-service/util"


var defaultDate time.Time

func init() {
    dt, _ := time.Parse("2006-01-02", "2018-01-01")
    defaultDate = dt
}

type ParseResult struct {
    Regions []db.Region
    Games []db.GameInfo
    Prices []db.Price
}
var parseResult ParseResult

func renderNode(n *html.Node) string {
    var buf bytes.Buffer
    w := io.Writer(&buf)
    html.Render(w, n)
    return buf.String()
}

func simpleNodeContent(n *html.Node, nodeName string) string {
    s := renderNode(n)
    tag := "<"+nodeName+">"
    a := strings.Index(s, tag)
    b := strings.Index(s, "</"+nodeName+">")
    r := ""
    if a == -1 {
        tag1 := ">"
        a1 := strings.Index(s, tag1)
        r = s[a1+1:b]
    } else {
        r = s[a+len(tag):b]
    }
    if len(r) > 0 {
        r = strings.Trim(r, " \n")
    }
    return r
}

func getNodeAttr(n *html.Node, key string) (string, error) {
    for _, a := range n.Attr {
        if a.Key == key {
            return a.Val, nil
        }
    }
    return "", errors.New("Not found")
}

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

func getFirstElementByNameAndAttr(n *html.Node, name string, attrKey string, attrVal string) (*html.Node, error) {
    if n.Type == html.ElementNode && n.Data == name {
        for _, a := range n.Attr {
            if a.Key == attrKey {
                if a.Val == attrVal {
                    return n, nil
                }
            }
        }
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        r, err := getFirstElementByNameAndAttr(c, name, attrKey, attrVal)
        if err==nil && r != nil {
            return r, err
        }
    }
    return nil, errors.New("Not found")
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
            var title string = ""
            img, err := getFirstElementByName(c, "img")
            if err == nil {
                for _, a := range img.Attr {
                    if a.Key == "title" {
                        title = a.Val
                        break
                    }
                }
            }
            
            if len(title) == 0 {
                continue
            }
            rgn := renderNode(c)
            a := strings.LastIndex(rgn, "/>")
            b := strings.LastIndex(rgn, "<")
            abbr := rgn[a+2:b]
            abbr = strings.Trim(abbr, " \n")
            //fmt.Println(abbr, title)
            parseResult.Regions = append(parseResult.Regions, db.NewRegion(abbr, title))
        }
    }
}

func ParseGamePrice(nodePrice *html.Node) {
    tbody, err := getFirstElementByName(nodePrice, "tbody")
    if err != nil {
        return
    }
    for row:=tbody.FirstChild; row!=nil; row=row.NextSibling {
        var id string
        var name string
        var ref string
        var price string
        var lrgn, hrgn int
        var lp, hp string
        i := 0
        for c:=row.FirstChild; c!=nil; c=c.NextSibling {
            if c.Data == "th" {
                na, err := getFirstElementByName(c, "a")
                if err == nil {
                    for _, attr := range na.Attr {
                        if attr.Key == "href" {
                            link := attr.Val
                            ref = link
                            a := strings.LastIndex(link, "/")
                            l := link[a+1:]
                            b := strings.Index(l, "-")
                            id = l[:b]
                            break
                        }
                    }
                }

                gname := renderNode(c)
                a := strings.Index(gname, "\">")
                b := strings.LastIndex(gname, "</a")
                gname1 := gname[a+2:b]
                span := strings.Index(gname1, "<span")
                if span == -1 {
                    name = strings.Trim(gname1, " \n")
                } else {
                    sp1 := strings.Index(gname1, ">")
                    sp2 := strings.Index(gname1, "</span>")
                    n1 := gname1[sp1+1:sp2]
                    n2 := gname1[sp2+7:]
                    name = strings.Trim(n1 + " " + n2, " \n")
                    name = strings.Trim(n1, " \n") + " " + strings.Trim(n2, " \n")
                }
            } else if c.Data == "td" {
                p := renderNode(c)
                a := strings.Index(p, ">")
                b := strings.LastIndex(p, "</")
                np := p[a+1:b]
                np = strings.Trim(np, "¥")
                if i == 0 {
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
            b, r := util.UnEscape(name)
            if b {
                name = r
            }
            //fmt.Println(id, name, price)
            //fmt.Println("lrgn:", regions[lrgn].Name, " lprice:", lp, " hrgn:", regions[hrgn].Name, " hprice:", hp)
            parseResult.Games = append(parseResult.Games, db.NewGameInfo(id, name, ref))
            parseResult.Prices = append(parseResult.Prices, db.NewPrice(id, price, lp, parseResult.Regions[lrgn].Region_id, hp, parseResult.Regions[hrgn].Region_id))
        }
    }
}

func parseDate(s string) time.Time {
    //January 1st, 2019
    r, err := regexp.Compile(`([a-zA-Z]+) ([0-9]+)[a-z]*, ([0-9]{4})`)
    if err != nil {
        return defaultDate
    }
    params := r.FindStringSubmatch(s)
    t := fmt.Sprintf("%s %s %s", params[1], params[2], params[3])
    layout := "January 2 2006"
    dt, err := time.Parse(layout, t)
    if err != nil {
        return defaultDate
    }
    return dt
}

func DeepParseSingleGame(g *db.GameInfo) bool {
    htm, err := FetchHtmlFromUrl(g.Ref)
    if err != nil {
        return false
    }
    doc, err := html.Parse(strings.NewReader(htm))
    if err != nil {
        return false
    }
    body, err := getFirstElementByName(doc, "body")
    if err != nil {
        return false
    }

    div, err := getFirstElementByNameAndAttr(body, "div", "class", "hero game-hero")
    if err != nil {
        return false
    }
    div, err = getFirstElementByNameAndAttr(div, "div", "class", "wrapper")
    if err != nil {
        return false
    }

    var title string
    var desc string
    var date string
    var img string
    var imgType string
    for c:=div.FirstChild; c!=nil; c=c.NextSibling {
        if c.Data == "picture" {
            src, err := getFirstElementByName(c, "source")
            if err != nil {
                continue
            }
            srcset, err := getNodeAttr(src, "srcset")
            if err != nil {
                continue
            }
            r := strings.Split(srcset, " ")
            img = r[0]
            _t, err := getNodeAttr(src, "type")
            s := strings.LastIndex(_t, "/")
            _t = _t[s+1:]
            imgType = _t
        } else if c.Data == "div" {
            for d:=c.FirstChild; d!=nil; d=d.NextSibling {
                if d.Data == "h1" {
                    title = simpleNodeContent(d, "h1")
                } else if d.Data == "p" {
                    desc = simpleNodeContent(d, "p")
                } else if d.Data == "small" {
                    date = simpleNodeContent(d, "small")
                    //len("Released on ") = 12
                    date = date[12:]
                }
            }
        }
    }

    b, r := util.UnEscape(title)
    if b {
        title = r
    }
    b, r = util.UnEscape(desc)
    if b {
        desc = r
    }

    g.Name = title
    g.Desc = desc
    g.ReleaseDate = parseDate(date)
    g.CoverUrl = img
    g.CoverType = imgType
    fmt.Println(g)
    //fmt.Println(g.Name, g.ReleaseDate, date)
    return true
}

func DeepParseGameInfo() {
    if !db.ReCheckGameDetail() {
        return
    }

    for i:=0; i<len(parseResult.Games); i++ {
        g := &parseResult.Games[i]
        /*if i > 5 {
            break
        }//*/
        if db.IsGameDetialed(g.Id) {
            continue 
        }

        DeepParseSingleGame(g)
    }
}

func Parse(htm string, deep bool) (*ParseResult, error) {
    parseResult.Regions = parseResult.Regions[:0]
    parseResult.Games = parseResult.Games[:0]
    parseResult.Prices = parseResult.Prices[:0]

    doc, _ := html.Parse(strings.NewReader(htm))
    nodePrice, err := getPriceTable(doc)
    if err != nil {
        return nil, errors.New("Parse failed")
    }
    
    ParseRegion(nodePrice)

    ParseGamePrice(nodePrice)

    if deep {
        DeepParseGameInfo()
    }
    return &parseResult, nil
}