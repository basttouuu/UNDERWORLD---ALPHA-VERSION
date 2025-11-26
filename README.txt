# 🎮 UNDERWORLD - ALPHA VERSION

## 📖 Description
UNDERWORLD est un RPG d'aventure développé en Go avec Ebiten. Explorez un monde fantastique, combattez des monstres, complétez des quêtes et débloquez des succès !

## ✨ Fonctionnalités Principales

### 🎯 Système de Jeu
- **5 Classes jouables** : Guerrier, Mage, Voleur, Assassin, Archer
- **Combat au tour par tour** avec IA avancée
- **Système de progression** : XP, niveaux, talents
- **Équipement** : 5 slots (Casque, Plastron, Bottes, Anneau, Arme)
- **Inventaire** : Stockage jusqu'à 10 items
- **Boutique** : Achat d'armes et équipements

### 📜 Progression
- **4 Quêtes** avec objectifs et récompenses
- **8 Achievements** déblocables
- **Arbre de talents** par classe (5 talents/classe)
- **Système de loot** après combats
- **Boss** : 3 boss principaux avec comportement avancé

### 🤖 IA et Combat
- **IA Ennemie avancée** : Marche, Attaque, Défense, Compétences
- **Comportement Boss** : Tactique adaptée selon HP
- **14 types de monstres** avec stats variées
- **Compétences spéciales** par classe
- **Système de défense** et réduction de dégâts

### 🎨 Interface
- **Popups animés** : Victoire, Défaite, Level Up, Quêtes, Achievements
- **Menus centrés** avec design soigné
- **Système de scroll** pour pages longues (PageUp/PageDown)
- **Barres de progression** visuelles
- **Code couleur** pour meilleure lisibilité

### 💾 Système de Sauvegarde
- **F5** : Sauvegarder la partie
- **F9** : Charger la sauvegarde
- Sauvegarde JSON complète (stats, équipement, progression)

## 🚀 Installation

### Prérequis
- **Go 1.18+** (https://go.dev/dl/)
- **Git** (optionnel)

### Lancement du Jeu

#### Option 1 : Exécutable (Windows)
```bash
# Double-cliquez sur underworld.exe
./underworld.exe
```

#### Option 2 : Compilation depuis les sources
```bash
# Cloner le projet
git clone https://github.com/basttouuu/UNDERWORLD---ALPHA-VERSION.git
cd UNDERWORLD---ALPHA-VERSION

# Installer les dépendances
go mod download

# Compiler
go build -o underworld.exe

# Lancer
./underworld.exe
```

#### Option 3 : Exécution directe
```bash
go run .
```

## 🎮 Contrôles

### Menu Principal
- **[ESPACE]** : Continuer/Valider
- **[G/N/V/A/E]** : Choisir classe (Guerrier/mage/Voleur/Assassin/archEr)
- **[S]** : Voir statistiques
- **[B]** : Boutique
- **[Q]** : Journal de quêtes
- **[T]** : Arbre de talents
- **[J]** : Achievements
- **[I]** : Inventaire
- **[M]** : Carte du monde
- **[ÉCHAP]** : Retour/Pause

### Combat
- **[Q]** : Attaque légère
- **[W]** : Attaque forte (25% chance de rater)
- **[E]** : Compétence spéciale (coûte mana)
- **[I]** : Ouvrir inventaire
- **[0-9]** : Utiliser item (en combat = action)
- **[ENTRÉE]** : Fuir le combat

### Navigation
- **[PageUp]** : Scroll vers le haut
- **[PageDown]** : Scroll vers le bas
- **[F5]** : Sauvegarder
- **[F9]** : Charger

## 📁 Structure du Projet

```
UNDERWORLD---ALPHA-VERSION/
├── main.go                    # Logique principale du jeu
├── systemes.go               # Quêtes, Talents, Achievements
├── sauvegarde.go             # Système de sauvegarde/chargement
├── beep.go                   # Gestion audio
├── initPseudo.go             # Initialisation joueur
├── map.go                    # Système de carte
├── medieval.ttf              # Police de caractères
├── musique.mp3               # Musique de fond
├── image/                    # Images de fond et UI
├── assets/
│   ├── perso/               # Sprites des personnages
│   └── monstres/            # Sprites des monstres
├── GUIDE_JOUEUR.md          # Guide détaillé du joueur
├── NOUVEAUTES.md            # Liste des nouveautés
├── CHANGELOG.md             # Historique des versions
├── NOUVELLES_FONCTIONNALITES.md  # Documentation technique
└── README.txt               # Ce fichier

```

## 📚 Documentation

- **GUIDE_JOUEUR.md** : Guide complet pour les joueurs
- **NOUVEAUTES.md** : Détails des systèmes implémentés
- **CHANGELOG.md** : Historique des mises à jour (v2.0 → v2.3)
- **NOUVELLES_FONCTIONNALITES.md** : Documentation technique v2.3

## 🎯 Progression Suggérée

1. **Démarrage** : Choisir une classe et nommer votre personnage
2. **Exploration** : Visitez les différents lieux (Champs, Forge, Tour, etc.)
3. **Combat** : Sortez par la Grande Porte pour combattre
4. **Quête 1** : Tuez 5 monstres (Crabauge → Vorlapin → Gobelin → Boss Lycaon)
5. **Amélioration** : Achetez équipement à la boutique
6. **Talents** : Dépensez points de talents (1 par niveau)
7. **Quête 2-3** : Continuez l'aventure contre des ennemis plus forts

## 🏆 Système d'Achievements

### Premier Sang
Tuez votre premier monstre
**Récompense** : +10 PV Max

### Tueur de Boss
Vainquez 5 boss
**Récompense** : +50 PV Max, +10 Force

### Collectionneur
Possédez 20 items différents
**Récompense** : +100 or, Accès marchand spécial

### Riche
Accumulez 10,000 or
**Récompense** : Couronne Dorée (+20 tous stats)

### Sans Pitié
Gagnez 50 combats sans fuir
**Récompense** : +10% chance critique permanent

### Invincible
Gagnez 100 combats sans mourir
**Récompense** : +15 Défense, Titre "L'Invincible"

### Maître Artisan
Utilisez tous les lieux 50 fois
**Récompense** : Accès forge légendaire

### Explorateur
Découvrez tous les secrets du monde
**Récompense** : Cape de l'Explorateur

## 🐛 Problèmes Connus

Aucun bug majeur connu. Si vous rencontrez un problème :
1. Vérifiez que tous les fichiers image/audio sont présents
2. Utilisez `verifier.ps1` pour vérifier les ressources (Windows)
3. Consultez les logs dans la console

## 🔄 Versions

- **v2.3** (Actuelle) : Popups, IA avancée, Scroll
- **v2.2** : Espacement UI, Équipement direct
- **v2.1** : Optimisations, UI redesign
- **v2.0** : Quêtes, Talents, Achievements
- **v1.x** : Version initiale

## 🤝 Contribution

Projet développé par **basttouuu**

## 📜 Licence

Projet personnel - Tous droits réservés

## 🎮 Bon Jeu !

Amusez-vous bien dans UNDERWORLD et n'oubliez pas de sauvegarder régulièrement avec F5 !

---

**Version** : 2.3 Alpha
**Date** : Novembre 2025
**Moteur** : Ebiten v2
**Langage** : Go 1.18+
