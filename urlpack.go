package urlpack

import (
    "net/url"
    "strings"
    "strconv"
    "encoding/hex"
    "crypto/sha256"
    "errors"
)

type URL struct {
    Hash        string      `json:"hash"`
    String      string      `json:"string"`
    Scheme      string      `json:"scheme"`
    Authority   Authority   `json:"authority"`
    Path        Path        `json:"path"`
    Query       Query       `json:"query"`
    Fragment    string      `json:"fragment"`
}

type Authority struct {
    Userinfo    Userinfo    `json:"userinfo"`
    Host        Host        `json:"host"`
}

type Userinfo struct {
    Username    string      `json:"username"`
    Password    string      `json:"password"`
    String      string      `json:"string"`
}

type Host struct {
    String      string      `json:"string"`
    Hostname    Hostname    `json:"hostname"`
    Port        int         `json:"port"`
}

type Hostname struct {
    FQDN            string      `json:"fqdn"`
    Domains         []string    `json:"domains"`
    DomainsCount    int         `json:"domainsCount"`
    WWW             bool        `json:"www"`
    TLD             string      `json:"tld"`
}

type Path struct {
    String  string      `json:"string"`
    Dir     []string    `json:"dir"`
    File    File        `json:"file"`
}

type File struct {
    Fullname    string      `json:"fullname"`
    Name        string      `json:"name"`
    Ext         []string    `json:"ext"`
}

type Query struct {
    String      string                  `json:"string"`
    List        map[string][]string     `json:"list"`
    Delimiter   string                  `json:"delimiter"`
}



func New(urlStr string) (URL, error) {
    var urlPack URL
    var err error
    var u *url.URL

    if u, err = url.Parse(urlStr); err != nil {
        return urlPack, err
    }

    if u.Path == "" {
        u.Path = "/"
    }

    u.Scheme = strings.ToLower(u.Scheme)
    u.Host = strings.ToLower(u.Host)

    if err = urlPack.setStandardFields(u); err != nil {
        return urlPack, err
    }

    if err = urlPack.Host().setFields(); err != nil {
        return urlPack, err
    }

    if err = urlPack.Path.setFields(); err != nil {
        return urlPack, err
    }

    if err = urlPack.Query.setDelimiter(); err != nil {
        return urlPack, err
    }

    urlPack.setHash()

    return urlPack, err

}

func (urlPack *URL) Userinfo() *Userinfo {
    return &urlPack.Authority.Userinfo
}

func (urlPack *URL) Host() *Host {
    return &urlPack.Authority.Host
}

func (urlPack *URL) Hostname() *Hostname {
    return &urlPack.Authority.Host.Hostname
}

func (urlPack *URL) File() *File {
    return &urlPack.Path.File
}

func (urlPack *URL) setStandardFields(u *url.URL) error {
    var err error
    urlPack.String = u.String()
    urlPack.Scheme = u.Scheme
    if u.User != nil {
        urlPack.Authority.Userinfo.String = u.User.String()
        urlPack.Authority.Userinfo.Username = u.User.Username()
        urlPack.Authority.Userinfo.Password, _ = u.User.Password()
    }
    urlPack.Authority.Host.String = u.Host
    urlPack.Path.String = u.Path
    urlPack.Query.String = u.RawQuery
    urlPack.Fragment = u.Fragment
    urlPack.Query.List = u.Query()
    return err
}

func (urlPack *URL) setHash() {
    h := sha256.New()
    b := []byte(urlPack.String)
    h.Write(b)
    urlPack.Hash = hex.EncodeToString(h.Sum(nil))
}

func (h *Host) setFields() error {
    var err error
    host := strings.Split(h.String, ":")
    if len(host) == 2 {
        if h.Port, err = strconv.Atoi(host[1]); err != nil {
            return err
        }
    } else {
        h.Port = 80
    }
    h.Hostname.FQDN = host[0]
    h.Hostname.Domains = strings.Split(h.Hostname.FQDN, ".")
    h.Hostname.DomainsCount = len(h.Hostname.Domains)
    h.Hostname.TLD = h.Hostname.Domains[h.Hostname.DomainsCount-1]
    if h.Hostname.Domains[0] == "www" {
        h.Hostname.WWW = true
    }
    return err
}

func (p *Path) setFields() error {
    var err error
    pathStr := strings.Trim(p.String, "/")
    if pathStr == "" {
        return err
    }
    pathStr = strings.Replace(pathStr, "//", "/", -1)
    path := strings.Split(pathStr, "/")
    pathLastKey := len(path) - 1
    if strings.ContainsAny(path[pathLastKey], ".") {
        p.File.Fullname = path[pathLastKey]
        file := strings.Split(p.File.Fullname, ".")
        if len(file) < 2 {
            return errors.New("Bad filename in URL path!")
        }
        p.File.Name = file[0]
        p.File.Ext = file[1:]
        p.Dir = path[:pathLastKey]
    } else {
        p.Dir = path
    }
    return err
}

func (q *Query) setDelimiter() error {
    var err error
    if q.String == "" {
        return err
    }
    switch {
        case strings.ContainsAny(q.String, "&"):
            q.Delimiter = "&"
        case strings.ContainsAny(q.String, ";"):
            q.Delimiter = ";"
    }
    return err
}