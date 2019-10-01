package components

import (
	"fmt"
	"html"

	"github.com/erroneousboat/termui"
	"github.com/renstrom/fuzzysearch/fuzzy"
)

const (
	IconOnline       = "●"
	IconOffline      = "○"
	IconChannel      = "#"
	IconGroup        = "☰"
	IconIM           = "●"
	IconMpIM         = "☰"
	IconNotification = "*"

	PresenceAway   = "away"
	PresenceActive = "active"

	ChannelTypeChannel = "channel"
	ChannelTypeGroup   = "group"
	ChannelTypeIM      = "im"
	ChannelTypeMpIM    = "mpim"
)

type ChannelItem struct {
	ID             string
	Name           string
	Topic          string
	Type           string
	UserID         string
	Presence       string
	Notification   bool
	StylePrefix    string
	StyleIcon      string
	StyleText      string
	IsSearchResult bool
}

// ToString will set the label of the channel, how it will be
// displayed on screen. Based on the type, different icons are
// shown, as well as an optional notification icon.
func (c ChannelItem) ToString() string {
	var prefix string
	if c.Notification {
		prefix = IconNotification
	} else {
		prefix = " "
	}

	var icon string
	switch c.Type {
	case ChannelTypeChannel:
		icon = IconChannel
	case ChannelTypeGroup:
		icon = IconGroup
	case ChannelTypeMpIM:
		icon = IconMpIM
	case ChannelTypeIM:
		switch c.Presence {
		case PresenceActive:
			icon = IconOnline
		case PresenceAway:
			icon = IconOffline
		default:
			icon = IconIM
		}
	}

	label := fmt.Sprintf(
		"[%s](%s) [%s](%s) [%s](%s)",
		prefix, c.StylePrefix,
		icon, c.StyleIcon,
		c.Name, c.StyleText,
	)

	return label
}

// GetChannelName will return a formatted representation of the
// name of the channel
func (c ChannelItem) GetChannelName() string {
	var channelName string
	if c.Topic != "" {
		channelName = fmt.Sprintf("%s - %s",
			html.UnescapeString(c.Name),
			html.UnescapeString(c.Topic),
		)
	} else {
		channelName = c.Name
	}
	return channelName
}

// Channels is the definition of a Channels component
type Channels struct {
	ChannelItems    []ChannelItem // sorted list of channels
	CursorPosition  int
	List            *termui.List // ui of visible channels
	Offset          int          // from what offset are channels rendered
	SearchPosition  int          // current position in search results
	SelectedChannel string       // index of which channel is selected from the List
	UnreadOnly      bool         // only show unread messages when on
}

// CreateChannels is the constructor for the Channels component
func CreateChannelsComponent(height int, unreadOnly bool) *Channels {
	channels := &Channels{
		List: termui.NewList(),
	}

	channels.List.BorderLabel = "Conversations"
	channels.List.Height = height

	channels.SelectedChannel = ""
	channels.Offset = 0
	channels.UnreadOnly = unreadOnly

	return channels
}

// List lists all visible channels
// always includes the currently selected channel (if one is selected)
// if the app is configured to only show messages with unread messages, filter out read channels
func (c *Channels) ListChannels() (items []ChannelItem) {
	if !c.UnreadOnly {
		items = c.ChannelItems
		return
	}

	for _, chn := range c.ChannelItems {
		if chn.ID == c.SelectedChannel || chn.Notification || chn.IsSearchResult {
			items = append(items, chn)
		}
	}

	return
}

// ListSearchResults lists all channels that are results of the current search
func (c *Channels) ListSearchResults() (items []ChannelItem) {
	for _, chn := range c.ListChannels() {
		if chn.IsSearchResult {
			items = append(items, chn)
		}
	}

	return
}

// Buffer implements interface termui.Bufferer
func (c *Channels) Buffer() termui.Buffer {
	buf := c.List.Buffer()

	c.locateCursor()

	for i, item := range c.ListChannels()[c.Offset:] {

		y := c.List.InnerBounds().Min.Y + i

		if y > c.List.InnerBounds().Max.Y-1 {
			break
		}

		// Set the visible cursor
		var cells []termui.Cell
		if y == c.CursorPosition {
			cells = termui.DefaultTxBuilder.Build(
				item.ToString(), c.List.ItemBgColor, c.List.ItemFgColor)
		} else {
			cells = termui.DefaultTxBuilder.Build(
				item.ToString(), c.List.ItemFgColor, c.List.ItemBgColor)
		}

		// Append ellipsis when overflows
		cells = termui.DTrimTxCls(cells, c.List.InnerWidth())

		x := c.List.InnerBounds().Min.X
		for _, cell := range cells {
			buf.Set(x, y, cell)
			x += cell.Width()
		}

		// When not at the end of the pane fill it up empty characters
		for x < c.List.InnerBounds().Max.X {
			if y == c.CursorPosition {
				buf.Set(x, y,
					termui.Cell{
						Ch: ' ',
						Fg: c.List.ItemBgColor,
						Bg: c.List.ItemFgColor,
					},
				)
			} else {
				buf.Set(
					x, y,
					termui.Cell{
						Ch: ' ',
						Fg: c.List.ItemFgColor,
						Bg: c.List.ItemBgColor,
					},
				)
			}
			x++
		}
	}

	return buf
}

// GetHeight implements interface termui.GridBufferer
func (c *Channels) GetHeight() int {
	return c.List.Block.GetHeight()
}

// SetWidth implements interface termui.GridBufferer
func (c *Channels) SetWidth(w int) {
	c.List.SetWidth(w)
}

// SetX implements interface termui.GridBufferer
func (c *Channels) SetX(x int) {
	c.List.SetX(x)
}

// SetY implements interface termui.GridBufferer
func (c *Channels) SetY(y int) {
	c.List.SetY(y)
}

func (c *Channels) SetChannels(channels []ChannelItem) {
	c.ChannelItems = channels

	c.locateCursor()

	// set the current channel to the first one in the list
	// when unread-only is off
	if !c.UnreadOnly && len(c.ChannelItems) > 0 {
		c.SetSelectedChannel(c.ChannelItems[0].ID)
	}
}

func (c *Channels) MarkAsRead(channelID string) {

	if index, ok := c.FindChannel(channelID); ok {
		c.ChannelItems[index].Notification = false
	}
}

func (c *Channels) MarkAsUnread(channelID string) {

	if index, ok := c.FindChannel(channelID); ok {
		c.ChannelItems[index].Notification = true
	}
}

func (c *Channels) SetPresence(channelID string, presence string) {
	if index, ok := c.FindChannel(channelID); ok {
		c.ChannelItems[index].Presence = presence
	}
}

// FindChannel finds the index of the channel in ChannelItems
// it is not necessarily visible
func (c *Channels) FindChannel(channelID string) (index int, ok bool) {
	for i, channel := range c.ChannelItems {
		if channel.ID == channelID {
			index = i
			ok = true
			break
		}
	}

	return index, ok
}

// FindVisibleChannels finds a channel by ID in the visible channel list
func (c *Channels) FindVisibleChannel(channelID string) (chn ChannelItem, ok bool) {

	for _, channel := range c.ListChannels() {
		if channel.ID == channelID {
			chn = channel
			ok = true
			break
		}
	}

	return chn, ok
}

// SetSelectedChannel sets the SelectedChannel given its ID
func (c *Channels) SetSelectedChannel(channelID string) {

	if _, ok := c.FindChannel(channelID); ok {
		c.SelectedChannel = channelID
	}

	c.locateCursor()
}

// Get SelectedChannel returns the ChannelItem that is currently selected
func (c *Channels) GetSelectedChannel() (selected ChannelItem, ok bool) {

	var index int
	if index, ok = c.FindChannel(c.SelectedChannel); ok {
		selected = c.ChannelItems[index]
	}

	return
}

func (c *Channels) locateCursor() (prev, curr, next int) {

	c.CursorPosition = c.List.InnerBounds().Min.Y
	channels := c.ListChannels()

	if c.SelectedChannel == "" && len(channels) > 0 {
		c.SetSelectedChannel(channels[0].ID)
		return
	}

	for i := 0; i < len(channels)-1; i++ {
		chn := channels[i]

		c.ScrollDown()
		c.ScrollUp()

		if chn.ID == c.SelectedChannel {
			curr = curr + i - c.Offset + 1
			prev = curr - 1
			next = curr + 1
			c.CursorPosition = curr
			return
		}

	}

	return
}

// get the channel item that preceeds the selected channel item in the visible channel list
// TODO use locateCursor
func (c *Channels) getPreviousItem() (prev ChannelItem, ok bool) {

	for _, curr := range c.ListChannels() {
		if curr.ID == c.SelectedChannel && prev.ID != "" {
			ok = true
			break
		}

		prev = curr
	}

	return
}

// get the channel item that proceeds the selected channel item in the visible channel list
// TODO use locateCursor
func (c *Channels) getNextItem() (next ChannelItem, pos int, ok bool) {

	channels := c.ListChannels()
	for i := len(channels) - 1; i >= 0; i-- {
		curr := channels[i]
		if curr.ID == c.SelectedChannel && next.ID != "" {
			pos = i
			ok = true
			break
		}

		next = curr
	}

	return
}

// MoveCursorUp will decrease the SelectedChannel by 1
// returns the item over which the cursor is now hovering
func (c *Channels) MoveCursorUp() (hovering ChannelItem, ok bool) {

	if hovering, ok = c.getPreviousItem(); ok {
		c.SetSelectedChannel(hovering.ID)
		c.ScrollUp()
	}

	return
}

// MoveCursorDown will increase the SelectedChannel by 1
// returns the item over which the cursor is now hovering
func (c *Channels) MoveCursorDown() (hovering ChannelItem, ok bool) {

	if hovering, _, ok = c.getNextItem(); ok {
		c.SetSelectedChannel(hovering.ID)
		c.ScrollDown()
	}

	return
}

// MoveCursorTop will move the cursor to the top of the channels
func (c *Channels) MoveCursorTop() {
	if list := c.ListChannels(); len(list) > 0 {
		c.SetSelectedChannel(list[0].ID)
	}

	c.Offset = 0
}

// MoveCursorBottom will move the cursor to the bottom of the channels
func (c *Channels) MoveCursorBottom() {
	if list := c.ListChannels(); len(list) > 0 {
		index := len(list) - 1
		c.SetSelectedChannel(list[index].ID)

		offset := len(list) - (c.List.InnerBounds().Max.Y - 1)

		if offset < 0 {
			c.Offset = 0
		} else {
			c.Offset = offset
		}
	}
}

// ScrollUp enables us to scroll through the channel list when it overflows
func (c *Channels) ScrollUp() {

	// Is cursor at the top of the channel view?
	if c.CursorPosition == c.List.InnerBounds().Min.Y {
		if c.Offset > 0 {
			c.Offset--
		}
	}
}

// ScrollDown enables us to scroll through the channel list when it overflows
// pos is the absolute position of the current item
func (c *Channels) ScrollDown() {

	// Is the cursor at the bottom of the channel view?
	if c.CursorPosition >= c.List.InnerBounds().Max.Y {
		if c.Offset < len(c.ChannelItems)-1 {
			c.Offset++
		}
	}
}

// Search will search through the channels to find a channel,
// when a match has been found the selected channel will then
// be the channel that has been found
func (c *Channels) Search(term string) (resultCount int) {

	for i, chn := range c.ChannelItems {
		if chn.IsSearchResult {
			c.ChannelItems[i].IsSearchResult = false
		}
	}

	targets := make([]string, 0)
	for _, c := range c.ChannelItems {
		targets = append(targets, c.Name)
	}

	matches := fuzzy.Find(term, targets)

	for _, m := range matches {
		for i, item := range c.ChannelItems {
			if m == item.Name {
				resultCount = resultCount + 1
				c.ChannelItems[i].IsSearchResult = true
				break
			}
		}
	}

	if resultCount > 0 {
		c.GotoPosition(0)
		c.SearchPosition = 0
	}

	return
}

// GotoPosition is used by to automatically scroll to a specific
// location in the channels component
func (c *Channels) GotoPosition(newPos int) (ok bool) {

	// there's nothing to be done if there are no search results, or the given position
	// is out of range
	var results []ChannelItem
	if results = c.ListSearchResults(); len(results) == 0 || newPos > len(results)-1 {
		return
	}

	// Is the new position in range of the current view?
	minRange := c.Offset
	maxRange := c.Offset + (c.List.InnerBounds().Max.Y - 2)

	var newChannelIndex int
	if newChannelIndex, ok = c.FindChannel(results[newPos].ID); !ok {
		return
	}

	newChannelID := c.ChannelItems[newChannelIndex].ID

	if newChannelIndex < minRange {
		// How much do we need to scroll to get it into range?
		c.Offset = c.Offset - (minRange - newChannelIndex)

		// newPos is above, we need to scroll up.
		c.SetSelectedChannel(newChannelID)
	} else if newChannelIndex > maxRange {
		// How much do we need to scroll to get it into range?
		c.Offset = c.Offset + (newChannelIndex - maxRange)

		// newPos is below, we need to scroll down
		c.SetSelectedChannel(newChannelID)
	} else {
		// newPos is inside range
		c.SetSelectedChannel(newChannelID)
	}

	return true
}

// SearchNext allows us to cycle through search results
func (c *Channels) SearchNext() {
	newPosition := c.SearchPosition + 1
	if ok := c.GotoPosition(newPosition); ok {
		c.SearchPosition = newPosition
	}
}

// SearchPrev allows us to cycle through resrch results
func (c *Channels) SearchPrev() {
	newPosition := c.SearchPosition - 1
	if ok := c.GotoPosition(newPosition); ok {
		c.SearchPosition = newPosition
	}
}

// Jump to the first channel with a notification
func (c *Channels) Jump() {
	for i, _ := range c.ListChannels() {
		c.GotoPosition(i)
		break
	}
}
