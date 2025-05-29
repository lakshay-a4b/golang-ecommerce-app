package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/your-username/golang-ecommerce-app/models"
	"github.com/your-username/golang-ecommerce-app/services"
	"github.com/your-username/golang-ecommerce-app/utils"
)

type ProductController struct {
	productService *services.ProductService
}

func NewProductController(productService *services.ProductService) *ProductController {
	return &ProductController{productService: productService}
}

func (pc *ProductController) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 5
	}

	products, err := pc.productService.GetPaginatedProducts(r.Context(), page, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data":  products,
		"page":  page,
		"limit": limit,
	})
}

func (pc *ProductController) GetProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := pc.productService.GetProductByID(r.Context(), id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch product")
		return
	}
	if product == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, product)
}

func (pc *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product *models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if product.Name == "" || product.Image == "" || product.Description == "" || product.Price <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "All fields are required")
		return
	}

	createdProduct, err := pc.productService.CreateProduct(r.Context(), product)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, createdProduct)
}

func (pc *ProductController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if product.Name == "" || product.Image == "" || product.Description == "" || product.Price <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "All fields are required")
		return
	}

	updatedProduct, err := pc.productService.UpdateProduct(r.Context(), id, product)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}
	if updatedProduct == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedProduct)
}

func (pc *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	deletedProduct, err := pc.productService.DeleteProduct(r.Context(), id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}
	if deletedProduct == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message":        "Product deleted successfully",
		"deletedProduct": deletedProduct,
	})
}
