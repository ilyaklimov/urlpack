package urlpack

import (
    "fmt"
    // "os"
    "net/url"
    "strings"
    "strconv"
    // "encoding/json"
    "encoding/hex"
    "crypto/sha256"
    "errors"
)

type URL struct {
    hash string
    string string
    scheme string
    authority Authority
    path Path
    query Query
    fragment string
}

type Authority struct {
    userinfo *url.Userinfo
    host Host
}

type Host struct {
    string string
    hostname Hostname
    port int
}

type Hostname struct {
    FQDN string
    Domains []string
    DomainsCount int
    WWW bool
    TLD string
}

type Path struct {
    string string
    Dir []string
    file File
}

type File struct {
    Fullname string
    Name string
    Ext []string
}

type Query struct {
    string string
    list map[string][]string
    Delimiter string
}



func New(urlStr string) (URL, error) {
    var urlPack URL
    var err error
    var u *url.URL

    if u, err = url.Parse(urlStr); err != nil {
        return urlPack, err
    }

    u.Scheme = strings.ToLower(u.Scheme)
    u.Host = strings.ToLower(u.Host)

    if err = urlPack.setStandardFields(u); err != nil {
        return urlPack, err
    }

    if err = urlPack.Host().setFields(); err != nil {
        return urlPack, err
    }

    if err = urlPack.Path().setFields(); err != nil {
        return urlPack, err
    }

    if err = urlPack.Query().setDelimiter(); err != nil {
        return urlPack, err
    }

    urlPack.setHash()

    return urlPack, err

}

func (urlPack *URL) Hash() string {
    return urlPack.hash
}

func (urlPack *URL) String() string {
    return urlPack.string
}

func (urlPack *URL) Userinfo() *url.Userinfo {
    return urlPack.authority.userinfo
}

func (urlPack *URL) Host() *Host {
    return &urlPack.authority.host
}

func (urlPack *URL) Hostname() *Hostname {
    return &urlPack.authority.host.hostname
}

func (urlPack *URL) Path() *Path {
    return &urlPack.path
}

func (urlPack *URL) Query() *Query {
    return &urlPack.query
}

func (q *Query) String() string {
    return q.string
}

// func (urlPack *URL) JSON() ([]byte, error) {

// }


func (urlPack *URL) setStandardFields(u *url.URL) error {
    var err error
    urlPack.string = u.String()
    urlPack.scheme = u.Scheme
    urlPack.authority.userinfo = u.User
    urlPack.authority.host.string = u.Host
    urlPack.path.string = u.Path
    urlPack.query.string = u.RawQuery
    urlPack.fragment = u.Fragment
    urlPack.query.list = u.Query()
    return err
}

func (urlPack *URL) setHash() {
    h := sha256.New()
    b := []byte(urlPack.String())
    h.Write(b)
    urlPack.hash = hex.EncodeToString(h.Sum(nil))
}

func (h *Host) setFields() error {
    var err error
    host := strings.Split(h.string, ":")
    if len(host) == 2 {
        if h.port, err = strconv.Atoi(host[1]); err != nil {
            return err
        }
    } else {
        h.port = 80
    }
    h.hostname.FQDN = host[0]
    h.hostname.Domains = strings.Split(h.hostname.FQDN, ".")
    h.hostname.DomainsCount = len(h.hostname.Domains)
    h.hostname.TLD = h.hostname.Domains[h.hostname.DomainsCount-1]
    if h.hostname.Domains[0] == "www" {
        h.hostname.WWW = true
    }
    return err
}

func (p *Path) setFields() error {
    var err error
    pathStr := strings.Trim(p.string, "/")
    if p.string == "" {
        return err
    }
    path := strings.Split(pathStr, "/")
    pathLastKey := len(path) - 1
    if strings.ContainsAny(path[pathLastKey], ".") {
        p.file.Fullname = path[pathLastKey]
        file := strings.Split(p.file.Fullname, ".")
        if len(file) < 2 {
            return errors.New("Bad filename in URL path!")
        }
        p.file.Name = file[0]
        p.file.Ext = file[1:]
        p.Dir = path[:pathLastKey]
    } else {
        p.Dir = path
    }

    fmt.Println(p.Dir)
    return err
}

func (q *Query) setDelimiter() error {
    var err error
    if q.string == "" {
        return err
    }
    switch {
        case strings.ContainsAny(q.string, "&"):
            q.Delimiter = "&"
        case strings.ContainsAny(q.string, ";"):
            q.Delimiter = ";"
        default:
            return errors.New("Bad delimiter in URL query!")
    }
    return err
}








// type URL struct {
//         Scheme     string
//         Opaque     string    // encoded opaque data
//         User       *Userinfo // username and password information
//         Host       string    // host or host:port
//         Path       string
//         RawPath    string // encoded path hint (Go 1.5 and later only; see EscapedPath method)
//         ForceQuery bool   // append a query ('?') even if RawQuery is empty
//         RawQuery   string // encoded query values, without '?'
//         Fragment   string // fragment for references, without '#'
// }

// func main() {
//  var urlStr string = "http://user:pass@www.Example.com:80/dir1/dir2/index.php?q1=hello&q2=world"

//  fmt.Println(urlStr)

//  u, err := url.Parse(urlStr)
//  if err != nil {
//      fmt.Println(err)
//  }

//  fmt.Println("Scheme: ", u.Scheme)
//  fmt.Println("Opaque: ", u.Opaque)
//  fmt.Println("User: ", u.User)
//  fmt.Println("Host: ", u.Host)
//  fmt.Println("Path: ", u.Path)
//  fmt.Println("RawQuery: ", u.RawQuery)
//  fmt.Println("Fragment: ", u.Fragment)
//  fmt.Println("User Password: ", u.User.Username())

//  pass, b :=  u.User.Password()

//  fmt.Println("User Password: ", pass)
//  fmt.Println("User Password (b): ", b)

//  fmt.Println("IsAbs: ", u.IsAbs())

//  fmt.Println("u: ", u)

//  fmt.Println("u query: ", u.Query())

//  u.Scheme = strings.ToLower(u.Scheme)
//  u.Host = strings.ToLower(u.Host)

//  fmt.Println("u: ", u)

//  urlStr2 := u.String()

//  fmt.Println(urlStr2)

//  var us URL
//  fmt.Println(us.string)




// }

/*
    domainStr := "www.api.ya.example.com"
    domain := strings.Split(domainStr, ".")
    fmt.Println(domain)

    domainCount := len(domain)

    fmt.Println(domainCount)

    tld := domain[domainCount-1]

    fmt.Println(tld)



*/