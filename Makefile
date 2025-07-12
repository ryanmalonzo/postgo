# Makefile pour PostGO

.PHONY: generate clean build run demo-typed demo-builder demo-full test help

# Génère le code typé automatiquement
generate:
	@echo "Génération du code typé..."
	@go run cmd/generate/main.go cmd/generate/generator.go -output=generated
	@echo "✓ Code généré avec succès!"

# Nettoie le code généré
clean:
	@echo "Nettoyage du code généré..."
	@rm -rf generated/
	@echo "✓ Code généré supprimé!"

# Construit le projet
build:
	@echo "Construction du projet..."
	@go build -o postgo main.go
	@echo "✓ Projet construit!"

# Lance la démo basique
run:
	@go run main.go

# Lance la démo avec le système typé
demo-typed: generate
	@echo "Lancement de la démo typée..."
	@go run main.go -demo=typed

# Lance la démo avec le builder pattern
demo-builder:
	@go run main.go -demo=builder

# Lance la démo complète
demo-full:
	@go run main.go -demo=full

# Teste la compilation du code généré
test: generate
	@echo "Test de compilation du code généré..."
	@go build ./generated/...
	@echo "✓ Code généré compile correctement!"

# Régénère et test
regen: clean generate test
	@echo "✓ Régénération complète terminée!"

# Affiche l'aide
help:
	@echo "Commandes disponibles:"
	@echo "  generate     - Génère le code typé à partir du schéma"
	@echo "  clean        - Supprime le code généré"
	@echo "  build        - Construit le projet"
	@echo "  run          - Lance la démo basique"
	@echo "  demo-typed   - Lance la démo avec le système typé"
	@echo "  demo-builder - Lance la démo avec le builder pattern"
	@echo "  demo-full    - Lance la démo complète"
	@echo "  test         - Teste la compilation du code généré"
	@echo "  regen        - Nettoie, régénère et teste"
	@echo "  help         - Affiche cette aide"

# Par défaut affiche l'aide
.DEFAULT_GOAL := help
