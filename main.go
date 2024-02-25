package main

import (
	"fmt"

	"time"

	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator/v10"

	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	swaggerFiles "github.com/swaggo/files" // swagger embed files

	_ "api/sal/docs"
)

// define product struct type
type Product struct {
	SkuId       string    `json:"skuId" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" validate:"gt=0"`
	CreatedAt   time.Time `json:"createdAt" swaggerignore:"true"`
}

var Products = make(map[string]map[string]Product)

func getMerchantId(ctx *gin.Context) string {
	return ctx.Param("merchantId")
}

func getSkuId(ctx *gin.Context) string {
	return ctx.Param("skuId")
}

func getMerchantProducts(merchantId string) (map[string]Product, bool) {
	merchantProducts, merchantIdExists := Products[merchantId]
	return merchantProducts, merchantIdExists
}

// GetProducts	 godoc
// @Summary      Get Merchant products
// @Description  List Merchant Products.
// @Tags         products
// @Produce      json
// @Success      200  {object}  map[string]Product
// @Router       /products/{merchantId} [get]
// @Param		merchantId path string true "merchant id"
func GetProducts(ctx *gin.Context) {
	merchantId := getMerchantId(ctx)
	merchantProducts, merchantIdExists := getMerchantProducts(merchantId)

	//check if merchant id exists
	if !merchantIdExists {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Merchant not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, merchantProducts)
}

// CreateProduct	 godoc
// @Summary Create product
// @Description  Create Merchant Product.
// @Tags         products
// @Accept		 json
// @Produce      json
// @Success      202  {object}  Product
// @Router       /products/{merchantId} [post]
// @Param		merchantId path string true "merchant id"
// @Param product body Product true "The product to create"
func CreateProduct(ctx *gin.Context) {
	var product Product
	err := ctx.BindJSON(&product)

	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid product data"})
		fmt.Println(err.Error())
		return
	}

	//validate the product data
	validate := validator.New()
	validationErr := validate.Struct(product)
	if validationErr != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": validationErr.Error()})
		return
	}

	//set the created at time to now
	product.CreatedAt = time.Now()

	//get merchantId
	merchantId := getMerchantId(ctx)

	// check if the inner map exists for the merchantId
	_, innerExists := Products[merchantId]
	if !innerExists {
		// create the inner map if it does not exist
		Products[merchantId] = make(map[string]Product)
	}

	//check if product sku exists
	_, skuExists := Products[merchantId][product.SkuId]
	if skuExists {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Product sku exists"})
		return
	}

	//save new product
	Products[merchantId][product.SkuId] = product
	ctx.IndentedJSON(http.StatusAccepted, product)
}

// UpdateProduct	 godoc
// @Summary Update product
// @Description  Update Merchant Product.
// @Tags         products
// @Accept		 json
// @Produce      json
// @Success      202  {object}  Product
// @Router       /products/{merchantId}/{skuId} [put]
// @Param		merchantId path string true "merchant id"
// @Param		skuId path string true "product sku"
// @Param product body Product true "The product to update"
func UpdateProduct(ctx *gin.Context) {
	merchantId := getMerchantId(ctx)
	merchantProducts, merchantIdExists := getMerchantProducts(merchantId)

	//check if merchant id exists
	if !merchantIdExists {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Merchant not found"})
		return
	}

	//get sku id
	skuId := getSkuId(ctx)

	//get product by sku
	product, skuExists := merchantProducts[skuId]

	//check if sku exists
	if !skuExists {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	err := ctx.BindJSON(&product)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid product data"})
		fmt.Println(err.Error())
		return
	}

	//validate the product data
	validate := validator.New()
	validationErr := validate.Struct(product)
	if validationErr != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": validationErr.Error()})
		return
	}

	//check if sku changed and it exists
	updatedSku := product.SkuId
	_, updatedSkuExists := merchantProducts[updatedSku]
	if skuId != updatedSku && updatedSkuExists {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Product sku exists"})
		return
	}

	//update product by sku
	merchantProducts[skuId] = product
	ctx.IndentedJSON(http.StatusAccepted, product)
}

// DeleteProduct	 godoc
// @Summary Delete product
// @Description  Delete Merchant Product.
// @Tags         products
// @Accept		 json
// @Produce      json
// @Success      200
// @Router       /products/{merchantId}/{skuId} [delete]
// @Param		merchantId path string true "merchant id"
// @Param		skuId path string true "product sku"
func DeleteProduct(ctx *gin.Context) {
	merchantId := getMerchantId(ctx)
	merchantProducts, merchantIdExists := getMerchantProducts(merchantId)

	//check if merchant id exists
	if !merchantIdExists {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Merchant not found"})
		return
	}

	//get sku id
	skuId := getSkuId(ctx)

	//get product by sku
	_, skuExists := merchantProducts[skuId]

	//check if sku exists
	if !skuExists {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	//delete product by sku
	delete(merchantProducts, skuId)
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

// @title           SAL Merchant Product API
// @version         1.0
// @description     A Product management API for SAL Merchants.

// @contact.name   SAL
// @contact.url    zraelwalker@gmail.com
// @contact.email  zraelwalker@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /api/v1
func main() {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		v1.GET("/products/:merchantId", GetProducts)
		v1.POST("/products/:merchantId", CreateProduct)
		v1.PUT("/products/:merchantId/:skuId", UpdateProduct)
		v1.DELETE("products/:merchantId/:skuId", DeleteProduct)
	}

	router.Run("localhost:8081")
}