package types

type ListData interface {
	MediaData
	isListData()
	GetDetails() MediaDetails
	GetItemList() ItemList
	SetItemList(ItemList)
	GetMediaType() MediaType
}

func (Collection) isListData() {}
func (Playlist) isListData()   {}

func (c Collection) GetItemList() ItemList          { return c.ItemList }
func (c Collection) GetDetails() MediaDetails       { return c.ItemList.Details }
func (c *Collection) SetItemList(itemList ItemList) { c.ItemList = itemList }
func (c Collection) GetMediaType() MediaType        { return MediaTypeCollection }

func (p Playlist) GetDetails() MediaDetails       { return p.ItemList.Details }
func (p Playlist) GetItemList() ItemList          { return p.ItemList }
func (p *Playlist) SetItemList(itemList ItemList) { p.ItemList = itemList }
func (p Playlist) GetMediaType() MediaType        { return MediaTypePlaylist }
