package types

type ListData interface {
	MediaData
	isListData()

	GetTitle() string
	GetDetails() MediaDetails
	GetItemList() ItemList

	SetItemList(ItemList)
	SetDetails(MediaDetails)
	AddListItem(ListItem)

	GetMediaType() MediaType
}

func (Collection) isListData() {}

func (c Collection) GetItemList() ItemList    { return c.ItemList }
func (c Collection) GetDetails() MediaDetails { return c.ItemList.Details }
func (c *Collection) SetDetails(details MediaDetails) {
	c.ItemList.Details = details
}
func (c *Collection) SetItemList(itemList ItemList) { c.ItemList = itemList }
func (c Collection) GetMediaType() MediaType        { return MediaTypeCollection }

func (c *Collection) AddListItem(item ListItem) {
	c.ItemList.AddItem(item)
}

func (c *Collection) AddListItemWithClientID(item ListItem, clientID uint64) {
	c.ItemList.AddItemWithClientID(item, clientID)
}

func (c *Collection) RemoveItem(itemID uint64, clientID uint64) error {
	return c.ItemList.RemoveItem(itemID, clientID)
}

func (p Collection) GetTitle() string { return p.Details.Title }

func (Playlist) isListData()                {}
func (p Playlist) GetDetails() MediaDetails { return p.ItemList.Details }
func (p *Playlist) SetDetails(details MediaDetails) {
	p.ItemList.Details = details
}
func (p Playlist) GetItemList() ItemList          { return p.ItemList }
func (p *Playlist) SetItemList(itemList ItemList) { p.ItemList = itemList }
func (p Playlist) GetMediaType() MediaType        { return MediaTypePlaylist }

func (p *Playlist) AddListItem(item ListItem) {
	p.ItemList.AddItem(item)
}

func (p *Playlist) AddListItemWithClientID(item ListItem, clientID uint64) {
	p.ItemList.AddItemWithClientID(item, clientID)
}

func (p *Playlist) RemoveItem(itemID uint64, clientID uint64) error {
	return p.ItemList.RemoveItem(itemID, clientID)
}

// GetTitle
func (p Playlist) GetTitle() string { return p.Details.Title }
