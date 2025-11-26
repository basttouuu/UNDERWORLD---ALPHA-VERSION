package main

import (
	"fmt"
	"strings"
)

func initPseudo() {
	fmt.Println("---------------------- Choisissez votre pseudo : -----------------------")
	fmt.Scanln(&Nom)
	if Nom == "" || Nom == "Noir" {
		fmt.Println("❌ Pseudo vide, recommencez.")
		initPseudo()
	}
	Nom = formatPseudo(Nom)
	fmt.Println("Votre pseudo sera donc :", Nom)
	fmt.Println("Parfait, continuons !")
	fmt.Println("✅ Parfait, lancement du jeu...")
}

func formatPseudo(p string) string {
	if len(p) == 0 {
		return ""
	}
	return strings.ToUpper(string(p[0])) + strings.ToLower(p[1:])
}

