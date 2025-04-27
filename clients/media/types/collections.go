package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Collection represents an unordered set of unique media items
type Collection struct {
	ItemList
	itemMap map[uint64]int // Maps IDs to indices for fast lookup
}

// Initialize the collection's item map
func (c *Collection) initItemMap() {
	if c.itemMap == nil {
		c.itemMap = make(map[uint64]int, len(c.Items))
		for i, item := range c.Items {
			c.itemMap[item.ItemID] = i
		}
	}
}

// FindItemByID optimized for Collections using the item map
func (c *Collection) FindItemByID(id uint64) (ListItem, int, bool) {
	c.initItemMap()
	if idx, exists := c.itemMap[id]; exists && idx < len(c.Items) && c.Items[idx].ItemID == id {
		return c.Items[idx], idx, true
	}
	// Fallback to linear search in case map is out of sync
	return c.ItemList.FindItemByID(id)
}

// AddItem overriden to ensure uniqueness in collections
func (c *Collection) AddItem(item ListItem, clientID uint64) {
	c.initItemMap()

	// Check for existing item first - ignore duplicates
	if _, exists := c.itemMap[item.ItemID]; exists {
		return
	}

	// Call parent implementation
	c.ItemList.AddItemWithClientID(item, clientID)

	// Update the map
	c.itemMap[item.ItemID] = len(c.Items) - 1
}

func (c *Collection) isListData() {}

func (c *Collection) GetItemList() ItemList    { return c.ItemList }
func (c *Collection) GetDetails() MediaDetails { return c.ItemList.Details }
func (c *Collection) SetDetails(details MediaDetails) {
	c.ItemList.Details = details
}
func (c *Collection) SetItemList(itemList ItemList) { c.ItemList = itemList }
func (c *Collection) GetMediaType() MediaType       { return MediaTypeCollection }

func (c *Collection) AddListItem(item ListItem) {
	c.ItemList.AddItem(item)
}

func (c *Collection) AddListItemWithClientID(item ListItem, clientID uint64) {
	c.ItemList.AddItemWithClientID(item, clientID)
}

func (c *Collection) RemoveItem(itemID uint64, clientID uint64) error {
	return c.ItemList.RemoveItem(itemID, clientID)
}

func (p *Collection) GetTitle() string { return p.Details.Title }
func (*Collection) isMediaData()       {}

// EnsureNoDuplicates is now much more efficient with the item map
func (c *Collection) EnsureNoDuplicates() {
	c.initItemMap()
	seen := make(map[uint64]bool)
	uniqueItems := make([]ListItem, 0, len(c.Items))

	for _, item := range c.Items {
		if !seen[item.ItemID] {
			seen[item.ItemID] = true
			uniqueItems = append(uniqueItems, item)
		}
	}

	c.Items = uniqueItems
	c.ItemCount = len(c.Items)

	// Rebuild the map
	c.itemMap = make(map[uint64]int, len(c.Items))
	for i, item := range c.Items {
		c.itemMap[item.ItemID] = i
	}
}

// Scan
func (m *Collection) Scan(value any) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, m)
}

func (m *Collection) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}
