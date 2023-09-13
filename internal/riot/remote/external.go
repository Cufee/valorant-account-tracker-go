package remote

type Episode struct {
	UUID            string `json:"uuid"`
	Tiers           []Tier `json:"tiers"`
	AssetPath       string `json:"assetPath"`
	AssetObjectName string `json:"assetObjectName"`
}

type Tier struct {
	Tier     int    `json:"tier"`
	TierName string `json:"tierName"`

	Division     string `json:"division"`
	DivisionName string `json:"divisionName"`

	Color           string `json:"color"`
	BackgroundColor string `json:"backgroundColor"`

	SmallIcon            string `json:"smallIcon"`
	LargeIcon            string `json:"largeIcon"`
	RankTriangleUpIcon   string `json:"rankTriangleUpIcon"`
	RankTriangleDownIcon string `json:"rankTriangleDownIcon"`
}

func GetCompetitiveTiers() (Episode, error) {
	var data struct {
		Data []Episode `json:"data"`
	}
	err := plainRequest("GET", "https://valorant-api.com/v1/competitivetiers", &data)
	if err != nil {
		return Episode{}, err
	}

	var episode Episode
	// Is there a better way to get the last item in a map?
	for _, v := range data.Data {
		episode = v
	}

	return episode, nil
}

type Versions struct {
	ManifestID        string `json:"manifestId"`
	Branch            string `json:"branch"`
	Version           string `json:"version"`
	BuildVersion      string `json:"buildVersion"`
	EngineVersion     string `json:"engineVersion"`
	RiotClientVersion string `json:"riotClientVersion"`
	RiotClientBuild   string `json:"riotClientBuild"`
	BuildDate         string `json:"buildDate"`
}

func GetCurrentVersions() (Versions, error) {
	var data struct {
		Data Versions `json:"data"`
	}
	err := plainRequest("GET", "https://valorant-api.com/v1/version", &data)
	if err != nil {
		return Versions{}, err
	}
	return data.Data, nil
}
