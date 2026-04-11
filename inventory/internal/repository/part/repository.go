package part

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	repoModel "github.com/ChopX4/raketka/inventory/internal/repository/model"
)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(ctx context.Context, db *mongo.Database) (*repository, error) {
	collection := db.Collection("parts")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "name", Value: "text"}},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys: bson.D{{Key: "category", Value: 1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return &repository{
		collection: collection,
	}, nil
}

func SeedParts(ctx context.Context, db *mongo.Database) error {
	return seedParts(ctx, db.Collection("parts"))
}

func seedParts(ctx context.Context, collection *mongo.Collection) error {
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	now := time.Now().UTC()

	parts := []interface{}{
		repoModel.Part{
			UUID:          "11111111-1111-1111-1111-111111111111",
			Name:          "Ion Engine MK-I",
			Description:   "Compact ion engine for light exploratory spacecraft.",
			Price:         125000,
			StockQuantity: 8,
			Category:      repoModel.CategoryEngine,
			Dimensions: repoModel.Dimensions{
				Length: 2.4,
				Width:  1.2,
				Height: 1.1,
				Weight: 340,
			},
			Manufacturer: repoModel.Manufacturer{
				Name:    "Nova Propulsion",
				Country: "Japan",
				WebSite: "https://nova-propulsion.example",
			},
			Tags:      []string{"engine", "ion", "lightweight"},
			Metadata:  map[string]any{"series": "MK-I", "power_kw": 480},
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		repoModel.Part{
			UUID:          "22222222-2222-2222-2222-222222222222",
			Name:          "Cryo Fuel Cell",
			Description:   "High-density fuel module for long-range journeys.",
			Price:         42000,
			StockQuantity: 15,
			Category:      repoModel.CategoryFuel,
			Dimensions: repoModel.Dimensions{
				Length: 1.5,
				Width:  0.8,
				Height: 0.8,
				Weight: 190,
			},
			Manufacturer: repoModel.Manufacturer{
				Name:    "Orbital Energy",
				Country: "Germany",
				WebSite: "https://orbital-energy.example",
			},
			Tags:      []string{"fuel", "cryo", "long-range"},
			Metadata:  map[string]any{"capacity_l": 950, "temperature_c": -183},
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		repoModel.Part{
			UUID:          "33333333-3333-3333-3333-333333333333",
			Name:          "Panorama Porthole",
			Description:   "Reinforced observation porthole with impact shielding.",
			Price:         17500,
			StockQuantity: 21,
			Category:      repoModel.CategoryPorthole,
			Dimensions: repoModel.Dimensions{
				Length: 0.9,
				Width:  0.9,
				Height: 0.2,
				Weight: 55,
			},
			Manufacturer: repoModel.Manufacturer{
				Name:    "Stellar Glassworks",
				Country: "France",
				WebSite: "https://stellar-glassworks.example",
			},
			Tags:      []string{"porthole", "glass", "shielded"},
			Metadata:  map[string]any{"material": "transparent aluminum"},
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		repoModel.Part{
			UUID:          "44444444-4444-4444-4444-444444444444",
			Name:          "Falcon Wing S",
			Description:   "Medium lift wing segment for cargo and scout ships.",
			Price:         88000,
			StockQuantity: 6,
			Category:      repoModel.CategoryWing,
			Dimensions: repoModel.Dimensions{
				Length: 3.7,
				Width:  1.4,
				Height: 0.4,
				Weight: 410,
			},
			Manufacturer: repoModel.Manufacturer{
				Name:    "AeroForge",
				Country: "United States",
				WebSite: "https://aeroforge.example",
			},
			Tags:      []string{"wing", "cargo", "lift"},
			Metadata:  map[string]any{"load_class": "S"},
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}

	_, err = collection.InsertMany(ctx, parts)
	return err
}
