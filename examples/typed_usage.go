package examples

import (
	"fmt"
	"postgo/db"
	"postgo/generated"
	"postgo/logging"

	_ "github.com/lib/pq"
)

// DemoTypedInserts démontre l'utilisation du système de typage généré pour Insert et Update
func DemoTypedInserts() {
	fmt.Println("=== Démonstration du système de typage généré (Insert & Update) ===")

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

	// === EXEMPLES D'UPDATE ===
	
	// 6. Update d'un utilisateur
	fmt.Println("\n--- Update d'un utilisateur ---")
	err = generated.Users.Update().
		SetName("John Doe Updated").
		SetEmail("john.updated@example.com").
		Where("name = 'John Doe'").
		Execute(conn)
	
	if err != nil {
		fmt.Printf("Erreur lors de l'update: %v\n", err)
	} else {
		fmt.Println("✓ Utilisateur mis à jour avec succès!")
	}

	// 7. Update d'une entreprise (colonnes optionnelles)
	fmt.Println("\n--- Update d'une entreprise ---")
	err = generated.Companies.Update().
		SetEmployeeCount(200).
		SetRevenue(2500000.75).
		Where("name = 'Tech Corp'").
		Execute(conn)
	
	if err != nil {
		fmt.Printf("Erreur lors de l'update: %v\n", err)
	} else {
		fmt.Println("✓ Entreprise mise à jour avec succès!")
	}

	// 8. Update d'une seule colonne
	fmt.Println("\n--- Update d'une seule colonne ---")
	err = generated.Posts.Update().
		SetPublished(false).
		Where("title = 'Mon premier article'").
		Execute(conn)
	
	if err != nil {
		fmt.Printf("Erreur lors de l'update: %v\n", err)
	} else {
		fmt.Println("✓ Post mis à jour avec succès!")
	}

	// 9. Test de validation "aucune colonne à mettre à jour"
	fmt.Println("\n--- Test de validation des updates vides ---")
	err = generated.Users.Update().
		Where("id = 1").
		Execute(conn)
	
	if err != nil {
		fmt.Printf("✓ Validation réussie - Erreur attendue: %v\n", err)
	} else {
		fmt.Println("❌ La validation a échoué - aucune erreur détectée")
	}

	// 10. Test de prévention de duplication de colonnes
	fmt.Println("\n--- Test de prévention de duplication ---")
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("✓ Validation réussie - Panic attendu: %v\n", r)
		}
	}()
	
	generated.Categories.Update().
		SetSlug("tech").
		SetSlug("technology"). // Tentative de redéfinir la même colonne
		Where("id = 1")
	
	fmt.Println("❌ La validation a échoué - aucun panic détecté")

	fmt.Println("\n=== Démonstration terminée ===")
}
