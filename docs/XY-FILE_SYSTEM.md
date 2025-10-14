# Module de gestion des fichiers de tickets

## Objectif
Ce module gère la création, l'impression et l'archivage des tickets d'impression.

## Flux de travail
1. Création d'un fichier temporaire pour l'impression
2. Envoi à l'imprimante via la commande
3. Suppression du fichier temporaire
4. Archivage des métadonnées et du contenu

## Structure des fichiers
- `/tmp`: Fichiers temporaires pour impression (nettoyés immédiatement)
- `/tickets`: Archive organisée par date (YYYY/MM/DD)
- `/tickets/latest.json`: Dernier ticket imprimé (pour debug rapide)

## Configuration
- `PRINT_KEEP_ARCHIVE`: Activer/désactiver l'archivage (default: true)
- `PRINT_ARCHIVE_RETENTION_DAYS`: Nombre de jours de conservation (default: 90)

## Nettoyage
- Fichiers temporaires: supprimés après impression
- Fichiers temporaires orphelins: supprimés après 1 heure
- Archives: supprimées après X jours (configurable)

## Utilisation
```go
// Injecter dans PrintService
printService.PrintTicket("Contenu", "XP-80C", map[string]string{
    "order_id": "12345",
    "type": "kitchen"
})
```