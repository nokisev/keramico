package com.ecom.keramico.service;

import com.ecom.keramico.model.Category;
import com.ecom.keramico.model.Product;
import com.ecom.keramico.repository.CategoryRepository;
import com.ecom.keramico.repository.ProductRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Optional;

@Service
public class ProductService {

    @Autowired
    private ProductRepository productRepository;
    @Autowired
    private CategoryRepository categoryRepository;

    public List<Product> getAllProducts() {
        return productRepository.findAll();
    }

    public Optional<Product> getProductById(Long id) {
        if (productRepository.findById(id).isEmpty()) {
            return Optional.empty();
        }
        return productRepository.findById(id);
    }

    public List<Product> getAllProductsByCategory(String category) {
        return productRepository.findAllByCategory(categoryRepository.findByName(category));
    }

    public List<Product> getProductByName(String name) {
        return productRepository.findAllByName(name);
    }

    public Product createProduct(Product product) {
        categoryRepository.save(product.getCategory());
        return productRepository.save(product);
    }

    public Product updateProduct(Long id, Product product) {
        Product oldProduct = productRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("Product not found"));

        if (product.getName() != null) {
            oldProduct.setName(product.getName());
        }
        if (product.getPrice() != 0) {
            oldProduct.setPrice(product.getPrice());
        }
        if (product.getCategory() != null) {
            oldProduct.setCategory(product.getCategory());
        }
        return productRepository.save(oldProduct);
    }

    public void deleteProduct(Long id) {
        productRepository.deleteById(id);
    }
}
