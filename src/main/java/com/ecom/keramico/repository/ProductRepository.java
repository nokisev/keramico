package com.ecom.keramico.repository;

import com.ecom.keramico.model.Product;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface ProductRepository extends JpaRepository<Product, Long> {
    List<Product> findAllByName(String name);
    List<Product> findAllByCategory(String category);
}
