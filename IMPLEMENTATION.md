# Pistes de réflexion et justification des choix techniques

## 1. Première implémentation : réflexion et tags de struct

### 1.1 Principe et fonctionnement

La première approche utilisait la réflexion Go native pour analyser des structures avec des tags, similaire à d'autres ORMs comme GORM ou Ent.

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

Le système analysait les tags `db:` pour extraire les contraintes et générait automatiquement le SQL correspondant via la réflexion.

### 1.2 Avantages pour la DX

- **Familiarité** : Approche similaire aux ORMs Go populaires
- **Compacité** : Définition de table en une seule struct
- **Lisibilité** : Structure claire avec types Go natifs
- **Convention** : Pattern établi dans l'écosystème Go

### 1.3 Inconvénients pour la DX

- **Autocomplétion limitée** : L'IDE ne peut pas proposer les noms de colonnes lors de l'écriture de requêtes
- **Validation à l'exécution uniquement** : Les erreurs de noms de colonnes ne sont détectées qu'au runtime
- **Tags magiques** : La syntaxe `db:"not_null,unique"` n'est pas vérifiée à la compilation
- **Réflexion opaque** : Comportement difficile à déboguer et à comprendre
- **Performance** : Coût de la réflexion à l'exécution

### 1.4 Raisons de l'abandon

Le problème principal était l'**absence d'autocomplétion intelligente**. Lors de l'écriture de requêtes, les développeurs devaient :

- Mémoriser les noms de colonnes
- Faire des allers-retours constants avec la définition de struct
- Découvrir les erreurs de frappe uniquement à l'exécution

Cette approche ne respectait pas l'objectif de DX optimale du projet.

## 2. Deuxième implémentation : query builder pattern

### 2.1 Principe et fonctionnement

La deuxième approche conservait les structs pour la définition de schéma mais introduisait un système de query builder fluide pour les opérations CRUD.

```go
// Requête SELECT avec builder pattern
selectQuery := query.NewSelectQuery("users").
    AddColumn("id").
    AddColumn("name").
    Where("id = 1")

rows, err := selectQuery.Execute(conn.GetDatabase())

// Requête INSERT avec builder pattern
insertQuery := query.NewInsertQuery("users").
    AddColumn("name").AddValue("John").
    AddColumn("email").AddValue("john@example.com").
    AddColumn("password").AddValue("secret123")

_, err = insertQuery.Execute(conn.GetDatabase())
```

Le système proposait des builders séparés pour chaque type d'opération (SELECT, INSERT, UPDATE, DELETE) avec une API fluide.

### 2.2 Avantages pour la DX

- **API fluide** : Syntaxe naturelle et lisible
- **Flexibilité** : Construction de requêtes complexes étape par étape
- **Séparation des responsabilités** : Schéma vs requêtes
- **Debugging** : SQL généré visible et compréhensible

### 2.3 Inconvénients pour la DX

- **Pas d'autocomplétion pour les colonnes** : Les noms de colonnes restaient des chaînes de caractères
- **Validation manuelle** : Aucune vérification de l'existence des colonnes
- **Erreurs à l'exécution** : Fautes de frappe détectées tardivement
- **Verbosité** : Code plus long pour des opérations simples
- **Cohérence** : Risque de divergence entre définition de schéma et requêtes

### 2.4 Pourquoi cette approche était insuffisante

Bien que plus flexible, cette approche ne résolvait pas le problème fondamental de **l'absence de validation des types et des noms de colonnes à la compilation**. Les développeurs étaient toujours obligés de mémoriser les noms de colonnes et découvraient les erreurs uniquement à l'exécution.

## 3. Implémentation actuelle : schema-based avec génération de code

### 3.1 Principe et fonctionnement

L'approche actuelle s'inspire du modèle Prisma : un fichier de schéma centralisé (`db/schema.go`) et un générateur qui produit du code Go typé avec autocomplétion complète.

```go
// Définition dans schema.go
func createUserTable() *TableBuilder {
    return NewTable("users").
        AddAttribute("name", String).NotNull().Build().
        AddAttribute("email", String).NotNull().Unique().Build().
        AddAttribute("password", String).NotNull().Build()
}

// Code généré automatiquement
generated.Users.Insert().
    SetName("John Doe").          // string avec autocomplétion
    SetEmail("john@example.com"). // string avec contrainte UNIQUE
    SetPassword("secret123").     // string avec contrainte NOT NULL
    Execute(conn)
```

### 3.2 Avantages pour la DX

- **Autocomplétion complète** : L'IDE propose automatiquement tous les noms de colonnes
- **Validation des types à la compilation** : Impossible de passer un `int` à une colonne `string`
- **Validation des contraintes** : Vérification automatique des champs NOT NULL
- **Feedback immédiat** : Les erreurs sont détectées pendant l'écriture du code
- **API intuitive** : Méthodes `SetColumnName()` générées automatiquement
- **Sécurité** : Prévention des erreurs de frappe et des incohérences
- **Régénération automatique** : Le code s'adapte aux changements de schéma

### 3.3 Architecture détaillée

Le système se compose de plusieurs éléments :

**Schéma centralisé (`db/schema.go`)**

```go
func registerAllTables() {
    registerTable("users", createUserTable())
    registerTable("companies", createCompanyTable())
    // ...
}
```

**Générateur de code (`cmd/generate/`)**

- Analyse le schéma défini
- Génère des structures Go typées
- Crée des méthodes Set pour chaque colonne
- Ajoute la validation des contraintes

**Code généré (`generated/`)**

```go
type UsersTable struct {
    Name string
}

var Users = &UsersTable{Name: "users"}

func (t *UsersTable) Insert() *UsersInsertBuilder {
    return &UsersInsertBuilder{
        query: query.NewInsertQuery("users"),
    }
}

func (b *UsersInsertBuilder) SetName(value string) *UsersInsertBuilder {
    // ...
}
```

### 3.4 Workflow de développement

Le workflow est optimisé pour la productivité :

1. **Modification du schéma** dans `db/schema.go`
2. **Régénération automatique** avec `make generate`
3. **Utilisation immédiate** avec autocomplétion complète

```bash
# Régénération simple
make generate

# Régénération + test de compilation
make regen
```

## 4. Comparaison des approches du point de vue DX

### 4.1 Tableau comparatif

| Critère DX                 | Réflexion + Tags | Query Builder        | Schema + Génération      |
| -------------------------- | ---------------- | -------------------- | ------------------------ |
| **Autocomplétion IDE**     | ❌ Aucune        | ❌ Aucune            | ✅ Complète              |
| **Validation types**       | ❌ Runtime       | ❌ Runtime           | ✅ Compilation           |
| **Validation contraintes** | ❌ Runtime       | ❌ Aucune            | ✅ Runtime + Compilation |
| **Détection erreurs**      | ❌ Tardive       | ❌ Tardive           | ✅ Immédiate             |
| **Simplicité API**         | ✅ Simple        | ⚠️ Verbeux           | ✅ Intuitive             |
| **Maintenabilité**         | ⚠️ Dispersé      | ⚠️ Risque divergence | ✅ Centralisé            |
| **Courbe d'apprentissage** | ✅ Familière     | ⚠️ Moyenne           | ✅ Faible                |

### 4.2 Critères d'évaluation DX

**Feedback immédiat** : L'approche actuelle permet de détecter les erreurs dès l'écriture du code, avant même la compilation.

**Productivité** : L'autocomplétion élimine les allers-retours avec la documentation et réduit les erreurs de frappe.

**Confiance** : La validation des types garantit que le code qui compile fonctionnera à l'exécution.

**Maintenabilité** : Le schéma centralisé facilite les modifications et assure la cohérence.

## 5. Justification de l'implémentation actuelle

### 5.1 Autocomplétion IDE

L'autocomplétion est le pilier de l'expérience développeur moderne. Avec l'approche actuelle :

```go
// L'IDE propose automatiquement : SetName, SetEmail, SetPassword
generated.Users.Insert().Set...
```

Cette fonctionnalité transforme l'écriture de code d'un exercice de mémorisation en une exploration guidée.

### 5.2 Validation des types à la compilation

```go
// ❌ Erreur de compilation - impossible
generated.Users.Insert().SetEmployeeCount("string") // Attend un int

// ✅ Validation automatique des types
generated.Companies.Insert().SetEmployeeCount(150) // Correct
```

### 5.3 Validation des contraintes à l'exécution

```go
// ❌ Erreur explicite si champ obligatoire manquant
err := generated.Users.Insert().
    SetName("John").
    // SetEmail manquant - champ NOT NULL
    Execute(conn)
// Erreur : "la colonne obligatoire 'email' n'a pas été définie"
```

### 5.4 Simplicité d'utilisation

L'API générée est plus simple que les approches précédentes :

```go
// Avant (query builder)
query.NewInsertQuery("users").
    AddColumn("name").AddValue("John").
    AddColumn("email").AddValue("john@example.com")

// Maintenant (généré)
generated.Users.Insert().
    SetName("John").
    SetEmail("john@example.com")
```

## 6. Perspectives d'amélioration future

### 6.1 Améliorations possibles de la DX

**Support des requêtes SELECT typées**

```go
// Objectif futur
users := generated.Users.Select().
    Where(generated.Users.Email.Equals("john@example.com")).
    Execute(conn)
```

**Validation des relations**

```go
// Relations typées
generated.Posts.Insert().
    SetTitle("Mon article").
    SetAuthor(generated.Users.ById(1)) // Validation de FK
```

**Migration automatique**

```bash
# Détection automatique des changements de schéma
make migrate
```
