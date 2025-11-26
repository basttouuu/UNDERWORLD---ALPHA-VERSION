package main

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
)

// ---------------------------- CONFIG FENETRE ----------------------------
const (
	WinW = 1920
	WinH = 1080
)

// ---------------------------- DONNEES ET STRUCTURES ----------------------------

type Arme struct {
	Nom   string
	Prix  float64
	Degat float64
}

type Equipement struct {
	Nom            string
	Prix           float64
	Type           string // "arme", "armure", "accessoire", "casque", "bottes", "anneau"
	Degat          float64
	Defense        float64
	BonusForce     float64
	BonusDex       float64
	BonusInt       float64
	BonusMana      float64
	BonusPV        float64
	BonusEndurance float64
}

type Monstre struct {
	Nom           string
	PVMax         float64
	DegatsMin     float64
	DegatsMax     float64
	XPGagneMin    int
	XPGagneMax    int
	Couleur       color.RGBA
	LootCategorie string
}

// Loots par monstre
var lootTable = map[string][]string{
	"Crabauge":            {"Carapace fissurée", "Pinces usées"},
	"Boss Lycaon":         {"Crocs de lycaon", "Peau de bête"},
	"Gobelin":             {"Dague rouillée", "Collier de dents"},
	"Muddig":              {"Éclat de pierre", "Fragment runique"},
	"Gros Serpent Vorace": {"Peau de Serpent", "Crocs de Serpent"},
	"Vorlapin":            {"Pattes de Lapin", "Couteau usée"},
	"Serpent Livestide":   {"Peau de Serpent Bleu", "Crocs de Serpent Bleu"},
	"Loosers Wood":        {"Morceaux de Bois", "Chaines"},
	"Boss Wezaemon":       {"Épée Technologique", "Casque Brisé"},
	"Poiscaille Zombie":   {"Écaille", "Épée Zombie Brisées"},
	"Lugia":               {"Écailles", "Griffes"},
	"Leviathan":           {"Écailles", "Crocs"},
	"Atlanticus Repunorca - Orque électrique": {"Écailles", "Crocs"},
	"Kthaanid - Maître des Abysses":           {"Écailles", "Crocs"},
}

var queteActuelle int = 1 // 1 = quête 1, 2 = quête 2 et.....

// Dernier loot obtenu
var dernierLoot string

type Competence struct {
	Nom        string
	CoutMana   float64
	MultDegats float64
	Effet      string
}

// ---------------------------- QUETES ----------------------------

type Quete struct {
	ID           int
	Nom          string
	Description  string
	Objectif     string
	Progres      int
	ProgresMax   int
	Recompense   string
	RecompenseOr float64
	RecompenseXP int
	Complete     bool
	Active       bool
}

var quetes []Quete
var queteActive *Quete

// ---------------------------- TALENTS ----------------------------

type Talent struct {
	Nom         string
	Description string
	Niveau      int
	NiveauMax   int
	Type        string // "passif" ou "actif"
	Effet       func()
	Cooldown    int
	CooldownMax int
}

var talents map[string][]Talent
var pointsTalents int = 0

// ---------------------------- ACHIEVEMENTS ----------------------------

type Achievement struct {
	ID           int
	Nom          string
	Description  string
	Progres      int
	Objectif     int
	Recompense   string
	Deverrouille bool
	Icone        string
}

var achievements []Achievement
var monstresTotalTues int = 0
var orTotalGagne float64 = 0
var combatsSansFuite int = 0

var armesDispo = []Arme{{"CURE DENT", 5, 1}}
var armesGuerrier = []Arme{
	{"Épée Longue", 30, 10},
	{"Hache de Guerre", 45, 14},
	{"Bouclier Renforcé", 25, 0},
}
var armesMage = []Arme{
	{"Baguette de Feu", 35, 8},
	{"Bâton Mystique", 50, 12},
	{"Grimoire Ancien", 40, 5},
}
var armesVoleur = []Arme{
	{"Dague Aiguisée", 20, 6},
	{"Arc Court", 25, 7},
	{"Arbalète Légère", 30, 8},
}
var armesAssassin = []Arme{
	{"Lames Jumelles", 35, 9},
	{"Kris Empoisonné", 40, 10},
	{"Arc Noir", 45, 11},
}
var armesARCHER = []Arme{
	{"Arc Elfique", 40, 12},
	{"Lance de Bois Sacré", 35, 9},
	{"Épée Fine", 30, 8},
}

var bestiaire = []Monstre{
	// Nom, PV, AttMin, AttMax, DefMin, DefMax, Couleur, Type

	// Quête 1
	{"Crabauge", 100, 3, 6, 5, 8, color.RGBA{0, 255, 128, 255}, "bete"},       // scale 1.0
	{"Vorlapin", 180, 4, 7, 6, 9, color.RGBA{200, 50, 50, 255}, "bete"},       // scale 1.8
	{"Gobelin", 260, 5, 8, 7, 10, color.RGBA{80, 200, 80, 255}, "humanoide"},  // scale 2.6
	{"Boss Lycaon", 420, 6, 10, 9, 12, color.RGBA{50, 180, 120, 255}, "bete"}, // scale 3.5

	// Quête 2
	{"Muddig", 600, 7, 11, 10, 13, color.RGBA{50, 180, 120, 255}, "humanoide"},         // scale 5.0
	{"Gros Serpent Vorace", 780, 8, 12, 11, 14, color.RGBA{50, 180, 120, 255}, "bete"}, // scale 6.0
	{"Serpent Livestide", 820, 9, 13, 12, 15, color.RGBA{50, 180, 120, 255}, "bete"},   // scale 6.2
	{"Loosers Wood", 900, 9, 13, 12, 15, color.RGBA{50, 180, 120, 255}, "bete"},        // scale 6.5
	{"Boss Wezaemon", 1100, 10, 14, 14, 17, color.RGBA{50, 180, 120, 255}, "bete"},     // scale 7.5

	// Quête 3
	{"Poiscaille Zombie", 950, 8, 12, 12, 15, color.RGBA{50, 180, 120, 255}, "humanoide"},          // scale 8.0
	{"Lugia", 1050, 9, 13, 13, 16, color.RGBA{50, 180, 120, 255}, "bete"},                          // scale 8.5
	{"Leviathan", 1200, 10, 14, 14, 17, color.RGBA{50, 180, 120, 255}, "bete"},                     // scale 9.0
	{"Altanticus Repunorca", 1250, 10, 15, 14, 18, color.RGBA{50, 180, 120, 255}, "bete"},          // scale 9.2
	{"Kthaanid - Maître des Abysses", 1600, 13, 17, 16, 20, color.RGBA{50, 180, 120, 255}, "bete"}, // 10
}

var pageBoutique int = 1

var shopItems = []struct {
	Nom  string
	Prix int
}{
	{"Épée rouillée", 30},
	{"Arc court", 45},
	{"Bâton de novice", 25},
}

var shopEquipements = []struct {
	Nom  string
	Prix int
}{
	{"Casque en fer", 50},
	{"Plastron de cuir", 80},
	{"Bottes renforcées", 40},
	{"Anneau magique", 120},
}

// Slots d’équipement
var equipement = struct {
	Casque   string
	Plastron string
	Bottes   string
	Anneau   string
}{}

var equipSlots Slots

var selection *Equipement

var mesEquipements []Equipement

type Slots struct {
	Casque *Equipement
	Armure *Equipement
	Anneau *Equipement
	Arme   *Equipement
	Bottes *Equipement
}

// Bonus conférés par chaque équipement
var equipementBonus = map[string]struct {
	PV           float64
	Force        float64
	Endurance    float64
	Dexterite    float64
	Intelligence float64
	Mana         float64
	Defense      float64
}{
	"Casque en fer":     {PV: 15, Endurance: 20, Defense: 5},
	"Plastron de cuir":  {PV: 25, Endurance: 30, Defense: 10},
	"Bottes renforcées": {PV: 10, Endurance: 10, Dexterite: 5},
	"Anneau magique":    {Mana: 50, Intelligence: 20, Force: 10},
	"Gants de combat":   {Force: 15, Dexterite: 10},
	"Cape mystique":     {Intelligence: 25, Mana: 30},
}

var statLoot []string

var competences = map[string][]Competence{
	"MAGE": {
		{"Boule de Feu", 10, 1.5, "brûlure"},
		{"Éclair Mystique", 15, 2.0, "étourdissement"},
	},
	"GUERRIER": {
		{"Coup Puissant", 5, 1.8, ""},
		{"Charge", 8, 2.2, "étourdissement"},
	},
	"VOLEUR": {
		{"Attaque Sournoise", 6, 1.7, "critique"},
		{"Pluie de Flèches", 12, 2.0, ""},
	},
	"ASSASSIN": {
		{"Coup Mortel", 10, 2.5, "saignement"},
		{"Poison", 8, 1.2, "empoisonnement"},
	},
	"ARCHER": {
		{"Flèche Enchantée", 8, 1.6, "perce-armure"},
		{"Lame Sylvestre", 10, 2.0, ""},
	},
}

var (
	medievalFont font.Face
)

// État du joueur
var (
	playerIsMoving  bool
	playerFacing    = "droite"
	playerAttacking bool
)

// Données de jeu
var (
	classe               string
	state                int
	prevState            int     // <- pour revenir au menu précédent (utilisé pour inventaire)
	argent               float64 = 200
	intelligence         float64
	force                float64
	agilite              float64
	endurance            float64
	mana                 float64
	dexterite            float64
	degats               float64
	arme                 []string
	equipementArme       *Equipement
	equipementArmure     *Equipement
	equipementAccessoire *Equipement
	defenseTotal         float64
	titre                string
	lieuactu             string
	toto                 string
	pv                   float64 = 100
	pvMAX                float64 = 100
	level                float64 = 1
	xp                   float64
	xpMax                float64 = 50
	Nom                  string  // Déclaré ici, supposé initialisé par initPseudo
	tata                 float64 = 0
	tota                 float64 = 0
	tato                 float64 = 0
	tto                  float64 = 0
	tta                  float64 = 0
)

// Combat
var (
	combatInit      bool
	tourDuJoueur    bool
	logCombat       []string
	ennemiNom       string
	ennemiPV        float64
	ennemiPVMax     float64
	ennemiDegMin    float64
	ennemiDegMax    float64
	ennemiColor     color.RGBA
	ennemiXPGainMin int
	ennemiXPGainMax int
	ennemiTypeLoot  string
	fuitePossible   bool = true
	monster         string
	monstreIndex    int
	ordreMonstres1  = []string{"Crabauge", "Vorlapin", "Gobelin", "Boss Lycaon", "Muddig", "Gros Serpent Vorace", "Serpent Livestide", "Loosers Wood", "Boss Wezaemon", "Poiscaille Zombie", "Lugia", "Leviathan", "Atlanticus Repunorca - Orque électrique", "Kthaanid - Maître des Abysses"}

	// IA Ennemi
	ennemiAction    string // "marcher", "attaquer", "défendre", "compétence"
	ennemiCooldown  int
	ennemiDefense   bool
	ennemiPosX      float64 = 1400
	ennemiPosY      float64 = 400
	ennemiVitesse   float64 = 2.0
	ennemiDirection int     = -1 // -1 = vers gauche (joueur), 1 = vers droite
)

// Popups
var (
	popupActive  bool
	popupType    string // "victoire", "defaite", "levelup", "quete", "achievement"
	popupMessage string
	popupTimer   int
	popupAlpha   float64
	popupScale   float64 = 1.0
)

// Scroll pour pages longues
var (
	scrollOffset int
	scrollMax    int
)

// Drop/Victoire
var (
	levelUpMsg       string
	armeSelectionnee Arme
)

// Images
var (
	accueilImage        *ebiten.Image
	introImage          *ebiten.Image
	selectionImage      *ebiten.Image
	champImage          *ebiten.Image
	fondImage           *ebiten.Image
	shopImage           *ebiten.Image
	statImage           *ebiten.Image
	inventaireImage     *ebiten.Image
	tourImage           *ebiten.Image
	forgeImage          *ebiten.Image
	hotelImage          *ebiten.Image
	sortieImage         *ebiten.Image
	pauseImage          *ebiten.Image
	mortImage           *ebiten.Image
	foretImage          *ebiten.Image
	victoireImage       *ebiten.Image
	defaiteImage        *ebiten.Image
	resultatImage       *ebiten.Image
	fondqueteImage      *ebiten.Image
	fondquete2Image     *ebiten.Image
	fondquete3Image     *ebiten.Image
	mcombatImage        *ebiten.Image
	gcombatImage        *ebiten.Image
	vcombatImage        *ebiten.Image
	acombatImage        *ebiten.Image
	ecombatImage        *ebiten.Image
	histoire1Image      *ebiten.Image
	histoire2Image      *ebiten.Image
	histoire3Image      *ebiten.Image
	creditImage         *ebiten.Image
	fondequipementImage *ebiten.Image
	jeuImage            *ebiten.Image
	artisantImage       *ebiten.Image
	combatImage         *ebiten.Image
)

// États
const (
	StateAccueil = iota
	StateIntro
	StateStart
	StateClasse
	StateArme
	StateJeu
	StateBoutique
	StateConfirmation
	StateVendre
	StateLieu
	StateSortie
	StatePause
	StateAttention
	StateAttentionPrix
	StateStat
	StateForge
	StateInventaire
	StateChamps
	StateHotel
	StateTour
	StateArtisant
	StateBackLife
	StateForet
	StateCombat
	StateVictoire
	StateDefaite
	StateResultat
	StateHistoire1
	StateThemeG
	StateThemeA
	StateThemeM
	StateThemeV
	StateThemeE
	StateCredit
	StateEquipement
	StateHistoire2
	StateHistoire3
	StateQuetes
	StateTalents
	StateAchievements
	StateProposerEquiper
)

// ---------------------------- JEU ET INPUT ----------------------------

type Game struct {
	prevKeys           map[ebiten.Key]bool
	JoueurImage        *ebiten.Image
	MonstreImage       *ebiten.Image
	X, Y               float64
	DamageTimer        int
	MonstreX, MonstreY float64
	MonstreDX          float64
	MonstreDamageTimer int
}

func (g *Game) JustPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key) && !g.prevKeys[key]
}

func (g *Game) Update() error {
	// Toggle plein écran
	if g.JustPressed(ebiten.KeyF5) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	// Init prevKeys
	if g.prevKeys == nil {
		g.prevKeys = make(map[ebiten.Key]bool)
		for k := ebiten.Key(0); k < ebiten.KeyMax; k++ {
			g.prevKeys[k] = ebiten.IsKeyPressed(k)
		}
	}

	// Gestion des états
	switch state {
	case StateAccueil:
		if g.JustPressed(ebiten.KeySpace) {
			state = StateIntro
		}
	case StateIntro:
		if g.JustPressed(ebiten.KeySpace) {
			state = StateStart
		}
	case StateStart:
		if g.JustPressed(ebiten.KeySpace) {
			state = StateClasse
		}
	case StateClasse:
		ClasseSelection(g)
	case StateArme:
		ArmeSelection(g)
	case StateJeu:
		Jeu(g)
	case StateBoutique:
		Boutique(g)
	case StateConfirmation:
		Confirmation(g)
	case StateVendre:
		Vente(g)
	case StateLieu:
		Lieu(g)
	case StateSortie:
		Sortie(g)
	case StatePause:
		Pause(g)
	case StateHotel:
		Hotel(g)
	case StateForge:
		Forge(g)
	case StateForet:
		Foret(g)
	case StateCombat:
		if g.JustPressed(ebiten.KeyI) {
			prevState = StateCombat
			state = StateInventaire
			return nil
		}
		if !combatInit {
			InitialiserCombat(g)
		}
		Combat(g)
		MettreAJourCombat()
	case StateResultat:
		Resultat(g)

	case StateAttention, StateAttentionPrix:
		if g.JustPressed(ebiten.KeyB) {
			state = StateBoutique
		}
	case StateStat:
		if g.JustPressed(ebiten.KeyEscape) {
			state = StateJeu
		}
	case StateInventaire:
		Inventaire(g)
	case StateChamps:
		Champs(g)
	case StateTour:
		Tour(g)
	case StateVictoire:
		Victoire(g)
	case StateDefaite:
		Defaite(g)
	case StateBackLife:
		BackLife(g)
	case StateHistoire1:
		Histoire1(g)
	case StateHistoire2:
		Histoire2(g)
	case StateHistoire3:
		Histoire3(g)
	case StateThemeG:
		ThemeG(g)
	case StateThemeA:
		ThemeA(g)
	case StateThemeV:
		ThemeV(g)
	case StateThemeE:
		ThemeE(g)
	case StateThemeM:
		ThemeM(g)
	case StateCredit:
		Credit(g)
	case StateEquipement:
		UpdateEquipementMenu(g)
	case StateQuetes:
		MenuQuetes(g)
	case StateTalents:
		MenuTalents(g)
	case StateAchievements:
		MenuAchievements(g)
	case StateProposerEquiper:
		ProposerEquiper(g)
	}

	// Vérifier quêtes et achievements à chaque frame
	verifierQuetes()
	verifierAchievements()

	// Mettre à jour les popups
	UpdatePopup()

	// Gestion du scroll avec PageUp/PageDown
	if g.JustPressed(ebiten.KeyPageDown) {
		scrollOffset += 100
		if scrollOffset > scrollMax {
			scrollOffset = scrollMax
		}
	}
	if g.JustPressed(ebiten.KeyPageUp) {
		scrollOffset -= 100
		if scrollOffset < 0 {
			scrollOffset = 0
		}
	}

	// Animation combat
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.X -= 10
		playerIsMoving = true
		playerFacing = "left"
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.X += 10
		playerIsMoving = true
		playerFacing = "right"
	}

	// Timer dégâts joueur
	if g.DamageTimer > 0 {
		g.DamageTimer--
	}

	// Timer dégâts monstre
	if g.MonstreDamageTimer > 0 {
		g.MonstreDamageTimer--
	}

	// Animation monstre (va-et-vient)
	g.MonstreX += g.MonstreDX
	screenW, _ := ebiten.WindowSize()
	if g.MonstreX < float64(screenW)/2 || g.MonstreX > float64(screenW)-300 {
		g.MonstreDX *= -1
	}

	// Détection attaque
	playerAttacking = ebiten.IsKeyPressed(ebiten.KeySpace)

	// Mise à jour des touches
	for k := ebiten.Key(0); k < ebiten.KeyMax; k++ {
		g.prevKeys[k] = ebiten.IsKeyPressed(k)
	}
	g.LimiterDeplacement()

	return nil
}

// ---------------------------- ETATS: SELECTIONS ----------------------------

func ClasseSelection(g *Game) {
	switch {
	case g.JustPressed(ebiten.KeyN):
		classe = "MAGE"
		force, agilite, intelligence, endurance, mana, dexterite, level = 2, 1, 8, 2, 20, 5, 1
		chargerSprites(g)
	case g.JustPressed(ebiten.KeyE):
		classe = "ARCHER"
		force, agilite, intelligence, endurance, mana, dexterite, level = 2, 3, 7, 2, 15, 4, 1
		chargerSprites(g)
	case g.JustPressed(ebiten.KeyG):
		classe = "GUERRIER"
		force, agilite, intelligence, endurance, mana, dexterite, level = 8, 4, 3, 4, 5, 1, 1
		chargerSprites(g)
	case g.JustPressed(ebiten.KeyV):
		classe = "VOLEUR"
		force, agilite, intelligence, endurance, mana, dexterite, level = 3, 8, 4, 2, 5, 6, 1
		chargerSprites(g)
	case g.JustPressed(ebiten.KeyQ):
		classe = "ASSASSIN"
		force, agilite, intelligence, endurance, mana, dexterite, level = 3, 8, 2, 2, 5, 8, 1
		chargerSprites(g)
	default:
		return
	}
	state = StateArme
}

func ArmeSelection(g *Game) {
	switch {
	case g.JustPressed(ebiten.KeyE):
		toto = "Épée"
		force += 5
		degats += 3
	case g.JustPressed(ebiten.KeyB):
		toto = "Baguette"
		mana += 10
		degats += 2
	case g.JustPressed(ebiten.KeyQ):
		toto = "Arc"
		agilite += 2
		degats += 3
	case g.JustPressed(ebiten.KeyC):
		toto = "Couteau de Cuistôt"
		force += 1
		degats += 1
	case g.JustPressed(ebiten.KeyD):
		toto = "Dague"
		dexterite += 2
		degats += 3
	}
	if toto != "" {
		arme = append(arme, toto)
		state = StateJeu
	}
}

// ---------------------------- ETAT: JEU PRINCIPAL ----------------------------

func Jeu(g *Game) {
	switch {
	case g.JustPressed(ebiten.KeyL):
		state = StateLieu
	case g.JustPressed(ebiten.KeyS):
		state = StateStat
	case g.JustPressed(ebiten.KeyI):
		prevState = StateJeu
		state = StateInventaire
	case g.JustPressed(ebiten.KeyB):
		state = StateBoutique
	case g.JustPressed(ebiten.KeyF5):
		// Sauvegarder la partie
		err := sauvegarderJeu("save.json")
		if err != nil {
			log.Printf("❌ Erreur de sauvegarde: %v\n", err)
		}
	case g.JustPressed(ebiten.KeyF9):
		// Charger la partie
		err := chargerJeu("save.json")
		if err != nil {
			log.Printf("❌ Erreur de chargement: %v\n", err)
		}
	case g.JustPressed(ebiten.KeyQ):
		state = StateQuetes
	case g.JustPressed(ebiten.KeyT):
		state = StateTalents
	case g.JustPressed(ebiten.KeyA):
		state = StateAchievements
	case g.JustPressed(ebiten.KeyP):
		state = StatePause
	case g.JustPressed(ebiten.KeyT):
		state = StateThemeG
	case g.JustPressed(ebiten.KeyC):
		state = StateCredit
	case g.JustPressed(ebiten.KeyE):
		state = StateEquipement

	case g.JustPressed(ebiten.KeyNumLock):
		force = 10000000
		agilite = 1000000
		endurance = 100000
		level = 10000
		argent = 10000000
		dexterite = 1000000
		intelligence = 1000000
		mana = 100000

		state = StateJeu
	}
}

// ---------------------------- ETATS: SOUS-MENUS ----------------------------

func UpdateEquipementMenu(g *Game) {
	if g.JustPressed(ebiten.KeyEscape) {
		state = StateJeu
		return
	}

	if g.JustPressed(ebiten.Key1) {
		toggleEquipement("Casque")
	}
	if g.JustPressed(ebiten.Key2) {
		toggleEquipement("Plastron")
	}
	if g.JustPressed(ebiten.Key3) {
		toggleEquipement("Bottes")
	}
	if g.JustPressed(ebiten.Key4) {
		toggleEquipement("Arme")
	}
	if g.JustPressed(ebiten.Key5) {
		toggleEquipement("Anneau")
	}
}

func toggleEquipement(slot string) {
	if equipSlots.Get(slot) != nil {
		Desequiper(slot)
	} else {
		EquiperObjetSlot(slot, "Équipement de base")
	}
}

func (s *Slots) Get(slot string) *Equipement {
	switch strings.ToLower(slot) {
	case "casque":
		return s.Casque
	case "armure", "plastron":
		return s.Armure
	case "anneau":
		return s.Anneau
	case "arme":
		return s.Arme
	case "bottes":
		return s.Bottes
	default:
		return nil
	}
}

func (s *Slots) Set(slot string, eq *Equipement) {
	switch strings.ToLower(slot) {
	case "casque":
		s.Casque = eq
	case "armure", "plastron":
		s.Armure = eq
	case "anneau":
		s.Anneau = eq
	case "arme":
		s.Arme = eq
	case "bottes":
		s.Bottes = eq
	}
}

func EquiperObjetSlot(slot, nom string) {
	// crée un Equipement minimal avec un Type cohérent avec le slot
	eq := &Equipement{
		Nom: nom,
		Type: map[string]string{
			"casque":   "Casque",
			"armure":   "Armure",
			"plastron": "Armure",
			"anneau":   "Anneau",
			"arme":     "Arme",
			"bottes":   "Bottes",
		}[strings.ToLower(slot)],
	}

	// positionne dans equipSlots
	equipSlots.Set(slot, eq)

	// applique les effets/bonus côté stats (utilise ta table equipementBonus via le nom)
	EquiperObjet(slot, nom)
}

func Boutique(g *Game) {
	if g.JustPressed(ebiten.KeyEscape) {
		state = StateJeu
		return
	}

	// --- PAGE 1 : armes + potions ---
	if pageBoutique == 1 {
		var armesClasse []Arme
		switch classe {
		case "MAGE":
			armesClasse = armesMage
		case "GUERRIER":
			armesClasse = armesGuerrier
		case "VOLEUR":
			armesClasse = armesVoleur
		case "ASSASSIN":
			armesClasse = armesAssassin
		case "ARCHER":
			armesClasse = armesARCHER
		default:
			armesClasse = armesDispo
		}

		// potions disponibles
		armesClasse = append(armesClasse,
			Arme{"Potion de Soin", 50, 0},
			Arme{"Potion de Soin Majeure", 120, 0},
			Arme{"Potion de Mana", 40, 0},
			Arme{"Potion de Poison", 30, 0},
			Arme{"Potion d'XP", 100, 0},
		)

		// mapping dynamique des touches de sélection (déclaré ici pour Update)
		keys := []ebiten.Key{ebiten.KeyC, ebiten.KeyF, ebiten.KeyO, ebiten.KeyP, ebiten.KeyS, ebiten.KeyE}
		for i := 0; i < len(armesClasse) && i < len(keys); i++ {
			if g.JustPressed(keys[i]) {
				armeSelectionnee = armesClasse[i]
				state = StateConfirmation
				return
			}
		}
	}

	// --- PAGE 2 : équipements ---
	if pageBoutique == 2 {
		keys := []ebiten.Key{ebiten.KeyC, ebiten.KeyF, ebiten.KeyO, ebiten.KeyP, ebiten.KeyS}
		for i := 0; i < len(shopEquipements) && i < len(keys); i++ {
			if g.JustPressed(keys[i]) {
				armeSelectionnee = Arme{shopEquipements[i].Nom, float64(shopEquipements[i].Prix), 0}
				switch shopEquipements[i].Nom {
				case "Casque en fer":
					if equipSlots.Casque != nil {
						mesEquipements = append(mesEquipements, *equipSlots.Casque)
					}
				case "Plastron de cuir ":
					if equipSlots.Armure != nil {
						mesEquipements = append(mesEquipements, *equipSlots.Armure)
					}
				case "Anneau magique":
					if equipSlots.Anneau != nil {
						mesEquipements = append(mesEquipements, *equipSlots.Anneau)
					}
				case "Bottes renforcées":
					if equipSlots.Bottes != nil {
						mesEquipements = append(mesEquipements, *equipSlots.Bottes)
					}
				}
				state = StateConfirmation
				return
			}
		}
	}
	// --- Navigation entre pages ---
	if g.JustPressed(ebiten.Key1) {
		pageBoutique++
		if pageBoutique > 2 { // 3 pages : armes/potions, équipements shop
			pageBoutique = 1
		}
	}
	if g.JustPressed(ebiten.KeyV) {
		state = StateVendre
	}
}

func Confirmation(g *Game) {
	if g.JustPressed(ebiten.KeyV) {
		acheterArme(armeSelectionnee)
		// Proposer d'équiper directement si c'est un équipement
		if estEquipement(armeSelectionnee.Nom) {
			state = StateProposerEquiper
		}
	} else if g.JustPressed(ebiten.KeyR) {
		state = StateBoutique
	}
}

func estEquipement(nom string) bool {
	// Vérifier si c'est un équipement plutôt qu'une potion
	return !strings.Contains(nom, "Potion")
}

func ProposerEquiper(g *Game) {
	if g.JustPressed(ebiten.KeyE) {
		// Équiper l'objet acheté
		equiperObjetAchete(armeSelectionnee.Nom)
		state = StateJeu
	} else if g.JustPressed(ebiten.KeyN) {
		state = StateJeu
	}
}

func equiperObjetAchete(nom string) {
	// Créer l'équipement et l'équiper directement
	eq := creerEquipement(nom, armeSelectionnee.Degat)

	// Déterminer le slot approprié
	var slot string
	switch {
	case strings.Contains(nom, "Casque"):
		slot = "Casque"
	case strings.Contains(nom, "Plastron") || strings.Contains(nom, "Armure"):
		slot = "Plastron"
	case strings.Contains(nom, "Bottes"):
		slot = "Bottes"
	case strings.Contains(nom, "Anneau"):
		slot = "Anneau"
	case strings.Contains(nom, "Épée") || strings.Contains(nom, "Hache") || strings.Contains(nom, "Arc") || strings.Contains(nom, "Baguette") || strings.Contains(nom, "Dague"):
		slot = "Arme"
	default:
		slot = "Arme"
	}

	// Déséquiper l'ancien si présent
	if equipSlots.Get(slot) != nil {
		Desequiper(slot)
	}

	// Équiper le nouvel objet
	equipSlots.Set(slot, eq)
	appliquerBonusEquipement(eq)
}

// ---------------------------- SYSTEME DE POPUP ----------------------------

func AfficherPopup(typePopup, message string) {
	popupActive = true
	popupType = typePopup
	popupMessage = message
	popupTimer = 120 // 2 secondes à 60 FPS
	popupAlpha = 0
	popupScale = 0.5
}

func UpdatePopup() {
	if !popupActive {
		return
	}

	// Animation d'entrée
	if popupTimer > 90 {
		popupAlpha += 0.08
		popupScale += 0.05
		if popupAlpha > 1.0 {
			popupAlpha = 1.0
		}
		if popupScale > 1.0 {
			popupScale = 1.0
		}
	}

	// Animation de sortie
	if popupTimer < 30 {
		popupAlpha -= 0.03
		if popupAlpha < 0 {
			popupAlpha = 0
		}
	}

	popupTimer--
	if popupTimer <= 0 {
		popupActive = false
		popupAlpha = 0
		popupScale = 1.0
	}
}

func DrawPopup(screen *ebiten.Image) {
	if !popupActive || popupAlpha <= 0 {
		return
	}

	// Fond semi-transparent
	overlay := ebiten.NewImage(WinW, WinH)
	overlay.Fill(color.RGBA{0, 0, 0, uint8(100 * popupAlpha)})
	screen.DrawImage(overlay, nil)

	// Couleur selon le type
	var popupColor color.RGBA
	var titre string

	switch popupType {
	case "victoire":
		popupColor = color.RGBA{100, 255, 100, uint8(255 * popupAlpha)}
		titre = "✦ VICTOIRE ✦"
	case "defaite":
		popupColor = color.RGBA{255, 100, 100, uint8(255 * popupAlpha)}
		titre = "☠ DEFAITE ☠"
	case "levelup":
		popupColor = color.RGBA{255, 215, 0, uint8(255 * popupAlpha)}
		titre = "★ NIVEAU SUPERIEUR ★"
	case "quete":
		popupColor = color.RGBA{100, 200, 255, uint8(255 * popupAlpha)}
		titre = "✓ QUETE COMPLETEE ✓"
	case "achievement":
		popupColor = color.RGBA{255, 165, 0, uint8(255 * popupAlpha)}
		titre = "⚡ SUCCES DEBLOQUE ⚡"
	default:
		popupColor = color.RGBA{255, 255, 255, uint8(255 * popupAlpha)}
		titre = "NOTIFICATION"
	}

	// Cadre de popup centré avec effet de scale
	boxW := 800.0
	boxH := 300.0
	boxX := float64(WinW)/2 - (boxW*popupScale)/2
	boxY := float64(WinH)/2 - (boxH*popupScale)/2

	// Fond du popup
	popupBox := ebiten.NewImage(int(boxW), int(boxH))
	popupBox.Fill(color.RGBA{20, 20, 40, uint8(230 * popupAlpha)})

	// Bordure dorée
	for i := 0; i < 5; i++ {
		drawRect(popupBox, i, i, int(boxW)-i*2, int(boxH)-i*2, color.RGBA{255, 215, 0, uint8(200 * popupAlpha)})
	}

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(popupScale, popupScale)
	opts.GeoM.Translate(boxX, boxY)
	screen.DrawImage(popupBox, opts)

	// Texte (sans scale pour rester lisible)
	textY := int(boxY) + 80
	drawCenteredTextWithColor(screen, titre, textY, medievalFont, popupColor)

	whiteColor := color.RGBA{255, 255, 255, uint8(255 * popupAlpha)}
	drawCenteredTextWithColor(screen, popupMessage, textY+80, medievalFont, whiteColor)
}

func drawRect(img *ebiten.Image, x, y, w, h int, c color.Color) {
	for i := x; i < x+w; i++ {
		img.Set(i, y, c)
		img.Set(i, y+h-1, c)
	}
	for i := y; i < y+h; i++ {
		img.Set(x, i, c)
		img.Set(x+w-1, i, c)
	}
}

func drawCenteredTextWithColor(screen *ebiten.Image, str string, y int, f font.Face, c color.Color) {
	bounds := text.BoundString(f, str)
	x := (WinW - bounds.Dx()) / 2
	text.Draw(screen, str, f, x, y, c)
}

// ---------------------------- IA ENNEMI AMELIOREE ----------------------------

func UpdateEnnemiIA() {
	if !combatInit || tourDuJoueur {
		return
	}

	// Cooldown entre les actions
	if ennemiCooldown > 0 {
		ennemiCooldown--
		return
	}

	// Distance entre ennemi et joueur (simulation)
	distance := ennemiPosX - 500 // Joueur supposé à x=500

	// Choisir une action selon le type d'ennemi et la situation
	actionRoll := rand.Intn(100)

	// Boss et ennemis forts ont des comportements spéciaux
	isBoss := strings.Contains(ennemiNom, "Boss") || strings.Contains(ennemiNom, "Maître")

	if isBoss {
		// Boss plus agressif et tactique
		if ennemiPV < ennemiPVMax*0.3 {
			// En dessous de 30% HP : se défendre ou attaque désespérée
			if actionRoll < 40 {
				ennemiAction = "défendre"
				ennemiDefense = true
			} else {
				ennemiAction = "compétence"
			}
		} else if distance > 200 {
			// Loin : s'approcher
			ennemiAction = "marcher"
		} else if actionRoll < 60 {
			// Près : attaquer souvent
			ennemiAction = "attaquer"
		} else {
			// Compétence spéciale
			ennemiAction = "compétence"
		}
	} else {
		// Ennemis normaux : comportement plus simple
		if distance > 300 {
			ennemiAction = "marcher"
		} else if actionRoll < 70 {
			ennemiAction = "attaquer"
		} else if actionRoll < 85 {
			ennemiAction = "défendre"
			ennemiDefense = true
		} else {
			ennemiAction = "compétence"
		}
	}

	// Exécuter l'action (sera géré dans Combat)
	ennemiCooldown = 30 + rand.Intn(30) // 0.5-1 seconde entre actions
}

func ExecuterActionEnnemi() string {
	switch ennemiAction {
	case "marcher":
		// Mouvement vers le joueur
		ennemiPosX -= ennemiVitesse * 3
		if ennemiPosX < 600 {
			ennemiPosX = 600 // Distance minimale
		}
		return fmt.Sprintf("%s s'approche !", ennemiNom)

	case "défendre":
		// Réduire les dégâts du prochain coup
		ennemiDefense = true
		return fmt.Sprintf("%s se met en position défensive !", ennemiNom)

	case "compétence":
		// Attaque spéciale selon le type
		var dmg float64
		var message string

		switch ennemiTypeLoot {
		case "bete":
			// Rugissement : dégâts + peur
			dmg = (ennemiDegMin + ennemiDegMax) / 2 * 1.5
			message = fmt.Sprintf("%s pousse un rugissement terrifiant ! -%.0f PV", ennemiNom, dmg)
		case "humanoide":
			// Attaque tactique : dégâts précis
			dmg = ennemiDegMax * 1.3
			message = fmt.Sprintf("%s lance une attaque calculée ! -%.0f PV", ennemiNom, dmg)
		default:
			dmg = (ennemiDegMin + ennemiDegMax) * 0.8
			message = fmt.Sprintf("%s utilise une compétence ! -%.0f PV", ennemiNom, dmg)
		}

		// Appliquer défense joueur
		reduction := defenseTotal * 0.5
		if reduction > dmg*0.75 {
			reduction = dmg * 0.75
		}
		dmg -= reduction
		if dmg < 0 {
			dmg = 0
		}

		pv -= dmg
		return message

	case "attaquer":
		// Attaque normale (géré dans Combat())
		return ""

	default:
		return ""
	}
}

func Resultat(g *Game) {
	if dernierLoot == "" {
		loots := lootTable[ennemiNom]
		if len(loots) > 0 {
			drop := loots[rand.Intn(len(loots))]
			arme = append(arme, drop)
			dernierLoot = drop
		}
	}
	if g.JustPressed(ebiten.KeySpace) {
		levelUpMsg = ""
		dernierLoot = "" // reset
		state = StateJeu
	}
}

func Vente(g *Game) {
	keys := []ebiten.Key{
		ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4, ebiten.Key5,
		ebiten.Key6, ebiten.Key7, ebiten.Key8, ebiten.Key9,
	}
	for i := 0; i < len(arme) && i < len(keys); i++ {
		if g.JustPressed(keys[i]) {
			vendreArme(i, 10)
		}
	}
	if g.JustPressed(ebiten.KeyB) {
		state = StateBoutique
	}
}

func Lieu(g *Game) {
	if g.JustPressed(ebiten.KeyC) {
		lieuactu = "AUX CHAMPS DE BLÉ"
		state = StateChamps
		chargerSprites(g)
		return
	}
	if g.JustPressed(ebiten.KeyF) {
		lieuactu = "À LA FORGE"
		state = StateForge
		chargerSprites(g)
		return
	}
	if g.JustPressed(ebiten.KeyT) {
		lieuactu = "À LA TOUR DE SORCIER"
		state = StateTour
		chargerSprites(g)
		return
	}
	if g.JustPressed(ebiten.KeyS) {
		lieuactu = "SORTIR DU VILLAGE"
		state = StateSortie
		chargerSprites(g)
		return
	}
	if g.JustPressed(ebiten.KeyH) {
		lieuactu = "RENTRER À L'HÔTEL DE DUDEBU"
		state = StateHotel
		chargerSprites(g)
		return
	}
	if g.JustPressed(ebiten.KeyEscape) {
		state = StateJeu
		return
	}
}

func ThemeG(g *Game) {
	switch {
	case g.JustPressed(ebiten.KeySpace):
		state = StateThemeA
	}
}

func ThemeA(g *Game) {
	switch {
	case g.JustPressed(ebiten.KeySpace):
		state = StateThemeM
	}
}

func ThemeM(g *Game) {
	switch {
	case g.JustPressed(ebiten.KeySpace):
		state = StateThemeE
	}
}

func ThemeE(g *Game) {
	switch {
	case g.JustPressed(ebiten.KeySpace):
		state = StateThemeV
	}
}

func ThemeV(g *Game) {
	switch {
	case g.JustPressed(ebiten.KeySpace):
		state = StateJeu
	}
}

func Credit(g *Game) {
	switch {
	case g.JustPressed(ebiten.KeySpace):
		state = StateJeu
	}
}

func Sortie(g *Game) {
	switch {
	case g.JustPressed(ebiten.KeyC):
		lieuactu = "DANS LA FORÊT DES MONSTRES"
		state = StateForet
	case g.JustPressed(ebiten.KeyR):
		lieuactu = "AU VILLAGE DU DUDEBU"
		state = StateJeu
	case g.JustPressed(ebiten.KeyEscape):
		state = StateLieu
	}
}

func Histoire1(g *Game) {
	if g.JustPressed(ebiten.KeySpace) {
		state = StateHotel
	}
}

func Histoire2(g *Game) {
	if g.JustPressed(ebiten.KeySpace) {
		state = StateTour
	}
}

func Histoire3(g *Game) {
	if g.JustPressed(ebiten.KeySpace) {
		state = StateForge
	}
}

func Forge(g *Game) {
	switch {
	case g.JustPressed(ebiten.KeyL):
		intelligence += 0.2
	case g.JustPressed(ebiten.KeyC):
		force += 0.2
	case g.JustPressed(ebiten.KeyN):
		dexterite += 0.1
	case g.JustPressed(ebiten.KeyD):
		state = StateHistoire3
	case g.JustPressed(ebiten.KeyEscape):
		state = StateJeu
	}
}

func Pause(g *Game) {
	if g.JustPressed(ebiten.KeyJ) {
		state = StateJeu
		if len(titre) > 1 {
			return
		}
		titre = ""
		titre += "REPOS DU HÉROS"
	}
}

func Hotel(g *Game) {
	tata, tota, tato, tto, tta = 0, 0, 0, 0, 0
	if g.JustPressed(ebiten.KeyEscape) {
		tata = 0
		tota = 0
		tato = 0
		state = StateJeu
		return
	}
	switch {
	case g.JustPressed(ebiten.KeyL):
		tata += 0.05
		intelligence += tata
		state = StateHotel
	case g.JustPressed(ebiten.KeyC):
		state = StateHistoire1
	case g.JustPressed(ebiten.KeyN):
		tota += 2
		argent -= tota
		state = StateHotel
	case g.JustPressed(ebiten.KeyD):
		tato += 1
		argent += tato
		state = StateHotel
	}

}

func Champs(g *Game) {
	tata, tota, tato, tto, tta = 0, 0, 0, 0, 0
	if g.JustPressed(ebiten.KeyEscape) {
		tata = 0
		tota = 0
		tato = 0
		tto = 0
		tta = 0

		state = StateJeu
		return
	}
	switch {
	case g.JustPressed(ebiten.KeyL):
		tata += 0.1
		force += tata
		state = StateChamps
	case g.JustPressed(ebiten.KeyC):
		tota += 0.1
		agilite += tota
		state = StateChamps
	case g.JustPressed(ebiten.KeyV):
		tato += 0.5
		argent += tato
		state = StateChamps
	case g.JustPressed(ebiten.KeyD):
		tto += 0.08
		intelligence += tto
		state = StateChamps
	case g.JustPressed(ebiten.KeyN):
		tta += 0.1
		endurance += tta
		state = StateChamps
	}
}

func Tour(g *Game) {
	tata, tota, tato, tto, tta = 0, 0, 0, 0, 0
	if g.JustPressed(ebiten.KeyEscape) {
		tata = 0
		tota = 0
		tato = 0
		tto = 0
		tta = 0

		state = StateJeu
		return
	}
	switch {
	case g.JustPressed(ebiten.KeyL):
		tto += 0.03
		intelligence += tto
	case g.JustPressed(ebiten.KeyC):
		tata += 0.3
		dexterite += tata
	case g.JustPressed(ebiten.KeyM):
		tota += 0.4
		mana += tota
	case g.JustPressed(ebiten.KeyD):
		state = StateHistoire2

	}
}

// ---------------------------- INVENTAIRE / SHOP ----------------------------

func acheterArme(a Arme) {
	if len(arme) >= 10 {
		state = StateAttention
		return
	}
	if argent < a.Prix {
		state = StateAttentionPrix
		return
	}
	argent -= a.Prix

	switch a.Nom {
	case "Potion de Soin", "Potion de Soin Majeure", "Potion de Mana", "Potion de Poison", "Potion d'XP":
		arme = append(arme, a.Nom) // item consommable
	default:
		// arme/équipement - juste ajouter à l'inventaire sans modifier les stats
		arme = append(arme, a.Nom)
		toto = a.Nom
	}
	state = StateBoutique
}

func creerEquipement(nom string, degat float64) *Equipement {
	eq := &Equipement{Nom: nom, Degat: degat}

	// Déterminer le type et les bonus selon le nom
	if strings.Contains(nom, "Épée") || strings.Contains(nom, "Hache") || strings.Contains(nom, "Dague") || strings.Contains(nom, "Arc") || strings.Contains(nom, "Baguette") || strings.Contains(nom, "Bâton") || strings.Contains(nom, "Lance") || strings.Contains(nom, "Kris") || strings.Contains(nom, "Arbalète") || strings.Contains(nom, "Lames") {
		eq.Type = "arme"
		eq.BonusForce = degat * 0.5
		eq.BonusDex = degat * 0.3
	} else if strings.Contains(nom, "Bouclier") || strings.Contains(nom, "Armure") || strings.Contains(nom, "Casque") {
		eq.Type = "armure"
		eq.Defense = degat * 2
		eq.BonusPV = degat * 5
	} else if strings.Contains(nom, "Grimoire") || strings.Contains(nom, "Anneau") || strings.Contains(nom, "Amulette") {
		eq.Type = "accessoire"
		eq.BonusInt = degat * 0.8
		eq.BonusMana = degat * 2
	} else {
		eq.Type = "arme"
	}

	return eq
}

func equiperItem(nom string) bool {
	// Trouver l'arme dans la liste des armes disponibles
	var armeData *Arme
	toutesArmes := [][]Arme{armesGuerrier, armesMage, armesVoleur, armesAssassin, armesARCHER}
	for _, liste := range toutesArmes {
		for _, a := range liste {
			if a.Nom == nom {
				armeData = &a
				break
			}
		}
		if armeData != nil {
			break
		}
	}

	if armeData == nil {
		return false
	}

	eq := creerEquipement(armeData.Nom, armeData.Degat)

	// Déséquiper l'ancien item du même type
	switch eq.Type {
	case "arme":
		if equipementArme != nil {
			retirerBonusEquipement(equipementArme)
		}
		equipementArme = eq
	case "armure":
		if equipementArmure != nil {
			retirerBonusEquipement(equipementArmure)
		}
		equipementArmure = eq
	case "accessoire":
		if equipementAccessoire != nil {
			retirerBonusEquipement(equipementAccessoire)
		}
		equipementAccessoire = eq
	}

	appliquerBonusEquipement(eq)
	return true
}

func appliquerBonusEquipement(eq *Equipement) {
	if eq == nil {
		return
	}
	degats += eq.Degat
	defenseTotal += eq.Defense
	force += eq.BonusForce
	dexterite += eq.BonusDex
	intelligence += eq.BonusInt
	mana += eq.BonusMana
	pvMAX += eq.BonusPV
	if pv > pvMAX {
		pv = pvMAX
	}
}

func retirerBonusEquipement(eq *Equipement) {
	if eq == nil {
		return
	}
	degats -= eq.Degat
	defenseTotal -= eq.Defense
	force -= eq.BonusForce
	dexterite -= eq.BonusDex
	intelligence -= eq.BonusInt
	mana -= eq.BonusMana
	pvMAX -= eq.BonusPV
	if pv > pvMAX {
		pv = pvMAX
	}
}

func vendreArme(index int, prix float64) {
	if index < 0 || index >= len(arme) {
		state = StateAttention
		return
	}
	arme = append(arme[:index], arme[index+1:]...)
	argent += prix
}

func _vendreEquipement(nom string) {
	for i, eq := range mesEquipements {
		if eq.Nom == nom {
			argent += float64(eq.Prix) / 2
			mesEquipements = append(mesEquipements[:i], mesEquipements[i+1:]...)
			return
		}
	}
}

func UseItem(index int) bool {
	if index < 0 || index >= len(arme) {
		return false
	}
	it := arme[index]
	switch it {
	case "Potion de Soin":
		pv += 20
		if pv > pvMAX {
			pv = pvMAX
		}
		arme = append(arme[:index], arme[index+1:]...)
		return true
	case "Potion de Soin Majeure":
		pv += 50
		if pv > pvMAX {
			pv = pvMAX
		}
		arme = append(arme[:index], arme[index+1:]...)
		return true
	case "Potion de Mana":
		mana += 30
		if mana > 100 {
			mana = 100
		}
		arme = append(arme[:index], arme[index+1:]...)
		return true
	case "Potion d'XP":
		xp += 20
		arme = append(arme[:index], arme[index+1:]...)
		return true
	case "Potion de Poison":
		if state == StateCombat || prevState == StateCombat {
			dmg := 8.0 + level*1.5
			ennemiPV -= dmg
			arme = append(arme[:index], arme[index+1:]...)
			return true
		}
		return false
	default:
		// Essayer d'équiper l'item
		if equiperItem(it) {
			arme = append(arme[:index], arme[index+1:]...)
			return true
		}
		return false
	}
}

// Inventaire
func Inventaire(g *Game) {

	keys := []ebiten.Key{
		ebiten.Key0, ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4, ebiten.Key5,
		ebiten.Key6, ebiten.Key7, ebiten.Key8, ebiten.Key9,
	}
	for i := 0; i < len(keys) && i < len(arme); i++ {
		if g.JustPressed(keys[i]) {
			consumed := UseItem(i)

			if consumed && prevState == StateCombat {

				state = StateCombat
				tourDuJoueur = false
				return
			}
		}
	}

	if g.JustPressed(ebiten.KeyEscape) {
		if prevState != 0 {
			state = prevState
			prevState = 0
		} else {
			state = StateJeu
		}
	}
}

// ---------------------------- FORET / MONSTRES ----------------------------

func Foret(g *Game) {
	switch {
	case g.JustPressed(ebiten.KeyC):
		state = StateCombat
	case g.JustPressed(ebiten.KeyF):
		titre = ""
		titre += "BAH LE NULOS TU FUIS HEIN"
		combatsSansFuite = 0 // Réinitialiser la série de combats sans fuite
		state = StateJeu
	}
}

// ---------------------------- COMBAT ----------------------------

func InitialiserCombat(g *Game) {
	combatInit = true
	tourDuJoueur = true
	logCombat = logCombat[:0]

	if monstreIndex >= len(ordreMonstres1) {
		monstreIndex = len(ordreMonstres1) - 1
	}
	monster = ordreMonstres1[monstreIndex]
	monstreIndex++

	for _, m := range bestiaire {
		if strings.EqualFold(m.Nom, monster) {
			ennemiNom = strings.ToUpper(m.Nom)

			scale := 1.0
			switch strings.ToLower(monster) {
			case "Crabauge":
				scale = 1.0
			case "Vorlapin":
				scale = 2
			case "Gobelin":
				scale = 3
			case "Boss Lycaon":
				scale = 4
			case "Muddig":
				scale = 5
			case "Gros Serpent Vorace":
				scale = 6
			case "Serpent Livestide":
				scale = 6.5
			case "Loosers Wood":
				scale = 7
			case "Boss Wezaemon":
				scale = 7.5
			case "Poiscalle Zombie":
				scale = 7.5
			case "Lugia":
				scale = 8
			default:

				scale = 1.0 + float64(monstreIndex)/10.0
			}
			ennemiPVMax = m.PVMax * scale
			ennemiPV = ennemiPVMax
			ennemiDegMin = m.DegatsMin * scale
			ennemiDegMax = m.DegatsMax * scale
			ennemiColor = m.Couleur
			ennemiXPGainMin = int(float64(m.XPGagneMin) * scale)
			ennemiXPGainMax = int(float64(m.XPGagneMax) * scale)
			if ennemiXPGainMin < 1 {
				ennemiXPGainMin = 1
			}
			if ennemiXPGainMax < ennemiXPGainMin {
				ennemiXPGainMax = ennemiXPGainMin
			}
			ennemiTypeLoot = m.LootCategorie
			break
		}
	}

	chargerSprites(g)
	logCombat = append(logCombat, fmt.Sprintf("Un %s surgit !", ennemiNom))
}

func Combat(g *Game) {

	if g.JustPressed(ebiten.KeyEnter) && fuitePossible {
		logCombat = append(logCombat, "Tu fuis le combat...")
		combatInit = false

		if monstreIndex > 0 {
			monstreIndex--
		}
		state = StateForet
		return
	}

	if tourDuJoueur {

		if g.JustPressed(ebiten.KeyQ) {
			min := 6.0 + force*0.3 + degats*0.8 + dexterite*0.1 + mana*0.4
			max := 12.0 + force*0.6 + degats*1.2 + dexterite*0.2 + mana*0.8
			dmg := randFloat(min, max)
			if crit(10) {
				dmg *= 1.5
				logCombat = append(logCombat, fmt.Sprintf("Coup critique ! -%.0f PV à %s", dmg, ennemiNom))
			} else {
				logCombat = append(logCombat, fmt.Sprintf("Attaque légère ! -%.0f PV à %s", dmg, ennemiNom))
			}
			ennemiPV -= dmg
			g.MonstreDamageTimer = 30
			tourDuJoueur = false
			return
		}
		if g.JustPressed(ebiten.KeyW) {
			if miss(25) {
				logCombat = append(logCombat, "Attaque forte ratée !")
			} else {
				min := 14.0 + force*0.7 + degats*1.0 + mana*0.4 + intelligence*0.2 + agilite*0.3 + endurance*0.2
				max := 26.0 + force*1.0 + degats*1.6 + mana*0.8 + intelligence*0.4 + agilite*0.6 + endurance*0.4
				dmg := randFloat(min, max)
				if crit(15) {
					dmg *= 1.5
					logCombat = append(logCombat, fmt.Sprintf("Attaque forte critique ! -%.0f PV", dmg))
				} else {
					logCombat = append(logCombat, fmt.Sprintf("Attaque forte ! -%.0f PV", dmg))
				}
				ennemiPV -= dmg
				g.MonstreDamageTimer = 30
			}
			tourDuJoueur = false
			return
		}
		if g.JustPressed(ebiten.KeyE) {
			skills := competences[classe]
			if len(skills) == 0 {
				logCombat = append(logCombat, "Pas de compétence disponible !")
				return
			}
			skill := skills[rand.Intn(len(skills))]
			if mana < skill.CoutMana {
				logCombat = append(logCombat, "Pas assez de mana !")
				return
			}
			mana -= skill.CoutMana
			base := force*0.5 + intelligence*0.8 + degats*1.0 + mana*0.5
			dmg := randFloat(base*0.8, base*skill.MultDegats)
			ennemiPV -= dmg
			g.MonstreDamageTimer = 30
			logCombat = append(logCombat, fmt.Sprintf("%s utilise %s ! -%.0f PV (%s)", classe, skill.Nom, dmg, skill.Effet))
			tourDuJoueur = false
			return
		}
	} else {
		// IA Ennemi : mise à jour et exécution d'action
		UpdateEnnemiIA()

		// Si l'action n'est pas "attaquer", exécuter l'action spéciale
		if ennemiAction != "attaquer" && ennemiAction != "" {
			actionMsg := ExecuterActionEnnemi()
			if actionMsg != "" {
				logCombat = append(logCombat, actionMsg)
			}
			g.DamageTimer = 5
			tourDuJoueur = true
			ennemiAction = "" // Reset action
			return
		}

		// Attaque normale
		min := ennemiDegMin
		max := ennemiDegMax
		dmg := randFloat(min, max)

		// Si ennemi en défense, réduire ses dégâts mais améliorer sa défense
		if ennemiDefense {
			dmg *= 0.6 // -40% dégâts en position défensive
			ennemiDefense = false
		}

		// Réduction des dégâts par la défense
		reduction := defenseTotal * 0.5
		if reduction > dmg*0.75 {
			reduction = dmg * 0.75 // Max 75% de réduction
		}
		dmg -= reduction
		if dmg < 0 {
			dmg = 0
		}

		switch ennemiTypeLoot {
		case "bete":
			if crit(8) {
				bleed := 4.0 + level
				if defenseTotal > 0 {
					logCombat = append(logCombat, fmt.Sprintf("%s inflige une morsure ! -%.0f PV (%.0f bloqués) + saignement -%.0f", ennemiNom, dmg, reduction, bleed))
				} else {
					logCombat = append(logCombat, fmt.Sprintf("%s inflige une morsure ! -%.0f PV + saignement -%.0f", ennemiNom, dmg, bleed))
				}
				pv -= dmg + bleed
			} else {
				if defenseTotal > 0 {
					logCombat = append(logCombat, fmt.Sprintf("%s attaque ! -%.0f PV (%.0f bloqués)", ennemiNom, dmg, reduction))
				} else {
					logCombat = append(logCombat, fmt.Sprintf("%s attaque ! -%.0f PV", ennemiNom, dmg))
				}
				pv -= dmg
			}
		default:
			if defenseTotal > 0 {
				logCombat = append(logCombat, fmt.Sprintf("%s attaque ! -%.0f PV (%.0f bloqués)", ennemiNom, dmg, reduction))
			} else {
				logCombat = append(logCombat, fmt.Sprintf("%s attaque ! -%.0f PV", ennemiNom, dmg))
			}
			pv -= dmg
		}
		g.DamageTimer = 5
		tourDuJoueur = true
		ennemiAction = "" // Reset action
	}
}

func MettreAJourCombat() {
	if ennemiPV <= 0 {
		logCombat = append(logCombat, fmt.Sprintf("%s vaincu !", ennemiNom))
		titre = ""
		titre = fmt.Sprintf("TUEUR DE %s", ennemiNom)
		combatInit = false

		// Si on vient de tuer le Boss Lycaon, passer à la quête 2

		switch ennemiNom {
		case "BOSS LYCAON":
			queteActuelle = 2
			// Charger un nouveau décor pour la quête 2
			InitQuete("image/quete2.png")

		// Si on vient de tuer le Boss Wezaemon, passer à la quête 3
		case "BOSS WEZAEMON":
			queteActuelle = 3
			InitQuete("image/quete3.png")

		}
		// Gain d'XP
		xpg := ennemiXPGainMin
		if diff := ennemiXPGainMax - ennemiXPGainMin; diff > 0 {
			xpg += rand.Intn(diff + 1)
		}
		xp += float64(xpg)

		// Tracking pour quêtes et achievements
		monstresTotalTues++
		combatsSansFuite++

		// Level up éventuel
		levelUpMsg = ""
		for xp >= xpMax {
			xp -= xpMax
			level++
			xpMax += 20
			force += 10
			endurance += 10
			pvMAX += 20
			pv = pvMAX
			pointsTalents++ // Gagner un point de talent par niveau
			levelUpMsg = fmt.Sprintf("Niveau %.0f atteint !", level)

			// Popup de level up
			AfficherPopup("levelup", fmt.Sprintf("NIVEAU %.0f ! +10 Force, +10 Endurance, +20 PV Max", level))
		} // Loot
		lootItem, lootStats := GenererLoot(ennemiTypeLoot)
		AppliquerStatUp(lootStats)
		dernierLoot = lootItem // utile pour Draw()

		// Popup de victoire
		AfficherPopup("victoire", fmt.Sprintf("%s vaincu ! +%d XP", ennemiNom, xpg))

		state = StateVictoire
		return
	}

	if pv <= 0 {
		combatInit = false

		ennemiPV = ennemiPVMax
		if monstreIndex > 0 {
			monstreIndex--
		}

		// Popup de défaite
		AfficherPopup("defaite", "Vous avez été vaincu...")

		state = StateDefaite
		return
	}
}

// ---------------------------- LOOT / STATS ----------------------------

func GenererLoot(categorie string) (string, string) {
	roll := rand.Intn(100)
	cat := strings.ToLower(strings.TrimSpace(categorie))

	switch cat {
	case "bete", "monstre", "bête":
		switch {
		case roll < 35:
			argent += 200 + float64(rand.Intn(15))
			statLoot = append(statLoot, "Argent +100")
			return "Bourse de pièces", "Argent"
		case roll < 65:
			arme = append(arme, "Éclat de carapace")
			dexterite += 10
			statLoot = append(statLoot, "Dextérité +10")
			return "Éclat de carapace", "Dextérité"
		case roll < 85:
			pvMAX += 20
			pv += 20
			if pv > pvMAX {
				pv = pvMAX
			}
			statLoot = append(statLoot, "PV Max +20")
			return "Morceau de vitalité", "PV Max"
		default:
			degats += 10
			statLoot = append(statLoot, "Dégâts +10")
			return "Griffe affûtée", "Dégâts arme"
		}

	case "humanoide":
		switch {
		case roll < 40:
			argent += 200 + float64(rand.Intn(20))
			statLoot = append(statLoot, "Argent +150")
			return "Bourse de brigand", "Argent"
		case roll < 70:
			intelligence += 10
			statLoot = append(statLoot, "Intelligence +10")
			return "Page arcanique", "Intelligence"
		default:
			degats += 10
			dexterite += 10
			statLoot = append(statLoot, "Dégâts +10, Dextérité +10")
			return "Dague rouillée", "Dégâts arme"
		}

	case "rocheux":
		switch {
		case roll < 50:
			pvMAX += 80
			pv += 80
			if pv > pvMAX {
				pv = pvMAX
			}
			statLoot = append(statLoot, "PV Max +80")
			return "Fragment de pierre magique", "PV Max"
		default:
			force += 10
			statLoot = append(statLoot, "Force +10")
			return "Éclat granitique", "Force"
		}

	default:
		if roll < 50 {
			argent += 200
			statLoot = append(statLoot, "Argent +100")
			return "Vieilles pièces", "Argent"
		}
		degats += 10
		statLoot = append(statLoot, "Dégâts +10")
		return "Petit talon", "Dégâts arme"
	}
}

func Equiper(eq Equipement) {
	switch eq.Type {
	case "Casque":
		if equipSlots.Casque != nil {
			mesEquipements = append(mesEquipements, *equipSlots.Casque)

		}
		equipSlots.Casque = &eq
	case "Plastron":
		if equipSlots.Armure != nil {
			mesEquipements = append(mesEquipements, *equipSlots.Armure)

		}
		equipSlots.Armure = &eq
	case "Anneau":
		if equipSlots.Anneau != nil {
			mesEquipements = append(mesEquipements, *equipSlots.Anneau)

		}
		equipSlots.Anneau = &eq
	case "Arme":
		if equipSlots.Arme != nil {
			mesEquipements = append(mesEquipements, *equipSlots.Arme)

		}
		equipSlots.Arme = &eq
	case "Bottes":
		if equipSlots.Bottes != nil {
			mesEquipements = append(mesEquipements, *equipSlots.Bottes)

		}
		equipSlots.Bottes = &eq
	}
}

func Desequiper(slot string) {
	switch slot {
	case "Casque":
		if equipSlots.Casque != nil {
			mesEquipements = append(mesEquipements, *equipSlots.Casque)
			equipSlots.Casque = nil
		}
	case "Plastron":
		if equipSlots.Armure != nil {
			mesEquipements = append(mesEquipements, *equipSlots.Armure)
			equipSlots.Armure = nil
		}
	case "Anneau":
		if equipSlots.Anneau != nil {
			mesEquipements = append(mesEquipements, *equipSlots.Anneau)
			equipSlots.Anneau = nil
		}
	case "Arme":
		if equipSlots.Arme != nil {
			mesEquipements = append(mesEquipements, *equipSlots.Arme)

			equipSlots.Arme = nil
		}
	case "Bottes":
		if equipSlots.Bottes != nil {
			mesEquipements = append(mesEquipements, *equipSlots.Bottes)
			equipSlots.Bottes = nil
		}
	}
}

func slotToString(eq *Equipement) string {
	if eq == nil {
		return "vide"
	}
	return eq.Nom
}

func Victoire(g *Game) {
	if g.JustPressed(ebiten.KeySpace) {
		state = StateResultat
	}
}

func Defaite(g *Game) {
	if g.JustPressed(ebiten.KeyEscape) {
		state = StateBackLife
	}
}

func BackLife(g *Game) {
	if g.JustPressed(ebiten.KeyEscape) {
		perte := argent * 0.1
		argent -= perte
		if argent < 0 {
			argent = 0
		}
		log.Printf("Vous avez perdu %.2f OKSA en tombant au combat.", perte)
		pv = pvMAX * 0.5
		// ensure monster not advanced on death
		if monstreIndex > 0 {
			monstreIndex--
		}
		state = StateJeu
	}
}

func AppliquerStatUp(statu string) {
	switch classe {
	case "GUERRIER":
		dexterite += 8
		force += 10
		agilite += 6
		mana += 5
		intelligence += 2

		pvMAX += 20
		pv += 20
		if pv > pvMAX {
			pv = pvMAX
		}
		if statu == "Dégâts arme" {
			degats += 2
		}
	case "ASSASSIN":
		dexterite += 10
		force += 6
		agilite += 10
		mana += 5
		intelligence += 4

		pvMAX += 20
		pv += 20
		if pv > pvMAX {
			pv = pvMAX
		}

		if statu == "Dégâts arme" {
			degats += 2
		}
	case "MAGE":
		dexterite += 5
		force += 3
		agilite += 3
		mana += 25
		intelligence += 10

		pvMAX += 20
		pv += 20
		if pv > pvMAX {
			pv = pvMAX
		}

		if statu == "Dégâts arme" {
			degats += 2
		}
	case "VOLEUR":
		dexterite += 10
		force += 4
		agilite += 10
		mana += 5
		intelligence += 6

		pvMAX += 20
		pv += 20
		if pv > pvMAX {
			pv = pvMAX
		}
		if statu == "Dégâts arme" {
			degats += 2
		}

	case "ARCHER":
		dexterite += 7
		force += 4
		agilite += 12
		mana += 15
		intelligence += 6

		pvMAX += 20
		pv += 20
		if pv > pvMAX {
			pv = pvMAX
		}
		if statu == "Dégâts arme" {
			degats += 2
		}
	}
}
func EquiperObjet(slot string, nom string) {
	switch nom {
	case "Casque en fer":
		equipement.Casque = nom
	case "Plastron de cuir":
		equipement.Plastron = nom
	case "Bottes renforcées":
		equipement.Bottes = nom
	case "Anneau magique":
		equipement.Anneau = nom
	default:
		// pas un équipement, donc rien à équiper
		return
	}

	// Appliquer les bonus
	if bonus, ok := equipementBonus[nom]; ok {
		pvMAX += bonus.PV
		force += bonus.Force
		endurance += bonus.Endurance
		pv = pvMAX // soigner au max quand on équipe
	}
}

func (g *Game) LimiterDeplacement() {
	// largeur/hauteur du sprite (à ajuster si tu connais la taille exacte)
	spriteW := 64.0
	spriteH := 64.0

	// bornes horizontales
	if g.X < 0 {
		g.X = 0
	}
	if g.X > float64(WinW)-spriteW {
		g.X = float64(WinW) - spriteW
	}

	// bornes verticales
	if g.Y < 0 {
		g.Y = 0
	}
	if g.Y > float64(WinH)-spriteH {
		g.Y = float64(WinH) - spriteH
	}
}

// ---------------------------- DESSIN / UI ----------------------------

func drawCombatScene(g *Game, screen *ebiten.Image, bg color.RGBA) {
	drawFull(screen, fondqueteImage)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.X, g.Y)
	if g.DamageTimer > 0 {
		opts.ColorM.Scale(1, 0, 0, 1)
	}
	if g.JoueurImage != nil {
		srcX := 230
		srcY := 270
		srcRect := image.Rect(srcX, srcY, 0, 0)
		partieJoueur := g.JoueurImage.SubImage(srcRect).(*ebiten.Image)
		screen.DrawImage(partieJoueur, opts)
	}

	mopts := &ebiten.DrawImageOptions{}
	mopts.GeoM.Translate(g.MonstreX, g.MonstreY)
	if g.MonstreDamageTimer > 0 {
		mopts.ColorM.Scale(1, 0, 0, 1)
	}
	if g.MonstreImage != nil {
		srcX := 230
		srcY := 270
		srcRect := image.Rect(srcX, srcY, 0, 0)
		partieMonstre := g.MonstreImage.SubImage(srcRect).(*ebiten.Image)
		screen.DrawImage(partieMonstre, mopts)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch state {
	case StateAccueil:
		drawBlack(screen)
		drawCenteredText(screen, "UNE HISTOIRE FOLLE", 310, medievalFont, color.White)
		drawCenteredText(screen, "APPUYEZ SUR ESPACE POUR COMMENCER", 1000, medievalFont, color.White)
	case StateIntro:
		drawBlack(screen)
		drawCenteredText(screen, "UNE ÉPOQUE FUT CONNUE SOUS LE NOM D'ÈRE DES DIEUX.", 80, medievalFont, color.White)
		drawCenteredText(screen, "LES ILLUSTRES ÊTRES DIVINS ONT UN JOUR DISPARU,", 160, medievalFont, color.White)
		drawCenteredText(screen, "LAISSANT LEUR MONDE À LEUR POSTÉRITÉ.", 240, medievalFont, color.White)
		drawCenteredText(screen, "NOUS SOMMES LEURS HÉRITIERS SPIRITUELS.", 340, medievalFont, color.White)
		drawCenteredText(screen, "CONFORMÉMENT À LEUR SOUHAIT,", 440, medievalFont, color.White)
		drawCenteredText(screen, "NOUS NOUS SOMMES DISPERSÉS PARTOUT SUR TERRE POUR CONTINUER À PERPÉTUER LA VIE...", 540, medievalFont, color.White)
		drawCenteredText(screen, "NOUS SUIVONS AUJOURD'HUI LEURS TRACES GRÂCE AU PATRIMOINE", 640, medievalFont, color.White)
		drawCenteredText(screen, "DU PASSÉ DONT LE SOUFFLE EST TOUJOURS PRÉSENT DANS L'HISTOIRE", 740, medievalFont, color.White)
		drawCenteredText(screen, "ET SES VESTIGES.", 800, medievalFont, color.White)
		drawCenteredText(screen, "EXPLORATEUR OU EXPLORATRICE, C'EST LE VENT D'EST QUI T'A MENÉ ICI.", 860, medievalFont, color.White)
		drawCenteredText(screen, "TON DESTIN EST D'AVANCER, TE RELEVANT CHAQUE FOIS QUE TU T'ÉCROULERAS...", 960, medievalFont, color.White)
		drawCenteredText(screen, "APPUYEZ SUR ESPACE POUR CONTINUER", 1060, medievalFont, color.White)
	case StateStart:
		drawBlack(screen)
		drawCenteredText(screen, "ET PEUT-ÊTRE DISPARAÎTRAS-TU AVEC LUI.", 80, medievalFont, color.White)
		drawCenteredText(screen, "QUELLE EST TA RAISON D'ÊTRE ?", 160, medievalFont, color.White)
		drawCenteredText(screen, "CONFIERAS-TU TA VIE À TON ÉPÉE ?", 240, medievalFont, color.White)
		drawCenteredText(screen, "ESSAIERAS-TU D'ATTEINDRE LES HAUTEURS DE LA MAGIE ?", 340, medievalFont, color.White)
		drawCenteredText(screen, "APRÈS AVOIR JETÉ UN ŒIL À SES ABYSSES ?", 440, medievalFont, color.White)
		drawCenteredText(screen, "TU POURRAS AUSSI CHOISIR DE NE PAS EMPRUNTER LA VOIE DU COMBAT.", 540, medievalFont, color.White)
		drawCenteredText(screen, "TOUT EST ICI. TOUT EST EN TOI.", 640, medievalFont, color.White)
		drawCenteredText(screen, "VA, FAIS LE PREMIER PAS !", 740, medievalFont, color.White)
		drawCenteredText(screen, "EXPLORE L'INCONNU, L'AVENIR EST TON POTENTIEL.", 800, medievalFont, color.White)
		drawCenteredText(screen, "C'EST LA MISSION DONT HÉRITENT, DEPUIS L'ÈRE DES DIEUX,", 860, medievalFont, color.White)
		drawCenteredText(screen, "TOUS LES ENFANTS NÉS DANS CETTE CONTRÉE !", 960, medievalFont, color.White)
		drawCenteredText(screen, "APPUYEZ SUR ESPACE POUR CONTINUER", 1060, medievalFont, color.White)
	case StateClasse:
		drawBlack(screen)
		drawCenteredText(screen, " [G] GUERRIER   [N] MAGE   [V] VOLEUR   [A] ASSASSIN   [E] ARCHER ", 860, medievalFont, color.White)
	case StateArme:
		drawBlack(screen)
		drawCenteredText(screen, " [E] ÉPÉE   [B] BAGUETTE   [C] COUTEAU DE CUISTÔT   [D] DAGUE   [A] ARC ", 860, medievalFont, color.White)
	case StateJeu:
		drawBlack(screen)
		drawCenteredText(screen, "=== UNDERWORLD - MENU PRINCIPAL ===", 80, medievalFont, color.RGBA{255, 215, 0, 255})
		drawCenteredText(screen, fmt.Sprintf("%s | Niveau %.0f | Or: %.0f | %s", classe, level, argent, Nom), 140, medievalFont, color.White)
		drawCenteredText(screen, "------------------------------------------------------------", 190, medievalFont, color.White)

		// Section Jeu
		drawCenteredText(screen, "EXPLORATION & COMBAT", 260, medievalFont, color.RGBA{255, 200, 100, 255})
		drawCenteredText(screen, "[L] Lieux  [B] Boutique  [E] Équipement", 310, medievalFont, color.White)

		// Section Progression
		drawCenteredText(screen, "PROGRESSION", 400, medievalFont, color.RGBA{100, 200, 255, 255})
		drawCenteredText(screen, "[S] Stats  [Q] Quêtes  [T] Talents  [A] Succès", 450, medievalFont, color.White)

		// Section Inventaire
		drawCenteredText(screen, "INVENTAIRE", 540, medievalFont, color.RGBA{150, 255, 150, 255})
		drawCenteredText(screen, "[I] Ouvrir l'inventaire", 590, medievalFont, color.White)

		// Section Système
		drawCenteredText(screen, "SYSTÈME", 680, medievalFont, color.RGBA{200, 200, 200, 255})
		drawCenteredText(screen, "[P] Pause  [C] Crédits  [F5] Sauvegarder  [F9] Charger", 730, medievalFont, color.White)

		drawCenteredText(screen, "------------------------------------------------------------", 810, medievalFont, color.White)
		drawCenteredText(screen, fmt.Sprintf("Lieu actuel: %s", lieuactu), 870, medievalFont, color.RGBA{255, 255, 100, 255})
	case StateLieu:
		drawBlack(screen)
		drawCenteredText(screen, "CHOISIS UN LIEU", 200, medievalFont, color.White)
		drawCenteredText(screen, "------------------------------------------------------------", 250, medievalFont, color.White)
		drawCenteredText(screen, "[C] CHAMPS DE BLÉ  [F] FORGE  [T] TOUR DE SORCIER", 360, medievalFont, color.White)
		drawCenteredText(screen, "[S] SORTIR DU VILLAGE  [H] HÔTEL DE DUDEBU", 460, medievalFont, color.White)
		drawCenteredText(screen, "------------------------------------------------------------", 510, medievalFont, color.White)
		drawCenteredText(screen, "[ECHAP] RETOUR", 900, medievalFont, color.White)
	case StateStat:
		drawBlack(screen)
		drawCenteredText(screen, "=== STATISTIQUES DU PERSONNAGE ===", 100, medievalFont, color.RGBA{255, 215, 0, 255})
		drawCenteredText(screen, fmt.Sprintf("%s | Niveau %.0f | XP: %.0f/%.0f", classe, level, xp, xpMax), 160, medievalFont, color.White)
		drawCenteredText(screen, "------------------------------------------------------------", 210, medievalFont, color.White)

		// Barres de vie et mana
		drawCenteredText(screen, fmt.Sprintf("PV: %.0f/%.0f", pv, pvMAX), 270, medievalFont, color.RGBA{255, 100, 100, 255})
		drawBar(screen, WinW/2-300, 310, 600, 30, pv, pvMAX, color.RGBA{200, 30, 30, 255})
		drawCenteredText(screen, fmt.Sprintf("Mana: %.0f | Défense: %.0f", mana, defenseTotal), 370, medievalFont, color.RGBA{100, 150, 255, 255})
		drawBar(screen, WinW/2-300, 410, 600, 30, mana, 100, color.RGBA{30, 30, 200, 255})
		y := 460
		y = afficherStat(screen, "Force", force, 0, y, color.RGBA{220, 80, 80, 255})
		y = afficherStat(screen, "Agilité", agilite, 0, y, color.RGBA{80, 220, 120, 255})
		y = afficherStat(screen, "Intelligence", intelligence, 0, y, color.RGBA{80, 160, 240, 255})
		y = afficherStat(screen, "Endurance", endurance, 0, y, color.RGBA{220, 180, 80, 255})
		y = afficherStat(screen, "Dextérité", dexterite, 0, y, color.RGBA{180, 120, 220, 255})
		y = afficherStat(screen, "Dégâts arme", degats, 0, y, color.RGBA{240, 240, 240, 255})

		// Afficher l'équipement actuel
		y += 40
		drawCenteredText(screen, "=== ÉQUIPEMENT ===", y, medievalFont, color.RGBA{255, 215, 0, 255})
		y += 40
		if equipementArme != nil {
			drawCenteredText(screen, fmt.Sprintf("Arme: %s (+%.0f dégâts)", equipementArme.Nom, equipementArme.Degat), y, medievalFont, color.White)
			y += 30
		}
		if equipementArmure != nil {
			drawCenteredText(screen, fmt.Sprintf("Armure: %s (+%.0f défense)", equipementArmure.Nom, equipementArmure.Defense), y, medievalFont, color.White)
			y += 30
		}
		if equipementAccessoire != nil {
			drawCenteredText(screen, fmt.Sprintf("Accessoire: %s", equipementAccessoire.Nom), y, medievalFont, color.White)
			y += 30
		}

		drawCenteredText(screen, "[ECHAP] Retour", 1000, medievalFont, color.White)

	case StateBoutique:
		drawFull(screen, shopImage)
		drawCenteredText(screen,
			fmt.Sprintf("BOUTIQUE (PAGE %d) ------------- ARGENT : %.2f", pageBoutique, argent),
			160, medievalFont, color.White)
		drawCenteredText(screen, "------------------------------------------------------------", 260, medievalFont, color.White)

		if pageBoutique == 1 {

			var armesClasse []Arme
			switch classe {
			case "MAGE":
				armesClasse = armesMage
			case "GUERRIER":
				armesClasse = armesGuerrier
			case "VOLEUR":
				armesClasse = armesVoleur
			case "ASSASSIN":
				armesClasse = armesAssassin
			case "ARCHER":
				armesClasse = armesARCHER
			default:
				armesClasse = armesDispo
			}
			armesClasse = append(armesClasse,
				Arme{"Potion de Soin", 50, 0},
				Arme{"Potion de Soin Majeure", 120, 0},
				Arme{"Potion de Mana", 40, 0},
				Arme{"Potion de Poison", 30, 0},
				Arme{"Potion d'XP", 100, 0},
			)

			labels := []string{"C", "F", "O", "P", "S", "E"}
			for i := 0; i < len(armesClasse) && i < len(labels); i++ {
				drawCenteredText(screen,
					fmt.Sprintf("[%s] %s - %.0f OKSA", labels[i], armesClasse[i].Nom, armesClasse[i].Prix),
					360+i*80, medievalFont, color.White)
			}
			drawCenteredText(screen, "[→] PAGE 2 : SHOP EQUIPEMENTS", 940, medievalFont, color.White)

		} else if pageBoutique == 2 {
			// équipements
			labels := []string{"C", "F", "O", "P", "S"}
			for i := 0; i < len(shopEquipements) && i < len(labels); i++ {
				drawCenteredText(screen,
					fmt.Sprintf("[%s] %s - %d OKSA", labels[i], shopEquipements[i].Nom, shopEquipements[i].Prix),
					360+i*80, medievalFont, color.White)
			}
			drawCenteredText(screen, "[←] PAGE 1 : ARMES/POTIONS", 940, medievalFont, color.White)

		}
	case StateChamps:
		drawBlack(screen)
		drawCenteredText(screen, "VOUS ÊTES AUX CHAMPS DE BLÉ ", 200, medievalFont, color.White)
		drawCenteredText(screen, "------------------------------------------------------------", 280, medievalFont, color.White)
		drawCenteredText(screen, fmt.Sprintf("[L] LABOURER (FORCE+) : %.2f", force), 400, medievalFont, color.White)
		drawCenteredText(screen, fmt.Sprintf("[C] COURIR (AGILITÉ+) : %.2f", agilite), 500, medievalFont, color.White)
		drawCenteredText(screen, fmt.Sprintf("[V] CULTIVER/VENDRE (ARGENT+) : %.2f", argent), 600, medievalFont, color.White)
		drawCenteredText(screen, fmt.Sprintf("[N] NAGER (ENDURANCE+) : %.2f", endurance), 700, medievalFont, color.White)
		drawCenteredText(screen, "------------------------------------------------------------", 880, medievalFont, color.White)
		drawCenteredText(screen, "[ECHAP] RETOUR", 960, medievalFont, color.White)
	case StateTour:
		drawBlack(screen)
		drawCenteredText(screen, "VOUS ÊTES À LA TOUR DE SORCIER ", 200, medievalFont, color.White)
		drawCenteredText(screen, "------------------------------------------------------------", 280, medievalFont, color.White)
		drawCenteredText(screen, fmt.Sprintf("[L] LIRE (INT+) : %.2f", intelligence), 420, medievalFont, color.White)
		drawCenteredText(screen, fmt.Sprintf("[C] STRATÉGIE (DEX+) : %.2f", dexterite), 520, medievalFont, color.White)
		drawCenteredText(screen, fmt.Sprintf("[M] MANA (MANA+) : %.2f", mana), 620, medievalFont, color.White)
		drawCenteredText(screen, "[D] CONSEILS DU GRAND MAGE", 720, medievalFont, color.White)
		drawCenteredText(screen, "------------------------------------------------------------", 880, medievalFont, color.White)
		drawCenteredText(screen, "[ECHAP] RETOUR", 960, medievalFont, color.White)
	case StateForge:
		drawBlack(screen)
		drawCenteredText(screen, "VOUS ÊTES À LA FORGE ", 200, medievalFont, color.White)
		drawCenteredText(screen, "------------------------------------------------------------", 280, medievalFont, color.White)
		drawCenteredText(screen, fmt.Sprintf("[L] COMPOSITION (INT+) : %.2f", intelligence), 420, medievalFont, color.White)
		drawCenteredText(screen, fmt.Sprintf("[C] FORGE (FORCE+) : %.2f", force), 520, medievalFont, color.White)
		drawCenteredText(screen, fmt.Sprintf("[M] MODELAGE (DEX+) : %.2f", dexterite), 620, medievalFont, color.White)
		drawCenteredText(screen, "[D] CONSEILS ", 720, medievalFont, color.White)
		drawCenteredText(screen, "------------------------------------------------------------", 880, medievalFont, color.White)
		drawCenteredText(screen, "[ECHAP] RETOUR", 960, medievalFont, color.White)
	case StateInventaire:
		drawBlack(screen)
		if len(arme) == 0 {
			drawCenteredText(screen, "(Vide)", 400, medievalFont, color.Black)
		} else {
			startY := 340
			for i, it := range arme {
				drawCenteredText(screen, fmt.Sprintf("[%d] %s", i, it), startY+i*70, medievalFont, color.Black)
			}
		}
		drawCenteredText(screen, "[ECHAP] RETOUR", 1000, medievalFont, color.White)
		drawCenteredText(screen, "UTILISER ITEM: [0-9] (En combat, les potions servent d'action)", 920, medievalFont, color.White)

	case StateHotel:
		drawBlack(screen)
		drawCenteredText(screen, "VOUS ÊTES À L'HÔTEL DE VILLE ", 200, medievalFont, color.White)
		drawCenteredText(screen, "------------------------------------------------------------", 280, medievalFont, color.White)
		drawCenteredText(screen, fmt.Sprintf("[L] REGISTRES (INT+)%.2f", intelligence), 420, medievalFont, color.White)
		drawCenteredText(screen, "[C] PARLER", 520, medievalFont, color.White)
		drawCenteredText(screen, fmt.Sprintf("[N] IMPÔTS (-ARGENT)%.2f", argent), 620, medievalFont, color.White)
		drawCenteredText(screen, fmt.Sprintf("[D] VOLER (+ARGENT)%.2f", argent), 720, medievalFont, color.White)
		drawCenteredText(screen, "------------------------------------------------------------", 880, medievalFont, color.White)
		drawCenteredText(screen, "[ECHAP] RETOUR", 960, medievalFont, color.White)
	case StateSortie:
		drawBlack(screen)
		drawCenteredText(screen, "VOUS ÊTES À LA GRANDE PORTE ", 200, medievalFont, color.White)
		drawCenteredText(screen, "------------------------------------------------------------", 280, medievalFont, color.White)
		drawCenteredText(screen, "[C] SORTIR COMBATTRE DES MONSTRES ", 480, medievalFont, color.White)
		drawCenteredText(screen, "[R] RENTRER DANS LE VILLAGE", 620, medievalFont, color.White)
		drawCenteredText(screen, "------------------------------------------------------------", 880, medievalFont, color.White)
		drawCenteredText(screen, "[ECHAP] RETOUR", 960, medievalFont, color.White)
	case StateAttention:
		drawFull(screen, shopImage)
		drawCenteredText(screen, "------------------------- ATTENTION -------------------------", 200, medievalFont, color.White)
		drawCenteredText(screen, "------------------------------------------------------------", 280, medievalFont, color.White)
		drawCenteredText(screen, "VOUS NE POUVEZ PAS AVOIR PLUS DE 10 ITEMS", 440, medievalFont, color.White)
		drawCenteredText(screen, "DANS VOTRE INVENTAIRE", 520, medievalFont, color.White)
		drawCenteredText(screen, "------------------------------------------------------------", 680, medievalFont, color.White)
		drawCenteredText(screen, "[B] RETOUR BOUTIQUE", 1000, medievalFont, color.White)
	case StateAttentionPrix:
		drawFull(screen, shopImage)
		drawCenteredText(screen, "------------------------- ATTENTION -------------------------", 160, medievalFont, color.White)
		drawCenteredText(screen, "VOUS N'AVEZ PAS ASSEZ D'OKSA", 540, medievalFont, color.White)
		drawCenteredText(screen, "-------------------------------------------------------------", 910, medievalFont, color.White)
		drawCenteredText(screen, "[B] RETOUR BOUTIQUE", 1000, medievalFont, color.White)
	case StateVendre:
		drawFull(screen, shopImage)
		drawCenteredText(screen, "------------------------- VENTES -------------------------", 60, medievalFont, color.White)
		drawCenteredText(screen, "APPUYEZ SUR [1-9] POUR VENDRE L'ITEM À CETTE POSITION", 120, medievalFont, color.White)
		for i, it := range arme {
			drawCenteredText(screen, fmt.Sprintf("[%d] %s", i+1, it), 200+i*40, medievalFont, color.White)
		}
		drawCenteredText(screen, "-------------------------------------------------------------", 910, medievalFont, color.White)
		drawCenteredText(screen, "[B] RETOUR BOUTIQUE", 1000, medievalFont, color.White)
	case StateConfirmation:
		drawFull(screen, shopImage)
		drawCenteredText(screen, "=== CONFIRMATION D'ACHAT ===", 200, medievalFont, color.RGBA{255, 215, 0, 255})
		drawCenteredText(screen, "------------------------------------------------------------", 260, medievalFont, color.White)

		drawCenteredText(screen, fmt.Sprintf("Objet: %s", armeSelectionnee.Nom), 380, medievalFont, color.White)
		drawCenteredText(screen, fmt.Sprintf("Prix: %.0f or", armeSelectionnee.Prix), 450, medievalFont, color.RGBA{255, 255, 100, 255})

		drawCenteredText(screen, "------------------------------------------------------------", 540, medievalFont, color.White)
		drawCenteredText(screen, "[V] VALIDER L'ACHAT", 650, medievalFont, color.RGBA{100, 255, 100, 255})
		drawCenteredText(screen, "[R] ANNULER", 720, medievalFont, color.RGBA{255, 100, 100, 255})

	case StateProposerEquiper:
		drawFull(screen, shopImage)
		drawCenteredText(screen, "=== OBJET ACHETE ===", 200, medievalFont, color.RGBA{100, 255, 100, 255})
		drawCenteredText(screen, "------------------------------------------------------------", 260, medievalFont, color.White)

		drawCenteredText(screen, fmt.Sprintf("%s", armeSelectionnee.Nom), 380, medievalFont, color.White)
		drawCenteredText(screen, "a été ajouté à votre inventaire !", 450, medievalFont, color.White)

		drawCenteredText(screen, "------------------------------------------------------------", 540, medievalFont, color.White)
		drawCenteredText(screen, "Voulez-vous l'équiper maintenant ?", 630, medievalFont, color.RGBA{255, 255, 100, 255})

		drawCenteredText(screen, "[E] ÉQUIPER MAINTENANT", 740, medievalFont, color.RGBA{100, 255, 100, 255})
		drawCenteredText(screen, "[N] GARDER DANS L'INVENTAIRE", 810, medievalFont, color.RGBA{200, 200, 200, 255})
	case StateBackLife:
		drawBlack(screen)
		drawCenteredText(screen, "------------------------- VOUS ÊTES MORT -------------------------", 80, medievalFont, color.White)
		drawCenteredText(screen, "APPUYEZ SUR [ECHAP] POUR REVENIR À LA VIE A 1/2 DE VOS PV MAX", 720, medievalFont, color.White)
	case StatePause:
		drawBlack(screen)
		drawCenteredText(screen, "PAUSE", 160, medievalFont, color.White)
		drawCenteredText(screen, "[J] REPRENDRE", 1000, medievalFont, color.White)
	case StateForet:
		drawFull(screen, foretImage)
		drawCenteredText(screen, " COMBATTRE [C]", 480, medievalFont, color.White)
		drawCenteredText(screen, "FUIR [F]", 620, medievalFont, color.White)
	case StateVictoire:
		drawBlack(screen)
		drawCenteredText(screen, "[ESPACE] CONTINUER  ", 1000, medievalFont, color.White)

	case StateResultat:
		drawBlack(screen)
		drawCenteredText(screen, "BUTIN OBTENU :", 80, medievalFont, color.Black)

		if dernierLoot != "" {
			drawCenteredText(screen, fmt.Sprintf(" - %s", dernierLoot), 300, medievalFont, color.Black)
		}

		drawCenteredText(screen, "[ESPACE] CONTINUER  ", 1000, medievalFont, color.Black)
	case StateDefaite:
		drawBlack(screen)
	case StateCombat:
		if queteActuelle == 1 {
			drawFull(screen, fondqueteImage)
		}
		if queteActuelle == 2 {
			drawFull(screen, fondquete2Image)
		}
		if queteActuelle == 3 {
			drawFull(screen, fondquete3Image)
		}
		switch monster {
		case "Crabauge":
			drawCombatScene(g, screen, color.RGBA{120, 180, 255, 255})
		case "Vorlapin":
			drawCombatScene(g, screen, color.RGBA{200, 240, 200, 255})
		case "Gobelin":
			drawCombatScene(g, screen, color.RGBA{80, 200, 80, 255})
		case "Boss Lycaon":
			drawCombatScene(g, screen, color.RGBA{180, 60, 60, 255})
		case "Muddig":
			drawCombatScene(g, screen, color.RGBA{180, 60, 60, 255})
		case "Gros Serpent Vorace":
			drawCombatScene(g, screen, color.RGBA{180, 60, 60, 255})
		case "Serpent Livestide":
			drawCombatScene(g, screen, color.RGBA{180, 60, 60, 255})
		case "Loosers Wood":
			drawCombatScene(g, screen, color.RGBA{180, 60, 60, 255})
		case "Boss Wezaemon":
			drawCombatScene(g, screen, color.RGBA{180, 60, 60, 255})
		case "Poiscaille Zombie":
			drawCombatScene(g, screen, color.RGBA{180, 60, 60, 255})
		case "Lugia":
			drawCombatScene(g, screen, color.RGBA{180, 60, 60, 255})
		case "Leviathan":
			drawCombatScene(g, screen, color.RGBA{180, 60, 60, 255})
		case "Atlanticus Repunorca - Orque électrique":
			drawCombatScene(g, screen, color.RGBA{180, 60, 60, 255})
		case "Kthaanid - Maître des Abysses":
			drawCombatScene(g, screen, color.RGBA{180, 60, 60, 255})

		}
		// === INTERFACE DE COMBAT AMÉLIORÉE ===

		// Barre de vie du joueur (en haut à droite)
		drawBar(screen, WinW-320, 30, 280, 25, pv, pvMAX, color.RGBA{200, 30, 30, 255})
		text.Draw(screen, fmt.Sprintf("Vous: %.0f/%.0f PV", pv, pvMAX), medievalFont, WinW-310, 25, color.White)

		// Barre de mana du joueur
		drawBar(screen, WinW-320, 70, 280, 20, mana, 100, color.RGBA{30, 100, 255, 255})
		text.Draw(screen, fmt.Sprintf("Mana: %.0f", mana), medievalFont, WinW-310, 65, color.RGBA{150, 200, 255, 255})

		// Stats du joueur (coin droit)
		text.Draw(screen, fmt.Sprintf("Defense: %.0f", defenseTotal), medievalFont, WinW-310, 110, color.RGBA{200, 200, 100, 255})
		text.Draw(screen, fmt.Sprintf("Niveau: %.0f", level), medievalFont, WinW-310, 140, color.White)

		// Barre de vie de l'ennemi (centre-haut, grande et visible)
		barWidth := 500.0
		barX := WinW/2 - barWidth/2
		drawCenteredText(screen, fmt.Sprintf("=== %s ===", ennemiNom), 160, medievalFont, ennemiColor)
		drawBar(screen, int(barX), 190, int(barWidth), 35, ennemiPV, ennemiPVMax, color.RGBA{255, 50, 50, 255})
		drawCenteredText(screen, fmt.Sprintf("%.0f / %.0f PV", ennemiPV, ennemiPVMax), 215, medievalFont, color.White)

		// Contrôles (en bas, centrés et visibles)
		drawCenteredText(screen, "=== ACTIONS DE COMBAT ===", 920, medievalFont, color.RGBA{255, 215, 0, 255})
		drawCenteredText(screen, "[A] Attaque Legere  [Z] Attaque Forte  [E] Competence", 970, medievalFont, color.White)
		drawCenteredText(screen, "[ENTER] Pause  [I] Inventaire", 1010, medievalFont, color.RGBA{150, 150, 150, 255})

		// Log de combat (bas gauche, mieux formaté)
		text.Draw(screen, "=== HISTORIQUE ===", medievalFont, 40, WinH-160, color.RGBA{255, 215, 0, 255})
		y := WinH - 120
		for i := max(0, len(logCombat)-3); i < len(logCombat); i++ {
			text.Draw(screen, "  "+logCombat[i], medievalFont, 40, y, color.RGBA{200, 200, 255, 255})
			y += 35
		}

	case StateHistoire1:
		drawBlack(screen)
		text.Draw(screen, "On raconte qu’il y a cent ans, notre cité fut bâtie sur les ruines d’un temple oublié.", medievalFont, WinW-1200, 60, color.Black)
		text.Draw(screen, "Les pierres de ses fondations portent encore les symboles des anciens dieux,", medievalFont, WinW-1200, 120, color.White)
		text.Draw(screen, "et certains disent qu’elles veillent toujours sur nous", medievalFont, WinW-1200, 180, color.Black)
		text.Draw(screen, "C’est peut-être pour cela que, malgré les guerres et les famines, notre ville n’a jamais plié.", medievalFont, WinW-1200, 240, color.Black)
		drawCenteredText(screen, "[ESPACE] CONTINUER  ", 1000, medievalFont, color.White)

	case StateHistoire2:
		drawBlack(screen)
		text.Draw(screen, "Mon Petit, nul ne protège mieux notre cité que moi, Le Grand Mage Hugo,", medievalFont, WinW-1000, 60, color.White)
		text.Draw(screen, "gardien des arcanes et dernier héritier du savoir des anciens prêtres.", medievalFont, WinW-1000, 120, color.White)
		text.Draw(screen, "On dit qu’il converse encore avec les esprits du temple,", medievalFont, WinW-1000, 180, color.White)
		text.Draw(screen, "et que son bâton renferme la lumière même des dieux oubliés.", medievalFont, WinW-1000, 240, color.White)
		drawCenteredText(screen, "[ESPACE] CONTINUER  ", 1000, medievalFont, color.White)

	case StateHistoire3:
		drawBlack(screen)
		text.Draw(screen, "Parmi ceux qui protègent notre cité, nul n’est plus respecté que le Maître Forgeron.", medievalFont, WinW-1800, 700, color.White)
		text.Draw(screen, "Ses mains ont façonné les lames qui ont repoussé les hordes venues des ténèbres,", medievalFont, WinW-1800, 750, color.White)
		text.Draw(screen, "et chaque coup de son marteau résonne comme un serment de protection.", medievalFont, WinW-1800, 800, color.White)
		text.Draw(screen, "On dit que dans les flammes de sa forge brûle un feu sacré, hérité des anciens dieux.", medievalFont, WinW-1800, 850, color.White)
		drawCenteredText(screen, "[ESPACE] CONTINUER  ", 1000, medievalFont, color.White)
	case StateThemeG:
		drawFull(screen, gcombatImage)
		drawCenteredText(screen, "[ESPACE] CONTINUER  ", 1000, medievalFont, color.RGBA{R: 255, G: 215, B: 0, A: 255})
	case StateThemeA:
		drawFull(screen, acombatImage)
		drawCenteredText(screen, "[ESPACE] CONTINUER  ", 1000, medievalFont, color.RGBA{R: 255, G: 215, B: 0, A: 255})
	case StateThemeM:
		drawFull(screen, mcombatImage)
		drawCenteredText(screen, "[ESPACE] CONTINUER  ", 1000, medievalFont, color.RGBA{R: 255, G: 215, B: 0, A: 255})
	case StateThemeV:
		drawFull(screen, vcombatImage)
		drawCenteredText(screen, "[ESPACE] CONTINUER  ", 1000, medievalFont, color.RGBA{R: 255, G: 215, B: 0, A: 255})
	case StateThemeE:
		drawFull(screen, ecombatImage)
		drawCenteredText(screen, "[ESPACE] CONTINUER  ", 1000, medievalFont, color.RGBA{R: 255, G: 215, B: 0, A: 255})

	case StateCredit:
		drawBlack(screen)
		text.Draw(screen, "EVAN FEAT MIEL POPS", medievalFont, 200, 550, color.White)
		text.Draw(screen, "GEOGEO EST SON VRAI BLAZE ", medievalFont, 200, 650, color.Black)
		text.Draw(screen, "HUGO LE SINGE SAVANT ", medievalFont, 200, 750, color.Black)
		drawCenteredText(screen, "[ESPACE] RETOUR  ", 1000, medievalFont, color.Black)

	case StateEquipement:
		drawBlack(screen)
		text.Draw(screen, fmt.Sprintf("Casque : %v", slotToString(equipSlots.Casque)), medievalFont, 300, 280, color.White)
		text.Draw(screen, "[1] Équiper / Déséquiper casque", medievalFont, 300, 340, color.White)
		text.Draw(screen, fmt.Sprintf("Plastron : %v", slotToString(equipSlots.Armure)), medievalFont, 300, 420, color.White)
		text.Draw(screen, "[2] Équiper / Déséquiper plastron", medievalFont, 300, 480, color.White)
		text.Draw(screen, fmt.Sprintf("Bottes : %v", slotToString(equipSlots.Bottes)), medievalFont, 300, 560, color.White)
		text.Draw(screen, "[3] Équiper / Déséquiper bottes", medievalFont, 300, 620, color.White)
		text.Draw(screen, fmt.Sprintf("Arme : %v", slotToString(equipSlots.Arme)), medievalFont, 300, 700, color.White)
		text.Draw(screen, "[4] Équiper / Déséquiper arme", medievalFont, 300, 760, color.White)
		text.Draw(screen, fmt.Sprintf("Anneau : %v", slotToString(equipSlots.Anneau)), medievalFont, 300, 840, color.White)
		text.Draw(screen, "[5] Équiper / Déséquiper anneau", medievalFont, 300, 900, color.White)

	case StateQuetes:
		drawFull(screen, fondqueteImage)
		drawCenteredText(screen, "=== JOURNAL DE QUETES ===", 80, medievalFont, color.RGBA{255, 215, 0, 255})
		drawCenteredText(screen, "Suivez vos objectifs et reclamez vos recompenses", 130, medievalFont, color.RGBA{200, 200, 200, 255})

		y := 200 - scrollOffset
		contentHeight := 0

		for i, q := range quetes {
			if y < 150 || y > 950 {
				// Ne pas dessiner si hors écran
				y += 195
				contentHeight += 195
				continue
			}

			var couleur color.Color = color.White
			statut := "EN COURS"
			if q.Complete {
				statut = "TERMINEE"
				couleur = color.RGBA{100, 255, 100, 255}
			} else if q.Active {
				statut = "ACTIVE"
				couleur = color.RGBA{255, 215, 0, 255}
			}

			drawCenteredText(screen, fmt.Sprintf("%d. %s [%s]", i+1, q.Nom, statut), y, medievalFont, couleur)
			y += 40
			drawCenteredText(screen, q.Description, y, medievalFont, color.RGBA{180, 180, 180, 255})
			y += 35

			// Barre de progression
			drawBar(screen, WinW/2-200, y, 400, 20, float64(q.Progres), float64(q.ProgresMax), color.RGBA{100, 150, 255, 255})
			drawCenteredText(screen, fmt.Sprintf("%d/%d", q.Progres, q.ProgresMax), y+15, medievalFont, color.White)
			y += 40

			if !q.Complete {
				drawCenteredText(screen, fmt.Sprintf("Recompense: %s + %.0f or + %d XP", q.Recompense, q.RecompenseOr, q.RecompenseXP), y, medievalFont, color.RGBA{255, 255, 100, 255})
			} else {
				drawCenteredText(screen, "RECOMPENSE RECLAMEE!", y, medievalFont, color.RGBA{100, 255, 100, 255})
			}
			y += 60
			contentHeight += 195
		}

		// Calculer scrollMax
		scrollMax = contentHeight - 600
		if scrollMax < 0 {
			scrollMax = 0
		}

		// Indicateur de scroll
		if scrollMax > 0 {
			drawCenteredText(screen, "[PageUp/PageDown] Defiler", 950, medievalFont, color.RGBA{150, 150, 150, 255})
		}

		drawCenteredText(screen, "[ECHAP] Retour au menu", 1000, medievalFont, color.White)

	case StateTalents:
		drawFull(screen, fondquete2Image)
		drawCenteredText(screen, "=== ARBRE DE TALENTS ===", 80, medievalFont, color.RGBA{255, 215, 0, 255})
		drawCenteredText(screen, fmt.Sprintf("Points disponibles: %d | Classe: %s", pointsTalents, classe), 130, medievalFont, color.RGBA{100, 255, 100, 255})
		drawCenteredText(screen, "Ameliorez vos competences pour devenir plus puissant!", 170, medievalFont, color.RGBA{200, 200, 200, 255})

		y := 240
		if classTalents, ok := talents[classe]; ok {
			for i, t := range classTalents {
				var couleur color.Color = color.White
				if t.Niveau == t.NiveauMax {
					couleur = color.RGBA{255, 215, 0, 255}
				} else if t.Niveau > 0 {
					couleur = color.RGBA{150, 200, 255, 255}
				}

				drawCenteredText(screen, fmt.Sprintf("[%d] %s [%d/%d]", i+1, t.Nom, t.Niveau, t.NiveauMax), y, medievalFont, couleur)
				y += 40
				drawCenteredText(screen, t.Description, y, medievalFont, color.RGBA{180, 180, 180, 255})
				y += 35

				// Barre de niveau
				drawBar(screen, WinW/2-150, y, 300, 15, float64(t.Niveau), float64(t.NiveauMax), color.RGBA{255, 215, 0, 255})
				y += 50
			}
		}

		drawCenteredText(screen, "Appuyez sur [1-9] pour ameliorer un talent", 950, medievalFont, color.RGBA{200, 200, 200, 255})
		drawCenteredText(screen, "[ECHAP] Retour au menu", 1000, medievalFont, color.White)

	case StateAchievements:
		drawFull(screen, fondquete3Image)
		drawCenteredText(screen, "=== SUCCES & ACHIEVEMENTS ===", 80, medievalFont, color.RGBA{255, 215, 0, 255})
		drawCenteredText(screen, "Debloquez des succes pour obtenir des recompenses permanentes!", 130, medievalFont, color.RGBA{200, 200, 200, 255})

		y := 200 - scrollOffset
		contentHeight := 0

		for _, a := range achievements {
			if y < 150 || y > 900 {
				// Ne pas dessiner si hors écran
				y += 150
				contentHeight += 150
				continue
			}

			var couleur color.Color = color.RGBA{128, 128, 128, 255}
			if a.Deverrouille {
				couleur = color.RGBA{255, 215, 0, 255}
			}

			drawCenteredText(screen, fmt.Sprintf("%s %s", a.Icone, a.Nom), y, medievalFont, couleur)
			y += 40
			drawCenteredText(screen, a.Description, y, medievalFont, color.RGBA{180, 180, 180, 255})
			y += 35

			if !a.Deverrouille {
				progres := fmt.Sprintf("    Progression: %d/%d", a.Progres, a.Objectif)
				text.Draw(screen, progres, medievalFont, 320, y, color.White)
				y += 30
				text.Draw(screen, fmt.Sprintf("    Recompense: %s", a.Recompense), medievalFont, 320, y, color.White)
			} else {
				text.Draw(screen, "    Termine!", medievalFont, 320, y, color.White)
			}
			y += 45
			contentHeight += 150
		}

		// Calculer scrollMax
		scrollMax = contentHeight - 550
		if scrollMax < 0 {
			scrollMax = 0
		}

		// Indicateur de scroll
		if scrollMax > 0 {
			text.Draw(screen, "[PageUp/PageDown] Defiler", medievalFont, 320, 900, color.RGBA{150, 150, 150, 255})
		}

		text.Draw(screen, "[ECHAP] Retour", medievalFont, 300, 950, color.White)

	}

	// Dessiner la popup par-dessus tout
	DrawPopup(screen)
}

func (g *Game) TakeDamage() {
	if g.DamageTimer == 0 {
		pv -= 5
		g.DamageTimer = 30
	}
}

func (g *Game) MonsterTakeDamage() {
	if g.MonstreDamageTimer == 0 {
		ennemiPV -= 5
		g.MonstreDamageTimer = 30
	}
}

func chargerSprites(g *Game) {
	var playerPath string
	switch classe {
	case "GUERRIER":
		playerPath = "assets/perso/chevalier.png"
	case "MAGE":
		playerPath = "assets/perso/mage.png"
	case "ASSASSIN":
		playerPath = "assets/perso/assassin.png"
	case "VOLEUR":
		playerPath = "assets/perso/voleur.png"
	case "ARCHER":
		playerPath = "assets/perso/elfe.png"
	default:
		playerPath = "assets/perso/chevalier.png"
	}
	playerImage, _, err := ebitenutil.NewImageFromFile(playerPath)
	if err != nil {
		log.Fatal(err)
	}
	g.JoueurImage = playerImage

	var monsterPath string
	switch monster {
	case "Crabauge":
		monsterPath = "assets/monstres/crabauge.png"
	case "Vorlapin":
		monsterPath = "assets/monstres/lapin.png"
	case "Gobelin":
		monsterPath = "assets/monstres/gobelin.png"
	case "Boss Lycaon":
		monsterPath = "assets/monstres/bossloup.png"
	case "Muddig":
		monsterPath = "assets/monstres/muddig.png"
	case "Gros Serpent Vorace":
		monsterPath = "assets/monstres/serpent.png"
	case "Serpent Livestide":
		monsterPath = "assets/monstres/serpentaqua.png"
	case "Loosers Wood":
		monsterPath = "assets/monstres/looserswood.png"
	case "Boss Wezaemon":
		monsterPath = "assets/monstres/wezaemon.png"
	case "Poiscaille Zombie":
		monsterPath = "assets/monstres/poissonzombie.png"
	case "Lugia":
		monsterPath = "assets/monstres/lugia.png"
	case "Leviathan":
		monsterPath = "assets/monstres/leviathan.png"
	case "Atlanticus Repunorca - Orque électrique":
		monsterPath = "assets/monstres/orque.png"
	case "Kthaanid - Maître des Abysses":
		monsterPath = "assets/monstres/pieuvre.png"

	default:
		monsterPath = "assets/monstres/crabauge.png"
	}
	monsterImage, _, err := ebitenutil.NewImageFromFile(monsterPath)
	if err != nil {
		log.Fatal(err)
	}
	g.MonstreImage = monsterImage
}

// ---------------------------- UTIL DRAW ----------------------------

func drawCenteredText(screen *ebiten.Image, str string, y int, face font.Face, clr color.Color) {
	lines := strings.Split(strings.ToUpper(str), "\n")
	metrics := face.Metrics()
	lineHeight := metrics.Height.Round()
	sw, _ := screen.Size()
	for i, line := range lines {
		bounds := text.BoundString(face, line)
		x := (sw - bounds.Dx()) / 2
		text.Draw(screen, line, face, x, y+i*lineHeight, clr)
	}
}

func drawBar(screen *ebiten.Image, x, y, width, height int, value, max float64, clr color.RGBA) {
	if max <= 0 {
		max = 1
	}
	ratio := value / max
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(width), float64(height), color.RGBA{30, 30, 30, 255})
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(width)*ratio, float64(height), clr)
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(width), 2, color.RGBA{255, 255, 255, 30})
	ebitenutil.DrawRect(screen, float64(x), float64(y+height-2), float64(width), 2, color.RGBA{255, 255, 255, 30})
}

func afficherStat(screen *ebiten.Image, nom string, valeur float64, bonus float64, y int, statColor color.RGBA) int {
	signe := ""
	if bonus > 0 {
		signe = "+"
	}
	texte := fmt.Sprintf("%-15s : %5.2f", nom, valeur)
	drawCenteredText(screen, texte, y, medievalFont, statColor)
	if bonus != 0 {
		drawCenteredText(screen, fmt.Sprintf("(Bonus: %s%.2f)", signe, bonus), y+30, medievalFont, color.RGBA{200, 200, 200, 255})
		return y + 70
	}
	return y + 50
}

func drawFull(screen *ebiten.Image, img *ebiten.Image) {
	if img == nil {
		screen.Fill(color.Black)
		return
	}
	sw, sh := screen.Size()
	iw, ih := img.Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(sw)/float64(iw), float64(sh)/float64(ih))
	screen.DrawImage(img, op)
}

func drawBlack(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 10, 20, 255}) // Fond bleu très foncé pour moins de fatigue visuelle
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return WinW, WinH
}

// ---------------------------- UTILS ALEATOIRES ----------------------------

func randFloat(min, max float64) float64 {
	if max < min {
		min, max = max, min
	}
	return min + rand.Float64()*(max-min)
}

func miss(pct int) bool {
	return rand.Intn(100) < pct
}

func crit(pct int) bool {
	return rand.Intn(100) < pct
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func _maxi(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ---------------------------- INIT IMAGES ----------------------------

func mustLoad(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Printf("⚠️ Image non trouvée : %s (une image par défaut sera utilisée)\n", path)
		// Créer une image vide de 1x1 pixel comme fallback
		img = ebiten.NewImage(1, 1)
	}
	return img
}

// Fonctions Init pour les images encore utilisées
func InitShop(path string)    { shopImage = mustLoad(path) }
func InitForet(path string)   { foretImage = mustLoad(path) }
func InitQuete(path string)   { fondqueteImage = mustLoad(path) }
func InitQuete2(path string)  { fondquete2Image = mustLoad(path) }
func InitQuete3(path string)  { fondquete3Image = mustLoad(path) }
func InitMCombat(path string) { mcombatImage = mustLoad(path) }
func InitGCombat(path string) { gcombatImage = mustLoad(path) }
func InitVCombat(path string) { vcombatImage = mustLoad(path) }
func InitACombat(path string) { acombatImage = mustLoad(path) }
func InitECombat(path string) { ecombatImage = mustLoad(path) }

// ---------------------------- MAIN ----------------------------

func main() {
	initPseudo()
	go PlaySong()

	// load font (try medieval.ttf, otherwise fall back to basicfont.Face7x13)
	medievalFont = nil
	if data, err := ioutil.ReadFile("medieval.ttf"); err == nil {
		if tt, err := opentype.Parse(data); err == nil {
			if face, err := opentype.NewFace(tt, &opentype.FaceOptions{
				Size:    32,
				DPI:     72,
				Hinting: font.HintingFull,
			}); err == nil {
				medievalFont = face
			}
		}
	}
	// fallback if font not loaded
	if medievalFont == nil {
		medievalFont = basicfont.Face7x13
	}

	// Charger les images essentielles (boutique, forêt, quêtes, combat)
	InitShop("image/shop.png")
	InitForet("image/foret.png")
	InitQuete("image/quete1.png")
	InitQuete2("image/quete2.png")
	InitQuete3("image/quete3.png")
	InitGCombat("image/gcombat.png")
	InitMCombat("image/mcombat.png")
	InitECombat("image/ecombat.png")
	InitVCombat("image/vcombat.png")
	InitACombat("image/acombat.png")

	// Initialiser les systèmes de jeu
	initQuetes()
	initTalents()
	initAchievements()

	// initial player sprite (fallback)
	playerPath := "assets/perso/chevalier.png"
	playerImage, _, err := ebitenutil.NewImageFromFile(playerPath)
	if err != nil {
		log.Printf("⚠️ Sprite du joueur non trouvé : %s (sprite par défaut utilisé)\n", playerPath)
		playerImage = ebiten.NewImage(32, 32)
	}

	// create game and set images
	g := &Game{
		X:         1200,
		Y:         620,
		MonstreX:  100,
		MonstreY:  620,
		MonstreDX: 2,
	}
	g.JoueurImage = playerImage

	// default monster sprite (will be overwritten when combat starts)
	monsterPath := "assets/monstres/crabauge.png"
	if img, _, err := ebitenutil.NewImageFromFile(monsterPath); err == nil {
		g.MonstreImage = img
	}

	// setup window
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Underworld - Un Monde Fantastique")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.IsFullscreen()

	// run game
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
