# PostGo - Simple PostgreSQL ORM avec Builder Pattern

Un ORM PostgreSQL simple et léger utilisant le builder pattern pour la définition de tables.

## Installation

```bash
go mod tidy
```

## Usage

The project includes multiple ways to run examples:

### Basic Demo (default)
```bash
go run .
```

### Builder Pattern Examples
```bash
go run . -demo=builder
```

### Full Comprehensive Demo
```bash
go run . -demo=full
```

## Utilisation de base

### Création d'une table simple

```go
package main

import (
    "postgo/db"
    _ "github.com/lib/pq"
)

func main() {
    // Connexion à la base
    conn, err := db.NewConnection("localhost", 5432, "user", "password", "database")
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    // Création d'une table avec le builder pattern
    userTable := db.NewTable("users").
        AddAttribute("name", db.String).NotNull().Build().
        AddAttribute("email", db.String).NotNull().Unique().Build().
        AddAttribute("password", db.String).NotNull().Build()

    // Création de la table dans la base
    err = conn.CreateTable(userTable)
    if err != nil {
        panic(err)
    }
}
```

### Types de données disponibles

- `db.String` - VARCHAR(255)
- `db.Integer` - INTEGER
- `db.Float` - FLOAT
- `db.Boolean` - BOOLEAN

### Contraintes disponibles

- `.NotNull()` - Ajoute NOT NULL
- `.Unique()` - Ajoute UNIQUE

### Exemples d'utilisation

```go
// Table avec différents types de données
productTable := db.NewTable("products").
    AddAttribute("name", db.String).NotNull().Build().
    AddAttribute("price", db.Float).NotNull().Build().
    AddAttribute("in_stock", db.Boolean).Build().
    AddAttribute("quantity", db.Integer).Build()

// Table avec contraintes multiples
categoryTable := db.NewTable("categories").
    AddAttribute("slug", db.String).NotNull().Unique().Build().
    AddAttribute("display_name", db.String).NotNull().Build()

// Table minimale (seulement l'ID auto-incrémenté)
simpleTable := db.NewTable("logs")
```

## Architecture

### Composants principaux

- **TableBuilder** : Constructeur de table avec le pattern builder
- **AttributeBuilder** : Constructeur d'attributs avec contraintes
- **Connection** : Gestionnaire de connexion PostgreSQL

### SQL généré

Le builder génère automatiquement :

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

## Test avec Docker

Pour tester rapidement avec PostgreSQL :

```bash
docker-compose up -d
```

Cela démarre PostgreSQL et Adminer sur http://localhost:8080

## Exemples

Voir les fichiers d'exemple :

- `main.go` - Exemple de base
- `example_builder.go` - Démonstration complète du builder pattern
- `examples/` - Définitions de tables d'exemple

## Philosophie

Cet ORM est volontairement simple :

- ✅ ID auto-incrémenté obligatoire
- ✅ Types de base (String, Integer, Float, Boolean)
- ✅ Contraintes essentielles (NOT NULL, UNIQUE)
- ❌ Pas de foreign keys
- ❌ Pas d'index personnalisés
- ❌ Pas de relations complexes

L'objectif est de fournir un outil simple et prévisible pour des cas d'usage basiques.
