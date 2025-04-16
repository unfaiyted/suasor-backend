package app

// import (
// 	"suasor/repository"
// 	"suasor/services"
// )

//
// // mediaServicesImpl implements MediaServices
// type mediaServicesImpl struct {
// 	personService *services.PersonService
// 	creditService *services.CreditService
// }
//
// func (m *mediaServicesImpl) PersonService() *services.PersonService {
// 	return m.personService
// }
//
// func (m *mediaServicesImpl) CreditService() *services.CreditService {
// 	return m.creditService
// }
//
// // NewMediaServices creates a new MediaServices instance
// func NewMediaServices(
// 	personRepo repository.PersonRepository,
// 	creditRepo repository.CreditRepository,
// ) MediaServices {
// 	return &mediaServicesImpl{
// 		personService: services.NewPersonService(personRepo, creditRepo),
// 		creditService: services.NewCreditService(creditRepo, personRepo),
// 	}
// }
