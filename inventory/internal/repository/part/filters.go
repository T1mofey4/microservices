package part

import (
	"github.com/T1mofey4/microservices/inventory/internal/model"
	"github.com/T1mofey4/microservices/inventory/internal/repository/converter"
	repoModel "github.com/T1mofey4/microservices/inventory/internal/repository/model"
)

func filterByUUID(filter []string, parts []repoModel.Part) []repoModel.Part {
	// Если фильтр пустой, возвращаем все части
	if len(filter) == 0 {
		return parts
	}

	set := make(map[string]struct{}, len(filter))
	for _, uuid := range filter {
		set[uuid] = struct{}{}
	}

	result := make([]repoModel.Part, 0, len(parts))
	for _, p := range parts {
		if _, ok := set[p.UUID]; ok {
			result = append(result, p)
		}
	}

	return result
}

func filterByName(filter []string, parts []repoModel.Part) []repoModel.Part {
	if len(filter) == 0 {
		return parts
	}

	set := make(map[string]struct{}, len(filter))
	for _, name := range filter {
		set[name] = struct{}{}
	}

	result := make([]repoModel.Part, 0, len(parts))
	for _, p := range parts {
		if _, ok := set[p.Name]; ok {
			result = append(result, p)
		}
	}

	return result
}

func filterByCategory(filter []model.Category, parts []repoModel.Part) []repoModel.Part {
	if len(filter) == 0 {
		return parts
	}

	set := make(map[repoModel.Category]struct{}, len(filter))
	for _, category := range filter {
		repoCategory := converter.CategoryToRepoModel(category)
		set[repoCategory] = struct{}{}
	}

	result := make([]repoModel.Part, 0, len(parts))
	for _, p := range parts {
		if _, ok := set[p.Category]; ok {
			result = append(result, p)
		}
	}

	return result
}

func filterByCountry(filter []string, parts []repoModel.Part) []repoModel.Part {
	if len(filter) == 0 {
		return parts
	}

	set := make(map[string]struct{}, len(filter))
	for _, country := range filter {
		set[country] = struct{}{}
	}

	result := make([]repoModel.Part, 0, len(parts))
	for _, p := range parts {
		if _, ok := set[p.Manufacturer.Country]; ok {
			result = append(result, p)
		}
	}

	return result
}

func filterByTag(filter []string, parts []repoModel.Part) []repoModel.Part {
	if len(filter) == 0 {
		return parts
	}

	set := make(map[string]struct{}, len(parts))
	for _, tag := range filter {
		set[tag] = struct{}{}
	}

	result := make([]repoModel.Part, 0, len(parts))
	for _, p := range parts {
		for _, tag := range p.Tags {
			if _, ok := set[tag]; ok {
				result = append(result, p)
				break
			}
		}
	}

	return result
}
