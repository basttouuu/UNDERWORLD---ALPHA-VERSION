# Script de vérification rapide - Underworld v2.1

Write-Host ""
Write-Host "=== UNDERWORLD v2.1 - VERIFICATION ===" -ForegroundColor Cyan
Write-Host ""

# Vérifier l'exécutable
if (Test-Path "underworld.exe") {
    $exe = Get-Item "underworld.exe"
    $sizeMB = [math]::Round($exe.Length/1MB, 2)
    Write-Host "Executable:" -ForegroundColor Yellow
    Write-Host "  Taille: $sizeMB MB" -ForegroundColor Green
    Write-Host "  Date: $($exe.LastWriteTime)" -ForegroundColor Green
    Write-Host ""
} else {
    Write-Host "  ERREUR: underworld.exe manquant!" -ForegroundColor Red
}

# Compter les fichiers
$goFiles = (Get-ChildItem -Filter "*.go").Count
$mdFiles = (Get-ChildItem -Filter "*.md").Count
$totalCode = (Get-ChildItem -Filter "*.go" | Measure-Object -Property Length -Sum).Sum / 1KB

Write-Host "Fichiers sources:" -ForegroundColor Yellow
Write-Host "  Fichiers .go: $goFiles" -ForegroundColor Green
Write-Host "  Fichiers .md: $mdFiles" -ForegroundColor Green
Write-Host "  Code total: $([math]::Round($totalCode, 2)) KB" -ForegroundColor Green
Write-Host ""

# Vérifier les dossiers
$assets = Test-Path "assets"
$images = Test-Path "image"
$audio = Test-Path "musique.mp3"
$font = Test-Path "medieval.ttf"

Write-Host "Ressources:" -ForegroundColor Yellow
if ($assets) { Write-Host "  Assets: OK" -ForegroundColor Green } else { Write-Host "  Assets: MANQUANT" -ForegroundColor Red }
if ($images) { Write-Host "  Images: OK" -ForegroundColor Green } else { Write-Host "  Images: MANQUANT" -ForegroundColor Red }
if ($audio) { Write-Host "  Musique: OK" -ForegroundColor Green } else { Write-Host "  Musique: MANQUANT" -ForegroundColor Red }
if ($font) { Write-Host "  Police: OK" -ForegroundColor Green } else { Write-Host "  Police: MANQUANT" -ForegroundColor Red }
Write-Host ""

# Résumé des améliorations
Write-Host "Ameliorations v2.1:" -ForegroundColor Yellow
Write-Host "  Interface principale redesignee" -ForegroundColor Cyan
Write-Host "  Combat ameliore avec HUD clair" -ForegroundColor Cyan
Write-Host "  Menus Quetes/Talents/Succes centres" -ForegroundColor Cyan
Write-Host "  Barres de progression visuelles" -ForegroundColor Cyan
Write-Host "  Code nettoye et optimise" -ForegroundColor Cyan
Write-Host "  Documentation reorganisee" -ForegroundColor Cyan
Write-Host ""

Write-Host "=== PRET A JOUER ===" -ForegroundColor Green
Write-Host "Lancez: .\underworld.exe" -ForegroundColor Yellow
Write-Host ""
