# AI Conversation Storage Enhancement Project

## Overview

This document outlines a plan to enhance the AI conversation capabilities in Suasor by implementing persistent storage for conversations, messages, and recommendations. The current system only stores conversations in memory, leading to data loss on service restarts and limiting analytical capabilities.

> **Implementation Note**: The implementation uses GORM for ORM functionality, integrated with the existing database system. The models are defined with appropriate GORM tags to enable automated table creation and relationship management.

## Current Limitations

1. **In-Memory Only Storage**: All conversation data is lost on service restarts
2. **No Historical Access**: Users cannot access or continue previous conversations
3. **Limited Analytics**: No ability to analyze conversation patterns or recommendation effectiveness
4. **No Recommendation Tracking**: Recommendations aren't permanently stored for future reference
5. **Scalability Issues**: Memory usage increases with active conversations

## Project Goals

1. Implement persistent storage for AI conversations
2. Enable users to view and resume past conversations
3. Track recommendations for future analysis
4. Provide analytics on conversation patterns and recommendation effectiveness
5. Maintain performance with scalable storage design

## Implementation Roadmap

### Phase 1: Data Model & Storage Design

- [x] Analyze current AI conversation implementation
- [x] Design database schema for conversation storage
- [x] Create repository interface definitions
- [x] Implement database models
- [x] Create database migrations

### Phase 2: Repository Implementation

- [x] Implement conversation repository
- [x] Implement message repository
- [x] Implement recommendation repository
- [ ] Add unit tests for repository layer

### Phase 3: Service Layer 

- [x] Create AI conversation service
- [x] Implement business logic for conversation management
- [x] Add recommendation tracking
- [x] Implement conversation history retrieval
- [ ] Add unit tests for service layer

### Phase 4: Handler & Client Integration

- [x] Update AI handlers to use new services
- [x] Modify Claude client to persist conversations
- [ ] Add endpoints for conversation history retrieval
    - [ ] GET /api/v1/user/ai/conversations
    - [ ] GET /api/v1/user/ai/conversations/:conversationId
    - [ ] GET /api/v1/user/ai/conversations/:conversationId/messages
    - [ ] GET /api/v1/user/ai/recommendations
    - [ ] POST /api/v1/user/ai/conversations/:conversationId/continue
    - [ ] PUT /api/v1/user/ai/conversations/:conversationId/archive
    - [ ] DELETE /api/v1/user/ai/conversations/:conversationId
- [x] Update service registrations in DI container

### Phase 5: Testing & Optimization

- [ ] Add integration tests for conversation flow
- [ ] Optimize database queries with indices
- [ ] Implement caching for active conversations
- [ ] Add cleanup job for old conversations

## Detailed Task List

### Database Models Implementation

1. Create `AIConversation` model:
   - ID (string): Unique conversation identifier
   - UserID (uint64): Owner of the conversation
   - ClientID (uint64): AI client used
   - ContentType (string): Type of content discussed
   - Status (string): active, archived, etc.
   - SystemPrompt (string): System instructions for AI
   - UserPreferences (string): JSON format of preferences
   - CreatedAt, UpdatedAt, ExpiresAt (time.Time): Timestamps
   - MessageCount (int): Number of messages
   - LastMessageTime (time.Time): Last activity timestamp

2. Create `AIMessage` model:
   - ID (string): Unique message identifier
   - ConversationID (string): Parent conversation
   - Role (string): "user" or "assistant"
   - Content (string): Message content
   - Timestamp (time.Time): When message was sent
   - Metadata (string): JSON format of additional data
   - TokenUsage (int): Tokens used for this message

3. Create `AIRecommendation` model:
   - ID (string): Unique recommendation identifier
   - MessageID (string): Message containing recommendation
   - ConversationID (string): Parent conversation
   - UserID (uint64): User who received recommendation
   - ItemType (string): Type of recommended item
   - Title (string): Title of recommended item
   - ExternalID (string): ID from external source
   - Data (string): Complete recommendation in JSON
   - Reason (string): Explanation for recommendation
   - CreatedAt (time.Time): When recommended
   - Selected (bool): If user selected this recommendation
   - SelectedAt (time.Time): When user selected it

4. Create `AIConversationAnalytics` model (optional):
   - ConversationID (string): Conversation reference
   - TotalUserMessages, TotalAssistantMessages (int): Message counts
   - TotalRecommendations, SelectedRecommendations (int): Recommendation tracking
   - TotalTokensUsed (int): Token usage tracking
   - AverageResponseTime (float64): Performance metric
   - ConversationDuration (int): Duration in seconds
   - LastUpdatedAt (time.Time): Last analytics update

### Repository Implementation

1. Create `AIConversationRepository` interface with methods:
   - CreateConversation
   - GetConversationByID 
   - GetConversationsByUserID
   - UpdateConversationStatus
   - DeleteConversation
   - ArchiveOldConversations
   - AddMessage
   - GetMessagesByConversationID
   - GetConversationHistory
   - AddRecommendation
   - GetRecommendationsByConversationID
   - GetRecommendationsByUserID
   - UpdateRecommendationSelection
   - GetConversationAnalytics
   - UpdateConversationAnalytics
   - GetUserConversationStats

2. Implement SQL repository that satisfies this interface:
   - Use transactions for multi-operation sequences
   - Add proper error handling
   - Implement efficient query patterns
   - Use prepared statements

### Service Layer Implementation

1. Create `AIConversationService` interface with methods:
   - StartConversation
   - SendMessage
   - GetConversationHistory
   - GetUserConversations
   - GetUserRecommendationHistory
   - ArchiveConversation
   - DeleteConversation
   - GetConversationInsights
   - GetUserAIInteractionSummary

2. Implement the service:
   - Add business logic for conversation management
   - Handle conversation state transitions
   - Implement recommendation extraction and storage
   - Add analytics calculations
   - Handle security checks (user ownership)

### Integration with Existing Code

1. Modify AI handlers:
   - Update StartConversation to use new service
   - Update SendConversationMessage to persist messages
   - Add new endpoints for history retrieval

2. Update Claude client:
   - Modify to work with persistence layer
   - Update conversation tracking to use database
   - Keep in-memory cache for active conversations

3. Update dependency injection:
   - Register new repository and service
   - Update handler constructor

## Database Schema

The models are defined with GORM tags to automatically create and manage the database schema:

### AIConversation Model

```go
// AIConversation represents a conversation between a user and an AI client
type AIConversation struct {
    ID               string         `gorm:"primaryKey;type:varchar(36)"`
    UserID           uint64         `gorm:"column:user_id;not null;index"`
    ClientID         uint64         `gorm:"column:client_id;not null"`
    ContentType      string         `gorm:"column:content_type;not null;type:varchar(50)"`
    Status           string         `gorm:"column:status;not null;default:active;index;type:varchar(20)"`
    SystemPrompt     string         `gorm:"column:system_prompt;type:text"`
    UserPreferences  string         `gorm:"column:user_preferences;type:json"`
    CreatedAt        time.Time      `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt        time.Time      `gorm:"column:updated_at;autoUpdateTime"`
    ExpiresAt        sql.NullTime   `gorm:"column:expires_at"`
    MessageCount     int            `gorm:"column:message_count;default:0"`
    LastMessageTime  sql.NullTime   `gorm:"column:last_message_time"`
    
    // Relationship fields
    Messages         []AIMessage         `gorm:"foreignKey:ConversationID"`
    Recommendations  []AIRecommendation  `gorm:"foreignKey:ConversationID"`
    Analytics        *AIConversationAnalytics `gorm:"foreignKey:ConversationID"`
}
```

### AIMessage Model

```go
// AIMessage represents a single message in an AI conversation
type AIMessage struct {
    ID              string         `gorm:"primaryKey;type:varchar(36)"`
    ConversationID  string         `gorm:"column:conversation_id;not null;index;type:varchar(36)"`
    Role            string         `gorm:"column:role;not null;type:varchar(20)"`
    Content         string         `gorm:"column:content;not null;type:text"`
    Timestamp       time.Time      `gorm:"column:timestamp;not null;index"`
    Metadata        sql.NullString `gorm:"column:metadata;type:json"`
    TokenUsage      int            `gorm:"column:token_usage;default:0"`
    
    // Relationship fields
    Conversation    *AIConversation    `gorm:"foreignKey:ConversationID"`
    Recommendations []AIRecommendation `gorm:"foreignKey:MessageID"`
}
```

### AIRecommendation Model

```go
// AIRecommendation represents a recommendation extracted from an AI conversation
type AIRecommendation struct {
    ID              string         `gorm:"primaryKey;type:varchar(36)"`
    MessageID       string         `gorm:"column:message_id;not null;type:varchar(36)"`
    ConversationID  string         `gorm:"column:conversation_id;not null;index;type:varchar(36)"`
    UserID          uint64         `gorm:"column:user_id;not null;index"`
    ItemType        string         `gorm:"column:item_type;not null;index;type:varchar(50)"`
    Title           string         `gorm:"column:title;not null;type:varchar(255)"`
    ExternalID      sql.NullString `gorm:"column:external_id;type:varchar(100)"`
    Data            string         `gorm:"column:data;not null;type:json"`
    Reason          sql.NullString `gorm:"column:reason;type:text"`
    CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime"`
    Selected        bool           `gorm:"column:selected;not null;default:false;index"`
    SelectedAt      sql.NullTime   `gorm:"column:selected_at"`
    
    // Relationship fields
    Message         *AIMessage      `gorm:"foreignKey:MessageID"`
    Conversation    *AIConversation `gorm:"foreignKey:ConversationID"`
}
```

### AIConversationAnalytics Model

```go
// AIConversationAnalytics tracks usage patterns and effectiveness
type AIConversationAnalytics struct {
    ConversationID         string    `gorm:"primaryKey;column:conversation_id;type:varchar(36)"`
    TotalUserMessages      int       `gorm:"column:total_user_messages;default:0"`
    TotalAssistantMessages int       `gorm:"column:total_assistant_messages;default:0"`
    TotalRecommendations   int       `gorm:"column:total_recommendations;default:0"`
    SelectedRecommendations int      `gorm:"column:selected_recommendations;default:0"`
    TotalTokensUsed        int       `gorm:"column:total_tokens_used;default:0"`
    AverageResponseTime    float64   `gorm:"column:average_response_time;default:0"`
    ConversationDuration   int       `gorm:"column:conversation_duration;default:0"`
    LastUpdatedAt          time.Time `gorm:"column:last_updated_at;autoUpdateTime"`
    
    // Relationship field
    Conversation    *AIConversation `gorm:"foreignKey:ConversationID"`
}
```

## API Endpoints

### New Endpoints

1. **GET** `/api/v1/user/ai/conversations`
   - List user's conversation history
   - Supports pagination and filtering

2. **GET** `/api/v1/user/ai/conversations/:conversationId`
   - Get details of a specific conversation
   - Includes metadata and statistics

3. **GET** `/api/v1/user/ai/conversations/:conversationId/messages`
   - Get all messages in a conversation
   - Supports pagination

4. **GET** `/api/v1/user/ai/recommendations`
   - Get user's recommendation history
   - Supports filtering by content type

5. **POST** `/api/v1/user/ai/conversations/:conversationId/continue`
   - Resume a previous conversation
   - Loads context from database

6. **PUT** `/api/v1/user/ai/conversations/:conversationId/archive`
   - Archive a conversation

7. **DELETE** `/api/v1/user/ai/conversations/:conversationId`
   - Delete a conversation

### Modified Endpoints

1. **POST** `/api/v1/client/:clientId/ai/conversation/start`
   - Updated to persist conversation to database

2. **POST** `/api/v1/client/:clientId/ai/conversation/message`
   - Updated to persist messages and extracted recommendations

## Performance Considerations

1. **Caching Strategy**:
   - Cache active conversations in memory
   - Consider Redis for distributed caching
   - Implement TTL for cache entries

2. **Database Optimization**:
   - Add appropriate indices for common queries
   - Consider partitioning for large datasets
   - Implement pagination for large result sets

3. **Clean-up Strategy**:
   - Implement scheduled job to archive old conversations
   - Add expiration policy for inactive conversations

## Analytics Potential

With this enhanced storage system, we can implement:

1. **User Engagement Metrics**:
   - Conversation frequency and duration
   - Response patterns
   - Feature usage

2. **Recommendation Effectiveness**:
   - Track which recommendations users select
   - Analyze patterns in accepted recommendations
   - Measure AI accuracy over time

3. **Content Insights**:
   - Discover trending content types
   - Identify recommendation gaps
   - Analyze user preference patterns

## Next Steps

1. Begin with database model implementation
2. Create migration scripts
3. Implement repository layer
4. Build service layer
5. Update handlers and clients
6. Add tests
7. Deploy and monitor performance