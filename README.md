# postgo - ORM PostgreSQL avec système de typage automatique

Un ORM PostgreSQL simple et léger offrant une **Developer Experience optimale** grâce à un système de génération de code qui fournit autocomplétion complète, validation des types et sécurité à la compilation.

## Installation

```bash
go mod tidy
```

## Démarrage rapide

### 1. Définir votre schéma dans `db/schema.go`

```go
func createUserTable() *TableBuilder {
    return NewTable("users").
        AddAttribute("name", String).NotNull().Build().
        AddAttribute("email", String).NotNull().Unique().Build().
        AddAttribute("password", String).NotNull().Build()
}
```

### 2. Générer le code typé

```bash
make generate
```

### 3. Utiliser avec autocomplétion complète

```go
package main

import (
    "postgo/db"
    "postgo/generated"
)

func main() {
    conn, err := db.NewConnection("localhost", 5432, "postgo", "postgo", "postgo")
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    // L'IDE propose automatiquement tous les noms de colonnes !
    err = generated.Users.Insert().
        SetName("John Doe").          // string
        SetEmail("john@example.com"). // string avec contrainte UNIQUE  
        SetPassword("secret123").     // string avec contrainte NOT NULL
        Execute(conn)
    
    if err != nil {
        fmt.Printf("Erreur: %v\n", err)
    }
}
```

## Utilisation

### Définition du schéma

Toutes les tables sont définies dans `db/schema.go` avec le builder pattern :

```go
// Table avec différents types de données
func createCompanyTable() *TableBuilder {
    return NewTable("companies").
        AddAttribute("name", String).NotNull().Unique().Build().
        AddAttribute("employee_count", Integer).Build().
        AddAttribute("revenue", Float).Build().
        AddAttribute("is_public", Boolean).NotNull().Build()
}
```

#### Types de données disponibles

- `String` - VARCHAR(255)
- `Integer` - INTEGER  
- `Float` - FLOAT
- `Boolean` - BOOLEAN

#### Contraintes disponibles

- `.NotNull()` - Ajoute NOT NULL
- `.Unique()` - Ajoute UNIQUE

### Génération et utilisation du code

```bash
# Générer le code typé
make generate

# Nettoyer et régénérer
make regen

# Tester la compilation
make test
```

Le code généré fournit :

- **Autocomplétion IDE complète** pour tous les noms de colonnes
- **Validation des types à la compilation** (impossible de passer un `int` à une colonne `string`)
- **Validation des contraintes à l'exécution** (champs NOT NULL obligatoires)
- **Prévention des erreurs** de frappe et d'incohérences

```go
// Types automatiquement détectés
generated.Companies.Insert().
    SetName("Tech Corp").           // string (NOT NULL)
    SetEmployeeCount(150).          // int
    SetRevenue(1250000.50).         // float64
    SetIsPublic(true).              // bool (NOT NULL)
    Execute(conn)

// ❌ Erreurs détectées à la compilation
generated.Users.Insert().
    SetEmployeeCount("string")      // Type incorrect
    SetInvalidColumn("value")       // Colonne inexistante
```

## Architecture

### Composants principaux

- **Schéma centralisé** (`db/schema.go`) : Définition de toutes les tables
- **Générateur de code** (`cmd/generate/`) : Analyse le schéma et génère le code Go typé
- **Code généré** (`generated/`) : Structures typées avec autocomplétion complète
- **Connection** : Gestionnaire de connexion PostgreSQL

### SQL généré

Le système génère automatiquement :

- Un ID SERIAL PRIMARY KEY pour chaque table
- Les définitions de colonnes avec leurs types
- Les contraintes NOT NULL et UNIQUE

Exemple de SQL généré :

```sql
CREATE TABLE IF NOT EXISTS "users" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "password" VARCHAR(255) NOT NULL
)
```

## Workflow de développement

Le workflow est optimisé pour la productivité :

1. **Modifier le schéma** dans `db/schema.go`
2. **Régénérer** avec `make generate` 
3. **Utiliser immédiatement** avec autocomplétion complète

```bash
# Cycle de développement
make regen    # Nettoie, régénère et teste
```

## Démonstration

```bash
# Démo complète avec le système typé
go run . -demo=typed
```

## Test avec Docker

Pour tester rapidement avec PostgreSQL :

```bash
docker compose up -d
```

Cela démarre PostgreSQL et Adminer sur http://localhost:8080

## Philosophie

postgo est volontairement simple tout en offrant une **Developer Experience moderne** :

- ✅ **ID auto-incrémenté obligatoire** pour chaque table
- ✅ **Types de base** (String, Integer, Float, Boolean) 
- ✅ **Contraintes essentielles** (NOT NULL, UNIQUE)
- ✅ **Autocomplétion complète** grâce au code généré
- ✅ **Validation à la compilation** pour éviter les erreurs
- ✅ **Simplicité d'usage** avec API intuitive

Limitations volontaires :
- ❌ Pas de foreign keys
- ❌ Pas d'index personnalisés  
- ❌ Pas de relations complexes

L'objectif est de fournir un outil **simple, sûr et productif** pour des cas d'usage basiques avec la meilleure expérience développeur possible.
