package repo

import (
	"context"
	"database/sql"
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
// 	"sku": "", // not null, minLength 1, maxLength 30
// 	"category": "", /** not null, enum of:
// 			- "Clothing"
// 			- "Accessories"
// 			- "Footwear"
// 			- "Beverages"
// 			*/
// 	"imageUrl": "", // not null, should be url
// 	"notes":"", // not null, minLength 1, maxLength 200
// 	"price":1, // not null, min: 1
// 	"stock": 1, // not null, min: 0, max: 100000
// 	"location": "", // not null, minLength 1, maxLength 200
// 	"isAvailable": true // not null
// }

func (cr *productRepo) AddProduct(ctx context.Context, sub string, product entity.Product) (dto.ResAddOrUpdateProduct, error) {
	// add product
	q := `INSERT INTO products (user_id, name, sku, category, image_url, notes, price, stock, location, is_available, created_at)
	VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, EXTRACT(EPOCH FROM now())::bigint) RETURNING id, created_at`

	var id string
	var createdAt int64
	err := cr.conn.QueryRow(ctx, q, sub, product.Name, product.SKU, product.Category, product.ImageUrl, product.Notes, product.Price, product.Stock, product.Location, product.IsAvailable).Scan(&id, &createdAt)
	if err != nil {
		return dto.ResAddOrUpdateProduct{}, err
	}

	createdAtTime := time.Unix(createdAt, 0) // Convert createdAt from int64 to time.Time
	return dto.ResAddOrUpdateProduct{ID: id, CreatedAt: timepkg.TimeToISO8601(createdAtTime)}, nil
}

func (cr *productRepo) GetProduct(ctx context.Context, param dto.ParamGetProduct, sub string) ([]dto.ResGetProduct, error) {
	var query strings.Builder

	query.WriteString("SELECT id, name, sku, category, image_url, stock, notes, price, location, is_available, created_at FROM products WHERE 1=1 ")

	// param id
	if param.ID != "" {
		id, err := strconv.Atoi(param.ID)
		if err != nil {
			return nil, err
		}
		query.WriteString(fmt.Sprintf("AND id = %d ", id))
	}

	// param name: case insensitive, if een in between the name it will be included
	if param.Name != "" {
		query.WriteString(fmt.Sprintf("AND LOWER(name) LIKE LOWER('%s') ", fmt.Sprintf("%%%s%%", param.Name)))
	}

	// param isAvailable value should be true / false / empty
	// if empty, it will show all product
	// show param isAvailable
	fmt.Printf("param.IsAvailable: %s\n", param.IsAvailable)
	if param.IsAvailable == "true" {
		query.WriteString("AND is_available = true ")
	} else if param.IsAvailable == "false" {
		query.WriteString("AND is_available = false ")
	} else {
		query.WriteString("AND is_available = true ")
	}

	// param category
	if param.Category != "" {
		query.WriteString(fmt.Sprintf("AND category = '%s' ", param.Category))
	}

	// param sku
	if param.SKU != "" {
		query.WriteString(fmt.Sprintf("AND sku = '%s' ", param.SKU))
	}

	// param price sort by asc or desc, if value is wrong, just ignore the param
	if param.Price == "asc" {
		query.WriteString("ORDER BY price ASC ")
	} else if param.Price == "desc" {
		query.WriteString("ORDER BY price DESC ")
	}

	// param inStock value should be true / false / empty
	// if empty, it will show all product
	// show param inStock
	fmt.Printf("param.InStock: %s\n", param.InStock)
	if param.InStock == "true" {
		query.WriteString("AND stock > 0 ")
	} else if param.InStock == "false" {
		query.WriteString("AND stock = 0 ")
	} else {
		query.WriteString("AND stock >= 0 ")
	}

	// param createdAt sort by created time asc or desc, if value is wrong, just ignore the param
	if param.CreatedAt == "asc" {
		query.WriteString("ORDER BY created_at ASC ")
	} else if param.CreatedAt == "desc" {
		query.WriteString("ORDER BY created_at DESC ")
	}

	if param.Search != "" {
		query.WriteString(fmt.Sprintf("AND LOWER(name) LIKE LOWER('%s') ", fmt.Sprintf("%%%s%%", param.Search)))
	}

	// limit and offset
	if param.Limit == 0 {
		param.Limit = 5
	}

	query.WriteString(fmt.Sprintf("LIMIT %d OFFSET %d", param.Limit, param.Offset))

	// show query
	fmt.Println(query.String())

	rows, err := cr.conn.Query(ctx, query.String()) // Replace $1 with sub
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]dto.ResGetProduct, 0, 10)
	for rows.Next() {
		var imageUrl sql.NullString
		var createdAt int64

		result := dto.ResGetProduct{}
		err := rows.Scan(
			&result.ID,
			&result.Name,
			&result.SKU,
			&result.Category,
			&imageUrl,
			&result.Stock,
			&result.Notes,
			&result.Price,
			&result.Location,
			&result.IsAvailable,
			&createdAt)
		if err != nil {
			return nil, err
		}

		result.ImageUrl = imageUrl.String
		result.CreatedAt = timepkg.TimeToISO8601(time.Unix(createdAt, 0))
		results = append(results, result)
	}

	return results, nil
}

func (cr *productRepo) GetProductShop(ctx context.Context, param dto.ParamGetProductShop) ([]dto.ResGetProductShop, error) {
	var query strings.Builder

	query.WriteString("SELECT id, name, sku, category, image_url, stock, price, location, created_at FROM products WHERE 1=1 ")

	// param name: case insensitive, if een in between the name it will be included
	if param.Name != "" {
		query.WriteString(fmt.Sprintf("AND LOWER(name) LIKE LOWER('%s') ", fmt.Sprintf("%%%s%%", param.Name)))
	}

	// param category
	if param.Category != "" {
		query.WriteString(fmt.Sprintf("AND category = '%s' ", param.Category))
	}

	// param sku
	if param.SKU != "" {
		query.WriteString(fmt.Sprintf("AND sku = '%s' ", param.SKU))
	}

	// param price sort by asc or desc, if value is wrong, just ignore the param
	if param.Price == "asc" {
		query.WriteString("ORDER BY price ASC ")
	} else if param.Price == "desc" {
		query.WriteString("ORDER BY price DESC ")
	}

	// param inStock value should be true / false
	if param.InStock {
		query.WriteString("AND stock > 0 ")
	} else if !param.InStock {
		query.WriteString("AND stock = 0 ")
	}

	// limit and offset
	if param.Limit == 0 {
		param.Limit = 5
	}

	query.WriteString(fmt.Sprintf("LIMIT %d OFFSET %d", param.Limit, param.Offset))

	// show query
	// fmt.Println(query.String())

	rows, err := cr.conn.Query(ctx, query.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]dto.ResGetProductShop, 0, 10)
	for rows.Next() {
		var imageUrl sql.NullString
		var createdAt int64

		result := dto.ResGetProductShop{}
		err := rows.Scan(
			&result.ID,
			&result.Name,
			&result.SKU,
			&result.Category,
			&imageUrl,
			&result.Stock,
			&result.Price,
			&result.Location,
			&createdAt)
		if err != nil {
			return nil, err
		}

		result.ImageUrl = imageUrl.String
		result.CreatedAt = timepkg.TimeToISO8601(time.Unix(createdAt, 0))
		results = append(results, result)
	}

	return results, nil
}

func (cr *productRepo) UpdateProduct(ctx context.Context, id, sub string, product entity.Product) error {
	q := `UPDATE products SET name=$1, sku=$2, category=$3, image_url=$4, notes=$5, price=$6, stock=$7, location=$8, is_available=$9 WHERE id=$10 AND user_id=$11`

	_, err := cr.conn.Exec(ctx, q, product.Name, product.SKU, product.Category, product.ImageUrl, product.Notes, product.Price, product.Stock, product.Location, product.IsAvailable, id, sub)
	if err != nil {
		return err
	}

	return nil
}

func (cr *productRepo) DeleteProduct(ctx context.Context, sub, id string) error {
	q := `DELETE FROM products WHERE id=$1 AND user_id=$2`

	_, err := cr.conn.Exec(ctx, q, id, sub)
	if err != nil {
		return err
	}

	return nil
}
