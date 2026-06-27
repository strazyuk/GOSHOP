package services

import (
    "errors"
    "time"

    "go-ecommerce/database"
    "go-ecommerce/models"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

// AuthService handles registration, login, and token generation
type AuthService struct {
    jwtSecret string
    jwtExpiry int // hours
}

func NewAuthService(secret string, expiryHours int) *AuthService {
    return &AuthService{
        jwtSecret: secret,
        jwtExpiry: expiryHours,
    }
}

// RegisterInput is the data a client sends to create an account
type RegisterInput struct {
    Email     string `json:"email"      binding:"required,email"`    // must be a valid email
    Password  string `json:"password"   binding:"required,min=8"`    // at least 8 characters
    FirstName string `json:"first_name" binding:"required"`
    LastName  string `json:"last_name"  binding:"required"`
    Phone     string `json:"phone"`
}

// LoginInput is the data a client sends to log in
type LoginInput struct {
    Email    string `json:"email"    binding:"required,email"`
    Password string `json:"password" binding:"required"`
}
type AuthResponse struct {
    Token string              `json:"token"`
    User  models.UserResponse `json:"user"`
}


func (s *AuthService) Register(input RegisterInput) (*AuthResponse, error) {
    
    var existing models.User
    result := database.DB.Where("email = ?", input.Email).First(&existing)
    if result.Error == nil{
        return nil, errors.New("email already registered")
    }

    

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, errors.New("failed to hash password")
    }

    // Create the user record
    user := models.User{
        Email:     input.Email,
        Password:  string(hashedPassword),
        FirstName: input.FirstName,
        LastName:  input.LastName,
        Phone:     input.Phone,
        Role:      models.RoleCustomer,
    }

    if err := database.DB.Create(&user).Error; err != nil {
        return nil, errors.New("failed to create user")
    }

    cart := models.Cart{UserID: user.ID}
    database.DB.Create(&cart)
    token, err := s.generateToken(user.ID, string(user.Role))
    if err != nil {
        return nil, err
    }

    return &AuthResponse{
        Token: token,
        User:  user.ToResponse(),
    }, nil
}


func (s *AuthService) Login(input LoginInput) (*AuthResponse, error) {
    var user models.User
    if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
        return nil, errors.New("invalid email or password")
    }

    
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        return nil, errors.New("invalid email or password")
    }

    token, err := s.generateToken(user.ID, string(user.Role))
    if err != nil {
        return nil, err
    }

    return &AuthResponse{
        Token: token,
        User:  user.ToResponse(),
    }, nil
}

// Claims is the data we embed inside the JWT token
type Claims struct {
    UserID uint   `json:"user_id"`
    Role   string `json:"role"`
    // jwt.RegisteredClaims adds standard fields like ExpiresAt, IssuedAt, etc.
    jwt.RegisteredClaims
}

func (s *AuthService) generateToken(userID uint, role string) (string, error) {
    claims := Claims{
        UserID: userID,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.jwtExpiry) * time.Hour)),
       
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    // jwt.NewWithClaims creates a token with HS256 signing method
    // HS256 = HMAC-SHA256 — a symmetric algorithm (same key to sign and verify)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // SignedString signs the token with our secret key and returns the string
    return token.SignedString([]byte(s.jwtSecret))
}

// ValidateToken parses and validates a JWT string, returning the claims
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
       
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(s.jwtSecret), nil
    })

    if err != nil || !token.Valid {
        return nil, errors.New("invalid or expired token")
    }

    return claims, nil
}
