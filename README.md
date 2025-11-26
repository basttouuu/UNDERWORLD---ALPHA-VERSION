# 🎮 UNDERWORLD - ALPHA VERSION

> Un RPG d'aventure épique développé en Go avec Ebiten

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://go.dev)
[![Ebiten](https://img.shields.io/badge/Ebiten-v2-orange?style=flat)](https://ebiten.org)
[![Version](https://img.shields.io/badge/Version-2.3%20Alpha-brightgreen?style=flat)](https://github.com/basttouuu/UNDERWORLD---ALPHA-VERSION)

## 📖 Description

UNDERWORLD est un RPG d'aventure dans un monde fantastique où vous incarnez un héros devant combattre des monstres, accomplir des quêtes et débloquer des succès. Avec un système de combat au tour par tour, une IA ennemie avancée et une progression riche, plongez dans une aventure épique !

## ✨ Fonctionnalités

### 🎯 Gameplay
- ⚔️ **5 Classes jouables** : Guerrier, Mage, Voleur, Assassin, Archer
- 🎲 **Combat tactique** au tour par tour avec IA intelligente
- 📈 **Progression** : Système d'XP, niveaux et talents
- 🛡️ **Équipement** : 5 slots d'équipement avec bonus
- 🎒 **Inventaire** : Gestion de 10 items max
- 🏪 **Boutique** : Achat d'armes et équipements avec équipement direct

### 📜 Contenu
- ✅ **4 Quêtes** principales avec objectifs et récompenses
- 🏆 **8 Achievements** déblocables
- 🌳 **Arbre de talents** : 5 talents uniques par classe
- 👹 **14 types de monstres** avec stats et comportements variés
- 🐉 **3 Boss** avec IA tactique avancée

### 🎨 Interface
- 🎭 **Popups animés** pour victoire, défaite, level up, quêtes, achievements
- 📱 **UI centrée** avec design épuré
- 📜 **Système de scroll** (PageUp/PageDown) pour menus longs
- 📊 **Barres de progression** visuelles partout
- 🎨 **Code couleur** intuitif

### 🤖 IA Avancée
- 🚶 **Déplacement** : Ennemis marchent vers le joueur
- ⚔️ **Attaques variées** : Normale, Lourde, Compétence spéciale
- 🛡️ **Défense tactique** : Position défensive (-40% dégâts)
- 🧠 **Boss intelligents** : Comportement adaptatif selon HP
- ⏱️ **Cooldowns** : Timing entre actions

### 💾 Sauvegarde
- **F5** : Sauvegarder
- **F9** : Charger
- Format JSON avec toutes les données

## 🚀 Installation

> 💡 **Débutant ?** Consultez le [Guide de Démarrage Rapide](QUICKSTART.md) pour être opérationnel en 30 secondes !

### Prérequis
- [Go 1.18+](https://go.dev/dl/)
- Git (optionnel)

### Méthode 1 : Exécutable Windows
```bash
# Téléchargez underworld.exe et lancez-le
./underworld.exe
```

### Méthode 2 : Compilation
```bash
# Cloner le dépôt
git clone https://github.com/basttouuu/UNDERWORLD---ALPHA-VERSION.git
cd UNDERWORLD---ALPHA-VERSION

# Télécharger les dépendances
go mod download

# Compiler
go build -o underworld.exe

# Lancer
./underworld.exe
```

### Méthode 3 : Exécution directe
```bash
go run .
```

## 🎮 Contrôles

### Navigation
| Touche | Action |
|--------|--------|
| `ESPACE` | Continuer/Valider |
| `ÉCHAP` | Retour/Pause |
| `PageUp/PageDown` | Scroll |
| `F5` | Sauvegarder |
| `F9` | Charger |

### Menu Principal
| Touche | Action |
|--------|--------|
| `S` | Statistiques |
| `B` | Boutique |
| `Q` | Journal de quêtes |
| `T` | Arbre de talents |
| `J` | Achievements |
| `I` | Inventaire |
| `M` | Carte |

### Combat
| Touche | Action |
|--------|--------|
| `Q` | Attaque légère |
| `W` | Attaque forte (25% miss) |
| `E` | Compétence (coûte mana) |
| `0-9` | Utiliser item |
| `ENTRÉE` | Fuir |

## 📊 Classes

### 🗡️ Guerrier
- **Style** : Corps à corps, tank
- **Talents** : Rage du Guerrier, Mur de Fer, Coup Dévastateur
- **Force** : Haute défense et dégâts physiques

### 🔮 Mage
- **Style** : Magie, distance
- **Talents** : Maîtrise Élémentaire, Bouclier Arcane, Nova Mystique
- **Force** : Compétences puissantes, mana élevé

### 🗡️ Voleur
- **Style** : Agilité, esquive
- **Talents** : Ombre Furtive, Combo Mortel, Pickpocket
- **Force** : Critiques fréquents, vitesse

### 🔪 Assassin
- **Style** : Burst damage, poison
- **Talents** : Lames Empoisonnées, Exécution, Maître des Ombres
- **Force** : Dégâts massifs instantanés

### 🏹 Archer
- **Style** : Distance, précision
- **Talents** : Œil de Faucon, Tir Perçant, Pluie de Flèches
- **Force** : Portée, dégâts constants

## 🏆 Achievements

| Succès | Objectif | Récompense |
|--------|----------|------------|
| 🩸 Premier Sang | Tuer 1 monstre | +10 PV Max |
| 👑 Tueur de Boss | Vaincre 5 boss | +50 PV Max, +10 Force |
| 📦 Collectionneur | 20 items | +100 or, Marchand spécial |
| 💰 Riche | 10,000 or | Couronne Dorée (+20 stats) |
| ⚔️ Sans Pitié | 50 combats sans fuite | +10% crit permanent |
| 🛡️ Invincible | 100 combats sans mort | +15 Défense |
| 🔨 Maître Artisan | 50 utilisations lieux | Forge légendaire |
| 🗺️ Explorateur | Tous les secrets | Cape Explorateur |

## 📸 Captures d'écran

```
[Menu Principal]    [Combat]         [Quêtes]        [Talents]
   Centré          IA Avancée       Progression      Arbre
   Design épuré    Popups animés    Barres visuelles Customisation
```

## 📁 Structure du Projet

```
UNDERWORLD---ALPHA-VERSION/
├── main.go                 # Logique principale (3200+ lignes)
├── systemes.go            # Quêtes, Talents, Achievements
├── sauvegarde.go          # Système de sauvegarde JSON
├── beep.go                # Gestion audio
├── initPseudo.go          # Init personnage
├── map.go                 # Système de carte
├── go.mod / go.sum        # Dépendances
├── medieval.ttf           # Police médiévale
├── musique.mp3            # Musique de fond
├── image/                 # 20+ images UI/décors
├── assets/
│   ├── perso/            # 5 sprites classes
│   └── monstres/         # 14 sprites monstres
└── docs/                  # Documentation complète
```

## 🔄 Historique des Versions

### v2.3 (Actuelle) - Novembre 2025
- ✅ Système de popups animés
- ✅ IA ennemie avancée (4 comportements)
- ✅ Scroll pour menus longs
- ✅ Boss avec tactiques adaptatives

### v2.2
- ✅ Espacement UI amélioré
- ✅ Équipement direct après achat
- ✅ Meilleurs espacements partout

### v2.1
- ✅ UI redesign complet
- ✅ Menus centrés
- ✅ Optimisations performance

### v2.0
- ✅ Système de quêtes (4)
- ✅ Arbre de talents (25)
- ✅ Achievements (8)
- ✅ Sauvegarde/Chargement

## 🛠️ Technologies

- **Langage** : [Go 1.18+](https://go.dev)
- **Moteur** : [Ebiten v2](https://ebiten.org) (game engine 2D)
- **Audio** : [Beep](https://github.com/faiface/beep)
- **Fonts** : [golang.org/x/image/font](https://pkg.go.dev/golang.org/x/image/font)

## 🐛 Bugs Connus

Aucun bug majeur. Le jeu est stable et testé.

Si problème :
1. Vérifiez présence des ressources (images/audio)
2. Utilisez `verifier.ps1` (Windows)
3. Consultez logs console

## 🤝 Contribution

Projet personnel développé par **basttouuu**

Suggestions bienvenues via Issues !

## 📚 Documentation

- 📖 [GUIDE_JOUEUR.md](GUIDE_JOUEUR.md) - Guide complet
- 🆕 [NOUVEAUTES.md](NOUVEAUTES.md) - Systèmes détaillés
- 📝 [CHANGELOG.md](CHANGELOG.md) - Historique complet
- 🔧 [NOUVELLES_FONCTIONNALITES.md](NOUVELLES_FONCTIONNALITES.md) - Docs technique v2.3

## 📜 Licence

© 2025 basttouuu - Tous droits réservés

## 🎮 Bon Jeu !

N'oubliez pas de **sauvegarder avec F5** régulièrement !

---

⭐ Si vous appréciez le jeu, laissez une étoile sur GitHub !

**Version** : 2.3 Alpha | **Date** : Novembre 2025 | **Moteur** : Ebiten v2
