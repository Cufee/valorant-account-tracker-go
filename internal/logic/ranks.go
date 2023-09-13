package logic

import "github.com/Cufee/valorant-account-tracker-go/internal/riot/remote"

type PlayerRank struct {
	SeasonId string `json:"seasonId"`
	Tier     int    `json:"tier"`
}

func ParseLastPlayerRank(mmr remote.PlayerMMR, sortedSeasons []remote.Season) (PlayerRank, error) {
	if mmr.LatestCompetitiveUpdate.TierAfterUpdate > 0 {
		return PlayerRank{
			SeasonId: mmr.LatestCompetitiveUpdate.SeasonID,
			Tier:     mmr.LatestCompetitiveUpdate.TierAfterUpdate,
		}, nil
	}

	for _, season := range sortedSeasons {
		skillGroup, ok := mmr.QueueSkills["competitive"]
		if !ok {
			continue
		}
		selectedSeason, ok := skillGroup.SeasonalInfoBySeasonID[season.ID]
		if !ok {
			continue
		}
		if selectedSeason.Rank > 0 {
			return PlayerRank{
				SeasonId: season.ID,
				Tier:     selectedSeason.Rank,
			}, nil
		}
	}

	return PlayerRank{Tier: 0, SeasonId: mmr.LatestCompetitiveUpdate.SeasonID}, nil
}
