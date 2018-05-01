package goinsta

import (
	"github.com/erikdubbelboer/fasthttp"
)

// ClientInfo ...
type ClientInfo struct {
	Username string
	// TODO: Is safe to store password in memory?
	Password string
	// TODO: Allow user change this fields?
	DeviceID  string
	UUID      string
	RankToken string
	Token     string
	PhoneID   string
}

type cookies map[string]*fasthttp.Cookie

func (ck *cookies) Set(key, value []byte) {
	ks := b2s(key)
	c, ok := (*ck)[ks]
	if !ok {
		c = fasthttp.AcquireCookie()
	}
	c.SetKeyBytes(key)
	c.SetValueBytes(value)
	(*ck)[ks] = c
}

func (ck *cookies) SetCookies(cks []*fasthttp.Cookie) {
	for _, c := range cks {
		(*ck)[b2s(c.Key())] = c
	}
}

func (ck *cookies) Release() {
	for k, c := range *ck {
		fasthttp.ReleaseCookie(c)
		delete(*ck, k)
	}
}

func (ck *cookies) Peek(v string) string {
	c := (*ck)[v]
	if c != nil {
		return b2s(c.Value())
	}
	return ""
}

func (ck *cookies) Cookies() map[string]*fasthttp.Cookie {
	return *ck
}

// Instagram ....
type Instagram struct {
	Logged bool
	Info   ClientInfo

	// DialFunc allows user to use proxy function.
	// See also: https://godoc.org/github.com/erikdubbelboer/fasthttp#Client.Dial
	DialFunc fasthttp.DialFunc

	client  *fasthttp.Client
	cookies *cookies

	// Account stores logged in user data and interactions
	Account *Accout `json:"user,logged_in_user"`

	// Instagram objects
	User    *User
	Media   *Media
	Search  *Search
	Explore *Explore
	Inbox   *Inbox

	StatusResponse
}

// NewViaProxy All requests will use proxy server (example http://<ip>:<port>)
func NewViaProxy(username, password string, dialFunc fasthttp.DialFunc) *Instagram {
	insta := New(username, password)
	insta.DialFunc = dialFunc
	return insta
}

// New creates instagram structure
func New(username, password string) *Instagram {
	insta := &Instagram{
		client: &fasthttp.Client{
			Name: goInstaUserAgent,
		},
		cookies: nil,
		Info: &ClientInfo{
			DeviceID: generateDeviceID(generateMD5Hash(username + password)),
			Username: username,
			Password: password,
			UUID:     generateUUID(true),
			PhoneID:  generateUUID(true),
		},
	}
	insta.fill()
	return insta
}

// TODO
func (insta *Instagram) fill() {
	if insta.User == nil {
		insta.User = NewUser(insta)
	}
	user := insta.User
	if insta.Current == nil {
		insta.Current = &ProfileData{}
	}

	if insta.Current.Feed == nil {
		insta.Current.Feed = NewFeed(user)
	}
	if insta.Inbox == nil {
		insta.Inbox = NewInbox(insta)
	}
	if insta.Current.Following == nil {
		insta.Current.Following = NewUsers(user, false)
	}
	if insta.Current.Followers == nil {
		insta.Current.Followers = NewUsers(user, true)
	}

	if insta.Current.insta == nil {
		insta.Current.insta = insta
	}
	if insta.Media == nil {
		insta.Media = NewMedia(insta)
	}
	if insta.Search == nil {
		insta.Search = NewSearch(insta)
	}
}
