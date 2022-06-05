package main

import (
	"flag"
	"fmt"
	"math/rand"

	"github.com/charmbracelet/lipgloss"
	"github.com/gempir/go-twitch-irc/v3"
)

var channel string

type messageHistory struct {
	users map[string]*lipgloss.Style
	msgs  []*twitch.PrivateMessage
}

func init() {
	flag.StringVar(&channel, "user", "Raven_Tech", "Join <users> chat channel.")

	flag.Parse()
}

func main() {
	client := twitch.NewAnonymousClient() // for an anonymous user (no write capabilities)
	// client := twitch.NewClient("yourtwitchusername", "oauth:123123123")

	state := NewMessageHistory()

	client.OnPrivateMessage(privateMessageFunc(state))
	client.OnConnect(connected)

	client.Join(channel)

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

func NewMessageHistory() *messageHistory {
	return &messageHistory{
		users: make(map[string]*lipgloss.Style),
	}
}

func privateMessageFunc(m *messageHistory) func(twitch.PrivateMessage) {
	return func(message twitch.PrivateMessage) {
		fmt.Printf("<%v>: %v\n",
			m.getStyle(message.User).Render(message.User.DisplayName),
			message.Message,
		)
	}
}

func (m *messageHistory) getStyle(user twitch.User) *lipgloss.Style {
	var style *lipgloss.Style
	var ok bool

	if style, ok = m.users[user.ID]; !ok {
		if len(user.Color) == 0 {
			user.Color = fmt.Sprintf(
				"#%x%x%x",
				rand.Intn(256),
				rand.Intn(256),
				rand.Intn(256),
			)
		}

		tstyle := lipgloss.
			NewStyle().
			Foreground(lipgloss.Color(user.Color))
		style = &tstyle

		m.users[user.ID] = style
	}

	return style
}

func connected() {
	fmt.Println("Connected to server.")
}
