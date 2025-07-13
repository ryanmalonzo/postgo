package examples

import (
	"fmt"
	"log"
	"postgo/db"
	"postgo/generated"
	"postgo/logging"

	_ "github.com/lib/pq"
)

// DemoTypedInserts démontre l'utilisation du système de typage généré pour Insert, Update et Delete
func DemoTypedInserts() {
	fmt.Println("=== Démonstration du système de typage généré (Insert, Update & Delete) ===")

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
	testDuplicationPrevention()

	// === EXEMPLES DE DELETE ===
	
	// Insérer des données spécifiques pour les tests de suppression
	fmt.Println("\n--- Préparation des données pour les tests Delete ---")
	
	// Insertion d'un utilisateur pour test de suppression
	err = generated.Users.Insert().
		SetName("User To Delete").
		SetEmail("delete.me@example.com").
		SetPassword("password123").
		Execute(conn)
	
	if err != nil {
		fmt.Printf("Note: Erreur lors de l'insertion du test user: %v\n", err)
	}

	// Insertion d'un post pour test de suppression
	err = generated.Posts.Insert().
		SetTitle("Post à supprimer").
		SetContent("Ce post sera supprimé dans les tests.").
		SetPublished(false).
		Execute(conn)
	
	if err != nil {
		fmt.Printf("Note: Erreur lors de l'insertion du test post: %v\n", err)
	}

	// 11. Suppression d'un utilisateur spécifique
	fmt.Println("\n--- Suppression d'un utilisateur ---")
	err = generated.Users.Delete().
		Where("email = 'delete.me@example.com'").
		Execute(conn)
	
	if err != nil {
		fmt.Printf("Erreur lors de la suppression: %v\n", err)
	} else {
		fmt.Println("✓ Utilisateur supprimé avec succès!")
	}

	// 12. Suppression d'une entreprise par nom (plus sûr que par ID)
	fmt.Println("\n--- Suppression d'une entreprise ---")
	err = generated.Companies.Delete().
		Where("name = 'Tech Corp'").
		Execute(conn)
	
	if err != nil {
		fmt.Printf("Erreur lors de la suppression: %v\n", err)
	} else {
		fmt.Println("✓ Entreprise supprimée avec succès!")
	}

	// 13. Suppression avec conditions multiples
	fmt.Println("\n--- Suppression avec conditions multiples ---")
	err = generated.Posts.Delete().
		Where("published = false").
		Where("title LIKE '%supprimer%'").
		Execute(conn)
	
	if err != nil {
		fmt.Printf("Erreur lors de la suppression: %v\n", err)
	} else {
		fmt.Println("✓ Posts supprimés avec succès!")
	}

	// 14. Suppression de catégories par slug
	fmt.Println("\n--- Suppression de catégories ---")
	err = generated.Categories.Delete().
		Where("slug = 'technology'").
		Execute(conn)
	
	if err != nil {
		fmt.Printf("Erreur lors de la suppression: %v\n", err)
	} else {
		fmt.Println("✓ Catégorie supprimée avec succès!")
	}

	// 15. Test de construction de requête sans exécution
	fmt.Println("\n--- Test de construction de requête DELETE ---")
	sqlQuery, args := generated.Users.Delete().
		Where("email LIKE '%@example.com'").
		Build()
	
	fmt.Printf("Requête SQL générée: %s\n", sqlQuery)
	fmt.Printf("Arguments: %v\n", args)

	// 16. Test de validation - DELETE sans condition WHERE
	fmt.Println("\n--- Test de validation DELETE sans WHERE ---")
	sqlQueryNoWhere, argsNoWhere := generated.Users.Delete().Build()
	fmt.Printf("Requête sans WHERE: %s\n", sqlQueryNoWhere)
	fmt.Printf("Arguments sans WHERE: %v\n", argsNoWhere)
	fmt.Println("⚠️  Attention: Cette requête supprimerait tous les utilisateurs!")

	// === EXEMPLES DE SELECT ===

	fmt.Println("\n--- Insertion d'un utilisateur ---")
	err = generated.Users.Insert().
		SetName("Alice Doe").
		SetEmail("alice.doe@example.com").
		SetPassword("securepassword123").
		Execute(conn)
	
	if err != nil {
		fmt.Printf("Erreur lors de l'insertion: %v\n", err)
	} else {
		fmt.Println("✓ Utilisateur inséré avec succès!")
	}

	// 17. SelectAll - Récupérer tous les utilisateurs
	fmt.Println("\n1. Test SelectAll:")
	users, err := generated.Users.Select().
	SelectAll().
	Execute(conn)
	if err != nil {
		log.Printf("Erreur lors du SelectAll: %v", err)
	} else {
		fmt.Printf("Nombre d'utilisateurs trouvés: %d\n", len(users))
		for _, user := range users {
			fmt.Printf("  - ID: %d, Name: %s, Email: %s\n", user.Id, user.Name, user.Email)
		}
	}

	// Test 2: Select avec colonnes spécifiques
	fmt.Println("\n2. Test Select avec colonnes spécifiques:")
	users, err = generated.Users.Select().
	SelectColumns("name", "email").
	Execute(conn)
	if err != nil {
		log.Printf("Erreur lors du Select avec colonnes: %v", err)
	} else {
		fmt.Printf("Nombre d'utilisateurs (name, email) trouvés: %d\n", len(users))
		for _, user := range users {
			fmt.Printf("  - Name: %s, Email: %s\n", user.Name, user.Email)
		}
	}

	// Test 3: Select avec condition WHERE
	fmt.Println("\n3. Test Select avec WHERE:")
	users, err = generated.Users.Select().
	SelectAll().
	Where("name LIKE '%John%'").
	Execute(conn)
	if err != nil {
		log.Printf("Erreur lors du Select avec WHERE: %v", err)
	} else {
		fmt.Printf("Nombre d'utilisateurs avec 'John' dans le nom: %d\n", len(users))
		for _, user := range users {
			fmt.Printf("  - ID: %d, Name: %s, Email: %s\n", user.Id, user.Name, user.Email)
		}
	}

	// Test 4: Select avec WHERE typé
	fmt.Println("\n4. Test Select avec WHERE typé:")
	users, err = generated.Users.Select().
	SelectAll().
	WhereName("Alice Doe").
	Execute(conn)
	if err != nil {
		log.Printf("Erreur lors du Select avec WHERE typé: %v", err)
	} else {
		fmt.Printf("Nombre d'utilisateurs nommés 'Alice Doe': %d\n", len(users))
		for _, user := range users {
			fmt.Printf("  - ID: %d, Name: %s, Email: %s\n", user.Id, user.Name, user.Email)
		}
	}

	// Test 5: ExecuteOne - Récupérer un seul utilisateur
	fmt.Println("\n5. Test récupérer user by ID:")
	user, err := generated.Users.Select().
	SelectAll().
	WhereId(1).
	ExecuteOne(conn)
	if err != nil {
		log.Printf("Erreur lors du GetById: %v", err)
	} else {
		fmt.Printf("Utilisateur avec ID n°1 => ID: %d, Name: %s, Email: %s\n", user.Id, user.Name, user.Email)
	}

	fmt.Println("\n=== Démonstration terminée ===")
}

// testDuplicationPrevention teste la prévention de duplication de colonnes
func testDuplicationPrevention() {
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
}
