package examples

import (
	"fmt"
	"postgo/db"
	"postgo/generated"
	"postgo/logging"

	_ "github.com/lib/pq"
)

// DemoTypedInserts démontre l'utilisation du système de typage généré
func DemoTypedInserts() {
	fmt.Println("=== Démonstration du système de typage généré ===")

	// Connexion à la base de données
	conn, err := db.NewConnection("localhost", 5432, "postgo", "postgo", "postgo")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Initialisation automatique de toutes les tables du schéma
	err = db.InitAllTables(conn)
	if err != nil {
		panic(err)
	}
	logging.Info.Println("All schema tables initialized successfully!")

	// Exemples d'utilisation avec le système typé généré

	// 1. Insertion dans la table users avec auto-complétion
	fmt.Println("\n--- Insertion d'un utilisateur ---")
	err = generated.Users.Insert().
		SetName("John Doe").
		SetEmail("john.doe@example.com").
		SetPassword("securepassword123").
		Execute(conn)
	
	if err != nil {
		fmt.Printf("Erreur lors de l'insertion: %v\n", err)
	} else {
		fmt.Println("✓ Utilisateur inséré avec succès!")
	}

	// 2. Insertion dans la table companies
	fmt.Println("\n--- Insertion d'une entreprise ---")
	err = generated.Companies.Insert().
		SetName("Tech Corp").
		SetDescription("Une entreprise technologique innovante").
		SetEmployeeCount(150).
		SetRevenue(1250000.50).
		SetIsPublic(true).
		Execute(conn)
	
	if err != nil {
		fmt.Printf("Erreur lors de l'insertion: %v\n", err)
	} else {
		fmt.Println("✓ Entreprise insérée avec succès!")
	}

	// 3. Insertion dans la table posts
	fmt.Println("\n--- Insertion d'un post ---")
	err = generated.Posts.Insert().
		SetTitle("Mon premier article").
		SetContent("Ceci est le contenu de mon premier article de blog.").
		SetPublished(true).
		Execute(conn)
	
	if err != nil {
		fmt.Printf("Erreur lors de l'insertion: %v\n", err)
	} else {
		fmt.Println("✓ Post inséré avec succès!")
	}

	// 4. Insertion dans la table categories
	fmt.Println("\n--- Insertion d'une catégorie ---")
	err = generated.Categories.Insert().
		SetSlug("technology").
		SetDisplayName("Technologie").
		Execute(conn)
	
	if err != nil {
		fmt.Printf("Erreur lors de l'insertion: %v\n", err)
	} else {
		fmt.Println("✓ Catégorie insérée avec succès!")
	}

	// 5. Exemple de validation des champs obligatoires
	fmt.Println("\n--- Test de validation des champs obligatoires ---")
	err = generated.Users.Insert().
		SetName("Jane Doe").
		// Oubli volontaire de l'email (champ obligatoire)
		SetPassword("password123").
		Execute(conn)
	
	if err != nil {
		fmt.Printf("✓ Validation réussie - Erreur attendue: %v\n", err)
	} else {
		fmt.Println("❌ La validation a échoué - aucune erreur détectée")
	}

	fmt.Println("\n=== Démonstration terminée ===")
}
