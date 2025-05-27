# postgo - ORM PostgreSQL pour Go

postgo est un ORM (Object-Relational Mapping) moderne et léger pour PostgreSQL, conçu spécifiquement pour Go. Il utilise la réflexion pour générer automatiquement des schémas SQL à partir de structures Go, offrant une approche simple et élégante pour la gestion des bases de données.

## 🚀 Fonctionnalités actuelles

### Gestion des connexions

- **Connexion PostgreSQL sécurisée** : Établissement de connexions à PostgreSQL avec gestion d'erreurs intégrée
- **Test de connectivité** : Vérification automatique de la connexion via `Ping()`
- **Fermeture propre** : Gestion automatique de la fermeture des connexions

### Modélisation de données

- **Interface Model** : Interface standardisée pour tous les modèles avec méthode `TableName()`
- **BaseModel** : Structure de base avec ID auto-incrémenté et clé primaire
- **Embedding de structures** : Support des structures imbriquées (comme `BaseModel`)
- **Tags personnalisés** : Utilisation de tags `db` pour définir les contraintes SQL

### Contraintes SQL supportées

- `primary_key` : Définition de clés primaires
- `auto_increment` : Auto-incrémentation avec type SERIAL PostgreSQL
- `not_null` : Champs obligatoires
- `unique` : Contraintes d'unicité
- Support des contraintes personnalisées

### Mapping de types

Mapping automatique des types Go vers PostgreSQL :

- `string` → `VARCHAR(255)`
- `int`, `int8`, `int16`, `int32`, `int64` → `INTEGER` (ou `SERIAL` si auto-increment)
- `float32`, `float64` → `FLOAT`
- `bool` → `BOOLEAN`

### Création de tables

- **Génération automatique de schémas** : Création de tables SQL à partir de structures Go
- **Analyse par réflexion** : Extraction automatique des métadonnées des modèles
- **CREATE TABLE IF NOT EXISTS** : Évite les erreurs de duplication
- **Noms de colonnes échappés** : Protection contre l'injection SQL

### Gestion de base de données

- **Création de bases de données** : Fonction pour créer de nouvelles bases de données PostgreSQL
- **Vérification d'existence** : Contrôle automatique avant création pour éviter les doublons

### Système de logging

- **Logging structuré** : Séparation des niveaux INFO, WARNING, ERROR
- **Horodatage automatique** : Logs avec date, heure et fichier source
- **Sortie configurable** : STDOUT pour info/warning, STDERR pour erreurs

### Infrastructure de développement

- **Docker Compose** : Configuration prête pour PostgreSQL 17.4 et Adminer
- **Variables d'environnement** : Configuration flexible via variables d'env
- **Adminer intégré** : Interface web pour l'administration de base de données

## 📋 Feuille de route

### Générateur de code Go pour l'autocomplétion

- Analyse des tables PostgreSQL existantes
- Génération automatique de structures Go avec tags appropriés
- Query builder typé avec autocomplétion IDE
- Support des relations entre tables (foreign keys)
- Génération de méthodes CRUD typées par table

### Query Builder

- Constructeur de requêtes fluide et typé
- Support des JOINs complexes
- Agrégations et fonctions SQL
- Sous-requêtes et CTEs (Common Table Expressions)

## 📦 Installation

```bash
go mod init votre-projet
go get github.com/lib/pq
```

## 🚀 Démarrage rapide

### 1. Configuration Docker (optionnel)

```bash
docker-compose up -d
```

### 2. Définition d'un modèle

```go
package main

import "postgo/db"

type User struct {
    db.BaseModel                    // ID auto-généré
    Name     string `db:"not_null"`
    Email    string `db:"not_null,unique"`
    Password string `db:"not_null"`
}

func (u *User) TableName() string {
    return "users"
}
```

### 3. Connexion et création de table

```go
package main

import (
    "postgo/db"
    _ "github.com/lib/pq"
)

func main() {
    // Connexion à PostgreSQL
    conn, err := db.NewConnection("localhost", 5432, "postgo", "postgo", "postgo")
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    // Création automatique de la table
    err = conn.CreateTable(&User{})
    if err != nil {
        panic(err)
    }
}
```

## 🏗️ Architecture

```
postgo/
├── db/               # Cœur de l'ORM
│   ├── connection.go # Gestion des connexions
│   ├── database.go   # Opérations sur les bases de données
│   ├── model.go      # Interfaces et métadonnées des modèles
│   └── table.go      # Création et gestion des tables
├── examples/         # Exemples d'utilisation
│   └── user.go       # Modèle d'exemple
├── logging/          # Système de logs
│   └── logging.go    # Configuration des loggers
└── compose.yaml      # Stack PostgreSQL + Adminer
```
