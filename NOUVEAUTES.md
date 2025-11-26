# 🎮 NOUVEAUTÉS - UNDERWORLD ALPHA v2.0

## 🆕 Systèmes Ajoutés

### 📜 Système de Quêtes (Touche Q)
- **4 quêtes principales** avec objectifs variés
- **Progression automatique** : tuez des monstres, accumulez de l'or
- **Récompenses** : objets spéciaux + or + XP
- **Interface dédiée** : consultez votre journal à tout moment

#### Quêtes disponibles :
1. **Tueur de Monstres** : Éliminer 10 ennemis
2. **Collectionneur d'Or** : Amasser 1000 pièces d'or
3. **Armurerie** : Acheter 3 armes différentes
4. **Guerrier Endurci** : Gagner 5 combats sans fuir

### 🌳 Arbre de Talents (Touche T)
- **1 point par niveau** gagné
- **5 classes** avec talents uniques
- **Talents améliorables** jusqu'au niveau 5
- **Bonus permanents** : dégâts, défense, régénération, etc.

#### Exemples de talents :
- **GUERRIER** : Rage de Combat (+10% dégâts max)
- **MAGE** : Maîtrise Arcane (+15% mana max)
- **VOLEUR** : Butin (+50% or trouvé max)
- **ASSASSIN** : Lame Empoisonnée (+15% dégâts poison)
- **ARCHER** : Tir Précis (+15% précision)

### 🏆 Système de Succès (Touche A)
- **8 succès débloquables**
- **Récompenses permanentes** : stats, capacités
- **Suivi de progression** : voyez combien il reste à faire
- **Icônes uniques** pour chaque succès

#### Succès disponibles :
- 🗡️ Premier Sang : Tuer 1 monstre → +5 Force
- ⚔️ Tueur en Série : Tuer 50 monstres → +15 Force
- 💰 Riche : 5000 or → +10% Or trouvé
- 🏆 Champion : 20 combats sans fuite → +20 Endurance
- 📚 Collectionneur : 20 objets → +5 slots d'inventaire
- 🛡️ Invincible : 100 défense → +25 Défense bonus
- ⚡ Puissance : 200 force → +30 Force bonus
- 🔮 Archimage : 200 intelligence → +50 Mana Max

## 🔧 Fichiers Modifiés

### main.go
- Ajout des structures `Quete`, `Talent`, `Achievement`
- 3 nouveaux états : `StateQuetes`, `StateTalents`, `StateAchievements`
- Raccourcis clavier : Q, T, A pour accéder aux menus
- Intégration des fonctions de vérification dans la boucle de jeu
- Tracking des statistiques : monstres tués, combats sans fuite
- Attribution de points de talents au level up

### systemes.go (NOUVEAU)
- `initQuetes()` : Initialise les 4 quêtes principales
- `initTalents()` : Crée l'arbre de talents pour les 5 classes
- `initAchievements()` : Configure les 8 succès
- `verifierQuetes()` : Vérifie la progression des quêtes
- `verifierAchievements()` : Vérifie si des succès sont débloqués
- `completerQuete()` : Distribue les récompenses de quête
- `debloquerAchievement()` : Distribue les récompenses de succès
- `MenuQuetes()`, `MenuTalents()`, `MenuAchievements()` : Interfaces utilisateur

## 🎨 Interface Utilisateur

### Menu Quêtes
- Liste de toutes les quêtes
- Statut : [EN COURS], [ACTIVE], [TERMINEE]
- Barre de progression (X/Y)
- Détails des récompenses

### Menu Talents
- Points disponibles affichés en haut
- Liste des talents de votre classe
- Niveau actuel / Niveau max
- Description et effet de chaque talent
- Amélioration avec touches 1-9

### Menu Succès
- Icône + Nom + Statut
- Progression détaillée
- Récompenses affichées
- Indication visuelle des succès débloqués

## 📊 Variables Globales Ajoutées

```go
// Quêtes
var quetes []Quete
var queteActive *Quete

// Talents
var talents map[string][]Talent
var talentsActifs []Talent
var pointsTalents int = 0

// Achievements
var achievements []Achievement
var monstresTotalTues int = 0
var orTotalGagne float64 = 0
var combatsSansFuite int = 0
```

## 🎯 Intégration dans le Jeu

### Initialisation (main())
```go
initQuetes()
initTalents()
initAchievements()
```

### Boucle de Jeu (Update())
```go
// Vérifier la progression à chaque frame
verifierQuetes()
verifierAchievements()

// Nouveaux états
case StateQuetes:
    MenuQuetes(g)
case StateTalents:
    MenuTalents(g)
case StateAchievements:
    MenuAchievements(g)
```

### Tracking dans le Combat
```go
// Après chaque victoire
monstresTotalTues++
combatsSansFuite++

// À chaque level up
pointsTalents++

// Si le joueur fuit
combatsSansFuite = 0
```

## 📈 Progression

### Quêtes → Récompenses
Complétez des objectifs pour obtenir objets uniques, or, et XP

### Niveaux → Talents
Gagnez des points de talents pour améliorer votre personnage

### Exploits → Succès
Accomplissez des défis pour des bonus permanents

## 🚀 Comment Jouer

1. **Lancez le jeu** : `underworld.exe`
2. **Pendant le jeu**, appuyez sur :
   - **[Q]** : Voir vos quêtes
   - **[T]** : Gérer vos talents
   - **[A]** : Consulter vos succès
3. **Combattez des monstres** pour progresser dans les quêtes
4. **Gagnez des niveaux** pour obtenir des points de talents
5. **Débloquez des succès** pour des bonus permanents

## 🎊 Prochaines Étapes

### Systèmes à Implémenter
- 🔨 **Craft/Forge** : Améliorer les armes et créer des objets
- 💥 **Effets Visuels** : Particules, screen shake, nombres de dégâts
- 🎵 **Sons d'Interface** : Feedback audio pour les actions
- 📦 **Plus de Quêtes** : Quêtes secondaires et quêtes d'histoire
- 🏪 **Marchand Spécialisé** : Acheter des matériaux de craft

## 🐛 Tests Nécessaires

- ✅ Compilation réussie
- ⏳ Test des menus Q/T/A
- ⏳ Vérification de la progression des quêtes
- ⏳ Test des points de talents au level up
- ⏳ Déblocage des succès
- ⏳ Sauvegarde/Chargement avec nouveaux systèmes

## 📝 Notes Techniques

- Les talents sont stockés dans un `map[string][]Talent` par classe
- Les quêtes se vérifient automatiquement à chaque frame
- Les achievements se débloquent dès que l'objectif est atteint
- Le compteur `combatsSansFuite` se réinitialise si le joueur fuit
- Les images `quete2.png` et `quete3.png` sont utilisées pour les menus Talents et Achievements

## 🎮 Bon Jeu !

Les nouveaux systèmes enrichissent considérablement l'expérience de jeu avec :
- Des objectifs clairs (quêtes)
- Une personnalisation profonde (talents)
- Des défis à relever (succès)

Explorez, combattez, et progressez dans le monde d'Underworld ! 🌟
