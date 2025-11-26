package main

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func PlaySong() {
	f, err := os.Open("musique.mp3")
	if err != nil {
		log.Println("⚠️ Musique non trouvée (musique.mp3), le jeu continue sans son")
		return
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Println("⚠️ Erreur de lecture de la musique, le jeu continue sans son")
		return
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	ctrlStreamer := &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}
	done := make(chan bool)
	speaker.Play(beep.Seq(ctrlStreamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}
