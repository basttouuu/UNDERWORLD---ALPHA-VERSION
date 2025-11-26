package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Sauvegarde contient toutes les données du joueur
type Sauvegarde struct {
	// Infos personnage
	Nom           string
	Classe        string
	Level         float64
	XP            float64
	XPMax         float64
	Argent        float64
	Titre         string
	LieuActu      string
	QueteActuelle int

	// Stats
	PV           float64
	PVMax        float64
	Mana         float64
	Force        float64
	Agilite      float64
	Endurance    float64
	Intelligence float64
	Dexterite    float64
	Degats       float64
	DefenseTotal float64

	// Inventaire
	Arme []string

	// Équipement
	EquipementArme       *Equipement
	EquipementArmure     *Equipement
	EquipementAccessoire *Equipement

	// Progression
	MonstreIndex int
}

func sauvegarderJeu(fichier string) error {
	save := Sauvegarde{
		Nom:                  Nom,
		Classe:               classe,
		Level:                level,
		XP:                   xp,
		XPMax:                xpMax,
		Argent:               argent,
		Titre:                titre,
		LieuActu:             lieuactu,
		QueteActuelle:        queteActuelle,
		PV:                   pv,
		PVMax:                pvMAX,
		Mana:                 mana,
		Force:                force,
		Agilite:              agilite,
		Endurance:            endurance,
		Intelligence:         intelligence,
		Dexterite:            dexterite,
		Degats:               degats,
		DefenseTotal:         defenseTotal,
		Arme:                 arme,
		EquipementArme:       equipementArme,
		EquipementArmure:     equipementArmure,
		EquipementAccessoire: equipementAccessoire,
		MonstreIndex:         monstreIndex,
	}

	data, err := json.MarshalIndent(save, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fichier, data, 0644)
	if err != nil {
		return err
	}

	log.Println("✅ Partie sauvegardée:", fichier)
	return nil
}

func chargerJeu(fichier string) error {
	data, err := ioutil.ReadFile(fichier)
	if err != nil {
		return err
	}

	var save Sauvegarde
	err = json.Unmarshal(data, &save)
	if err != nil {
		return err
	}

	// Restaurer les variables
	Nom = save.Nom
	classe = save.Classe
	level = save.Level
	xp = save.XP
	xpMax = save.XPMax
	argent = save.Argent
	titre = save.Titre
	lieuactu = save.LieuActu
	queteActuelle = save.QueteActuelle
	pv = save.PV
	pvMAX = save.PVMax
	mana = save.Mana
	force = save.Force
	agilite = save.Agilite
	endurance = save.Endurance
	intelligence = save.Intelligence
	dexterite = save.Dexterite
	degats = save.Degats
	defenseTotal = save.DefenseTotal
	arme = save.Arme
	equipementArme = save.EquipementArme
	equipementArmure = save.EquipementArmure
	equipementAccessoire = save.EquipementAccessoire
	monstreIndex = save.MonstreIndex

	log.Println("✅ Partie chargée:", fichier)
	return nil
}
