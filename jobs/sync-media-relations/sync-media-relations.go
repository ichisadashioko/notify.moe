package main

import (
	"fmt"
	"strings"

	"github.com/animenotifier/arn"

	"github.com/animenotifier/kitsu"
	"github.com/fatih/color"
)

func main() {
	color.Yellow("Syncing media relations with Kitsu DB")

	kitsuMediaRelations := kitsu.StreamMediaRelations()
	relations := map[arn.AnimeID]*arn.AnimeRelations{}

	for mediaRelation := range kitsuMediaRelations {
		// We only care about anime for now
		if mediaRelation.Relationships.Source.Data.Type != "anime" || mediaRelation.Relationships.Destination.Data.Type != "anime" {
			continue
		}

		relationType := strings.Replace(mediaRelation.Attributes.Role, "_", " ", -1)
		animeID := mediaRelation.Relationships.Source.Data.ID
		destinationAnimeID := mediaRelation.Relationships.Destination.Data.ID

		// Confirm that the anime IDs are valid
		exists, _ := arn.DB.Exists("Anime", animeID)

		if !exists {
			continue
		}

		exists, _ = arn.DB.Exists("Anime", destinationAnimeID)

		if !exists {
			continue
		}

		fmt.Printf(
			"%s %s has %s which is %s %s\n",
			mediaRelation.Relationships.Source.Data.Type,
			animeID,
			color.GreenString(relationType),
			mediaRelation.Relationships.Destination.Data.Type,
			destinationAnimeID,
		)

		// Add anime to the global map
		relationsList, found := relations[animeID]

		if !found {
			relationsList = &arn.AnimeRelations{
				AnimeID: animeID,
				Items:   []*arn.AnimeRelation{},
			}
			relations[animeID] = relationsList
		}

		relationsList.Items = append(relationsList.Items, &arn.AnimeRelation{
			AnimeID: destinationAnimeID,
			Type:    relationType,
		})

		// for _, item := range relationsList.Items {
		// 	fmt.Println("*", item.Type, item.AnimeID)
		// }
	}

	// Save relations map
	for _, animeRelations := range relations {
		err := animeRelations.Save()

		if err != nil {
			color.Red(err.Error())
		}
	}

	color.Green("Finished.")
}