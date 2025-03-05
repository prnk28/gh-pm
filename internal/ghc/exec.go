package ghc

import "github.com/prnk28/gh-pm/internal/models"

func GetProjects() ([]models.ProjectsJson, error) {
	var projects models.ProjectsListJson
	err := QueryProjectList.ExecUnmarshal(&projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func GetProjectItems(owner string) ([]models.CardsJson, error) {
	var projectCards models.CardsListJson
	err := QueryProjectItemList.ExecUnmarshal(&projectCards)
	if err != nil {
		return nil, err
	}
	return projectCards, nil
}

func GetWhoami() (*models.UserJson, error) {
	var userInfo models.UserJson
	err := QueryUserWhoami.ExecUnmarshal(&userInfo)
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}
