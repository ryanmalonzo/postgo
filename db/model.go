package db

// Interface de base à implémenter par tous les modèles (conservée pour compatibilité)
type Model interface {
	TableName() string
}

// BaseModel struct de base avec ID auto-incrémenté (conservé pour compatibilité)
type BaseModel struct {
	ID int64
}
