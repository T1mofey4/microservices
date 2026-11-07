package main

import (
	"context"
	"errors"
	"sync"

	inventoryV1 "github.com/T1mofey4/microservices/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrPartNotFound = errors.New("part not found")

type inventoryStorage struct {
	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func NewInventoryStorage(parts []*inventoryV1.Part) *inventoryStorage {
	partsByUuid := make(map[string]*inventoryV1.Part, len(parts))

	for _, p := range parts {
		partsByUuid[p.Uuid] = p
	}

	return &inventoryStorage{
		parts: partsByUuid,
	}
}

// Возвращает информацию о детали по её UUID
func (s *inventoryStorage) GetPart(ctx context.Context, uuid string) (*inventoryV1.Part, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.parts[uuid]
	if !ok {
		return nil, ErrPartNotFound
	}

	return p, nil
}

// Возвращает список деталей с возможностью фильтрации
func (s *inventoryStorage) ListParts(ctx context.Context, filter *inventoryV1.PartsFilter) ([]*inventoryV1.Part, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Получаем все детали из хранилища
	partSlice := make([]*inventoryV1.Part, 0, len(s.parts))
	for _, p := range s.parts {
		partSlice = append(partSlice, p)
	}

	// Если фильтр не задан (пустой), то возвращаем весь список запчастей
	if filter == nil {
		return partSlice, nil
	}

	// Проходим по всем фильтрам
	result := filterByUUID(filter.Uuids, partSlice)
	result = filterByName(filter.Names, result)
	result = filterByCategory(filter.Categories, result)
	result = filterByCountry(filter.ManufacturerCountries, result)
	result = filterByTag(filter.Tags, result)

	return result, nil
}

// Фильтр по UUID
func filterByUUID(filter []string, parts []*inventoryV1.Part) []*inventoryV1.Part {
	// Если фильтр пустой, возвращаем все части
	if len(filter) == 0 {
		return parts
	}

	set := make(map[string]struct{}, len(filter))
	for _, uuid := range filter {
		set[uuid] = struct{}{}
	}

	result := make([]*inventoryV1.Part, 0, len(parts))
	for _, p := range parts {
		if _, ok := set[p.Uuid]; ok {
			result = append(result, p)
		}
	}

	return result
}

// Фильтр по имени
func filterByName(filter []string, parts []*inventoryV1.Part) []*inventoryV1.Part {
	if len(filter) == 0 {
		return parts
	}

	set := make(map[string]struct{}, len(filter))
	for _, name := range filter {
		set[name] = struct{}{}
	}

	result := make([]*inventoryV1.Part, 0, len(parts))
	for _, p := range parts {
		if _, ok := set[p.Name]; ok {
			result = append(result, p)
		}
	}

	return result
}

// Фильтр по категории
func filterByCategory(filter []inventoryV1.Category, parts []*inventoryV1.Part) []*inventoryV1.Part {
	if len(filter) == 0 {
		return parts
	}

	set := make(map[inventoryV1.Category]struct{}, len(filter))
	for _, category := range filter {
		set[category] = struct{}{}
	}

	result := make([]*inventoryV1.Part, 0, len(parts))
	for _, p := range parts {
		if _, ok := set[p.Category]; ok {
			result = append(result, p)
		}
	}

	return result
}

// Фильтр по странам
func filterByCountry(filter []string, parts []*inventoryV1.Part) []*inventoryV1.Part {
	if len(filter) == 0 {
		return parts
	}

	set := make(map[string]struct{}, len(filter))
	for _, country := range filter {
		set[country] = struct{}{}
	}

	result := make([]*inventoryV1.Part, 0, len(parts))
	for _, p := range parts {
		if _, ok := set[p.Manufacturer.Country]; ok {
			result = append(result, p)
		}
	}

	return result
}

// Фильтр по тегам
func filterByTag(filter []string, parts []*inventoryV1.Part) []*inventoryV1.Part {
	if len(filter) == 0 {
		return parts
	}

	set := make(map[string]struct{}, len(parts))
	for _, tag := range filter {
		set[tag] = struct{}{}
	}

	result := make([]*inventoryV1.Part, 0, len(parts))
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

func testParts() []*inventoryV1.Part {
	parts := []*inventoryV1.Part{
		{Uuid: "11111111-1111-1111-1111-111111111111",
			Name:          "Main Engine",
			Price:         111.00,
			StockQuantity: 111,
			Category:      inventoryV1.Category_CATEGORY_ENGINE,
			Dimensions: &inventoryV1.Dimensions{
				Length: 0.1,
				Width:  0.2,
				Height: 0.3,
				Weight: 0.5,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    "Test Manufacturer Corp",
				Country: "USA",
				Website: "http://testcorp.com",
			},
			Tags: []string{"electronics", "tools"},
			Metadata: map[string]*inventoryV1.MetadataValue{
				"Warranty": {
					Value: &inventoryV1.MetadataValue_StringValue{
						StringValue: "1 year"},
				},
			},
			CreatedAt: timestamppb.Now(),
			UpdatedAt: timestamppb.Now(),
		},
		{Uuid: "22222222-2222-2222-2222-222222222222",
			Name:          "Main Fuel Tank",
			Price:         55.50,
			StockQuantity: 500,
			Category:      inventoryV1.Category_CATEGORY_FUEL,
			Dimensions: &inventoryV1.Dimensions{
				Length: 1.0,
				Width:  0.5,
				Height: 0.2,
				Weight: 1.2,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    "WoodWorks Inc",
				Country: "CAN",
				Website: "http://woodworks.ca",
			},
			Tags: []string{"wood", "hardware"},
			Metadata: map[string]*inventoryV1.MetadataValue{
				"Material": {
					Value: &inventoryV1.MetadataValue_StringValue{
						StringValue: "Oak"},
				},
			},
			CreatedAt: timestamppb.Now(),
			UpdatedAt: timestamppb.Now(),
		},
		{Uuid: "33333333-3333-3333-3333-333333333333",
			Name:          "Part3",
			Price:         1200.00,
			StockQuantity: 10,
			Category:      inventoryV1.Category_CATEGORY_PORTHOLE,
			Dimensions: &inventoryV1.Dimensions{
				Length: 2.0,
				Width:  1.5,
				Height: 0.5,
				Weight: 25.0,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    "ElectroTech Global",
				Country: "JPN",
				Website: "http://electrotech.jp",
			},
			Tags: []string{"electronics", "high-end"},
			Metadata: map[string]*inventoryV1.MetadataValue{
				"Voltage": {
					Value: &inventoryV1.MetadataValue_StringValue{
						StringValue: "220V"},
				},
			},
			CreatedAt: timestamppb.Now(),
			UpdatedAt: timestamppb.Now(),
		},
		{Uuid: "44444444-4444-4444-4444-444444444444",
			Name:          "Part4",
			Price:         12.99,
			StockQuantity: 2000,
			Category:      inventoryV1.Category_CATEGORY_FUEL,
			Dimensions: &inventoryV1.Dimensions{
				Length: 0.1,
				Width:  0.1,
				Height: 0.1,
				Weight: 0.05,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    "Fasteners Ltd",
				Country: "USA",
				Website: "http://fasteners.com",
			},
			Tags: []string{"fasteners", "small"},
			Metadata: map[string]*inventoryV1.MetadataValue{
				"Size": {
					Value: &inventoryV1.MetadataValue_StringValue{
						StringValue: "M3 x 10mm"},
				},
			},
			CreatedAt: timestamppb.Now(),
			UpdatedAt: timestamppb.Now(),
		},
	}
	return parts
}
