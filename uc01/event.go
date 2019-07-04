package uc01

import (
	"fmt"
	"encoding/json"
	"github.com/go-redis/redis"

	"github.com/mberlanda/ru101-go/utils"
)

// Event ...
type Event struct {
	Sku            string `json:"sku"`
	Name           string `json:"name"`
	DisabledAccess bool   `json:"disabled_access"`
	MedalEvent     bool   `json:"medal_event"`
	Venue          string `json:"venue"`
	Category       string `json:"category"`
}

// DefaultEventList ...
const DefaultEventList = []Event{
	Event{
		Sku: "123-ABC-723",
		Name: "Men's 100m Final",
		DisabledAccess: true,
		MedalEvent: true,
		Venue: "Olympic Stadium",
		Category: "Track & Field"
   },
   Event{
	   Sku: "737-DEF-911",
		Name: "Women's 4x100m Heats",
		DisabledAccess: true,
		MedalEvent: false,
		Venue: "Olympic Stadium",
		Category: "Track & Field"
   },
   Event{
	   Sku: "320-GHI-921",
		Name: "Womens Judo Qualifying",
		DisabledAccess: false,
		MedalEvent: false,
		Venue: "Nippon Budokan",
		Category: "Martial Arts"
   }
}

// ToJSON ...
func (e *Event) ToJSON() string {
	bs, _ := json.Marshal(e)
	return string(bs)
}

// ParseEvent ...
func ParseEvent(je string) *Event {
	e := Event{}
	json.Unmarshal([]byte(je), &e)
	return &e
}

// Uc01 is a struct to wrap the module dependency
type Uc01 struct {
	Client *redis.Client
	KeyNameHelper *utils.KeyNameHelper
}

func (u *Uc01) EventKeyName(sku string) string {
	return u.KeyNameHelper.CreateKeyName([]string{"event", sku})
}

// CreateEvents ...
func (u *Uc01) CreateEvents(evts []Event) {
	for _, e := range evts {
		u.Client.Set(u.EventKeyName(e.Sku), e.ToJSON(), 0)
	}
}

// PrintEventName ...
func (u *Uc01) PrintEventName(sku string) {
	e := ParseEvent(u.Client.Get(u.EventKeyName(sku)))
	if e.Name != "" {
		fmt.Print(e.Name)
		return
	}
	fmt.Print(e.Sku)
}

// PrintEventName ...
func (u *Uc01) PrintEventName(sku string) []string{
	matches := []string{}
}