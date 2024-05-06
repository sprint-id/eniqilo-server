package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/entity"
	timepkg "github.com/sprint-id/eniqilo-server/pkg/time"
)

type productRepo struct {
	conn *pgxpool.Pool
}

func newProductRepo(conn *pgxpool.Pool) *productRepo {
	return &productRepo{conn}
}

// {
// 	"name": "", // not null, minLength 1, maxLength 30
// 	"race": "", /** not null, enum of:
// 			- "Persian"
// 			- "Maine Coon"
// 			- "Siamese"
// 			- "Ragdoll"
// 			- "Bengal"
// 			- "Sphynx"
// 			- "British Shorthair"
// 			- "Abyssinian"
// 			- "Scottish Fold"
// 			- "Birman" */
// 	"sex": "", // not null, enum of: "male" / "female"
// 	"ageInMonth": 1, // not null, min: 1, max: 120082
// 	"description":"" // not null, minLength 1, maxLength 200
// 	"imageUrls":[ // not null, minItems: 1, items: not null, should be url
// 		"","",""
// 	]
// }

// ResAddCat struct {
// 	ID        string `json:"id"`
// 	CreatedAt string `json:"createdAt"`
// }

func (cr *productRepo) AddCat(ctx context.Context, sub string, cat entity.Cat) (dto.ResAddCat, error) {
	// add cat
	q := `INSERT INTO cats (user_id, name, race, sex, age_in_month, description, image_urls, created_at)
	VALUES ( $1, $2, $3, $4, $5, $6, $7, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`

	image_urls := "{" + strings.Join(cat.ImageUrls, ",") + "}" // Format image URLs as a PostgreSQL array

	var id string
	err := cr.conn.QueryRow(ctx, q, sub,
		cat.Name, cat.Race, cat.Sex, cat.AgeInMonth, cat.Description, image_urls).Scan(&id)
	if err != nil {
		return dto.ResAddCat{}, err
	}

	createdAt := time.Now()
	return dto.ResAddCat{ID: id, CreatedAt: timepkg.TimeToISO8601(createdAt)}, nil
}

func (cr *productRepo) GetCat(ctx context.Context, param dto.ParamGetCat, sub string) ([]dto.ResGetCat, error) {
	var query strings.Builder

	query.WriteString("SELECT id, name, race, sex, age_in_month, description, image_urls, has_matched, created_at FROM cats WHERE 1=1 ")

	if param.Owned {
		query.WriteString(fmt.Sprintf("AND user_id = '%s' ", sub))
	}

	// param id
	if param.ID != "" {
		id, err := strconv.Atoi(param.ID)
		if err != nil {
			return nil, err
		}
		query.WriteString(fmt.Sprintf("AND id = %d ", id))
	}

	// param race
	// Define a map to store the mappings between race names with spaces and underscores
	var raceMap = map[string]string{
		"Maine Coon":        "Maine_Coon",
		"British Shorthair": "British_Shorthair",
		"Scottish Fold":     "Scottish_Fold",
		// Add more mappings as needed
	}

	if param.Race != "" {
		// Check if the race name needs to be transformed
		if mappedRace, ok := raceMap[param.Race]; ok {
			param.Race = mappedRace
		}

		query.WriteString(fmt.Sprintf("AND race = '%s' ", param.Race))
	}

	// param sex
	if param.Sex != "" {
		query.WriteString(fmt.Sprintf("AND sex = '%s' ", param.Sex))
	}

	// param age in month can be equal, grater than and lower than
	// Assuming param.AgeInMonth is a string containing the ageInMonth parameter with comparison operator, e.g., ">39941"
	if param.AgeInMonth != "" {
		operator := ""
		ageStr := ""

		// Check for the comparison operator at the beginning of the string
		switch {
		case strings.HasPrefix(param.AgeInMonth, ">"):
			operator = ">"
			ageStr = param.AgeInMonth[1:] // Remove the ">" prefix
		case strings.HasPrefix(param.AgeInMonth, "<"):
			operator = "<"
			ageStr = param.AgeInMonth[1:] // Remove the "<" prefix
		case strings.HasPrefix(param.AgeInMonth, "="):
			operator = "="
			ageStr = param.AgeInMonth[1:] // Remove the "=" prefix
		default:
			// Handle exact value case if needed
			operator = "="
			ageStr = param.AgeInMonth
		}

		// Convert ageStr to an integer
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			return nil, err
		}

		// Add the ageInMonth condition to the query based on the operator
		switch operator {
		case ">":
			query.WriteString(fmt.Sprintf("AND age_in_month > %d ", age))
		case "<":
			query.WriteString(fmt.Sprintf("AND age_in_month < %d ", age))
		case "=":
			query.WriteString(fmt.Sprintf("AND age_in_month = %d ", age))
		default:
			// Handle unsupported operator
			return nil, errors.New("unsupported operator")
		}
	}

	if param.Search != "" {
		query.WriteString(fmt.Sprintf("AND LOWER(name) LIKE LOWER('%s') ", fmt.Sprintf("%%%s%%", param.Search)))
	}

	if param.HasMatched {
		query.WriteString(fmt.Sprintf("AND has_matched = %t ", param.HasMatched))
	}

	// limit and offset
	if param.Limit == 0 {
		param.Limit = 5
	}

	query.WriteString(fmt.Sprintf("LIMIT %d OFFSET %d", param.Limit, param.Offset))

	// show query
	// fmt.Println(query.String())

	rows, err := cr.conn.Query(ctx, query.String()) // Replace $1 with sub
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]dto.ResGetCat, 0, 10)
	for rows.Next() {
		var imageUrl sql.NullString
		var createdAt int64
		var description string
		var race string

		result := dto.ResGetCat{}
		err := rows.Scan(
			&result.ID,
			&result.Name,
			&race,
			&result.Sex,
			&result.AgeInMonth,
			&description,
			&imageUrl,
			&result.HasMatched,
			&createdAt)
		if err != nil {
			return nil, err
		}

		// Transform race
		var raceMap = map[string]string{
			"Maine_Coon":        "Maine Coon",
			"British_Shorthair": "British Shorthair",
			"Scottish_Fold":     "Scottish Fold",
			// Add more mappings as needed
		}

		if mappedRace, ok := raceMap[race]; ok {
			result.Race = mappedRace
		} else {
			result.Race = race // Assign the original race value if it is not found in the map
		}

		result.ImageUrls = strings.Split(imageUrl.String, ",")
		result.CreatedAt = timepkg.TimeToISO8601(time.Unix(createdAt, 0))
		result.Description = description
		results = append(results, result)
	}

	return results, nil
}

func (cr *productRepo) GetCatByID(ctx context.Context, id, sub string) (dto.ResGetCat, error) {
	q := `SELECT id,
		name,
		race,
		sex,
		age_in_month,
		description,
		image_urls,
		EXISTS (
			SELECT 1 FROM match_cats m WHERE m.user_cat_id = c.id AND m.user_id = $1
		) AS has_matched,
		created_at
	FROM cats c WHERE id = $2`

	var imageUrl sql.NullString
	var createdAt int64
	var description string

	result := dto.ResGetCat{}
	err := cr.conn.QueryRow(ctx, q, sub, id).Scan(&result.ID, &result.Name, &result.Race, &result.Sex, &result.AgeInMonth, &description, &imageUrl, &result.HasMatched, &createdAt)
	if err != nil {
		return dto.ResGetCat{}, err
	}

	result.ImageUrls = strings.Split(imageUrl.String, ",")
	result.CreatedAt = timepkg.TimeToISO8601(time.Unix(createdAt, 0))
	result.Description = description

	return result, nil
}

func (cr *productRepo) UpdateCat(ctx context.Context, id, sub string, cat entity.Cat) error {
	q := `UPDATE cats SET 
		name = $1,
		race = $2,
		sex = $3,
		age_in_month = $4,
		description = $5,
		image_urls = $6
	WHERE
		id = $7 AND user_id = $8`

	image_urls := "{" + strings.Join(cat.ImageUrls, ",") + "}" // Format image URLs as a PostgreSQL array

	_, err := cr.conn.Exec(ctx, q,
		cat.Name, cat.Race, cat.Sex, cat.AgeInMonth, cat.Description, image_urls, id, sub)

	if err != nil {
		return err
	}

	return nil
}

func (cr *productRepo) DeleteCat(ctx context.Context, id string, sub string) error {
	q := `DELETE FROM cats WHERE id = $1 AND user_id = $2`

	// log id and sub
	// fmt.Println("id: ", id)
	// fmt.Println("sub: ", sub)

	_, err := cr.conn.Exec(ctx, q, id, sub)
	if err != nil {
		return err
	}

	return nil
}
