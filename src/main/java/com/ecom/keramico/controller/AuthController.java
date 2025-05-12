package com.ecom.keramico.controller;

import com.ecom.keramico.model.User;
import com.ecom.keramico.model.UserRole;
import com.ecom.keramico.repository.UserRepository;
import org.springframework.http.ResponseEntity;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/users")
public class AuthController {

//    TODO
    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;

    public AuthController(UserRepository userRepository, PasswordEncoder passwordEncoder) {
        this.userRepository = userRepository;
        this.passwordEncoder = passwordEncoder;
    }

    @PostMapping("/register")
    public ResponseEntity<Object> registerUser(@RequestBody User regUser) {
        User user = new User();
        user.setUsername(regUser.getUsername());
        user.setHashed_password(passwordEncoder.encode(regUser.getPassword()));
        user.setRole(UserRole.USER);

        userRepository.save(user);
        return ResponseEntity.ok("User successful register");
    }
}
