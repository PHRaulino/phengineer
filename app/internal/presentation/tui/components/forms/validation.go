package forms

import (
    "errors"
    "fmt"
    "net/mail"
    "net/url"
    "regexp"
    "strconv"
    "strings"
)

// Erros comuns
var (
    ErrFieldRequired = errors.New("campo obrigatório")
    ErrInvalidEmail  = errors.New("email inválido")
    ErrInvalidURL    = errors.New("URL inválida")
    ErrTooShort      = errors.New("valor muito curto")
    ErrTooLong       = errors.New("valor muito longo")
    ErrInvalidFormat = errors.New("formato inválido")
    ErrWeakPassword  = errors.New("senha muito fraca")
)

// ValidationFunc é uma função de validação
type ValidationFunc func(value string) error

// Validações pré-definidas
var (
    // Required valida campo obrigatório
    Required ValidationFunc = func(value string) error {
        if strings.TrimSpace(value) == "" {
            return ErrFieldRequired
        }
        return nil
    }
    
    // Email valida formato de email
    Email ValidationFunc = func(value string) error {
        _, err := mail.ParseAddress(value)
        if err != nil {
            return ErrInvalidEmail
        }
        return nil
    }
    
    // URL valida formato de URL
    ValidateURL ValidationFunc = func(value string) error {
        u, err := url.Parse(value)
        if err != nil || u.Scheme == "" || u.Host == "" {
            return ErrInvalidURL
        }
        return nil
    }
    
    // Numeric valida se é número
    Numeric ValidationFunc = func(value string) error {
        if _, err := strconv.ParseFloat(value, 64); err != nil {
            return errors.New("deve ser um número")
        }
        return nil
    }
    
    // Alpha valida se contém apenas letras
    Alpha ValidationFunc = func(value string) error {
        if !regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(value) {
            return errors.New("deve conter apenas letras")
        }
        return nil
    }
    
    // AlphaNumeric valida se é alfanumérico
    AlphaNumeric ValidationFunc = func(value string) error {
        if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(value) {
            return errors.New("deve conter apenas letras e números")
        }
        return nil
    }
)

// MinLength retorna validador de comprimento mínimo
func MinLength(min int) ValidationFunc {
    return func(value string) error {
        if len(value) < min {
            return fmt.Errorf("deve ter pelo menos %d caracteres", min)
        }
        return nil
    }
}

// MaxLength retorna validador de comprimento máximo
func MaxLength(max int) ValidationFunc {
    return func(value string) error {
        if len(value) > max {
            return fmt.Errorf("deve ter no máximo %d caracteres", max)
        }
        return nil
    }
}

// Range valida se está dentro de um intervalo
func Range(min, max int) ValidationFunc {
    return func(value string) error {
        num, err := strconv.Atoi(value)
        if err != nil {
            return errors.New("deve ser um número")
        }
        if num < min || num > max {
            return fmt.Errorf("deve estar entre %d e %d", min, max)
        }
        return nil
    }
}

// Pattern valida com regex
func Pattern(pattern string, message string) ValidationFunc {
    re := regexp.MustCompile(pattern)
    return func(value string) error {
        if !re.MatchString(value) {
            return errors.New(message)
        }
        return nil
    }
}

// StrongPassword valida senha forte
func StrongPassword() ValidationFunc {
    return func(value string) error {
        if len(value) < 8 {
            return errors.New("senha deve ter pelo menos 8 caracteres")
        }
        
        var hasUpper, hasLower, hasNumber, hasSpecial bool
        
        for _, char := range value {
            switch {
            case 'A' <= char && char <= 'Z':
                hasUpper = true
            case 'a' <= char && char <= 'z':
                hasLower = true
            case '0' <= char && char <= '9':
                hasNumber = true
            case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
                hasSpecial = true
            }
        }
        
        if !hasUpper {
            return errors.New("senha deve conter pelo menos uma letra maiúscula")
        }
        if !hasLower {
            return errors.New("senha deve conter pelo menos uma letra minúscula")
        }
        if !hasNumber {
            return errors.New("senha deve conter pelo menos um número")
        }
        if !hasSpecial {
            return errors.New("senha deve conter pelo menos um caractere especial")
        }
        
        return nil
    }
}

// Combine combina múltiplas validações
func Combine(validators ...ValidationFunc) ValidationFunc {
    return func(value string) error {
        for _, validator := range validators {
            if err := validator(value); err != nil {
                return err
            }
        }
        return nil
    }
}

// Custom permite criar validação customizada
func Custom(fn func(string) error) ValidationFunc {
    return ValidationFunc(fn)
}