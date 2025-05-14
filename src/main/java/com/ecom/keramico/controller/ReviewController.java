package com.ecom.keramico.controller;

import com.ecom.keramico.model.Product;
import com.ecom.keramico.model.Review;
import com.ecom.keramico.repository.ProductRepository;
import com.ecom.keramico.repository.ReviewRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/reviews")
public class ReviewController {

    @Autowired
    private ReviewRepository reviewRepository;
    @Autowired
    private ProductRepository productRepository;

    @GetMapping("/product/{product_id}")
    public List<Review> showReviewOfProduct(@PathVariable Long product_id) {
        Product product = productRepository.findById(product_id)
                .orElseThrow(() -> new RuntimeException("Product or reviews not found"));
        List<Review> reviews = product.getReviews();
        return reviews;
    }

    @PostMapping("/{product_id}")
    public String addReviewToProduct(@PathVariable Long product_id, @RequestBody Review review) {
        Product product = productRepository.findById(product_id)
                .orElseThrow(() -> new RuntimeException("Product not found"));
        product.getReviews().add(review);
        return "Review added";
    }
}
