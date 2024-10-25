package store

type StoreModule struct {
	Create             useCaseCreate
	Update             useCaseUpdate
	AddBusinessHour    useCaseAddBusinessHour
	DeleteBusinessHour useCaseDeleteBusinessHour
	GetByID            useCaseGetByID
	GetByFilter        useCaseGetByFilter
	SetProfileImage    useCaseSetProfileImage
	SetHeaderImage     useCaseSetHeaderImage
}

func NewStoreModule(storeRepository Repository) StoreModule {
	return StoreModule{
		Create:             newUseCaseCreate(storeRepository),
		Update:             newUseCaseUpdate(storeRepository),
		AddBusinessHour:    newUseCaseAddBusinessHour(storeRepository),
		DeleteBusinessHour: newUseCaseDeleteBusinessHour(storeRepository),
		GetByID:            newUseCaseGetByID(storeRepository),
		GetByFilter:        newUseCaseGetByFilter(storeRepository),
		SetProfileImage:    newUseCaseSetProfileImage(storeRepository),
		SetHeaderImage:     newUseCaseSetHeaderImage(storeRepository),
	}
}
