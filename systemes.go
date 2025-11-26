package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// ---------------------------- INITIALISATION QUETES ----------------------------

func initQuetes() {
	quetes = []Quete{
		{
			ID:           1,
			Nom:          "Première Chasse",
			Description:  "Tuez 5 monstres pour prouver votre valeur",
			Objectif:     "Tuer 5 monstres",
			Progres:      0,
			ProgresMax:   5,
			RecompenseOr: 100,
			RecompenseXP: 50,
			Recompense:   "Potion de Soin",
			Complete:     false,
			Active:       true,
		},
		{
			ID:           2,
			Nom:          "Le Premier Boss",
			Description:  "Affrontez et vainquez le Boss Lycaon",
			Objectif:     "Vaincre Boss Lycaon",
			Progres:      0,
			ProgresMax:   1,
			RecompenseOr: 500,
			RecompenseXP: 200,
			Recompense:   "Épée de l'Alpha",
			Complete:     false,
			Active:       false,
		},
		{
			ID:           3,
			Nom:          "Chasseur de Trésors",
			Description:  "Accumulez 1000 pièces d'or",
			Objectif:     "Gagner 1000 or",
			Progres:      0,
			ProgresMax:   1000,
			RecompenseOr: 0,
			RecompenseXP: 100,
			Recompense:   "Anneau de Fortune",
			Complete:     false,
			Active:       true,
		},
		{
			ID:           4,
			Nom:          "Maître d'Armes",
			Description:  "Achetez 10 armes différentes",
			Objectif:     "Acheter 10 armes",
			Progres:      0,
			ProgresMax:   10,
			RecompenseOr: 200,
			RecompenseXP: 150,
			Recompense:   "Cape du Guerrier",
			Complete:     false,
			Active:       true,
		},
	}

	if len(quetes) > 0 {
		queteActive = &quetes[0]
	}
}

// ---------------------------- INITIALISATION TALENTS ----------------------------

func initTalents() {
	talents = make(map[string][]Talent)

	// Talents Guerrier
	talents["GUERRIER"] = []Talent{
		{
			Nom:         "Rage du Guerrier",
			Description: "+10% dégâts pour chaque ennemi tué (max 50%)",
			Niveau:      0,
			NiveauMax:   5,
			Type:        "passif",
		},
		{
			Nom:         "Mur de Fer",
			Description: "+100% défense pendant 3 tours de combat",
			Niveau:      0,
			NiveauMax:   3,
			Type:        "actif",
			CooldownMax: 5,
		},
		{
			Nom:         "Coup Dévastateur",
			Description: "Inflige 400% des dégâts normaux",
			Niveau:      0,
			NiveauMax:   3,
			Type:        "actif",
			CooldownMax: 4,
		},
		{
			Nom:         "Endurance Titanesque",
			Description: "+50 PV Max par niveau",
			Niveau:      0,
			NiveauMax:   5,
			Type:        "passif",
		},
	}

	// Talents Mage
	talents["MAGE"] = []Talent{
		{
			Nom:         "Méditation Profonde",
			Description: "+5 mana par niveau au début de chaque combat",
			Niveau:      0,
			NiveauMax:   5,
			Type:        "passif",
		},
		{
			Nom:         "Bouclier Magique",
			Description: "Absorbe 100 dégâts pendant 2 tours",
			Niveau:      0,
			NiveauMax:   3,
			Type:        "actif",
			CooldownMax: 6,
		},
		{
			Nom:         "Météore",
			Description: "Dégâts massifs basés sur Intelligence × 3",
			Niveau:      0,
			NiveauMax:   3,
			Type:        "actif",
			CooldownMax: 5,
		},
		{
			Nom:         "Siphon de Mana",
			Description: "Vole 20% de la mana utilisée comme PV",
			Niveau:      0,
			NiveauMax:   3,
			Type:        "passif",
		},
	}

	// Talents Voleur
	talents["VOLEUR"] = []Talent{
		{
			Nom:         "Mains Lestes",
			Description: "+50 or par ennemi vaincu",
			Niveau:      0,
			NiveauMax:   5,
			Type:        "passif",
		},
		{
			Nom:         "Coup Critique Fatal",
			Description: "+20% chance de coup critique",
			Niveau:      0,
			NiveauMax:   5,
			Type:        "passif",
		},
		{
			Nom:         "Ombre Furtive",
			Description: "50% d'esquiver la prochaine attaque",
			Niveau:      0,
			NiveauMax:   3,
			Type:        "actif",
			CooldownMax: 4,
		},
	}

	// Talents Assassin
	talents["ASSASSIN"] = []Talent{
		{
			Nom:         "Lames Empoisonnées",
			Description: "Les attaques infligent +20 dégâts poison sur 3 tours",
			Niveau:      0,
			NiveauMax:   3,
			Type:        "passif",
		},
		{
			Nom:         "Exécution",
			Description: "Tue instantanément si ennemi < 20% PV",
			Niveau:      0,
			NiveauMax:   1,
			Type:        "actif",
			CooldownMax: 10,
		},
		{
			Nom:         "Frappe Précise",
			Description: "+30% dégâts avec les dagues",
			Niveau:      0,
			NiveauMax:   5,
			Type:        "passif",
		},
	}

	// Talents Archer
	talents["ARCHER"] = []Talent{
		{
			Nom:         "Flèche Perforante",
			Description: "Ignore 30% de la défense ennemie",
			Niveau:      0,
			NiveauMax:   3,
			Type:        "passif",
		},
		{
			Nom:         "Pluie de Flèches",
			Description: "5 attaques rapides de 60% dégâts",
			Niveau:      0,
			NiveauMax:   3,
			Type:        "actif",
			CooldownMax: 6,
		},
		{
			Nom:         "Œil de Faucon",
			Description: "+15% précision et dégâts à distance",
			Niveau:      0,
			NiveauMax:   5,
			Type:        "passif",
		},
	}
}

// ---------------------------- INITIALISATION ACHIEVEMENTS ----------------------------

func initAchievements() {
	achievements = []Achievement{
		{
			ID:           1,
			Nom:          "Premier Sang",
			Description:  "Tuez votre premier monstre",
			Progres:      0,
			Objectif:     1,
			Recompense:   "+10 PV Max permanent",
			Deverrouille: false,
			Icone:        "🗡️",
		},
		{
			ID:           2,
			Nom:          "Tueur de Boss",
			Description:  "Vainquez 5 boss",
			Progres:      0,
			Objectif:     5,
			Recompense:   "+50 PV Max, +10 Force",
			Deverrouille: false,
			Icone:        "👑",
		},
		{
			ID:           3,
			Nom:          "Collectionneur",
			Description:  "Possédez 20 items différents",
			Progres:      0,
			Objectif:     20,
			Recompense:   "+100 or, Accès marchand spécial",
			Deverrouille: false,
			Icone:        "📦",
		},
		{
			ID:           4,
			Nom:          "Riche",
			Description:  "Accumulez 10,000 or",
			Progres:      0,
			Objectif:     10000,
			Recompense:   "Couronne Dorée (+20 tous stats)",
			Deverrouille: false,
			Icone:        "💰",
		},
		{
			ID:           5,
			Nom:          "Sans Pitié",
			Description:  "Gagnez 50 combats sans fuir",
			Progres:      0,
			Objectif:     50,
			Recompense:   "+10% chance critique permanent",
			Deverrouille: false,
			Icone:        "⚔️",
		},
		{
			ID:           6,
			Nom:          "Invincible",
			Description:  "Terminez un combat sans recevoir de dégâts",
			Progres:      0,
			Objectif:     1,
			Recompense:   "+15 Défense permanent",
			Deverrouille: false,
			Icone:        "🛡️",
		},
		{
			ID:           7,
			Nom:          "Maître Artisan",
			Description:  "Craftez 10 items à la forge",
			Progres:      0,
			Objectif:     10,
			Recompense:   "Recettes légendaires débloquées",
			Deverrouille: false,
			Icone:        "🔨",
		},
		{
			ID:           8,
			Nom:          "Explorateur",
			Description:  "Visitez tous les lieux du jeu",
			Progres:      0,
			Objectif:     8,
			Recompense:   "Bottes de Voyage (+20% vitesse)",
			Deverrouille: false,
			Icone:        "🗺️",
		},
	}
}

// ---------------------------- GESTION QUETES ----------------------------

func verifierQuetes() {
	for i := range quetes {
		if quetes[i].Active && !quetes[i].Complete {
			switch quetes[i].ID {
			case 1: // Tuer 5 monstres
				quetes[i].Progres = monstresTotalTues
			case 3: // Accumuler 1000 or
				quetes[i].Progres = int(orTotalGagne)
			}

			if quetes[i].Progres >= quetes[i].ProgresMax {
				completerQuete(&quetes[i])
			}
		}
	}
}

func completerQuete(q *Quete) {
	if q.Complete {
		return
	}

	q.Complete = true
	xp += float64(q.RecompenseXP)
	argent += q.RecompenseOr

	if q.Recompense != "" && q.Recompense != "Aucune" {
		arme = append(arme, q.Recompense)
	}

	log.Printf("🎉 Quête complétée: %s! +%d XP, +%.0f Or", q.Nom, q.RecompenseXP, q.RecompenseOr)

	// Afficher popup de quête complétée
	AfficherPopup("quete", q.Nom+" : "+q.Recompense)

	// Activer la quête suivante
	if q.ID < len(quetes) {
		quetes[q.ID].Active = true
		queteActive = &quetes[q.ID]
	}
}

// ---------------------------- GESTION TALENTS ----------------------------

func obtenirTalent(classe string, index int) {
	if pointsTalents <= 0 {
		return
	}

	if talentsList, ok := talents[classe]; ok {
		if index >= 0 && index < len(talentsList) {
			if talentsList[index].Niveau < talentsList[index].NiveauMax {
				talentsList[index].Niveau++
				pointsTalents--
				appliquerTalent(&talentsList[index])
				talents[classe] = talentsList // Update map
				log.Printf("✨ Talent amélioré: %s (Niveau %d)", talentsList[index].Nom, talentsList[index].Niveau)
			}
		}
	}
}

func appliquerTalent(t *Talent) {
	switch t.Nom {
	case "Endurance Titanesque":
		pvMAX += 50 * float64(t.Niveau)
		pv += 50 * float64(t.Niveau)
	case "Méditation Profonde":
		mana += 5 * float64(t.Niveau)
	}
}

// ---------------------------- GESTION ACHIEVEMENTS ----------------------------

func verifierAchievements() {
	for i := range achievements {
		if achievements[i].Deverrouille {
			continue
		}

		switch achievements[i].ID {
		case 1: // Premier Sang
			achievements[i].Progres = monstresTotalTues
		case 2: // Tueur de Boss
			// Géré manuellement lors des boss
		case 3: // Collectionneur
			achievements[i].Progres = len(arme)
		case 4: // Riche
			if int(orTotalGagne) > achievements[i].Progres {
				achievements[i].Progres = int(orTotalGagne)
			}
		case 5: // Sans Pitié
			achievements[i].Progres = combatsSansFuite
		}

		if achievements[i].Progres >= achievements[i].Objectif {
			debloquerAchievement(&achievements[i])
		}
	}
}

func debloquerAchievement(a *Achievement) {
	if a.Deverrouille {
		return
	}

	a.Deverrouille = true

	// Afficher popup d'achievement
	AfficherPopup("achievement", a.Nom+" : "+a.Recompense)

	// Appliquer les récompenses
	switch a.ID {
	case 1: // Premier Sang
		pvMAX += 10
		pv += 10
	case 2: // Tueur de Boss
		pvMAX += 50
		pv += 50
		force += 10
	case 3: // Collectionneur
		argent += 100
	case 6: // Invincible
		defenseTotal += 15
	}

	log.Printf("🏆 Achievement débloqué: %s %s - %s", a.Icone, a.Nom, a.Recompense)
}

// ---------------------------- MENUS ----------------------------

func MenuQuetes(g *Game) {
	if g.JustPressed(ebiten.KeyEscape) || g.JustPressed(ebiten.KeyQ) {
		state = StateJeu
		scrollOffset = 0 // Reset scroll
		return
	}
}

func MenuTalents(g *Game) {
	if g.JustPressed(ebiten.KeyEscape) || g.JustPressed(ebiten.KeyT) {
		state = StateJeu
		return
	}

	// Sélectionner talents avec touches 1-9
	keys := []ebiten.Key{
		ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4,
		ebiten.Key5, ebiten.Key6, ebiten.Key7, ebiten.Key8, ebiten.Key9,
	}

	if talentsList, ok := talents[classe]; ok {
		for i, key := range keys {
			if i < len(talentsList) && g.JustPressed(key) {
				obtenirTalent(classe, i)
				break
			}
		}
	}
}

func MenuAchievements(g *Game) {
	if g.JustPressed(ebiten.KeyEscape) || g.JustPressed(ebiten.KeyA) {
		state = StateJeu
		scrollOffset = 0 // Reset scroll
		return
	}
}
