# Changements apportés à l'ORM PostGo

## Résumé

Le code source de l'ORM PostgreSQL a été complètement réécrit pour utiliser le **builder pattern** au lieu de la réflexion et des tags de struct.

## Avant (tags et réflexion)

```go
type User struct {
    db.BaseModel
    Name     string `db:"not_null"`
    Email    string `db:"not_null,unique"`
    Password string `db:"not_null"`
}

func (u *User) TableName() string {
    return "users"
}

// Utilisation
err = conn.CreateTable(&User{})
```

## Après (builder pattern)

```go
// Définition avec builder pattern
userTable := db.NewTable("users").
    AddAttribute("name", db.String).NotNull().Build().
    AddAttribute("email", db.String).NotNull().Unique().Build().
    AddAttribute("password", db.String).NotNull().Build()

// Utilisation
err = conn.CreateTable(userTable)
```

## Fichiers modifiés

### `/db/table.go` - Complètement réécrit

- Suppression de la logique de réflexion
- Ajout des types `AttributeType`, `Attribute`, `AttributeBuilder`, `TableBuilder`
- Nouvelle méthode `CreateTable(tableBuilder *TableBuilder)`
- Builder pattern fluide avec `.AddAttribute().NotNull().Build()`

### `/db/model.go` - Simplifié

- Conservation de `Model` interface et `BaseModel` pour compatibilité
- Suppression de toute la logique de réflexion (`GetMetadata`, `extractFields`, etc.)
- Plus de tags `db` nécessaires

### `/main.go` - Adapté au nouveau pattern

- Utilisation de `examples.CreateUserTable()` au lieu de `&examples.User{}`

### `/examples/user_table.go` - Nouveau fichier

- Définition de la table users avec le builder pattern
- Remplacement de la struct avec tags

### `/examples/advanced_tables.go` - Nouveau fichier

- Exemples de tables plus complexes (companies, posts, categories)

## Fonctionnalités conservées

✅ **ID auto-incrémenté obligatoire** - Chaque table a automatiquement un `id SERIAL PRIMARY KEY`
✅ **Types supportés** - String, Integer, Float, Boolean (mêmes mappings SQL)
✅ **Contraintes** - NOT NULL, UNIQUE
✅ **Simplicité** - Pas de foreign keys, pas d'index (comme demandé)

## Nouvelles fonctionnalités (syntaxe uniquement)

✅ **Builder pattern fluide** - Syntaxe `.AddAttribute("name", db.String).NotNull().Build()`
✅ **Types explicites** - `db.String`, `db.Integer`, `db.Float`, `db.Boolean`
✅ **Chaînage de contraintes** - `.NotNull().Unique().Build()`

## Exemples de migration

### Table simple

```go
// Avant
type Post struct {
    db.BaseModel
    Title     string `db:"not_null"`
    Content   string
    Published bool
}

// Après
postTable := db.NewTable("posts").
    AddAttribute("title", db.String).NotNull().Build().
    AddAttribute("content", db.String).Build().
    AddAttribute("published", db.Boolean).Build()
```

### Table avec contraintes multiples

```go
// Avant
type Category struct {
    db.BaseModel
    Slug        string `db:"not_null,unique"`
    DisplayName string `db:"not_null"`
}

// Après
categoryTable := db.NewTable("categories").
    AddAttribute("slug", db.String).NotNull().Unique().Build().
    AddAttribute("display_name", db.String).NotNull().Build()
```

## SQL généré (identique)

Le SQL généré reste exactement le même :

```sql
CREATE TABLE IF NOT EXISTS "users" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "password" VARCHAR(255) NOT NULL
)
```

## Avantages du nouveau système

1. **Syntaxe plus claire** - `.AddAttribute("name", db.String).NotNull()` vs `Name string \`db:"not_null"\``
2. **Pas de réflexion** - Plus simple, plus prévisible
3. **Types explicites** - `db.String` vs inférence depuis `string`
4. **Fluide** - Chaînage naturel des méthodes
5. **Moins magique** - Plus explicite sur ce qui se passe

## Tests

Tous les exemples ont été testés et fonctionnent :

- `go run main.go` - ✅
- `go run examples/demo/example_builder.go` - ✅
- Compilation sans erreurs - ✅
- SQL généré correct - ✅
