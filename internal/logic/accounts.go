package logic

import (
	"github.com/Cufee/valorant-account-tracker-go/internal/riot/local"
	"github.com/Cufee/valorant-account-tracker-go/internal/riot/remote"
	"github.com/Cufee/valorant-account-tracker-go/internal/types"
)

func GetCurrentPlayerAccount() (types.Account, error) {
	// local api
	credentials, err := local.GetAccessTokens()
	if err != nil {
		return types.Account{}, err
	}
	session, err := local.GetGameSessionInfo()
	if err != nil {
		return types.Account{}, err
	}

	// external api
	versions, err := remote.GetCurrentVersions()
	if err != nil {
		return types.Account{}, err
	}
	compTiers, err := remote.GetCompetitiveTiers()
	if err != nil {
		return types.Account{}, err
	}

	// remote api
	remoteOptions := remote.RequestOptions{
		Shard:          session.Shard,
		Region:         session.Region,
		ClientPlatform: session.Platform,
		ClientVersion:  versions.RiotClientVersion,
	}
	content, err := remote.GetGameContent(remoteOptions)
	if err != nil {
		return types.Account{}, err
	}
	player, err := remote.GetPlayerInfo(credentials)
	if err != nil {
		return types.Account{}, err
	}
	playerMMR, err := remote.GetPlayerMMR(player.Sub, credentials, remoteOptions)
	if err != nil {
		return types.Account{}, err
	}

	// parse rank data
	rank, err := ParseLastPlayerRank(playerMMR, content.Seasons)
	if err != nil {
		return types.Account{}, err
	}

	var tierData remote.Tier
	for _, tier := range compTiers.Tiers {
		if tier.Tier == rank.Tier {
			tierData = tier
			break
		}
	}

	var account types.Account = types.Account{
		ID:       player.Sub,
		Username: player.Username,
		Name:     player.Acct.GameName,
		Tag:      player.Acct.TagLine,
		LastRank: types.Rank{
			Tier:  rank.Tier,
			Color: tierData.Color,
			Name:  tierData.TierName,
			Icon:  tierData.LargeIcon,
		},
	}
	return account, nil
}
