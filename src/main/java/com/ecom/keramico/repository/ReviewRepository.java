package com.ecom.keramico.repository;

import com.ecom.keramico.model.Product;
import com.ecom.keramico.model.Review;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface ReviewRepository extends JpaRepository<Review, Long> {
    List<Review> findReviewByProductId(Product product);
}
