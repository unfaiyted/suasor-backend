package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// AIConversation represents a conversation between a user and an AI client
type AIConversation struct {
	ID               string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	UserID           uint64         `json:"userId" gorm:"column:user_id;not null;index"`
	ClientID         uint64         `json:"clientId" gorm:"column:client_id;not null"`
	ContentType      string         `json:"contentType" gorm:"column:content_type;not null;type:varchar(50)"`
	Status           string         `json:"status" gorm:"column:status;not null;default:active;index;type:varchar(20)"`
	SystemPrompt     string         `json:"systemPrompt" gorm:"column:system_prompt;type:text"`
	UserPreferences  string         `json:"userPreferences" gorm:"column:user_preferences;type:json"`
	CreatedAt        time.Time      `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time      `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
	ExpiresAt        sql.NullTime   `json:"expiresAt,omitempty" gorm:"column:expires_at"`
	MessageCount     int            `json:"messageCount" gorm:"column:message_count;default:0"`
	LastMessageTime  sql.NullTime   `json:"lastMessageTime,omitempty" gorm:"column:last_message_time"`
	
	// Relationship fields
	Messages         []AIMessage         `json:"messages,omitempty" gorm:"foreignKey:ConversationID"`
	Recommendations  []AIRecommendation  `json:"recommendations,omitempty" gorm:"foreignKey:ConversationID"`
	Analytics        *AIConversationAnalytics `json:"analytics,omitempty" gorm:"foreignKey:ConversationID"`
}

// TableName specifies the table name for GORM
func (AIConversation) TableName() string {
	return "ai_conversations"
}

// BeforeCreate hook to set updated_at if not set
func (c *AIConversation) BeforeCreate(tx *gorm.DB) error {
	if c.UpdatedAt.IsZero() {
		c.UpdatedAt = time.Now()
	}
	return nil
}

// GetUserPreferencesMap returns the user preferences as a map
func (c *AIConversation) GetUserPreferencesMap() (map[string]any, error) {
	if c.UserPreferences == "" {
		return make(map[string]any), nil
	}
	
	var prefs map[string]any
	err := json.Unmarshal([]byte(c.UserPreferences), &prefs)
	if err != nil {
		return nil, err
	}
	
	return prefs, nil
}

// SetUserPreferencesMap sets the user preferences from a map
func (c *AIConversation) SetUserPreferencesMap(prefs map[string]any) error {
	data, err := json.Marshal(prefs)
	if err != nil {
		return err
	}
	
	c.UserPreferences = string(data)
	return nil
}

// NewAIConversation creates a new conversation with default values
func NewAIConversation(id string, userID, clientID uint64, contentType, systemPrompt string) *AIConversation {
	now := time.Now()
	return &AIConversation{
		ID:              id,
		UserID:          userID,
		ClientID:        clientID,
		ContentType:     contentType,
		Status:          "active",
		SystemPrompt:    systemPrompt,
		UserPreferences: "{}",
		CreatedAt:       now,
		UpdatedAt:       now,
		MessageCount:    0,
	}
}

// AIMessage represents a single message in an AI conversation
type AIMessage struct {
	ID              string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ConversationID  string         `json:"conversationId" gorm:"column:conversation_id;not null;index;type:varchar(36)"`
	Role            string         `json:"role" gorm:"column:role;not null;type:varchar(20)"` // "user" or "assistant"
	Content         string         `json:"content" gorm:"column:content;not null;type:text"`
	Timestamp       time.Time      `json:"timestamp" gorm:"column:timestamp;not null;index"`
	Metadata        sql.NullString `json:"metadata,omitempty" gorm:"column:metadata;type:json"` // JSON storage for optional metadata
	TokenUsage      int            `json:"tokenUsage" gorm:"column:token_usage;default:0"` // Track token usage
	
	// Relationship fields
	Conversation    *AIConversation    `json:"-" gorm:"foreignKey:ConversationID"`
	Recommendations []AIRecommendation `json:"recommendations,omitempty" gorm:"foreignKey:MessageID"`
}

// TableName specifies the table name for GORM
func (AIMessage) TableName() string {
	return "ai_messages"
}

// GetMetadataMap returns the metadata as a map
func (m *AIMessage) GetMetadataMap() (map[string]any, error) {
	if !m.Metadata.Valid || m.Metadata.String == "" {
		return make(map[string]any), nil
	}
	
	var metadata map[string]any
	err := json.Unmarshal([]byte(m.Metadata.String), &metadata)
	if err != nil {
		return nil, err
	}
	
	return metadata, nil
}

// SetMetadataMap sets the metadata from a map
func (m *AIMessage) SetMetadataMap(metadata map[string]any) error {
	if metadata == nil {
		m.Metadata = sql.NullString{Valid: false}
		return nil
	}
	
	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	
	m.Metadata = sql.NullString{
		String: string(data),
		Valid:  true,
	}
	return nil
}

// NewAIMessage creates a new message with default values
func NewAIMessage(id string, conversationID string, role string, content string) *AIMessage {
	return &AIMessage{
		ID:             id,
		ConversationID: conversationID,
		Role:           role,
		Content:        content,
		Timestamp:      time.Now(),
		TokenUsage:     0,
	}
}

// AIRecommendation represents a recommendation extracted from an AI conversation
type AIRecommendation struct {
	ID              string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	MessageID       string         `json:"messageId" gorm:"column:message_id;not null;type:varchar(36)"`
	ConversationID  string         `json:"conversationId" gorm:"column:conversation_id;not null;index;type:varchar(36)"`
	UserID          uint64         `json:"userId" gorm:"column:user_id;not null;index"`
	ItemType        string         `json:"itemType" gorm:"column:item_type;not null;index;type:varchar(50)"` // Movie, Series, Music, etc.
	Title           string         `json:"title" gorm:"column:title;not null;type:varchar(255)"`
	ExternalID      sql.NullString `json:"externalId,omitempty" gorm:"column:external_id;type:varchar(100)"` // ID from external source if available
	Data            string         `json:"data" gorm:"column:data;not null;type:json"` // JSON storage of complete recommendation
	Reason          sql.NullString `json:"reason,omitempty" gorm:"column:reason;type:text"` // Why it was recommended
	CreatedAt       time.Time      `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	Selected        bool           `json:"selected" gorm:"column:selected;not null;default:false;index"` // If user selected this recommendation
	SelectedAt      sql.NullTime   `json:"selectedAt,omitempty" gorm:"column:selected_at"`
	
	// Relationship fields
	Message         *AIMessage      `json:"-" gorm:"foreignKey:MessageID"`
	Conversation    *AIConversation `json:"-" gorm:"foreignKey:ConversationID"`
}

// TableName specifies the table name for GORM
func (AIRecommendation) TableName() string {
	return "ai_recommendations"
}

// GetDataMap returns the recommendation data as a map
func (r *AIRecommendation) GetDataMap() (map[string]any, error) {
	if r.Data == "" {
		return make(map[string]any), nil
	}
	
	var data map[string]any
	err := json.Unmarshal([]byte(r.Data), &data)
	if err != nil {
		return nil, err
	}
	
	return data, nil
}

// SetDataMap sets the recommendation data from a map
func (r *AIRecommendation) SetDataMap(data map[string]any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	r.Data = string(jsonData)
	return nil
}

// NewAIRecommendation creates a new recommendation with default values
func NewAIRecommendation(id string, messageID string, conversationID string, userID uint64, itemType string, title string) *AIRecommendation {
	return &AIRecommendation{
		ID:             id,
		MessageID:      messageID,
		ConversationID: conversationID,
		UserID:         userID,
		ItemType:       itemType,
		Title:          title,
		Data:           "{}",
		CreatedAt:      time.Now(),
		Selected:       false,
	}
}

// AIConversationAnalytics tracks usage patterns and effectiveness
type AIConversationAnalytics struct {
	ConversationID         string    `json:"conversationId" gorm:"primaryKey;column:conversation_id;type:varchar(36)"`
	TotalUserMessages      int       `json:"totalUserMessages" gorm:"column:total_user_messages;default:0"`
	TotalAssistantMessages int       `json:"totalAssistantMessages" gorm:"column:total_assistant_messages;default:0"`
	TotalRecommendations   int       `json:"totalRecommendations" gorm:"column:total_recommendations;default:0"`
	SelectedRecommendations int      `json:"selectedRecommendations" gorm:"column:selected_recommendations;default:0"`
	TotalTokensUsed        int       `json:"totalTokensUsed" gorm:"column:total_tokens_used;default:0"`
	AverageResponseTime    float64   `json:"averageResponseTime" gorm:"column:average_response_time;default:0"`
	ConversationDuration   int       `json:"conversationDuration" gorm:"column:conversation_duration;default:0"` // In seconds
	LastUpdatedAt          time.Time `json:"lastUpdatedAt" gorm:"column:last_updated_at;autoUpdateTime"`
	
	// Relationship field
	Conversation    *AIConversation `json:"-" gorm:"foreignKey:ConversationID"`
}

// TableName specifies the table name for GORM
func (AIConversationAnalytics) TableName() string {
	return "ai_conversation_analytics"
}

// NewAIConversationAnalytics creates a new analytics record with default values
func NewAIConversationAnalytics(conversationID string) *AIConversationAnalytics {
	return &AIConversationAnalytics{
		ConversationID:         conversationID,
		TotalUserMessages:      0,
		TotalAssistantMessages: 0,
		TotalRecommendations:   0,
		SelectedRecommendations: 0,
		TotalTokensUsed:        0,
		AverageResponseTime:    0,
		ConversationDuration:   0,
		LastUpdatedAt:          time.Now(),
	}
}

// UserConversationStats represents aggregated statistics for a user's AI interactions
type UserConversationStats struct {
	UserID                 uint64    `json:"userId" gorm:"column:user_id"`
	TotalConversations     int       `json:"totalConversations" gorm:"column:total_conversations"`
	TotalMessages          int       `json:"totalMessages" gorm:"column:total_messages"`
	TotalRecommendations   int       `json:"totalRecommendations" gorm:"column:total_recommendations"`
	SelectedRecommendations int      `json:"selectedRecommendations" gorm:"column:selected_recommendations"`
	AverageMessagesPerConv float64   `json:"averageMessagesPerConv" gorm:"column:average_messages_per_conv"`
	FavoriteContentType    string    `json:"favoriteContentType" gorm:"column:favorite_content_type"`
	LastConversation       time.Time `json:"lastConversation" gorm:"column:last_conversation"`
}

// TableName specifies the table name for GORM 
// This is a view-like object that doesn't have a DB table
func (UserConversationStats) TableName() string {
	return ""
}

// AIInteractionSummary provides a high-level view of user's AI usage
// This is a DTO that combines multiple data sources
type AIInteractionSummary struct {
	UserID              uint64                  `json:"userId"`
	Stats               UserConversationStats   `json:"stats"`
	RecentConversations []*AIConversation       `json:"recentConversations"`
	TopRecommendations  []*AIRecommendation     `json:"topRecommendations"`
	ContentTypeBreakdown map[string]int         `json:"contentTypeBreakdown"`
}

// ConversationInsights provides detailed analytics for a specific conversation
// This is a DTO that combines multiple data sources
type ConversationInsights struct {
	Conversation    *AIConversation          `json:"conversation"`
	Analytics       *AIConversationAnalytics `json:"analytics"`
	MessageSummary  map[string]int           `json:"messageSummary"`
	TimingData      map[string]any           `json:"timingData"`
	RecommendationEffectiveness float64      `json:"recommendationEffectiveness"`
}