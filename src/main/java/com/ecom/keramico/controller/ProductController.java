package com.ecom.keramico.controller;

import com.ecom.keramico.model.Product;
import com.ecom.keramico.service.ProductService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Optional;

@RestController
@RequestMapping("/api/products/")
@Slf4j
public class ProductController {

    @Autowired
    private ProductService productService;

    @GetMapping
    public ResponseEntity<List<Product>> showAllProducts() {
        log.info("ProductController showAllProducts()");
        return ResponseEntity.ok(productService.getAllProducts());
    }

    @GetMapping("/{id}")
    public ResponseEntity<Optional<Product>> showProductDetails(@PathVariable Long id) {
        Optional<Product> product = productService.getProductById(id);
        log.info("/product/{}", id);
        return product.isEmpty() ? ResponseEntity.notFound().build() : ResponseEntity.ok(product);
    }

    @GetMapping("/category/{category}")
    public ResponseEntity<List<Product>> showProductsByCategory(@PathVariable String category) {
        log.info("/category/{}", category);
        return ResponseEntity.ok(productService.getAllProductsByCategory(category));
    }

    @GetMapping("/search/{name}")
    public ResponseEntity<List<Product>> searchProductByName(@PathVariable String name) {
        log.info("search/{}", name);
        return ResponseEntity.ok(productService.getProductByName(name));
    }


    /* TODO: ADMIN */
    @PostMapping
    public ResponseEntity<Product> addNewProduct(@RequestBody Product product) {
        log.info("/products");
        return ResponseEntity.ok(productService.createProduct(product));
    }

    @PutMapping("/{id}")
    public ResponseEntity<Product> updateProduct(@PathVariable Long id, @RequestBody Product product) {
        return ResponseEntity.ok(productService.updateProduct(id, product));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteProduct(@PathVariable Long id) {
        productService.deleteProduct(id);
        return ResponseEntity.noContent().build();
    }
}
