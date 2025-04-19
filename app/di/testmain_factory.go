package di

func TestMediaFactory() {
	// Just a test function
	db := &gorm.DB{}
	clientFactory := &client.ClientFactoryService{}
	_ = createMediaDataFactory(db, clientFactory)
}
