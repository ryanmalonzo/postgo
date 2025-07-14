package main

import (
	"flag"
	"fmt"
	"os"
	"postgo/db"
)

func main() {
	var outputDir = flag.String("output", "generated", "Répertoire de sortie pour les fichiers générés")
	flag.Parse()

	fmt.Println("=== Générateur de code PostGO ===")
	
	// Créer le répertoire de sortie s'il n'existe pas
	err := os.MkdirAll(*outputDir, 0755)
	if err != nil {
		panic(fmt.Errorf("impossible de créer le répertoire %s: %v", *outputDir, err))
	}

	// Obtenir toutes les tables du schéma
	tables := db.GetAllTables()
	if len(tables) == 0 {
		fmt.Println("Aucune table trouvée dans le schéma")
		return
	}

	fmt.Printf("Génération du code pour %d table(s)...\n", len(tables))

	// Générer le fichier principal avec les types et constantes
	err = generateMainTypes(*outputDir)
	if err != nil {
		panic(fmt.Errorf("erreur lors de la génération des types: %v", err))
	}

	// Générer un fichier pour chaque table
	for tableName, table := range tables {
		err = generateTableFile(*outputDir, tableName, table)
		if err != nil {
			panic(fmt.Errorf("erreur lors de la génération de la table %s: %v", tableName, err))
		}
		fmt.Printf("✓ Table '%s' générée\n", tableName)
	}

	fmt.Printf("✓ Génération terminée dans le répertoire '%s'\n", *outputDir)
}
