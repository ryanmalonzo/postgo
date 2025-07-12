# PostGO - Système de Typage Automatique

Ce document explique comment utiliser le système de génération de code typé de PostGO.

## 🚀 Vue d'ensemble

Le système de génération de code PostGO analyse automatiquement votre schéma de base de données et génère des structures Go typées avec :

- ✅ **Auto-complétion IDE** complète pour les noms de colonnes
- ✅ **Validation des types** à la compilation
- ✅ **Validation des contraintes** à l'exécution (NOT NULL, UNIQUE)
- ✅ **API fluide** pour l'insertion de données
- ✅ **Génération automatique** - s'adapte à tout nouveau schéma

## 📦 Installation et Configuration

### Prérequis

- Go 1.18+
- PostgreSQL
- Module PostGO initialisé

### Première génération

```bash
# Générer le code typé
make generate

# Ou directement avec go run
go run cmd/generate/main.go cmd/generate/generator.go -output=generated
```

## 🎯 Utilisation

### 1. Après génération, utilisation simple :

```go
package main

import (
    "postgo/db"
    "postgo/generated"
)

func main() {
    conn, err := db.NewConnection("localhost", 5432, "dbname", "user", "password")
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    // L'IDE vous propose automatiquement les colonnes disponibles !
    err = generated.Users.Insert().
        SetName("John Doe").          // string
        SetEmail("john@example.com"). // string avec contrainte UNIQUE
        SetPassword("secret123").     // string avec contrainte NOT NULL
        Execute(conn)

    if err != nil {
        // Gestion automatique des erreurs (contraintes, types, etc.)
        fmt.Printf("Erreur: %v\n", err)
    }
}
```

### 2. Types automatiquement détectés :

```go
// Pour la table companies
generated.Companies.Insert().
    SetName("Tech Corp").           // string (NOT NULL)
    SetDescription("...").          // string (nullable)
    SetEmployeeCount(150).          // int
    SetRevenue(1250000.50).         // float64
    SetIsPublic(true).              // bool (NOT NULL)
    Execute(conn)
```

### 3. Validation automatique :

```go
// ❌ Erreur de compilation - méthode inexistante
generated.Users.Insert().
    SetInvalidColumn("value") // L'IDE ne proposera pas cette méthode

// ❌ Erreur de compilation - type incorrect
generated.Users.Insert().
    SetEmployeeCount("string") // Attend un int, pas une string

// ❌ Erreur d'exécution - contrainte NOT NULL
generated.Users.Insert().
    SetName("John").
    // SetEmail("...") manquant - champ obligatoire
    Execute(conn) // Retourne une erreur
```

## 🛠 Commandes Make

```bash
# Génération
make generate     # Génère le code typé
make clean       # Supprime le code généré
make regen       # Nettoie et régénère

# Démonstrations
make demo-typed   # Démo du système typé (recommandé)
make demo-builder # Démo builder pattern classique
make demo-full    # Démo complète

# Test et construction
make test        # Teste la compilation
make build       # Construit le projet
make help        # Affiche l'aide
```

## 📁 Structure Générée

```
generated/
├── types.go       # Types et interfaces communs
├── tables.go      # Export des tables
├── users.go       # Structure typée pour users
├── companies.go   # Structure typée pour companies
├── posts.go       # Structure typée pour posts
└── categories.go  # Structure typée pour categories
```

## 🔄 Workflow de Développement

### 1. Modifier le schéma

Editez `db/schema.go` pour ajouter/modifier des tables :

```go
// Dans registerAllTables()
registerTable("products", createProductTable())

// Nouvelle fonction
func createProductTable() *TableBuilder {
    return NewTable("products").
        AddAttribute("name", String).NotNull().Build().
        AddAttribute("price", Float).NotNull().Build().
        AddAttribute("in_stock", Boolean).Build()
}
```

### 2. Régénérer automatiquement

```bash
make regen
```

### 3. Utiliser immédiatement

```go
// Le nouveau code est automatiquement disponible !
generated.Products.Insert().
    SetName("Laptop").     // Auto-complétion disponible
    SetPrice(999.99).      // Type float64 automatiquement détecté
    SetInStock(true).      // Type bool automatiquement détecté
    Execute(conn)
```

## 🔍 Fonctionnalités Avancées

### Types supportés

- `String` → `string`
- `Integer` → `int`
- `Float` → `float64`
- `Boolean` → `bool`
- `SERIAL` → `int` (auto-incrémenté, ignoré dans l'insertion)

### Contraintes validées

- **NOT NULL** : Validation à l'exécution si non défini
- **UNIQUE** : Géré par PostgreSQL
- **PRIMARY KEY** : ID auto-généré (ignoré)

### Prévention des erreurs

- ✅ Empêche la définition multiple d'une même colonne
- ✅ Valide que tous les champs NOT NULL sont définis
- ✅ Types vérifiés à la compilation
- ✅ Noms de colonnes vérifiés par l'IDE

## 🚨 Notes Importantes

1. **Code généré** : Ne jamais modifier manuellement les fichiers dans `generated/`
2. **Régénération** : Toujours régénérer après modification du schéma
3. **ID automatique** : L'ID est auto-généré et n'apparaît pas dans l'API
4. **Noms de colonnes** : Les snake_case sont automatiquement convertis en CamelCase

## 🎯 Exemple Complet

```go
package main

import (
    "fmt"
    "postgo/db"
    "postgo/generated"
)

func main() {
    // Connexion
    conn, _ := db.NewConnection("localhost", 5432, "postgo", "postgo", "postgo")
    defer conn.Close()

    // Initialisation du schéma
    db.InitAllTables(conn)

    // Insertions typées avec validation complète

    // Utilisateur complet
    err := generated.Users.Insert().
        SetName("Alice Dupont").
        SetEmail("alice@example.com").
        SetPassword("motdepasse123").
        Execute(conn)
    fmt.Printf("Utilisateur: %v\n", err)

    // Entreprise avec champs optionnels
    err = generated.Companies.Insert().
        SetName("InnovTech").
        SetEmployeeCount(50).
        SetIsPublic(false).
        // Description et Revenue optionnels
        Execute(conn)
    fmt.Printf("Entreprise: %v\n", err)

    // Post de blog
    err = generated.Posts.Insert().
        SetTitle("Guide PostGO").
        SetContent("Comment utiliser le système de typage...").
        SetPublished(true).
        Execute(conn)
    fmt.Printf("Post: %v\n", err)
}
```

Voilà ! Vous avez maintenant un système de typage automatique qui s'adapte à tout schéma et fournit une expérience de développement optimale avec validation complète ! 🎉
